package middleware

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "strings"
)

func AuthJWT() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        header := ctx.GetHeader("Authorization")
        if header == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
            ctx.Abort()
            return
        }

        authParts := strings.Split(header, " ")
        if len(authParts) != 2 || authParts[0] != "Bearer" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            ctx.Abort()
            return
        }

        token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
            return []byte("todo-gen-secret"), nil
        })
        if err != nil || !token.Valid {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            ctx.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            ctx.Abort()
            return
        }

        employeeID, ok := claims["employeeID"].(int)
        if !ok {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid employee ID in token"})
            ctx.Abort()
            return
        }

        ctx.Set("employeeID", employeeID)

        ctx.Next()
    }
}

