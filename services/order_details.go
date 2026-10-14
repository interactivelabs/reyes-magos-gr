package services

import (
	"database/sql"
	"errors"
	"reyes-magos-gr/store"
	"reyes-magos-gr/store/models"
)

// enrichOrder attaches the toy and code data shown alongside an order. Missing
// toys or codes are tolerated so a deleted row doesn't hide the order.
func enrichOrder(
	toysStore store.ToysStore,
	codesStore store.CodesStore,
	order models.Order,
) (models.OrderDetails, error) {
	toy, err := toysStore.GetToyByID(order.ToyID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.OrderDetails{}, err
	}

	code, err := codesStore.GetCodeByID(order.CodeID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return models.OrderDetails{}, err
	}

	return models.OrderDetails{
		Order:        order,
		ToyName:      toy.ToyName,
		ToySourceURL: toy.SourceURL,
		Code:         code.Code,
		Note:         code.Note,
	}, nil
}

func enrichOrders(
	toysStore store.ToysStore,
	codesStore store.CodesStore,
	orders []models.Order,
) ([]models.OrderDetails, error) {
	details := make([]models.OrderDetails, 0, len(orders))
	for _, order := range orders {
		detail, err := enrichOrder(toysStore, codesStore, order)
		if err != nil {
			return nil, err
		}
		details = append(details, detail)
	}
	return details, nil
}
