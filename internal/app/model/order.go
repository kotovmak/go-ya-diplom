package model

import (
	"go-ya-diplom/internal/app/errors"
	"strconv"
	"time"
)

type Order struct {
	ID         int       `json:"-"`
	Number     string    `json:"number" validate:"required"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual,omitempty"`
	UserID     int       `json:"-"`
	UploatedAt time.Time `json:"uploated_at"`
}

func (o *Order) Validate() error {
	if o.Number == "" {
		return errors.ErrOrderNumberInvalid
	}
	num, err := strconv.Atoi(o.Number)
	if err != nil {
		return err
	}
	ok := o.valid(num)
	if !ok {
		return errors.ErrOrderNumberInvalid
	}
	return nil
}

func (o *Order) valid(number int) bool {
	return (number%10+o.checksum(number/10))%10 == 0
}

func (o *Order) checksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
