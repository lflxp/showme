package admin

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lflxp/lflxp-music/core/middlewares/template"

	log "github.com/go-eden/slf4go"
	"github.com/lflxp/tools/orm/sqlite"
)

func init() {
	// vpn := Vpn{}s
	template.Register(new(Vpn), new(Machine), new(Cdn), new(More), new(User), new(Claims), new(Groups), new(Userauth), new(History))

	user := User{Username: "admin"}
	has, err := sqlite.NewOrm().Get(&user)
	if err != nil {
		log.Error(err)
	}

	if !has {
		claims := Claims{
			Auth:  "admin",
			Type:  "nav",
			Value: "dashboard",
		}

		sqlite.NewOrm().Insert(&claims)

		log.Info("init admin user")
		sql := "insert into user('username','password','claims_id') values ('admin','admin','1');"
		n, err := sqlite.NewOrm().Query(sql)
		if err != nil {
			log.Errorf("init admin user err %s", err.Error())
		}
		log.Infof("insert admin user count: %d", len(n))
	}
}

/*
name  字段名
verbose_name 标识
list_display 显示字段
search_fields 查询字段
manytomany 一对多字段 指定表明
colType 字段类型 -> string|int|file|textarea|radio|m2m|otm|o2o|time|select|multiselect|password
radio|select -> Name|value,Name|value,...
o2o -> "tablename|showColumns,showColumns" -> first columns is id
*/
type Vpn struct {
	Id   int64  `xorm:"id notnull unique pk autoincr" name:"id"`
	Vpn  string `xorm:"vpn" name:"vpn" verbose_name:"Vpn字段测试" list:"true" search:"true"`
	Name string `xorm:"name" name:"name" verbose_name:"姓名" list:"true" search:"false"`
	Ip   string `xorm:"ip" name:"ip" verbose_name:"ip信息" list:"true" search:"false"`
}

type Cdn struct {
	Id           int64     `xorm:"id pk not null autoincr" name:"id" search:"true"`
	Num          int64     `xorm:"num" verbose_name:"数字" name:"num" colType:"int" list:"true" search:"true"`
	Cdn_name     string    `xorm:"cdn_name" name:"cdn_name" verbose_name:"cdn的名称" search:"true"`
	Type         string    `xorm:"type" name:"type" verbose_name:"类型" search:"false" colType:"textarea"`
	Detail       string    `xorm:"detail" name:"detail" verbose_name:"VPN信息" list:"false" search:"false" o2m:"vpn|id,vpn" colType:"o2m"`
	Radio        string    `xorm:"raidodas" name:"raidodas" verbose_name:"Radio单选" list:"true" search:"false" colType:"radio" radio:"男|man,女|girl,人妖|none"`
	Select       string    `xorm:"ss" name:"ss" verbose_name:"Select单选固定" list:"true" search:"false" colType:"select" select:"男11111111111111111111111111|man,女|girl,人妖|none"`
	MultiSelect  string    `xorm:"ss1" name:"ss1" verbose_name:"Multiselect多选" list:"true" search:"false" colType:"multiselect" multiselect:"男|man,女|girl,人妖|none,中|zhong,国|guo,人|ren,重|chong,Qing|qing"`
	MultiSelect2 string    `xorm:"ss2" name:"ss2" verbose_name:"Multiselect多选" list:"true" search:"false" colType:"multiselect" multiselect:"男|man,女|girl,人妖|none,中|zhong,国|guo,人|ren,重|chong,Qing|qing"`
	Files        string    `xorm:"file" name:"file" verbose_name:"cdn的名称" search:"true" colType:"file"`
	Times        time.Time `xorm:"times" name:"times" verbose_name:"时间" list:"true" search:"true"`
	Create       time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	Update       time.Time `xorm:"updated"` //这个Field将在Insert或Update时自动赋值为当前时间
}

type Machine struct {
	Id     int64     `xorm:"id pk not null autoincr" name:"id" search:"true"`
	Sn     string    `xorm:"sn" name:"sn" search:"true"`
	Mac    string    `xorm:"mac" name:"mac" search:"true"`
	Ip     string    `xorm:"ip" name:"ip" search:"true"`
	Name   string    `xorm:"name" name:"name" search:"true"`
	Create time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	Update time.Time `xorm:"updated"` //这个Field将在Insert或Update时自动赋值为当前时间
}

type More struct {
	Uid      int64  `xorm:"id pk not null autoincr" name:"id"`
	Username string `xorm:"username unique" name:"username" search:"true"`
	Alias    string `xorm:"-"`
	Vpn      `xorm:"vpn_id int(11)" colType:"o2o" o2o:"vpn|id,name,ip,vpn" verbose_name:"vpn外键" name:"vpn_id"`
	MoreVpn  string `xorm:"more" colType:"o2m" o2m:"vpn|id,name,ip,vpn" verbose_name:"vpn一对多" name:"more"` //id1,id2,id3,id4
}

