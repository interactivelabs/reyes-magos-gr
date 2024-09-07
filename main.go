package main

import (
	"fmt"
	"log"
	"os"
	"reyes-magos-gr/api"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/handlers"
	"reyes-magos-gr/handlers/admin"
	"reyes-magos-gr/handlers/volunteers"
	"reyes-magos-gr/middleware"
	"reyes-magos-gr/platform/authenticator"
	"reyes-magos-gr/platform/database"
	"reyes-magos-gr/services"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

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

	// CREATE REPOSITORY INSTANCES
	codesRepository := repository.CodesRepository{DB: db}
	toysRepository := repository.ToysRepository{DB: db}
	ordersRepository := repository.OrdersRepository{DB: db}
	volunteerCodesRepository := repository.VolunteerCodesRepository{DB: db}
	volunteersRepository := repository.VolunteersRepository{DB: db}

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
	e.GET("/support", homeHandler.SupportViewHandler)
	e.GET("/401", homeHandler.Error401)
	e.GET("/404", homeHandler.Error404)
	e.GET("/500", homeHandler.Error500)
	e.GET("/health", homeHandler.HealthViewHandler)

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

	myCodesHandler := volunteers.MyCodesHandler{
		VolunteersService: volunteersService,
		CodesRepository:   codesRepository,
	}
	vg.GET("/mycodes", myCodesHandler.MyCodesViewHandler)
	vg.POST("/mycodes/give/:code_id", myCodesHandler.GiveCode)

	myOrdersHandler := volunteers.MyOrdersHandler{
		VolunteersService: volunteersService,
		Ordersrepository:  ordersRepository,
	}
	vg.GET("/myorders", myOrdersHandler.MyOrdersViewHandler)
	vg.POST("/myorders/:order_id/completed", myOrdersHandler.MyOrdersCompleteHandler)

	// ADMIN ENDPOINTS
	ag := e.Group("/admin")

	ag.Use(middleware.IsAdmin())

	toyHandler := api.ToyHandler{
		ToysRepository: toysRepository,
	}
	ag.POST("/api/toy", toyHandler.CreateToyApiHandler)
	ag.POST("/api/toys", toyHandler.CreateBatchToysApiHandler)
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
	codesHandler := admin.CodesHandler{
		CodesRepository:          codesRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
		CodesService:             codesService,
	}
	ag.GET("/codes", codesHandler.CodesViewHandler)
	ag.POST("/codes/assign", codesHandler.AssignCodesHandler)
	ag.POST("/codes/remove", codesHandler.RemoveCodesHandler)
	ag.POST("/codes/create", codesHandler.CreateCodesHandler)

	adminOrdersHandler := admin.OrdersHandler{
		OrdersRepository:     ordersRepository,
		ToysRepository:       toysRepository,
		VolunteersRepository: volunteersRepository,
	}
	ag.GET("/orders", adminOrdersHandler.OrdersViewHandler)
	ag.GET("/order/:order_id", adminOrdersHandler.OrderCardViewHandler)
	ag.GET("/order/:order_id/edit", adminOrdersHandler.UpdateOrderViewHandler)
	ag.PUT("/order/:order_id/save", adminOrdersHandler.SaveOrderChangesHandler)

	volunteersHandler := admin.VolunteersHandler{
		VolunteersService:    volunteersService,
		VolunteersRepository: volunteersRepository,
	}
	ag.GET("/volunteers", volunteersHandler.VolunteersViewHandler)
	ag.GET("/volunteers/create", volunteersHandler.VolunteersCreateHandler)
	ag.POST("/volunteers", volunteersHandler.VolunteersCreatePostHandler)
	ag.GET("/volunteers/:volunteer_id", volunteersHandler.VolunteersUpdateViewHandler)
	ag.PUT("/volunteers/:volunteer_id/save", volunteersHandler.VolunteersUpdatePutHandler)
	ag.DELETE("/volunteers/:volunteer_id/delete", volunteersHandler.VolunteersDeleteHandler)

	toysHandler := admin.ToysHandler{
		ToysRepository: toysRepository,
	}
	ag.GET("/toys", toysHandler.ToysViewHandler)
	ag.GET("/toys/create", toysHandler.CreateToyFormHandler)
	ag.POST("/toys", toysHandler.CreateToyPostHandler)
	ag.GET("/toys/:toy_id", toysHandler.UpdateToyFormHandler)
	ag.PUT("/toys/:toy_id/save", toysHandler.UpdateToyPutHandler)
	ag.DELETE("/toys/:toy_id/delete", toysHandler.DeleteToyHandler)

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
