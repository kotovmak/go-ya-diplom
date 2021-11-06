package interfaces

import "go-ya-diplom/internal/app/model"

type Store interface {
	User() UserRepository
}

type UserRepository interface {
	Create(*model.User) error
	Find(int) (*model.User, error)
	FindByLogin(string) (*model.User, error)
}
