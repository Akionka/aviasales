package mysqlstore

import (
	"github.com/akionka/aviasales/internal/store"
)

type FlightRepository struct {
	store *Store
}

func (r *FlightRepository) Create(f *store.FlightModel) error {
	_, err := r.store.db.Exec("INSERT INTO flight (dep_date, line_code, is_hot, liner_code) VALUES (?, ?, ?, ?)",
		f.DepDate,
		f.LineCode,
		f.IsHot,
		f.LinerCode,
	)
	return err
}

func (r *FlightRepository) Find(id int) (*store.FlightModel, error) {
	flight := &store.FlightModel{}
	if err := r.store.db.Get(flight, "SELECT * FROM flight WHERE id = ?", id); err != nil {
		return nil, err
	}
	return flight, nil
}

func (r *FlightRepository) FindAll(row_count, offset int) (*[]store.FlightModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	flights := &[]store.FlightModel{}
	if err := r.store.db.Select(flights, "SELECT * FROM flight ORDER BY dep_date, line_code LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return flights, nil
}

func (r *FlightRepository) Update(id int, f *store.FlightModel) error {
	res, err := r.store.db.Exec("UPDATE flight SET dep_date = ?, line_code = ?, is_hot = ?, liner_code = ? WHERE id = ?",
		f.DepDate,
		f.LineCode,
		f.IsHot,
		f.LinerCode,
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

func (r *FlightRepository) Delete(id int) error {
	res, err := r.store.db.Exec("DELETE FROM flight WHERE id = ?", id)
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

func (r *FlightRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from flight")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
