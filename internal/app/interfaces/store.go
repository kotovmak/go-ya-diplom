package interfaces

import "go-ya-diplom/internal/app/model"

type Store interface {
	User() UserRepository
	Order() OrderRepository
	Withdraw() WithdrawRepository
}

type UserRepository interface {
	Create(model.User) error
	Update(model.User) error
	FindByLogin(string) (model.User, error)
}

type OrderRepository interface {
	Create(model.Order) error
	FindByNumber(string) (model.Order, error)
	FindByUser(int) ([]model.Order, error)
}

type WithdrawRepository interface {
	Create(model.Withdraw) error
	FindByUser(int) ([]model.Withdraw, error)
}
