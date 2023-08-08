package postgres

import (
	"errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(connectionString string) (*Storage, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, errors.New("db connection error: " + err.Error())
	}

	return &Storage{db}, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
