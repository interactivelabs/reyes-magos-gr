package models

type CartItem struct {
	CartID      int64 `json:"cart_id"`
	CodeID      int64 `json:"code_id"`
	ToyID       int64 `json:"toy_id"`
	VolunteerID int64 `json:"volunteer_id"`
	Used        int64 `json:"used"`
	Deleted     int64 `json:"deleted"`
}
