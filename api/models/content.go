package models

import (
	"time"
)

type Bucket struct {
	Dbname    string `json:dbname`
	Tablename string `json:"tablename"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

/* 一篇文章包含如下：
1. 文章的标题、用户、创建时间、修改时间、用户UID
2. 文章的内容
3. 留言板信息、留言用户信息
4. 留言用户信息包含：级别、威望、注册时间、财富、UID、称昵、回复时间、几楼
包括id、label、message
*/
type Info struct {
	Id       int64     `json:"id"`
	Index    []int     `json:"index"`
	Label    string    `json:"label"`
	Message  string    `json:"message"`
	UserName string    `json:"username"`
	UserId   string    `json:"userid"`
	Count    int64     `json:"count"`
	Create   time.Time `json:"create"`
	Update   time.Time `json:"update"`
	// MessageBoard []Info  `json:"messageboard"`
}

type BucketInfo struct {
	Dbname    string `json:dbname`
	Tablename string `json:"tablename"`
	Key       string `json:"key"`
	Value     Info   `json:"value"`
}

// 新增帖子或评论的models
type Message struct {
	Tablename string `json:"tablename"`
	Key       string `json:"key"`
	ParentId  string `json:"parentid"` // 1-1-0-1
	NewId     int64  `json:"newid"`
	Label     string `json:"label"`
	Value     Info   `json:"value"`
}
