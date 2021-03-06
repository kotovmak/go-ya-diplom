package model

import (
	"go-ya-diplom/internal/app/errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type Withdraw struct {
	ID          int       `json:"-"`
	Order       string    `json:"order" validate:"required"`
	Sum         float32   `json:"sum" validate:"required"`
	Status      string    `json:"status"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      int       `json:"-"`
}

func (w *Withdraw) Validate() error {
	validate := validator.New()
	err := validate.Struct(w)
	if err != nil {
		return err
	}
	num, err := strconv.Atoi(w.Order)
	if err != nil {
		return err
	}
	ok := valid(num)
	if !ok {
		return errors.ErrOrderNumberInvalid
	}
	return nil
}
