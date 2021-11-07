package handlers

import (
	"database/sql"
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/errors"
	"go-ya-diplom/internal/app/interfaces"
	"go-ya-diplom/internal/app/model"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store        interfaces.Store
	cfg          *config.Config
	tokenManager interfaces.TokenManager
}

func New(s interfaces.Store, cfg *config.Config, t interfaces.TokenManager) *Handler {
	return &Handler{
		store:        s,
		cfg:          cfg,
		tokenManager: t,
	}
}

func (h *Handler) HelloHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	}
}

func (h *Handler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		u := model.User{}

		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		validate := u.Validate()
		if validate != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, validate.Error())
		}

		u1, err := h.store.User().FindByLogin(u.Login)
		if err != nil || !u1.ComparePassword(u.Password) {
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
		u := model.User{}

		if err := c.Bind(&u); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		validate := u.Validate()
		if validate != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, validate.Error())
		}

		u1, err := h.store.User().FindByLogin(u.Login)
		if u1.ID > 0 {
			return echo.NewHTTPError(http.StatusConflict, errors.ErrAlreadyExists.Error())
		}
		if err != nil && err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

func (h *Handler) OrderUpload() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("user")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		u, err := h.store.User().FindByLogin(user.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		body, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		o := model.Order{
			Status:     "NEW",
			UploatedAt: time.Now(),
			Number:     string(body),
			UserID:     u.ID,
		}

		validate := o.Validate()
		if validate != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, validate.Error())
		}

		o1, err := h.store.Order().FindByNumber(o.Number)
		if o1.ID > 0 {
			if o1.UserID == u.ID {
				return echo.NewHTTPError(http.StatusOK, errors.ErrAlreadyExists.Error())
			}
			return echo.NewHTTPError(http.StatusConflict, errors.ErrAlreadyExists.Error())
		}
		if err != nil && err != sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = h.store.Order().Create(o)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusAccepted, o)
	}
}

func (h *Handler) OrderList() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("user")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		u, err := h.store.User().FindByLogin(user.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		orderList, err := h.store.Order().FindByUser(u.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNoContent, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, orderList)
	}
}

func (h *Handler) Balance() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("user")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		u, err := h.store.User().FindByLogin(user.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		b := model.Balance{
			Balance:   float32(u.Balance) / 100,
			Withdrawn: float32(u.Withdrawn) / 100,
		}

		return c.JSON(http.StatusOK, b)
	}
}

func (h *Handler) Withdraw() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("user")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		u, err := h.store.User().FindByLogin(user.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		wr := model.Withdraw{}

		if err := c.Bind(&wr); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		validate := wr.Validate()
		if validate != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, validate.Error())
		}

		sum := int(wr.Sum * 100)
		if u.Balance < sum {
			return echo.NewHTTPError(http.StatusPaymentRequired, errors.ErrNotEnoughMoney)
		}

		w := model.Withdraw{
			Status: "PROCESSED",
			UserID: u.ID,
			Sum:    wr.Sum,
			Order:  wr.Order,
		}

		err = h.store.Withdraw().Create(w)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		u.Balance -= sum
		u.Withdrawn += sum
		h.store.User().Update(u)

		return c.JSON(http.StatusOK, w)
	}
}

func (h *Handler) WithdrawList() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := c.Cookie("user")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		u, err := h.store.User().FindByLogin(user.Value)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
		}

		withdrawList, err := h.store.Withdraw().FindByUser(u.ID)
		if err != nil {
			if err == sql.ErrNoRows {
				return echo.NewHTTPError(http.StatusNoContent, err.Error())
			}
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, withdrawList)
	}
}
