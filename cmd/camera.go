// What it does:
//
// This example opens a video capture device, then streams MJPEG from it.
// Once running point your browser to the hostname/port you passed in the
// command line (for example http://localhost:8080) and you should see
// the live video stream.
//
// How to run:
//
// mjpeg-streamer [camera ID] [host:port]
//
//		go get -u github.com/hybridgroup/mjpeg
// 		go run ./cmd/mjpeg-streamer/main.go 1 0.0.0.0:8080
//
// # 安装

// 参考：https://gocv.io/getting-started/macos/

// # 摄像头启动命令

// - go get -u -d gocv.io/x/gocv
// - brew upgrade opencv
// - brew install opencv
// - brew install pkgconfig
// - go run ./cmd/version/main.go

// # 人脸识别

// - brew install tbb numpy vtk
// - brew install opencv
// - go run face-detect.go 0 data/haarcascade_frontalface_default.xml

// ## 参考：

// - https://gocv.io/writing-code/face-detect/
// - https://github.com/hybridgroup/gocv/tree/release/cmd
// - https://gocv.io/writing-code/more-examples/
// - https://kebingzao.com/2021/05/27/opencv-gocv/

package cmd

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	log "log/slog"

	"github.com/hybridgroup/mjpeg"
	"github.com/skratchdot/open-golang/open"
	"gocv.io/x/gocv"

	"github.com/lflxp/showme/asset"
	"github.com/lflxp/showme/utils"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/homedir"
)

const MinimumArea = 3000
const videoName = "output.avi"

var (
	cam_d        string
	cam_a        string
	cam_x        string
	cam_motion   bool
	cam_windows  bool
	cam_web      bool
	cam_pic      bool
	cam_video    bool
	cam_detect   bool
	err          error
	webcam       *gocv.VideoCapture
	stream       *mjpeg.Stream
	targetPath   string
	email_user   string
	email_pwd    string
	email_sendto string
	email_host   string
	email_port   int
	cam_count    int // 摄像头计数
)

type Motion string

const (
	MotionReady  Motion = "Ready"
	MotionDetect Motion = "Motion Detect"
)

// cameraCmd represents the camera command
var cameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "opencv本地视频流处理",
	Long: `前置条件：编译安装opencv和gocv 
	安装文档:
	1.【https://gocv.io/getting-started/】
	2. [https://kebingzao.com/2021/05/27/opencv-gocv/]
	支持：人脸识别｜人脸侦测｜运动侦测等
	支持：webcam ｜ window GUI`,
	Run: func(cmd *cobra.Command, args []string) {
		if cam_detect {
			log.Debug("请指定人脸识别模型文件")
			_, err := os.Stat(cam_x)
			if os.IsNotExist(err) {
				// 制造数据
				if home := homedir.HomeDir(); home != "" {
					targetPath = filepath.Join(home, ".face.xml")
				} else {
					targetPath = "./face.xml"
				}

				f, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
				if err != nil {
					log.Error("create file error", "file", targetPath, "error", err.Error())
					return
				}
				defer f.Close()

				_, err = f.Write(asset.FaceXml)
				if err != nil {
					log.Error("write file error", "file", targetPath, "error", err.Error())
					return
				}
			} else {
				log.Error("file not found", "file", cam_x)
				return
			}
		}

		// open webcam
		webcam, err = gocv.OpenVideoCapture(cam_d)
		if err != nil {
			fmt.Printf("Error opening capture device: %v\n", cam_d)
			return
		}
		defer webcam.Close()

		if !cam_windows {
			// create the mjpeg stream
			stream = mjpeg.NewStream()
		}

		if cam_motion {
			// start motion capturing
			go motionCapture()
		} else {
			// start capturing
			go mjpegCapture()
		}

		if cam_web {
			log.Info("Capturing. Point your browser to " + cam_a)

			// start http server
			http.Handle("/", stream)

			server := &http.Server{
				Addr:         cam_a,
				ReadTimeout:  60 * time.Second,
				WriteTimeout: 60 * time.Second,
			}

			var url string
			if strings.HasPrefix(cam_a, "http") {
				url = cam_a
			} else if strings.Contains(cam_a, "0.0.0.0") {
				url = "http://" + strings.Replace(cam_a, "0.0.0.0", "localhost", -1)
			} else {
				url = "http://" + cam_a
			}
			open.Start(url)
			err = server.ListenAndServe()
			if err != nil {
				log.Error(err.Error())
			}
		} else {
			quit := make(chan os.Signal)
			signal.Notify(quit, []os.Signal{os.Interrupt, syscall.SIGTERM}...)
			<-quit
			log.Warn("receive interrupt signal")
		}
	},
}

