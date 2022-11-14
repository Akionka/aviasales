package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type TicketRepository struct {
	store *Store
}

func (r *TicketRepository) Create(t *store.TicketModel) error {
	_, err := r.store.db.Exec("INSERT INTO ticket (number, pass_last_name, pass_given_name, pass_birth_date, pass_passport_number, pass_sex, purchase_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		t.Number,
		t.PassengerLastName,
		t.PassengerGivenName,
		t.PassengerBirthDate,
		t.PassengerPassportNumber,
		t.PassengerSex,
		t.PurchaseID,
	)
	return err
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

func (r *TicketRepository) Update(number int64, t *store.TicketModel) error {
	res, err := r.store.db.Exec("UPDATE ticket SET number = ?, pass_last_name = ?, pass_given_name = ?, pass_birth_date = ?, pass_passport_number = ?, pass_sex = ?, purchase_id = ? WHERE number = ?",
		t.Number,
		t.PassengerLastName,
		t.PassengerGivenName,
		t.PassengerBirthDate,
		t.PassengerPassportNumber,
		t.PassengerSex,
		t.PurchaseID,
		number,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrUpdatedItemDoesNotExist
	}
	return err
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

func (r *TicketRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from ticket")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}