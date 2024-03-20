//go:build linux
// +build linux

/*
Copyright © 2024 NAME lflxp <382023823@qq.com>
*/
package cmd

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	log "log/slog"

	"github.com/korandiz/v4l"
	"github.com/korandiz/v4l/fmt/mjpeg"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

var (
	cam_d string
	cam_w int
	cam_h int
	cam_f int
	cam_a string
	cam_l bool
	cam_r bool
)

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
		fmt.Fprintln(os.Stderr, "Using device", cam_d)
		cam, err := v4l.Open(cam_d)
		fatal("Open", err)

		if cam_l {
			configs, err := cam.ListConfigs()
			fatal("ListConfigs", err)
			fmt.Fprintln(os.Stderr, "Supported device configs:")
			found := false
			for _, cfg := range configs {
				if cfg.Format != mjpeg.FourCC {
					continue
				}
				found = true
				fmt.Fprintln(os.Stderr, " ", cfg2str(cfg))
			}
			if !found {
				fmt.Fprintln(os.Stderr, "  (none)")
			}
			os.Exit(0)
		}

		cfg, err := cam.GetConfig()
		fatal("GetConfig", err)
		cfg.Format = mjpeg.FourCC
		if cam_w > 0 {
			cfg.Width = cam_w
		}
		if cam_h > 0 {
			cfg.Height = cam_h
		}
		if cam_f > 0 {
			cfg.FPS = v4l.Frac{uint32(cam_f), 1}
		}
		fmt.Fprintln(os.Stderr, "Requested config:", cfg2str(cfg))
		err = cam.SetConfig(cfg)
		fatal("SetConfig", err)
		err = cam.TurnOn()
		fatal("TurnOn", err)
		cfg, err = cam.GetConfig()
		fatal("GetConfig", err)
		if cfg.Format != mjpeg.FourCC {
			fmt.Fprintln(os.Stderr, "Failed to set MJPEG format.")
			os.Exit(1)
		}
		fmt.Fprintln(os.Stderr, "Actual device config:", cfg2str(cfg))

		if cam_r {
			ctrls, err := cam.ListControls()
			fatal("ListControls", err)
			for _, ctrl := range ctrls {
				cam.SetControl(ctrl.CID, ctrl.Default)
			}
		}

		blankImg := image.NewRGBA(image.Rect(0, 0, cfg.Width, cfg.Height))
		buf := new(bytes.Buffer)
		jpeg.Encode(buf, blankImg, nil)
		blank = buf.Bytes()

		go handleInterrupt()
		go stream(cam)

		openUrl := fmt.Sprintf("http://127.0.0.1%s", cam_a)
		open.Start(openUrl)
		log.Info("Listening on address", cam_a)
		srv := http.Server{
			Addr:    cam_a,
			Handler: http.HandlerFunc(serveHTTP),
		}
		err = srv.ListenAndServe()
		fatal("ListenAndServe", err)
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
	cameraCmd.Flags().IntVarP(&cam_w, "width", "w", 0, "image width")
	cameraCmd.Flags().IntVarP(&cam_h, "height", "H", 0, "image height")
	cameraCmd.Flags().IntVarP(&cam_f, "frame", "f", 0, "frame rate")
	cameraCmd.Flags().StringVarP(&cam_a, "address", "a", ":8080", "address to listen on")
	cameraCmd.Flags().BoolVarP(&cam_l, "log", "l", false, "print supported device configs and quit")
	cameraCmd.Flags().BoolVarP(&cam_r, "reset", "r", false, "reset all controls to default")
}

func cfg2str(cfg v4l.DeviceConfig) string {
	return fmt.Sprintf("%dx%d @ %.4g FPS", cfg.Width, cfg.Height,
		float64(cfg.FPS.N)/float64(cfg.FPS.D))
}

var (
	mu      sync.Mutex
	clients []*client
	stopped bool
	quit    = make(chan int, 1)
	blank   []byte
)

type client struct {
	i  int
	ch chan []byte
}

func newClient() *client {
	mu.Lock()
	defer mu.Unlock()
	if stopped {
		return nil
	}
	clt := &client{
		i:  len(clients),
		ch: make(chan []byte, 1),
	}
	clt.ch <- blank
	clients = append(clients, clt)
	return clt
}

func (clt *client) remove() {
	mu.Lock()
	defer mu.Unlock()
	i := clt.i
	last := len(clients) - 1
	clients[i] = clients[last]
	clients[i].i = i
	clients[last] = nil
	clients = clients[:last]
	clt.i = -1
	if stopped && len(clients) == 0 {
		quit <- 1
	}
}

func handleInterrupt() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Info("Stopping...")
	mu.Lock()
	stopped = true
	if len(clients) == 0 {
		quit <- 1
	} else {
		for _, clt := range clients {
			close(clt.ch)
		}
	}
	mu.Unlock()
	<-quit
	os.Exit(0)
}

func stream(cam *v4l.Device) {
	for {
		buf, err := cam.Capture()
		if err != nil {
			log.Info("Capture:", err)
			proc, _ := os.FindProcess(os.Getpid())
			proc.Signal(os.Interrupt)
			// os.Exit(1)
			break
		}
		b := make([]byte, buf.Size())
		buf.ReadAt(b, 0)
		mu.Lock()
		if stopped {
			mu.Unlock()
			break
		}
		for _, clt := range clients {
			select {
			case clt.ch <- b:
			case <-clt.ch:
				clt.ch <- b
			}
		}
		mu.Unlock()
	}
}

func serveHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info("[%s] New connection\n", r.RemoteAddr)
	hj, ok := w.(http.Hijacker)
	if !ok {
		log.Info("[%s] ResponseWrite is not a Hijacker", r.RemoteAddr)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	conn, _, err := hj.Hijack()
	if err != nil {
		log.Info("[%s] Hijack: %v\n", r.RemoteAddr, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	defer conn.Close()

	clt := newClient()
	if clt == nil {
		return
	}
	defer clt.remove()

	const B = "45c7pIy0cxa4vWtwGuVuAkbzKAQGpRjz9eyhyHTv"

	_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\n" +
		"Date: " + time.Now().UTC().Format(http.TimeFormat) + "\r\n" +
		"Content-Type: multipart/x-mixed-replace; boundary=" + B + "\r\n" +
		"Cache-Control: no-cache, no-store, max-age=0, must-revalidate\r\n" +
		"Pragma: no-cache\r\n" +
		"\r\n" +
		"--" + B + "\r\n"))
	if err != nil {
		log.Info("[%s] %v\n", r.RemoteAddr, err)
		return
	}

	for {
		buf := <-clt.ch
		conn.SetDeadline(time.Now().Add(time.Second))
		_, err := conn.Write([]byte("Content-Type: image/jpeg\r\n\r\n"))
		if err != nil {
			log.Info("[%s] %v\n", r.RemoteAddr, err)
			return
		}
		if buf == nil {
			_, err := conn.Write(blank)
			if err == nil {
				conn.Write([]byte("--" + B + "--\r\n"))
			} else {
				log.Info("[%s] %v\n", r.RemoteAddr, err)
			}
			log.Info("[%s] Quitting\n", r.RemoteAddr)
			return
		}
		_, err = conn.Write(buf)
		if err != nil {
			log.Info("[%s] %v\n", r.RemoteAddr, err)
			return
		}
		_, err = conn.Write([]byte("--" + B + "\r\n"))
		if err != nil {
			log.Info("[%s] %v\n", r.RemoteAddr, err)
			return
		}
	}
}

func fatal(p string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", p, err)
		os.Exit(1)
	}
}
