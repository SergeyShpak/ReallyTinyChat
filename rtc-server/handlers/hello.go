package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleHello(ws *websocket.Conn, msg *types.Hello) {
	if err := addToConnections(ws, msg); err != nil {
		errMsg := &types.Error{
			Hint: err.Error(),
			Code: http.StatusConflict,
		}
		ws.WriteJSON(errMsg)
		return
	}
	return
}

func addToConnections(ws *websocket.Conn, msg *types.Hello) error {
	room, ok := rooms[msg.Room]
	if !ok {
		createRoom(ws, msg)
		return sendHelloOKMessage(ws, msg)
	}
	if err := enterRoom(ws, msg, room); err != nil {
		return err
	}
	return sendRoomInfoMessage(room)
}

func createRoom(ws *websocket.Conn, msg *types.Hello) {
	rooms[msg.Room] = &room{
		name: msg.Room,
		connector: &connection{
			login: msg.Login,
			conn:  ws,
		},
	}
	log.Printf("%s created a room %s\n", msg.Login, msg.Room)
}

func enterRoom(ws *websocket.Conn, msg *types.Hello, r *room) error {
	if r.connectee != nil {
		return fmt.Errorf("conflict")
	}
	r.connectee = &connection{
		login: msg.Login,
		conn:  ws,
	}
	log.Printf("%s entered a room with %s\n", r.connectee.login, r.connector.login)
	return nil
}

func sendHelloOKMessage(ws *websocket.Conn, msg *types.Hello) error {
	okMsg, err := types.NewMessageHelloOK(msg.Login, msg.Room)
	if err != nil {
		log.Println("Could not create a HelloOK message")
		return err
	}
	ws.WriteJSON(okMsg)
	return nil
}

func sendRoomInfoMessage(r *room) error {
	log.Println("Sending room info messages")
	roomInfoMsg, err := types.NewMessageRoomInfo(r.connector.login, r.connectee.login, r.name)
	if err != nil {
		log.Println("Could not create a RoomInfo message")
		return err
	}
	r.connector.conn.WriteJSON(roomInfoMsg)
	r.connectee.conn.WriteJSON(roomInfoMsg)
	return nil
}
