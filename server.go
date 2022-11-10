package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/akionka/aviasales/internal/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	tokenSecret                  = []byte("thisisthemostsecrettokenintheentireuniverse")
	errBadAuthorizationToken     = errors.New("bad authorization token")
	errIncorrectLoginOrPassword  = errors.New("incorrect login or password")
	errRequestedItemDoesNotExist = errors.New("requested item does not exist")
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
		handlers.AllowedOrigins([]string{"http://127.0.0.1:5500"})))

	s.router.HandleFunc("/session", s.handleSessionsCreate()).Methods("POST", "OPTIONS")

	secured := s.router.NewRoute().Subrouter()
	secured.Use(s.authenticateUser)

	secured.HandleFunc("/whoami", func(w http.ResponseWriter, r *http.Request) {
		c, ok := r.Context().Value(ctxKeyCashier).(*store.CashierModel)
		if ok {
			cashierResponse := &Cashier{
				Login:      c.Login,
				LastName:   c.LastName,
				FirstName:  c.FirstName,
				MiddleName: c.MiddleName,
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

	// secured.HandleFunc("/booking_offices", s.handleBookingOfficesCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/airports", s.handleAirportsCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/cashiers", s.handleCashiersCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/flight_in_tickets", s.handleFlightInTicketsCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/flights", s.handleFlightsCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/lines", s.handleLinesCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/liner_models", s.handleLinerModelsCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/liners", s.handleLinersCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/purchases", s.handlePurchasesCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/seats", s.handleSeatsCreate()).Methods(http.MethodPost, http.MethodOptions)
	// secured.HandleFunc("/tickets", s.handleTicketsCreate()).Methods(http.MethodPost, http.MethodOptions)

	secured.HandleFunc("/airports/{code}", s.handleAirportGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/booking_offices/{id:[0-9]+}", s.handleBookingOfficeGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/cashiers/{login}", s.handleCashierGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/flight_in_tickets/{dep_date}/{line_code}/{seat_id:[0-9]+}/{ticket_no:[0-9]+}", s.handleFlightInTicketGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/flights/{dep_date}/{line_code}", s.handleFlightGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/lines/{code}", s.handleLineGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/liner_models/{code}", s.handleLinerModelGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/liners/{code}", s.handleLinerGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/purchases/{id:[0-9]+}", s.handlePurchaseGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/seats/{id:[0-9]+}", s.handleSeatGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
	secured.HandleFunc("/tickets/{id:[0-9]+}", s.handleTicketGetDeleteUpdate()).Methods(http.MethodGet, http.MethodDelete, http.MethodPut, http.MethodOptions)
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

		cashierLogin, ok := claims["login"].(string)
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errBadAuthorizationToken)
			return
		}
		cashier, err := s.store.Cashier().Find(cashierLogin)
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
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		c, err := s.store.Cashier().Find(req.Login)
		if err != nil || !c.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectLoginOrPassword)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"login": c.Login,
		})
		tokenString, err := token.SignedString(tokenSecret)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, 200, map[string]string{
			"token": tokenString,
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

		airPortResponse := make([]Airport, len(*airports))
		for i, v := range *airports {
			airPortResponse[i] = Airport{
				IATACode: v.IATACode,
				City:     v.City,
				Timezone: v.Timezone,
			}
		}

		s.respond(w, r, 200, airPortResponse)
	}
}

func (s *server) handleAirportGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

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

		officeResponse := make([]BookingOffice, len(*offices))
		for i, v := range *offices {
			officeResponse[i] = BookingOffice{
				ID:          v.ID,
				Address:     v.Address,
				PhoneNumber: v.PhoneNumber,
			}
		}

		s.respond(w, r, 200, officeResponse)
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

		cashierResponse := make([]Cashier, len(*cashiers))
		for i, v := range *cashiers {
			cashierResponse[i] = Cashier{
				Login:      v.Login,
				LastName:   v.LastName,
				FirstName:  v.FirstName,
				MiddleName: v.MiddleName,
			}
		}
		s.respond(w, r, 200, cashierResponse)
	}
}

func (s *server) handleCashierGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			c, err := s.store.Cashier().Find(vars["login"])
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Cashier{
				Login:      c.Login,
				LastName:   c.LastName,
				FirstName:  c.FirstName,
				MiddleName: c.MiddleName,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.Cashier().Delete(vars["login"])
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {
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

		flightInTicketsResponse := make([]FlightInTicket, len(*flightInTickets))
		for i, v := range *flightInTickets {
			flightInTicketsResponse[i] = FlightInTicket{
				DepDate:  v.DepDate,
				LineCode: v.LineCode,
				SeatID:   v.SeatID,
				TicketNo: v.TicketNo,
			}
		}
		s.respond(w, r, 200, flightInTicketsResponse)
	}
}

