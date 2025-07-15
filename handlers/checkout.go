package handlers

import (
	"net/http"
	"reyes-magos-gr/lib"
	"reyes-magos-gr/store"
	checkout "reyes-magos-gr/views/checkout"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RedeemToyHandler struct {
	ToysStore store.ToysStore
}

func NewRedeemToyHandler(toysStore store.ToysStore) *RedeemToyHandler {
	return &RedeemToyHandler{
		ToysStore: toysStore,
	}
}

type RedeemToyViewRequest struct {
	Code string `query:"code"`
}

func (h *RedeemToyHandler) RedeemToyViewHandler(ctx echo.Context) error {
	toyIDStr := ctx.Param("toy_id")
	toyID, err := strconv.ParseInt(toyIDStr, 10, 64)
	toy, err := h.ToysStore.GetToyByID(toyID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	cr := new(CatalogHandlerRequest)
	if err := ctx.Bind(cr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	code := cr.Code

	return lib.Render(ctx, checkout.Checkout(toy, code))
}
