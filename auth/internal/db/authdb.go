package db

type AuthData struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}
