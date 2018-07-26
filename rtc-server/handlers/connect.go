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
var wsRooms sync.Map

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
	wsRooms.Store(ws, &room{})
	for {
		msg := &types.Message{}
		_, ok := wsRooms.Load(ws)
		if !ok {
			ws.Close()
			break
		}
		if err := ws.ReadJSON(msg); err != nil {
			closeErr, ok := err.(*websocket.CloseError)
			if ok {
				log.Println("need to close the room: ", closeErr.Code)
				r, ok := wsRooms.Load(ws)
				if !ok {
					log.Println("Oops, that does not look good...")
					break
				}
				if err := closeRoom(ws, r.(*room)); err != nil {
					log.Println("error during room close: ", err)
				}
				break
			}
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

func closeRoom(ws *websocket.Conn, r *room) error {
	closeMsg, err := types.NewMessageClose("That's all, folks!")
	if err != nil {
		return err
	}
	var connToInform *websocket.Conn
	if r != nil {
		if r.connector != nil {
			connToInform = r.connector.conn
		}
		if r.connectee != nil && r.connectee.conn == ws {
			connToInform = r.connector.conn
		}
	}
	if connToInform != nil {
		if err := connToInform.WriteJSON(closeMsg); err != nil {
			return err
		}
	}
	if err := ws.Close(); err != nil {
		return err
	}
	wsRooms.Delete(connToInform)
	wsRooms.Delete(ws)
	rooms.Delete(r.name)
	return nil
}
