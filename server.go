package main

import (
	"database/sql"
	"log"
	"reyes-magos-gr/api"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/handlers"
	"reyes-magos-gr/middleware"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const filename = "./db/reyes.sqlite"

func main() {
	e := echo.New()

	// Middleware
	e.Validator = middleware.NewValidator()

	// Initialize DB
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Repositories
	toysRepository := repository.ToysRepository{
		DB: db,
	}

	// HTML VIEWS
	homeHandler := handlers.HomeHandler{}
	e.GET("/", homeHandler.HomeViewHandler)

	redeemHandler := handlers.RedeemHandler{
		DB: db,
	}
	e.GET("/redeem", redeemHandler.RedeemViewHandler)

	e.Static("/public", "public")

	// API ENDPOINTS
	toyHandler := api.ToyHandler{
		ToysRepository: toysRepository,
	}
	e.POST("/api/toys", toyHandler.CreateToyApiHandler)
	e.PUT("/api/toys", toyHandler.UpdateToyApiHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
