/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/lflxp/screego/auth"
	"github.com/lflxp/screego/config"
	"github.com/lflxp/screego/logger"
	"github.com/lflxp/screego/router"
	"github.com/lflxp/screego/server"
	"github.com/lflxp/screego/turn"
	"github.com/lflxp/screego/ws"
	"github.com/lflxp/showme/utils"
	"github.com/rs/zerolog"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/client-go/util/homedir"
)

const screego_config = `SCREEGO_SECRET=secure
SCREEGO_LOG_LEVEL=debug
SCREEGO_CORS_ALLOWED_ORIGINS=http://localhost:3000
SCREEGO_USERS_FILE=./users
SCREEGO_EXTERNAL_IP=127.0.0.1`
const users = `# Password: admin
admin:$2a$12$kNgc2ZYAXzIL6SHY.8PHAOQ8Casi0s1bKatYoG/jupt2yV1M5K5nO
`

var (
	screego_hash bool
	screego_name string
	screego_pass string
)

// 设置screego环境变量
// 如果存在则不修改，不存在则设置
func checkScreeGoEnv(key, value string) error {
	val, ok := os.LookupEnv(key)
	if ok && val != "" {
		return nil
	}

	slog.Debug("Setting screego env", slog.Any("key", key), slog.Any("value", value))
	if key == "SCREEGO_USERS_FILE" {
		var path string
		if home := homedir.HomeDir(); home != "" {
			path = filepath.Join(home, ".screego_users")
			// 写文件
			if err := os.WriteFile(path, []byte(users), 0644); err != nil {
				slog.Error("Failed to write users file", slog.Any("path", path), slog.Any("err", err))
				return err
			}
			return os.Setenv(key, path)
		} else {
			return os.Setenv(key, value)
		}
	} else if key == "SCREEGO_EXTERNAL_IP" {
		ips := utils.GetIPs()
		return os.Setenv(key, ips[0])
	} else {
		return os.Setenv(key, value)
	}
}

// screegoCmd represents the screego command
var screegoCmd = &cobra.Command{
	Use:   "screego",
	Short: "web在线屏幕共享",
	Long: `基于浏览器的视频共享服务,
	默认端口【5050】
	默认账号密码【admin/admin】
	默认用户配置路径 【~/.screego_users】`,
	Run: func(cmd *cobra.Command, args []string) {
		if screego_hash {
			pass := []byte(screego_pass)
			if screego_name == "" {
				slog.Error("--name must be set")
			}

			hashedPw, err := bcrypt.GenerateFromPassword(pass, 12)
			if err != nil {
				slog.Error("could not generate password", err)
			}

			fmt.Printf("%s:%s", screego_name, string(hashedPw))
			fmt.Println("")
		} else {
			checkScreeGoEnv("SCREEGO_SECRET", "secure")
			checkScreeGoEnv("SCREEGO_LOG_LEVEL", "debug")
			checkScreeGoEnv("SCREEGO_CORS_ALLOWED_ORIGINS", "http://localhost:3000")
			checkScreeGoEnv("SCREEGO_USERS_FILE", "./users")
			checkScreeGoEnv("SCREEGO_EXTERNAL_IP", "127.0.0.1")
			conf, errs := config.Get()
			logger.Init(conf.LogLevel.AsZeroLogLevel())

			exit := false
			for _, err := range errs {
				slog.Error("errors level %s error: %s", err.Level, err.Msg)
				exit = exit || err.Level == zerolog.FatalLevel || err.Level == zerolog.PanicLevel
			}
			if exit {
				os.Exit(1)
			}

			if _, _, err := conf.TurnIPProvider.Get(); err != nil {
				// error is already logged by .Get()
				os.Exit(1)
			}

			users, err := auth.ReadPasswordsFile(conf.UsersFile, conf.Secret, conf.SessionTimeoutSeconds)
			if err != nil {
				slog.Error("errors", err)
			}

			auth, err := turn.Start(conf)
			if err != nil {
				slog.Error("errors", err)
			}

			rooms := ws.NewRooms(auth, users, conf)

			go rooms.Start()

			// 自动打开浏览器
			url := "http://127.0.0.1:5050"
			open.Start(url)
			r := router.Router(conf, rooms, users, version)
			if err := server.Start(r, conf.ServerAddress, conf.TLSCertFile, conf.TLSKeyFile); err != nil {
				slog.Error("errors", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(screegoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// screegoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// screegoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	screegoCmd.Flags().BoolVarP(&screego_hash, "hash", "H", false, "是否开启hash密码生成")
	screegoCmd.Flags().StringVarP(&screego_name, "name", "n", "admin", "登录用户名")
	screegoCmd.Flags().StringVarP(&screego_pass, "pass", "p", "admin", "登录密码")
}
