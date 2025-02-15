package storage_test

import (
    "context"
    "testing"
    "database/sql"
    "errors"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestGetMerchCost_ValidCost(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    merchStorage := storage.NewMerchStorage(db)
    ctx := context.Background()
    item := "cup"

    mock.ExpectQuery("SELECT cost FROM merch WHERE name = ?").
        WithArgs(item).
            WillReturnRows(sqlmock.NewRows([]string{"cost"}).AddRow(20))

    cost, err := merchStorage.GetMerchCost(ctx, item)
    assert.NoError(t, err)
    assert.Equal(t, 20, cost)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMerchCost_UnknownItem(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    merchStorage := storage.NewMerchStorage(db)
    ctx := context.Background()
    item := "lgkdfklgj"

    mock.ExpectQuery("SELECT cost FROM merch WHERE name = ?").
        WithArgs(item).
            WillReturnError(sql.ErrNoRows)

    cost, err := merchStorage.GetMerchCost(ctx, item)
    assert.ErrorIs(t, err, sql.ErrNoRows)
    assert.Equal(t, 0, cost)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetMerchCost_QueryFailed(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    merchStorage := storage.NewMerchStorage(db)
    ctx := context.Background()
    item := "hoody"

    mock.ExpectQuery("SELECT cost FROM merch WHERE name = ?").
        WithArgs(item).
            WillReturnError(errors.New("query failed"))

    cost, err := merchStorage.GetMerchCost(ctx, item)
    assert.Error(t, err)
    assert.Equal(t, 0, cost)
    
    assert.NoError(t, mock.ExpectationsWereMet())
}
