package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/MosinEvgeny/unilib/backend/internal/database"
	"github.com/MosinEvgeny/unilib/backend/internal/models"
	"github.com/gin-gonic/gin"
	_ "github.com/jmoiron/sqlx"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"net/http"
	"strconv"
)

func CreateContract(c *gin.Context) {
	var contractData struct {
		AdminFullName     string                `json:"admin_full_name"`
		LibrarianFullName string                `json:"librarian_full_name"`
		TotalBooks        int                   `json:"total_books"`
		TotalSum          float64               `json:"total_sum"`
		Books             []models.ContractBook `json:"books"`
	}

	if err := c.BindJSON(&contractData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Вставка нового контракта в базу данных
	var contractID int
	err := db.DB.QueryRow(`
        INSERT INTO contracts (admin_full_name, librarian_full_name, total_books, total_sum) 
        VALUES ($1, $2, $3, $4) RETURNING contract_id
    `, contractData.AdminFullName, contractData.LibrarianFullName, contractData.TotalBooks, contractData.TotalSum).Scan(&contractID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Вставка данных о книгах в контракт
	for _, book := range contractData.Books {
		_, err := db.DB.Exec(`INSERT INTO contract_books (contract_id, book_id, title, author, price, copies, sum)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			contractID, book.BookID, book.Title, book.Author, book.Price, book.Copies, book.Sum)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Ошибка при добавлении книги в договор: %s", err.Error())})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Договор успешно создан",
		"contract_id": contractID,
	})
}

// GenerateContractPDF генерирует PDF-файл договора.
func GenerateContractPDF(c *gin.Context) {
	contractIDStr := c.Param("contractId")
	contractID, err := strconv.Atoi(contractIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID договора"})
		return
	}

	// Получение данных договора из базы данных
	var contract models.Contract
	err = db.DB.QueryRowx("SELECT * FROM contracts WHERE contract_id = $1", contractID).StructScan(&contract)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Договор не найден"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных договора: " + err.Error()})
		return
	}

	// Получение ФИО администратора и библиотекаря
	var adminName, librarianName string
	err = db.DB.QueryRow("SELECT full_name FROM readers WHERE role = 'admin' LIMIT 1").Scan(&adminName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении ФИО администратора: " + err.Error()})
		return
	}
	err = db.DB.QueryRow("SELECT full_name FROM readers WHERE role = 'librarian' LIMIT 1").Scan(&librarianName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении ФИО библиотекаря: " + err.Error()})
		return
	}

	contract.AdminFullName = adminName
	contract.LibrarianFullName = librarianName

	// Создание PDF-документа
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 15, 20)

	// Добавление шрифтов (без проверки ошибки)
	m.AddUTF8Font("customFont", consts.Normal, "assets/fonts/DejaVuSans.ttf")
	m.AddUTF8Font("customFontBold", consts.Bold, "assets/fonts/DejaVuSans-Bold.ttf")

	// Формирование заголовка акта (используем Family вместо SetFont)
	m.Text("ФОРМА АКТА О ПРИЕМЕ ДОКУМЕНТОВ В БИБЛИОТЕКУ", props.Text{Align: consts.Center, Family: "customFontBold", Size: 14})
	m.Text(fmt.Sprintf("Акт N %d", contract.ContractID), props.Text{Align: consts.Center, Family: "customFontBold", Size: 14})

	m.Text(fmt.Sprintf("Настоящий акт составлен %s", contract.CreationDate.Format("02.01.2006")), props.Text{Align: consts.Left, Family: "customFont", Size: 10})
	m.Text(fmt.Sprintf("%s и %s", contract.LibrarianFullName, contract.AdminFullName), props.Text{Align: consts.Left, Family: "customFont", Size: 10})
	m.Text(fmt.Sprintf("о приеме в библиотеку книг в количестве %d экземпляров на общую сумму %.2f", contract.TotalBooks, contract.TotalSum), props.Text{Align: consts.Left, Family: "customFont", Size: 10})

	//  Список книг
	m.RegisterHeader(func() {
		m.Text("Список книг:", props.Text{Align: consts.Left, Family: "customFont", Top: 10, Size: 10})
	})

	// Получение списка книг, связанных с договором
	var contractBooks []models.ContractBook

	query := `SELECT books.book_id, books.title, books.author, $1 AS price, $2 AS copies, $3 AS sum
            FROM books
            WHERE books.book_id = $4`

	for _, book := range contract.Books {
		err = db.DB.Select(&contractBooks, query, book.Price, book.Copies, book.Sum, book.BookID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка книг для договора: " + err.Error()})
			return
		}
	}

	for _, book := range contractBooks {
		m.Text(fmt.Sprintf("- %s - %s", book.Author, book.Title), props.Text{Align: consts.Left, Family: "customFont"})
	}

	// Подписи
	m.Row(40, func() {
		m.Col(6, func() {
			m.Text("Библиотекарь:", props.Text{Align: consts.Left, Family: "customFont"})
			m.Text(contract.LibrarianFullName, props.Text{Align: consts.Left, Family: "customFont"})
			m.Text("Подпись: _______", props.Text{Top: 10, Align: consts.Left, Family: "customFont"})
			m.Text("Дата: _______", props.Text{Align: consts.Left, Family: "customFont"})
		})
		m.Col(6, func() {
			m.Text("Администратор:", props.Text{Align: consts.Left, Family: "customFont"})
			m.Text(contract.AdminFullName, props.Text{Align: consts.Left, Family: "customFont"})
			m.Text("Подпись: _______", props.Text{Top: 10, Align: consts.Left, Family: "customFont"})
			m.Text("Дата: _______", props.Text{Align: consts.Left, Family: "customFont"})
		})
	})

	// Сохранение PDF-документа
	reportFileName := fmt.Sprintf("reports/contract_%d.pdf", contractID)
	err = m.OutputFileAndClose(reportFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании PDF-файла договора: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "PDF-файл договора успешно создан", "filename": reportFileName})
}
