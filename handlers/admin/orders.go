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

	return lib.Render(ctx, ordersView.Orders(orders))
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
	ToyID int64  `form:"toy_id" validate:"required"`
	Code  string `form:"code" validate:"required"`
}

func (h OrdersHandler) SaveOrderChangesHandler(ctx echo.Context) error {
	orderIDStr := ctx.Param("order_id")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}

	order, err := h.OrdersRepository.GetOrderByID(orderID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = h.OrdersRepository.UpdateOrder(order)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.Redirect(http.StatusSeeOther, "/admin/orders")
}
