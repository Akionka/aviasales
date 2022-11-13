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

func (r *FlightRepository) Find(depDate string, lineCode string) (*store.FlightModel, error) {
	flight := &store.FlightModel{}
	if err := r.store.db.Get(flight, "SELECT * FROM flight WHERE dep_date = ? AND line_code = ?", depDate, lineCode); err != nil {
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

func (r *FlightRepository) Update(depDate string, lineCode string, f *store.FlightModel) error {
	res, err := r.store.db.Exec("UPDATE flight SET dep_date = ?, line_code = ?, is_hot = ?, liner_code = ? WHERE dep_date = ? AND line_code = ?",
		f.DepDate,
		f.LineCode,
		f.IsHot,
		f.LinerCode,
		depDate,
		lineCode,
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

func (r *FlightRepository) Delete(depDate string, lineCode string) error {
	res, err := r.store.db.Exec("DELETE FROM flight WHERE dep_date = ? AND line_code = ?", depDate, lineCode)
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
