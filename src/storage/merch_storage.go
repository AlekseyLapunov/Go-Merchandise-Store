package storage

import (
	"database/sql"
)

type MerchStorage struct {
	db *sql.DB
}

func NewMerchStorage(db *sql.DB) MerchStorage {
	return MerchStorage{db: db}
}
