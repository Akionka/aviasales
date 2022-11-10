package main

import (
	"context"
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
	tokenSecret                 = []byte("thisisthemostsecrettokenintheentireuniverse")
	errBadAuthorizationToken    = errors.New("bad authorization token")
	errIncorrectLoginOrPassword = errors.New("incorrect login or password")
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
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
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
		s.error(w, r, http.StatusInternalServerError, errors.New("bad"))
	}).Methods("GET", "OPTIONS")

	securedGet := secured.Methods("GET", "OPTIONS").Subrouter()
	securedGet.Use(s.paginateMiddleware)

	securedGet.HandleFunc("/airports", s.handleAirportsGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/booking_offices", s.handleBookingOfficesGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/cashiers", s.handleCashiersGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/flight_in_tickets", s.handleFlightInTicketsGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/flights", s.handleFlightsGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/lines", s.handleLinesGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/liner_models", s.handleLinerModelsGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/liners", s.handleLinersGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/purchases", s.handlePurchasesGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/seats", s.handleSeatsGet()).Methods("GET", "OPTIONS")
	securedGet.HandleFunc("/tickets", s.handleTicketsGet()).Methods("GET", "OPTIONS")

	secured.HandleFunc("/airports/{code}", s.handleAirportsDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/booking_offices/{id:[0-9]+}", s.handleBookingOfficesDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/cashiers/{login}", s.handleCashiersDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/flight_in_tickets/{dep_date}/{line_code}/{seat_id:[0-9]+}/{ticket_no:[0-9]+}", s.handleFlightInTicketsDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/flights/{dep_date}/{line_code}", s.handleFlightsDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/lines/{code}", s.handleLinesDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/liner_models/{code}", s.handleLinerModelsDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/liners/{code}", s.handleLinersDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/purchases/{id:[0-9]+}", s.handlePurchasesDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/seats/{id:[0-9]+}", s.handleSeatsDelete()).Methods("DELETE", "OPTIONS")
	secured.HandleFunc("/tickets/{id:[0-9]+}", s.handleTicketsDelete()).Methods("DELETE", "OPTIONS")
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

func (s *server) handleAirportsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := s.store.Airport().Delete(vars["code"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleBookingOfficesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		offices, err := s.store.BookingOffice().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleBookingOfficesDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		err = s.store.BookingOffice().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleCashiersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		cashiers, err := s.store.Cashier().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleCashiersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := s.store.Cashier().Delete(vars["login"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleFlightInTicketsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		flightInTickets, err := s.store.FlightInTicket().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleFlightInTicketsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["seat_id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		ticketNo, err := strconv.ParseInt(vars["ticket_no"], 10, 64)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		err = s.store.FlightInTicket().Delete(vars["dep_date"], vars["line_code"], id, ticketNo)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleFlightsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		flights, err := s.store.Flight().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleFlightsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := s.store.Flight().Delete(vars["dep_date"], vars["line_code"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleLinesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		lines, err := s.store.Line().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleLinesDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := s.store.Line().Delete(vars["code"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleLinerModelsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		models, err := s.store.LinerModel().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleLinerModelsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := s.store.LinerModel().Delete(vars["code"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleLinersGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		liners, err := s.store.Liner().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleLinersDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		err := s.store.Liner().Delete(vars["code"])
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handlePurchasesGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		purchases, err := s.store.Purchase().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handlePurchasesDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err = s.store.Purchase().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleSeatsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		seats, err := s.store.Seat().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleSeatsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err = s.store.Seat().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
	}
}

func (s *server) handleTicketsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := r.Context().Value(ctxKeyPagination).(paginationInfo)
		ticket, err := s.store.Ticket().FindAll(p.rowCount, (p.page-1)*p.rowCount)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
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

func (s *server) handleTicketsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		err = s.store.Seat().Delete(id)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusNoContent, nil)
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
