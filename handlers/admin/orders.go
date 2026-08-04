package admin

import (
	"database/sql"
	"errors"
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
	ordersView "reyes-magos-gr/views/admin/orders"
	"reyes-magos-gr/views/components"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersStore     store.OrdersStore
	ToysStore       store.ToysStore
	VolunteersStore store.VolunteersStore
	CodesStore      store.CodesStore
}

func NewOrdersHandler(
	ordersStore store.OrdersStore,
	toysStore store.ToysStore,
	volunteersStore store.VolunteersStore,
	codesStore store.CodesStore,
) *OrdersHandler {
	return &OrdersHandler{
		OrdersStore:     ordersStore,
		ToysStore:       toysStore,
		VolunteersStore: volunteersStore,
		CodesStore:      codesStore,
	}
}

func (h *OrdersHandler) toOrderDetails(order models.Order) (models.OrderDetails, error) {
	toy, err := h.ToysStore.GetToyByID(order.ToyID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.OrderDetails{}, err
	}

	code, err := h.CodesStore.GetCodeByID(order.CodeID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.OrderDetails{}, err
	}

	return models.OrderDetails{
		Order:        order,
		ToyName:      toy.ToyName,
		ToySourceURL: toy.SourceURL,
		Code:         code.Code,
	}, nil
}

func (h *OrdersHandler) toOrdersDetails(orders []models.Order) ([]models.OrderDetails, error) {
	details := make([]models.OrderDetails, 0, len(orders))
	for _, order := range orders {
		detail, err := h.toOrderDetails(order)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}

type Order interface {
	Param(name string) string
}

func getOrderId(o Order) (int64, error) {
	orderIDStr := o.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}
	return orderID, nil
}

func (h *OrdersHandler) OrdersViewHandler(ctx echo.Context) error {
	orders, err := h.OrdersStore.GetAllActiveOrders()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	completedOrders, err := h.OrdersStore.GetCompletedOrders()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	orderDetails, err := h.toOrdersDetails(orders)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	completedOrderDetails, err := h.toOrdersDetails(completedOrders)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.Orders(orderDetails, completedOrderDetails))
}

func (h *OrdersHandler) OrderCardViewHandler(ctx echo.Context) error {
	orderID, err := getOrderId(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersStore.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	details, err := h.toOrderDetails(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, components.OrderCard(details, true))
}

func (h *OrdersHandler) UpdateOrderViewHandler(ctx echo.Context) error {
	orderID, err := getOrderId(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersStore.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	toy, err := h.ToysStore.GetToyByID(order.ToyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	volunteer, err := h.VolunteersStore.GetVolunteerByID(order.VolunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.UpdateOrder(order, toy, volunteer))
}

type SaveOrderChangesrRequest struct {
	ShippedDate    string `form:"shipped_date"    validate:"iso_8601_date"`
	OrderCompleted int64  `form:"order_completed" validate:"number"`
}

func (h *OrdersHandler) SaveOrderChangesHandler(ctx echo.Context) error {
	saveOrderRequest := new(SaveOrderChangesrRequest)
	if err := ctx.Bind(saveOrderRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Validate(saveOrderRequest); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	orderID, err := getOrderId(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersStore.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if saveOrderRequest.ShippedDate != "" {
		order.Shipped = 1
	}

	shippedDate, err := time.Parse(lib.YYYYMMDD, saveOrderRequest.ShippedDate)
	order.ShippedDate = shippedDate.Format(time.RFC3339)

	order.Completed = saveOrderRequest.OrderCompleted

	if order.Completed == 1 {
		order.CompletedDate = time.Now().Format(time.RFC3339)
	}

	err = h.OrdersStore.UpdateOrder(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	details, err := h.toOrderDetails(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, components.OrderCard(details, true))
}
