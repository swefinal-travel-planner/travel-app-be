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
	friendHandler *FriendHandler,
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

		invitationFriend := v1.Group("/invitation-friends")
		{
			invitationFriend.POST("", authMiddleware.VerifyAccessToken, invitationFriendHandler.AddFriend)
			invitationFriend.GET("", authMiddleware.VerifyAccessToken, invitationFriendHandler.GetAllInvitations)
			invitationFriend.PUT("/accept/:invitationId", authMiddleware.VerifyAccessToken, invitationFriendHandler.AcceptInvitation)
			invitationFriend.PUT("/deny/:invitationId", authMiddleware.VerifyAccessToken, invitationFriendHandler.DenyInvitation)
		}
		friend := v1.Group("/friends")
		{
			friend.GET("", authMiddleware.VerifyAccessToken, friendHandler.ViewFriends)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
