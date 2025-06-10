package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func DbConnect() *sql.DB {
	if db != nil {
		return db
	}

	var err error
	db, err = sql.Open("mysql", "root:@tcp(127.0.0.1:3307)/ordermatch?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Connected Successfully")

	return db
}
