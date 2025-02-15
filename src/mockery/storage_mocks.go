package mockery

import (
	"context"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
	"github.com/stretchr/testify/mock"
)

func unused(...interface{}) {}

type MockEmployeeStorage struct {
	mock.Mock
}

func (m *MockEmployeeStorage) GetEmployee(ctx context.Context, login string) (*entity.Employee, error) {
	unused(ctx, login)
	return nil, nil
}

func (m *MockEmployeeStorage) GetEmployeeOrRegister(ctx context.Context, login, password string) (*entity.Employee, error) {
	args := m.Called(ctx, login, password)
	return args.Get(0).(*entity.Employee), args.Error(1)
}

func (m *MockEmployeeStorage) GetEmployeeID(ctx context.Context, login string) (int, error) {
	args := m.Called(ctx, login)
	return args.Int(0), args.Error(1)
}

func (m *MockEmployeeStorage) GetEmployeeLogin(ctx context.Context, employeeID int) (string, error) {
	args := m.Called(ctx, employeeID)
	return args.String(0), args.Error(1)
}

func (m *MockEmployeeStorage) RegisterEmployee(ctx context.Context, login, password string) (*entity.Employee, error) {
	unused(ctx, login, password)
	return nil, nil
}

type MockManagementStorage struct {
	mock.Mock
}

func (m *MockManagementStorage) GetCoins(ctx context.Context, employeeID int) (int, error) {
	args := m.Called(ctx, employeeID)
	return args.Int(0), args.Error(1)
}

func (m *MockManagementStorage) GetInventory(ctx context.Context, employeeID int) ([]entity.InventoryItem, error) {
	args := m.Called(ctx, employeeID)
	return args.Get(0).([]entity.InventoryItem), args.Error(1)
}

func (m *MockManagementStorage) GetCoinHistory(ctx context.Context, employeeID int) (*entity.CoinHistory, error) {
	args := m.Called(ctx, employeeID)
	return args.Get(0).(*entity.CoinHistory), args.Error(1)
}

func (m *MockManagementStorage) ProvidePurchase(ctx context.Context, employeeID int, item string, cost int) error {
	unused(ctx, employeeID, item, cost)
	return nil
}

func (m *MockManagementStorage) ProvideOperation(ctx context.Context, senderID, receiverID, amount int) error {
	unused(ctx, senderID, receiverID, amount)
	return nil
}

func (m *MockManagementStorage) FetchReceivedHistory(ctx context.Context, receiverID int) ([]entity.RecvEntry, error) {
	unused(ctx, receiverID)
	return nil, nil
}

func (m *MockManagementStorage) FetchSentHistory(ctx context.Context, senderID int) ([]entity.SentEntry, error) {
	unused(ctx, senderID)
	return nil, nil
}
