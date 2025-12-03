package store

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() error {

	// CHANGE PASSWORD HERE
	username := "root"
	password := "sairenu"
	host := "127.0.0.1"
	port := "3306"
	database := "movie_booking"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username, password, host, port, database)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	return DB.Ping()
}
