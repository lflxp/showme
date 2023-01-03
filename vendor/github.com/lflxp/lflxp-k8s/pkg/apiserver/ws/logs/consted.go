package logs

import "time"

const (
	// Time allowed to read the next pong message from the client.
	pongWait = 2 * 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// pingPeriod = 10 * time.Second
)
