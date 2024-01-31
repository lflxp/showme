package music

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bogem/id3v2/v2"
	"github.com/lflxp/lflxp-music/core/middlewares/template"
	"github.com/lflxp/tools/orm/sqlite"
	"k8s.io/client-go/util/homedir"
)

func init() {
	template.Register(new(Musichistory))
}

type Musichistory struct {
	Id       int64   `xorm:"id" name:"id" json:"id"`
	Duration float64 `xorm:"duration" name:"duration" verbose_name:"duration字段测试" search:"true" json:"duration"`
	Album    string  `xorm:"album" name:"album" json:"album"`
	Image    string  `xorm:"image" name:"image" json:"image"`
	Name     string  `xorm:"name notnull unique" name:"name" json:"name"`
	Url      string  `xorm:"url" name:"url" json:"url"`
	Singer   string  `xorm:"singer" name:"singer" json:"singer"`
	User     string  `xorm:"user" name:"user" json:"user"`
}

func (m *Musichistory) Post() (int64, error) {
	return sqlite.NewOrm().Insert(m)
}

// download
func (m *Musichistory) Download() (error, bool) {
	path := fmt.Sprintf("%s/%s.mp3", filepath.Join(homedir.HomeDir(), fmt.Sprintf(".music/%s", m.User)), m.Name)
	_, err := os.Lstat(path)
	if !os.IsNotExist(err) {
		// 返回false表示没有进行下载
		return errors.New("文件已存在"), false
	}

	res, err := http.Get(m.Url)
	if err != nil {
		slog.Error(err.Error())
		return err, false
	}

	f, err := os.Create(path)
	if err != nil {
		slog.Error(err.Error())
		return err, false
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		slog.Error(err.Error())
		return err, false
	} else {
		// 写入信息
		tag, err := id3v2.Open(path, id3v2.Options{Parse: true})
		if err != nil {
			slog.Error("Error while opening mp3 file", "error", err.Error())
		}
		defer tag.Close()
		tag.SetTitle(m.Name)
		tag.SetAlbum(m.Album)
		tag.SetArtist(m.Singer)
		tag.SetGenre(m.Image)
		tag.SetYear(fmt.Sprintf("%.2f", m.Duration))
	}

	return nil, true
}

type PostData struct {
	Param struct {
		Data Musichistory `json:"data"`
	} `json:"param"`
}
