package auth

import (
	"go-ya-diplom/internal/app/errors"
	"go-ya-diplom/internal/app/model"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	AccessTokenCookieName  = "access-token"
	RefreshTokenCookieName = "refresh-token"
)

type JWTCookie struct {
	signingKey string
	ttl        time.Duration
	refreshTTL time.Duration
	refreshKey string
}

type Claims struct {
	Login string `json:"login"`
	jwt.StandardClaims
}

func (t *JWTCookie) GenerateTokensAndSetCookies(user *model.User, c echo.Context) error {
	accessToken, exp, err := t.generateAccessToken(user)
	if err != nil {
		return err
	}

	t.setTokenCookie(AccessTokenCookieName, accessToken, exp, c)
	t.setUserCookie(user, exp, c)

	refreshToken, exp, err := t.generateRefreshToken(user)
	if err != nil {
		return err
	}
	t.setTokenCookie(RefreshTokenCookieName, refreshToken, exp, c)

	return nil
}

func (t *JWTCookie) generateRefreshToken(user *model.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(t.refreshTTL)

	return t.generateToken(user, expirationTime, []byte(t.refreshKey))
}

func (t *JWTCookie) generateAccessToken(user *model.User) (string, time.Time, error) {
	// Declare the expiration time of the token (1h).
	expirationTime := time.Now().Add(t.ttl)

	return t.generateToken(user, expirationTime, []byte(t.signingKey))
}

func (t *JWTCookie) generateToken(user *model.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	claims := &Claims{
		Login: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

func (t *JWTCookie) setTokenCookie(name, token string, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func (t *JWTCookie) setUserCookie(user *model.User, expiration time.Time, c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "user"
	cookie.Value = user.Login
	cookie.Expires = expiration
	cookie.Path = "/"
	c.SetCookie(cookie)
}

func (t *JWTCookie) JWTErrorChecker(err error, c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, errors.ErrIncorrectEmailOrPassword)
}
