//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http"
	v1 "github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/v1"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	repositoryimplement "github.com/swefinal-travel-planner/travel-app-be/internal/repository/implement"
	serviceimplement "github.com/swefinal-travel-planner/travel-app-be/internal/service/implement"
)

var container = wire.NewSet(
	controller.NewApiContainer,
)

// may have grpc server in the future
var serverSet = wire.NewSet(
	http.NewServer,
)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(
	v1.NewStudentHandler,
	v1.NewAuthHandler,
)

var serviceSet = wire.NewSet(
	serviceimplement.NewStudentService,
	serviceimplement.NewAuthService,
)

var repositorySet = wire.NewSet(
	repositoryimplement.NewStudentRepository,
	repositoryimplement.NewUserRepository,
)

func InitializeContainer(
	db database.Db,
) *controller.ApiContainer {
	wire.Build(serverSet, handlerSet, serviceSet, repositorySet, container)
	return &controller.ApiContainer{}
}
