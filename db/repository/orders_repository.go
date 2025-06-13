package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
	utils "reyes-magos-gr/db/repository/utils"
)

type OrdersRepository struct {
	DB *sql.DB
}

func (r OrdersRepository) CreateOrder(order model.Order) (id int64, err error) {
	queryStr, params, err := utils.BuildInsertQuery(order, "orders")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteMutationQuery(r.DB, queryStr, params...)
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

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
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

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r OrdersRepository) GetOrderByID(orderID int64) (order model.Order, err error) {
	row := r.DB.QueryRow(`
		SELECT `+orderAllFields+`
		FROM orders
		WHERE order_id = ?
	`, orderID)
	return scanAllOrder(row)
}

func (r OrdersRepository) GetPendingOrdersByVolunteerID(volunteerID int64) (orders []model.Order, err error) {
	rows, err := r.DB.Query(`
		SELECT `+orderAllFields+`
		FROM orders
		WHERE volunteer_id = ?
			AND deleted = 0
			AND completed = 0
	`, volunteerID)
	if err != nil {
		return nil, err
	}

	return GetOrdersFromQuery(rows)
}

func (r OrdersRepository) GetAllActiveOrders() (orders []model.Order, err error) {
	rows, err := r.DB.Query(`
		SELECT ` + orderAllFields + `
		FROM orders
		WHERE completed = 0`)
	if err != nil {
		return nil, err
	}

	return GetOrdersFromQuery(rows)
}

func (r OrdersRepository) GetCompletedOrders() (orders []model.Order, err error) {
	rows, err := r.DB.Query(`
		SELECT ` + orderAllFields + `
		FROM orders
		WHERE completed = 1`)
	if err != nil {
		return nil, err
	}

	return GetOrdersFromQuery(rows)
}

func GetOrdersFromQuery(rows *sql.Rows) (orders []model.Order, err error) {
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

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
	COALESCE(completed_date, ''),
	cancelled,
	deleted`

func scanAllOrder(s utils.Scanner) (order model.Order, err error) {
	err = s.Scan(
		&order.OrderID,
		&order.ToyID,
		&order.VolunteerID,
		&order.CodeID,
		&order.OrderDate,
		&order.Shipped,
		&order.ShippedDate,
		&order.Completed,
		&order.CompletedDate,
		&order.Cancelled,
		&order.Deleted,
	)
	return order, err
}
