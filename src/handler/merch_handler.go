package handler

import (
    "log"
    "github.com/gin-gonic/gin"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "net/http"
)

type MerchHandler struct {
    usecase usecase.IMerchUsecase
}

func NewMerchHandler(u usecase.IMerchUsecase) *MerchHandler {
    return &MerchHandler{usecase: u} 
}

func (h *MerchHandler) BuyItem(ctx *gin.Context) {
    item := ctx.Param("item")
    if item == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request format"})
        return
    }

    employeeID := ctx.GetInt("employeeID")

    if err, isInternal := h.usecase.BuyItem(ctx.Request.Context(), employeeID, item); err != nil {
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
