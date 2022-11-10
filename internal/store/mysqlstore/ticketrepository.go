package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type TicketRepository struct {
	store *Store
}

func (r *TicketRepository) Find(number int64) (*store.TicketModel, error) {
	ticket := &store.TicketModel{}
	if err := r.store.db.Get(ticket, "SELECT * FROM ticket WHERE number = ?", number); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) FindAll(row_count, offset int) (*[]store.TicketModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}
	tickets := &[]store.TicketModel{}
	if err := r.store.db.Select(tickets, "SELECT * FROM ticket ORDER BY number LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) Delete(number int64) error {
	res, err := r.store.db.Exec("DELETE FROM ticket WHERE number = ?", number)
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
