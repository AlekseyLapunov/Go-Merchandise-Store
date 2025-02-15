package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
    "strconv"
    "log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
        if len(authParts) != 2 || authParts[0] != "BearerAuth" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            ctx.Abort()
            return
        }

        secretJWT, err := FetchSecretJWT()
        if err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong on our side"})
            ctx.Abort()
            log.Println(err.Error())
            return
        }

        token, err := jwt.Parse(authParts[1], func(token *jwt.Token) (interface{}, error) {
            return []byte(secretJWT), nil
        })
        if err != nil || !token.Valid {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            ctx.Abort()
            log.Println(err.Error())
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            ctx.Abort()
            return
        }

        employeeIDstr, ok := claims["employeeID"].(string)
        if !ok {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid employee ID in token"})
            ctx.Abort()
            return
        }

        employeeID, err := strconv.Atoi(employeeIDstr)
        if err != nil {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid employee ID format in token"})
            ctx.Abort()
            log.Println(err.Error())
            return
        }

        ctx.Set("employeeID", employeeID)

        ctx.Next()
    }
}

func FetchSecretJWT() (string, error) {
    const ENV_JWT_SECRET = "JWT_SECRET"

    jwt_secret := os.Getenv(ENV_JWT_SECRET)
    if jwt_secret == "" {
        return "", fmt.Errorf("\"%s\" environment variable is not set", ENV_JWT_SECRET)
    }
    return jwt_secret, nil
}

