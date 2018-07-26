package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleOffer(ws *websocket.Conn, msg *types.Offer) error {
	rInterface, ok := rooms.Load(msg.Room)
	if !ok {
		errMsg := "Room not found"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	r := rInterface.(*room)
	repacked, err := types.NewMessageOffer(msg)
	if err != nil {
		errMsg := "Can't create a new message offer"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	if msg.IsResponse {
		if r.connectee == nil {
			errMsg := "Room is not full, no connectee"
			log.Println(errMsg)
			return fmt.Errorf(errMsg)
		}
		log.Println("Sending to: ", r.connectee.login)
		return r.connectee.conn.WriteJSON(repacked)
	}
	if r.connector == nil {
		errMsg := "Room is not full, no connector"
		log.Println(errMsg)
		return fmt.Errorf(errMsg)
	}
	log.Println("Sending to: ", r.connector.login)
	return r.connector.conn.WriteJSON(repacked)
}
