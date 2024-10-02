package models

type Reader struct {
	ReaderID     int    `json:"reader_id" db:"reader_id"`
	FullName     string `json:"full_name" db:"full_name"`
	Faculty      string `json:"faculty" db:"faculty"`
	Course       int    `json:"course" db:"course"`
	StudentID    string `json:"student_id" db:"student_id"`
	Phone_number string `json:"phone_number" db:"phone_number"`
	Username     string `json:"username" db:"username"`
	Password     string `json:"password" db:"password"`
	Role         string `json:"role" db:"role"`
}
