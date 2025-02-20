package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(router *gin.Engine, studentHandler *StudentHandler, authAuthHandler *AuthHandler) {
	v1 := router.Group("/api/v1")
	{
		students := v1.Group("/students")
		{
			students.GET("/", studentHandler.GetAll)
		}
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authAuthHandler.Register)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
