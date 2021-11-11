package worker

import (
	"context"
	"database/sql"
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
	recordCh chan model.Order
}

func New(cfg *config.Config, s interfaces.Store) *Worker {
	return &Worker{
		cfg:      cfg,
		store:    s,
		recordCh: make(chan model.Order),
	}
}

func (w *Worker) Init(ctx context.Context) error {
	g, _ := errgroup.WithContext(ctx)

	for i := 0; i < w.cfg.NumOfWorkers; i++ {
		g.Go(func() error {
			var err error
			for ch := range w.recordCh {
				err = w.handle(ctx, ch)
			}
			return err
		})
	}
	return nil
}

func (w *Worker) handle(ctx context.Context, o model.Order) error {
	// TODO Если надо отправлять POST
	// out, err := json.Marshal(o)
	// if err != nil {
	// 	return err
	// }
	// url := w.cfg.AccrualSystemAddress + "/api/orders/" + o.Number
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(out))
	// if err != nil {
	// 	return err
	// }
	// req = req.WithContext(ctx)
	// http.DefaultClient.Do(req)
	err := w.store.Query().Create(ctx, model.Query{
		Order: o.Number,
	})
	log.Println("add to query new Order")

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (w *Worker) check(ctx context.Context, q model.Query) error {
	log.Println("starting check")
	url := w.cfg.AccrualSystemAddress + "/api/orders/" + q.Order
	resp, err := http.Get(url)
	log.Println(url)
	log.Println(resp)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	a := &model.AccrualResponse{}
	if err := json.NewDecoder(resp.Body).Decode(a); err != nil {
		return err
	}
	log.Println(a)

	if a.Status == "PROCESSED" || a.Status == "INVALID" {
		o, err := w.store.Order().FindByNumber(ctx, q.Order)
		if err != nil {
			return err
		}
		o.Status = a.Status
		o.Accrual = int(a.Accrual * 100)
		err = w.store.Order().Update(ctx, o)
		if err != nil {
			return err
		}
		err = w.store.Query().Delete(ctx, q)
		if err != nil {
			return err
		}
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == 429 {
			err := w.store.Query().Create(ctx, model.Query{
				Order:        q.Order,
				ProcessingAt: time.Now().Add(1 * time.Minute),
			})
			if err != nil {
				return err
			}
		}
		err := w.store.Query().Create(ctx, model.Query{
			Order: q.Order,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) Add(o model.Order) {
	w.recordCh <- o
	log.Println("send Order to chan")
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			queryList, err := w.store.Query().Find(ctx)
			if err != nil && err != sql.ErrNoRows {
				log.Println(err)
				return err
			}
			for _, q := range queryList {
				log.Println("found new query row")
				go w.check(ctx, q)
			}
			log.Println("sleep")
			time.Sleep(3 * time.Second)
		}
	}
}
