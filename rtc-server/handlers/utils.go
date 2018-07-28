package handlers

import (
	"fmt"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
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

func getConnectionInRoom(roomName string, partner string) (*types.Connection, error) {
	r := getRoom(roomName)
	if r == nil {
		return nil, errors.NewServerError(404, fmt.Sprintf("could not find the room %s", roomName))
	}
	conn := r.GetConnection(partner)
	if conn != nil {
		return conn, nil
	}
	return nil, errors.NewServerError(404, fmt.Sprintf("user %s is not found in the room %s", partner, roomName))
}
