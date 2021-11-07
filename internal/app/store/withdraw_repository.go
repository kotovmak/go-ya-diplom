package store

import (
	"go-ya-diplom/internal/app/model"
)

type WithdrawRepository struct {
	store *store
}

// Create ...
func (r *WithdrawRepository) Create(w model.Withdraw) error {
	sum := int(w.Sum * 100)
	return r.store.db.QueryRow(
		"INSERT INTO withdraws (\"order\", status, sum, user_id) VALUES ($1, $2, $3, $4)",
		w.Order,
		w.Status,
		sum,
		w.UserID,
	).Err()
}

func (r *WithdrawRepository) FindByUser(userID int) ([]model.Withdraw, error) {
	ol := []model.Withdraw{}
	data, err := r.store.db.Query(
		"SELECT withdraw_id, \"order\", status, sum, processed_at, user_id FROM withdraws WHERE user_id=$1",
		userID,
	)
	if err != nil {
		return ol, err
	}
	for data.Next() {
		o := model.Withdraw{}
		var sum int
		err = data.Scan(
			&o.ID,
			&o.Order,
			&o.Status,
			&sum,
			&o.ProcessedAt,
			&o.UserID,
		)
		if err != nil {
			return nil, err
		}
		o.Sum = float32(sum) / 100
		ol = append(ol, o)
	}

	return ol, nil
}
