package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	v1 "github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/v1"
)

type Server struct {
	authAuthHandler         *v1.AuthHandler
	invitationFriendHandler *v1.InvitationFriendHandler
	friendHandler           *v1.FriendHandler
	userHandler             *v1.UserHandler
	authMiddleware          *middleware.AuthMiddleware
}

func NewServer(authAuthHandler *v1.AuthHandler,
	invitationFriendHandler *v1.InvitationFriendHandler,
	friendHandler *v1.FriendHandler,
	userHandler *v1.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
) *Server {
	return &Server{
		authAuthHandler:         authAuthHandler,
		invitationFriendHandler: invitationFriendHandler,
		friendHandler:           friendHandler,
		userHandler:             userHandler,
		authMiddleware:          authMiddleware,
	}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)

	v1.MapRoutes(router, s.authAuthHandler, s.invitationFriendHandler, s.friendHandler, s.userHandler, s.authMiddleware)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		fmt.Println("There is error: " + err.Error())
		return
	}
}
