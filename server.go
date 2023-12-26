package main

import (
	"reyes-magos-gr/handlers"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	homeHandler := handlers.HomeHandler{}
	e.GET("/", homeHandler.HomeViewHandler)

	redeemHandler := handlers.RedeemHandler{}
	e.GET("/redeem", redeemHandler.RedeemViewHandler)

	e.Static("/public", "public")

	e.Logger.Fatal(e.Start(":1323"))
}
