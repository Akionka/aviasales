package main

import (
	"time"
)

type Airport struct {
	IATACode string `json:"iata_code"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

type BookingOffice struct {
	ID          int    `json:"id"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

type Cashier struct {
	Login      string `json:"login"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
}

type Flight struct {
	DepDate   *time.Time `json:"dep_date"`
	LineCode  string    `json:"line_code"`
	IsHot     bool      `json:"is_hot"`
	LinerCode string    `json:"liner_code"`
}

type FlightInTicket struct {
	DepDate  string `json:"dep_date"`
	LineCode string `json:"line_code"`
	SeatID   int    `json:"seat_id"`
	TicketNo int64  `json:"ticket_no"`
}

type Line struct {
	LineCode   string  `json:"line_code"`
	DepTime    string  `json:"dep_time"`
	ArrTime    string  `json:"arr_time"`
	BasePrice  float64 `json:"base_price"`
	DepAirport string  `json:"dep_airport"`
	ArrAirport string  `json:"arr_airport"`
}

type Liner struct {
	IATACode  string `json:"iata_code"`
	ModelCode string `json:"model_code"`
}

type LinerModel struct {
	IATATypeCode string `json:"iata_type_code"`
	Name         string `json:"name"`
}

type Purchase struct {
	ID              int     `json:"id"`
	Date            string  `json:"date"`
	BookingOfficeID int     `json:"booking_office_id"`
	TotalPrice      float64 `json:"total_price"`
	ContactPhone    string  `json:"contact_phone"`
	ContactEmail    string  `json:"contact_email"`
	CashierLogin    string  `json:"cashier_login"`
}

type Seat struct {
	ID             int    `json:"id"`
	Number         string `json:"number"`
	Class          string `json:"class"`
	LinerModelCode string `json:"liner_model_code"`
}

type Ticket struct {
	Number                  int64  `json:"number"`
	PassengerLastName       string `json:"passenger_last_name"`
	PassengerGivenName      string `json:"passenger_given_name"`
	PassengerBirthDate      string `json:"passenger_birth_date"`
	PassengerPassportNumber string `json:"passenger_passport_number"`
	PassengerSex            uint8  `json:"passenger_sex"`
	PurchaseID              int    `json:"purchase_id"`
}
