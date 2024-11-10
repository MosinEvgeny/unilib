package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB //  Изменено: *sqlx.DB

func InitDB() {
	connStr := "user=postgres password=pmi-gotosw dbname=unilib sslmode=disable" //  Подключись  к  своей  БД
	var err error
	DB, err = sqlx.Open("postgres", connStr) //  Изменено: sqlx.Open
	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Успешное  подключение  к  базе  данных!")
}
