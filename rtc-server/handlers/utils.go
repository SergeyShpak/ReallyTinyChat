package handlers

import (
	"github.com/gorilla/websocket"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

func getRoom(name string) *types.Room {
	rIface, ok := rooms.Load(name)
	if !ok {
		return nil
	}
	// TODO(SSH): throw error
	r, ok := rIface.(*types.Room)
	if !ok {
		return nil
	}
	return r
}

func addRoom(room *types.Room) {
	rooms.Store(room.Name, room)
}

func getUser(ws *websocket.Conn) *types.User {
	userIface, ok := users.Load(ws)
	if !ok {
		return nil
	}
	user, ok := userIface.(*types.User)
	// TODO(SSH): throw error
	if !ok {
		return nil
	}
	return user
}