func init() {
	rootCmd.AddCommand(cameraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cameraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cameraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cameraCmd.Flags().StringVarP(&cam_d, "id", "d", "0", "[camera ID] eg. /dev/video0")
	cameraCmd.Flags().StringVarP(&cam_a, "address", "a", "0.0.0.0:8080", "[host:port] eg. 127.0.0.1:8080 ")
	cameraCmd.Flags().StringVarP(&cam_x, "xml", "x", "", "人脸识别 [classifier XML file]")
	cameraCmd.Flags().BoolVarP(&cam_motion, "motion", "m", false, "是否开启运动侦测")
	cameraCmd.Flags().BoolVarP(&cam_windows, "windows", "w", false, "是否开启app显示 ｜ 浏览器显示")
	cameraCmd.Flags().BoolVarP(&cam_web, "web", "B", false, "是否开启浏览器显示")
	cameraCmd.Flags().BoolVarP(&cam_pic, "pic", "p", false, "是否保存图片")
	cameraCmd.Flags().BoolVarP(&cam_video, "video", "v", false, "是否保存视频")
	cameraCmd.Flags().BoolVarP(&cam_detect, "detect", "D", false, "是否支持人脸识别")
	cameraCmd.Flags().StringVarP(&email_user, "email_user", "U", "", "邮件发送用户名")
	cameraCmd.Flags().StringVarP(&email_pwd, "email_pwd", "W", "", "邮件发送密码")
	cameraCmd.Flags().StringVarP(&email_sendto, "email_sendto", "T", "", "邮件接收地址")
	cameraCmd.Flags().StringVarP(&email_host, "email_host", "H", "smtp.163.com", "邮件发送用户名")
	cameraCmd.Flags().IntVarP(&email_port, "email_port", "P", 465, "邮件发送端口")
	cameraCmd.Flags().IntVarP(&cam_count, "count", "C", 100, "摄像头计数")
}

func mjpegCapture() {
	var window *gocv.Window
	if cam_windows {
		window = gocv.NewWindow("Cam Window")
		defer window.Close()
	}

	img := gocv.NewMat()
	defer img.Close()

	var writer *gocv.VideoWriter
	if cam_video {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Cannot read device %v\n", cam_d)
			return
		}

		writer, err = gocv.VideoWriterFile(videoName, "MJPG", 25.0, img.Cols(), img.Rows(), true)
		if err != nil {
			log.Error("VideoWriterFile error", "err", err)
			return
		}
		defer writer.Close()
	}

	var blue color.RGBA
	var classifier gocv.CascadeClassifier
	if cam_detect {
		// color for the rect when faces detected
		blue = color.RGBA{0, 0, 255, 0}

		// load classifier to recognize faces
		classifier = gocv.NewCascadeClassifier()
		defer classifier.Close()

		if !classifier.Load(targetPath) {
			log.Info("Error reading cascade file: %v\n", targetPath)
			return
		}
	}

	count := 0
	for {
		count++
		if ok := webcam.Read(&img); !ok {
			log.Info("Device closed: %v\n", cam_d)
			return
		}
		if img.Empty() {
			continue
		}

		var isPerson bool
		if cam_detect {
			// detect faces
			rects := classifier.DetectMultiScale(img)
			if len(rects) > 0 {
				log.Info("found faces", "Num", len(rects), "Count", count)
				isPerson = true
			}

			// draw a rectangle around each face on the original image,
			// along with text identifying as "Human"
			for _, r := range rects {
				gocv.Rectangle(&img, r, blue, 3)

				size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
				pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
				gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
			}
		}

		// 窗口模式无论是否侦测到人脸或者运动都显示
		if cam_windows {
			window.IMShow(img)
			if window.WaitKey(1) == 27 {
				break
			}
		} else if cam_web {
			// 未探测到人依然可以视频或者
			buf, _ := gocv.IMEncode(".jpg", img)
			stream.UpdateJPEG(buf.GetBytes())
			buf.Close()
		}

		// 前提条件： 必须侦测到人脸且有运动侦测启动的情况下才会执行以下操作
		// 1. 如果侦测到人脸且处于运动状态，则保存图片
		// 2. 如果侦测到人脸且处于运动状态，则保存图片且发送邮件
		// 3. 如果侦测到人脸且处于运动状态，则保存视频
		// 4. 如果侦测到人脸且处于运动状态，则保存视频且发送邮件
		// 5. 如果侦测到人脸且处于运动状态，则保存视频且保存图片
		// 6. 如果侦测到人脸且处于运动状态，则保存视频且保存图片且发送邮件
		if cam_pic && !cam_video {
			if count%cam_count == 0 {
				// 读取图片
				// src := gocv.IMRead("image.png", gocv.IMReadColor)
				// croppedMat := src.Region(image.Rect(0, 0, src.Cols(), src.Rows()/2))
				if cam_detect && isPerson {
					log.Debug("检测到人脸")
					rs := img.Clone()
					pic_path := fmt.Sprintf("./%d-%d.jpg", time.Now().UnixMicro(), count)
					gocv.IMWrite(pic_path, rs)
					sendEmail(pic_path, fmt.Sprintf("%s 检测到人脸", time.Now().String()), "", "")
					log.Debug("保存图片", "COUNT", count)
				} else if cam_detect && !isPerson {
					log.Debug("No person detected")
				} else {
					rs := img.Clone()
					gocv.IMWrite(fmt.Sprintf("./%d-%d.jpg", time.Now().UnixMicro(), count), rs)
					log.Debug("保存图片", "COUNT", count)
				}
			}
		} else if !cam_pic && cam_video {
			log.Debug("保存视频", "COUNT", count)
			if cam_detect && isPerson {
				if _, err := os.Stat(videoName); err != nil {
					if !os.IsExist(err) {
						writer, err = gocv.VideoWriterFile(videoName, "MJPG", 25.0, img.Cols(), img.Rows(), true)
						if err != nil {
							log.Error("VideoWriterFile error", "err", err)
							return
						}
					}
				}
				writer.Write(img)
				if count%cam_count == 0 {
					log.Info("发送视频邮件", "FILE", videoName)
					sendEmail("", fmt.Sprintf("%s 检测到人脸", time.Now().String()), "", videoName)
				}
			} else if cam_detect && !isPerson {
				log.Debug("No person detected")
			} else {
				writer.Write(img)
			}
		} else if cam_pic && cam_video {
			log.Debug("保存视频", "COUNT", count)
			if cam_detect && isPerson {
				if _, err := os.Stat(videoName); err != nil {
					if !os.IsExist(err) {
						writer, err = gocv.VideoWriterFile(videoName, "MJPG", 25.0, img.Cols(), img.Rows(), true)
						if err != nil {
							log.Error("VideoWriterFile error", "err", err)
							return
						}
					}
				}
				writer.Write(img)
				if count%cam_count == 0 {
					log.Info("发送视频邮件", "FILE", videoName)
					rs := img.Clone()
					pic_path := fmt.Sprintf("./%d-%d.jpg", time.Now().UnixMicro(), count)
					gocv.IMWrite(pic_path, rs)
					sendEmail(pic_path, fmt.Sprintf("%s 检测到人脸", time.Now().String()), "", videoName)
				}
			} else if cam_detect && !isPerson {
				log.Debug("No person detected")
			} else {
				writer.Write(img)
				rs := img.Clone()
				gocv.IMWrite(fmt.Sprintf("./%d-%d.jpg", time.Now().UnixMicro(), count), rs)
				log.Debug("保存图片", "COUNT", count)
			}
		}
	}
}

