package types

import (
	"encoding/json"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
)

type Message struct {
	Type  string
	Login string
	Room  string
	Token string
}

type ResponseMessage struct {
	Type    string
	Payload string
}

type Hello struct {
	Login string
	Room  string
}

type HelloOK struct {
	Login    string
	Secret   []byte
	Room     string
	Partners []string
}

type Offer struct {
	Login      string
	Partner    string
	IsResponse bool
	Offer      string
}

type Ice struct {
	Candidate string
	Partner   string
}

type Close struct {
	Message string
}

type Error struct {
	Hint string
	Code int
}

func NewMessageHelloOK(login string, secret []byte, room *Room) (*ResponseMessage, error) {
	payload := &HelloOK{
		Login:    login,
		Secret:   secret,
		Room:     room.Name,
		Partners: room.ListConnections(),
	}
	return createMessage(payload, "HELLOOK")
}

func NewMessageOffer(payload *Offer) (*ResponseMessage, error) {
	return createMessage(payload, "OFFER")
}

func NewMessageIce(payload *Ice) (*ResponseMessage, error) {
	return createMessage(payload, "ICE")
}

func NewMessageClose(message string) (*ResponseMessage, error) {
	payload := &Close{
		Message: message,
	}
	return createMessage(payload, "CLOSE")
}

func NewMessageError(err *errors.ServerError) (*ResponseMessage, error) {
	payload := &Error{
		Code: err.Code,
		Hint: err.Hint,
	}
	return createMessage(payload, "ERROR")
}

func createMessage(payload interface{}, msgType string) (*ResponseMessage, error) {
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.NewServerError(500, err.Error())
	}
	msg := &ResponseMessage{
		Type:    msgType,
		Payload: string(payloadB),
	}
	return msg, nil
}
