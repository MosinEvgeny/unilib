package models

import "time"

type Issue struct {
	IssueID    int       `json:"issue_id" db:"issue_id"`
	CopyID     int       `json:"copy_id" db:"copy_id"`
	ReaderID   int       `json:"reader_id" db:"reader_id"`
	IssueDate  time.Time `json:"issue_date" db:"issue_date"`
	DueDate    time.Time `json:"due_date" db:"due_date"`
	ReturnDate time.Time `json:"return_date" db:"return_date"`
}
