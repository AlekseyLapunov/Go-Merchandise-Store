package storage

import (
    "context"
    "database/sql"
)

type IMerchStorage interface {
    GetMerchCost(ctx context.Context, item string) (int, error)
}

type MerchStorage struct {
    db *sql.DB
}

func NewMerchStorage(db *sql.DB) MerchStorage {
    return MerchStorage{db: db}
}

func (s *MerchStorage) GetMerchCost(ctx context.Context, item string) (int, error) {
    var cost int

    err := s.db.QueryRowContext(ctx, "SELECT cost FROM merch WHERE name = $1", item).Scan(&cost)

    return cost, err
}
