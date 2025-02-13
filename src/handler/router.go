package handler

import ( 
    "github.com/gin-gonic/gin"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/middleware"

)

func NewRouter(employeeUsecase usecase.EmployeeUsecase, merchUsecase usecase.MerchUsecase) *gin.Engine {
    router := gin.Default()
    router.Use(gin.Logger())

    employeeHandler := NewEmployeeHandler(employeeUsecase)
    if employeeHandler == nil {

    }

    merchHandler := NewMerchHandler(merchUsecase)
    if merchHandler == nil {

    }

    authGroup := router.Group("/api")
    {
        authGroup.POST("/auth", employeeHandler.Auth)
    }

    apiGroup := router.Group("/api")
    apiGroup.Use(middleware.AuthJWT()) // jwt
    {
        apiGroup.GET( "/info",      merchHandler.Info)
        apiGroup.GET( "/sendCoin",  employeeHandler.SendCoin)
        apiGroup.POST("/buy/:item", merchHandler.BuyItem)
    }

    return router
}
