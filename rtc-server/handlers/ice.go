package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleIce(ws *websocket.Conn, msg *types.Ice) error {
	room, ok := rooms[msg.Room]
	if !ok {
		errMsg := "Room not found"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	repacked, err := types.NewMessageIce(msg)
	if err != nil {
		errMsg := "Can't create a new message ice"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	if ws == room.connectee.conn {
		return room.connector.conn.WriteJSON(repacked)
	}
	return room.connectee.conn.WriteJSON(repacked)
}
