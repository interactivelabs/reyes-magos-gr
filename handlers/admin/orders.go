package admin

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	ordersView "reyes-magos-gr/views/admin/orders"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersRepository     repository.OrdersRepository
	ToysRepository       repository.ToysRepository
	VolunteersRepository repository.VolunteersRepository
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

func (h OrdersHandler) OrdersViewHandler(ctx echo.Context) error {
	orders, err := h.OrdersRepository.GetAllActiveOrders()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	completedOrders, err := h.OrdersRepository.GetCompletedOrders()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.Orders(orders, completedOrders))
}

func (h OrdersHandler) OrderCardViewHandler(ctx echo.Context) error {
	orderID, err := getOrderId(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.LinkOrderCard(order))
}

func (h OrdersHandler) UpdateOrderViewHandler(ctx echo.Context) error {
	orderID, err := getOrderId(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	order, err := h.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	toy, err := h.ToysRepository.GetToyByID(order.ToyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	volunteer, err := h.VolunteersRepository.GetVolunteerByID(order.VolunteerID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.UpdateOrder(order, toy, volunteer))
}

type SaveOrderChangesrRequest struct {
	ShippedDate    string `form:"shipped_date" validate:"iso_8601_date"`
	OrderCompleted int64  `form:"order_completed" validate:"number"`
}

func (h OrdersHandler) SaveOrderChangesHandler(ctx echo.Context) error {
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

	order, err := h.OrdersRepository.GetOrderByID(orderID)
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

	err = h.OrdersRepository.UpdateOrder(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.LinkOrderCard(order))
}
