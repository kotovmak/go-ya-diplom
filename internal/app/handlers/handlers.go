package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ctxKey int8

type store interface {
}

type Handler struct {
	store   store
	baseURL string
}

const (
	ctxKeyUserID ctxKey = iota
)

func New(s store, baseURL string) *Handler {
	return &Handler{
		store:   s,
		baseURL: baseURL,
	}
}

func (h *Handler) HelloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
