package middleware

import (
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/error_utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/env"
	"github.com/swefinal-travel-planner/travel-app-be/internal/utils/jwt"
)

type AuthMiddleware struct {
	authService              service.AuthService
	authenticationRepository repository.AuthenticationRepository
	userRepository           repository.UserRepository
}

func NewAuthMiddleware(
	authService service.AuthService,
	authenticationRepository repository.AuthenticationRepository,
	userRepository repository.UserRepository,
) *AuthMiddleware {
	return &AuthMiddleware{
		authService:              authService,
		authenticationRepository: authenticationRepository,
		userRepository:           userRepository,
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

func GetUserIdHelper(c *gin.Context) int64 {
	userId, exists := c.Get("userId")
	if !exists {
		return 0
	}
	return userId.(int64)
}

func (a *AuthMiddleware) VerifyAccessToken(c *gin.Context) {
	// Get the JWT secret from the environment
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		log.Error("AuthMiddleware.VerifyAccessToken Error getting JWT secret: " + err.Error())
		statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.INTERNAL_SERVER_ERROR, "")
		c.AbortWithStatusJSON(statusCode, errResponse)
		return
	}

	// Retrieve the access token from the header
	accessToken := getAccessToken(c)

	claims, err := jwt.VerifyToken(accessToken, jwtSecret)
	if err == nil {
		// If the access token is valid, extract user Id and proceed
		if payload, ok := claims.Payload.(map[string]interface{}); ok {
			userId := int64(payload["id"].(float64))
			c.Set("userId", userId)
			c.Next()
			return
		}
	}

	statusCode, errResponse := error_utils.ErrorCodeToHttpResponse(error_utils.ErrorCode.ACCESS_TOKEN_INVALID, "accessToken")
	c.AbortWithStatusJSON(statusCode, errResponse)
}
