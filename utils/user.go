package utils

import (
	"database/sql"
	"fmt"

	"todo-auth/handler"
	"todo-auth/middlewares"

	_ "github.com/lib/pq"
)

var db *sql.DB

func SetDatabase(DB *sql.DB) {
	handler.SetDb(DB)
	middlewares.SetDB(DB)
	db = DB
	fmt.Println("DB HAS BEEN INITIALIZED IN UTILS")

}

func GetDb() *sql.DB {
	if db == nil {
		fmt.Println("it is nil ")
	}
	return db
}
