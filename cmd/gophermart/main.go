package main

import (
	"context"
	"go-ya-diplom/internal/app/auth"
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/errors"
	"go-ya-diplom/internal/app/handlers"
	"go-ya-diplom/internal/app/store"
	"go-ya-diplom/internal/app/worker"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	cfg.InitFlags()

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
		//middleware.Recover(),
		middleware.Decompress(),
		middleware.Gzip(),
	)

	s := store.New(db)
	t := auth.New(cfg)
	w := worker.New(cfg, s)
	if cfg.AccrualSystemAddress == "" {
		log.Fatal(errors.ErrAccrualSystemAddressEmpty)
	}
	w.Init(ctx)
	go w.Run(ctx)
	h := handlers.New(s, cfg, t, w)

	e.GET("/", h.HelloHandler())
	v1 := e.Group("/api")
	{
		user := v1.Group("/user")
		{
			user.POST("/register", h.Register())
			user.POST("/login", h.Login())
			authGroup := user.Group("")
			{
				authGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
					Claims:                  &auth.Claims{},
					SigningKey:              []byte(cfg.SigningKey),
					TokenLookup:             "cookie:access-token",
					ErrorHandlerWithContext: t.JWT().JWTErrorChecker,
				}))
				authGroup.POST("/orders", h.OrderUpload())
				authGroup.GET("/orders", h.OrderList())
				balance := authGroup.Group("/balance")
				{
					balance.GET("", h.Balance())
					balance.POST("/withdraw", h.Withdraw())
					balance.GET("/withdrawals", h.WithdrawList())
				}
			}
		}
	}

	log.Printf("[INIT] ServerAddress '%s'", cfg.ServerAddress)
	log.Printf("[INIT] BaseURL '%s'", cfg.BaseURL)
	if err := e.Start(cfg.ServerAddress); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
