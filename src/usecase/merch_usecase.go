package usecase

import (
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
	"context"
)

type MerchUsecase struct {
	coinStorage storage.CoinStorage
	storage     storage.MerchStorage
}

func NewMerchUsecase(s storage.MerchStorage, c storage.CoinStorage) MerchUsecase {
	return MerchUsecase{storage: s, coinStorage: c}
}

func (u MerchUsecase) BuyItem(ctx context.Context, employeeID int, item string) error {
    return nil
}
