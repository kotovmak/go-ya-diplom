package store

import (
	"context"
	"go-ya-diplom/internal/app/model"
)

type QueryRepository struct {
	store *store
}

func (r *QueryRepository) Create(ctx context.Context, q model.Query) error {
	return r.store.db.QueryRowContext(ctx,
		"INSERT INTO \"query\" (\"order\", processing_at) VALUES ($1, $2)",
		q.Order,
		q.ProcessingAt,
	).Err()
}

func (r *QueryRepository) FindByOrder(ctx context.Context, order string) (model.Query, error) {
	o := model.Query{}
	if err := r.store.db.QueryRowContext(ctx,
		"SELECT query_id, \"order\", processing_at FROM \"query\" WHERE \"order\" = $1",
		order,
	).Scan(
		&o.ID,
		&o.Order,
		&o.ProcessingAt,
	); err != nil {
		return o, err
	}

	return o, nil
}

func (r *QueryRepository) Find(ctx context.Context) ([]model.Query, error) {
	ol := []model.Query{}
	data, err := r.store.db.QueryContext(ctx,
		"SELECT query_id, \"order\", processing_at FROM \"query\" WHERE processing_at <= now()",
	)
	if err != nil {
		return ol, err
	}
	for data.Next() {
		o := model.Query{}
		err = data.Scan(
			&o.ID,
			&o.Order,
			&o.ProcessingAt,
		)
		if err != nil {
			return nil, err
		}
		ol = append(ol, o)
	}

	return ol, nil
}

func (r *QueryRepository) Delete(ctx context.Context, q model.Query) error {
	return r.store.db.QueryRowContext(ctx,
		"DELETE FROM \"query\" WHERE query_id = $1",
		q.ID,
	).Err()
}
