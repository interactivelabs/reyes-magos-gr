package model

type Volunteer struct {
	VolunteerID int64  `json:"volunteer_id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	Address2    string `json:"address2"`
	Country     string `json:"country"`
	State       string `json:"state"`
	City        string `json:"city"`
	Province    string `json:"province"`
	ZipCode     string `json:"zip_code"`
	Secret      string `json:"secret"`
}
