package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type PurchaseRepository struct {
	store *Store
}

func (r *PurchaseRepository) Create(p *store.PurchaseModel) error {
	_, err := r.store.db.Exec("INSERT INTO purchase (id, date, booking_office_id, total_price, contact_phone, contact_email, cashier_login) VALUES (?, ?, ?, ?, ?, ?, ?)",
		p.ID,
		p.Date,
		p.BookingOfficeID,
		p.TotalPrice,
		p.ContactPhone,
		p.ContactEmail,
		p.CashierLogin,
	)
	return err
}

func (r *PurchaseRepository) Find(id int) (*store.PurchaseModel, error) {
	purchase := &store.PurchaseModel{}
	if err := r.store.db.Get(purchase, "SELECT * FROM purchase WHERE id = ?", id); err != nil {
		return nil, err
	}
	return purchase, nil
}

func (r *PurchaseRepository) FindAll(row_count, offset int) (*[]store.PurchaseModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}
	purchases := &[]store.PurchaseModel{}
	if err := r.store.db.Select(purchases, "SELECT * FROM purchase ORDER BY id LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return purchases, nil
}

func (r *PurchaseRepository) Update(id int, p *store.PurchaseModel) error {
	res, err := r.store.db.Exec("UPDATE purchase SET id = ?, date = ?, booking_office_id = ?, total_price = ?, contact_phone = ?, contact_email = ?, cashier_login = ? WHERE id = ?",
		p.ID,
		p.Date,
		p.BookingOfficeID,
		p.TotalPrice,
		p.ContactPhone,
		p.ContactEmail,
		p.CashierLogin,
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
		return ErrUpdatedItemDoesNotExist
	}
	return err
}

func (r *PurchaseRepository) Delete(id int) error {
	res, err := r.store.db.Exec("DELETE FROM purchase WHERE id = ?", id)
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

func (r *PurchaseRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from purchase")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}