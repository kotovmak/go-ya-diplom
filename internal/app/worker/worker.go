package worker

import (
	"context"
	"encoding/json"
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/interfaces"
	"go-ya-diplom/internal/app/model"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
)

type Worker struct {
	cfg      *config.Config
	store    interfaces.Store
	query    map[string]model.Order
	recordCh chan bool
}

func New(cfg *config.Config, s interfaces.Store) *Worker {
	return &Worker{
		cfg:      cfg,
		store:    s,
		query:    make(map[string]model.Order),
		recordCh: make(chan bool),
	}
}

func (w *Worker) Init(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < w.cfg.NumOfWorkers; i++ {
		g.Go(func() error {
			var err error
			for ch := range w.recordCh {
				if ch {
					err = w.handle()
				}
			}
			return err
		})
	}
	return nil
}

func (w *Worker) handle() error {
	if len(w.query) == 0 {
		return nil
	}

	var o model.Order
	for _, v := range w.query {
		o = v
		delete(w.query, v.Number)
		break
	}

	url := w.cfg.AccrualSystemAddress + "/api/orders/" + o.Number
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			w.Sleep(1*time.Minute, o)
		}
		w.Sleep(3*time.Second, o)
		return nil
	}

	a := &model.AccrualResponse{}
	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return err
	}

	if a.Status == "PROCESSED" || a.Status == "INVALID" {
		o.Status = a.Status
		o.Accrual = int(a.Accrual * 100)
		w.store.Order().Update(o)
		return nil
	}

	w.Sleep(3*time.Second, o)
	return nil
}

func (w *Worker) Sleep(d time.Duration, o model.Order) {
	time.Sleep(d)
	w.query[o.Number] = o
	w.recordCh <- true
}

func (w *Worker) Add(o model.Order) {
	w.query[o.Number] = o
	log.Printf("add query %v", w.query)
	w.recordCh <- true
}
