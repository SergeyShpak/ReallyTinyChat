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
			fmt.Println("error occurred: ", err)
			return
		}
		ws.WriteJSON(errMsg)
		log.Println(servErr)
		return
	}
	return
}

func addToConnections(ws *websocket.Conn, msg *types.Hello) error {
	r, ok := rooms.Load(msg.Room)
	if !ok {
		createRoom(ws, msg)
		return sendHelloOKMessage(ws, msg)
	}
	if err := enterRoom(ws, msg, r.(*room)); err != nil {
		return err
	}
	return sendRoomInfoMessage(r.(*room))
}

func createRoom(ws *websocket.Conn, msg *types.Hello) {
	r := &room{
		name: msg.Room,
		connector: &connection{
			login: msg.Login,
			conn:  ws,
		},
	}
	rooms.Store(msg.Room, r)
	wsRooms.Store(ws, r)
	log.Printf("%s created a room %s\n", msg.Login, msg.Room)
}

func enterRoom(ws *websocket.Conn, msg *types.Hello, r *room) error {
	wsRooms.Store(ws, r)
	if r.connectee != nil {
		return fmt.Errorf("cannot join the room %s as it is already full", r.name)
	}
	if r.connector.login == msg.Login {
		return fmt.Errorf("change your login \"%s\" as it is the same as that of the room owner", msg.Login)
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
