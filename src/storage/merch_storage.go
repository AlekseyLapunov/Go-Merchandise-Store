package storage

import "database/sql"

type MerchStorage struct {
	db *sql.DB
}

func NewMerchStorage(db *sql.DB) *MerchStorage {
	if db == nil {
		return nil
	}

	return &MerchStorage{db: db}
}