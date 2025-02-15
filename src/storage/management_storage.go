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
        FROM purchases AS p
        JOIN merch AS m ON p.merch_id = m.id
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
    var received []entity.RecvEntry
    var err error
    if received, err = s.fetchReceivedHistory(ctx, employeeID); err != nil {
        return nil, err
    }

    var sent []entity.SentEntry
    if sent, err = s.fetchSentHistory(ctx, employeeID); err != nil {
        return nil, err
    }

    return &entity.CoinHistory{
        Received: received,
        Sent:     sent,
    }, nil
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

func (s *ManagementStorage) fetchReceivedHistory(ctx context.Context, receiverID int) ([]entity.RecvEntry, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT e.login, o.amount 
        FROM operations AS o
        JOIN employees AS e ON o.send_user_id = e.id
        WHERE o.recv_user_id = $1
    `, receiverID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var received []entity.RecvEntry
    for rows.Next() {
        var entry entity.RecvEntry
        if err := rows.Scan(&entry.FromUser, &entry.Amount); err != nil {
            return nil, err
        }
        received = append(received, entry)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return received, nil  
}

func (s *ManagementStorage) fetchSentHistory(ctx context.Context, senderID int) ([]entity.SentEntry, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT e.login, o.amount 
        FROM operations AS o
        JOIN employees AS e ON o.recv_user_id = e.id
        WHERE o.send_user_id = $1
    `, senderID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sent []entity.SentEntry
    for rows.Next() {
        var entry entity.SentEntry
        if err := rows.Scan(&entry.ToUser, &entry.Amount); err != nil {
            return nil, err
        }
        sent = append(sent, entry)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }
    
    return sent, nil
}

