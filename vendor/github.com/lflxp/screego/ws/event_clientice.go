package ws

import (
	"fmt"

	"github.com/lflxp/screego/ws/outgoing"
	"github.com/rs/zerolog/log"
)

func init() {
	register("clientice", func() Event {
		return &ClientICE{}
	})
}

type ClientICE outgoing.P2PMessage

func (e *ClientICE) Execute(rooms *Rooms, current ClientInfo) error {
	if current.RoomID == "" {
		return fmt.Errorf("not in a room")
	}

	room, ok := rooms.Rooms[current.RoomID]
	if !ok {
		return fmt.Errorf("room with id %s does not exist", current.RoomID)
	}

	session, ok := room.Sessions[e.SID]

	if !ok {
		log.Debug().Str("id", e.SID.String()).Msg("unknown session")
		return nil
	}

	if session.Client != current.ID {
		return fmt.Errorf("permission denied for session %s", e.SID)
	}

	room.Users[session.Host].Write <- outgoing.ClientICE(*e)

	return nil
}
