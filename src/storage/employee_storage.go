package storage

import (
    "context"
    "database/sql"
    "errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
)

type IEmployeeStorage interface {
    GetEmployee(ctx context.Context, login string) (*entity.Employee, error)
    GetEmployeeID(ctx context.Context, login string) (int, error)
    GetEmployeeLogin(ctx context.Context, employeeID int) (string, error)
    GetEmployeeOrRegister(ctx context.Context, login, password string) (*entity.Employee, error)

    RegisterEmployee(ctx context.Context, login, password string) (*entity.Employee, error)
}

type EmployeeStorage struct {
    db *sql.DB
}

func NewEmployeeStorage(db *sql.DB) EmployeeStorage {
    return EmployeeStorage{db: db}
}

func (s *EmployeeStorage) GetEmployee(ctx context.Context, login string) (*entity.Employee, error) {
    var employee entity.Employee

    err := s.db.QueryRowContext(ctx, "SELECT id, login, password, coins FROM employees WHERE login = $1", login).Scan(
        &employee.ID, &employee.Login, &employee.Password, &employee.Coins,
    )

    return &employee, err
}

func (s *EmployeeStorage) GetEmployeeID(ctx context.Context, login string) (int, error) {
    var id int

    err := s.db.QueryRowContext(ctx, "SELECT id FROM employees WHERE login = $1", login).Scan(&id)

    return id, err
}

func (s *EmployeeStorage) GetEmployeeLogin(ctx context.Context, employeeID int) (string, error) {
    var login string

    err := s.db.QueryRowContext(ctx, "SELECT login FROM employees WHERE id = $1", employeeID).Scan(&login)

    return login, err
}

func (s *EmployeeStorage) GetEmployeeOrRegister(ctx context.Context, login, password string) (*entity.Employee, error) {
    employee, err := s.GetEmployee(ctx, login)

    if employee != nil && err == nil {
        return employee, nil
    }

    var regErr error
    if errors.Is(err, sql.ErrNoRows) {
        employee, regErr = s.RegisterEmployee(ctx, login, password)
    } else {
        return nil, err
    }

    if regErr != nil {
        return employee, regErr
    }

    return employee, nil
}

func (s *EmployeeStorage) RegisterEmployee(ctx context.Context, login, password string) (*entity.Employee, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    var employee entity.Employee
    employee.Login    = login
    employee.Password = string(hashedPassword)
    employee.Coins    = 1000

    err = s.db.QueryRowContext(ctx, `
        INSERT INTO employees (login, password, coins) 
        VALUES ($1, $2, $3)
        RETURNING id
    `, employee.Login, employee.Password, employee.Coins).Scan(&employee.ID)
    if err != nil {
        return nil, err
    }

    return &employee, nil
}

