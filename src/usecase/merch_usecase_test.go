package usecase_test

import (
	"context"
	"errors"
	"testing"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/mockery"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
	"github.com/stretchr/testify/assert"
)

func TestMerchUsecase_BuyItem(t *testing.T) {
	mockManagementStorage := new(mockery.MockManagementStorage)
	mockMerchStorage      := new(mockery.MockMerchStorage)

	u := usecase.NewMerchUsecase(mockMerchStorage, mockManagementStorage)
	ctx := context.Background()
	employeeID := 1
	itemType := "powerbank"
    cost     := 200

	mockManagementStorage.On("GetCoins", ctx, employeeID).Return(500, nil)
    mockMerchStorage.On("GetMerchCost", ctx, itemType).Return(cost, nil)
    mockManagementStorage.On("ProvidePurchase", ctx, employeeID, itemType, cost).Return(nil)

	err, isInternal := u.BuyItem(ctx, employeeID, itemType)
	assert.NoError(t, err)
	assert.False(t, isInternal)

    mockManagementStorage.ExpectedCalls = nil
	mockManagementStorage.On("GetCoins", ctx, employeeID).Return(199, nil)
    mockMerchStorage.On("GetMerchCost", ctx, itemType).Return(cost, nil)
	err, _ = u.BuyItem(ctx, employeeID, itemType)
	assert.Error(t, err)
	assert.False(t, isInternal)

    mockManagementStorage.ExpectedCalls = nil
	mockManagementStorage.On("GetCoins", ctx, employeeID).Return(0, errors.New("db error"))
	err, isInternal = u.BuyItem(ctx, employeeID, itemType)
	assert.Error(t, err)
	assert.True(t, isInternal)

    mockManagementStorage.ExpectedCalls = nil
	mockManagementStorage.On("GetCoins", ctx, employeeID).Return(500, nil)
    mockMerchStorage.On("GetMerchCost", ctx, itemType).Return(cost, nil)
    mockManagementStorage.On("ProvidePurchase", ctx, employeeID, itemType, cost).Return(errors.New("db error"))
	err, isInternal = u.BuyItem(ctx, employeeID, itemType)
	assert.Error(t, err)
	assert.True(t, isInternal)
}