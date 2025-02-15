package mockery

import (
	"context"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
	"github.com/stretchr/testify/mock"
)

type MockEmployeeUsecase struct {
	mock.Mock
}

func (m *MockEmployeeUsecase) Auth(ctx context.Context, login string, password string) (string, error) {
    args := m.Called(ctx, login, password)
	
	return args.Get(0).(string), args.Error(1)
}

func (m *MockEmployeeUsecase) Info(ctx context.Context, employeeID int) (*entity.InfoResponse, error) {
	args := m.Called(ctx, employeeID)
	if args.Get(0) != nil {
		return args.Get(0).(*entity.InfoResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockEmployeeUsecase) SendCoin(ctx context.Context, senderID int, toUser string, amount int) (e error, isInternal bool) {
    args := m.Called(ctx, senderID, toUser, amount)
	
	return args.Error(0), args.Get(1).(bool)
}

type MockMerchUsecase struct {
	mock.Mock
}

func (m *MockMerchUsecase) BuyItem(ctx context.Context, employeeID int, item string) (err error, isInternal bool) {
    args := m.Called(ctx, employeeID, item)
	
	return args.Error(0), args.Get(1).(bool)
}
