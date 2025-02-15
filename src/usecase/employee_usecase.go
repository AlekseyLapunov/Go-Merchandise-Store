package usecase

import (
    "context"
    "database/sql"
    "errors"
    "log"
    "strconv"

    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/middleware"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/storage"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type IEmployeeUsecase interface {
    Auth(ctx context.Context, login, password string) (string, error)
    Info(ctx context.Context, employeeID int) (*entity.InfoResponse, error)
    SendCoin(ctx context.Context, senderID int, toUser string, amount int) (e error, isInternal bool)
}

type EmployeeUsecase struct {
    storage           storage.IEmployeeStorage
    managementStorage storage.IManagementStorage
}

func NewEmployeeUsecase(s storage.IEmployeeStorage, c storage.IManagementStorage) EmployeeUsecase {
    return EmployeeUsecase{storage: s, managementStorage: c}
}

func (u *EmployeeUsecase) Auth(ctx context.Context, login, password string) (string, error) {
    errString := "auth went wrong on our side"

    employee, err := u.storage.GetEmployeeOrRegister(ctx, login, password)
    if err != nil {
        log.Println(err)
        return "", errors.New(errString)
    }
    
    if err := bcrypt.CompareHashAndPassword([]byte(employee.Password), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "employeeID": strconv.Itoa(employee.ID),
    })

    secretJWT, err := middleware.FetchSecretJWT()
    if err != nil {
        log.Println(err)
        return "", errors.New(errString)
    }

    tokenString, err := token.SignedString([]byte(secretJWT))
    if err != nil {
        log.Println(err)
        return "", errors.New(errString)
    }

    return tokenString, nil
}

func (u *EmployeeUsecase) Info(ctx context.Context, employeeID int) (*entity.InfoResponse, error) {
    errString := "trouble getting information (on our side)"

    balance, err := u.managementStorage.GetCoins(ctx, employeeID)
    if err != nil {
        log.Println(err.Error())
        return nil, errors.New(errString)
    }

    inventory, err := u.managementStorage.GetInventory(ctx, employeeID)
    if err != nil {
        log.Println(err.Error())
        return nil, errors.New(errString)
    }

    coinHistory, err := u.managementStorage.GetCoinHistory(ctx, employeeID)
    if err != nil {
        log.Println(err.Error())
        return nil, errors.New(errString)
    }
    
    if coinHistory == nil {
        log.Println("coinHistory ptr was nil")
        return nil, errors.New(errString)
    }

    return &entity.InfoResponse{
        Coins:       balance,
        Inventory:   inventory,
        CoinHistory: *coinHistory,
    }, nil
}

func (u *EmployeeUsecase) SendCoin(ctx context.Context, senderID int, toUser string, amount int) (e error, isInternal bool) {
    errString := "trouble sending coins (on our side)"

    if amount < 0 {
        return errors.New("negative coins amount prohibited"), false
    }

    balance, err := u.managementStorage.GetCoins(ctx, senderID)
    if err != nil {
        log.Println(err.Error())
        return errors.New(errString), true
    }

    if balance < amount {
        return errors.New("not enough coins"), false
    }

    senderLogin, err := u.storage.GetEmployeeLogin(ctx, senderID)
    if err != nil {
        return errors.New(errString), true
    }

    if senderLogin == toUser {
        return errors.New("can't send coins to yourself"), false
    }

    receiverID, err := u.storage.GetEmployeeID(ctx, toUser)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return errors.New("receiver user name not found"), false
        } else {
            return errors.New(errString), true
        }
    }

    if err := u.managementStorage.ProvideOperation(ctx, senderID, receiverID, amount); err != nil {
        log.Println(err.Error())
        return errors.New(errString), true
    }

    return nil, false
}
