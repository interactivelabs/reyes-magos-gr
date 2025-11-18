package services

import (
	"errors"
	"net/http"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	"time"

	"github.com/labstack/echo/v4"
)

type OrdersServiceApp struct {
	CartsStore          store.CartsStore
	CodesStore          store.CodesStore
	OrdersStore         store.OrdersStore
	VolunteerCodesStore store.VolunteerCodesStore
}

func NewOrdersService(
	cartsStore store.CartsStore,
	codesStore store.CodesStore,
	ordersStore store.OrdersStore,
	volunteerCodesStore store.VolunteerCodesStore,
) *OrdersServiceApp {
	return &OrdersServiceApp{
		CartsStore:          cartsStore,
		CodesStore:          codesStore,
		OrdersStore:         ordersStore,
		VolunteerCodesStore: volunteerCodesStore,
	}
}

type OrdersService interface {
	CreateOrder(toyID int64, code string) (order models.Order, err error)
	CreateOrders(cartIDs []int64, quantity []int, codes []models.Code) error
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

func (s *OrdersServiceApp) CreateOrders(
	cartIDs []int64,
	quantities []int,
	codes []models.Code,
) error {
	// Calculate total codes needed
	totalCodesNeeded := 0
	for _, qty := range quantities {
		totalCodesNeeded += qty
	}

	if totalCodesNeeded > len(codes) {
		return errors.New("not enough codes for the order")
	}

	// Pre-allocate with known capacity
	cartItemsToUpdate := make([]models.CartItem, 0, len(cartIDs))
	ordersToCreate := make([]models.Order, 0, totalCodesNeeded)
	codesToUpdate := make([]models.Code, 0, totalCodesNeeded)
	codeIndex := 0

	for i, cartID := range cartIDs {
		item, err := s.CartsStore.GetCartItemByID(cartID)
		if err != nil {
			return err
		}

		// Create orders for this cart item
		for range quantities[i] {
			code := codes[codeIndex]
			ordersToCreate = append(ordersToCreate, models.Order{
				ToyID:       item.ToyID,
				CodeID:      code.CodeID,
				VolunteerID: item.VolunteerID,
				OrderDate:   time.Now().Format(time.RFC3339),
			})

			code.Given = 1
			code.Used = 1

			codesToUpdate = append(codesToUpdate, code)
			codeIndex++
		}

		// Update cart item with the LAST code used (or rethink this logic)
		codeID := codes[codeIndex-1].CodeID
		item.CodeID = &codeID
		item.Used = 1

		cartItemsToUpdate = append(cartItemsToUpdate, item)
	}

	// Update all codes in a batch
	if err := s.CodesStore.UpdateCodes(codesToUpdate); err != nil {
		return err
	}

	// Update all cart items atomically
	if err := s.CartsStore.UpdateCartItems(cartItemsToUpdate); err != nil {
		return err
	}

	// Create all orders first
	return s.OrdersStore.CreateOrders(ordersToCreate)
}
