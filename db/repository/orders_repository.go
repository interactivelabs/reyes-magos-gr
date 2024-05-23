package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
	utils "reyes-magos-gr/db/repository/utils"
)

type OrdersRepository struct {
	DB *sql.DB
}

func (r OrdersRepository) CreateOrder(order model.Order) (int64, error) {
	queryStr, params, err := utils.BuildInsertQuery(order, "orders")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r OrdersRepository) UpdateOrder(order model.Order) error {
	queryStr, params, err := utils.BuildUpdateQuery(order, "orders", "order_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r OrdersRepository) DeleteOrder(orderID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(orderID, "orders", "order_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r OrdersRepository) GetOrderByID(orderID int64) (model.Order, error) {
	row := r.DB.QueryRow(`
		SELECT `+orderAllFields+`
		FROM orders
		WHERE order_id = ?
	`, orderID)
	return scanAllOrder(row)
}

func (r OrdersRepository) GetOrdersByVolunteerID(volunteerID int64) ([]model.Order, error) {
	rows, err := r.DB.Query(`
		SELECT `+orderAllFields+`
		FROM orders
		WHERE volunteer_id = ?
	`, volunteerID)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var orders []model.Order
	for rows.Next() {
		order, err := scanAllOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r OrdersRepository) GetAllActiveOrders() ([]model.Order, error) {
	rows, err := r.DB.Query(`
		SELECT ` + orderAllFields + `
		FROM orders
		WHERE deleted = 0 AND completed = 0`)
	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var orders []model.Order
	for rows.Next() {
		order, err := scanAllOrder(rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

const orderAllFields string = `
	order_id,
	toy_id,
	volunteer_id,
	code_id,
	order_date,
	shipped,
	COALESCE(shipped_date, ''),
	completed,
	deleted`

type OrderScanner interface {
	Scan(dest ...interface{}) error
}

func scanAllOrder(s OrderScanner) (model.Order, error) {
	var order model.Order
	err := s.Scan(
		&order.OrderID,
		&order.ToyID,
		&order.VolunteerID,
		&order.CodeID,
		&order.OrderDate,
		&order.Shipped,
		&order.ShippedDate,
		&order.Completed,
		&order.Deleted,
	)
	return order, err
}
