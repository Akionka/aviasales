// Файл internal\store\mysqlstore\linermodelrepository.go содержит код  для работы с таблицей Модель самолета
package mysqlstore

import "github.com/akionka/aviasales/internal/store"

type LinerModelRepository struct {
	store *Store
}

func (r *LinerModelRepository) Create(m *store.LinerModelModel) error {
	_, err := r.store.db.Exec("INSERT INTO liner_model (iata_type_code, name) VALUES (?, ?)",
		m.IATATypeCode,
		m.Name,
	)
	return err
}

func (r *LinerModelRepository) Find(code string) (*store.LinerModelModel, error) {
	linerModel := &store.LinerModelModel{}
	if err := r.store.db.Get(linerModel, "SELECT * FROM liner_model WHERE iata_type_code = ?", code); err != nil {
		return nil, err
	}
	return linerModel, nil
}

func (r *LinerModelRepository) FindAll(row_count, offset int) (*[]store.LinerModelModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}

	linerModels := &[]store.LinerModelModel{}
	if err := r.store.db.Select(linerModels, "SELECT * FROM liner_model ORDER BY iata_type_code LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return linerModels, nil
}

func (r *LinerModelRepository) Update(code string, m *store.LinerModelModel) error {
	res, err := r.store.db.Exec("UPDATE liner_model SET iata_type_code = ?, name = ? WHERE iata_type_code = ?",
		m.IATATypeCode,
		m.Name,
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

func (r *LinerModelRepository) Delete(code string) error {
	res, err := r.store.db.Exec("DELETE FROM liner_model WHERE iata_type_code = ?", code)
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

func (r *LinerModelRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from liner_model")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}
