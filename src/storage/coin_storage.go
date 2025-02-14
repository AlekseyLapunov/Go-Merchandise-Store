package storage

import (
	"context"
	"database/sql"
)

type CoinStorage struct {
	db *sql.DB
}

func NewCoinStorage(db *sql.DB) CoinStorage {
	return CoinStorage{db: db}
}

func (s *CoinStorage) GetBalance(ctx context.Context, employeeID int) (int, error) {
	return 0, nil
}

func (s *CoinStorage) ProvideOperation(ctx context.Context, senderID, receiverID, amount int) error {
    return nil
}

