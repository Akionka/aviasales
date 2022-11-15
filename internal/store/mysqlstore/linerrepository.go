package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type LinerRepository struct {
	store *Store
}

func (r *LinerRepository) Create(l *store.LinerModel) error {
	_, err := r.store.db.Exec("INSERT INTO liner (iata_code, model_code) VALUES (?, ?)",
		l.IATACode,
		l.ModelCode,
	)
	return err
}

func (r *LinerRepository) Find(code string) (*store.LinerModel, error) {
	liner := &store.LinerModel{}
	if err := r.store.db.Get(liner, "SELECT * FROM liner WHERE iata_code = ?", code); err != nil {
		return nil, err
	}
	return liner, nil
}

func (r *LinerRepository) FindAll(row_count, offset int) (*[]store.LinerModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}
	liners := &[]store.LinerModel{}
	if err := r.store.db.Select(liners, "SELECT * FROM liner ORDER BY iata_code LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return liners, nil
}

func (r *LinerRepository) Update(code string, l *store.LinerModel) error {
	res, err := r.store.db.Exec("UPDATE liner SET iata_code = ?, model_code = ? WHERE iata_code = ?",
		l.IATACode,
		l.ModelCode,
		code,
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

func (r *LinerRepository) Delete(code string) error {
	res, err := r.store.db.Exec("DELETE FROM liner WHERE iata_code = ?", code)
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

func (r *LinerRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from liner")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
