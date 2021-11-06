package auth

import (
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/interfaces"
	"log"
	"time"
)

type auth struct {
	cfg       *config.Config
	jwtCookie *JWTCookie
}

func New(cfg *config.Config) *auth {
	return &auth{
		cfg: cfg,
	}
}

func (a *auth) JWT() interfaces.JWTCookie {
	if a.jwtCookie != nil {
		return a.jwtCookie
	}

	ttl, err := time.ParseDuration(a.cfg.TokenTTL)
	if err != nil {
		log.Fatal(err)
	}
	rttl, err := time.ParseDuration(a.cfg.RefreshTTL)
	if err != nil {
		log.Fatal(err)
	}

	a.jwtCookie = &JWTCookie{
		signingKey: a.cfg.SigningKey,
		refreshTTL: rttl,
		refreshKey: a.cfg.RefreshKey,
		ttl:        ttl,
	}

	return a.jwtCookie
}
