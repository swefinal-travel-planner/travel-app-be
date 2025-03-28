package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
)

func MapRoutes(router *gin.Engine,
	authHandler *AuthHandler,
	invitationFriendHandler *InvitationFriendHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google-login", authHandler.FirebaseLogin)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/send-otp", authHandler.SendOTPToMail)
			auth.POST("/verify-otp", authHandler.VerifyOTP)
			auth.POST("/forgot-password", authHandler.SetPassword)

			auth.GET("/test", authMiddleware.VerifyAccessToken, authHandler.Test)
		}

		invitationFriend := v1.Group("/invitation-friend")
		{
			invitationFriend.POST("/add", authMiddleware.VerifyAccessToken, invitationFriendHandler.AddFriend)
			invitationFriend.GET("/get-all", authMiddleware.VerifyAccessToken, invitationFriendHandler.GetAllInvitations)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
