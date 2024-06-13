package handlers_admin

import (
	"net/http"
	"reyes-magos-gr/db/repository"
	"reyes-magos-gr/lib"
	ordersView "reyes-magos-gr/views/admin/orders"
	"strconv"

	"github.com/labstack/echo/v4"
)

type OrdersHandler struct {
	OrdersRepository     repository.OrdersRepository
	ToysRepository       repository.ToysRepository
	VolunteersRepository repository.VolunteersRepository
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
	orderIDStr := ctx.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.OrderCard(order))
}

func (h OrdersHandler) UpdateOrderViewHandler(ctx echo.Context) error {
	orderIDStr := ctx.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
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
		return err
	}

	orderIDStr := ctx.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if saveOrderRequest.ShippedDate != "" {
		order.Shipped = 1
	}
	order.ShippedDate = saveOrderRequest.ShippedDate
	order.Completed = saveOrderRequest.OrderCompleted

	err = h.OrdersRepository.UpdateOrder(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return lib.Render(ctx, ordersView.OrderCard(order))
}
