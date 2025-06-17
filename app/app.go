package app

import (
	"database/sql"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
)

type App struct {
	CartsStore          store.CartsStore
	CodesStore          store.CodesStore
	OrdersStore         store.OrdersStore
	ToysStore           store.ToysStore
	VolunteerCodesStore store.VolunteerCodesStore
	VolunteersStore     store.VolunteersStore
	CodesService        services.CodesService
	OrderService        services.OrdersService
	VolunteersService   services.VolunteersService
}

func NewApp(db *sql.DB) *App {
	cartsRepository := store.NewCartsStore(db)
	codesStore := store.NewCodesStore(db)
	ordersStore := store.NewOrdersStore(db)
	toysRepository := store.NewToysStore(db)
	volunteerCodesStore := store.NewVolunteerCodesStore(db)
	volunteersStore := store.NewVolunteersStore(db)

	codesService := services.NewCodesService(codesStore)

	orderService := services.NewOrdersService(
		codesStore,
		ordersStore,
		volunteerCodesStore,
	)

	volunteersService := services.NewVolunteersService(
		cartsRepository,
		codesStore,
		ordersStore,
		volunteersStore,
		volunteerCodesStore,
	)

	return &App{
		CodesStore:          codesStore,
		OrdersStore:         ordersStore,
		ToysStore:           toysRepository,
		VolunteerCodesStore: volunteerCodesStore,
		VolunteersStore:     volunteersStore,
		CodesService:        codesService,
		OrderService:        orderService,
		VolunteersService:   volunteersService,
	}
}
