package controllers

import (
	"database/sql"
	"net/http"

	"github.com/MosinEvgeny/unilib/backend/db"
	"github.com/MosinEvgeny/unilib/backend/models"
	"github.com/gin-gonic/gin"
)

func GetReaderOrders(c *gin.Context) {
	readerID := c.Param("readerId")

	rows, err := db.DB.Query(`
      SELECT b.title, i.issue_date, i.due_date, i.return_date
      FROM issue AS i
      JOIN copies AS c ON i.copy_id = c.copy_id
      JOIN books AS b ON c.book_id = b.book_id
      WHERE i.reader_id = $1
    `, readerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var orders []gin.H
	for rows.Next() {
		var title string
		var issueDate, dueDate, returnDate sql.NullTime
		err = rows.Scan(&title, &issueDate, &dueDate, &returnDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		order := gin.H{
			"title":      title,
			"issue_date": issueDate.Time.Format("2006-01-02"),
			"due_date":   dueDate.Time.Format("2006-01-02"),
		}
		if returnDate.Valid {
			order["return_date"] = returnDate.Time.Format("2006-01-02")
		}

		orders = append(orders, order)
	}

	c.JSON(http.StatusOK, orders)
}

func IssueBook(c *gin.Context) {
	var issueData struct {
		ReaderID  int    `json:"reader_id"`
		BookTitle string `json:"book_title"`
	}
	if err := c.BindJSON(&issueData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bookID int
	err := db.DB.QueryRow("SELECT book_id FROM books WHERE title = $1", issueData.BookTitle).Scan(&bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Книга не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверка доступности книги
	var availableCopies int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM copies WHERE book_id = $1 AND status = 'Доступен'", bookID).Scan(&availableCopies)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if availableCopies == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нет доступных экземпляров книги"})
		return
	}

	// Получение первого доступного экземпляра
	var copyID int
	err = db.DB.QueryRow("SELECT copy_id FROM copies WHERE book_id = $1 AND status = 'Доступен' LIMIT 1", bookID).Scan(&copyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Создание записи о выдаче
	_, err = db.DB.Exec("INSERT INTO issue (copy_id, reader_id, issue_date, due_date) VALUES ($1, $2, CURRENT_DATE, CURRENT_DATE + INTERVAL '14 days')", copyID, issueData.ReaderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновление статуса экземпляра
	_, err = db.DB.Exec("UPDATE copies SET status = 'Выдан' WHERE copy_id = $1", copyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Книга успешно выдана"})
}

func GetIssueByID(c *gin.Context) {
	issueID := c.Param("issueId")

	var issue models.Issue
	err := db.DB.QueryRow("SELECT * FROM issue WHERE issue_id = $1", issueID).Scan(
		&issue.IssueID, &issue.CopyID, &issue.ReaderID, &issue.IssueDate, &issue.DueDate, &issue.ReturnDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Выдача не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, issue)
}

func ReturnBook(c *gin.Context) {
	issueID := c.Param("issueId")

	//  Обновление  записи  о  выдаче
	_, err := db.DB.Exec("UPDATE issue SET return_date = CURRENT_DATE WHERE issue_id = $1 AND return_date IS NULL", issueID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//  Получение  ID  экземпляра  из  записи  о  выдаче
	var copyID int
	err = db.DB.QueryRow("SELECT copy_id FROM issue WHERE issue_id = $1", issueID).Scan(&copyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Обновление статуса экземпляра
	_, err = db.DB.Exec("UPDATE copies SET status = 'Доступен' WHERE copy_id = $1", copyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Книга успешно принята"})
}

func GetIssuesByStudentID(c *gin.Context) {
	studentID := c.Param("studentId")

	rows, err := db.DB.Query(`
        SELECT i.issue_id, b.title, i.issue_date, i.due_date, i.return_date
        FROM issue AS i
        JOIN copies AS c ON i.copy_id = c.copy_id
        JOIN books AS b ON c.book_id = b.book_id
        JOIN readers AS r ON i.reader_id = r.reader_id
        WHERE r.student_id = $1
    `, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var issues []gin.H
	for rows.Next() {
		var issueID int
		var bookTitle string
		var issueDate, dueDate, returnDate sql.NullTime
		err = rows.Scan(&issueID, &bookTitle, &issueDate, &dueDate, &returnDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		issue := gin.H{
			"issue_id":   issueID,
			"book_title": bookTitle,
			"issue_date": issueDate.Time.Format("2006-01-02"),
			"due_date":   dueDate.Time.Format("2006-01-02"),
		}
		if returnDate.Valid {
			issue["return_date"] = returnDate.Time.Format("2006-01-02")
		}

		issues = append(issues, issue)
	}

	c.JSON(http.StatusOK, issues)
}

func GetReaderIssues(c *gin.Context) {
	studentID := c.Param("studentId")

	rows, err := db.DB.Query(`
	  SELECT i.issue_id, b.title, i.issue_date, i.due_date, i.return_date
	  FROM issue AS i
	  JOIN copies AS c ON i.copy_id = c.copy_id
	  JOIN books AS b ON c.book_id = b.book_id
	  JOIN readers AS r ON i.reader_id = r.reader_id
	  WHERE r.student_id = $1
	  ORDER BY i.issue_date DESC -- Сортировка по дате выдачи (по убыванию)
	`, studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var issues []gin.H
	for rows.Next() {
		var issueID int
		var bookTitle string
		var issueDate, dueDate, returnDate sql.NullTime
		err = rows.Scan(&issueID, &bookTitle, &issueDate, &dueDate, &returnDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		issue := gin.H{
			"issue_id":   issueID,
			"book_title": bookTitle,
			"issue_date": issueDate.Time.Format("2006-01-02"),
			"due_date":   dueDate.Time.Format("2006-01-02"),
		}
		if returnDate.Valid {
			issue["return_date"] = returnDate.Time.Format("2006-01-02")
		}

		issues = append(issues, issue)
	}

	c.JSON(http.StatusOK, issues)
}
