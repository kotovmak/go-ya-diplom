package errors

import "errors"

var (
	ErrNotFound                  = errors.New("not found")
	ErrRowDeleted                = errors.New("row deleted")
	ErrIncorrectEmailOrPassword  = errors.New("incorrect email or password")
	ErrNotAuthenticated          = errors.New("not authenticated")
	ErrEmptySigningKey           = errors.New("empty signing key")
	ErrUnexpectedSigningMethod   = errors.New("unexpected signing method")
	ErrAlreadyExists             = errors.New("already exist")
	ErrOrderNumberInvalid        = errors.New("order number is invalid")
	ErrNotEnoughMoney            = errors.New("not enough money")
	ErrAccrualSystemAddressEmpty = errors.New("accrual system address empty")
)
