package store

import (
	"go-ya-diplom/internal/app/model"
)

type WithdrawRepository struct {
	store *store
}

// Create ...
func (r *WithdrawRepository) Create(w model.Withdraw) error {
	return r.store.db.QueryRow(
		"INSERT INTO withdraws (\"order\", status, sum, user_id) VALUES ($1, $2, $3, $4)",
		w.Order,
		w.Status,
		w.Sum,
		w.UserID,
	).Err()
}

func (r *WithdrawRepository) Find() ([]model.Withdraw, error) {
	ol := []model.Withdraw{}
	data, err := r.store.db.Query(
		"SELECT withdraw_id, \"order\", status, sum, processed_at, user_id FROM withdraws",
	)
	if err != nil {
		return ol, err
	}
	for data.Next() {
		o := model.Withdraw{}
		err = data.Scan(
			&o.ID,
			&o.Order,
			&o.Status,
			&o.Sum,
			&o.ProcessedAt,
			&o.UserID,
		)
		if err != nil {
			return nil, err
		}
		ol = append(ol, o)
	}

	return ol, nil
}
