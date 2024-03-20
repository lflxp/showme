//go:build linux
// +build linux

/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	log "log/slog"

	"github.com/blackjack/webcam"
	"github.com/korandiz/v4l"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var (
	cam_d string
	cam_f string
	cam_s string
	cam_m bool
	cam_l string
	cam_p bool
)

const (
	V4L2_PIX_FMT_PJPG = 0x47504A50
	V4L2_PIX_FMT_YUYV = 0x56595559
)

type FrameSizes []webcam.FrameSize

func (slice FrameSizes) Len() int {
	return len(slice)
}

// For sorting purposes
func (slice FrameSizes) Less(i, j int) bool {
	ls := slice[i].MaxWidth * slice[i].MaxHeight
	rs := slice[j].MaxWidth * slice[j].MaxHeight
	return ls < rs
}

// For sorting purposes
func (slice FrameSizes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

var supportedFormats = map[webcam.PixelFormat]bool{
	V4L2_PIX_FMT_PJPG: true,
	V4L2_PIX_FMT_YUYV: true,
}

// cameraCmd represents the camera command
var cameraCmd = &cobra.Command{
	Use:   "camera",
	Short: "本地摄像头webcam",
	Long:  `浏览器webcam`,
	Run: func(cmd *cobra.Command, args []string) {
		if cam_d == "" {
			devs := v4l.FindDevices()
			if len(devs) != 1 {
				fmt.Fprintln(os.Stderr, "Use -d to select device.")
				for _, info := range devs {
					fmt.Fprintln(os.Stderr, " ", info.Path)
				}
				os.Exit(1)
			}
			cam_d = devs[0].Path
		}
		cam, err := webcam.Open(cam_d)
		if err != nil {
			panic(err.Error())
		}
		defer cam.Close()

		// select pixel format
		format_desc := cam.GetSupportedFormats()

		fmt.Println("Available formats:")
		for _, s := range format_desc {
			fmt.Fprintln(os.Stderr, s)
		}

		var format webcam.PixelFormat
	FMT:
		for f, s := range format_desc {
			if cam_f == "" {
				if supportedFormats[f] {
					format = f
					break FMT
				}

			} else if cam_f == s {
				if !supportedFormats[f] {
					log.Info(format_desc[f], "format is not supported, exiting")
					return
				}
				format = f
				break
			}
		}
		if format == 0 {
			log.Info("No format found, exiting")
			return
		}

		// select frame size
		frames := FrameSizes(cam.GetSupportedFrameSizes(format))
		sort.Sort(frames)

		fmt.Fprintln(os.Stderr, "Supported frame sizes for format", format_desc[format])
		for _, f := range frames {
			fmt.Fprintln(os.Stderr, f.GetString())
		}
		var size *webcam.FrameSize
		if cam_s == "" {
			size = &frames[len(frames)-1]
		} else {
			for _, f := range frames {
				if cam_s == f.GetString() {
					size = &f
					break
				}
			}
		}
		if size == nil {
			log.Info("No matching frame size, exiting")
			return
		}

		fmt.Fprintln(os.Stderr, "Requesting", format_desc[format], size.GetString())
		f, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
		if err != nil {
			log.Info("SetImageFormat return error", err)
			return

		}
		fmt.Fprintf(os.Stderr, "Resulting image format: %s %dx%d\n", format_desc[f], w, h)

		// start streaming
		err = cam.StartStreaming()
		if err != nil {
			log.Error(err.Error())
			return
		}

		var (
			li   chan *bytes.Buffer = make(chan *bytes.Buffer)
			fi   chan []byte        = make(chan []byte)
			back chan struct{}      = make(chan struct{})
		)
		go encodeToImage(cam, back, fi, li, w, h, f)
		if cam_m {
			go httpImage(cam_l, li)
		} else {
			go httpVideo(cam_l, li)
		}

		timeout := uint32(5) //5 seconds
		start := time.Now()
		var fr time.Duration

		for {
			err = cam.WaitForFrame(timeout)
			if err != nil {
				log.Error(err.Error())
				return
			}

			switch err.(type) {
			case nil:
			case *webcam.Timeout:
				log.Error(err.Error())
				continue
			default:
				log.Error(err.Error())
				return
			}

			frame, err := cam.ReadFrame()
			if err != nil {
				log.Error(err.Error())
				return
			}
			if len(frame) != 0 {

				// print framerate info every 10 seconds
				fr++
				if cam_p {
					if d := time.Since(start); d > time.Second*10 {
						fmt.Println(float64(fr)/(float64(d)/float64(time.Second)), "fps")
						start = time.Now()
						fr = 0
					}
				}

				select {
				case fi <- frame:
					<-back
				default:
				}
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
	cameraCmd.Flags().StringVarP(&cam_d, "device", "d", "", "摄像头设备 eg. /dev/video0")
	cameraCmd.Flags().StringVarP(&cam_f, "fmtstr", "f", "", "video format to use, default first supported")
	cameraCmd.Flags().StringVarP(&cam_s, "szstr", "s", "", "frame size to use, default largest one")
	cameraCmd.Flags().StringVarP(&cam_l, "address", "l", ":8080", "addr to listien")
	cameraCmd.Flags().BoolVarP(&cam_p, "fps", "p", false, "print fps info")
	cameraCmd.Flags().BoolVarP(&cam_m, "single", "m", false, "single image http mode, default mjpeg video")
}

func encodeToImage(wc *webcam.Webcam, back chan struct{}, fi chan []byte, li chan *bytes.Buffer, w, h uint32, format webcam.PixelFormat) {

	var (
		frame []byte
		img   image.Image
	)
	for {
		bframe := <-fi
		// copy frame
		if len(frame) < len(bframe) {
			frame = make([]byte, len(bframe))
		}
		copy(frame, bframe)
		back <- struct{}{}

		switch format {
		case V4L2_PIX_FMT_YUYV:
			yuyv := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
			for i := range yuyv.Cb {
				ii := i * 4
				yuyv.Y[i*2] = frame[ii]
				yuyv.Y[i*2+1] = frame[ii+2]
				yuyv.Cb[i] = frame[ii+1]
				yuyv.Cr[i] = frame[ii+3]

			}
			img = yuyv
		default:
			log.Error("invalid format ?")
		}
		//convert to jpeg
		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, img, nil); err != nil {
			log.Error(err.Error())
			return
		}

		const N = 50
		// broadcast image up to N ready clients
		nn := 0
	FOR:
		for ; nn < N; nn++ {
			select {
			case li <- buf:
			default:
				break FOR
			}
		}
		if nn == 0 {
			li <- buf
		}

	}
}

func httpImage(addr string, li chan *bytes.Buffer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("connect from", r.RemoteAddr, r.URL)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		//remove stale image
		<-li

		img := <-li

		w.Header().Set("Content-Type", "image/jpeg")

		if _, err := w.Write(img.Bytes()); err != nil {
			log.Error(err.Error())
			return
		}

	})

	var url string
	if strings.HasPrefix(addr, ":") {
		url = "http://localhost" + addr
	} else {
		url = "http://" + addr
	}
	open.Start(url)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Error(err.Error())
	}
}

func httpVideo(addr string, li chan *bytes.Buffer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info("connect from", r.RemoteAddr, r.URL)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		//remove stale image
		<-li
		const boundary = `frame`
		w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
		multipartWriter := multipart.NewWriter(w)
		multipartWriter.SetBoundary(boundary)
		for {
			img := <-li
			image := img.Bytes()
			iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
				"Content-type":   []string{"image/jpeg"},
				"Content-length": []string{strconv.Itoa(len(image))},
			})
			if err != nil {
				log.Error(err.Error())
				return
			}
			_, err = iw.Write(image)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
	})

	var url string
	if strings.HasPrefix(addr, ":") {
		url = "http://localhost" + addr
	} else {
		url = "http://" + addr
	}
	open.Start(url)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Error(err.Error())
	}
}
