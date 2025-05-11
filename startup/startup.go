package startup

import (
	"os"

	"github.com/gammazero/workerpool"
	log "github.com/sirupsen/logrus"
	"github.com/swefinal-travel-planner/travel-app-be/internal"
	"github.com/swefinal-travel-planner/travel-app-be/internal/controller"
	"github.com/swefinal-travel-planner/travel-app-be/internal/database"
)

func Migrate() {
	// Open the database connection
	db := database.Open()

	database.MigrateUp(db)
}

func registerDependencies() *controller.ApiContainer {
	// Open database connection
	db := database.Open()

	return internal.InitializeContainer(db)
}

func Execute() {
	// Configure logrus
	log.SetOutput(os.Stdout)

	container := registerDependencies()

	wp := workerpool.New(2)

	wp.Submit(container.HttpServer.Run)

	wp.StopWait()
}
