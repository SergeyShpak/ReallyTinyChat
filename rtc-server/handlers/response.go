package handlers

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int
	Message string
}

func WriteResponse(w http.ResponseWriter, resp *Response) {
	status := resp.Status
	msg, err := json.Marshal(resp.Message)
	if err != nil {
		msg = nil
		status = http.StatusInternalServerError
	}
	w.WriteHeader(status)
	w.Write(msg)
}
