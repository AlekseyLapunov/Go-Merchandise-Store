package storage

import (
	"context"
	"database/sql"
)

type MerchStorage struct {
	db *sql.DB
}

func NewMerchStorage(db *sql.DB) MerchStorage {
	return MerchStorage{db: db}
}

func (s *MerchStorage) GetMerchPrice(ctx context.Context, item string) (int, error) {
	return 0, nil
}

func (s *MerchStorage) AddPurchase(ctx context.Context, employeeID int, item string) error {
	return nil
}
