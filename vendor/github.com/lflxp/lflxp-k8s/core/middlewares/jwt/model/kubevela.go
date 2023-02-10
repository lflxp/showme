package model

import "time"

type KubeVelaToken struct {
	User struct {
		CreateTime    time.Time `json:"createTime"`
		LastLoginTime time.Time `json:"lastLoginTime"`
		Name          string    `json:"name"`
		Email         string    `json:"email"`
		Disabled      bool      `json:"disabled"`
	} `json:"user"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
