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
	userHandler *UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	healHandler *HealthHandler,
	notificationHandler *NotificationHandler,
	tripHandler *TripHandler,
	invitationTripHandler *InvitationTripHandler,
) {
	v1 := router.Group("/api/v1")
	{
		health := v1.Group("/health")
		{
			health.GET("", healHandler.Check)
		}
		notification := v1.Group("/notifications")
		{
			notification.POST("/test", notificationHandler.TestNotification)
			notification.GET("", authMiddleware.VerifyAccessToken, notificationHandler.GetAllNotification)
			notification.POST("/:id/seen", authMiddleware.VerifyAccessToken, notificationHandler.SeenNotification)
		}
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/google-login", authHandler.FirebaseLogin)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/register/send-otp", authHandler.SendOTPToEmailForRegister)
			auth.POST("/register/verify-otp", authHandler.VerifyOTPForRegister)
			auth.POST("/reset-password/send-otp", authHandler.SendOTPToEmailForResetPassword)
			auth.POST("/reset-password/verify-otp", authHandler.VerifyOTPForResetPassword)
			auth.POST("/reset-password", authHandler.SetPassword)

			auth.GET("/test", authMiddleware.VerifyAccessToken, authHandler.Test)
		}

		invitationFriend := v1.Group("/invitation-friends")
		{
			invitationFriend.POST("", authMiddleware.VerifyAccessToken, invitationFriendHandler.AddFriend)
			invitationFriend.GET("/received", authMiddleware.VerifyAccessToken, invitationFriendHandler.GetAllReceivedInvitations)
			invitationFriend.GET("/requested", authMiddleware.VerifyAccessToken, invitationFriendHandler.GetAllRequestedInvitations)
			invitationFriend.PUT("/accept/:invitationId", authMiddleware.VerifyAccessToken, invitationFriendHandler.AcceptInvitation)
			invitationFriend.PUT("/deny/:invitationId", authMiddleware.VerifyAccessToken, invitationFriendHandler.DenyInvitation)
			invitationFriend.DELETE("/:invitationId", authMiddleware.VerifyAccessToken, invitationFriendHandler.WithdrawInvitation)
		}
		friend := v1.Group("/friends")
		{
			friend.GET("", authMiddleware.VerifyAccessToken, friendHandler.ViewFriends)
			friend.DELETE("/:friendId", authMiddleware.VerifyAccessToken, friendHandler.RemoveFriend)
		}
		user := v1.Group("/users")
		{
			user.GET("/", authMiddleware.VerifyAccessToken, userHandler.SearchUser)
			user.PUT("/notification-token", authMiddleware.VerifyAccessToken, userHandler.UpdateNotificationToken)
			user.PATCH("/me", authMiddleware.VerifyAccessToken, userHandler.UpdateProfile)
		}
		trip := v1.Group("/trips")
		{
			trip.POST("", authMiddleware.VerifyAccessToken, tripHandler.CreateTripManually)
			trip.GET("", authMiddleware.VerifyAccessToken, tripHandler.GetAllTrips)
			trip.GET("/:tripId", authMiddleware.VerifyAccessToken, tripHandler.GetTrip)
			trip.PATCH("/:tripId", authMiddleware.VerifyAccessToken, tripHandler.UpdateTrip)
			trip.POST("/:tripId/trip-items", authMiddleware.VerifyAccessToken, tripHandler.CreateTripItems)
			trip.GET("/:tripId/trip-items", authMiddleware.VerifyAccessToken, tripHandler.GetTripItems)
			trip.POST("/ai", authMiddleware.VerifyAccessToken, tripHandler.CreateTripByAI)
		}
		tripInvitation := v1.Group("/invitation-trips")
		{
			tripInvitation.POST("", authMiddleware.VerifyAccessToken, invitationTripHandler.SendInvitation)
			tripInvitation.PUT("/accept/:invitationId", authMiddleware.VerifyAccessToken, invitationTripHandler.AcceptInvitation)
			tripInvitation.PUT("/deny/:invitationId", authMiddleware.VerifyAccessToken, invitationTripHandler.DenyInvitation)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
