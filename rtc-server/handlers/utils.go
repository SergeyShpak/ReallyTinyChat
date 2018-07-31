package handlers

import (
	"fmt"

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
