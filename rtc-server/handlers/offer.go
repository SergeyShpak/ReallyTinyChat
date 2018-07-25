package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleOffer(ws *websocket.Conn, msg *types.Offer) error {
	room, ok := rooms[msg.Room]
	if !ok {
		errMsg := "Room not found"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}

	repacked, err := types.NewMessageOffer(msg)
	if err != nil {
		errMsg := "Can't create a new message offer"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	if msg.IsResponse {
		if room.connectee == nil {
			errMsg := "Room is not full, no connectee"
			log.Println(errMsg)
			return fmt.Errorf(errMsg)
		}
		log.Println("Sending to: ", room.connectee.login)
		return room.connectee.conn.WriteJSON(repacked)
	}
	if room.connector == nil {
		errMsg := "Room is not full, no connector"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	log.Println("Sending to: ", room.connector.login)
	return room.connector.conn.WriteJSON(repacked)
}
