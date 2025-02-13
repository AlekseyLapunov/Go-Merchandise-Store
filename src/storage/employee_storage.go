package storage

import "database/sql"

type EmployeeStorage struct {
	db *sql.DB
}

func NewEmployeeStorage(db *sql.DB) *EmployeeStorage {
	if db == nil {
		return nil
	}

	return &EmployeeStorage{db: db}
}

