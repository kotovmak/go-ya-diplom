package interfaces

import (
	"context"
	"go-ya-diplom/internal/app/model"
)

type Store interface {
	User() UserRepository
	Order() OrderRepository
	Withdraw() WithdrawRepository
	Query() QueryRepository
}

type UserRepository interface {
	Create(context.Context, model.User) error
	Update(context.Context, model.User) error
	FindByLogin(context.Context, string) (model.User, error)
	FindByID(context.Context, int) (model.User, error)
}

type OrderRepository interface {
	Create(context.Context, model.Order) error
	Update(context.Context, model.Order) error
	FindByNumber(context.Context, string) (model.Order, error)
	FindByUser(context.Context, int) ([]model.OrderList, error)
}

type WithdrawRepository interface {
	Create(context.Context, model.Withdraw) error
	FindByUser(context.Context, int) ([]model.Withdraw, error)
}

type QueryRepository interface {
	Create(context.Context, model.Query) error
	FindByOrder(context.Context, string) (model.Query, error)
	Find(context.Context) ([]model.Query, error)
	Delete(context.Context, model.Query) error
}
