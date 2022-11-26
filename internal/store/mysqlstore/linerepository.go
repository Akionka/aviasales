// Файл internal\store\mysqlstore\linerepository.go содержит код для работы с таблицей Рейс
package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type LineRepository struct {
	store *Store
}

func (r *LineRepository) Create(l *store.LineModel) error {
	_, err := r.store.db.Exec("INSERT INTO line (line_code, dep_time, arr_time, base_price, dep_airport, arr_airport) VALUES (?, ?, ?, ?, ?, ?)",
		l.LineCode,
		l.DepTime,
		l.ArrTime,
		l.BasePrice,
		l.DepAirport,
		l.ArrAirport,
	)
	return err
}

func (r *LineRepository) Find(code string) (*store.LineModel, error) {
	line := &store.LineModel{}
	if err := r.store.db.Get(line, "SELECT * FROM line WHERE line_code = ?", code); err != nil {
		return nil, err
	}
	return line, nil
}

func (r *LineRepository) FindAll(row_count, offset int) (*[]store.LineModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	lines := &[]store.LineModel{}
	if err := r.store.db.Select(lines, "SELECT * FROM line ORDER BY line_code LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return lines, nil
}

func (r *LineRepository) Update(code string, l *store.LineModel) error {
	res, err := r.store.db.Exec("UPDATE line SET line_code = ?, dep_time = ?, arr_time = ?, base_price = ?, dep_airport = ?, arr_airport = ? WHERE line_code = ?",
		l.LineCode,
		l.DepTime,
		l.ArrTime,
		l.BasePrice,
		l.DepAirport,
		l.ArrAirport,
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

func (r *LineRepository) Delete(code string) error {
	res, err := r.store.db.Exec("DELETE FROM line WHERE line_code = ?", code)
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

func (r *LineRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from line")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