// 设置内联图片
func sendEmail(path, title, body, video string) error {
	if email_user == "" || email_pwd == "" || email_host == "" || email_sendto == "" {
		return errors.New("Email config error")
	}
	to := strings.Split(email_sendto, ",")
	if body == "" {
		body = `<html>
		<body>
			<p>TITLE:</p>
			<img src="cid:myImage">
		</body>
	</html>`
		body = strings.Replace(body, "TITLE", title, -1)
	}

	hostname, err := os.Hostname()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	title = fmt.Sprintf("%s: %s", hostname, title)
	if path == "" {
		return utils.NewEmail(email_user, email_pwd, email_host, email_port).Send(to, title, body, video)
	} else {
		return utils.NewEmail(email_user, email_pwd, email_host, email_port).Send(to, title, body, video, path)
	}
}

func motionCapture() {
	var window *gocv.Window
	if cam_windows {
		window = gocv.NewWindow("Motion Window")
		defer window.Close()
	}

	img := gocv.NewMat()
	defer img.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	status := MotionReady

	log.Info("Start reading device: %v\n", cam_d)

	var writer *gocv.VideoWriter
	if cam_video {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Cannot read device %v\n", cam_d)
			return
		}

		writer, err = gocv.VideoWriterFile(videoName, "MJPG", 20, img.Cols(), img.Rows(), true)
		if err != nil {
			log.Error("VideoWriterFile error", "err", err)
			return
		}
		defer writer.Close()
	}

	// 人脸识别
	var blue color.RGBA
	var classifier gocv.CascadeClassifier
	if cam_detect {
		// color for the rect when faces detected
		blue = color.RGBA{0, 0, 255, 0}

		// load classifier to recognize faces
		classifier = gocv.NewCascadeClassifier()
		defer classifier.Close()

		if !classifier.Load(targetPath) {
			log.Info("Error reading cascade file: %v\n", targetPath)
			return
		}
	}

	count := 0

	for {
		if ok := webcam.Read(&img); !ok {
			log.Info("Device closed: %v\n", cam_d)
			return
		}
		if img.Empty() {
			continue
		}

		status = MotionReady
		statusColor := color.RGBA{0, 255, 0, 0}

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(img, &imgDelta)

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		kernel.Close()

		// now find contours
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)

		for i := 0; i < contours.Size(); i++ {
			area := gocv.ContourArea(contours.At(i))
			if area < MinimumArea {
				continue
			}

			status = MotionDetect
			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(contours.At(i))
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}

		contours.Close()

		gocv.PutText(&img, string(status), image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		var isPerson bool
		if cam_detect {
			// detect faces
			rects := classifier.DetectMultiScale(img)
			log.Debug("found faces", "Num", len(rects), "Count", count)
			if len(rects) > 0 {
				isPerson = true
			}

			// draw a rectangle around each face on the original image,
			// along with text identifying as "Human"
			for _, r := range rects {
				gocv.Rectangle(&img, r, blue, 3)

				size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
				pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
				gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
			}
		}

		if cam_windows {
			window.IMShow(img)
			if window.WaitKey(1) == 27 {
				break
			}
		} else if cam_web {
			buf, _ := gocv.IMEncode(".jpg", img)
			stream.UpdateJPEG(buf.GetBytes())
			buf.Close()
		}

		count++
		// save image to file
		// https://blog.csdn.net/m0_55708805/article/details/115467324
		if cam_pic && !cam_video {
			if count%cam_count == 0 {
				// 读取图片
				// src := gocv.IMRead("image.png", gocv.IMReadColor)
				// croppedMat := src.Region(image.Rect(0, 0, src.Cols(), src.Rows()/2))
				if cam_detect && isPerson && status == MotionDetect {
					log.Debug("保存图片", "COUNT", count)
					rs := img.Clone()
					pic_path := fmt.Sprintf("./%d-%d.jpg", time.Now().UnixMicro(), count)
					gocv.IMWrite(pic_path, rs)
					sendEmail(pic_path, fmt.Sprintf("%s 检测到人脸", time.Now().String()), "", "")
				} else {
					log.Debug("No person detected")
				}
			}
		} else if !cam_pic && cam_video {
			log.Debug("保存视频", "COUNT", count)
			if cam_detect && isPerson && status == MotionDetect {
				if _, err := os.Stat(videoName); err != nil {
					if !os.IsExist(err) {
						writer, err = gocv.VideoWriterFile(videoName, "MJPG", 25.0, img.Cols(), img.Rows(), true)
						if err != nil {
							log.Error("VideoWriterFile error", "err", err)
							return
						}
					}
				}
				writer.Write(img)
				if count%cam_count == 0 {
					log.Info("发送视频邮件", "FILE", videoName)
					sendEmail("", fmt.Sprintf("%s 检测到人脸", time.Now().String()), "", videoName)
				}
			} else {
				log.Debug("No person detected")
			}
		} else if cam_pic && cam_video {
			// 读取图片
			// src := gocv.IMRead("image.png", gocv.IMReadColor)
			// croppedMat := src.Region(image.Rect(0, 0, src.Cols(), src.Rows()/2))
			if cam_detect && isPerson && status == MotionDetect {
				if _, err := os.Stat(videoName); err != nil {
					if !os.IsExist(err) {
						writer, err = gocv.VideoWriterFile(videoName, "MJPG", 25.0, img.Cols(), img.Rows(), true)
						if err != nil {
							log.Error("VideoWriterFile error", "err", err)
							return
						}
					}
				}
				writer.Write(img)
				if count%cam_count == 0 {
					rs := img.Clone()
					pic_path := fmt.Sprintf("./%d-%d.jpg", time.Now().UnixMicro(), count)
					gocv.IMWrite(pic_path, rs)
					sendEmail(pic_path, fmt.Sprintf("%s 检测到人脸", time.Now().String()), "", videoName)
				}

			} else {
				log.Debug("No person detected")
			}
			log.Debug("保存图片", "COUNT", count)
		}
	}
}
