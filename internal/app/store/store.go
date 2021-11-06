package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-ya-diplom/internal/app/interfaces"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type store struct {
	db             *sql.DB
	userRepository *UserRepository
}

func New(db *sql.DB) *store {
	return &store{
		db: db,
	}
}

func NewDB(ctx context.Context, databaseDSN string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DATABASE_URL '%s'", err)
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to create connection pool '%s'", err)
	}

	err = initMigrations(databaseDSN)
	if err != nil && err != migrate.ErrNoChange && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("unable to create database '%s'", err)
	}
	return db, nil
}

func initMigrations(databaseDSN string) error {
	m, err := migrate.New(
		"file://internal/app/store/migrations",
		databaseDSN)
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

// User ...
func (s *store) User() interfaces.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
