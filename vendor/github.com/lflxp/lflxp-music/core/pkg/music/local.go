package music

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bogem/id3v2/v2"
	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/lflxp/lflxp-music/core/model/music"
	"github.com/lflxp/lflxp-music/core/utils"
	"github.com/lflxp/tools/httpclient"
	"k8s.io/client-go/util/homedir"
)

func init() {
	// 创建music文件夹
	home := homedir.HomeDir()
	if _, err := utils.IsExistAndCreateDir(filepath.Join(home, ".music")); err != nil {
		panic(err)
	}
}

// 下载mp3文件
func download(c *gin.Context) {
	var data music.PostData
	if err := c.BindJSON(&data); err != nil {
		httpclient.SendErrorMessage(c, http.StatusBadRequest, "not found", err.Error())
		return
	}

	data.Param.Data.User = c.Request.Header.Get("username")

	if err, isok := data.Param.Data.Download(); !isok {
		httpclient.SendErrorMessage(c, 200, "文件下载失败或已存在", err.Error())
	} else {
		httpclient.SendSuccessMessage(c, 200, "success download")
	}
}

func music_local_list(c *gin.Context) {
	username := c.Request.Header.Get("username")
	home := homedir.HomeDir()
	types := strings.Split(".avi,.wma,.rmvb,.rm,.wav,.mp3,.mp4,.mov,.3gp,.mpeg,.mpg,.mpe,.m4v,.mkv,.flv,.vob,.wmv,.asf,.asx", ",")

	if _, err := utils.IsExistAndCreateDir(filepath.Join(home, fmt.Sprintf(".music/%s", username))); err != nil {
		httpclient.SendErrorMessage(c, 500, "文件夹不存在或者创建文件夹失败", err.Error())
		return
	}

	filelist, err := GetAllFiles(filepath.Join(home, fmt.Sprintf(".music/%s", username)), types)
	if err != nil {
		httpclient.SendErrorMessage(c, 500, "获取音乐文件列表错误", fmt.Sprintf("PATH: %s ERR: %s", filepath.Join(home, fmt.Sprintf(".music/%s", username)), err.Error()))
		return
	}

	var list []music.Musichistory
	if len(filelist) > 0 {
		list = []music.Musichistory{}
		for index, x := range filelist {
			tag, err := id3v2.Open(fmt.Sprintf("%s/%s", filepath.Join(home, fmt.Sprintf(".music/%s", username)), x), id3v2.Options{Parse: true})
			if err != nil {
				log.Fatal("Error while opening mp3 file: ", err)
			}
			defer tag.Close()

			if tag.Title() == "" {
				tag.SetTitle(x)
			}

			second := 300.01
			year := tag.Year()
			if year != "" {
				second, err = strconv.ParseFloat(year, 64)
				if err != nil {
					log.Error(err)
					second = 300
				}
			}

			list = append(list, music.Musichistory{
				Id:       int64(index),
				Album:    tag.Album(),
				Duration: second,
				Image:    "https://p3.music.126.net/YglUhn-RRq6KM7Dfm6VUZw==/109951168255550269.jpg",
				Name:     strings.Replace(tag.Title(), ".mp3", "", -1),
				Singer:   tag.Artist(),
				Url:      fmt.Sprintf("/static/%s/%s", username, x),
				// Url: "https://music.163.com/song/media/outer/url?id=1997438791.mp3",
			})
		}
	}

	httpclient.SendSuccessMessage(c, 200, list)
}

func Upload(c *gin.Context) {
	username := c.Request.Header.Get("username")
	home := homedir.HomeDir()

	// 判断文件夹是否存在
	if _, err := utils.IsExistAndCreateDir(filepath.Join(home, fmt.Sprintf(".music/%s", username))); err != nil {
		httpclient.SendErrorMessage(c, 500, "文件夹不存在或者创建文件夹失败", err.Error())
		return
	}

	// 多文件
	form, _ := c.MultipartForm()
	files := form.File["file"]

	for _, file := range files {

		// 上传文件到指定的路径
		if err := c.SaveUploadedFile(file, fmt.Sprintf("%s/%s", filepath.Join(home, fmt.Sprintf(".music/%s", username)), file.Filename)); err != nil {
			httpclient.SendErrorMessage(c, http.StatusBadRequest, "upload err", err.Error())
			return
		}
	}
	c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
}

// 获取指定目录下的所有文件,包含子目录下的文件
// strings.Replace(dirPth+PthSep+fi.Name(), ".", "/static", 1)
// 只能查找当前目录下的所有视频文件
func GetAllFiles(dirPth string, suff []string) (files []string, err error) {
	// var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			// dirs = append(dirs, dirPth+PthSep+fi.Name())
			// GetAllFiles(dirPth+PthSep+fi.Name(), suff)
			fmt.Println("不递归")
		} else {
			// 过滤指定格式
			for _, x := range suff {
				ok := strings.HasSuffix(fi.Name(), x)
				if ok {
					files = append(files, fi.Name())
				}
			}

		}
	}

	// 读取子目录下文件
	// for _, table := range dirs {
	// 	temp, _ := GetAllFiles(table, suff)
	// 	for _, temp1 := range temp {
	// 		files = append(files, temp1)
	// 	}
	// }

	return files, nil
}
