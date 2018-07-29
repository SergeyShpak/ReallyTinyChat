package handlers

import (
	"fmt"

	"github.com/gorilla/websocket"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

func getRoom(name string) (*types.Room, error) {
	rIface, ok := rooms.Load(name)
	if !ok {
		return nil, errors.NewServerError(404, fmt.Sprintf("room \"%s\" was not found", name))
	}
	r, ok := rIface.(*types.Room)
	if !ok {
		return nil, errors.NewServerError(500, fmt.Sprintf("could not cast the room \"%s\" to its original type", name))
	}
	return r, nil
}

func addRoom(room *types.Room) {
	rooms.Store(room.Name, room)
}

func getUser(ws *websocket.Conn) (*types.User, error) {
	userIface, ok := users.Load(ws)
	if !ok {
		return nil, errors.NewServerError(404, fmt.Sprintf("a user that corresponds to the given WebSocket connection was not found"))
	}
	user, ok := userIface.(*types.User)
	if !ok {
		return nil, errors.NewServerError(500, fmt.Sprintf("could not cast the object \"%v\" to the User type", userIface))
	}
	return user, nil
}
