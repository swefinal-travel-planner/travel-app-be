//go:build wireinject
// +build wireinject

package internal

import (
	"github.com/google/wire"
	beanimplement "github.com/swefinal-travel-planner/travel-app-be/internal/bean/implement"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
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
	v1.NewAuthHandler,
	v1.NewInvitationFriendHandler,
	v1.NewFriendHandler,
)

var serviceSet = wire.NewSet(
	serviceimplement.NewAuthService,
	serviceimplement.NewInvitationFriendService,
	serviceimplement.NewFriendService,
)

var repositorySet = wire.NewSet(
	repositoryimplement.NewUserRepository,
	repositoryimplement.NewAuthenticationRepository,
	repositoryimplement.NewInvitationFriendRepository,
	repositoryimplement.NewFriendRepository,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)

var beanSet = wire.NewSet(
	beanimplement.NewBcryptPasswordEncoder,
	beanimplement.NewRedisService,
	beanimplement.NewMailClient,
)

func InitializeContainer(
	db database.Db,
) *controller.ApiContainer {
	wire.Build(serverSet, handlerSet, serviceSet, repositorySet, middlewareSet, beanSet, container)
	return &controller.ApiContainer{}
}
