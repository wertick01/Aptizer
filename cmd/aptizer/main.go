package main

import (
	"database/sql"
	"fmt"
	"log"

	"aptizer.com/internal/app/db"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	database, err := OpenDB("aptizer_user:my_password@/aptizer_test?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	newsStorage := db.NewNewsStorage(database)
	usersStorage := db.NewUsersStorage(database)
	store := db.New(
		database,
		newsStorage,
		usersStorage,
	)
	news, err := store.NewsStorage.List()
	fmt.Println(news)
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
