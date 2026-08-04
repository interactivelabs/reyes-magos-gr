package models

type OrderDetails struct {
	Order        Order
	ToyName      string
	ToySourceURL string
	Code         string
}
