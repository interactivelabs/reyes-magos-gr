package main

import (
	"database/sql"
	"log"
	"os"
	"reyes-magos-gr/api"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/handlers"
	"reyes-magos-gr/middleware"
	"reyes-magos-gr/platform/authenticator"
	"reyes-magos-gr/services"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

const filename = "./db/reyes.sqlite"

func main() {
	e := echo.New()

	// Middleware
	e.Validator = middleware.NewValidator()

	cookieSecret := os.Getenv("REYES_COOKIE_SECRET")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cookieSecret))))

	auth, err := authenticator.New()
	if err != nil {
		log.Fatalf("Failed to initialize the authenticator: %v", err)
	}

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

	ordersRepository := repository.OrdersRepository{
		DB: db,
	}

	volunteerCodesRepository := repository.VolunteerCodesRepository{
		DB: db,
	}

	volunteersRepository := repository.VolunteersRepository{
		DB: db,
	}

	// CREATE SERVICES INSTANCES
	codesService := services.CodesService{
		CodesRepository: codesRepository,
	}

	orderService := services.OrdersService{
		CodesRepository:          codesRepository,
		OrdersRepository:         ordersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
	}

	volunteersService := services.VolunteersService{
		CodesRepository:          codesRepository,
		OrdersRepository:         ordersRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
	}

	// PUBLIC ENDPOINTS
	homeHandler := handlers.HomeHandler{}
	e.GET("/", homeHandler.HomeViewHandler)

	loginHandler := handlers.LoginHandler{
		Auth: auth,
	}
	e.GET("/login", loginHandler.LoginRedirectHandler)
	e.GET("/callback", loginHandler.LoginCallbackHandler)
	e.GET("/logout", loginHandler.LogoutRedirectHandler)

	catalogHandler := handlers.CatalogHandler{
		ToysRepository: toysRepository,
	}
	e.GET("/catalog", catalogHandler.CatalogViewHandler)

	redeemToyHandler := handlers.RedeemToyHandler{
		ToysRepository: toysRepository,
	}
	e.GET("/redeem/:toy_id", redeemToyHandler.RedeemToyViewHandler)

	redeemMultipleHandler := handlers.RedeemMultipleHandler{}
	e.GET("/redeem/multiple", redeemMultipleHandler.RedeemMultipleViewHandler)

	ordersHandler := handlers.OrdersHandler{
		OrdersService:        orderService,
		VolunteersRepository: volunteersRepository,
	}
	e.POST("/orders/create", ordersHandler.CreateOrderViewHandler)

	// Serve static files (css, js, images)
	e.Static("/public", "public")

	// VOLUNTEER ENDPOINTS
	vg := e.Group("/volunteer")

	vg.Use(middleware.IsAuthenticated())

	myCodesHandler := handlers.MyCodesHandler{
		VolunteersService: volunteersService,
	}
	vg.GET("/mycodes", myCodesHandler.MyCodesViewHandler)

	// ADMIN ENDPOINTS
	ag := e.Group("/admin")

	ag.Use(middleware.IsAdmin())

	toyHandler := api.ToyHandler{
		ToysRepository: toysRepository,
	}
	ag.POST("/api/toy", toyHandler.CreateToyApiHandler)
	ag.PUT("/api/toy", toyHandler.UpdateToyApiHandler)
	ag.DELETE("/api/toy/:toy_id", toyHandler.DeleteToyApiHandler)

	volunteerHandler := api.VolunteerHandler{
		VolunteersRepository: volunteersRepository,
	}
	ag.POST("/api/volunteer", volunteerHandler.CreateVolunteerApiHandler)
	ag.PUT("/api/volunteer", volunteerHandler.UpdateVolunteerApiHandler)
	ag.DELETE("/api/volunteer/:volunteer_id", volunteerHandler.DeleteVolunteerApiHandler)

	codeHandler := api.CodeHandler{
		CodesService: codesService,
	}
	ag.POST("/api/code", codeHandler.CreateCodeApiHandler)
	ag.POST("/api/code/batch", codeHandler.CreateCodeBatchApiHandler)

	// ADMIN VIEWS
	codesHTMLHandler := handlers.CodesHandler{
		CodesRepository:          codesRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
		CodesService:             codesService,
	}
	ag.GET("/codes", codesHTMLHandler.CodesViewHandler)
	ag.POST("/codes/assign", codesHTMLHandler.AssignCodesHandler)
	ag.POST("/codes/remove", codesHTMLHandler.RemoveCodesHandler)
	ag.POST("/codes/create", codesHTMLHandler.CreateCodesHandler)

	var port = "localhost:8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	e.Logger.Fatal(e.Start(port))
}
