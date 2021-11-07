package store

import (
	"go-ya-diplom/internal/app/model"
)

type UserRepository struct {
	store *store
}

// Create ...
func (r *UserRepository) Create(u model.User) error {
	if err := u.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO users (login, password) VALUES ($1, $2) RETURNING user_id",
		u.Login,
		u.EncryptedPassword,
	).Scan(&u.ID)
}

func (r *UserRepository) Update(u model.User) error {
	return r.store.db.QueryRow(
		"UPDATE users SET (balance, withdrawn) = ($1, $2) WHERE user_id = $3",
		u.Balance,
		u.Withdrawn,
		u.ID,
	).Err()
}

// FindByLogin ...
func (r *UserRepository) FindByLogin(login string) (model.User, error) {
	u := model.User{}
	if err := r.store.db.QueryRow(
		"SELECT user_id, login, password, balance, withdrawn FROM users WHERE login = $1",
		login,
	).Scan(
		&u.ID,
		&u.Login,
		&u.EncryptedPassword,
		&u.Balance,
		&u.Withdrawn,
	); err != nil {
		return u, err
	}

	return u, nil
}
