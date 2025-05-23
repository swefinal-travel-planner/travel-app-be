package main

import (
	"os"

	_ "github.com/swefinal-travel-planner/travel-app-be/docs"

	"github.com/swefinal-travel-planner/travel-app-be/startup"
)

// @title Travel App API
// @version 1.0
// @description This is a travel planning application API
// @BasePath /api/v1
func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate-up" {
		startup.Migrate()
		return
	}

	startup.Execute()
}
