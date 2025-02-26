package handler

import (
    "log"
    "net/http"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
    usecase usecase.IEmployeeUsecase
}

func NewEmployeeHandler(u usecase.IEmployeeUsecase) *EmployeeHandler {
    return &EmployeeHandler{usecase: u}
}

func (h *EmployeeHandler) Auth(ctx *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
        log.Println(err.Error())
        return
    }

    token, err := h.usecase.Auth(ctx.Request.Context(), req.Username, req.Password)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        log.Println(err.Error())
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *EmployeeHandler) Info(ctx *gin.Context) {
    employeeID := ctx.GetInt("employeeID")

    info, err := h.usecase.Info(ctx.Request.Context(), employeeID)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        log.Println(err.Error())
        return
    }

    ctx.JSON(http.StatusOK, info)
}

func (h *EmployeeHandler) SendCoin(ctx *gin.Context) {
    var req struct {
        ToUser string `json:"toUser"`
        Amount int    `json:"amount"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
        log.Println(err.Error())
        return
    }

    senderID := ctx.GetInt("employeeID")

    if err, isInternal := h.usecase.SendCoin(ctx.Request.Context(), senderID, req.ToUser, req.Amount); err != nil {
        if isInternal {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        } else {
            ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        }
        log.Println(err.Error())
        
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "operation successful"})
}
