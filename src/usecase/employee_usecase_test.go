package usecase_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/mockery"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth_Success(t *testing.T) {
	mockEmployeeStorage := new(mockery.MockEmployeeStorage)
	mockManagementStorage := new(mockery.MockManagementStorage)
	u := usecase.NewEmployeeUsecase(mockEmployeeStorage, mockManagementStorage)

	ctx := context.Background()
	login := "test_user"
	password := "test_password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	employee := &entity.Employee{ID: 1, Login: login, Password: string(hashedPassword)}

	mockEmployeeStorage.On("GetEmployeeOrRegister", ctx, login, password).Return(employee, nil)

	os.Setenv("JWT_SECRET", "test-secret")
	token, err := u.Auth(ctx, login, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockEmployeeStorage.AssertExpectations(t)
}

func TestAuth_WrongPassword(t *testing.T) {
	mockEmployeeStorage := new(mockery.MockEmployeeStorage)
	mockManagementStorage := new(mockery.MockManagementStorage)
	u := usecase.NewEmployeeUsecase(mockEmployeeStorage, mockManagementStorage)

	ctx := context.Background()
	login := "test_user"
	password := "wrongpassword"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	employee := &entity.Employee{ID: 1, Login: login, Password: string(hashedPassword)}

	mockEmployeeStorage.On("GetEmployeeOrRegister", ctx, login, password).Return(employee, nil)

	os.Setenv("JWT_SECRET", "test-secret")
	token, err := u.Auth(ctx, login, password)

	assert.Error(t, err)
	assert.Empty(t, token)
	mockEmployeeStorage.AssertExpectations(t)
}

func TestEmployeeUsecase_Info(t *testing.T) {
	mockStorage := new(mockery.MockManagementStorage)
	u := usecase.NewEmployeeUsecase(nil, mockStorage)
	ctx := context.Background()
	employeeID := 1
	mockStorage.On("GetCoins", ctx, employeeID).Return(100, nil)
	mockStorage.On("GetInventory", ctx, employeeID).Return(
        []entity.InventoryItem{
            {Type: "t-shirt", Quantity: 10},
            {Type: "book",    Quantity: 5 },
        },
        nil,
    )
	mockStorage.On("GetCoinHistory", ctx, employeeID).Return(
		&entity.CoinHistory {
			Received: []entity.RecvEntry{{FromUser: "aleksey", Amount: 5}, {FromUser: "sveta", Amount: 45}},
            Sent:     []entity.SentEntry{{ToUser: "maksim", Amount: 10}},
		},
		nil,
	)

	info, err := u.Info(ctx, employeeID)
	assert.NoError(t, err)
	assert.Equal(t, 100, info.Coins)
	assert.Equal(t, info.Inventory[0].Type,     "t-shirt")
	assert.Equal(t, info.Inventory[0].Quantity, 10)
	assert.Equal(t, info.Inventory[1].Type,     "book")
	assert.Equal(t, info.Inventory[1].Quantity, 5)
	assert.Len(t, info.CoinHistory.Received, 2)
    assert.Len(t, info.CoinHistory.Sent, 1)

    mockStorage.ExpectedCalls = nil
	mockErrStr := "internal db error"
	mockStorage.On("GetCoins", ctx, employeeID).Return(0, errors.New(mockErrStr))
	_, err = u.Info(ctx, employeeID)
	assert.Error(t, err)

	mockStorage.On("GetInventory", ctx, employeeID).Return(nil, errors.New(mockErrStr))
	_, err = u.Info(ctx, employeeID)
	assert.Error(t, err)

	mockStorage.On("GetCoinHistory", ctx, employeeID).Return(nil, errors.New(mockErrStr))
	_, err = u.Info(ctx, employeeID)
	assert.Error(t, err)

	mockStorage.On("GetCoinHistory", ctx, employeeID).Return((*entity.CoinHistory)(nil), nil)
	_, err = u.Info(ctx, employeeID)
	assert.Error(t, err)
}

func TestEmployeeUsecase_SendCoin(t *testing.T) {
	mockManagementStorage := new(mockery.MockManagementStorage)
	mockEmployeeStorage := new(mockery.MockEmployeeStorage)
	u := usecase.NewEmployeeUsecase(mockEmployeeStorage, mockManagementStorage)
	ctx := context.Background()
	senderID := 1
	receiverLogin := "receiver"
	receiverID := 2
	amount := 50

	mockManagementStorage.On("GetCoins", ctx, senderID).Return(100, nil)
	mockEmployeeStorage.On("GetEmployeeLogin", ctx, senderID).Return("sender", nil)
	mockEmployeeStorage.On("GetEmployeeID", ctx, receiverLogin).Return(receiverID, nil)
	mockManagementStorage.On("ProvideOperation", ctx, senderID, receiverID, amount).Return(nil)

	err, isInternal := u.SendCoin(ctx, senderID, receiverLogin, amount)
	assert.NoError(t, err)
	assert.False(t, isInternal)

    mockManagementStorage.ExpectedCalls = nil
	mockManagementStorage.On("GetCoins", ctx, senderID).Return(40, nil)
	err, _ = u.SendCoin(ctx, senderID, receiverLogin, amount)
	assert.Error(t, err)

	mockEmployeeStorage.On("GetEmployeeLogin", ctx, senderID).Return(receiverLogin, nil)
	err, _ = u.SendCoin(ctx, senderID, receiverLogin, amount)
	assert.Error(t, err)

    mockManagementStorage.ExpectedCalls = nil
	mockManagementStorage.On("GetCoins", ctx, senderID).Return(0, errors.New("db error"))
	err, isInternal = u.SendCoin(ctx, senderID, receiverLogin, amount)
	assert.Error(t, err)
	assert.True(t, isInternal)

	mockEmployeeStorage.On("GetEmployeeID", ctx, receiverLogin).Return(0, errors.New("receiver not found"))
	err, _ = u.SendCoin(ctx, senderID, receiverLogin, amount)
	assert.Error(t, err)

	mockManagementStorage.On("ProvideOperation", ctx, senderID, receiverID, amount).Return(errors.New("operation error"))
	err, isInternal = u.SendCoin(ctx, senderID, receiverLogin, amount)
	assert.Error(t, err)
	assert.True(t, isInternal)
}
