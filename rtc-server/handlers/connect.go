package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
	"github.com/gorilla/websocket"
)

var rooms = make(map[string]*room)

type room struct {
	name      string
	connector *connection
	connectee *connection
}

type connection struct {
	login string
	conn  *websocket.Conn
}

func Connect(w http.ResponseWriter, r *http.Request) {
	if err := createConnection(w, r); err != nil {
		log.Printf("could not create a connection: %v", err)
		WriteResponse(w, &Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("could not create a connection: %v", err),
		})
		return
	}
	log.Printf("OK!")
}

func createConnection(w http.ResponseWriter, r *http.Request) error {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		listenToMessages(ws)
		wg.Done()
	}()
	log.Println("Here1")
	wg.Wait()
	return nil
}

func listenToMessages(ws *websocket.Conn) {
	for {
		msg := &types.Message{}
		if err := ws.ReadJSON(msg); err != nil {
			log.Println("an error occurred during message reading: ", err)
			continue
		}
		if err := handleMessage(ws, msg); err != nil {
			log.Println("an error occurred during message handling: ", err)
			continue
		}
	}
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
		return fmt.Errorf("type \"%s\" is unknown", t)
	}
	return nil
}
