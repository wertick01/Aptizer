package main

import (
	"database/sql"
	"fmt"
	"log"

	"aptizer.com/internal/app/db"
	"aptizer.com/internal/app/processors"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database, err := OpenDB("aptizer_user:my_password@/aptizer_test?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	processor := processors.NewProcessor(
		database,
		db.NewNewsStorage(database),
		db.NewUsersStorage(database),
	)
	users, err := processor.ListUsers()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range users {
		fmt.Println(i)
	}
}

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
