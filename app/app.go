package app

import (
	"database/sql"
	"reyes-magos-gr/services"
	"reyes-magos-gr/store"
)

type App struct {
	CartsStore               store.CartsStore
	CodesRepository          store.CodesRepository
	OrdersRepository         store.OrdersRepository
	ToysRepository           store.ToysRepository
	VolunteerCodesRepository store.VolunteerCodesRepository
	VolunteersRepository     store.VolunteersRepository
	CodesService             services.CodesService
	OrderService             services.OrdersService
	VolunteersService        services.VolunteersService
}

func NewApp(db *sql.DB) *App {
	cartsRepository := store.NewCartsStore(db)
	codesRepository := store.CodesRepository{DB: db}
	toysRepository := store.ToysRepository{DB: db}
	ordersRepository := store.OrdersRepository{DB: db}
	volunteerCodesRepository := store.VolunteerCodesRepository{DB: db}
	volunteersRepository := store.VolunteersRepository{DB: db}

	codesService := services.CodesService{
		CodesRepository: codesRepository,
	}

	orderService := services.OrdersService{
		CodesRepository:          codesRepository,
		OrdersRepository:         ordersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
	}

	volunteersService := services.NewVolunteersService(
		cartsRepository,
		codesRepository,
		ordersRepository,
		volunteersRepository,
		volunteerCodesRepository,
	)

	return &App{
		CodesRepository:          codesRepository,
		ToysRepository:           toysRepository,
		OrdersRepository:         ordersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
		VolunteersRepository:     volunteersRepository,
		CodesService:             codesService,
		OrderService:             orderService,
		VolunteersService:        volunteersService,
	}
}
