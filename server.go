// Файл server.go содержит обработчики HTTP запросов сервера: просмотр, добавление, удаление, изменение таблиц, генерация отчета.
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/akionka/aviasales/internal/store"
	"github.com/akionka/aviasales/internal/store/mysqlstore"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	tokenSecret                  = []byte("thisisthemostsecrettokenintheentireuniverse")
	errBadAuthorizationToken     = errors.New("некорректный токен авторизации")
	errIncorrectLoginOrPassword  = errors.New("неправильный логин или пароль")
	errRequestedItemDoesNotExist = errors.New("запрошенная сущность не существует")
	errRegularUserPutDelete      = errors.New("вы не можете удалять или обновлять данные")
)

const (
	ctxKeyCashier ctxKey = iota
	ctxKeyPagination
)

type paginationInfo struct {
	rowCount int
	page     int
}

type ctxKey uint8

type server struct {
	router *mux.Router
	store  store.Store
}

func newServer(store store.Store) *server {
	s := &server{
		store: store,
	}
	s.configureRouter()
	return s
}

func (s *server) start() error {
	return http.ListenAndServe(":8080", s.router)
}

func (s *server) configureRouter() {
	s.router = mux.NewRouter().PathPrefix("/api").Subrouter().StrictSlash(true)

	s.router.Use(handlers.CORS(
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodPost, http.MethodOptions}),
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5500", "http://localhost:3000"})))

	s.router.HandleFunc("/user", s.handleCashiersCreate()).Methods(http.MethodPost, http.MethodOptions)
	s.router.HandleFunc("/session", s.handleSessionsCreate()).Methods(http.MethodPost, http.MethodOptions)

	s.router.HandleFunc("/timezones", func(w http.ResponseWriter, r *http.Request) {
		timezones, err := s.store.Timezone().FindAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, timezones)
	}).Methods(http.MethodGet, http.MethodOptions)

	secured := s.router.NewRoute().Subrouter()
	secured.Use(s.authenticateUser)

	secured.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		c, ok := r.Context().Value(ctxKeyCashier).(*store.CashierModel)
		if ok {
			cashierResponse := &Cashier{
				ID:         c.ID,
				Login:      c.Login,
				LastName:   c.LastName,
				FirstName:  c.FirstName,
				MiddleName: c.MiddleName,
				RoleID:     c.RoleID,
			}
			s.respond(w, r, 200, cashierResponse)
			return
		}
		s.error(w, r, http.StatusInternalServerError, nil)
	}).Methods(http.MethodGet, http.MethodOptions)

	securedGet := secured.Methods(http.MethodGet, http.MethodOptions).Subrouter()
	securedGet.Use(s.paginateMiddleware)

	securedGet.HandleFunc("/airports", s.handleAirportsGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/booking_offices", s.handleBookingOfficesGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/cashiers", s.handleCashiersGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/flight_in_tickets", s.handleFlightInTicketsGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/flights", s.handleFlightsGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/lines", s.handleLinesGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/liner_models", s.handleLinerModelsGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/liners", s.handleLinersGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/purchases", s.handlePurchasesGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/seats", s.handleSeatsGet()).Methods(http.MethodGet, http.MethodOptions)
	securedGet.HandleFunc("/tickets", s.handleTicketsGet()).Methods(http.MethodGet, http.MethodOptions)

	securedGet.HandleFunc("/airports", s.handleAirportsCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/booking_offices", s.handleBookingOfficesCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/cashiers", s.handleCashiersCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/flight_in_tickets", s.handleFlightInTicketsCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/flights", s.handleFlightsCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/lines", s.handleLinesCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/liner_models", s.handleLinerModelsCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/liners", s.handleLinersCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/purchases", s.handlePurchasesCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/seats", s.handleSeatsCreate()).Methods(http.MethodPost, http.MethodOptions)
	securedGet.HandleFunc("/tickets", s.handleTicketsCreate()).Methods(http.MethodPost, http.MethodOptions)

	adminOnlyUpdateDelete := secured.NewRoute().Subrouter()
	adminOnlyUpdateDelete.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, ok := r.Context().Value(ctxKeyCashier).(*store.CashierModel)
			if !ok {
				h.ServeHTTP(w, r)
				return
			}
			if (r.Method == "PUT" || r.Method == "DELETE") && c.RoleID != 2 {
				s.error(w, r, http.StatusUnauthorized, errRegularUserPutDelete)
				return
			}
			h.ServeHTTP(w, r)
		})
	})

	adminOnlyUpdateDelete.HandleFunc("/airports/{code}", s.handleAirportGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/booking_offices/{id:[0-9]+}", s.handleBookingOfficeGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/cashiers/{id:[0-9]+}", s.handleCashierGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/cashiers/{id:[0-9]+}/password", s.handleCashierPasswordUpdate()).Methods(http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/flight_in_tickets/{id:[0-9]+}", s.handleFlightInTicketGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/flights/{id:[0-9]+}", s.handleFlightGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/lines/{code}", s.handleLineGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/liner_models/{code}", s.handleLinerModelGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/liners/{code}", s.handleLinerGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/purchases/{id:[0-9]+}", s.handlePurchaseGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/seats/{id:[0-9]+}", s.handleSeatGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/tickets/{id:[0-9]+}", s.handleTicketGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	adminOnlyUpdateDelete.HandleFunc("/tickets/{id:[0-9]+}/report", s.handleTicketReportGet()).Methods(http.MethodGet, http.MethodOptions)
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}
		if headerParts[0] != "Bearer" {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}

		token, err := jwt.Parse(headerParts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errBadAuthorizationToken
			}
			return tokenSecret, nil
		})
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !(ok && token.Valid) {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}

		login, ok := claims["login"].(string)
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}
		cashier, err := s.store.Cashier().FindByLogin(login)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyCashier, cashier)))
	})
}

