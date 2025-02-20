package main

import (
	"os"

	_ "github.com/swefinal-travel-planner/travel-app-be/docs"
	"github.com/swefinal-travel-planner/travel-app-be/startup"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate-up" {
		startup.Migrate()
		return
	}

	startup.Execute()
}
