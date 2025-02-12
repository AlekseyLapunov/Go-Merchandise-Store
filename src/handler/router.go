package handler

import ( 
    "github.com/gin-gonic/gin"
    "github.com/AlekseyLapunov/Go-Merchandise-Store/usecase"
)

func NewRouter() *gin.Engine {
    router := gin.Default()
    router.Use(gin.Logger())

    employeeHandler := NewEmployeeHandler(/*merchUsecase*/)
    if employeeHandler == nil {

    }

    merchHandler := NewMerchHandler(/*employeeUsecase*/)
    if merchHandler == nil {

    }

    // no jwt
    authGroup := router.Group("/api")
    {
        authGroup.POST("/auth", employeeHandler.Auth())
    }

    // jwt
    apiGroup := router.Group("/api")
    apiGroup.Use(/*middleware.Auth()*/)
    {
        apiGroup.GET("/info" /*, TODO merchHandler info */)
        apiGroup.GET("/sendCoin" /*, TODO merchHandler sendCoin*/ )
        apiGroup.POST("/buy/:item" /*, TODO employeeHandler sendCoin*/ )
    }

    return router
}
