package interfaces

import (
	"go-ya-diplom/internal/app/model"

	"github.com/labstack/echo/v4"
)

type TokenManager interface {
	JWT() JWTCookie
}

type JWTCookie interface {
	GenerateTokensAndSetCookies(user *model.User, c echo.Context) error
	JWTErrorChecker(err error, c echo.Context) error
}
