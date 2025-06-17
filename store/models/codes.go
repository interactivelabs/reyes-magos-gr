package models

type Code struct {
	CodeID     int64  `json:"code_id"`
	Code       string `json:"code"`
	Expiration string `json:"expiration"`
	Given      int64  `json:"given"`
	Used       int64  `json:"used"`
	Cancelled  int64  `json:"cancelled"`
	Deleted    int64  `json:"deleted"`
}
