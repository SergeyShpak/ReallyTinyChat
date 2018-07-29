package types

import (
	"encoding/json"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
)

type Message struct {
	Type    string
	Message string
}

type Hello struct {
	Login string
	Room  string
}

type HelloOK struct {
	Login    string
	Room     string
	Partners []string
}

type Offer struct {
	Login      string
	Room       string
	Partner    string
	IsResponse bool
	Offer      string
}

type Ice struct {
	Candidate string
	Room      string
	Partner   string
}

type Close struct {
	Message string
}

type Error struct {
	Hint string
	Code int
}

func NewMessageHelloOK(login string, room *Room) (*Message, error) {
	payload := &HelloOK{
		Login:    login,
		Room:     room.Name,
		Partners: room.ListConnections(),
	}
	return createMessage(payload, "HELLOOK")
}

/*
func NewMessageRoomInfo(connector string, connectee string, room string) (*Message, error) {
	payload := &RoomInfo{
		Connector: connector,
		Connectee: connectee,
		Room:      room,
	}
	return createMessage(payload, "ROOMINFO")
}
*/

func NewMessageOffer(payload *Offer) (*Message, error) {
	return createMessage(payload, "OFFER")
}

func NewMessageIce(payload *Ice) (*Message, error) {
	return createMessage(payload, "ICE")
}

func NewMessageClose(message string) (*Message, error) {
	payload := &Close{
		Message: message,
	}
	return createMessage(payload, "CLOSE")
}

func NewMessageError(err *errors.ServerError) (*Message, error) {
	payload := &Error{
		Code: err.Code,
		Hint: err.Hint,
	}
	return createMessage(payload, "ERROR")
}

func createMessage(payload interface{}, msgType string) (*Message, error) {
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.NewServerError(500, err.Error())
	}
	msg := &Message{
		Type:    msgType,
		Message: string(payloadB),
	}
	return msg, nil
}
