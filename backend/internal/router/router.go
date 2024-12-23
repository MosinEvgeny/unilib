package router

import (
	"github.com/MosinEvgeny/unilib/backend/internal/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.Default())

	// Маршруты для контрактов
	router.POST("/contracts", controllers.CreateContract)
	router.GET("/contracts/:contractId/pdf", controllers.GenerateContractPDF)
	router.GET("/get-admin", controllers.GetAdmin)
	router.GET("/get-librarian", controllers.GetLibrarian)

	// Маршруты для книг
	router.GET("/books", controllers.GetAllBooks)
	router.GET("/books/:bookId/available", controllers.GetBookAvailableCopies)
	router.POST("/books", controllers.CreateBooks)
	router.PUT("/books/:bookId", controllers.UpdateBook)
	router.DELETE("/books/:bookId", controllers.DeleteBook)

	// Маршруты для читателей
	router.POST("/register", controllers.RegisterReader)
	router.POST("/login", controllers.LoginReader)
	router.GET("/readers/:readerId", controllers.GetReaderByID)
	router.GET("/readers", controllers.GetAllReaders)
	router.PUT("/readers/:readerId", controllers.UpdateReader)
	router.DELETE("/readers/:readerId", controllers.DeleteReader)
	router.GET("/readers/by-student-id/:studentId", controllers.GetReaderByStudentID)

	// Маршруты для заказов
	router.GET("/orders/:readerId", controllers.GetReaderOrders)
	router.POST("/issue", controllers.IssueBook)
	router.GET("/issues/by-student-id/:studentId", controllers.GetIssuesByStudentID)
	router.GET("/reader/:studentId/issues", controllers.GetReaderIssues)

	// Ещё для чего-то
	router.GET("/issue/:issueId", controllers.GetIssueByID)
	router.POST("/librarian/register-reader", controllers.LibrarianRegisterReader)
	router.PUT("/issue/:issueId/return", controllers.ReturnBook)

	// Маршруты для отчётов
	router.GET("/reports/operations", controllers.GenerateOperationsReportData)
	router.POST("/reports/operations/generate", controllers.GenerateOperationsReportFile)
	router.GET("/reports", controllers.GetAllReports)
	router.DELETE("/reports/:reportId", controllers.DeleteReport)
	router.POST("/removal-act", controllers.GenerateRemovalAct)

	return router
}
