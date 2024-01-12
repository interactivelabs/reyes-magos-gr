package model

type Order struct {
	OrderID     int64  `json:"order_id"`
	ToyID       int64  `json:"toy_id"`
	VolunteerID int64  `json:"volunteer_id"`
	CodeID      int64  `json:"code_id"`
	OrderDate   string `json:"order_date"`
	Shipped     int64  `json:"shipped"`
	ShippedDate string `json:"shipped_date"`
	Deleted     int64  `json:"deleted"`
}
