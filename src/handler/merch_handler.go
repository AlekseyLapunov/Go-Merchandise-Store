package handler

import ( 
	"github.com/gin-gonic/gin"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
)

type MerchHandler struct {
	usecase usecase.MerchUsecase
}

func NewMerchHandler(u usecase.MerchUsecase) *MerchHandler {
	return &MerchHandler{usecase: u} 
}

func (h *MerchHandler) Info(ctx *gin.Context) {

}

func (h *MerchHandler) BuyItem(ctx *gin.Context) {

}
