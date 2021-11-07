package worker

import (
	"context"
	"encoding/json"
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/interfaces"
	"go-ya-diplom/internal/app/model"
	"net/http"

	"golang.org/x/sync/errgroup"
)

var (
	recordCh chan model.Order
)

type Worker struct {
	cfg   *config.Config
	store interfaces.Store
}

func New(cfg *config.Config, s interfaces.Store) *Worker {
	return &Worker{
		cfg:   cfg,
		store: s,
	}
}

func (w *Worker) Init(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < w.cfg.NumOfWorkers; i++ {
		g.Go(func() error {
			for order := range recordCh {
				w.handle(order)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (w *Worker) handle(o model.Order) error {
	resp, err := http.Get(w.cfg.AccrualSystemAddress + "/api/orders/" + o.Number)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return err
	}

	a := &model.AccrualResponse{}
	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return err
	}

	o.Status = a.Status
	o.Accrual = int(a.Accrual * 100)
	w.store.Order().Update(o)

	return nil
}

func (w *Worker) Run(o model.Order) {
	recordCh <- o
}
