package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
	"net/http"
)

type EmployeeHandler struct {
	usecase usecase.EmployeeUsecase
}

func NewEmployeeHandler(u usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

func (h *EmployeeHandler) Auth(ctx *gin.Context) {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }

    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    token, err := h.usecase.Auth(ctx.Request.Context(), req.Username, req.Password)
    if err != nil {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *EmployeeHandler) SendCoin(ctx *gin.Context) {

}
