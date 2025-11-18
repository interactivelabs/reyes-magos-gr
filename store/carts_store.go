package store

import (
	"database/sql"
	"reyes-magos-gr/store/dtos"
	"reyes-magos-gr/store/models"
	"reyes-magos-gr/store/utils"
)

type LibSQLCartsStore struct {
	DB *sql.DB
}

func NewCartsStore(db *sql.DB) *LibSQLCartsStore {
	return &LibSQLCartsStore{DB: db}
}

type CartsStore interface {
	GetCartToys(volunteerID int64) (cartItems []dtos.CartItem, err error)
	CreateCartItem(item models.CartItem) (cartID int64, err error)
	DeleteCartItem(cartId int64) (err error)
}

func (r LibSQLCartsStore) GetCartToys(volunteerID int64) (cartItems []dtos.CartItem, err error) {
	rows, err := r.DB.Query(`
		SELECT `+cartItemFields+`
		FROM toys
		INNER JOIN carts ON carts.toy_id = toys.toy_id
		WHERE
			toys.deleted = 0
			AND carts.used = 0
			AND carts.deleted = 0
			AND volunteer_id = ?;`, volunteerID)

	if err != nil {
		return nil, err
	}

	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var toy dtos.CartItem
		toy, err = scanCartItem(rows)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, toy)
	}

	return cartItems, nil
}

func (r LibSQLCartsStore) CreateCartItem(item models.CartItem) (cartID int64, err error) {
	queryStr, params, err := utils.BuildInsertQuery(item, "carts")
	if err != nil {
		return 0, err
	}

	res, err := utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r LibSQLCartsStore) DeleteCartItem(cartID int64) error {
	queryStr, params, err := utils.BuildDeleteQuery(cartID, "carts", "cart_id")
	if err != nil {
		return err
	}

	_, err = utils.ExecuteMutationQuery(r.DB, queryStr, params...)
	if err != nil {
		return err
	}
	return nil
}

const cartItemFields string = `
	carts.cart_id,
	toys.toy_id,
	toys.toy_name,
	COALESCE(toys.category, ''),
	toys.image1`

func scanCartItem(s utils.Scanner) (cartItem dtos.CartItem, err error) {
	err = s.Scan(
		&cartItem.CartID,
		&cartItem.ToyID,
		&cartItem.ToyName,
		&cartItem.Category,
		&cartItem.Image1)
	return cartItem, err
}
