package music

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/lflxp/lflxp-music/core/model/music"
	"github.com/lflxp/tools/httpclient"
	"github.com/lflxp/tools/orm/sqlite"
	"k8s.io/client-go/util/homedir"
)

func init() {
	sqlite.NewOrm().Sync2(new(music.Musichistory))
}

const (
	REPO_LIST       = "/music/history/list"
	REPO_ADD        = "/music/history/add"
	REPO_DELETE     = "/music/history/delete"
	REPO_UPLOAD     = "/music/local/upload"
	REPO_LOCAL_LIST = "/music/local/list"
	REPO_STATIC     = "/static"
)

func RegisterMusic(router *gin.Engine) {
	// FSStatic
	router.StaticFS(REPO_STATIC, http.Dir(filepath.Join(homedir.HomeDir(), "/.music")))
	shopGroup := router.Group("/api")
	{
		shopGroup.GET(REPO_LIST, repo_list)
		shopGroup.POST(REPO_ADD, repo_add)
		shopGroup.DELETE(REPO_DELETE, repo_delete)
		shopGroup.POST(REPO_UPLOAD, Upload)
		shopGroup.GET(REPO_LOCAL_LIST, music_local_list)
	}
}

func repo_list(c *gin.Context) {
	user := c.Request.Header.Get("username")
	data := make([]music.Musichistory, 0)
	err := sqlite.NewOrm().Where("user=?", user).Asc("id").Find(&data)
	if err != nil {
		httpclient.SendErrorMessage(c, 500, "查询错误", err.Error())
		return
	}

	// for i := 0; i < 10; i++ {
	// 	tmp := music.MusicHistory{
	// 		Id:       1997438791,
	// 		Duration: 257.251,
	// 		Album:    fmt.Sprintf("album%d", i),
	// 		Image:    "https://p3.music.126.net/dvBE3I5IYmDTmZq9SyKoRA==/109951168159576909.jpg",
	// 		Name:     fmt.Sprintf("name%d", i),
	// 		Url:      "https://music.163.com/song/media/outer/url?id=1997438791.mp3",
	// 		Singer:   fmt.Sprintf("singer%d", i),
	// 		User:     fmt.Sprintf("user%d", i),
	// 	}
	// 	tmp.Post()
	// }
	httpclient.SendSuccessMessage(c, 200, data)
}

func repo_add(c *gin.Context) {
	var data music.PostData
	if err := c.BindJSON(&data); err != nil {
		httpclient.SendErrorMessage(c, http.StatusBadRequest, "not found", err.Error())
		return
	}

	data.Param.Data.User = c.Request.Header.Get("username")

	list, err := data.Param.Data.Post()
	if err != nil {
		log.Errorf("新增音乐历史错误: %s", err.Error())
		httpclient.SendErrorMessage(c, 500, "新增音乐历史错误", err.Error())
		return
	}

	httpclient.SendSuccessMessage(c, 200, fmt.Sprintf("success add %d", list))
}

func repo_delete(c *gin.Context) {
	name, isok := c.GetQuery("id")
	if !isok {
		httpclient.SendErrorMessage(c, http.StatusBadRequest, "id not found", "id not found")
		return
	}

	user := c.Request.Header.Get("username")

	var data music.Musichistory
	n, err := sqlite.NewOrm().Where("name=? and user=?", name, user).Delete(&data)
	if err != nil {
		httpclient.SendErrorMessage(c, http.StatusBadRequest, "删除失败", err.Error())
		return
	}

	httpclient.SendSuccessMessage(c, 200, fmt.Sprintf("success delete %d", n))
}
