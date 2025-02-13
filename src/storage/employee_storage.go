package storage

import "database/sql"

type EmployeeStorage struct {
	db *sql.DB
}

func NewEmployeeStorage(db *sql.DB) EmployeeStorage {
	return EmployeeStorage{db: db}
}

