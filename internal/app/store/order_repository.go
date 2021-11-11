package store

import (
	"context"
	"go-ya-diplom/internal/app/model"
)

type OrderRepository struct {
	store *store
}

// Create ...
func (r *OrderRepository) Create(ctx context.Context, o model.Order) error {
	return r.store.db.QueryRowContext(ctx,
		"INSERT INTO orders (number, status, user_id) VALUES ($1, $2, $3) RETURNING order_id",
		o.Number,
		o.Status,
		o.UserID,
	).Scan(&o.ID)
}

func (r *OrderRepository) Update(ctx context.Context, u model.Order) error {
	return r.store.db.QueryRowContext(ctx,
		"UPDATE orders SET (status, accrual) = ($1, $2) WHERE order_id = $3",
		u.Status,
		u.Accrual,
		u.ID,
	).Err()
}

// FindByLogin ...
func (r *OrderRepository) FindByNumber(ctx context.Context, number string) (model.Order, error) {
	o := model.Order{}
	if err := r.store.db.QueryRowContext(ctx,
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

func (r *OrderRepository) FindByUser(ctx context.Context, userID int) ([]model.Order, error) {
	ol := []model.Order{}
	data, err := r.store.db.QueryContext(ctx,
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
