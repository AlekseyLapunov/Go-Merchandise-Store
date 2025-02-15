package handler_test

import (
    "errors"
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/handler"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/mockery"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestMerchHandler_BuyItem(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockUsecase := new(mockery.MockMerchUsecase)
    r := gin.Default()
    h := handler.NewMerchHandler(mockUsecase)
    r.POST("/api/buy/:item", h.BuyItem)

    employeeID := 0
    merchName := "book"

    // --- case successful purchase
    mockUsecase.On("BuyItem", mock.Anything, employeeID, merchName).Return(nil, false)

    req, _ := http.NewRequest(http.MethodPost, "/api/buy/book", nil)
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"status": "operation successful"}`, w.Body.String())

    // --- case no such merch
    mockUsecase.ExpectedCalls = nil
    mockUsecase.On("BuyItem", mock.Anything, employeeID, merchName).Return(errors.New("item not found"), false)

    req, _ = http.NewRequest(http.MethodPost, "/api/buy/book", nil)
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    assert.JSONEq(t, `{"error": "item not found"}`, w.Body.String())

    // --- case internal error
    expectedString := "something went wrong on our end"
    mockUsecase.ExpectedCalls = nil
    mockUsecase.On("BuyItem", mock.Anything, employeeID, merchName).Return(errors.New(expectedString), true)

    req, _ = http.NewRequest(http.MethodPost, "/api/buy/book", nil)
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    assert.JSONEq(t, fmt.Sprintf(`{"error": "%s"}`, expectedString), w.Body.String())

    // --- case missing item parameter
    req, _ = http.NewRequest(http.MethodPost, "/api/buy/", nil)
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusNotFound, w.Code)
}

