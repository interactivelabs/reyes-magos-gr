package main

import (
	"fmt"
	"log"
	"os"
	"reyes-magos-gr/app"
	"reyes-magos-gr/middleware"
	"reyes-magos-gr/platform/authenticator"
	"reyes-magos-gr/platform/database"
	"reyes-magos-gr/platform/flags"
	"reyes-magos-gr/router"
)

func main() {
	// Initialize DB
	db, connector, dir, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := connector.Sync(); err != nil {
		fmt.Println("Error syncing database:", err)
	}
	defer os.RemoveAll(dir)
	defer connector.Close()
	defer db.Close()

	// Initialize the application store and services
	app := app.NewApp(db)

	// Initialize the Auth0 authenticator
	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	// Initialize the PostHog client
	flagsClient, err := flags.NewFlagsClient()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

	e := router.SetupRouter(app, auth, flagsClient)

	e.HTTPErrorHandler = middleware.CustomHTTPErrorHandler

	var host = "localhost"
	var port = "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if os.Getenv("RAILWAY_PROJECT_ID") != "" {
		host = "0.0.0.0"
	}

	e.Logger.Fatal(e.Start(host + ":" + port))
}
