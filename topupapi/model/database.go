package model

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

var DB *pop.Connection

func init() {
	initDB()
}

func initDB() {
	dsn := os.Getenv("DSN_STRING")

	if dsn == "" {
		dsn = "root:dm4rk88@tcp(127.0.0.1:3306)/api?parseTime=true&sql_mode=TRADITIONAL&multiStatements=true"
	}
	db, err := pop.NewConnection(&pop.ConnectionDetails{
		Dialect: "mysql",
		URL:     dsn,
	})

	if err != nil {
		panic(err)
	}

	if err := db.Open(); err != nil {
		panic(errors.Wrap(err, "checking database connection"))
	}

	DB = db
	return
}
