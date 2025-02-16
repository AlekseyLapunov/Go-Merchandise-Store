package handler

import ( 
    "github.com/gin-gonic/gin"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/usecase"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/src/middleware"

)

func NewRouter(employeeUsecase usecase.IEmployeeUsecase, merchUsecase usecase.IMerchUsecase) *gin.Engine {
    router := gin.New()
    gin.SetMode(gin.ReleaseMode)

    employeeHandler := NewEmployeeHandler(employeeUsecase)
    merchHandler := NewMerchHandler(merchUsecase)

    authGroup := router.Group("/api")
    {
        authGroup.POST("/auth", employeeHandler.Auth)
    }

    apiGroup := router.Group("/api")
    apiGroup.Use(middleware.AuthJWT()) // jwt
    {
        apiGroup.GET( "/info",      employeeHandler.Info)
        apiGroup.POST("/sendCoin",  employeeHandler.SendCoin)
        apiGroup.GET( "/buy/:item", merchHandler.BuyItem)
    }

    return router
}
