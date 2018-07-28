package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleHello(ws *websocket.Conn, msg *types.Hello) {
	if err := addToConnections(ws, msg); err != nil {
		servErr, ok := err.(*errors.ServerError)
		if !ok {
			servErr = errors.NewServerError(500, fmt.Sprintf("%s", err))
		}
		errMsg, err := types.NewMessageError(servErr)
		if err != nil {
			log.Println("error occurred: ", err)
			return
		}
		ws.WriteJSON(errMsg)
		log.Println(servErr)
		return
	}
	if err := sendHelloOKMessage(ws, msg); err != nil {
		log.Println("error occurred: ", err)
		return
	}
	log.Println("HELLOOK message sent")
	return
}

func addToConnections(ws *websocket.Conn, msg *types.Hello) error {
	r := getRoom(msg.Room)
	if r == nil {
		r = types.NewRoom(msg.Room)
	}
	if r.IsConnected(msg.Login) {
		return errors.NewServerError(409, fmt.Sprintf("user %s is already connected", msg.Login))
	}
	conn := &types.Connection{
		Login: msg.Login,
		WS:    ws,
	}
	r.AddConnection(conn)
	return nil
}

func sendHelloOKMessage(ws *websocket.Conn, msg *types.Hello) error {
	okMsg, err := types.NewMessageHelloOK(msg.Login, msg.Room)
	if err != nil {
		return err
	}
	ws.WriteJSON(okMsg)
	return nil
}
