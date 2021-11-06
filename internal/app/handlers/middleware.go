package handlers

import (
	"go-ya-diplom/internal/app/auth"
	"go-ya-diplom/internal/app/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func (h *Handler) TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Get("user") == nil {
			return next(c)
		}
		u := c.Get("user").(*jwt.Token)

		claims := u.Claims.(*auth.Claims)

		if time.Until(time.Unix(claims.ExpiresAt, 0)) < 15*time.Minute {
			rc, err := c.Cookie(auth.RefreshTokenCookieName)
			if err == nil && rc != nil {
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(h.cfg.RefreshKey), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						c.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					_ = h.tokenManager.JWT().GenerateTokensAndSetCookies(&model.User{
						Login: claims.Login,
					}, c)
				}
			}
		}

		return next(c)
	}
}
