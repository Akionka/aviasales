package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type BookingOfficeRepository struct {
	store *Store
}

func (r *BookingOfficeRepository) Create(o *store.BookingOfficeModel) error {
	_, err := r.store.db.Exec("INSERT INTO booking_office (id, address, phone_number) VALUES (?, ?, ?)", o.ID, o.Address, o.PhoneNumber)
	return err
}

func (r *BookingOfficeRepository) Find(id int) (*store.BookingOfficeModel, error) {
	office := &store.BookingOfficeModel{}
	if err := r.store.db.Get(office, "SELECT * FROM booking_office WHERE id = ?", id); err != nil {
		return nil, err
	}
	return office, nil
}

func (r *BookingOfficeRepository) FindAll(row_count, offset int) (*[]store.BookingOfficeModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	offices := &[]store.BookingOfficeModel{}
	if err := r.store.db.Select(offices, "SELECT * FROM booking_office ORDER BY id LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return offices, nil
}

func (r *BookingOfficeRepository) Delete(id int) error {
	res, err := r.store.db.Exec("DELETE FROM booking_office WHERE id = ?", id)
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
