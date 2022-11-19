package main

import (
	"errors"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var ErrTooYoung = errors.New("person is too young")

func checkAgeOver18(value interface{}) error {
	v, _ := value.(time.Time)
	today := time.Now().In(v.Location())
	ty, tm, td := today.Date()
	today = time.Date(ty, tm, td, 0, 0, 0, 0, time.UTC)
	by, bm, bd := v.Date()
	v = time.Date(by, bm, bd, 0, 0, 0, 0, time.UTC)
	if today.Before(v) {
		return ErrTooYoung
	}
	age := ty - by
	if v.AddDate(age, 0, 0).After(today) {
		age--
	}
	if age < 18 {
		return ErrTooYoung
	}
	return nil
}

type Airport struct {
	IATACode string `json:"iata_code"`
	City     string `json:"city"`
	Timezone string `json:"timezone"`
}

func (a *Airport) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.IATACode, validation.Required, validation.Length(3, 3), is.Alpha),
		validation.Field(&a.City, validation.Required, validation.Length(4, 64)),
		validation.Field(&a.Timezone, validation.Required, validation.Length(4, 64)),
	)
}

type BookingOffice struct {
	ID          int    `json:"id"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
}

func (o *BookingOffice) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ID),
		validation.Field(&o.Address, validation.Required),
		validation.Field(&o.PhoneNumber, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{11,15}$"))),
	)
}

type Cashier struct {
	ID         int    `json:"id"`
	Login      string `json:"login"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	Password   string `json:"password,omitempty"`
}

func (c *Cashier) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ID),
		validation.Field(&c.Login, validation.Required, validation.Length(3, 32), is.Alphanumeric),
		validation.Field(&c.FirstName, validation.Required, validation.Length(3, 64)),
		validation.Field(&c.LastName, validation.Required, validation.Length(3, 64)),
		validation.Field(&c.MiddleName, validation.Length(3, 64)),
		validation.Field(&c.Password, validation.Length(6, 72)),
	)
}

type Flight struct {
	ID        int        `json:"id"`
	DepDate   *time.Time `json:"dep_date"`
	LineCode  string     `json:"line_code"`
	IsHot     bool       `json:"is_hot"`
	LinerCode string     `json:"liner_code"`
}

func (f *Flight) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.ID),
		validation.Field(&f.DepDate, validation.Required), //  validation.Date(time.RFC3339)

		validation.Field(&f.LineCode, validation.Required, validation.Length(3, 6), validation.Match(regexp.MustCompile("^[A-Z]{2}[0-9]{1,4}$"))),
		validation.Field(&f.IsHot),
		validation.Field(&f.LinerCode, validation.Required, validation.Length(3, 7), validation.Match(regexp.MustCompile(("^[A-Z]{2}[0-9]{1,5}$")))),
	)
}

type FlightInTicket struct {
	ID       int `json:"id"`
	FlightID int `json:"flight_id"`
	SeatID   int `json:"seat_id"`
	TicketID int `json:"ticket_id"`
}

func (f *FlightInTicket) Validate() error {
	return validation.ValidateStruct(f,
		validation.Field(&f.ID),
		validation.Field(&f.FlightID, validation.Required),
		validation.Field(&f.SeatID, validation.Required),
		validation.Field(&f.TicketID, validation.Required),
	)
}

type Line struct {
	LineCode   string  `json:"line_code"`
	DepTime    string  `json:"dep_time"`
	ArrTime    string  `json:"arr_time"`
	BasePrice  float64 `json:"base_price"`
	DepAirport string  `json:"dep_airport"`
	ArrAirport string  `json:"arr_airport"`
}

func (l *Line) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.LineCode, validation.Required, validation.Length(3, 6), validation.Match(regexp.MustCompile("^[A-Z]{2}[0-9]{1,4}$"))),
		validation.Field(&l.DepTime, validation.Required),
		validation.Field(&l.ArrTime, validation.Required),
		validation.Field(&l.BasePrice, validation.Required, validation.Min(0.0)),
		validation.Field(&l.DepAirport, validation.Required, validation.Length(3, 3), is.Alpha),
		validation.Field(&l.ArrAirport, validation.Required, validation.Length(3, 3), is.Alpha),
	)
}

type Liner struct {
	IATACode  string `json:"iata_code"`
	ModelCode string `json:"model_code"`
}

