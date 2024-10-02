package controllers

import (
	"fmt"
	"net/http"

	"github.com/MosinEvgeny/unilib/backend/db"     // Замени your-username
	"github.com/MosinEvgeny/unilib/backend/models" // Замени your-username
	"github.com/gin-gonic/gin"
)

func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, существует ли уже книга с таким ISBN (если ISBN указан)
	if book.ISBN != "" {
		var count int
		err := db.DB.QueryRow("SELECT COUNT(*) FROM books WHERE isbn = $1", book.ISBN).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Книга с таким ISBN уже существует"})
			return
		}
	}

	// Вставка новой книги в базу данных
	err := db.DB.QueryRow("INSERT INTO books (title, author, isbn, publisher, publication_year, total_copies, category, description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING book_id",
		book.Title, book.Author, book.ISBN, book.Publisher, book.PublicationYear, book.TotalCopies, book.Category, book.Description).Scan(&book.BookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Создание  записей  о  копиях  книги
	for i := 0; i < book.TotalCopies; i++ {
		_, err = db.DB.Exec("INSERT INTO copies (book_id, inventory_number, status, acquisition_date) VALUES ($1, $2, 'Доступен', CURRENT_DATE)", book.BookID, fmt.Sprintf("%d-%d", book.BookID, i+1))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Книга успешно добавлена", "book_id": book.BookID})
}

func GetAllBooks(c *gin.Context) {
	searchQuery := c.Query("search")

	var query string
	if searchQuery != "" {
		query = fmt.Sprintf("SELECT * FROM books WHERE title ILIKE '%%%s%%' OR author ILIKE '%%%s%%'", searchQuery, searchQuery)
	} else {
		query = "SELECT * FROM books"
	}

	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.BookID, &book.Title, &book.Author, &book.ISBN, &book.Publisher, &book.PublicationYear, &book.TotalCopies, &book.Category, &book.Description)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		books = append(books, book)
	}

	// Получение количества доступных экземпляров для каждой книги
	for i := range books {
		err = db.DB.QueryRow("SELECT COUNT(*) FROM copies WHERE book_id = $1 AND status = 'Доступен'", books[i].BookID).Scan(&books[i].AvailableCopies)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, books)
}

func GetBookAvailableCopies(c *gin.Context) {
	bookID := c.Param("bookId")

	var availableCopies int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM copies WHERE book_id = $1 AND status = 'Доступен'", bookID).Scan(&availableCopies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"available_copies": availableCopies})
}

func UpdateBook(c *gin.Context) {
	bookID := c.Param("bookId")

	var book models.Book
	if err := c.BindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, существует ли книга с таким ID
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM books WHERE book_id = $1", bookID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
		return
	}

	// Проверка, существует ли уже другая книга с таким ISBN (если ISBN изменен)
	if book.ISBN != "" {
		var isbnCount int
		err = db.DB.QueryRow("SELECT COUNT(*) FROM books WHERE isbn = $1 AND book_id != $2", book.ISBN, bookID).Scan(&isbnCount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if isbnCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Книга с таким ISBN уже существует"})
			return
		}
	}

	// Обновление информации о книге в базе данных
	_, err = db.DB.Exec(`
	  UPDATE books 
	  SET title = $1, author = $2, isbn = $3, publisher = $4, publication_year = $5, total_copies = $6, category = $7, description = $8 
	  WHERE book_id = $9
	`, book.Title, book.Author, book.ISBN, book.Publisher, book.PublicationYear, book.TotalCopies, book.Category, book.Description, bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Информация о книге успешно обновлена"})
}

func DeleteBook(c *gin.Context) {
	bookID := c.Param("bookId")

	// Проверка, существует ли книга с таким ID
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM books WHERE book_id = $1", bookID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
		return
	}

	// Удаление всех копий книги (ON DELETE CASCADE в БД)
	_, err = db.DB.Exec("DELETE FROM copies WHERE book_id = $1", bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Удаление книги
	_, err = db.DB.Exec("DELETE FROM books WHERE book_id = $1", bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно списана"})
}
