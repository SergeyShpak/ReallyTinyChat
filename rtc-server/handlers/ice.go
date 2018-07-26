package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleIce(ws *websocket.Conn, msg *types.Ice) error {
	rInterface, ok := rooms.Load(msg.Room)
	if !ok {
		errMsg := "Room not found"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	r := rInterface.(*room)
	repacked, err := types.NewMessageIce(msg)
	if err != nil {
		errMsg := "Can't create a new message ice"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	if ws == r.connectee.conn {
		return r.connector.conn.WriteJSON(repacked)
	}
	return r.connectee.conn.WriteJSON(repacked)
}
