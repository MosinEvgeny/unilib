package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/MosinEvgeny/unilib/backend/db"
	"github.com/MosinEvgeny/unilib/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// Report - модель для хранения информации об отчете в БД
type Report struct {
	ReportID      int       `db:"report_id"`
	ReportDate    time.Time `db:"report_date"`
	ReportContent string    `db:"report_content"`
	NewBooks      int       `json:"newBooks"`
	FormattedDate string    `json:"formatted_date"`
}

func GenerateOperationsReportData(c *gin.Context) {
	// Получаем даты начала и окончания из параметров запроса
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Парсим даты
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты начала"})
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты окончания"})
		return
	}

	//  Статистика
	var registeredReaders, issuedBooks, returnedBooks int

	//  Зарегистрированных читателей
	err = db.DB.QueryRow("SELECT COUNT(*) FROM readers WHERE registration_date BETWEEN $1 AND $2", startDate, endDate).Scan(&registeredReaders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества зарегистрированных читателей: " + err.Error()})
		return
	}

	//  Выданных книг
	err = db.DB.QueryRow("SELECT COUNT(*) FROM issue WHERE issue_date BETWEEN $1 AND $2", startDate, endDate).Scan(&issuedBooks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества выданных книг: " + err.Error()})
		return
	}

	//  Принятых книг
	err = db.DB.QueryRow("SELECT COUNT(*) FROM issue WHERE return_date BETWEEN $1 AND $2", startDate, endDate).Scan(&returnedBooks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества принятых книг: " + err.Error()})
		return
	}

	// Количество новых книг
	var newBooks int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM copies WHERE acquisition_date BETWEEN $1 AND $2", startDate, endDate).Scan(&newBooks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества новых книг: " + err.Error()})
		return
	}

	//  Получаем ФИО администратора и библиотекаря
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

	c.JSON(http.StatusOK, gin.H{
		"message":           "Данные  для  отчета  успешно  сгенерированы",
		"registeredReaders": registeredReaders,
		"issuedBooks":       issuedBooks,
		"returnedBooks":     returnedBooks,
		"adminName":         adminName,
		"librarianName":     librarianName,
		"newBooks":          newBooks,
	})
}

