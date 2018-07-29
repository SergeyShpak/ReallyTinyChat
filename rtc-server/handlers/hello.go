package handlers

import (
	"fmt"
	"log"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

func HandleHello(ws *websocket.Conn, msg *types.Hello) {
	if err := addToConnections(ws, msg); err != nil {
		servErr, ok := err.(*errors.ServerError)
		if !ok {
			servErr = errors.NewServerError(500, fmt.Sprintf("%s", err))
		}
		errMsg, err := types.NewMessageError(servErr)
		if err != nil {
			log.Println("error occurred: ", err)
			return
		}
		ws.WriteJSON(errMsg)
		log.Println(servErr)
		return
	}
	if err := sendHelloOKMessage(ws, msg); err != nil {
		log.Println("error occurred: ", err)
		return
	}
	log.Println("HELLOOK message sent")
	return
}

func addToConnections(ws *websocket.Conn, msg *types.Hello) error {
	r, err := getRoom(msg.Room)
	if err != nil {
		log.Println("error occurred: ", err)
		servErr, ok := err.(*errors.ServerError)
		if !ok {
			return errors.NewServerError(500, fmt.Sprintf("could not cast error %v to a ServerError", err))
		}
		if servErr.Code == 404 {
			r, err = createNewRoom(msg.Room)
			if err != nil {
				return err
			}
			addRoom(r)
		}
		if servErr.Code != 404 {
			return servErr
		}
	}
	if r.IsConnected(msg.Login) {
		return errors.NewServerError(409, fmt.Sprintf("user %s is already connected", msg.Login))
	}
	conn := &types.Connection{
		Login: msg.Login,
		WS:    ws,
	}
	user := &types.User{
		Login: msg.Login,
		Room:  msg.Room,
	}
	users.Store(ws, user)
	if err = r.AddConnection(conn); err != nil {
		return err
	}
	return nil
}

func createNewRoom(name string) (*types.Room, error) {
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

func sendHelloOKMessage(ws *websocket.Conn, msg *types.Hello) error {
	r, err := getRoom(msg.Room)
	if err != nil {
		return err
	}
	okMsg, err := types.NewMessageHelloOK(msg.Login, r)
	if err != nil {
		return err
	}
	ws.WriteJSON(okMsg)
	return nil
}
