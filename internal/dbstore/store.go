package dbstore

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func GetDB() *sqlx.DB {
	db, err := sqlx.Connect("mysql", "root:@(localhost:3306)/konachan")
	if err != nil {
		log.Fatalln(err)
	}

	return db
}
