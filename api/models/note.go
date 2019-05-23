package models

import (
	"time"

	"github.com/astaxie/beego"
	pkg "github.com/lflxp/showme/utils"
)

func init() {
	err := pkg.Engine.Sync2(new(Group), new(Note), new(Common))
	if err != nil {
		beego.Critical(err.Error())
	}
	beego.Informational("初始化表 Group Note Common")
}

// 目录
type Group struct {
	Id        int64     `json:"id"`
	Gname     string    `xorm:"gname" json:"gname"`   // 目录名
	Parent    string    `xorm:"parent" json:"parent"` // 多级目录关系
	Created   time.Time `xorm:"created"`
	Updated   time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (this *Group) Insert() (int64, error) {
	n, err := pkg.Engine.Insert(this)
	return n, err
}

func (this *Group) Update(id string) (int64, error) {
	a, err := pkg.Engine.Id(id).Update(this)
	return a, err
}

// 顶级目录
func GetAllGroups() ([]Group, error) {
	all := make([]Group, 0)
	err := pkg.Engine.Where("parent = ?", "").Find(&all)
	return all, err
}

// 二级目录
func GetAllGroupsSecond(gname string) ([]Group, error) {
	all := make([]Group, 0)
	err := pkg.Engine.Where("parent = ?", gname).Find(&all)
	return all, err
}

func DeleteGroup(id string) (int64, error) {
	g := new(Group)
	n, err := pkg.Engine.Id(id).Delete(g)
	return n, err
}

func NewGroup() *Group {
	return &Group{}
}

// 笔记
type Note struct {
	Id        int64     `json:"id"`
	Nname     string    `xorm:"nname" json:"nname"`       // 笔记标题
	Gname     string    `xorm:"gname index" json:"gname"` // 组名
	Username  string    `xorm:"username" json:"username"` // 用户名
	Title     string    `xorm:"title" json:"title"`
	Data      string    `xorm:"data" json:"data"`   // 笔记内容
	Score     int64     `xorm:"score" json:"score"` // 文章分数 浏览次数
	Created   time.Time `xorm:"created" json:"created"`
	Updated   time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (this *Note) Insert() (int64, error) {
	n, err := pkg.Engine.Insert(this)
	_, err = pkg.Engine.Get(this)
	if err != nil {
		return n, err
	}
	// if has {
	// 	go AddNote(this)
	// }
	return n, err
}

func (this *Note) Update(id string) (int64, error) {
	a, err := pkg.Engine.Id(id).Update(this)
	return a, err
}

func DistinctNote(name string) ([]Note, error) {
	all := make([]Note, 0)
	err := pkg.Engine.Distinct(name).Find(&all)
	return all, err
}

func GetNoteById(id string) (*Note, bool, error) {
	note := new(Note)
	has, err := pkg.Engine.Id(id).Get(note)
	return note, has, err
}

func GetAllNotesWithoutId() ([]Note, error) {
	all := make([]Note, 0)
	err := pkg.Engine.Find(&all)
	return all, err
}

func GetAllNotes(gname string) ([]Note, error) {
	all := make([]Note, 0)
	err := pkg.Engine.Where("gname = ?", gname).Find(&all)
	return all, err
}

func DeleteNote(id string) (int64, error) {
	g := new(Note)
	n, err := pkg.Engine.Id(id).Delete(g)
	return n, err
}

func NewNote() *Note {
	return &Note{}
}

// 评论
type Common struct {
	Id        int64     ` json:"id"`
	Cname     string    `xorm:"cname index" json:"cname"`
	Type      string    `xorm:"type" json:"type"`   // 评论类型
	Nname     string    `xorm:"nname" json:"nname"` // 笔记名
	Created   time.Time `xorm:"created" json:"created"`
	Updated   time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`
}

func (this *Common) Insert() (int64, error) {
	n, err := pkg.Engine.Insert(this)
	return n, err
}

func (this *Common) Update(id string) (int64, error) {
	a, err := pkg.Engine.Id(id).Update(this)
	return a, err
}

func GetAllCommon(nname string) ([]Common, error) {
	all := make([]Common, 0)
	err := pkg.Engine.Where("nname = ?", nname).Find(&all)
	return all, err
}

func DeleteCommon(id string) (int64, error) {
	g := new(Common)
	n, err := pkg.Engine.Id(id).Delete(g)
	return n, err
}

func NewCommon() *Common {
	return &Common{}
}
