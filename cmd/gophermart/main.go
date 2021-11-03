package main

import (
	"context"
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/handlers"
	"go-ya-diplom/internal/app/store"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.New()

	ctx := context.Background()

	db, err := store.NewDB(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(
		middleware.RequestID(),
		middleware.Logger(),
		middleware.Recover(),
		middleware.Decompress(),
		middleware.Gzip(),
	)

	s := store.New(db)
	h := handlers.New(s, cfg.BaseURL)

	e.GET("/", h.HelloHandler)

	log.Printf("[INIT] ServerAddress '%s'", cfg.ServerAddress)
	log.Printf("[INIT] BaseURL '%s'", cfg.BaseURL)
}
