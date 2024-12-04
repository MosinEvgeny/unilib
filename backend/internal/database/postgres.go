package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB() {
	connStr := "user=postgres password=pmi-gotosw dbname=unilib sslmode=disable"
	var err error
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Успешное  подключение  к  базе  данных!")
}
