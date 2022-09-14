package db

import "database/sql"

type Store struct {
	db           *sql.DB
	NewsStorage  *NewsStorage
	UsersStorage *UsersStorage
}

func New(
	db *sql.DB,
	newsStorage *NewsStorage,
	usersStorage *UsersStorage,
) *Store {
	return &Store{
		db:           db,
		NewsStorage:  newsStorage,
		UsersStorage: usersStorage,
	}
}

// NewNewsStorage - Creating new copy of NewsStorage struct.
func NewNewsStorage(db *sql.DB) *NewsStorage {
	storage := new(NewsStorage)
	storage.database = db
	return storage
}

// NewUsersStorage - Creating new copy of UsersStorage struct.
func NewUsersStorage(db *sql.DB) *UsersStorage {
	storage := new(UsersStorage)
	storage.DB = db
	return storage
}
