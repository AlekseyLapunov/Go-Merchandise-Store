package storage

import (
	"context"
	"database/sql"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
)

type EmployeeStorage struct {
	db *sql.DB
}

func NewEmployeeStorage(db *sql.DB) EmployeeStorage {
	return EmployeeStorage{db: db}
}

func (s *EmployeeStorage) GetEmployeeByLogin(ctx context.Context, login string) (*entity.Employee, error) {
	return nil, nil
}

func (s *EmployeeStorage) GetBalance(ctx context.Context, employeeID int) (int, error) {
	return 0, nil
}

func (s *EmployeeStorage) GetInventory(ctx context.Context, employeeID int) ([]entity.InventoryItem, error) {
	return []entity.InventoryItem{}, nil
}

func (s *EmployeeStorage) GetCoinHistory(ctx context.Context, employeeID int) (*entity.CoinHistory, error) {
	return nil, nil
}

func (s *EmployeeStorage) ProvideOperation(ctx context.Context, senderID, receiverID, amount int) error {
    return nil
}
