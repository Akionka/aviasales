// Файл internal\store\store.go содержит описания таблиц на уровне БД
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
	Timezone() TimezoneRepository
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
	ID         int    `db:"id"`
	Login      string `db:"login"`
	LastName   string `db:"last_name"`
	FirstName  string `db:"first_name"`
	MiddleName string `db:"middle_name"`
	Password   string `db:"password,omitempty"`
	RoleID     int    `db:"role_id"`
}

// ComparePassword returns true if the password matches and false otherwise
func (c *CashierModel) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.Password), []byte(password)) == nil
}

func (c *CashierModel) SetPassword(password string) error {
	p, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	c.Password = string(p)
	return nil
}

type FlightModel struct {
	ID        int       `db:"id"`
	DepDate   time.Time `db:"dep_date"`
	LineCode  string    `db:"line_code"`
	IsHot     bool      `db:"is_hot"`
	LinerCode string    `db:"liner_code"`
}

type FlightInTicketModel struct {
	ID       int `db:"id"`
	FlightID int `db:"flight_id"`
	SeatID   int `db:"seat_id"`
	TicketID int `db:"ticket_id"`
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
	ID              int       `db:"id"`
	Date            time.Time `db:"date"`
	BookingOfficeID int       `db:"booking_office_id"`
	TotalPrice      float64   `db:"total_price"`
	ContactPhone    string    `db:"contact_phone"`
	ContactEmail    string    `db:"contact_email"`
	CashierID       string    `db:"cashier_id"`
}

type SeatModel struct {
	ID             int    `db:"id"`
	Number         string `db:"number"`
	Class          string `db:"class"`
	LinerModelCode string `db:"model_code"`
}

type TicketModel struct {
	ID                      int       `db:"id"`
	PassengerLastName       string    `db:"pass_last_name"`
	PassengerGivenName      string    `db:"pass_given_name"`
	PassengerBirthDate      time.Time `db:"pass_birth_date"`
	PassengerPassportNumber string    `db:"pass_passport_number"`
	PassengerSex            uint8     `db:"pass_sex"`
	PurchaseID              int       `db:"purchase_id"`
}

type TicketReportFlightModel struct {
	DepCity      string    `db:"dep_city" json:"dep_city"`
	ArrCity      string    `db:"arr_city" json:"arr_city"`
	DepTimeLocal time.Time `db:"dep_time_local" json:"dep_time_local"`
	DepTimeGMT   time.Time `db:"dep_time_gmt" json:"dep_time_gmt"`
	ArrTimeLocal time.Time `db:"arr_time_local" json:"arr_time_local"`
	ArrTimeGMT   time.Time `db:"arr_time_gmt" json:"arr_time_gmt"`
	LineCode     string    `db:"line_code" json:"line_code"`
	SeatNumber   string    `db:"number" json:"number"`
	SeatClass    string    `db:"class" json:"class"`
	Price        float64   `db:"price" json:"price"`
}

type RoleModel struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}
