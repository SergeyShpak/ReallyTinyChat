package handlers

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/jwt"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

func (h *Handler) HandleHello(ws *websocket.Conn, login string, room string) {
	if err := h.addToConnections(ws, login, room); err != nil {
		servErr, ok := err.(*errors.ServerError)
		if !ok {
			servErr = errors.NewServerError(500, fmt.Sprintf("%s", err))
		}
		errMsg, err := types.NewMessageError(servErr)
		if err != nil {
			log.Println("error occurred: ", err)
			return
		}
		log.Println(servErr)
		ws.WriteJSON(errMsg)
		return
	}
	if err := h.sendHelloOKMessage(ws, login, room); err != nil {
		log.Println("error occurred: ", err)
		return
	}
	log.Println("HELLOOK message sent")
	return
}

func (h *Handler) addToConnections(ws *websocket.Conn, login string, room string) error {
	r, err := getRoom(room)
	var roomToAdd bool
	if err != nil {
		servErr, ok := err.(*errors.ServerError)
		if !ok {
			return errors.NewServerError(500, fmt.Sprintf("could not cast error %v to a ServerError", err))
		}
		if servErr.Code == 404 {
			r, err = h.createNewRoom(room)
			if err != nil {
				return err
			}
			log.Printf("created new room \"%s\"", r.Name)
			roomToAdd = true
		}
		if servErr.Code != 404 {
			return servErr
		}
	}
	if r.IsConnected(login) {
		return errors.NewServerError(409, fmt.Sprintf("user %s is already connected", login))
	}
	conn := &types.Connection{
		Login: login,
		WS:    ws,
	}
	if err = r.AddConnection(conn); err != nil {
		return err
	}
	secret, err := jwt.GenerateSecret()
	if err != nil {
		return err
	}
	user := &types.User{
		Login: login,
		Room:  room,
	}
	if err = h.StoreUser(ws, user, secret); err != nil {
		return err
	}
	users.Store(ws, user)
	if roomToAdd {
		rooms.Store(room, r)
	}
	return nil
}

func (h *Handler) createNewRoom(name string) (*types.Room, error) {
	r, err := types.NewRoom(name)
	if err != nil {
		servErr, ok := err.(*errors.ServerError)
		if !ok {
			return nil, errors.NewServerError(500, fmt.Sprintf("could not cast an error %v to the ServerError type", err))
		}
		return nil, servErr
	}
	return r, nil
}

func (h *Handler) sendHelloOKMessage(ws *websocket.Conn, login string, room string) error {
	r, err := getRoom(room)
	if err != nil {
		return err
	}
	secret, err := h.GetUserSecretWithConn(ws)
	if err != nil {
		return err
	}
	okMsg, err := types.NewMessageHelloOK(login, secret, r)
	if err != nil {
		return err
	}
	ws.WriteJSON(okMsg)
	return nil
}
