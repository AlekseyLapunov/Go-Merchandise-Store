package usecase

import (
	"context"
	"errors"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
)

type EmployeeUsecase struct {
	storage           storage.EmployeeStorage
    managementStorage storage.ManagementStorage
}

func NewEmployeeUsecase(s storage.EmployeeStorage, c storage.ManagementStorage) EmployeeUsecase {
	return EmployeeUsecase{storage: s, managementStorage: c}
}

func (u *EmployeeUsecase) Auth(ctx context.Context, login, password string) (string, error) {
    employee, err := u.storage.GetEmployee(ctx, login)
    if err != nil {
        return "", errors.New("invalid credentials")
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "employeeID": employee.ID,
    })

    tokenString, err := token.SignedString([]byte("todo-gen-secret"))
    if err != nil {
        return "", errors.New("failed to generate token")
    }

    return tokenString, nil
}

func (u *EmployeeUsecase) Info(ctx context.Context, employeeID int) (*entity.InfoResponse, error) {
    balance, err := u.managementStorage.GetCoins(ctx, employeeID)
    if err != nil {
        return nil, err
    }

    inventory, err := u.managementStorage.GetInventory(ctx, employeeID)
    if err != nil {
        return nil, err
    }

    coinHistory, err := u.managementStorage.GetCoinHistory(ctx, employeeID)
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

func (u *EmployeeUsecase) SendCoin(ctx context.Context, senderID int, toUser string, amount int) (e error, isInternal bool) {
    if amount < 0 {
        return errors.New("negative coins amount prohibited"), false
    }

    balance, err := u.managementStorage.GetCoins(ctx, senderID)
    if err != nil {
        return err, true
    }

    if balance < amount {
        return errors.New("not enough coins"), false
    }

    receiver, err := u.storage.GetEmployee(ctx, toUser)
    if err != nil || receiver == nil {
        return errors.New("receiver not found"), true
    }

    if err := u.managementStorage.ProvideOperation(ctx, senderID, receiver.ID, amount); err != nil {
        return err, true
    }

    return nil, false
}
