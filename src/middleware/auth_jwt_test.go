package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlekseyLapunov/Go-Merchandise-Store/src/middleware"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// Мок FetchSecretJWT для тестов
func mockFetchSecretJWT(secret string, err error) func() (string, error) {
    return func() (string, error) {
        return secret, err
    }
}

type TestCase struct {
    name           string
    authHeader     string
    secretJWT      string
    tokenClaims    jwt.MapClaims
    expectedStatus int
}

func TestAuthJWT(t *testing.T) {
    gin.SetMode(gin.TestMode)

    testcases := []TestCase{
        {
            name:           "No Authorization Header",
            authHeader:     "",
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "Invalid Authorization Header Format",
            authHeader:     "InvalidFormatToken",
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "Invalid Bearer Prefix",
            authHeader:     "BearerWrong abc.def.ghi",
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "FetchSecretJWT Error",
            authHeader:     "BearerAuth abc.def.ghi",
            expectedStatus: http.StatusInternalServerError,
        },
    }

    for _, tc := range testcases {
        t.Run(tc.name, func(t *testing.T) {
            tokenString := tc.authHeader
            if tc.tokenClaims != nil {
                token := jwt.NewWithClaims(jwt.SigningMethodHS256, tc.tokenClaims)
                signedToken, _ := token.SignedString([]byte(tc.secretJWT))
                tokenString = "BearerAuth " + signedToken
            }

            req := httptest.NewRequest(http.MethodGet, "/", nil)
            req.Header.Set("Authorization", tokenString)

            w := httptest.NewRecorder()

            r := gin.Default()
            r.Use(middleware.AuthJWT())
            r.GET("/", func(ctx *gin.Context) {
                ctx.JSON(http.StatusOK, gin.H{"message": "success"})
            })

            r.ServeHTTP(w, req)
            assert.Equal(t, tc.expectedStatus, w.Code)
        })
    }
}