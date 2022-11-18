package mysqlstore

type TimezoneRepository struct {
	store *Store
}

func (r *TimezoneRepository) FindAll() ([]string, error) {
	var timezones []string
	if err := r.store.db.Select(&timezones, "SELECT name FROM mysql.time_zone_name"); err != nil {
		return timezones, err
	}
	return timezones, nil
}
