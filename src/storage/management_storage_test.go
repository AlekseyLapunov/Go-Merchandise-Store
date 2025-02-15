package storage_test

import (
    "context"
    "testing"
    //"database/sql"
    //"errors"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestGetCoins(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"coins"}).AddRow(1000)
    mock.ExpectQuery("SELECT coins FROM employees WHERE id = ?").
        WithArgs(1).
            WillReturnRows(rows)

    storage := storage.NewManagementStorage(db)
    coins, err := storage.GetCoins(context.Background(), 1)
    assert.NoError(t, err)
    assert.Equal(t, 1000, coins)

    assert.NoError(t, mock.ExpectationsWereMet())
}


func TestGetInventory(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"name", "count"}).
        AddRow("hoody", 3).
            AddRow("book", 5)
    mock.ExpectQuery(`
        SELECT m\.name, COUNT\(p\.id\)
        FROM purchases AS p
        JOIN merch AS m ON p\.merch_id = m\.id
        WHERE p\.user_id = \$1
        GROUP BY m\.name`).
            WithArgs(1).
                WillReturnRows(rows)

    storage := storage.NewManagementStorage(db)
    inventory, err := storage.GetInventory(context.Background(), 1)

    assert.NoError(t, err)
    assert.Equal(t, []entity.InventoryItem{
        {Type: "hoody", Quantity: 3},
        {Type: "book", Quantity: 5},
    }, inventory)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestProvidePurchase(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    mock.ExpectBegin()
    mock.ExpectExec(`UPDATE employees SET coins = coins - \$1 WHERE id = \$2 AND coins >= \$1`).
        WithArgs(100, 1).
            WillReturnResult(sqlmock.NewResult(0, 1))
    mock.ExpectQuery(`SELECT id FROM merch WHERE name = \$1`).
        WithArgs("wallet").
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
    mock.ExpectExec("INSERT INTO purchases .+ VALUES .+").
        WithArgs(1, 1).
            WillReturnResult(sqlmock.NewResult(0, 1))
    mock.ExpectCommit()

    storage := storage.NewManagementStorage(db)

    err = storage.ProvidePurchase(context.Background(), 1, "wallet", 100)

    assert.NoError(t, err)

    assert.NoError(t, mock.ExpectationsWereMet())    
}

func TestProvideOperation_Success(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    mock.ExpectBegin()
    mock.ExpectExec(`UPDATE employees SET coins = coins - \$1 WHERE id = \$2 AND coins >= \$1`).
        WithArgs(111, 1).
            WillReturnResult(sqlmock.NewResult(0, 1))

    mock.ExpectExec(`UPDATE employees SET coins = coins \+ \$1 WHERE id = \$2`).
        WithArgs(111, 2).
            WillReturnResult(sqlmock.NewResult(0, 1))

    mock.ExpectExec("INSERT INTO operations .+ VALUES .+").
        WithArgs(1, 2, 111).
            WillReturnResult(sqlmock.NewResult(0, 1))
    mock.ExpectCommit()

    storage := storage.NewManagementStorage(db)

    err = storage.ProvideOperation(context.Background(), 1, 2, 111)

    assert.NoError(t, err)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFetchReceivedHistory(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    storage := storage.NewManagementStorage(db)

    mock.ExpectQuery(`
        SELECT e\.login, o\.amount
        FROM operations AS o JOIN employees AS e ON o\.send_user_id = e\.id
        WHERE o\.recv_user_id = \$1
    `).
        WithArgs(2).
            WillReturnRows(sqlmock.NewRows([]string{"login", "amount"}).
                AddRow("aleksey", 111).
                AddRow("maksim",  555),
            )

    result, err := storage.FetchReceivedHistory(context.Background(), 2)

    assert.NoError(t, err)
    assert.Len(t, result, 2)
    assert.Equal(t, "aleksey", result[0].FromUser)
    assert.Equal(t, 111,       result[0].Amount)
    assert.Equal(t, "maksim",  result[1].FromUser)
    assert.Equal(t, 555,       result[1].Amount)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFetchSentHistory(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    storage := storage.NewManagementStorage(db)

    mock.ExpectQuery(`
        SELECT e\.login, o\.amount
        FROM operations AS o JOIN employees AS e ON o\.recv_user_id = e\.id
        WHERE o\.send_user_id = \$1
    `).
        WithArgs(1).
            WillReturnRows(sqlmock.NewRows([]string{"login", "amount"}).
                AddRow("sveta", 742).
                AddRow("maria", 123),
            )

    result, err := storage.FetchSentHistory(context.Background(), 1)

    assert.NoError(t, err)
    assert.Len(t, result, 2)
    assert.Equal(t, "sveta", result[0].ToUser)
    assert.Equal(t, 742,     result[0].Amount)
    assert.Equal(t, "maria", result[1].ToUser)
    assert.Equal(t, 123,     result[1].Amount)

    assert.NoError(t, mock.ExpectationsWereMet())
}

