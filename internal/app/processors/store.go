package processors

import (
	"database/sql"

	database "aptizer.com/internal/app/db"
)

type Processor struct {
	store *database.Store
}

func NewProcessor(
	db *sql.DB,
	newsStorage *database.NewsStorage,
	usersStorage *database.UsersStorage,
) *Processor {
	return &Processor{
		store: database.New(
			db,
			newsStorage,
			usersStorage,
		),
	}
}
