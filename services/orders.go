package services

import (
	"net/http"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	"time"

	"github.com/labstack/echo/v4"
)

type OrdersServiceApp struct {
	CodesStore          store.CodesStore
	OrdersStore         store.OrdersStore
	VolunteerCodesStore store.VolunteerCodesStore
}

func NewOrdersService(
	codesStore store.CodesStore,
	ordersStore store.OrdersStore,
	volunteerCodesStore store.VolunteerCodesStore,
) *OrdersServiceApp {
	return &OrdersServiceApp{
		CodesStore:          codesStore,
		OrdersStore:         ordersStore,
		VolunteerCodesStore: volunteerCodesStore,
	}
}

type OrdersService interface {
	CreateOrder(toyID int64, code string) (order models.Order, err error)
}

func (s *OrdersServiceApp) CreateOrder(toyID int64, code string) (order models.Order, err error) {
	codeResult, err := s.CodesStore.GetCode(code)
	if err != nil {
		return order, echo.NewHTTPError(http.StatusBadRequest, "Code Not Found")
	}

	if codeResult.Used == 1 {
		return order, echo.NewHTTPError(http.StatusConflict, "Code already used")
	}

	volunteerID, err := s.VolunteerCodesStore.GetVolunteerIdByCodeId(codeResult.CodeID)
	if err != nil {
		return order, echo.NewHTTPError(http.StatusBadRequest, "Code not assigned to volunteer")
	}

	order = models.Order{
		ToyID:       toyID,
		CodeID:      codeResult.CodeID,
		VolunteerID: volunteerID,
		OrderDate:   time.Now().Format(time.RFC3339),
	}

	orderID, err := s.OrdersStore.CreateOrder(order)
	if err != nil {
		return order, err
	}

	codeResult.Used = 1
	err = s.CodesStore.UpdateCode(codeResult)
	if err != nil {
		return order, err
	}

	order, err = s.OrdersStore.GetOrderByID(orderID)
	if err != nil {
		return order, err
	}

	return order, nil
}
