package tty

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"golang.org/x/net/websocket"
)

func EncodeBase64(in string) string {
	return base64.StdEncoding.EncodeToString([]byte(in))
}

func Test_client(t *testing.T) {
	var wsurl = "ws://127.0.0.1:8080/ws"
	var origin = "http://127.0.0.1:8080/"
	ws, err := websocket.Dial(wsurl, "", origin)
	if err != nil {
		t.Fatal(err)
	}

	count := 0
	go func() {
		for {
			count++
			data := make([]byte, 1024)
			n, err := ws.Read(data)
			if err != nil {
				t.Fatal("receive", err)
				break
			}
			t.Log(string(data[:n]))

			if count > 100 {
				break
			}
		}
	}()

	var input string

	for {
		count++
		fmt.Scanln(&input)
		t.Logf("input %s\n", strings.TrimSpace(input))
		_, err := ws.Write([]byte(fmt.Sprintf("0pauDwXsnMEUerqd6psMKC1gWpSq3JLCJ%s", EncodeBase64(strings.TrimSpace(input)))))
		if err != nil {
			t.Fatal("send", err.Error())
		}

		if count > 100 {
			break
		}
	}
}