func (s *server) paginateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 0
		}
		count, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil {
			count = 10
		}
		if count == -1 {
			count = math.MaxInt64
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyPagination, paginationInfo{
			page:     page,
			rowCount: count,
		})))
	})
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	type response struct {
		Token string   `json:"token"`
		User  *Cashier `json:"user"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		c, err := s.store.Cashier().FindByLogin(req.Login)
		if err != nil || !c.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectLoginOrPassword)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"login": c.Login,
			"role":  c.RoleID,
		})
		tokenString, err := token.SignedString(tokenSecret)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, 200, response{
			Token: tokenString,
			User: &Cashier{
				ID:         c.ID,
				Login:      c.Login,
				LastName:   c.LastName,
				FirstName:  c.FirstName,
				MiddleName: c.MiddleName,
				RoleID:     c.RoleID,
			},
		})
	}
}

func (s *server) handleAirportsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		airports, err := s.store.Airport().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Airport().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := AirportList{
			Items:      make([]Airport, len(*airports)),
			TotalCount: totalCount,
		}

		for i, v := range *airports {
			response.Items[i] = Airport{
				IATACode: v.IATACode,
				City:     v.City,
				Timezone: v.Timezone,
			}
		}

		s.respond(w, r, 200, response)
	}
}

func (s *server) handleAirportGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if _, err := s.store.Airport().Find(vars["code"]); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			a, err := s.store.Airport().Find(vars["code"])
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Airport{
				IATACode: a.IATACode,
				City:     a.City,
				Timezone: a.Timezone,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.Airport().Delete(vars["code"])
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			a := &Airport{}
			if err := json.NewDecoder(r.Body).Decode(a); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := a.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Airport().Update(vars["code"], &store.AirportModel{
				IATACode: a.IATACode,
				City:     a.City,
				Timezone: a.Timezone,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, a)
		}
	}
}

func (s *server) handleBookingOfficesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		offices, err := s.store.BookingOffice().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.BookingOffice().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := BookingOfficeList{
			Items:      make([]BookingOffice, len(*offices)),
			TotalCount: totalCount,
		}

		for i, v := range *offices {
			response.Items[i] = BookingOffice{
				ID:          v.ID,
				Address:     v.Address,
				PhoneNumber: v.PhoneNumber,
			}
		}

		s.respond(w, r, 200, response)
	}
}

func (s *server) handleBookingOfficeGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.BookingOffice().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			o, err := s.store.BookingOffice().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &BookingOffice{
				ID:          o.ID,
				Address:     o.Address,
				PhoneNumber: o.PhoneNumber,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err = s.store.BookingOffice().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			b := &BookingOffice{}
			if err := json.NewDecoder(r.Body).Decode(b); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := b.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.BookingOffice().Update(id, &store.BookingOfficeModel{
				ID:          b.ID,
				Address:     b.Address,
				PhoneNumber: b.PhoneNumber,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, b)
		}
	}
}

func (s *server) handleCashiersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		cashiers, err := s.store.Cashier().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Cashier().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := CashierList{
			Items:      make([]Cashier, len(*cashiers)),
			TotalCount: totalCount,
		}

		for i, v := range *cashiers {
			response.Items[i] = Cashier{
				ID:         v.ID,
				Login:      v.Login,
				LastName:   v.LastName,
				FirstName:  v.FirstName,
				MiddleName: v.MiddleName,
			}
		}

		s.respond(w, r, 200, response)
	}
}

func (s *server) handleCashierGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Cashier().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			c, err := s.store.Cashier().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Cashier{
				ID:         c.ID,
				Login:      c.Login,
				LastName:   c.LastName,
				FirstName:  c.FirstName,
				MiddleName: c.MiddleName,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.Cashier().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			c := &Cashier{}
			if err := json.NewDecoder(r.Body).Decode(c); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := c.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Cashier().Update(id, &store.CashierModel{
				ID:         c.ID,
				Login:      c.Login,
				LastName:   c.LastName,
				FirstName:  c.FirstName,
				MiddleName: c.MiddleName,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, c)
		}

	}
}

func (s *server) handleCashierPasswordUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Cashier().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		c := &store.CashierModel{}
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := validation.Validate(c.Password, validation.Required, validation.Length(6, 72)); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		cModel := store.CashierModel{
			ID: id,
		}

		cModel.SetPassword(c.Password)
		if err := s.store.Cashier().UpdatePassword(&cModel); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) handleFlightInTicketsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		flightInTickets, err := s.store.FlightInTicket().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.FlightInTicket().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := FlightInTicketList{
			Items:      make([]FlightInTicket, len(*flightInTickets)),
			TotalCount: totalCount,
		}

		for i, v := range *flightInTickets {
			response.Items[i] = FlightInTicket{
				ID:       v.ID,
				FlightID: v.FlightID,
				SeatID:   v.SeatID,
				TicketID: v.TicketID,
			}
		}
		s.respond(w, r, 200, response)
	}
}

func (s *server) handleFlightInTicketGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.FlightInTicket().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			f, err := s.store.FlightInTicket().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &FlightInTicket{
				ID:       id,
				FlightID: f.FlightID,
				SeatID:   f.SeatID,
				TicketID: f.SeatID,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err = s.store.FlightInTicket().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			f := &FlightInTicket{}
			if err := json.NewDecoder(r.Body).Decode(f); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := f.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.FlightInTicket().Update(id, &store.FlightInTicketModel{
				FlightID: f.FlightID,
				SeatID:   f.SeatID,
				TicketID: f.TicketID,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, f)
		}
	}
}

func (s *server) handleFlightsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		flights, err := s.store.Flight().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Flight().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := FlightList{
			Items:      make([]Flight, len(*flights)),
			TotalCount: totalCount,
		}

		for i, v := range *flights {
			response.Items[i] = Flight{
				ID:        v.ID,
				DepDate:   v.DepDate,
				LineCode:  v.LineCode,
				LinerCode: v.LinerCode,
				IsHot:     v.IsHot,
			}
		}
		s.respond(w, r, 200, response)
	}
}

func (s *server) handleFlightGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Flight().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			f, err := s.store.Flight().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Flight{
				ID:        f.ID,
				DepDate:   f.DepDate,
				LineCode:  f.LineCode,
				IsHot:     f.IsHot,
				LinerCode: f.LinerCode,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.Flight().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			f := &Flight{}
			if err := json.NewDecoder(r.Body).Decode(f); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := f.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Flight().Update(id, &store.FlightModel{
				DepDate:   f.DepDate,
				LineCode:  f.LineCode,
				IsHot:     f.IsHot,
				LinerCode: f.LinerCode,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, f)
		}
	}
}

func (s *server) handleLinesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		lines, err := s.store.Line().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Line().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := LineList{
			Items:      make([]Line, len(*lines)),
			TotalCount: totalCount,
		}

		for i, v := range *lines {
			response.Items[i] = Line{
				LineCode:   v.LineCode,
				DepTime:    v.DepTime,
				ArrTime:    v.ArrTime,
				BasePrice:  v.BasePrice,
				DepAirport: v.DepAirport,
				ArrAirport: v.ArrAirport,
			}
		}
		s.respond(w, r, 200, response)
	}
}

func (s *server) handleLineGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if _, err := s.store.Line().Find(vars["code"]); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			l, err := s.store.Line().Find(vars["code"])
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Line{
				LineCode:   l.LineCode,
				DepTime:    l.DepTime,
				ArrTime:    l.ArrTime,
				BasePrice:  l.BasePrice,
				DepAirport: l.DepAirport,
				ArrAirport: l.ArrAirport,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.Line().Delete(vars["code"])
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			l := &Line{}
			if err := json.NewDecoder(r.Body).Decode(l); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := l.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Line().Update(vars["code"], &store.LineModel{
				LineCode:   l.LineCode,
				DepTime:    l.DepTime,
				ArrTime:    l.ArrTime,
				BasePrice:  l.BasePrice,
				DepAirport: l.DepAirport,
				ArrAirport: l.ArrAirport,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, l)
		}
	}
}

func (s *server) handleLinerModelsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		models, err := s.store.LinerModel().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.LinerModel().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := LinerModelList{
			Items:      make([]LinerModel, len(*models)),
			TotalCount: totalCount,
		}

		for i, v := range *models {
			response.Items[i] = LinerModel{
				IATATypeCode: v.IATATypeCode,
				Name:         v.Name,
			}
		}
		s.respond(w, r, 200, response)
	}
}

func (s *server) handleLinerModelGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if _, err := s.store.LinerModel().Find(vars["code"]); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			m, err := s.store.LinerModel().Find(vars["code"])
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &LinerModel{
				IATATypeCode: m.IATATypeCode,
				Name:         m.Name,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.LinerModel().Delete(vars["code"])
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			m := &LinerModel{}
			if err := json.NewDecoder(r.Body).Decode(m); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := m.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.LinerModel().Update(vars["code"], &store.LinerModelModel{
				IATATypeCode: m.IATATypeCode,
				Name:         m.Name,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, m)
		}
	}
}

func (s *server) handleLinersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		liners, err := s.store.Liner().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Liner().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := LinerList{
			Items:      make([]Liner, len(*liners)),
			TotalCount: totalCount,
		}

		for i, v := range *liners {
			response.Items[i] = Liner{
				IATACode:  v.IATACode,
				ModelCode: v.ModelCode,
			}
		}
		s.respond(w, r, 200, response)
	}
}

func (s *server) handleLinerGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if _, err := s.store.Liner().Find(vars["code"]); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			l, err := s.store.Liner().Find(vars["code"])
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Liner{
				IATACode:  l.IATACode,
				ModelCode: l.ModelCode,
			})
		}

		if r.Method == http.MethodDelete {
			err := s.store.Liner().Delete(vars["code"])
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
		}

		if r.Method == http.MethodPut {
			l := &Liner{}
			if err := json.NewDecoder(r.Body).Decode(l); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := l.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Liner().Update(vars["code"], &store.LinerModel{
				IATACode:  l.IATACode,
				ModelCode: l.ModelCode,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, l)
		}
	}
}

func (s *server) handlePurchasesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		purchases, err := s.store.Purchase().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Purchase().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := PurchaseList{
			Items:      make([]Purchase, len(*purchases)),
			TotalCount: totalCount,
		}

		for i, v := range *purchases {
			response.Items[i] = Purchase{
				ID:              v.ID,
				Date:            v.Date,
				BookingOfficeID: v.BookingOfficeID,
				TotalPrice:      v.TotalPrice,
				ContactPhone:    v.ContactPhone,
				ContactEmail:    v.ContactEmail,
				CashierID:       v.CashierID,
			}
		}
		s.respond(w, r, 200, response)
	}

}

func (s *server) handlePurchaseGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Purchase().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			p, err := s.store.Purchase().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Purchase{
				ID:              p.ID,
				Date:            p.Date,
				BookingOfficeID: p.BookingOfficeID,
				TotalPrice:      p.TotalPrice,
				ContactPhone:    p.ContactPhone,
				ContactEmail:    p.ContactEmail,
				CashierID:       p.CashierID,
			})
		}

		if r.Method == http.MethodDelete {
			err = s.store.Purchase().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
		}

		if r.Method == http.MethodPut {
			p := &Purchase{}
			if err := json.NewDecoder(r.Body).Decode(p); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := p.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Purchase().Update(id, &store.PurchaseModel{
				ID:              p.ID,
				Date:            p.Date,
				BookingOfficeID: p.BookingOfficeID,
				TotalPrice:      p.TotalPrice,
				ContactPhone:    p.ContactPhone,
				ContactEmail:    p.ContactEmail,
				CashierID:       p.CashierID,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, p)
		}
	}
}

func (s *server) handleSeatsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		seats, err := s.store.Seat().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Seat().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := SeatList{
			Items:      make([]Seat, len(*seats)),
			TotalCount: totalCount,
		}

		for i, v := range *seats {
			response.Items[i] = Seat{
				ID:             v.ID,
				Number:         v.Number,
				LinerModelCode: v.LinerModelCode,
				Class:          v.Class,
			}
		}
		s.respond(w, r, 200, response)

	}

}

func (s *server) handleSeatGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Seat().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			p, err := s.store.Seat().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Seat{
				ID:             p.ID,
				Number:         p.Number,
				Class:          p.Class,
				LinerModelCode: p.LinerModelCode,
			})
		}

		if r.Method == http.MethodDelete {
			err = s.store.Seat().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			seat := &Seat{}
			if err := json.NewDecoder(r.Body).Decode(seat); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := seat.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Seat().Update(id, &store.SeatModel{
				ID:             seat.ID,
				Number:         seat.Number,
				Class:          seat.Class,
				LinerModelCode: seat.LinerModelCode,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, seat)
		}
	}
}

func (s *server) handleTicketsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		ticket, err := s.store.Ticket().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		totalCount, err := s.store.Ticket().TotalCount()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		response := TicketList{
			Items:      make([]Ticket, len(*ticket)),
			TotalCount: totalCount,
		}

		for i, v := range *ticket {
			response.Items[i] = Ticket{
				ID:                      v.ID,
				PassengerLastName:       v.PassengerLastName,
				PassengerGivenName:      v.PassengerGivenName,
				PassengerBirthDate:      v.PassengerBirthDate,
				PassengerPassportNumber: v.PassengerPassportNumber,
				PassengerSex:            v.PassengerSex,
				PurchaseID:              v.PurchaseID,
			}
		}
		s.respond(w, r, 200, response)
	}
}

func (s *server) handleTicketGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if _, err := s.store.Ticket().Find(id); err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if r.Method == http.MethodGet {
			p, err := s.store.Ticket().Find(id)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Ticket{
				ID:                      p.ID,
				PassengerLastName:       p.PassengerLastName,
				PassengerGivenName:      p.PassengerGivenName,
				PassengerBirthDate:      p.PassengerBirthDate,
				PassengerPassportNumber: p.PassengerPassportNumber,
				PassengerSex:            p.PassengerSex,
				PurchaseID:              p.PurchaseID,
			})
		}

		if r.Method == http.MethodDelete {
			err = s.store.Ticket().Delete(id)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
			t := &Ticket{}
			if err := json.NewDecoder(r.Body).Decode(t); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			if err := t.Validate(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			if err := s.store.Ticket().Update(id, &store.TicketModel{
				ID:                      t.ID,
				PassengerLastName:       t.PassengerLastName,
				PassengerGivenName:      t.PassengerGivenName,
				PassengerBirthDate:      t.PassengerBirthDate,
				PassengerPassportNumber: t.PassengerPassportNumber,
				PassengerSex:            t.PassengerSex,
				PurchaseID:              t.PurchaseID,
			}); err != nil {
				if err == mysqlstore.ErrNoChanges {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, t)
		}
	}
}

func (s *server) handleAirportsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := &Airport{}
		if err := json.NewDecoder(r.Body).Decode(a); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := a.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Airport().Create(&store.AirportModel{
			IATACode: a.IATACode,
			City:     a.City,
			Timezone: a.Timezone,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, a)
	}
}

func (s *server) handleBookingOfficesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		o := &BookingOffice{}
		if err := json.NewDecoder(r.Body).Decode(o); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := o.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.BookingOffice().Create(&store.BookingOfficeModel{
			ID:          o.ID,
			Address:     o.Address,
			PhoneNumber: o.PhoneNumber,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, o)
	}
}

func (s *server) handleCashiersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &Cashier{}
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := c.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		cModel := &store.CashierModel{
			Login:      c.Login,
			LastName:   c.LastName,
			FirstName:  c.FirstName,
			MiddleName: c.MiddleName,
		}
		if err := cModel.SetPassword(c.Password); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := s.store.Cashier().Create(cModel); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		c.Password = ""
		s.respond(w, r, http.StatusOK, c)
	}
}

func (s *server) handleFlightInTicketsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := &FlightInTicket{}
		if err := json.NewDecoder(r.Body).Decode(f); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := f.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.FlightInTicket().Create(&store.FlightInTicketModel{
			FlightID: f.FlightID,
			SeatID:   f.SeatID,
			TicketID: f.TicketID,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, f)
	}
}

func (s *server) handleFlightsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f := &Flight{}
		if err := json.NewDecoder(r.Body).Decode(f); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := f.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Flight().Create(&store.FlightModel{
			DepDate:   f.DepDate,
			LineCode:  f.LineCode,
			IsHot:     f.IsHot,
			LinerCode: f.LinerCode,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, f)
	}
}

func (s *server) handleLinesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &Line{}
		if err := json.NewDecoder(r.Body).Decode(l); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := l.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Line().Create(&store.LineModel{
			LineCode:   l.LineCode,
			DepTime:    l.DepTime,
			ArrTime:    l.ArrTime,
			BasePrice:  l.BasePrice,
			DepAirport: l.DepAirport,
			ArrAirport: l.ArrAirport,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, l)
	}
}

func (s *server) handleLinerModelsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := &LinerModel{}
		if err := json.NewDecoder(r.Body).Decode(m); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := m.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.LinerModel().Create(&store.LinerModelModel{
			IATATypeCode: m.IATATypeCode,
			Name:         m.Name,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, m)
	}
}

func (s *server) handleLinersCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &Liner{}
		if err := json.NewDecoder(r.Body).Decode(l); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := l.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Liner().Create(&store.LinerModel{
			IATACode:  l.IATACode,
			ModelCode: l.ModelCode,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, l)
	}
}

func (s *server) handlePurchasesCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &Purchase{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := p.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Purchase().Create(&store.PurchaseModel{
			ID:              p.ID,
			Date:            p.Date,
			BookingOfficeID: p.BookingOfficeID,
			TotalPrice:      p.TotalPrice,
			ContactPhone:    p.ContactPhone,
			ContactEmail:    p.ContactEmail,
			CashierID:       p.CashierID,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, p)
	}
}

func (s *server) handleSeatsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		seat := &Seat{}
		if err := json.NewDecoder(r.Body).Decode(seat); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := seat.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Seat().Create(&store.SeatModel{
			ID:             seat.ID,
			Number:         seat.Number,
			Class:          seat.Class,
			LinerModelCode: seat.LinerModelCode,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, seat)
	}
}

func (s *server) handleTicketsCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := &Ticket{}
		if err := json.NewDecoder(r.Body).Decode(t); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := t.Validate(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := s.store.Ticket().Create(&store.TicketModel{
			ID:                      t.ID,
			PassengerLastName:       t.PassengerLastName,
			PassengerGivenName:      t.PassengerGivenName,
			PassengerBirthDate:      t.PassengerBirthDate,
			PassengerPassportNumber: t.PassengerPassportNumber,
			PassengerSex:            t.PassengerSex,
			PurchaseID:              t.PurchaseID,
		}); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, t)
	}
}

func (s *server) handleTicketReportGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		t, err := s.store.Ticket().Find(id)
		if err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errRequestedItemDoesNotExist)
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		flights, office, cashier, purchase, totalTime, err := s.store.Ticket().Report(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		report := &TicketReport{
			Ticket: Ticket{
				ID:                      t.ID,
				PassengerLastName:       t.PassengerLastName,
				PassengerGivenName:      t.PassengerGivenName,
				PassengerBirthDate:      t.PassengerBirthDate,
				PassengerPassportNumber: "******" + t.PassengerPassportNumber[6:10],
				PassengerSex:            t.PassengerSex,
			},
			BookingOffice: BookingOffice{
				ID:          office.ID,
				Address:     office.Address,
				PhoneNumber: office.PhoneNumber,
			},
			Cashier: Cashier{
				ID:         cashier.ID,
				FirstName:  cashier.FirstName,
				LastName:   cashier.LastName,
				MiddleName: cashier.MiddleName,
				Login:      cashier.Login,
			},
			Purchase: Purchase{
				ID:              purchase.ID,
				Date:            purchase.Date,
				TotalPrice:      purchase.TotalPrice,
				ContactPhone:    purchase.ContactPhone,
				ContactEmail:    purchase.ContactEmail,
				BookingOfficeID: purchase.BookingOfficeID,
				CashierID:       purchase.CashierID,
			},
			Flights:   flights,
			TotalTime: int(totalTime.Seconds()),
		}
		s.respond(w, r, http.StatusOK, report)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	if err, ok := err.(validation.Errors); ok {
		s.respond(w, r, code, map[string]validation.Errors{"error": err})
		return
	}
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
