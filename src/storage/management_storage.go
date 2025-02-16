package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
)

type IManagementStorage interface {
    GetCoins(ctx context.Context, employeeID int) (int, error)
    GetInventory(ctx context.Context, employeeID int) ([]entity.InventoryItem, error)
    GetCoinHistory(ctx context.Context, employeeID int) (*entity.CoinHistory, error)
    ProvidePurchase(ctx context.Context, employeeID int, item string, cost int) error
    ProvideOperation(ctx context.Context, senderID, receiverID, amount int) error

    FetchReceivedHistory(ctx context.Context, receiverID int) ([]entity.RecvEntry, error)
    FetchSentHistory(ctx context.Context, senderID int) ([]entity.SentEntry, error)
}

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
        WHERE p.emp_id = $1
        GROUP BY m.name
    `, employeeID)
    if errors.Is(err, sql.ErrNoRows) {
        return []entity.InventoryItem{}, nil
    }

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
    if received, err = s.FetchReceivedHistory(ctx, employeeID); err != nil {
        return nil, err
    }

    var sent []entity.SentEntry
    if sent, err = s.FetchSentHistory(ctx, employeeID); err != nil {
        return nil, err
    }

    return &entity.CoinHistory{
        Received: received,
        Sent:     sent,
    }, nil
}

func (s *ManagementStorage) ProvidePurchase(ctx context.Context, employeeID int, item string, cost int) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.ExecContext(ctx, `
        UPDATE employees 
        SET coins = coins - $1 
        WHERE id = $2 AND coins >= $1
    `, cost, employeeID)
    if err != nil {
        return err
    }

    var merchID int
    err = tx.QueryRowContext(ctx, `
        SELECT id 
        FROM merch 
        WHERE name = $1
    `, item).Scan(&merchID)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx, `
        INSERT INTO purchases (emp_id, merch_id) 
        VALUES ($1, $2)
    `, employeeID, merchID)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *ManagementStorage) ProvideOperation(ctx context.Context, senderID, receiverID, amount int) error {
    tx, err := s.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    _, err = tx.ExecContext(ctx, `
        UPDATE employees 
        SET coins = coins - $1 
        WHERE id = $2 AND coins >= $1
    `, amount, senderID)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx, `
        UPDATE employees 
        SET coins = coins + $1 
        WHERE id = $2
    `, amount, receiverID)
    if err != nil {
        return err
    }

    _, err = tx.ExecContext(ctx, `
        INSERT INTO operations (send_emp_id, recv_emp_id, amount) 
        VALUES ($1, $2, $3)
    `, senderID, receiverID, amount)
    if err != nil {
        return err
    }

    return tx.Commit()
}

func (s *ManagementStorage) FetchReceivedHistory(ctx context.Context, receiverID int) ([]entity.RecvEntry, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT e.login, o.amount 
        FROM operations AS o
        JOIN employees AS e ON o.send_emp_id = e.id
        WHERE o.recv_emp_id = $1
    `, receiverID)
    if errors.Is(err, sql.ErrNoRows) {
        return []entity.RecvEntry{}, nil
    }

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

func (s *ManagementStorage) FetchSentHistory(ctx context.Context, senderID int) ([]entity.SentEntry, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT e.login, o.amount 
        FROM operations AS o
        JOIN employees AS e ON o.recv_emp_id = e.id
        WHERE o.send_emp_id = $1
    `, senderID)
    if errors.Is(err, sql.ErrNoRows) {
        return []entity.SentEntry{}, nil
    }

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

