package mysqlstore

import (
	"errors"

	"github.com/akionka/aviasales/internal/store"
	"github.com/jmoiron/sqlx"
)

var ErrDeletedItemDoesNotExist = errors.New("the item you delete does not exist")
var ErrUpdatedItemDoesNotExist = errors.New("the item you update does not exist")

type Store struct {
	db                       *sqlx.DB
	airportRepository        *AirportRepository
	bookingOfficeRepository  *BookingOfficeRepository
	cashierRepository        *CashierRepository
	flightRepository         *FlightRepository
	flightInTicketRepository *FlightInTicketRepository
	lineRepository           *LineRepository
	linerRepository          *LinerRepository
	linerModelRepository     *LinerModelRepository
	purchaseRepository       *PurchaseRepository
	seatRepository           *SeatRepository
	ticketRepository         *TicketRepository
}

func New(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Airport() store.AirportRepository {
	if s.airportRepository != nil {
		return s.airportRepository
	}
	s.airportRepository = &AirportRepository{
		store: s,
	}
	return s.airportRepository
}

func (s *Store) BookingOffice() store.BookingOfficeRepository {
	if s.bookingOfficeRepository != nil {
		return s.bookingOfficeRepository
	}
	s.bookingOfficeRepository = &BookingOfficeRepository{
		store: s,
	}
	return s.bookingOfficeRepository
}

func (s *Store) Cashier() store.CashierRepository {
	if s.cashierRepository != nil {
		return s.cashierRepository
	}
	s.cashierRepository = &CashierRepository{
		store: s,
	}
	return s.cashierRepository
}

func (s *Store) Flight() store.FlightRepository {
	if s.flightRepository != nil {
		return s.flightRepository
	}
	s.flightRepository = &FlightRepository{
		store: s,
	}
	return s.flightRepository
}

func (s *Store) FlightInTicket() store.FlightInTicketRepository {
	if s.flightInTicketRepository != nil {
		return s.flightInTicketRepository
	}
	s.flightInTicketRepository = &FlightInTicketRepository{
		store: s,
	}
	return s.flightInTicketRepository
}

func (s *Store) Line() store.LineRepository {
	if s.lineRepository != nil {
		return s.lineRepository
	}
	s.lineRepository = &LineRepository{
		store: s,
	}
	return s.lineRepository
}

func (s *Store) Liner() store.LinerRepository {
	if s.linerRepository != nil {
		return s.linerRepository
	}
	s.linerRepository = &LinerRepository{
		store: s,
	}
	return s.linerRepository
}

func (s *Store) LinerModel() store.LinerModelRepository {
	if s.linerModelRepository != nil {
		return s.linerModelRepository
	}
	s.linerModelRepository = &LinerModelRepository{
		store: s,
	}
	return s.linerModelRepository
}

func (s *Store) Purchase() store.PurchaseRepository {
	if s.purchaseRepository != nil {
		return s.purchaseRepository
	}
	s.purchaseRepository = &PurchaseRepository{
		store: s,
	}
	return s.purchaseRepository
}

func (s *Store) Seat() store.SeatRepository {
	if s.seatRepository != nil {
		return s.seatRepository
	}
	s.seatRepository = &SeatRepository{
		store: s,
	}
	return s.seatRepository
}

func (s *Store) Ticket() store.TicketRepository {
	if s.ticketRepository != nil {
		return s.ticketRepository
	}
	s.ticketRepository = &TicketRepository{
		store: s,
	}
	return s.ticketRepository
}
