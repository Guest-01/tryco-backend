package handlers

import "github.com/Guest-01/tryco-backend/db/sqlc"

// sqlc의 Queries를 주입 받는 핸들러(컨트롤러)
type Handler struct {
	queries *sqlc.Queries
}

func New(queries *sqlc.Queries) *Handler {
	return &Handler{
		queries: queries,
	}
}
