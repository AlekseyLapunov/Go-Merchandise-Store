package storage_test

import (
    "context"
    "testing"
    "database/sql"
    "errors"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
)

func TestGetEmployee(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"id", "login", "password", "coins"}).AddRow(1, "aleksey", "superpass", 1000)
    mock.ExpectQuery(`
        SELECT id, login, password, coins
        FROM employees
        WHERE login = ?
    `).WithArgs("aleksey").WillReturnRows(rows)

    storage := storage.NewEmployeeStorage(db)
    employee, err := storage.GetEmployee(context.Background(), "aleksey")

    assert.NoError(t, err)
    assert.Equal(t, &entity.Employee{
        ID:       1,
        Login:    "aleksey",
        Password: "superpass",
        Coins:    1000,
    }, employee)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterEmployee(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    mock.ExpectQuery("INSERT INTO employees .+ RETURNING id").
        WithArgs("aleksey", sqlmock.AnyArg(), 1000).
            WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    storage := storage.NewEmployeeStorage(db)


    employee, err := storage.RegisterEmployee(context.Background(), "aleksey", "superpass")
    assert.NoError(t, err)
    assert.Equal(t, &entity.Employee{
        ID:       1,
        Login:    "aleksey",
        Password: employee.Password, // password will be hashed so i need to do this 'wrong thing' here
        Coins:    1000,
    }, employee)

    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRegisterEmployee_DuplicateLogin(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Could not create sql mock: %v", err)
    }
    defer db.Close()

    mock.ExpectQuery("INSERT INTO employees .+ RETURNING id").
         WithArgs("aleksey", sqlmock.AnyArg(), 1000).
            WillReturnError(sql.ErrNoRows)

    storage := storage.NewEmployeeStorage(db)
    _, err = storage.RegisterEmployee(context.Background(), "aleksey", "superpass") // already in the base

    assert.Error(t, err)
    assert.True(t, errors.Is(err, sql.ErrNoRows))

    assert.NoError(t, mock.ExpectationsWereMet())
}
