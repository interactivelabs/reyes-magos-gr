package services

import (
	"net/http"
	"reyes-magos-gr/db/model"
	"reyes-magos-gr/db/repository"
	"time"

	"github.com/labstack/echo/v4"
)

type OrdersService struct {
	CodesRepository          repository.CodesRepository
	OrdersRepository         repository.OrdersRepository
	VolunteerCodesRepository repository.VolunteerCodesRepository
}

func (s OrdersService) CreateOrder(toyID int64, code string) (order model.Order, err error) {
	codeResult, err := s.CodesRepository.GetCode(code)
	if err != nil {
		return order, echo.NewHTTPError(http.StatusBadRequest, "Code Not Found")
	}

	if codeResult.Used == 1 {
		return order, echo.NewHTTPError(http.StatusConflict, "Code already used")
	}

	volunteerID, err := s.VolunteerCodesRepository.GetVolunteerIdByCodeId(codeResult.CodeID)
	if err != nil {
		return order, echo.NewHTTPError(http.StatusBadRequest, "Code not assigned to volunteer")
	}

	order = model.Order{
		ToyID:       toyID,
		CodeID:      codeResult.CodeID,
		VolunteerID: volunteerID,
		OrderDate:   time.Now().Format(time.RFC3339),
	}

	orderID, err := s.OrdersRepository.CreateOrder(order)
	if err != nil {
		return order, err
	}

	codeResult.Used = 1
	err = s.CodesRepository.UpdateCode(codeResult)
	if err != nil {
		return order, err
	}

	order, err = s.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return order, err
	}

	return order, nil
}
