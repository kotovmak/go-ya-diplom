package store

import (
	"go-ya-diplom/internal/app/model"
)

type OrderRepository struct {
	store *store
}

// Create ...
func (r *OrderRepository) Create(o model.Order) error {
	return r.store.db.QueryRow(
		"INSERT INTO orders (number, status, user_id) VALUES ($1, $2, $3) RETURNING order_id",
		o.Number,
		o.Status,
		o.UserID,
	).Scan(&o.ID)
}

// FindByLogin ...
func (r *OrderRepository) FindByNumber(number string) (model.Order, error) {
	o := model.Order{}
	if err := r.store.db.QueryRow(
		"SELECT order_id, number, status, accrual, uploaded_at, user_id FROM orders WHERE number = $1",
		number,
	).Scan(
		&o.ID,
		&o.Number,
		&o.Status,
		&o.Accrual,
		&o.UploatedAt,
		&o.UserID,
	); err != nil {
		return o, err
	}

	return o, nil
}

func (r *OrderRepository) FindByUser(userID int) ([]model.Order, error) {
	ol := []model.Order{}
	data, err := r.store.db.Query(
		"SELECT order_id, number, status, accrual, uploaded_at, user_id FROM orders WHERE user_id = $1",
		userID,
	)
	if err != nil {
		return ol, err
	}
	for data.Next() {
		o := model.Order{}
		err = data.Scan(
			&o.ID,
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.UploatedAt,
			&o.UserID,
		)
		if err != nil {
			return nil, err
		}
		ol = append(ol, o)
	}

	return ol, nil
}
