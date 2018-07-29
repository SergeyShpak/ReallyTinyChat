package types

import "sync"

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
	r.ConnectionsMux.Lock()
	_, ok := r.Connections[login]
	r.ConnectionsMux.Unlock()
	return ok
}

func (r *Room) GetConnection(login string) *Connection {
	r.ConnectionsMux.Lock()
	conn, ok := r.Connections[login]
	r.ConnectionsMux.Unlock()
	if !ok {
		return nil
	}
	return conn
}

func (r *Room) ListConnections() []string {
	r.ConnectionsMux.Lock()
	connections := make([]string, 0, len(r.Connections))
	for c := range r.Connections {
		connections = append(connections, c)
	}
	r.ConnectionsMux.Unlock()
	return connections
}
