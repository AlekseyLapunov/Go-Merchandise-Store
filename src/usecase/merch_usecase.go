package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
)

type MerchUsecase struct {
    storage           storage.IMerchStorage
    managementStorage storage.IManagementStorage
}

func NewMerchUsecase(s storage.IMerchStorage, c storage.IManagementStorage) MerchUsecase {
    return MerchUsecase{storage: s, managementStorage: c}
}

func (u MerchUsecase) BuyItem(ctx context.Context, employeeID int, item string) (err error, isInternal bool) {
    errString := "problem buying item (on our side)"

    cost, err := u.storage.GetMerchCost(ctx, item)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.New("wrong item name"), false
        }

        log.Println(err)
        return errors.New(errString), true
    }

    balance, err := u.managementStorage.GetCoins(ctx, employeeID)
    if err != nil {
        log.Println(err)
        return errors.New(errString), true
    }
    if balance < cost {
        return errors.New("not enough coins"), false
    }

    if err := u.managementStorage.ProvidePurchase(ctx, employeeID, item, cost); err != nil {
        log.Println(err)
        return errors.New(errString), true
    }

    return nil, false
}
