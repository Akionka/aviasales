package store

type AirportRepository interface {
	Find(code string) (*AirportModel, error)
	FindAll(row_count, offset int) (*[]AirportModel, error)
	Delete(code string) error
}

type BookingOfficeRepository interface {
	Find(id int) (*BookingOfficeModel, error)
	FindAll(row_count, offset int) (*[]BookingOfficeModel, error)
	Delete(id int) error
}

type CashierRepository interface {
	Find(login string) (*CashierModel, error)
	FindAll(row_count, offset int) (*[]CashierModel, error)
	Delete(login string) error
}

type FlightInTicketRepository interface {
	Find(depDate string, lineCode string, seatID int, ticketNo int64) (*FlightInTicketModel, error)
	FindAll(row_count, offset int) (*[]FlightInTicketModel, error)
	Delete(depDate string, lineCode string, seatID int, ticketNo int64) error
}

type FlightRepository interface {
	Find(depDate string, lineCode string) (*FlightModel, error)
	FindAll(row_count, offset int) (*[]FlightModel, error)
	Delete(depDate string, lineCode string) error
}

type LineRepository interface {
	Find(code string) (*LineModel, error)
	FindAll(row_count, offset int) (*[]LineModel, error)
	Delete(code string) error
}

type LinerModelRepository interface {
	Find(code string) (*LinerModelModel, error)
	FindAll(row_count, offset int) (*[]LinerModelModel, error)
	Delete(code string) error
}

type LinerRepository interface {
	Find(code string) (*LinerModel, error)
	FindAll(row_count, offset int) (*[]LinerModel, error)
	Delete(code string) error
}

type PurchaseRepository interface {
	Find(id int) (*PurchaseModel, error)
	FindAll(row_count, offset int) (*[]PurchaseModel, error)
	Delete(id int) error
}

type SeatRepository interface {
	Find(id int) (*SeatModel, error)
	FindAll(row_count, offset int) (*[]SeatModel, error)
	Delete(id int) error
}

type TicketRepository interface {
	Find(number int64) (*TicketModel, error)
	FindAll(row_count, offset int) (*[]TicketModel, error)
	Delete(number int64) error
}
