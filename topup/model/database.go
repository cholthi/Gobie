package model

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	DB_USER = ""
	DB_PASS = ""
	DB_HOST = ""
	DB_NAME = ""
)

var DB *sqlx.DB

func InitDB(username, password, host, database string) {

	dsn := fmt.Sprintf("%s:%s@(%s:%s)/%s", username, password, host, "3306", database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic(err)
	}
	DB = db
	return
}
