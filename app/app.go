package app

import (
	"database/sql"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/services"
)

type App struct {
	CartsRepository          repository.CartsRepository
	CodesRepository          repository.CodesRepository
	OrdersRepository         repository.OrdersRepository
	ToysRepository           repository.ToysRepository
	VolunteerCodesRepository repository.VolunteerCodesRepository
	VolunteersRepository     repository.VolunteersRepository
	CodesService             services.CodesService
	OrderService             services.OrdersService
	VolunteersService        services.VolunteersService
}

func NewApp(db *sql.DB) *App {
	cartsRepository := repository.CartsRepository{DB: db}
	codesRepository := repository.CodesRepository{DB: db}
	toysRepository := repository.ToysRepository{DB: db}
	ordersRepository := repository.OrdersRepository{DB: db}
	volunteerCodesRepository := repository.VolunteerCodesRepository{DB: db}
	volunteersRepository := repository.VolunteersRepository{DB: db}

	codesService := services.CodesService{
		CodesRepository: codesRepository,
	}

	orderService := services.OrdersService{
		CodesRepository:          codesRepository,
		OrdersRepository:         ordersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
	}

	volunteersService := services.VolunteersService{
		CartsRepository:          cartsRepository,
		CodesRepository:          codesRepository,
		OrdersRepository:         ordersRepository,
		VolunteersRepository:     volunteersRepository,
		VolunteerCodesRepository: volunteerCodesRepository,
	}

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
