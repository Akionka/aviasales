package mysqlstore

import (
	"github.com/akionka/aviasales/internal/store"
)

type FlightInTicketRepository struct {
	store *Store
}

func (r *FlightInTicketRepository) Create(f *store.FlightInTicketModel) error {
	_, err := r.store.db.Exec("INSERT INTO flight_in_ticket (dep_date, line_code, seat_id, ticket_no) VALUES (?, ?, ?, ?)",
		f.DepDate,
		f.LineCode,
		f.SeatID,
		f.TicketNo,
	)
	return err
}

func (r *FlightInTicketRepository) Find(depDate string, lineCode string, seatID int, ticketNo int64) (*store.FlightInTicketModel, error) {
	flightInTicket := &store.FlightInTicketModel{}
	if err := r.store.db.Get(flightInTicket, "SELECT * FROM flight_in_ticket WHERE dep_date = ? AND line_code = ? AND seat_id = ? AND ticket_no = ?",
		depDate,
		lineCode,
		seatID,
		ticketNo,
	); err != nil {
		return nil, err
	}
	return flightInTicket, nil
}

func (r *FlightInTicketRepository) FindAll(row_count, offset int) (*[]store.FlightInTicketModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	flightInTickets := &[]store.FlightInTicketModel{}
	if err := r.store.db.Select(flightInTickets, "SELECT * FROM flight_in_ticket ORDER BY dep_date, line_code, seat_id, ticket_no LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return flightInTickets, nil
}

func (r *FlightInTicketRepository) Delete(depDate string, lineCode string, seatID int, ticketNo int64) error {
	res, err := r.store.db.Exec("DELETE FROM flight_in_ticket WHERE dep_date = ? AND line_code = ? AND seat_id = ? AND ticket_no = ?",
		depDate,
		lineCode,
		seatID,
		ticketNo,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrDeletedItemDoesNotExist
	}
	return nil
}
