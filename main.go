package main

import (
	"os"

	"github.com/swefinal-travel-planner/travel-app-be/startup"
)

// @BasePath /api/v1
func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate-up" {
		startup.Migrate()
		return
	}

	startup.Execute()
}
