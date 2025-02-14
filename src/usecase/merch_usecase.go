package usecase

import (
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
	"context"
	"errors"
)

type MerchUsecase struct {
	coinStorage storage.ManagementStorage
	storage     storage.MerchStorage
}

func NewMerchUsecase(s storage.MerchStorage, c storage.ManagementStorage) MerchUsecase {
	return MerchUsecase{storage: s, coinStorage: c}
}

func (u MerchUsecase) BuyItem(ctx context.Context, employeeID int, item string) error {
    cost, err := u.storage.GetMerchPrice(ctx, item)
    if err != nil {
        return errors.New("item not found")
    }

    balance, err := u.coinStorage.GetBalance(ctx, employeeID)
    if err != nil {
        return err
    }
    if balance < cost {
        return errors.New("not enough coins")
    }

	if err := u.coinStorage.ProvidePurchase(ctx, employeeID, item, cost); err != nil {
		return err
	}

    return nil
}
