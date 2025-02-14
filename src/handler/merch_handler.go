package handler

import ( 
    "github.com/gin-gonic/gin"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "net/http"
)

type MerchHandler struct {
    usecase usecase.MerchUsecase
}

func NewMerchHandler(u usecase.MerchUsecase) *MerchHandler {
    return &MerchHandler{usecase: u} 
}

func (h *MerchHandler) BuyItem(ctx *gin.Context) {
    item := ctx.Param("item")
    if item == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
        return
    }

    employeeID := ctx.GetInt("employeeID")

    if err := h.usecase.BuyItem(ctx.Request.Context(), employeeID, item); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "Operation successful"})
}
