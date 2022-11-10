package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type CashierRepository struct {
	store *Store
}

func (r *CashierRepository) Create(c *store.CashierModel) error {
	_, err := r.store.db.Exec("INSERT INTO сashier (login, last_name, first_name, middle_name, password) VALUES (?, ?, ?, ?, ?)",
		c.Login,
		c.LastName,
		c.FirstName,
		c.MiddleName,
		c.Password,
	)
	return err
}

func (r *CashierRepository) Find(login string) (*store.CashierModel, error) {
	cashier := &store.CashierModel{}
	if err := r.store.db.Get(cashier, "SELECT * FROM cashier WHERE login = ?", login); err != nil {
		return nil, err
	}
	return cashier, nil
}

func (r *CashierRepository) FindAll(row_count, offset int) (*[]store.CashierModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	cashiers := &[]store.CashierModel{}
	if err := r.store.db.Select(cashiers, "SELECT * FROM cashier ORDER BY login LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return cashiers, nil
}

func (r *CashierRepository) Delete(login string) error {
	res, err := r.store.db.Exec("DELETE FROM cashier WHERE login = ?", login)
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