func (l *Liner) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.IATACode, validation.Required, validation.Length(7, 7), is.Alphanumeric),
		validation.Field(&l.ModelCode, validation.Required, validation.Length(4, 4), is.Alphanumeric),
	)
}

type LinerModel struct {
	IATATypeCode string `json:"iata_type_code"`
	Name         string `json:"name"`
}

func (m *LinerModel) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.IATATypeCode, validation.Required, validation.Length(4, 4), is.Alphanumeric),
		validation.Field(&m.Name, validation.Required, validation.Length(3, 64)),
	)
}

type Purchase struct {
	ID              int        `json:"id"`
	Date            *time.Time `json:"date"`
	BookingOfficeID int        `json:"booking_office_id"`
	TotalPrice      float64    `json:"total_price"`
	ContactPhone    string     `json:"contact_phone"`
	ContactEmail    string     `json:"contact_email"`
	CashierID       string     `json:"cashier_id"`
}

func (p *Purchase) Validate() error {
	return validation.ValidateStruct(p,
		validation.Field(&p.ID),
		validation.Field(&p.Date, validation.Required),
		validation.Field(&p.BookingOfficeID, validation.Required),
		validation.Field(&p.TotalPrice, validation.Required, validation.Min(0.0)),
		validation.Field(&p.ContactPhone, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{11,15}$"))),
		validation.Field(&p.ContactEmail, validation.Required, is.Email),
		validation.Field(&p.CashierID, validation.Required),
	)
}

type Seat struct {
	ID             int    `json:"id"`
	Number         string `json:"number"`
	Class          string `json:"class"`
	LinerModelCode string `json:"model_code"`
}

func (s *Seat) Validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.ID),
		validation.Field(&s.Number, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{1,2}[A-Z]{1}$"))),
		validation.Field(&s.Class, validation.Required, validation.In("J", "W", "Y")),
		validation.Field(&s.LinerModelCode, validation.Required, validation.Required, validation.Length(4, 4), is.Alphanumeric),
	)
}

type Ticket struct {
	ID                      int        `json:"id"`
	PassengerLastName       string     `json:"passenger_last_name"`
	PassengerGivenName      string     `json:"passenger_given_name"`
	PassengerBirthDate      *time.Time `json:"passenger_birth_date"`
	PassengerPassportNumber string     `json:"passenger_passport_number"`
	PassengerSex            uint8      `json:"passenger_sex"`
	PurchaseID              int        `json:"purchase_id"`
}

func (t *Ticket) Validate() error {
	return validation.ValidateStruct(t,
		validation.Field(&t.ID),
		validation.Field(&t.PassengerLastName, validation.Required, validation.Length(3, 64)),
		validation.Field(&t.PassengerGivenName, validation.Required, validation.Length(3, 128)),
		validation.Field(&t.PassengerBirthDate, validation.Required, validation.By(checkAgeOver18)),
		validation.Field(&t.PassengerPassportNumber, validation.Required, validation.Length(10, 10), is.Digit),
		validation.Field(&t.PassengerSex, validation.Required, validation.In(uint8(1), uint8(2))),
		validation.Field(&t.PurchaseID, validation.Required),
	)
}

type AirportList struct {
	Items      []Airport `json:"items"`
	TotalCount int       `json:"total_count"`
}

type BookingOfficeList struct {
	Items      []BookingOffice `json:"items"`
	TotalCount int             `json:"total_count"`
}

type CashierList struct {
	Items      []Cashier `json:"items"`
	TotalCount int       `json:"total_count"`
}

type FlightList struct {
	Items      []Flight `json:"items"`
	TotalCount int      `json:"total_count"`
}

type FlightInTicketList struct {
	Items      []FlightInTicket `json:"items"`
	TotalCount int              `json:"total_count"`
}

type LineList struct {
	Items      []Line `json:"items"`
	TotalCount int    `json:"total_count"`
}

type LinerList struct {
	Items      []Liner `json:"items"`
	TotalCount int     `json:"total_count"`
}

type LinerModelList struct {
	Items      []LinerModel `json:"items"`
	TotalCount int          `json:"total_count"`
}

type PurchaseList struct {
	Items      []Purchase `json:"items"`
	TotalCount int        `json:"total_count"`
}

type SeatList struct {
	Items      []Seat `json:"items"`
	TotalCount int    `json:"total_count"`
}

type TicketList struct {
	Items      []Ticket `json:"items"`
	TotalCount int      `json:"total_count"`
}
