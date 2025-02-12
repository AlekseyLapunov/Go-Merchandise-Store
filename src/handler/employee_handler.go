package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/AlekseyLapunov/Go-Merchandise-Store/usecase"
)

type EmployeeHandler struct {
	usecase usecase.EmployeeUsecase
}

func NewEmployeeHandler(u usecase.EmployeeUsecase) *EmployeeHandler {
	return &EmployeeHandler{usecase: u}
}

func (h *EmployeeHandler) Auth(c *gin.Context) {

}

func (h *EmployeeHandler) SendCoin(c *gin.Context) {
	
}
