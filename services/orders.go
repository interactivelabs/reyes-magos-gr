package services

import (
	"database/sql"
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
	ToysStore           store.ToysStore
	VolunteerCodesStore store.VolunteerCodesStore
}

func NewOrdersService(
	cartsStore store.CartsStore,
	codesStore store.CodesStore,
	ordersStore store.OrdersStore,
	toysStore store.ToysStore,
	volunteerCodesStore store.VolunteerCodesStore,
) *OrdersServiceApp {
	return &OrdersServiceApp{
		CartsStore:          cartsStore,
		CodesStore:          codesStore,
		OrdersStore:         ordersStore,
		ToysStore:           toysStore,
		VolunteerCodesStore: volunteerCodesStore,
	}
}

type OrdersService interface {
	CreateOrder(toyID int64, code string) (order models.Order, err error)
	CreateOrders(volunteerID int64, items []CheckoutItem, codes []models.Code) error
	EnrichOrder(order models.Order) (models.OrderDetails, error)
	EnrichOrders(orders []models.Order) ([]models.OrderDetails, error)
}

func (s *OrdersServiceApp) EnrichOrder(order models.Order) (models.OrderDetails, error) {
	return enrichOrder(s.ToysStore, s.CodesStore, order)
}

func (s *OrdersServiceApp) EnrichOrders(orders []models.Order) ([]models.OrderDetails, error) {
	return enrichOrders(s.ToysStore, s.CodesStore, orders)
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

// CheckoutItem is one cart row being checked out: every row uses exactly one
// code, and its note is saved on that code so the volunteer can trace it back.
type CheckoutItem struct {
	CartID int64
	Note   string
}

func (s *OrdersServiceApp) CreateOrders(
	volunteerID int64,
	items []CheckoutItem,
	codes []models.Code,
) error {
	if len(items) == 0 {
		return ErrCartItemNotFound
	}
	if len(items) > len(codes) {
		return ErrNotEnoughCodes
	}

	cartItemsToUpdate := make([]models.CartItem, 0, len(items))
	ordersToCreate := make([]models.Order, 0, len(items))
	codesToUpdate := make([]models.Code, 0, len(items))
	seen := make(map[int64]bool, len(items))
	orderDate := time.Now().Format(time.RFC3339)

	for i, checkoutItem := range items {
		if seen[checkoutItem.CartID] {
			return ErrCartItemNotFound
		}
		seen[checkoutItem.CartID] = true

		item, err := s.CartsStore.GetCartItemByID(checkoutItem.CartID)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCartItemNotFound
		}
		if err != nil {
			return err
		}
		if item.VolunteerID != volunteerID {
			return ErrCartItemNotFound
		}

		code := codes[i]
		code.Given = 1
		code.Used = 1
		code.Note = checkoutItem.Note
		codesToUpdate = append(codesToUpdate, code)

		ordersToCreate = append(ordersToCreate, models.Order{
			ToyID:       item.ToyID,
			CodeID:      code.CodeID,
			VolunteerID: item.VolunteerID,
			OrderDate:   orderDate,
		})

		codeID := code.CodeID
		item.CodeID = &codeID
		item.Used = 1
		cartItemsToUpdate = append(cartItemsToUpdate, item)
	}

	if err := s.CodesStore.UpdateCodes(codesToUpdate); err != nil {
		return err
	}

	if err := s.CartsStore.UpdateCartItems(cartItemsToUpdate); err != nil {
		return err
	}

	return s.OrdersStore.CreateOrders(ordersToCreate)
}
