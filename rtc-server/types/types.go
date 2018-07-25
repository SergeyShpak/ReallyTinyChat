package types

import (
	"encoding/json"
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
	Login string
	Room  string
}

type RoomInfo struct {
	Connector string
	Connectee string
	Room      string
}

type Offer struct {
	Offer      string
	IsResponse bool
	Room       string
}

type Ice struct {
	Candidate string
	Room      string
}

type Close struct {
	Message string
}

type Error struct {
	Hint string
	Code int
}

func NewMessageHelloOK(login string, room string) (*Message, error) {
	payload := &HelloOK{
		Login: login,
		Room:  room,
	}
	return createMessage(payload, "HELLOOK")
}

func NewMessageRoomInfo(connector string, connectee string, room string) (*Message, error) {
	payload := &RoomInfo{
		Connector: connector,
		Connectee: connectee,
		Room:      room,
	}
	return createMessage(payload, "ROOMINFO")
}

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

func createMessage(payload interface{}, msgType string) (*Message, error) {
	payloadB, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		Type:    msgType,
		Message: string(payloadB),
	}
	return msg, nil
}
