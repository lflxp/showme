package proxy

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var logger = zerolog.Logger{}

const bind = "127.0.0.1:6001"

func init() {
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("***%s****", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	logger = zerolog.New(output).With().Timestamp().Logger()
}

// 4层 TCP反向代理
// from: 127.0.0.1:9999 to: 10.111.1.1:9999
func NewTCPProxy(from, to string) error {
	logger.Level(zerolog.DebugLevel)

	if to == "" {
		return errors.New("proxy backend is empty")
	}

	if from == "" {
		logger.Info().Str("bind", bind).Msg("use default bind")
		from = bind
	}

	err := runProxy(from, to)
	if err != nil {
		logger.Error().Err(err).Send()
		return err
	}

	return nil
}

func runProxy(bind, backend string) error {
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	defer listener.Close()
	logger.Info().Str("bind", bind).Str("backend", backend).Msg("tcp-proxy started.")
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error().Err(err).Send()
		} else {
			go connectionHandler(conn, backend)
		}
	}
}

// 拨号链接
func connectionHandler(conn net.Conn, backend string) {
	logger.Info().Str("conn", conn.RemoteAddr().String()).Msg("client connected.")
	target, err := net.Dial("tcp", backend)
	defer conn.Close()
	if err != nil {
		logger.Error().Err(err).Send()
	} else {
		defer target.Close()
		logger.Info().Str("conn", conn.RemoteAddr().String()).Str("backend", target.LocalAddr().String()).Msg("backend connected.")
		closed := make(chan bool, 2)
		go proxy(conn, target, closed)
		go proxy(target, conn, closed)
		<-closed
		logger.Info().Str("conn", conn.RemoteAddr().String()).Msg("Connection closed.")
	}
}

// 反向代理
func proxy(from net.Conn, to net.Conn, closed chan bool) {
	buffer := make([]byte, 4096)
	for {
		n1, err := from.Read(buffer)
		if err != nil {
			closed <- true
			return
		}
		n2, err := to.Write(buffer[:n1])
		logger.Debug().Str("from", from.RemoteAddr().String()).Int("recv", n1).Str("to", to.RemoteAddr().String()).Int("send", n2).Msg("Proxyying")
		if err != nil {
			closed <- true
			return
		}
	}
}
