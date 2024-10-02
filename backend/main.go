package main

import (
	"github.com/MosinEvgeny/unilib/backend/db"
	"github.com/MosinEvgeny/unilib/backend/router"
)

func main() {
	db.InitDB()
	defer db.DB.Close()

	router := router.SetupRouter()
	router.Run(":8080")
}
