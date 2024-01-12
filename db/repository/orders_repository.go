package repository

import (
	"database/sql"
	"reyes-magos-gr/db/model"
)

type OrdersRepository struct {
	DB *sql.DB
}

func (r OrdersRepository) CreateOrder(order model.Order) (int64, error) {
	queryStr, params, err := buildInsertQuery(order, "orders")
	if err != nil {
		return 0, err
	}

	res, err := executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r OrdersRepository) UpdateOrder(order model.Order) error {
	queryStr, params, err := buildUpdateQuery(order, "orders", "order_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

func (r OrdersRepository) DeleteOrder(orderID int64) error {
	queryStr, params, err := buildDeleteQuery(orderID, "orders", "order_id")
	if err != nil {
		return err
	}

	_, err = executeQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}
