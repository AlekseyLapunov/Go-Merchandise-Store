package usecase

import "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"

type EmployeeUsecase struct {
	storage storage.EmployeeStorage
}

func NewEmployeeUsecase(s storage.EmployeeStorage) EmployeeUsecase {
	return EmployeeUsecase{storage: s}
}
