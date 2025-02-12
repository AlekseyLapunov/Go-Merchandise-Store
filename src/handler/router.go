package handler

import ( 
    "github.com/gin-gonic/gin"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
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

    // no jwt
    authGroup := router.Group("/api")
    {
        authGroup.POST("/auth", employeeHandler.Auth)
    }

    // jwt
    apiGroup := router.Group("/api")
    apiGroup.Use(/*middleware.Auth()*/)
    {
        apiGroup.GET("/info",       merchHandler.Info)
        apiGroup.GET("/sendCoin",   employeeHandler.SendCoin)
        apiGroup.POST("/buy/:item", merchHandler.BuyItem)
    }

    return router
}
