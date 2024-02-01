package model

import "fmt"

type UrlCache struct {
	Name        string // require 保证唯一性
	Method      string // require 接口请求方法 GET POST PUT DELETE等
	Description string // require 接口描述，前端按钮中文名称
	IsAudit     bool   // 是否开启审计，查询功能不要开启审计
	Path        string // require 接口路由字符
	Group       string // require 参照Group映射表，取groupMapping的key
}

func (u *UrlCache) Vaild() error {
	if u.Name == "" || u.Method == "" || u.Description == "" || u.Path == "" || u.Group == "" {
		return fmt.Errorf("存在必填字段为空 %v", u)
	}
	return nil
}
