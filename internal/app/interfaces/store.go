package interfaces

import (
	"context"
	"go-ya-diplom/internal/app/model"
)

type Store interface {
	User() UserRepository
	Order() OrderRepository
	Withdraw() WithdrawRepository
}

type UserRepository interface {
	Create(context.Context, model.User) error
	Update(context.Context, model.User) error
	FindByLogin(context.Context, string) (model.User, error)
}

type OrderRepository interface {
	Create(context.Context, model.Order) error
	Update(context.Context, model.Order) error
	FindByNumber(context.Context, string) (model.Order, error)
	FindByUser(context.Context, int) ([]model.Order, error)
}

type WithdrawRepository interface {
	Create(context.Context, model.Withdraw) error
	FindByUser(context.Context, int) ([]model.Withdraw, error)
}
