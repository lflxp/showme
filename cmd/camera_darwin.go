//go:build darwin
// +build darwin

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

package cmd

import (
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	log "log/slog"

	"github.com/hybridgroup/mjpeg"
	"github.com/skratchdot/open-golang/open"
	"gocv.io/x/gocv"

	"github.com/spf13/cobra"
)

const MinimumArea = 3000

var (
	cam_d       string
	cam_a       string
	cam_x       string
	cam_motion  bool
	cam_windows bool
	err         error
	webcam      *gocv.VideoCapture
	stream      *mjpeg.Stream
)

// cameraCmd represents the camera command
var cameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "本地摄像头webcam",
	Long:  `浏览器webcam`,
	Run: func(cmd *cobra.Command, args []string) {
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

		if cam_windows {
			quit := make(chan os.Signal)
			signal.Notify(quit, []os.Signal{os.Interrupt, syscall.SIGTERM}...)
			<-quit
			log.Warn("receive interrupt signal")
		} else {
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
}

func mjpegCapture() {
	var window *gocv.Window
	if cam_windows {
		window = gocv.NewWindow("Cam Window")
		defer window.Close()
	}

	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			log.Info("Device closed: %v\n", cam_d)
			return
		}
		if img.Empty() {
			continue
		}

		if cam_windows {
			window.IMShow(img)
			if window.WaitKey(1) == 27 {
				break
			}
		} else {
			buf, _ := gocv.IMEncode(".jpg", img)
			stream.UpdateJPEG(buf.GetBytes())
			buf.Close()
		}
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

	status := "Ready"

	log.Info("Start reading device: %v\n", cam_d)

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", cam_d)
			return
		}
		if img.Empty() {
			continue
		}

		status = "Ready"
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

			status = "Motion detected"
			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(contours.At(i))
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}

		contours.Close()

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		if cam_windows {
			window.IMShow(img)
			if window.WaitKey(1) == 27 {
				break
			}
		} else {
			buf, _ := gocv.IMEncode(".jpg", img)
			stream.UpdateJPEG(buf.GetBytes())
			buf.Close()
		}
	}
}
