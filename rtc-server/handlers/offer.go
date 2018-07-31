package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

func HandleOffer(login string, room string, msg *types.Offer) error {
	repacked, err := types.NewMessageOffer(msg)
	if err != nil {
		return errors.NewServerError(500, "cannot forward an OFFER message")
	}
	r, err := getRoom(room)
	if err != nil {
		return err
	}
	log.Printf("Sending to: %s\n", msg.Partner)
	if err := r.Send(msg.Partner, repacked); err != nil {
		return errors.NewServerError(500, fmt.Sprintf("error occurred when sending an ICE message: %v", err))
	}
	return nil
}
