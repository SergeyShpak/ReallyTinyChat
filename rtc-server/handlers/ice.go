package handlers

import (
	"fmt"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleIce(ws *websocket.Conn, msg *types.Ice) error {
	partnerConn, err := getConnectionInRoom(msg.Room, msg.Partner)
	if err != nil {
		return err
	}
	repacked, err := types.NewMessageIce(msg)
	if err != nil {
		return errors.NewServerError(500, "cannot forward the ICE message")
	}
	if err := partnerConn.WS.WriteJSON(repacked); err != nil {
		return errors.NewServerError(500, fmt.Sprintf("error occurred when sending an ICE message: %v", err))
	}
	return nil
}
