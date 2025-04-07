package db

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type DB struct {
	db     *sql.DB
	logger *zerolog.Logger
}

func New(db *sql.DB, logger *zerolog.Logger) *DB {
	return &DB{db, logger}
}
