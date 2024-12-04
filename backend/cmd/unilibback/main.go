package main

import (
	"github.com/MosinEvgeny/unilib/backend/internal/database"
	"github.com/MosinEvgeny/unilib/backend/internal/router"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := router.SetupRouter()
	router.Run(":8080")
}
