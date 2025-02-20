package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/v1"
)

type Server struct {
	studentHandler  *v1.StudentHandler
	authAuthHandler *v1.AuthHandler
}

func NewServer(studentHandler *v1.StudentHandler, authAuthHandler *v1.AuthHandler) *Server {
	return &Server{studentHandler: studentHandler, authAuthHandler: authAuthHandler}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	v1.MapRoutes(router, s.studentHandler, s.authAuthHandler)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		return
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)
}
