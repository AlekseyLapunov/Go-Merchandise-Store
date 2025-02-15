package usecase_test

import (
	"context"
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
