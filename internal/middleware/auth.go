package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"h-two/internal/errors"
	"net/http"
	"os"
	"strings"
	"time"
)

func AuthMiddleware(c *gin.Context) {
	tokenStr := c.GetHeader("Authorization")
	if tokenStr == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{
			Message:    "Invalid Token",
			StatusCode: http.StatusUnauthorized,
			Status:     errors.UnAuthorized,
		})
		return
	}
	if !strings.HasPrefix(tokenStr, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{
			Message:    "Invalid Token",
			StatusCode: http.StatusUnauthorized,
			Status:     errors.UnAuthorized,
		})
		return
	}
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ApiError{
			Message:    "Internal server error",
			StatusCode: http.StatusInternalServerError,
			Status:     errors.InternalServerError,
		})
		return
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{
			Message:    "Unauthorized",
			StatusCode: http.StatusUnauthorized,
			Status:     errors.UnAuthorized,
		})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{
					Message:    "Token has expired",
					StatusCode: http.StatusUnauthorized,
					Status:     errors.UnAuthorized,
				})
				return
			}
		} else {
			c.AbortWithStatusJSON(http.StatusBadRequest, errors.ApiError{
				Message:    "Invalid Token",
				StatusCode: http.StatusBadRequest,
				Status:     errors.UnAuthorized,
			})
			return
		}
		c.Set("userId", claims["userId"])

	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, errors.ApiError{
			Message:    "Unauthorized",
			StatusCode: http.StatusUnauthorized,
			Status:     errors.UnAuthorized,
		})
		return
	}

	// Call the next handler
	c.Next()
}
