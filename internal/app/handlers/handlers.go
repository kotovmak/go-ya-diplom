package handlers

import (
	"database/sql"
	"encoding/json"
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/errors"
	"go-ya-diplom/internal/app/interfaces"
	"go-ya-diplom/internal/app/model"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	store        interfaces.Store
	cfg          *config.Config
	tokenManager interfaces.TokenManager
	validator    *validator.Validate
}

func New(s interfaces.Store, cfg *config.Config, t interfaces.TokenManager, v *validator.Validate) *Handler {
	return &Handler{
		store:        s,
		cfg:          cfg,
		tokenManager: t,
		validator:    v,
	}
}

func (h *Handler) HelloHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func (h *Handler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &model.LoginRequest{}

		validate, err := h.validate(c.Request(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if validate != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, validate.Error())
		}

		u, err := h.store.User().FindByLogin(req.Login)
		if err != nil || !u.ComparePassword(req.Password) {
			return echo.NewHTTPError(http.StatusUnauthorized, errors.ErrIncorrectEmailOrPassword)
		}

		err = h.tokenManager.JWT().GenerateTokensAndSetCookies(u, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	}
}

func (h *Handler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := &model.LoginRequest{}

		validate, err := h.validate(c.Request(), req)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if validate != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, validate.Error())
		}

		u1, err := h.store.User().FindByLogin(req.Login)
		if u1 != nil {
			return echo.NewHTTPError(http.StatusConflict, errors.ErrAlreadyExists.Error())
		}
		if err != nil && err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		u := &model.User{
			Login:    req.Login,
			Password: req.Password,
		}
		if err := h.store.User().Create(u); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		u.Sanitize()

		err = h.tokenManager.JWT().GenerateTokensAndSetCookies(u, c)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, u)
	}
}

func (h *Handler) validate(r *http.Request, req interface{}) (validate error, err error) {
	if err = json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}
	validate = h.validator.Struct(req)

	return validate, nil
}
