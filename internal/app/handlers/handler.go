package handlers

import (
	"database/sql"

	database "aptizer.com/internal/app/db"
	"aptizer.com/internal/app/handlers/authorization"
	"aptizer.com/internal/app/processors"
)

type Handler struct {
	Processor  *processors.Processor
	Authorizer *authorization.Authoriser
}

func NewHandler(
	db *sql.DB,
	newsStorage *database.NewsStorage,
	usersStorage *database.UsersStorage,
) *Handler {
	processor := processors.NewProcessor(
		db,
		newsStorage,
		usersStorage,
	)
	return &Handler{
		Processor: processor,
		Authorizer: authorization.NewAuthoriser(
			processor,
		),
	}
}
