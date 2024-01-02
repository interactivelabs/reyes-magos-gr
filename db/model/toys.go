package model

type Toy struct {
	ToyID     int64  `json:"toy_id"`
	ToyName   string `json:"toy_name"`
	AgeMin    int64  `json:"age_min"`
	AgeMax    int64  `json:"age_max"`
	Image1    string `json:"image1"`
	Image2    string `json:"image2"`
	SourceURL string `json:"source_url"`
}
