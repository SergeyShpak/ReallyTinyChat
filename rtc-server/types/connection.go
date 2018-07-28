package types

import "github.com/gorilla/websocket"

type Connection struct {
	Login string
	WS    *websocket.Conn
}
