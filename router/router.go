package router

import (
	"net/http"
	"os"
	"reyes-magos-gr/app"
	"reyes-magos-gr/handlers"
	"reyes-magos-gr/handlers/admin"
	"reyes-magos-gr/handlers/volunteers"
	reyes_middleware "reyes-magos-gr/middleware"
	"reyes-magos-gr/platform/authenticator"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

func SetupRouter(app *app.App, auth *authenticator.Authenticator) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Validator = reyes_middleware.NewValidator()

	// Security and session configuration
	cookieSecret := os.Getenv("REYES_COOKIE_SECRET")
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(cookieSecret))))

	env := os.Getenv("ENV")
	csrfDomain := "dl-toys.com"
	if env == "development" {
		csrfDomain = "localhost"
	}

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "cookie:_csrf",
		CookiePath:     "/",
		CookieDomain:   csrfDomain,
		CookieSecure:   true,
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteStrictMode,
	}))

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(10))))
	e.Use(middleware.Gzip())

	allowOrigins := []string{
		"https://dl-toys.com",
		"https://www.dl-toys.com",
		"https://static.dl-toys.com",
	}
	if env == "development" {
		allowOrigins = append(allowOrigins, "http://localhost:8080")
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowOrigins,
	}))

	// PUBLIC ENDPOINTS
	homeHandler := handlers.NewHomeHandler()
	e.GET("/", homeHandler.HomeViewHandler)
	e.GET("/support", homeHandler.SupportViewHandler)
	e.GET("/401", homeHandler.Error401)
	e.GET("/404", homeHandler.Error404)
	e.GET("/500", homeHandler.Error500)
	e.GET("/health", homeHandler.HealthViewHandler)
	e.GET("/verifyemail", homeHandler.VerifyEmailHandler)
	e.GET("/notvolunteer", homeHandler.NotVolunteerHandler)

	loginHandler := handlers.NewLoginHandler(auth)
	e.GET("/login/redirect", loginHandler.LoginRedirectHandler)
	e.GET("/login", loginHandler.LoginRedirectHandler)
	e.GET("/callback", loginHandler.LoginCallbackHandler)
	e.GET("/logout", loginHandler.LogoutRedirectHandler)

	catalogHandler := handlers.NewCatalogHandler(app.ToysStore)
	e.GET("/catalog", catalogHandler.CatalogViewHandler)

	redeemToyHandler := handlers.NewRedeemToyHandler(app.ToysStore)
	e.GET("/redeem/:toy_id", redeemToyHandler.RedeemToyViewHandler)

	ordersHandler := handlers.NewOrdersHandler(app.VolunteersStore, app.OrderService)
	e.POST("/orders/create", ordersHandler.CreateOrderViewHandler)

	// Serve static files (css, js, images)
	e.Static("/public", "public")

	// VOLUNTEER ENDPOINTS
	vg := e.Group("/volunteer")

	vg.Use(reyes_middleware.IsAuthenticated())

	myCodesHandler := volunteers.NewMyCodesHandler(app.CodesStore, app.VolunteersService)
	vg.GET("/mycodes", myCodesHandler.MyCodesViewHandler)
	vg.POST("/mycodes/give/:code_id", myCodesHandler.GiveCode)

	myOrdersHandler := volunteers.NewMyOrdersHandler(app.OrdersStore, app.VolunteersService)
	vg.GET("/myorders", myOrdersHandler.MyOrdersViewHandler)
	vg.POST("/myorders/:order_id/completed", myOrdersHandler.MyOrdersCompleteHandler)

	myCartHandler := volunteers.NewCartHandler(app.VolunteersService)
	vg.GET("/mycart", myCartHandler.CartViewHandler)

	// ADMIN ENDPOINTS
	ag := e.Group("/admin")

	ag.Use(reyes_middleware.IsAdmin())

	codesHandler := admin.NewCodesHandler(
		app.CodesStore,
		app.VolunteersStore,
		app.VolunteerCodesStore,
		app.CodesService,
	)
	ag.GET("/codes", codesHandler.CodesViewHandler)
	ag.POST("/codes/assign", codesHandler.AssignCodesHandler)
	ag.POST("/codes/remove", codesHandler.RemoveCodesHandler)
	ag.POST("/codes/create", codesHandler.CreateCodesHandler)

	adminOrdersHandler := admin.NewOrdersHandler(
		app.OrdersStore,
		app.ToysStore,
		app.VolunteersStore,
	)
	ag.GET("/orders", adminOrdersHandler.OrdersViewHandler)
	ag.GET("/order/:order_id", adminOrdersHandler.OrderCardViewHandler)
	ag.GET("/order/:order_id/edit", adminOrdersHandler.UpdateOrderViewHandler)
	ag.PUT("/order/:order_id/save", adminOrdersHandler.SaveOrderChangesHandler)

	volunteersHandler := admin.NewVolunteersHandler(app.VolunteersStore, app.VolunteersService)
	ag.GET("/volunteers", volunteersHandler.VolunteersViewHandler)
	ag.GET("/volunteers/create", volunteersHandler.VolunteersCreateHandler)
	ag.POST("/volunteers", volunteersHandler.VolunteersCreatePostHandler)
	ag.GET("/volunteers/:volunteer_id", volunteersHandler.VolunteersUpdateViewHandler)
	ag.PUT("/volunteers/:volunteer_id/save", volunteersHandler.VolunteersUpdatePutHandler)
	ag.DELETE("/volunteers/:volunteer_id/delete", volunteersHandler.VolunteersDeleteHandler)

	toysHandler := admin.NewToysHandler(app.ToysStore)
	ag.GET("/toys", toysHandler.ToysViewHandler)
	ag.GET("/toys/create", toysHandler.CreateToyFormHandler)
	ag.POST("/toys", toysHandler.CreateToyPostHandler)
	ag.GET("/toys/:toy_id", toysHandler.UpdateToyFormHandler)
	ag.PUT("/toys/:toy_id/save", toysHandler.UpdateToyPutHandler)
	ag.DELETE("/toys/:toy_id/delete", toysHandler.DeleteToyHandler)
	ag.GET("/toys/categories", toysHandler.ToysCategoriesViewHandler)

	return e
}
