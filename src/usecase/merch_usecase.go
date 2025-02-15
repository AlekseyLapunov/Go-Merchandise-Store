package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
)

type MerchUsecase struct {
    storage           storage.MerchStorage
    managementStorage storage.ManagementStorage
}

func NewMerchUsecase(s storage.MerchStorage, c storage.ManagementStorage) MerchUsecase {
    return MerchUsecase{storage: s, managementStorage: c}
}

func (u MerchUsecase) BuyItem(ctx context.Context, employeeID int, item string) (err error, isInternal bool) {
    cost, err := u.storage.GetMerchCost(ctx, item)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.New("wrong item name"), false
        }
        return err, true
    }

    balance, err := u.managementStorage.GetCoins(ctx, employeeID)
    if err != nil {
        return err, true
    }
    if balance < cost {
        return errors.New("not enough coins"), false
    }

    if err := u.managementStorage.ProvidePurchase(ctx, employeeID, item, cost); err != nil {
        return err, true
    }

    return nil, false
}
