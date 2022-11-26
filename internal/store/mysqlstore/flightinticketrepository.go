// Файл internal\store\mysqlstore\flightinticketrepository.go содержит код для работы с таблицей Полёт в билете
package mysqlstore

import (
	"github.com/akionka/aviasales/internal/store"
)

type FlightInTicketRepository struct {
	store *Store
}

func (r *FlightInTicketRepository) Create(f *store.FlightInTicketModel) error {
	_, err := r.store.db.Exec("INSERT INTO flight_in_ticket (flight_id, seat_id, ticket_id) VALUES (?, ?, ?)",
		f.FlightID,
		f.SeatID,
		f.TicketID,
	)
	return err
}

func (r *FlightInTicketRepository) Find(id int) (*store.FlightInTicketModel, error) {
	flightInTicket := &store.FlightInTicketModel{}
	if err := r.store.db.Get(flightInTicket, "SELECT * FROM flight_in_ticket WHERE id = ?", id); err != nil {
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
	if err := r.store.db.Select(flightInTickets, "SELECT * FROM flight_in_ticket ORDER BY id LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return flightInTickets, nil
}

func (r *FlightInTicketRepository) Update(id int, f *store.FlightInTicketModel) error {
	res, err := r.store.db.Exec("UPDATE flight_in_ticket SET flight_id = ?, seat_id = ?, ticket_id = ? WHERE id = ?",
		f.FlightID,
		f.SeatID,
		f.TicketID,
		id,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrNoChanges
	}
	return err
}

func (r *FlightInTicketRepository) Delete(id int) error {
	res, err := r.store.db.Exec("DELETE FROM flight_in_ticket WHERE id = ?", id)
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

func (r *FlightInTicketRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from flight_in_ticket")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
