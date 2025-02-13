package usecase

import (
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
	"context"
)

type EmployeeUsecase struct {
	storage storage.EmployeeStorage
}

func NewEmployeeUsecase(s storage.EmployeeStorage) EmployeeUsecase {
	return EmployeeUsecase{storage: s}
}

func (u EmployeeUsecase) Auth(ctx context.Context, username, password string) (string, error) {
	return "", nil
}
