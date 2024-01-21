package model

type VolunteerCode struct {
	VolunteerCodeID int64     `json:"volunteer_code_id"`
	VolunteerID     int64     `json:"volunteer_id"`
	CodeID          int64     `json:"code_id"`
	Deleted         int64     `json:"deleted"`
	Volunteer       Volunteer `json:"volunteer"`
	Code            Code      `json:"code"`
}
