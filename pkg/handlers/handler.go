package handlers

import (
	"github.com/MobasirSarkar/BeTask/database"
	"github.com/MobasirSarkar/BeTask/pkg/auth"
)

// created this handler for better approach

type Handler struct {
	db   *database.Postgres
	auth *auth.AuthService
}

func New(db *database.Postgres, auth *auth.AuthService) *Handler {
	return &Handler{
		db:   db,
		auth: auth,
	}
}
