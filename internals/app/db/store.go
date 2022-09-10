package db

import "database/sql"

type Store struct {
	db          *sql.DB
	NewsStorage *NewsStorage
}

func New(db *sql.DB, newsStorage *NewsStorage) *Store {
	return &Store{
		db:          db,
		NewsStorage: newsStorage,
	}
}
