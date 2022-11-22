package mysqlstore

import (
	"time"

	"github.com/akionka/aviasales/internal/store"
)

type TicketRepository struct {
	store *Store
}

func (r *TicketRepository) Create(t *store.TicketModel) error {
	_, err := r.store.db.Exec("INSERT INTO ticket (pass_last_name, pass_given_name, pass_birth_date, pass_passport_number, pass_sex, purchase_id) VALUES (?, ?, ?, ?, ?, ?)",
		t.PassengerLastName,
		t.PassengerGivenName,
		t.PassengerBirthDate,
		t.PassengerPassportNumber,
		t.PassengerSex,
		t.PurchaseID,
	)
	return err
}

func (r *TicketRepository) Find(id int) (*store.TicketModel, error) {
	ticket := &store.TicketModel{}
	if err := r.store.db.Get(ticket, "SELECT * FROM ticket WHERE id = ?", id); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *TicketRepository) FindAll(row_count, offset int) (*[]store.TicketModel, error) {
	if row_count < 0 {
		row_count = 0
	}
	if offset < 0 {
		offset = 0
	}
	tickets := &[]store.TicketModel{}
	if err := r.store.db.Select(tickets, "SELECT * FROM ticket ORDER BY id LIMIT ?, ?", offset, row_count); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *TicketRepository) Update(id int, t *store.TicketModel) error {
	res, err := r.store.db.Exec("UPDATE ticket SET pass_last_name = ?, pass_given_name = ?, pass_birth_date = ?, pass_passport_number = ?, pass_sex = ?, purchase_id = ? WHERE id = ?",
		t.PassengerLastName,
		t.PassengerGivenName,
		t.PassengerBirthDate,
		t.PassengerPassportNumber,
		t.PassengerSex,
		t.PurchaseID,
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

func (r *TicketRepository) Delete(id int) error {
	res, err := r.store.db.Exec("DELETE FROM ticket WHERE id = ?", id)
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

func (r *TicketRepository) TotalCount() (int, error) {
	var count int
	row := r.store.db.QueryRow("SELECT COUNT(*) from ticket")
	if err := row.Scan(&count); err != nil {
		return -1, err
	}
	return count, nil
}

func (r *TicketRepository) Report(id int) ([]*store.TicketReportFlightModel, *store.BookingOfficeModel, *store.CashierModel, *store.PurchaseModel, time.Duration, error) {
	var flights []*store.TicketReportFlightModel
	var office store.BookingOfficeModel
	var cashier store.CashierModel
	var purchase store.PurchaseModel
	var totalTime time.Duration

	rows, err := r.store.db.Queryx(`SELECT
	a1.city dep_city,
	a2.city arr_city,
	ADDTIME(f.dep_date, l.dep_time) dep_time_local,
	IF(CONVERT_TZ(ADDTIME(f.dep_date, l.dep_time),
							a1.timezone,
							'GMT') > CONVERT_TZ(ADDTIME(f.dep_date, l.arr_time),
							a2.timezone,
							'GMT'),
			DATE_ADD(ADDTIME(f.dep_date, l.arr_time),
					INTERVAL 1 DAY),
			ADDTIME(f.dep_date, l.arr_time)) arr_time_local_fixed,
	l.line_code,
	s.number,
	s.class
FROM
	ticket t
			INNER JOIN
	flight_in_ticket fit ON t.id = fit.ticket_id
			INNER JOIN
	flight f ON fit.flight_id = f.id
			INNER JOIN
	line l ON f.line_code = l.line_code
			INNER JOIN
	airport a1 ON l.dep_airport = a1.iata_code
			INNER JOIN
	airport a2 ON l.arr_airport = a2.iata_code
			INNER JOIN
	seat s ON s.id = fit.seat_id
WHERE t.id = ?
ORDER BY dep_time`, id)
	if err != nil {
		return flights, &office, &cashier, &purchase, totalTime, err
	}

	for rows.Next() {
		var f store.TicketReportFlightModel
		rows.Scan(&f.DepCity, &f.ArrCity, &f.DepTime, &f.ArrTime, &f.LineCode, &f.SeatNumber, &f.SeatClass)
		flights = append(flights, &f)
	}

	r.store.db.Get(&purchase, "SELECT p.* FROM ticket t INNER JOIN purchase p ON t.purchase_id = p.id WHERE t.id = ?", id)
	r.store.db.Get(&office, "SELECT b.* FROM ticket t INNER JOIN purchase p ON t.purchase_id = p.id INNER JOIN booking_office b ON p.booking_office_id = b.id WHERE t.id = ?", id)
	r.store.db.Get(&cashier, "SELECT c.* FROM ticket t INNER JOIN purchase p ON t.purchase_id = p.id INNER JOIN cashier c ON p.cashier_id = c.id WHERE t.id = ?", id)

	totalTime = flights[len(flights)-1].ArrTime.Sub(*flights[0].DepTime)

	return flights, &office, &cashier, &purchase, totalTime, nil
}
