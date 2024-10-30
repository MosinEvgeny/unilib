package controllers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/MosinEvgeny/unilib/backend/db"
	"github.com/MosinEvgeny/unilib/backend/models"
	"github.com/gin-gonic/gin"
)

func RegisterReader(c *gin.Context) {
	var reader models.Reader
	if err := c.BindJSON(&reader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validFaculties := map[string]bool{
		"Математический":      true,
		"Филологический":      true,
		"Физической культуры": true,
		"Педагогики и методики начального образования": true,
		"Педагогики и психологии детства":              true,
		"Естественнонаучный":                           true,
		"Иностранных языков":                           true,
		"Физический":                                   true,
		"Исторический":                                 true,
		"Музыки":                                       true,
		"Психологии":                                   true,
		"Информатики и экономики":                      true,
		"Правового и социально-педагогического образования": true,
		"Кафедра педагогики и психологии":                   true,
		"Отдел подготовки научно-педагогических кадров":     true,
		"Международное образование":                         true,
		"Открытый университет":                              true,
	}
	if !validFaculties[reader.Faculty] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный факультет"})
		return
	}

	// Проверка, существует ли уже пользователь с таким логином
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM readers WHERE username = $1", reader.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким логином уже существует"})
		return
	}

	// Вставка нового читателя в базу данных
	_, err = db.DB.Exec("INSERT INTO readers (full_name, faculty, course, student_id, phone_number, username, password, registration_date) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_DATE)",
		reader.FullName, reader.Faculty, reader.Course, reader.StudentID, reader.Phone_number, reader.Username, reader.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Читатель успешно зарегистрирован"})
}

func LoginReader(c *gin.Context) {
	// Получи логин и пароль из тела запроса
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Запрос к базе данных для проверки пользователя
	var reader models.Reader
	err := db.DB.QueryRow("SELECT reader_id, password, role FROM readers WHERE username = $1", credentials.Username).Scan(&reader.ReaderID, &reader.Password, &reader.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверка пароля (в реальном приложении используй хэширование!)
	if credentials.Password != reader.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный логин или пароль"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Успешная авторизация", "role": reader.Role, "reader_id": reader.ReaderID})
}

func LibrarianRegisterReader(c *gin.Context) {
	var reader models.Reader
	if err := c.BindJSON(&reader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, существует ли уже пользователь с таким логином
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM readers WHERE username = $1", reader.Username).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Пользователь с таким логином уже существует"})
		return
	}

	// Вставка нового читателя в базу данных
	_, err = db.DB.Exec("INSERT INTO readers (full_name, faculty, course, student_id, phone_number, username, password, registration_date) VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_DATE)",
		reader.FullName, reader.Faculty, reader.Course, reader.StudentID, reader.Phone_number, reader.Username, reader.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Читатель успешно зарегистрирован"})
}

func GetReaderByID(c *gin.Context) {
	readerID := c.Param("readerId")

	var reader models.Reader
	err := db.DB.QueryRow("SELECT * FROM readers WHERE reader_id = $1", readerID).Scan(
		&reader.ReaderID, &reader.FullName, &reader.Faculty, &reader.Course, &reader.StudentID, &reader.Phone_number, &reader.Username, &reader.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Читатель не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reader)
}

func GetAllReaders(c *gin.Context) {
	searchQuery := c.Query("search")
	var query string
	if searchQuery != "" {
		query = fmt.Sprintf(`
      SELECT * 
      FROM readers 
      WHERE role = 'reader' 
        AND (full_name ILIKE '%%%s%%' OR student_id ILIKE '%%%s%%')
    `, searchQuery, searchQuery)
	} else {
		query = "SELECT * FROM readers WHERE role = 'reader'"
	}

	rows, err := db.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var readers []models.Reader
	for rows.Next() {
		var reader models.Reader
		err := rows.Scan(&reader.ReaderID, &reader.FullName, &reader.Faculty, &reader.Course, &reader.StudentID, &reader.Phone_number, &reader.Username, &reader.Password, &reader.Role, &reader.RegistrationDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		readers = append(readers, reader)
	}

	c.JSON(http.StatusOK, readers)
}

func UpdateReader(c *gin.Context) {
	readerID := c.Param("readerId")

	var reader models.Reader
	if err := c.BindJSON(&reader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка, существует ли читатель с таким ID
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM readers WHERE reader_id = $1", readerID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Читатель не найден"})
		return
	}

	// Обновление данных читателя в базе данных
	_, err = db.DB.Exec(`
    UPDATE readers 
    SET full_name = $1, faculty = $2, course = $3, student_id = $4, phone_number = $5, username = $6, password = $7, role = $8 
    WHERE reader_id = $9
  `, reader.FullName, reader.Faculty, reader.Course, reader.StudentID, reader.Phone_number, reader.Username, reader.Password, reader.Role, readerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Данные читателя успешно обновлены"})
}

func DeleteReader(c *gin.Context) {
	readerID := c.Param("readerId")

	// Проверка, существует ли читатель с таким ID
	var count int
	err := db.DB.QueryRow("SELECT COUNT(*) FROM readers WHERE reader_id = $1", readerID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Читатель не найден"})
		return
	}

	// Удаление читателя
	_, err = db.DB.Exec("DELETE FROM readers WHERE reader_id = $1", readerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Читатель успешно удален"})
}

func GetReaderByStudentID(c *gin.Context) {
	studentID := c.Param("studentId")

	var reader models.Reader
	err := db.DB.QueryRow("SELECT * FROM readers WHERE student_id = $1", studentID).Scan(
		&reader.ReaderID, &reader.FullName, &reader.Faculty, &reader.Course, &reader.StudentID, &reader.Phone_number, &reader.Username, &reader.Password, &reader.Role, &reader.RegistrationDate,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Читатель не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reader)
}
