package repository

import (
	"github.com/jmoiron/sqlx"
)

type ReturnCartPostgres struct {
	db *sqlx.DB
}

func NewReturnCartPostgrers(db *sqlx.DB) *ReturnCartPostgres {
	return &ReturnCartPostgres{db: db}
}

func (r *ReturnCartPostgres) CreateRtCart() error {

	return nil
}
