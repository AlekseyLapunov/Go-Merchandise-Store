package storage

import (
    "context"
    "database/sql"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
)

type ManagementStorage struct {
    db *sql.DB
}

func NewManagementStorage(db *sql.DB) ManagementStorage {
    return ManagementStorage{db: db}
}

func (s *ManagementStorage) GetCoins(ctx context.Context, employeeID int) (int, error) {
    return 0, nil
}

func (s *ManagementStorage) GetInventory(ctx context.Context, employeeID int) ([]entity.InventoryItem, error) {
    return []entity.InventoryItem{}, nil
}

func (s *ManagementStorage) GetCoinHistory(ctx context.Context, employeeID int) (*entity.CoinHistory, error) {
    return nil, nil
}

func (s *ManagementStorage) WithdrawCoins(ctx context.Context, employeeID, amount int) error {
    return nil
}

func (s *ManagementStorage) ProvidePurchase(ctx context.Context, employeeID int, item string, cost int) error {
    return nil
}

func (s *ManagementStorage) ProvideOperation(ctx context.Context, senderID, receiverID, amount int) error {
    return nil
}

func (s *ManagementStorage) AddPurchase(ctx context.Context, employeeID int, item string) error {
    return nil
}
