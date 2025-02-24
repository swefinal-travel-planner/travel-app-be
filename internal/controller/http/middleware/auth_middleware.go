package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	httpcommon "github.com/swefinal-travel-planner/travel-app-be/internal/domain/http_common"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/constants"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/jwt"
)

type AuthMiddleware struct {
	authService              service.AuthService
	authenticationRepository repository.AuthenticationRepository
}

func NewAuthMiddleware(authService service.AuthService, authenticationRepository repository.AuthenticationRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authService:              authService,
		authenticationRepository: authenticationRepository,
	}
}

func getAccessToken(c *gin.Context) (token string) {
	authHeader := c.GetHeader("Authorization")
	var accessToken string
	parts := strings.Split(authHeader, " ")
	if len(parts) == 2 {
		accessToken = parts[1]
	}
	return accessToken
}

func (a *AuthMiddleware) getRefreshToken(c *gin.Context) (token string) {
	return c.GetHeader("X-Refresh-Token")
}

func GetUserIdHelper(c *gin.Context) int64 {
	userId, exists := c.Get("userId")
	fmt.Println("check 1 ", userId)
	if !exists {
		return 0
	}
	return userId.(int64)
}

func (a *AuthMiddleware) VerifyToken(c *gin.Context) {
	// Get the JWT secret from the environment
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			},
		))
		return
	}

	// Retrieve the access token from the header or cookies
	accessToken := getAccessToken(c)

	claims, err := jwt.VerifyToken(accessToken, jwtSecret)
	if err == nil {
		// If the access token is valid, extract user Id and proceed
		if payload, ok := claims.Payload.(map[string]interface{}); ok {
			userId := int64(payload["id"].(float64))
			fmt.Println("check 2 ", userId)
			c.Set("userId", userId)
			c.Next()
			return
		}
	}

	// If the access token is expired, check the refresh token
	if err.Error() == httpcommon.ErrorMessage.TokenExpired {
		refreshToken := a.getRefreshToken(c)
		refreshClaims, errRf := jwt.VerifyToken(refreshToken, jwtSecret)
		if errRf != nil {
			// If the refresh token is invalid or expired, abort with unauthorized
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				},
			))
			return
		}

		// Extract user Id from refresh token claims
		payload, ok := refreshClaims.Payload.(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				},
			))
			return
		}
		userId := int64(payload["id"].(float64))

		// Check if the refresh token exists and is still valid in the database
		refreshTokenEntity, err := a.authService.ValidateRefreshToken(c, userId)
		if err != nil || refreshTokenEntity == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				},
			))
			return
		}

		// Generate a new access token
		newAccessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
			"id": userId,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.InternalServerError,
				},
			))
			return
		}

		// Return new access token in response
		c.JSON(http.StatusOK, gin.H{
			"accessToken": newAccessToken,
		})
		c.Abort()
		return
	}

	// For all other errors, abort with unauthorized
	c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
		httpcommon.Error{
			Message: err.Error(),
			Code:    httpcommon.ErrorResponseCode.Unauthorized,
		},
	))
}
