package handlers

import (
	db "github.com/haarshmap/go-url/cmd/db/generated"
)

type Handler struct {
	queries *db.Queries
}

func NewHandler(q *db.Queries) *Handler {
	return &Handler{
		queries: q,
	}
}

// handlers
