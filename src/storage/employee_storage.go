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

func (s *EmployeeStorage) GetEmployee(ctx context.Context, login string) (*entity.Employee, error) {
    var employee entity.Employee

    err := s.db.QueryRowContext(ctx, "SELECT id, login, password, coins FROM employees WHERE login = $1", login).Scan(
        &employee.ID, &employee.Login, &employee.Login, &employee.Password, &employee.Coins,
    )

    return &employee, err
}

// register employee
