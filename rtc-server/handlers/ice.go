package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

// TODO: cannot merge with HandleOffer?
func HandleIce(login string, room string, msg *types.Ice) error {
	log.Printf("Received ICE message, forwarding to %s\n", msg.Partner)
	repacked, err := types.NewMessageIce(msg)
	if err != nil {
		return errors.NewServerError(500, "cannot forward the ICE message")
	}
	r, err := getRoom(room)
	if err != nil {
		return err
	}
	if err := r.Send(msg.Partner, repacked); err != nil {
		return errors.NewServerError(500, fmt.Sprintf("error occurred when sending an ICE message: %v", err))
	}
	return nil
}
