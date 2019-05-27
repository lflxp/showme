package models

import "errors"

type AdminLogin struct {
	Code int64 `json:"code"`
	Data struct {
		Token  string   `json:"token"`
		Roles  []string `json:"roles"`
		Name   string   `json:"name"`
		Avatar string   `json:"avatar"`
	} `json:"data"`
}

type AdminUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (this *AdminUser) GetToken() (*AdminLogin, error) {
	if this.Username == "admin" && this.Password == "admin" {
		tmp := &AdminLogin{}
		tmp.Code = 20000
		tmp.Data.Token = "admin"
		return tmp, nil
	}
	return nil, errors.New("no such person" + this.Username)
}

func (this *AdminLogin) GetInfo() (*AdminLogin, error) {
	if this.Data.Token == "admin" {
		tmp := &AdminLogin{}
		tmp.Code = 20000
		tmp.Data.Roles = []string{"admin"}
		tmp.Data.Name = "admin"
		tmp.Data.Avatar = "http://www.baidu.com"
		return tmp, nil
	}
	return nil, errors.New("nothing to tell you")
}
