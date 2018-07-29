package types

import (
	"fmt"
	"sync"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
)

type Room struct {
	Name string

	Connections    map[string]*Connection
	ConnectionsMux sync.RWMutex
}

func NewRoom(roomName string) *Room {
	r := &Room{
		Name:        roomName,
		Connections: make(map[string]*Connection),
	}
	return r
}

func (r *Room) AddConnection(conn *Connection) {
	if conn == nil {
		return
	}
	r.ConnectionsMux.Lock()
	r.Connections[conn.Login] = conn
	r.ConnectionsMux.Unlock()
	return
}

func (r *Room) RemoveConnection(login string) {
	r.ConnectionsMux.Lock()
	delete(r.Connections, login)
	r.ConnectionsMux.Unlock()
}

func (r *Room) IsConnected(login string) bool {
	r.ConnectionsMux.RLock()
	_, ok := r.Connections[login]
	r.ConnectionsMux.RUnlock()
	return ok
}

func (r *Room) ListConnections() []string {
	r.ConnectionsMux.RLock()
	connections := make([]string, 0, len(r.Connections))
	for c := range r.Connections {
		connections = append(connections, c)
	}
	r.ConnectionsMux.RUnlock()
	return connections
}

func (r *Room) IsEmpty() bool {
	r.ConnectionsMux.RLock()
	isEmpty := len(r.Connections) == 0
	r.ConnectionsMux.RUnlock()
	return isEmpty
}

func (r *Room) Send(login string, msg interface{}) error {
	var err error
	r.ConnectionsMux.RLock()
	err = r.send(login, msg)
	r.ConnectionsMux.RUnlock()
	return err
}

func (r *Room) send(login string, msg interface{}) error {
	conn, ok := r.Connections[login]
	if !ok {
		return errors.NewServerError(404, fmt.Sprintf("user %s was not found in the room %s", login, r.Name))
	}
	conn.WS.WriteJSON(msg)
	return nil
}
