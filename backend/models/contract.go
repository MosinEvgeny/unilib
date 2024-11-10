package models

import "time"

type ContractBook struct {
	ContractBookID int     `db:"contract_book_id"`
	ContractID     int     `db:"contract_id"`
	BookID         int     `db:"book_id"`
	Title          string  `db:"title"`
	Author         string  `db:"author"`
	Price          float64 `db:"price"`
	Copies         int     `db:"copies"`
	Sum            float64 `db:"sum"`
}

type Contract struct {
	ContractID        int            `json:"contract_id" db:"contract_id"`
	CreationDate      time.Time      `json:"creation_date" db:"creation_date"`
	AdminFullName     string         `json:"admin_full_name" db:"admin_full_name"`
	LibrarianFullName string         `json:"librarian_full_name" db:"librarian_full_name"`
	TotalBooks        int            `json:"total_books" db:"total_books"`
	TotalSum          float64        `json:"total_sum" db:"total_sum"`
	Books             []ContractBook `json:"books"` //  Убрали  db:"books",  так  как  это  поле  не  в  базе
}

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}
