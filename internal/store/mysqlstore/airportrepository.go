package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type AirportRepository struct {
	store *Store
}

func (r *AirportRepository) Create(a *store.AirportModel) error {
	_, err := r.store.db.Exec("INSERT INTO airport (iata_code, city, timezone) VALUES (?, ?, ?)",
		a.IATACode,
		a.City,
		a.Timezone,
	)
	return err
}

func (r *AirportRepository) Find(code string) (*store.AirportModel, error) {
	airport := &store.AirportModel{}
	if err := r.store.db.Get(airport, "SELECT * FROM airport WHERE iata_code = ?", code); err != nil {
		return nil, err
	}
	return airport, nil
}

func (r *AirportRepository) FindAll(row_count, offset int) (*[]store.AirportModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	airports := &[]store.AirportModel{}
	if err := r.store.db.Select(airports, "SELECT * FROM airport ORDER BY iata_code LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return airports, nil
}

func (r *AirportRepository) Delete(code string) error {
	res, err := r.store.db.Exec("DELETE FROM airport WHERE iata_code = ?", code)
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
