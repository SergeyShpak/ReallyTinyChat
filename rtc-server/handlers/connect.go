package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/jwt"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

var rooms sync.Map
var users sync.Map

func Connect(w http.ResponseWriter, r *http.Request) {
	if err := createConnection(w, r); err != nil {
		log.Printf("could not create a connection: %v", err)
		WriteResponse(w, &Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("could not create a connection: %v", err),
		})
		return
	}
}

func createConnection(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		// TODO(SSH): what about origins?
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return errors.NewServerError(500, fmt.Sprintf("upgrader failed: %v", err))
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		listenToMessages(ws)
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func listenToMessages(ws *websocket.Conn) {
	for {
		msg := &types.Message{}
		if err := ws.ReadJSON(msg); err != nil {
			handleListenMsgError(ws, err)
			break
		}
		if err := handleMessage(ws, msg); err != nil {
			log.Println("an error occurred during message handling: ", err)
			continue
		}
	}
}

func handleListenMsgError(ws *websocket.Conn, err error) {
	removeConnection(ws)
	log.Println("an error occurred during message reading: ", err)
	return
}

func handleMessage(ws *websocket.Conn, msg *types.Message) error {
	payload, err := verifyMessage(ws, msg)
	if err != nil {
		return err
	}
	t := msg.Type
	switch t {
	case "HELLO":
		log.Println("HELLO message received")
		HandleHello(ws, msg.Login, msg.Room)
	case "OFFER":
		log.Println("OFFER message received")
		offerMsg := &types.Offer{}
		if err := json.Unmarshal([]byte(payload), offerMsg); err != nil {
			log.Println("could not unmarshal as OFFER: ", payload)
			return err
		}
		if err := HandleOffer(msg.Login, msg.Room, offerMsg); err != nil {
			return err
		}
	case "ICE":
		log.Println("ICE message received")
		iceMsg := &types.Ice{}
		if err := json.Unmarshal([]byte(payload), iceMsg); err != nil {
			log.Println("could not unmarshal as ICE: ", payload)
			return err
		}
		if err := HandleIce(msg.Login, msg.Room, iceMsg); err != nil {
			return err
		}
	default:
		servErr := errors.NewServerError(400, fmt.Sprintf("the server does not know about \"%s\" message type", t))
		errMsg, err := types.NewMessageError(servErr)
		if err != nil {
			return err
		}
		log.Println("An error occurred: ", servErr)
		ws.WriteJSON(errMsg)
		return nil
	}
	return nil
}

func removeConnection(ws *websocket.Conn) {
	log.Println("Removing connection")
	ws.Close()
	u, err := getUser(ws)
	if err != nil {
		log.Println("an error occurred: ", err)
		return
	}
	r, err := getRoom(u.Room)
	if err != nil {
		log.Println("an error occurred: ", err)
	}
	r.RemoveConnection(u.Login)
	log.Printf("Removed user \"%s\" from room \"%s\"\n", u.Login, u.Room)
	users.Delete(ws)
	log.Printf("Closed connection to user \"%s\"\n", u.Login)
	if r.IsEmpty() {
		rooms.Delete(u.Room)
		log.Printf("Room \"%s\" went empty and was removed\n", u.Room)
		return
	}
}

func verifyMessage(ws *websocket.Conn, msg *types.Message) (payload string, err error) {
	if msg.Type == "HELLO" {
		return "", nil
	}
	u, err := getUser(ws)
	if err != nil {
		return "", err
	}
	if u.Login != msg.Login || u.Room != msg.Room {
		return "", errors.NewServerError(http.StatusBadRequest,
			fmt.Sprintf("WebSocket connection is associated with the user \"%s\" in the room \"%s\", not with user \"%s\" in the room \"%s\"", u.Login, u.Room, msg.Login, msg.Room))
	}
	return jwt.Verify(msg, u.Secret)
}
