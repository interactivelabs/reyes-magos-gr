package main

import (
	"database/sql"
	"log"
	"os"
	"reyes-magos-gr/api"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/handlers"
	"reyes-magos-gr/middleware"
	"reyes-magos-gr/services"

	echojwt "github.com/labstack/echo-jwt/v4"

	"github.com/golang-jwt/jwt/v5"
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

	// CREATE REPOSITORY INSTANCES
	codesRepository := repository.CodesRepository{
		DB: db,
	}

	volunteersRepository := repository.VolunteersRepository{
		DB: db,
	}

	volunteersCodesRepository := repository.VolunteerCodesRepository{
		DB: db,
	}

	// CREATE SERVICES INSTANCES
	codesService := services.CodesService{
		CodesRepository: codesRepository,
	}

	// CREATE HANDLERS INSTANCES
	toyHandler := api.ToyHandler{
		ToysRepository: repository.ToysRepository{
			DB: db,
		},
	}
	volunteerHandler := api.VolunteerHandler{
		VolunteersRepository: repository.VolunteersRepository{
			DB: db,
		},
	}
	codeHandler := api.CodeHandler{
		CodesService: codesService,
	}

	// HTML VIEWS
	codesHTMLHander := handlers.CodesHandler{
		CodesRepository:          codesRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteersCodesRepository,
		CodesService:             codesService,
	}

	homeHandler := handlers.HomeHandler{}
	e.GET("/", homeHandler.HomeViewHandler)

	redeemHandler := handlers.RedeemHandler{}
	e.GET("/redeem", redeemHandler.RedeemViewHandler)

	redeemMultipleHandler := handlers.RedeemMultipleHandler{}
	e.GET("/redeem_multiple", redeemMultipleHandler.RedeemMultipleViewHandler)

	loginHandler := api.LoginHandler{}
	// Login route
	e.POST("/login", loginHandler.UserLoginHandler)

	e.Static("/public", "public")

	// PUBLIC API ENDPOINTS

	// API ADMIN ENDPOINTS
	r := e.Group("/admin")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(api.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))

	r.POST("/api/toy", toyHandler.CreateToyApiHandler)
	r.PUT("/api/toy", toyHandler.UpdateToyApiHandler)
	r.DELETE("/api/toy/:toy_id", toyHandler.DeleteToyApiHandler)

	r.POST("/api/volunteer", volunteerHandler.CreateVolunteerApiHandler)
	r.PUT("/api/volunteer", volunteerHandler.UpdateVolunteerApiHandler)
	r.DELETE("/api/volunteer/:volunteer_id", volunteerHandler.DeleteVolunteerApiHandler)

	r.POST("/api/code", codeHandler.CreateCodeApiHandler)
	r.POST("/api/code/batch", codeHandler.CreateCodeBatchApiHandler)

	// ADMIN VIEWS
	r.GET("/codes", codesHTMLHander.CodesViewHandler)
	r.POST("/codes/assign", codesHTMLHander.AssignCodesHandler)
	r.POST("/codes/remove", codesHTMLHander.RemoveCodesHandler)
	r.POST("/codes/create", codesHTMLHander.CreateCodesHandler)

	var port = "localhost:8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	e.Logger.Fatal(e.Start(port))
}
