package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type SeatRepository struct {
	store *Store
}

func (r *SeatRepository) Create(s *store.SeatModel) error {
	_, err := r.store.db.Exec("INSERT INTO seat (id, number, class, liner_model_class) VALUES (?, ?, ?, ?)",
		s.ID,
		s.Number,
		s.Class,
		s.LinerModelCode,
	)
	return err
}

func (r *SeatRepository) Find(id int) (*store.SeatModel, error) {
	seat := &store.SeatModel{}
	if err := r.store.db.Get(seat, "SELECT * FROM seat WHERE id = ?", id); err != nil {
		return nil, err
	}
	return seat, nil
}

func (r *SeatRepository) FindAll(row_count, offset int) (*[]store.SeatModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}
	seats := &[]store.SeatModel{}
	if err := r.store.db.Select(seats, "SELECT * FROM seat ORDER BY id LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return seats, nil
}

func (r *SeatRepository) Delete(id int) error {
	res, err := r.store.db.Exec("DELETE FROM seat WHERE id = ?", id)
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