func (s *server) handleFlightInTicketGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["seat_id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		ticketNo, err := strconv.ParseInt(vars["ticket_no"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if r.Method == http.MethodGet {
			f, err := s.store.FlightInTicket().Find(vars["dep_date"], vars["line_code"], id, ticketNo)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &FlightInTicket{
				DepDate:  f.DepDate,
				LineCode: f.LineCode,
				SeatID:   f.SeatID,
				TicketNo: f.TicketNo,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err = s.store.FlightInTicket().Delete(vars["dep_date"], vars["line_code"], id, ticketNo)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {

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

		flightsResponse := make([]Flight, len(*flights))
		for i, v := range *flights {
			flightsResponse[i] = Flight{
				DepDate:   v.DepDate,
				LineCode:  v.LineCode,
				LinerCode: v.LinerCode,
				IsHot:     v.IsHot,
			}
		}
		s.respond(w, r, 200, flightsResponse)
	}
}

func (s *server) handleFlightGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		if r.Method == http.MethodGet {
			f, err := s.store.Flight().Find(vars["dep_date"], vars["line_code"])
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Flight{
				DepDate:   f.DepDate,
				LineCode:  f.LineCode,
				IsHot:     f.IsHot,
				LinerCode: f.LinerCode,
			})
			return
		}

		if r.Method == http.MethodDelete {
			err := s.store.Flight().Delete(vars["dep_date"], vars["line_code"])
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {

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

		linesResponse := make([]Line, len(*lines))
		for i, v := range *lines {
			linesResponse[i] = Line{
				LineCode:   v.LineCode,
				DepTime:    v.DepTime,
				ArrTime:    v.ArrTime,
				BasePrice:  v.BasePrice,
				DepAirport: v.DepAirport,
				ArrAirport: v.ArrAirport,
			}
		}
		s.respond(w, r, 200, linesResponse)
	}
}

func (s *server) handleLineGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

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

		modelsResponse := make([]LinerModel, len(*models))
		for i, v := range *models {
			modelsResponse[i] = LinerModel{
				IATATypeCode: v.IATATypeCode,
				Name:         v.Name,
			}
		}
		s.respond(w, r, 200, modelsResponse)
	}
}

func (s *server) handleLinerModelGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

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

		linersResponse := make([]Liner, len(*liners))
		for i, v := range *liners {
			linersResponse[i] = Liner{
				IATACode:  v.IATACode,
				ModelCode: v.ModelCode,
			}
		}
		s.respond(w, r, 200, linersResponse)
	}
}

func (s *server) handleLinerGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

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

		purchasesResponse := make([]Purchase, len(*purchases))
		for i, v := range *purchases {
			purchasesResponse[i] = Purchase{
				ID:              v.ID,
				Date:            v.Date,
				BookingOfficeID: v.BookingOfficeID,
				TotalPrice:      v.TotalPrice,
				ContactPhone:    v.ContactPhone,
				ContactEmail:    v.ContactEmail,
				CashierLogin:    v.CashierLogin,
			}
		}
		s.respond(w, r, 200, purchasesResponse)
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
				CashierLogin:    p.CashierLogin,
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

		seatsResponse := make([]Seat, len(*seats))
		for i, v := range *seats {
			seatsResponse[i] = Seat{
				ID:             v.ID,
				Number:         v.Number,
				LinerModelCode: v.LinerModelCode,
				Class:          v.Class,
			}
		}
		s.respond(w, r, 200, seatsResponse)

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

		ticketsResponse := make([]Ticket, len(*ticket))
		for i, v := range *ticket {
			ticketsResponse[i] = Ticket{
				Number:                  v.Number,
				PassengerLastName:       v.PassengerLastName,
				PassengerGivenName:      v.PassengerGivenName,
				PassengerBirthDate:      v.PassengerBirthDate,
				PassengerPassportNumber: v.PassengerPassportNumber,
				PassengerSex:            v.PassengerSex,
				PurchaseID:              v.PurchaseID,
			}
		}
		s.respond(w, r, 200, ticketsResponse)
	}
}

func (s *server) handleTicketGetDeleteUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ticketNo, err := strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if r.Method == http.MethodGet {
			p, err := s.store.Ticket().Find(ticketNo)
			if err != nil {
				if err == sql.ErrNoRows {
					s.error(w, r, http.StatusNotFound, errRequestedItemDoesNotExist)
					return
				}
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusOK, &Ticket{
				Number:                  p.Number,
				PassengerLastName:       p.PassengerLastName,
				PassengerGivenName:      p.PassengerGivenName,
				PassengerBirthDate:      p.PassengerBirthDate,
				PassengerPassportNumber: p.PassengerPassportNumber,
				PassengerSex:            p.PassengerSex,
				PurchaseID:              p.PurchaseID,
			})
		}

		if r.Method == http.MethodDelete {
			err = s.store.Ticket().Delete(ticketNo)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			s.respond(w, r, http.StatusNoContent, nil)
			return
		}

		if r.Method == http.MethodPut {

		}
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
