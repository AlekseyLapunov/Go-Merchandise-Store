package usecase

import (
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
	"context"
)

type MerchUsecase struct {
	storage storage.MerchStorage
}

func NewMerchUsecase(s storage.MerchStorage) MerchUsecase {
	return MerchUsecase{storage: s}
}

func (u MerchUsecase) BuyItem(ctx context.Context, employeeID int, item string) error {
	return nil
}
