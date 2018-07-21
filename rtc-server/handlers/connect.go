package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/SergeyShpak/ReallyTinyChat/rtc-server/types"
)

var connections = make(map[string]*connection, 0)

type connection struct {
	username string
}

func Connect(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("err while reading request body: %v", err)
		WriteResponse(w, &Response{
			Status:  http.StatusInternalServerError,
			Message: fmt.Sprintf("err while reading request body: %v", err),
		})
		return
	}
	req := &types.ConnectionRequest{}
	if err = json.Unmarshal(b, req); err != nil {
		log.Printf("err while unmarshalling request body: %v", err)
		WriteResponse(w, &Response{
			Status:  http.StatusBadRequest,
			Message: fmt.Sprintf("err while unmarshalling request body: %v", err),
		})
		return
	}
	if err = checkConnectionRequest(req); err != nil {
		log.Printf("request is invalid: %v", err)
		WriteResponse(w, &Response{
			Status:  http.StatusConflict,
			Message: fmt.Sprintf("request is invalid: %v", err),
		})
		return
	}
	createConnection(req)
	WriteResponse(w, &Response{
		Status: http.StatusCreated,
	})
	log.Printf("OK!")
}

func checkConnectionRequest(req *types.ConnectionRequest) error {
	if _, ok := connections[req.Username]; ok {
		return fmt.Errorf("username already exists")
	}
	return nil
}

func createConnection(req *types.ConnectionRequest) {
	conn := &connection{
		username: req.Username,
	}
	connections[req.Username] = conn
}
