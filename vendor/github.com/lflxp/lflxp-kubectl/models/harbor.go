package models

type HarborUsers struct {
	Uername        string `json:"username"`
	Email          string `json:"email"`
	Comment        string `json:"comment"`
	Password       string `json:"password"`
	Realname       string `json:"realname"`
	Has_admin_role int64  `json:"has_admin_role"`
}
