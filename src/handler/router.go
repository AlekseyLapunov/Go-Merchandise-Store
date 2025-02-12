package handler

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
    router := gin.Default()
    router.Use(gin.Logger())

    // employeeHandler := TODO
    // merchHandler    := TODO

    // no jwt
    authGroup := router.Group("/api")
    {
        authGroup.POST("/auth" /*, TODO employeeHandler auth */)
    }

    // jwt
    apiGroup := router.Group("/api")
    apiGroup.Use(/*middleware.Auth()*/)
    {
        apiGroup.GET("/info" /*, TODO merchHandler info */)
        apiGroup.GET("/sendCoin" /*, TODO merchHandler sendCoin*/ )
        apiGroup.POST("/buy/:item" /*, TODO merchHandler sendCoin*/ )
    }

    return router
}
