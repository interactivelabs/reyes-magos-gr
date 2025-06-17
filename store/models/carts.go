package models

type CartItem struct {
	CartID          int64 `json:"cart_id"`
	ToyID           int64 `json:"toy_id"`
	VolunteerID     int64 `json:"volunteer_id"`
	VolunteerCodeID int64 `json:"volunteer_code_id"`
	Used            int64 `json:"used"`
	Deleted         int64 `json:"deleted"`
}
