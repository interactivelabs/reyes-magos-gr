package model

type Code struct {
	CodeID     int64  `json:"code_id"`
	Code       string `json:"code"`
	Expiration string `json:"expiration"`
	Used       int64  `json:"used"`
	Cancelled  int64  `json:"cancelled"`
	Deleted    int64  `json:"deleted"`
}
