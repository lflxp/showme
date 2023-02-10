package model

// 用户表
type User struct {
	Id           int64  `json:"id" xorm:"id pk not null autoincr" name:"id"`
	Username     string `json:"username" xorm:"username" name:"username" verbose_name:"用户名" list:"true" search:"true"`
	Password     string `json:"password" xorm:"password" name:"password" verbose_name:"密码" colType:"password" list:"true" search:"true"`
	Name         string `json:"name" xorm:"name" name:"name" verbose_name:"名字" list:"true" search:"true"`
	FirstName    string `json:"firstName" xorm:"firstname" name:"firstname" verbose_name:"姓氏" list:"true" search:"true"`
	Email        string `json:"email" xorm:"email" name:"email" verbose_name:"电子邮件" list:"true" search:"true"`
	IsVaild      string `json:"isVaild" xorm:"isvaild" name:"isvaild" verbose_name:"有效" list:"true" search:"false" colType:"radio" radio:"有效|1,无效|0"`
	Status       string `json:"status" xorm:"status" name:"status" verbose_name:"状态" list:"true" search:"false" colType:"radio" radio:"有效|1,无效|0"`
	IsAdmin      string `json:"isAdmin" xorm:"isadmin" name:"isadmin" verbose_name:"超级用户状态" list:"true" search:"false" colType:"radio" radio:"是|1,不是|0"`
	Token        string `json:"token" xorm:"token" name:"token" verbose_name:"rancher token"`
	Tenant       string `json:"tenant"`
	AuthProvider string `json:"authProvider"`
	UserId       string `json:"userId"`
	Role         string `json:"role"`
	RoleLevel    string `json:"roleLevel"`
	RoleReal     string `json:"roleReal"`
	IsGlobal     string `json:"isGlobal"`
	DisplayName  string `json:"displayName"` // 用户显示名称
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// Resp 登录后返回resp结构体
type Resp struct {
	BaseType string `json:"baseType"`
	Code     string `json:"code"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
	Type     string `json:"type"`
}
