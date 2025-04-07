package controller

import "github.com/swefinal-travel-planner/travel-app-be/internal/controller/http"

type ApiContainer struct {
	HttpServer *http.Server
}

func NewApiContainer(httpServer *http.Server) *ApiContainer {
	return &ApiContainer{HttpServer: httpServer}
}
