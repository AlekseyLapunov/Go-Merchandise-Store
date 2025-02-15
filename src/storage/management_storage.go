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
    var coins int

    err := s.db.QueryRowContext(ctx, "SELECT coins FROM employees WHERE id = $1", employeeID).Scan(&coins)

    return coins, err
}

func (s *ManagementStorage) GetInventory(ctx context.Context, employeeID int) ([]entity.InventoryItem, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT m.name, COUNT(p.id) 
        FROM purchases p
        JOIN merch m ON p.merch_id = m.id
        WHERE p.user_id = $1
        GROUP BY m.name
    `, employeeID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var inventory []entity.InventoryItem
    for rows.Next() {
        var item entity.InventoryItem

        if err := rows.Scan(&item.Type, &item.Quantity); err != nil {
            return nil, err
        }
        
        inventory = append(inventory, item)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return inventory, nil
}

func (s *ManagementStorage) GetCoinHistory(ctx context.Context, employeeID int) (*entity.CoinHistory, error) {
    return nil, nil
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
