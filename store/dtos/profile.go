package dtos

type Profile struct {
	IsAdmin    bool
	IsLoggedIn bool
	Nickname   string
	Email      string
	Picture    string
}
