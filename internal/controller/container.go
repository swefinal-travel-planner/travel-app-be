package controller

import (
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/cronjob"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http"
)

type ApiContainer struct {
	HttpServer      *http.Server
	CronJobRegister *cronjob.CronJobRegister
}

func NewApiContainer(httpServer *http.Server, cronJobRegister *cronjob.CronJobRegister) *ApiContainer {
	return &ApiContainer{HttpServer: httpServer, CronJobRegister: cronJobRegister}
}
