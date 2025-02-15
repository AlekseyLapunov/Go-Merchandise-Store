package handler_test

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
    "strings"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/entity"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/handler"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/mockery"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestEmployeeHandler_Info(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockUsecase := new(mockery.MockEmployeeUsecase)
    r := gin.Default()
    h := handler.NewEmployeeHandler(mockUsecase)
    r.GET("/api/info", h.Info)

    employeeID := 0
    info := &entity.InfoResponse {
        Coins: 100,
        Inventory: []entity.InventoryItem{
            {Type: "t-shirt", Quantity: 10},
            {Type: "book",    Quantity: 5},
        },
        CoinHistory: entity.CoinHistory {
            Received: []entity.RecvEntry{{FromUser: "aleksey", Amount: 5}},
            Sent:     []entity.SentEntry{{ToUser:   "maksim",  Amount: 10}},
        },
    }

    mockUsecase.On("Info", mock.Anything, employeeID).Return(info, nil)

    req, _ := http.NewRequest(http.MethodGet, "/api/info", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    var response entity.InfoResponse
    json.Unmarshal(w.Body.Bytes(), &response)
    assert.Equal(t, info, &response)

    mockUsecase.ExpectedCalls = nil
    mockUsecase.On("Info", mock.Anything, employeeID).Return(nil, errors.New("not found"))

    req, _ = http.NewRequest(http.MethodGet, "/api/info", nil)
    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEmployeeHandler_SendCoin(t *testing.T) {
    query := "/api/sendCoin"
    gin.SetMode(gin.TestMode)
    mockUsecase := new(mockery.MockEmployeeUsecase)
    r := gin.Default()
    h := handler.NewEmployeeHandler(mockUsecase)
    r.POST(query, h.SendCoin)
    employeeID := 0

    // --- case successful
    mockUsecase.On("SendCoin", mock.Anything, employeeID, "maksim", 10).Return(nil, false)

    body := `{"toUser": "maksim", "amount": 10}`
    req, _ := http.NewRequest(http.MethodPost, query, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"status": "operation successful"}`, w.Body.String())

    // --- case bad request format
    req, _ = http.NewRequest(http.MethodPost, query, strings.NewReader(`{"toUser": 123}`))
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    assert.JSONEq(t, `{"error": "bad request format"}`, w.Body.String())

    // --- case coins insufficient
    mockUsecase.ExpectedCalls = nil
    mockUsecase.On("SendCoin", mock.Anything, employeeID, "maksim", 1000).Return(errors.New("not enough coins"), false)

    body = `{"toUser": "maksim", "amount": 1000}`
    req, _ = http.NewRequest(http.MethodPost, query, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusBadRequest, w.Code)
    assert.JSONEq(t, `{"error": "not enough coins"}`, w.Body.String())

    // --- case internal error
    mockUsecase.ExpectedCalls = nil
    mockUsecase.On("SendCoin", mock.Anything, employeeID, "maksim", 10).Return(errors.New("internal error"), true)

    body = `{"toUser": "maksim", "amount": 10}`
    req, _ = http.NewRequest(http.MethodPost, query, strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    req = setEmployeeID(req, employeeID)

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusInternalServerError, w.Code)
    assert.JSONEq(t, `{"error": "internal error"}`, w.Body.String())
}

func TestEmployeeHandler_Auth(t *testing.T) {
    gin.SetMode(gin.TestMode)
    mockUsecase := new(mockery.MockEmployeeUsecase)
    r := gin.Default()
    h := handler.NewEmployeeHandler(mockUsecase)
    r.POST("/api/auth", h.Auth)

    // --- case successful auth
    mockUsecase.On("Auth", mock.Anything, "aleksey", "securepassword").Return("valid_token", nil)

    body := `{"username": "aleksey", "password": "securepassword"}`
    req, _ := http.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
    assert.JSONEq(t, `{"token": "valid_token"}`, w.Body.String())

    // --- case invalid credentials
    mockUsecase.ExpectedCalls = nil
    mockUsecase.On("Auth", mock.Anything, "aleksey", "wrongpassword").Return("", errors.New("invalid credentials"))

    body = `{"username": "aleksey", "password": "wrongpassword"}`
    req, _ = http.NewRequest(http.MethodPost, "/api/auth", strings.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    w = httptest.NewRecorder()
    r.ServeHTTP(w, req)

    assert.Equal(t, http.StatusUnauthorized, w.Code)
    assert.JSONEq(t, `{"error": "invalid credentials"}`, w.Body.String())
}

func setEmployeeID(req *http.Request, id int) *http.Request {
    ctx := req.Context()
    ctx = context.WithValue(ctx, "employeeID", id)
    return req.WithContext(ctx)
}

