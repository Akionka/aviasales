package store

type AirportRepository interface {
	Create(*AirportModel) error
	Find(code string) (*AirportModel, error)
	FindAll(row_count, offset int) (*[]AirportModel, error)
	Delete(code string) error
}

type BookingOfficeRepository interface {
	Create(*BookingOfficeModel) error
	Find(id int) (*BookingOfficeModel, error)
	FindAll(row_count, offset int) (*[]BookingOfficeModel, error)
	Delete(id int) error
}

type CashierRepository interface {
	Create(*CashierModel) error
	Find(login string) (*CashierModel, error)
	FindAll(row_count, offset int) (*[]CashierModel, error)
	Delete(login string) error
}

type FlightInTicketRepository interface {
	Create(*FlightInTicketModel) error
	Find(depDate string, lineCode string, seatID int, ticketNo int64) (*FlightInTicketModel, error)
	FindAll(row_count, offset int) (*[]FlightInTicketModel, error)
	Delete(depDate string, lineCode string, seatID int, ticketNo int64) error
}

type FlightRepository interface {
	Create(*FlightModel) error
	Find(depDate string, lineCode string) (*FlightModel, error)
	FindAll(row_count, offset int) (*[]FlightModel, error)
	Delete(depDate string, lineCode string) error
}

type LineRepository interface {
	Create(*LineModel) error
	Find(code string) (*LineModel, error)
	FindAll(row_count, offset int) (*[]LineModel, error)
	Delete(code string) error
}

type LinerModelRepository interface {
	Create(*LinerModelModel) error
	Find(code string) (*LinerModelModel, error)
	FindAll(row_count, offset int) (*[]LinerModelModel, error)
	Delete(code string) error
}

type LinerRepository interface {
	Create(*LinerModel) error
	Find(code string) (*LinerModel, error)
	FindAll(row_count, offset int) (*[]LinerModel, error)
	Delete(code string) error
}

type PurchaseRepository interface {
	Create(*PurchaseModel) error
	Find(id int) (*PurchaseModel, error)
	FindAll(row_count, offset int) (*[]PurchaseModel, error)
	Delete(id int) error
}

type SeatRepository interface {
	Create(*SeatModel) error
	Find(id int) (*SeatModel, error)
	FindAll(row_count, offset int) (*[]SeatModel, error)
	Delete(id int) error
}

type TicketRepository interface {
	Create(*TicketModel) error
	Find(number int64) (*TicketModel, error)
	FindAll(row_count, offset int) (*[]TicketModel, error)
	Delete(number int64) error
}
