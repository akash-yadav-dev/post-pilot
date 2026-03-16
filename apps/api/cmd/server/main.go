package main

import (
	"log"
	"post-pilot/apps/api/cmd/server/bootstrap"
)

func main() {

	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("Error initializing app: %v", err)
	}
	app.Start()
}
