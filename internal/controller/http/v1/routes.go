package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
)

func MapRoutes(router *gin.Engine, authHandler *AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/test", authMiddleware.VerifyAccessToken, authHandler.Test)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
