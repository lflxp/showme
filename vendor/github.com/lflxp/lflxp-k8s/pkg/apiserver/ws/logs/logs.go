package logs

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/go-eden/slf4go"
	"github.com/gorilla/websocket"
	"github.com/lflxp/lflxp-k8s/pkg/apiserver/model"
)

var (
	// homeTempl = template.Must(template.New("").Parse(homeHTML))
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func GetLogs(data *model.CoreV1, c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Error(err.Error())
		}
		return
	}

	go writer(ws, data)
	reader(ws)
}

func writer(ws *websocket.Conn, data *model.CoreV1) {
	defer ws.Close()

	stream, err := data.GetLogs()
	if err != nil {
		log.Error(err.Error())
		return
	}
	defer stream.Close()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	// https://stackoverflow.com/questions/53852530/how-to-get-logs-from-kubernetes-using-go
	for {
		select {
		case <-pingTicker.C:
			return
		default:
			buf := make([]byte, upgrader.ReadBufferSize)
			numBytes, err := stream.Read(buf)
			if numBytes == 0 {
				continue
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}
			// message := string(buf[:numBytes])
			// log.Debug(message)
			if err := ws.WriteMessage(websocket.TextMessage, buf[:numBytes]); err != nil {
				log.Error(err.Error())
				return
			}
		}
	}
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(int64(upgrader.ReadBufferSize))
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}
