package usecase

import (
	"context"
	"errors"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
)

type EmployeeUsecase struct {
	storage storage.EmployeeStorage
}

func NewEmployeeUsecase(s storage.EmployeeStorage) EmployeeUsecase {
	return EmployeeUsecase{storage: s}
}

func (u *EmployeeUsecase) Auth(ctx context.Context, username, password string) (string, error) {
	return "", nil
}

func (u *EmployeeUsecase) Info(ctx context.Context, employeeID int) (*entity.InfoResponse, error) {
    balance, err := u.storage.GetBalance(ctx, employeeID)
    if err != nil {
        return nil, err
    }

    inventory, err := u.storage.GetInventory(ctx, employeeID)
    if err != nil {
        return nil, err
    }

    coinHistory, err := u.storage.GetCoinHistory(ctx, employeeID)
    if err != nil {
        return nil, err
    }
	
	if coinHistory == nil {
		return nil, errors.New("coinHistory ptr was nil")
	}

    return &entity.InfoResponse{
        Coins:       balance,
        Inventory:   inventory,
        CoinHistory: *coinHistory,
    }, nil
}

func (u *EmployeeUsecase) SendCoin(ctx context.Context, senderID int, toUser string, amount int) error {
	return nil
}
