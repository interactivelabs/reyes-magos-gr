package store

import (
	"database/sql"
	"reyes-magos-gr/store/models"
	utils "reyes-magos-gr/store/utils"
)

type LibSQLOrdersStore struct {
	DB *sql.DB
}

func NewOrdersStore(db *sql.DB) *LibSQLOrdersStore {
	return &LibSQLOrdersStore{DB: db}
}

type OrdersStore interface {
	CreateOrder(order models.Order) (id int64, err error)
	CreateOrders(orders []models.Order) error
	UpdateOrder(order models.Order) error
	DeleteOrder(orderID int64) error
	GetOrderByID(orderID int64) (order models.Order, err error)
	GetPendingOrdersByVolunteerID(volunteerID int64) (orders []models.Order, err error)
	GetAllActiveOrders() (orders []models.Order, err error)
	GetCompletedOrders() (orders []models.Order, err error)
}

func (r LibSQLOrdersStore) CreateOrder(order models.Order) (id int64, err error) {
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

func (r LibSQLOrdersStore) CreateOrders(orders []models.Order) error {
	if len(orders) == 0 {
		return nil
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, order := range orders {
		queryStr, params, err := utils.BuildInsertQuery(order, "orders")
		if err != nil {
			return err
		}

		_, err = tx.Exec(queryStr, params...)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r LibSQLOrdersStore) UpdateOrder(order models.Order) error {
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

func (r LibSQLOrdersStore) DeleteOrder(orderID int64) error {
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

func (r LibSQLOrdersStore) GetOrderByID(orderID int64) (order models.Order, err error) {
	row := r.DB.QueryRow(`
		SELECT `+orderAllFields+`
		FROM orders
		WHERE order_id = ?
	`, orderID)
	return scanAllOrder(row)
}

func (r LibSQLOrdersStore) GetPendingOrdersByVolunteerID(
	volunteerID int64,
) (orders []models.Order, err error) {
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

func (r LibSQLOrdersStore) GetAllActiveOrders() (orders []models.Order, err error) {
	rows, err := r.DB.Query(`
		SELECT ` + orderAllFields + `
		FROM orders
		WHERE completed = 0`)
	if err != nil {
		return nil, err
	}

	return GetOrdersFromQuery(rows)
}

func (r LibSQLOrdersStore) GetCompletedOrders() (orders []models.Order, err error) {
	rows, err := r.DB.Query(`
		SELECT ` + orderAllFields + `
		FROM orders
		WHERE completed = 1`)
	if err != nil {
		return nil, err
	}

	return GetOrdersFromQuery(rows)
}

func GetOrdersFromQuery(rows *sql.Rows) (orders []models.Order, err error) {
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

func scanAllOrder(s utils.Scanner) (order models.Order, err error) {
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
