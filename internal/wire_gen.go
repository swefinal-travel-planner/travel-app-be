// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/google/wire"
	"github.com/swefinal-travel-planner/travel-app-be/internal/bean/implement"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/middleware"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller/http/v1"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
	"github.com/swefinal-travel-planner/travel-app-be/internal/repository/implement"
	"github.com/swefinal-travel-planner/travel-app-be/internal/service/implement"
)

// Injectors from wire.go:

func InitializeContainer(db database.Db) *controller.ApiContainer {
	userRepository := repositoryimplement.NewUserRepository(db)
	authenticationRepository := repositoryimplement.NewAuthenticationRepository(db)
	passwordEncoder := beanimplement.NewBcryptPasswordEncoder()
	redisClient := beanimplement.NewRedisService()
	mailClient := beanimplement.NewMailClient()
	authService := serviceimplement.NewAuthService(userRepository, authenticationRepository, passwordEncoder, redisClient, mailClient)
	authHandler := v1.NewAuthHandler(authService)
	invitationFriendRepository := repositoryimplement.NewInvitationFriendRepository(db)
	friendRepository := repositoryimplement.NewFriendRepository(db)
	invitationFriendService := serviceimplement.NewInvitationFriendService(invitationFriendRepository, userRepository, friendRepository)
	invitationFriendHandler := v1.NewInvitationFriendHandler(invitationFriendService)
	friendService := serviceimplement.NewFriendService(friendRepository, userRepository)
	friendHandler := v1.NewFriendHandler(friendService)
	authMiddleware := middleware.NewAuthMiddleware(authService, authenticationRepository, userRepository)
	server := http.NewServer(authHandler, invitationFriendHandler, friendHandler, authMiddleware)
	apiContainer := controller.NewApiContainer(server)
	return apiContainer
}

// wire.go:

var container = wire.NewSet(controller.NewApiContainer)

// may have grpc server in the future
var serverSet = wire.NewSet(http.NewServer)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(v1.NewAuthHandler, v1.NewInvitationFriendHandler, v1.NewFriendHandler)

var serviceSet = wire.NewSet(serviceimplement.NewAuthService, serviceimplement.NewInvitationFriendService, serviceimplement.NewFriendService)

var repositorySet = wire.NewSet(repositoryimplement.NewUserRepository, repositoryimplement.NewAuthenticationRepository, repositoryimplement.NewInvitationFriendRepository, repositoryimplement.NewFriendRepository)

var middlewareSet = wire.NewSet(middleware.NewAuthMiddleware)

var beanSet = wire.NewSet(beanimplement.NewBcryptPasswordEncoder, beanimplement.NewRedisService, beanimplement.NewMailClient)
