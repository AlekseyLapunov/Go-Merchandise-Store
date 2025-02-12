package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
)

type EmployeeHandler struct {
	usecase *usecase.EmployeeUsecase
}

func NewEmployeeHandler(u *usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

func (h *EmployeeHandler) Auth(ctx *gin.Context) {

}

func (h *EmployeeHandler) SendCoin(ctx *gin.Context) {

}
