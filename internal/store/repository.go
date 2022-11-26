// Файл internal\store\repository.go содержит интерфейсы для работы с БД
package store

import "time"

type AirportRepository interface {
	Create(*AirportModel) error
	Find(code string) (*AirportModel, error)
	FindAll(row_count, offset int) (*[]AirportModel, error)
	Update(code string, a *AirportModel) error
	Delete(code string) error
	TotalCount() (int, error)
}

type BookingOfficeRepository interface {
	Create(*BookingOfficeModel) error
	Find(id int) (*BookingOfficeModel, error)
	FindAll(row_count, offset int) (*[]BookingOfficeModel, error)
	Update(id int, o *BookingOfficeModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type CashierRepository interface {
	Create(*CashierModel) error
	Find(id int) (*CashierModel, error)
	FindByLogin(login string) (*CashierModel, error)
	FindAll(row_count, offset int) (*[]CashierModel, error)
	Update(id int, c *CashierModel) error
	UpdatePassword(*CashierModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type FlightInTicketRepository interface {
	Create(*FlightInTicketModel) error
	Find(id int) (*FlightInTicketModel, error)
	FindAll(row_count, offset int) (*[]FlightInTicketModel, error)
	Update(id int, f *FlightInTicketModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type FlightRepository interface {
	Create(*FlightModel) error
	Find(id int) (*FlightModel, error)
	FindAll(row_count, offset int) (*[]FlightModel, error)
	Update(id int, f *FlightModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type LineRepository interface {
	Create(*LineModel) error
	Find(code string) (*LineModel, error)
	FindAll(row_count, offset int) (*[]LineModel, error)
	Update(code string, l *LineModel) error
	Delete(code string) error
	TotalCount() (int, error)
}

type LinerModelRepository interface {
	Create(*LinerModelModel) error
	Find(code string) (*LinerModelModel, error)
	FindAll(row_count, offset int) (*[]LinerModelModel, error)
	Update(code string, m *LinerModelModel) error
	Delete(code string) error
	TotalCount() (int, error)
}

type LinerRepository interface {
	Create(*LinerModel) error
	Find(code string) (*LinerModel, error)
	FindAll(row_count, offset int) (*[]LinerModel, error)
	Update(code string, l *LinerModel) error
	Delete(code string) error
	TotalCount() (int, error)
}

type PurchaseRepository interface {
	Create(*PurchaseModel) error
	Find(id int) (*PurchaseModel, error)
	FindAll(row_count, offset int) (*[]PurchaseModel, error)
	Update(id int, p *PurchaseModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type SeatRepository interface {
	Create(*SeatModel) error
	Find(id int) (*SeatModel, error)
	FindAll(row_count, offset int) (*[]SeatModel, error)
	Update(id int, s *SeatModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type TicketRepository interface {
	Report(id int) ([]*TicketReportFlightModel, *BookingOfficeModel, *CashierModel, *PurchaseModel, time.Duration, error)
	Create(*TicketModel) error
	Find(id int) (*TicketModel, error)
	FindAll(row_count, offset int) (*[]TicketModel, error)
	Update(id int, t *TicketModel) error
	Delete(id int) error
	TotalCount() (int, error)
}

type TimezoneRepository interface {
	FindAll() ([]string, error)
}