type User struct {
	Id        int64    `xorm:"id pk not null autoincr" name:"id"`
	Username  string   `xorm:"username" name:"username" verbose_name:"用户名" list:"true" search:"true"`
	Password  string   `xorm:"password" name:"password" verbose_name:"密码" colType:"password" list:"true" search:"true"`
	Name      string   `xorm:"name" name:"name" verbose_name:"名字" list:"true" search:"true"`
	FirstName string   `xorm:"firstname" name:"firstname" verbose_name:"姓氏" list:"true" search:"true"`
	Email     string   `xorm:"email" name:"email" verbose_name:"电子邮件" list:"true" search:"true"`
	IsVaild   string   `xorm:"isvaild" name:"isvaild" verbose_name:"有效" list:"true" search:"false" colType:"radio" radio:"有效|1,无效|0"`
	Status    string   `xorm:"status" name:"status" verbose_name:"状态" list:"true" search:"false" colType:"radio" radio:"有效|1,无效|0"`
	IsAdmin   string   `xorm:"isadmin" name:"isadmin" verbose_name:"超级用户状态" list:"true" search:"false" colType:"radio" radio:"是|1,不是|0"`
	Claims    []Claims `xorm:"claims_id int(11)" colType:"o2m" o2m:"claims|id,auth,type,value" verbose_name:"权限配置" name:"claims_id"`
	Token     string   `xorm:"token" name:"token" verbose_name:"rancher token"`
}

// 用户权限表
type Claims struct {
	Id    int64  `xorm:"id pk not null autoincr" name:"id"`
	Auth  string `xorm:"auth varchar(255) unique(only)" name:"auth" verbose_name:"权限" list:"true" search:"true"`               // 对应Auth => Username  eg: admin
	Type  string `json:"type" xorm:"type varchar(255) unique(only)" name:"type" verbose_name:"类型" list:"true" search:"true"`   // 权限类型 eg: nav
	Value string `json:"value" xorm:"value varchar(255) unique(only)" name:"value" verbose_name:"值" list:"true" search:"true"` // 权限指 eg: dashboard
}

type Groups struct {
	Id   int64  `xorm:"id pk not null autoincr" name:"id"`
	Name string `xorm:"name" name:"name" verbose_name:"名称" list:"true" search:"true"`
	Auth string `xorm:"auth" name:"auth" verbose_name:"权限" colType:"o2m" o2m:"userauth|name,group,content"`
	User `xorm:"user" colType:"o2m" o2m:"user|id,username,name,email" verbose_name:"用户组" name:"user"`
}

type Userauth struct {
	Id      int64  `xorm:"id pk not null autoincr" name:"id"`
	Name    string `xorm:"name" name:"name" verbose_name:"名称" list:"true" search:"true"`
	Group   string `xorm:"group" name:"group" verbose_name:"分组" list:"true" search:"true"`
	Content string `xorm:"content" name:"content" verbose_name:"内容" list:"false" search:"false"`
}

// 继承xorm.Engine 扩展InsertHistory自动插入History历史
func InsertHistory(beans ...interface{}) (int64, error) {
	defer func() {
		data, err := json.Marshal(beans)
		if err != nil {
			log.Error(err)
			return
		}
		info := History{
			Name:   "FHST操作",
			Op:     fmt.Sprintf("添加 %s", string(data)),
			Common: "fhst op",
		}
		_, err = AddHistory(&info)
		if err != nil {
			log.Error(err)
		}

	}()
	return sqlite.NewOrm().Insert(beans...)
}

type History struct {
	Id     int64     `xorm:"id pk not null autoincr" name:"id"`
	Name   string    `xorm:"name" name:"name" verbose_name:"操作历史" list:"true"`
	Op     string    `xorm:"op" name:"op" verbose_name:"操作"`
	Common string    `xorm:"common" name:"common" verbose_name:"备注"`
	Create time.Time `xorm:"created"` //这个Field将在Insert时自动赋值为当前时间
	Update time.Time `xorm:"updated"` //这个Field将在Insert或Update时自动赋值为当前时间
}

func getByUUIDHistory(uuid string) (*History, bool, error) {
	data := new(History)
	has, err := sqlite.NewOrm().Where("uuid = ?", uuid).Get(data)
	return data, has, err
}

func AddHistory(data *History) (int64, error) {
	affected, err := sqlite.NewOrm().Insert(data)
	return affected, err
}

func DelHistory(id string) (int64, error) {
	data := new(History)
	affected, err := sqlite.NewOrm().ID(id).Delete(data)
	return affected, err
}

func UpdateHistory(id string, data *History) (int64, error) {
	affected, err := sqlite.NewOrm().Table(new(History)).ID(id).Update(data)
	return affected, err
}
