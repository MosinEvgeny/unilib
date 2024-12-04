package models

type Book struct {
	BookID          int    `json:"book_id" db:"book_id"`
	Title           string `json:"title" db:"title"`
	Author          string `json:"author" db:"author"`
	ISBN            string `json:"isbn" db:"isbn"`
	Publisher       string `json:"publisher" db:"publisher"`
	PublicationYear int    `json:"publication_year" db:"publication_year"`
	TotalCopies     int    `json:"total_copies" db:"total_copies"`
	Category        string `json:"category" db:"category"`
	Description     string `json:"description" db:"description"`
	AvailableCopies int    `json:"available_copies" db:"available_copies"` // Добавим поле для кол-ва доступных экземпляров
}
