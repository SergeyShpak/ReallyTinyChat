package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
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
	users.Store(ws, &types.User{})
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
	t := msg.Type
	switch t {
	case "HELLO":
		log.Println("HELLO message received")
		helloMsg := &types.Hello{}
		json.Unmarshal([]byte(msg.Message), helloMsg)
		HandleHello(ws, helloMsg)
	case "OFFER":
		log.Println("OFFER message received")
		offerMsg := &types.Offer{}
		json.Unmarshal([]byte(msg.Message), offerMsg)
		HandleOffer(ws, offerMsg)
	case "ICE":
		log.Println("ICE message received")
		iceMsg := &types.Ice{}
		json.Unmarshal([]byte(msg.Message), iceMsg)
		HandleIce(ws, iceMsg)
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
	u := getUser(ws)
	if u == nil {
		return
	}
	r := getRoom(u.Room)
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
