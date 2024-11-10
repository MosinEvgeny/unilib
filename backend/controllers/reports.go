package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/MosinEvgeny/unilib/backend/db"
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
	//  Получаем  дату  месяц  назад
	aMonthAgo := time.Now().AddDate(0, -1, 0)

	//  Статистика
	var registeredReaders, issuedBooks, returnedBooks int

	//  Зарегистрированных читателей
	err := db.DB.QueryRow("SELECT COUNT(*) FROM readers WHERE registration_date >= $1", aMonthAgo).Scan(&registeredReaders)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества зарегистрированных читателей: " + err.Error()})
		return
	}

	//  Выданных книг
	err = db.DB.QueryRow("SELECT COUNT(*) FROM issue WHERE issue_date >= $1", aMonthAgo).Scan(&issuedBooks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества выданных книг: " + err.Error()})
		return
	}

	//  Принятых книг
	err = db.DB.QueryRow("SELECT COUNT(*) FROM issue WHERE return_date >= $1", aMonthAgo).Scan(&returnedBooks)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении количества принятых книг: " + err.Error()})
		return
	}

	// Количество новых книг
	var newBooks int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM copies WHERE acquisition_date >= $1", aMonthAgo).Scan(&newBooks)
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

	//  Период отчета
	aMonthAgo := time.Now().AddDate(0, -1, 0).Format("02.01.2006")
	currentDate := time.Now().Format("02.01.2006")
	m.Row(10, func() {
		m.Col(12, func() {
			m.Text(fmt.Sprintf("Период:  с  %s  по  %s", aMonthAgo, currentDate), props.Text{
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

	//  Ответственные лица
	m.Row(20, func() {
		m.Col(12, func() {
			m.Text("Ответственные  лица", props.Text{
				Top:    5,
				Style:  consts.Bold,
				Align:  consts.Center,
				Family: "customFontBold",
			})
		})
	})

	m.Row(10, func() {
		m.Col(6, func() {
			m.Text("Библиотекарь: ", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
		m.Col(6, func() {
			m.Text("Администратор: ", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})

	m.Row(10, func() {
		m.Col(6, func() {
			m.Text(reportData.LibrarianName, props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
		m.Col(6, func() {
			m.Text(reportData.AdminName, props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})

	m.Row(10, func() {
		m.Col(6, func() {
			m.Text("Дата: _____________________", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
		m.Col(6, func() {
			m.Text("Дата: _____________________", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
	})

	m.Row(10, func() {
		m.Col(6, func() {
			m.Text("Подпись: _____________________", props.Text{
				Align:  consts.Left,
				Family: "customFont",
			})
		})
		m.Col(6, func() {
			m.Text("Подпись: _____________________", props.Text{
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
