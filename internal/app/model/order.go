package model

import (
	"go-ya-diplom/internal/app/errors"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

type Order struct {
	ID         int       `json:"-"`
	Number     string    `json:"number" validate:"required"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UserID     int       `json:"-"`
	UploatedAt time.Time `json:"uploated_at"`
}

type AccrualResponse struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual"`
}

func (o *Order) Validate() error {
	validate := validator.New()
	err := validate.Struct(o)
	if err != nil {
		return err
	}
	num, err := strconv.Atoi(o.Number)
	if err != nil {
		return err
	}
	ok := valid(num)
	if !ok {
		return errors.ErrOrderNumberInvalid
	}
	return nil
}
