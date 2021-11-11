package auth

import (
	"go-ya-diplom/internal/app/config"
	"go-ya-diplom/internal/app/interfaces"
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

	a.jwtCookie = &JWTCookie{
		signingKey: a.cfg.SigningKey,
		refreshTTL: a.cfg.TokenTTL,
		refreshKey: a.cfg.RefreshKey,
		ttl:        a.cfg.RefreshTTL,
	}

	return a.jwtCookie
}
