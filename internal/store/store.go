package store

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	Airport() AirportRepository
	BookingOffice() BookingOfficeRepository
	Cashier() CashierRepository
	Flight() FlightRepository
	FlightInTicket() FlightInTicketRepository
	Line() LineRepository
	Liner() LinerRepository
	LinerModel() LinerModelRepository
	Purchase() PurchaseRepository
	Seat() SeatRepository
	Ticket() TicketRepository
}

type AirportModel struct {
	IATACode string `db:"iata_code"`
	City     string `db:"city"`
	Timezone string `db:"timezone"`
}

type BookingOfficeModel struct {
	ID          int    `db:"id"`
	Address     string `db:"address"`
	PhoneNumber string `db:"phone_number"`
}

type CashierModel struct {
	Login      string `db:"login"`
	LastName   string `db:"last_name"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	Password   string `db:"password,omitempty"`
}

type FlightModel struct {
	DepDate   *time.Time `db:"dep_date"`
	LineCode  string     `db:"line_code"`
	IsHot     bool       `db:"is_hot"`
	LinerCode string     `db:"liner_code"`
}

type FlightInTicketModel struct {
	DepDate  *time.Time `db:"dep_date"`
	LineCode string     `db:"line_code"`
	SeatID   int        `db:"seat_id"`
	TicketNo int64      `db:"ticket_no"`
}

type LineModel struct {
	LineCode   string  `db:"line_code"`
	DepTime    string  `db:"dep_time"`
	ArrTime    string  `db:"arr_time"`
	BasePrice  float64 `db:"base_price"`
	DepAirport string  `db:"dep_airport"`
	ArrAirport string  `db:"arr_airport"`
}

type LinerModel struct {
	IATACode  string `db:"iata_code"`
	ModelCode string `db:"model_code"`
}

type LinerModelModel struct {
	IATATypeCode string `db:"iata_type_code"`
	Name         string `db:"name"`
}

type PurchaseModel struct {
	ID              int     `db:"id"`
	Date            string  `db:"date"`
	BookingOfficeID int     `db:"booking_office_id"`
	TotalPrice      float64 `db:"total_price"`
	ContactPhone    string  `db:"contact_phone"`
	ContactEmail    string  `db:"contact_email"`
	CashierLogin    string  `db:"cashier_login"`
}

type SeatModel struct {
	ID             int    `db:"id"`
	Number         string `db:"number"`
	Class          string `db:"class"`
	LinerModelCode string `db:"liner_model_code"`
}

type TicketModel struct {
	Number                  int64  `db:"number"`
	PassengerLastName       string `db:"pass_last_name"`
	PassengerGivenName      string `db:"pass_given_name"`
	PassengerBirthDate      string `db:"pass_birth_date"`
	PassengerPassportNumber string `db:"pass_passport_number"`
	PassengerSex            uint8  `db:"pass_sex"`
	PurchaseID              int    `db:"purchase_id"`
}

// ComparePassword returns true if the password matches and false otherwise
func (c *CashierModel) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)) == nil
}
