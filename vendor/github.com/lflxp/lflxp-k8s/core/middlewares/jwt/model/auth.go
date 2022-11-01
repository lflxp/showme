package model

type Auth struct {
	Description  string `json:"description"`
	Password     string `json:"password"`
	Username     string `json:"username"`
	ResponseType string `json:"responseType"`
}