func GenerateOperationsReportFile(c *gin.Context) {
	//  Получаем данные из запроса
	var reportData struct {
		StartDate         string `json:"startDate"`
		EndDate           string `json:"endDate"`
		RegisteredReaders int    `json:"registeredReaders"`
		IssuedBooks       int    `json:"issuedBooks"`
		ReturnedBooks     int    `json:"returnedBooks"`
		AdminName         string `json:"adminName"`
		LibrarianName     string `json:"librarianName"`
		NewBooks          int    `json:"newBooks"`
	}
	if err := c.BindJSON(&reportData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//  Создание PDF-документа
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 15, 20)

	//  Добавляем шрифт с поддержкой кириллицы
	m.AddUTF8Font("customFont", consts.Normal, "assets/fonts/DejaVuSans.ttf")
	m.AddUTF8Font("customFontBold", consts.Bold, "assets/fonts/DejaVuSans-BoldOblique.ttf")

	//  Заголовок отчета
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("Отчет о выполненных операциях", props.Text{
				Top:    5,
				Style:  consts.Bold,
				Align:  consts.Center,
				Size:   14,
				Family: "customFontBold",
			})
		})
	})

	// Период отчета (получаем даты из reportData)
	startDate := reportData.StartDate
	endDate := reportData.EndDate

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Период: с %s по %s", startDate, endDate), props.Text{
				Align:  consts.Center,
				Family: "customFont",
			})
		})
	})

	//  Статистика
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("Статистика", props.Text{
				Top:    5,
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: "customFontBold",
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Книг  поступило:  %d", reportData.NewBooks), props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Зарегистрированных  читателей:  %d", reportData.RegisteredReaders), props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Книг  выдано  для  чтения:  %d", reportData.IssuedBooks), props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Книг  возвращено:  %d", reportData.ReturnedBooks), props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})

	// Ответственные лица
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Ответственные лица:", props.Text{
				Top:    5,
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: "customFontBold",
			})
		})
	})

	m.Row(40, func() {
		m.Col(7, func() {
			m.Text("Библиотекарь:", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text(reportData.LibrarianName, props.Text{
				Top:    5,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Подпись: _____________________", props.Text{
				Top:    15,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Дата: _____________________", props.Text{
				Top:    25,
				Align:  consts.Left,
				Family: "customFont",
			})
		})
		m.Col(7, func() {
			m.Text("Администратор:", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text(reportData.AdminName, props.Text{
				Top:    5,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Подпись: _____________________", props.Text{
				Top:    15,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Дата: _____________________", props.Text{
				Top:    25,
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})

	//  Сохранение PDF-документа
	reportFileName := fmt.Sprintf("reports/operations_report_%s.pdf", time.Now().Format("2006-01-02_15-04-05"))
	err := m.OutputFileAndClose(reportFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании файла отчета: " + err.Error()})
		return
	}

	//  Сохранение информации об отчете в БД
	reportContent := fmt.Sprintf("Зарегистрированных читателей: %d\nКниг выдано для чтения: %d\nКниг возвращено: %d\nНовых книг: %d", reportData.RegisteredReaders, reportData.IssuedBooks, reportData.ReturnedBooks, reportData.NewBooks)
	_, err = db.DB.Exec("INSERT INTO reports (report_date, report_content) VALUES (CURRENT_DATE, $1)", reportContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении отчета в базу данных: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отчет успешно создан", "filename": reportFileName})
}

func GetAllReports(c *gin.Context) {
	rows, err := db.DB.Query("SELECT * FROM reports ORDER BY report_date DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var reports []Report
	for rows.Next() {
		var report Report
		err := rows.Scan(&report.ReportID, &report.ReportDate, &report.ReportContent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		report.FormattedDate = report.ReportDate.Format("02.01.2006")
		reports = append(reports, report)
	}

	c.JSON(http.StatusOK, reports)
}

func DeleteReport(c *gin.Context) {
	reportID := c.Param("reportId")

	//  Удаление отчета из базы данных
	_, err := db.DB.Exec("DELETE FROM reports WHERE report_id = $1", reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Отчет успешно удален"})
}

func GenerateRemovalAct(c *gin.Context) {
	var requestData struct {
		InventoryNumbers []string `json:"inventory_numbers"`
	}

	if err := c.BindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 15, 20)

	//  Добавляем шрифт с поддержкой кириллицы
	m.AddUTF8Font("customFont", consts.Normal, "assets/fonts/DejaVuSans.ttf")
	m.AddUTF8Font("customFontBold", consts.Bold, "assets/fonts/DejaVuSans-BoldOblique.ttf")

	m.Row(16, func() {
		m.Col(12, func() {
			m.Text("ФОРМА АКТА О СПИСАНИИ КНИГ ИЗ БИБЛИОТЕКИ", props.Text{
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: "customFontBold",
				Size:   12,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Акт N ___", props.Text{
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: "customFontBold",
				Size:   12,
			})
		})
	})

	currentDate := time.Now().Format("02.01.2006")
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Настоящий акт составлен \"%s\"", currentDate), props.Text{
				Align:  consts.Left,
				Family: "customFont",
				Size:   10,
			})
		})
	})

	// Получаем ФИО администратора и библиотекаря для акта
	var adminName, librarianName string
	err := db.DB.QueryRow("SELECT full_name FROM readers WHERE role = 'admin' LIMIT 1").Scan(&adminName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении ФИО администратора: " + err.Error()})
		return
	}
	err = db.DB.QueryRow("SELECT full_name FROM readers WHERE role = 'librarian' LIMIT 1").Scan(&librarianName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении ФИО библиотекаря: " + err.Error()})
		return
	}

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("о списании из библиотеки книг в количестве _____________ экземпляров", props.Text{
				Align:  consts.Left,
				Family: "customFont",
				Size:   10,
			})
		})
	})

	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("с причиной __________________________________________________________________", props.Text{
				Align:  consts.Left,
				Family: "customFont",
				Size:   10,
			})
		})
	})

	m.Row(1, func() {
		m.Col(1, func() {
			m.TableList(
				[]string{"№", "Автор и заглавие книги", "Инвентарный номер"},
				func() [][]string { //  Изменено:  возвращаем [][]string
					var booksData [][]string
					for i, inventoryNumber := range requestData.InventoryNumbers {
						var book models.Book
						err := db.DB.QueryRow(`SELECT b.title, b.author FROM books AS b JOIN copies AS c ON b.book_id = c.book_id WHERE c.inventory_number = $1`, inventoryNumber).Scan(&book.Title, &book.Author)
						if err != nil {
							fmt.Println("Ошибка при получении данных книги:", err)
							continue
						}
						booksData = append(booksData, []string{fmt.Sprintf("%d", i+1), fmt.Sprintf("%s - %s", book.Author, book.Title), inventoryNumber})

						_, err = db.DB.Exec("UPDATE copies SET status = 'Списан' WHERE inventory_number = $1", inventoryNumber)
						if err != nil {
							fmt.Println("Ошибка при обновлении статуса экземпляра:", err)
							continue
						}
					}
					return booksData // Возвращаем booksData
				}(),
				props.TableList{
					Align:       consts.Left,
					HeaderProp:  props.TableListContent{Size: 10, GridSizes: []uint{1, 8, 3}, Style: consts.Bold, Family: "customFontBold"},
					ContentProp: props.TableListContent{Size: 8, GridSizes: []uint{1, 8, 3}, Family: "customFont"},
				},
			)
		})
	})

	// Ответственные лица
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text("Ответственные лица:", props.Text{
				Top:    5,
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: "customFontBold",
			})
		})
	})

	m.Row(40, func() {
		m.Col(7, func() {
			m.Text("Библиотекарь:", props.Text{
				Top:    5,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text(librarianName, props.Text{
				Top:    11,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Подпись: _____________________", props.Text{
				Top:    20,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Дата: _____________________", props.Text{
				Top:    30,
				Align:  consts.Left,
				Family: "customFont",
			})
		})
		m.Col(7, func() {
			m.Text("Администратор:", props.Text{
				Top:    5,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text(adminName, props.Text{
				Top:    11,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Подпись: _____________________", props.Text{
				Top:    20,
				Align:  consts.Left,
				Family: "customFont",
			})
			m.Text("Дата: _____________________", props.Text{
				Top:    30,
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})

	reportFileName := fmt.Sprintf("reports/removal_act_%s.pdf", time.Now().Format("2006-01-02_15-04-05"))
	err = m.OutputFileAndClose(reportFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании PDF файла: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Акт о списании успешно создан", "filename": reportFileName})
}
