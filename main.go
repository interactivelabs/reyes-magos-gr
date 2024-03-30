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

	toysRepository := repository.ToysRepository{
		DB: db,
	}

	// ordersRepository := repository.OrdersRepository{
	// 	DB: db,
	// }

	volunteersCodesRepository := repository.VolunteerCodesRepository{
		DB: db,
	}

	volunteersRepository := repository.VolunteersRepository{
		DB: db,
	}

	// CREATE SERVICES INSTANCES
	codesService := services.CodesService{
		CodesRepository: codesRepository,
	}

	// CREATE HANDLERS INSTANCES
	homeHandler := handlers.HomeHandler{}

	toyHandler := api.ToyHandler{
		ToysRepository: toysRepository,
	}

	volunteerHandler := api.VolunteerHandler{
		VolunteersRepository: volunteersRepository,
	}

	codeHandler := api.CodeHandler{
		CodesService: codesService,
	}

	redeemHandler := handlers.RedeemHandler{
		ToysRepository: toysRepository,
	}

	redeemToyHandler := handlers.RedeemToyHandler{
		ToysRepository: toysRepository,
	}

	redeemMultipleHandler := handlers.RedeemMultipleHandler{}

	// HTML VIEWS
	codesHTMLHandler := handlers.CodesHandler{
		CodesRepository:          codesRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteersCodesRepository,
		CodesService:             codesService,
	}

	e.GET("/", homeHandler.HomeViewHandler)

	e.GET("/redeem", redeemHandler.RedeemViewHandler)

	e.GET("/redeem/:toy_id", redeemToyHandler.RedeemToyViewHandler)

	e.GET("/redeem_multiple", redeemMultipleHandler.RedeemMultipleViewHandler)

	// Login route
	loginHandler := api.LoginHandler{}
	e.POST("/login", loginHandler.UserLoginHandler)

	e.Static("/public", "public")

	// PUBLIC API ENDPOINTS

	// API ADMIN ENDPOINTS
	r := e.Group("/admin")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
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
	r.GET("/codes", codesHTMLHandler.CodesViewHandler)
	r.POST("/codes/assign", codesHTMLHandler.AssignCodesHandler)
	r.POST("/codes/remove", codesHTMLHandler.RemoveCodesHandler)
	r.POST("/codes/create", codesHTMLHandler.CreateCodesHandler)

	var port = "localhost:8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	e.Logger.Fatal(e.Start(port))
}
