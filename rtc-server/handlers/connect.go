package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/cache"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/config"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/errors"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/jwt"
	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

var rooms sync.Map
var users sync.Map

type Handler struct {
	cache cache.Client
}

func NewHandler(c *config.Config) (*Handler, error) {
	if c == nil {
		return nil, errors.NewServerError(http.StatusInternalServerError, "configuration object passed to the handler in nil")
	}
	h := &Handler{}
	var err error
	h.cache, err = cache.NewClient(c.Cache)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Handler) Connect(w http.ResponseWriter, r *http.Request) {
	if err := h.createConnection(w, r); err != nil {
		log.Printf("could not create a connection: %v", err)
		WriteResponse(w, &Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("could not create a connection: %v", err),
		})
		return
	}
}

func (h *Handler) createConnection(w http.ResponseWriter, r *http.Request) error {
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
		h.listenToMessages(ws)
		wg.Done()
	}()
	wg.Wait()
	return nil
}

func (h *Handler) listenToMessages(ws *websocket.Conn) {
	for {
		msg := &types.Message{}
		if err := ws.ReadJSON(msg); err != nil {
			h.handleListenMsgError(ws, err)
			break
		}
		if err := h.handleMessage(ws, msg); err != nil {
			log.Println("an error occurred during message handling: ", err)
			continue
		}
	}
}

func (h *Handler) handleListenMsgError(ws *websocket.Conn, err error) {
	h.removeConnection(ws)
	log.Println("an error occurred during message reading: ", err)
	return
}

func (h *Handler) handleMessage(ws *websocket.Conn, msg *types.Message) error {
	payload, err := h.verifyMessage(ws, msg)
	if err != nil {
		return err
	}
	t := msg.Type
	switch t {
	case "HELLO":
		log.Println("HELLO message received")
		h.HandleHello(ws, msg.Login, msg.Room)
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

func (h *Handler) removeConnection(ws *websocket.Conn) {
	log.Println("Removing connection")
	ws.Close()
	u, _, err := h.GetUserWithConn(ws)
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
	if err = h.RemoveUserWithConn(ws); err != nil {
		log.Println("error occurred during user removal: ", err)
	}
	log.Printf("Closed connection to user \"%s\"\n", u.Login)
	if r.IsEmpty() {
		rooms.Delete(u.Room)
		log.Printf("Room \"%s\" went empty and was removed\n", u.Room)
		return
	}
}

func (h *Handler) verifyMessage(ws *websocket.Conn, msg *types.Message) (payload string, err error) {
	if msg.Type == "HELLO" {
		return "", nil
	}
	user, secret, err := h.GetUserWithConn(ws)
	if err != nil {
		return "", err
	}
	if user.Login != msg.Login || user.Room != msg.Room {
		return "", errors.NewServerError(http.StatusBadRequest,
			fmt.Sprintf("WebSocket connection is associated with the user \"%s\" in the room \"%s\", not with user \"%s\" in the room \"%s\"",
				user.Login, user.Room, msg.Login, msg.Room))
	}
	return jwt.Verify(msg, secret)
}
