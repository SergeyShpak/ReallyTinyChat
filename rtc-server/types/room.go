package types

import (
	"fmt"
	"sync"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
)

type RoomOpts struct {
	Capacity int
}

type Room struct {
	Capacity int
	Name     string

	Connections    map[string]*Connection
	ConnectionsMux sync.RWMutex
}

func NewRoom(roomName string, opts ...*RoomOpts) (*Room, error) {
	o, err := getRoomOpts(opts)
	if err != nil {
		return nil, err
	}
	r := &Room{
		Name:        roomName,
		Connections: make(map[string]*Connection),
	}
	if o.Capacity < 2 {
		return nil, errors.NewServerError(400, "room's capacity cannot be less than 2")
	}
	r.Capacity = o.Capacity
	return r, nil
}

func getRoomOpts(opts []*RoomOpts) (*RoomOpts, error) {
	if len(opts) == 0 {
		return &RoomOpts{
			Capacity: 2,
		}, nil
	}
	if len(opts) == 1 {
		return opts[1], nil
	}
	return nil, errors.NewServerError(500, fmt.Sprintf("too many (%d) options passed to the Room constructor", len(opts)))
}

func (r *Room) AddConnection(conn *Connection) error {
	if conn == nil {
		return errors.NewServerError(500, fmt.Sprintf("trying to add a nil connection to the room \"%s\"", r.Name))
	}
	r.ConnectionsMux.Lock()
	err := r.addConnection(conn)
	r.ConnectionsMux.Unlock()
	return err
}

func (r *Room) addConnection(conn *Connection) error {
	if len(r.Connections) >= r.Capacity {
		return errors.NewServerError(400, fmt.Sprintf("room \"%s\" is full", r.Name))
	}
	r.Connections[conn.Login] = conn
	return nil
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
