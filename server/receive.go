package server

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type PrimusPayload struct {
	Channel string              `json:"channel"`
	Method  string              `json:"string"`
	Query   string              `json:"query"`
	Body    string              `json:"body"`
	Header  map[string][]string `json:"header"`
}

func ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	channel := strings.TrimLeft(r.URL.Path, "/receive/")

	if channel == "" {
		sendResponse(w, "invalid url", http.StatusInternalServerError)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendResponse(w, "failed to read request body", http.StatusInternalServerError)
		return
	}

	payload := &PrimusPayload{
		Channel: channel,
		Method:  r.Method,
		Body:    string(reqBody),
		Header:  r.Header,
		Query:   r.URL.RawQuery,
	}
	SocketIO.BroadcastTo(channel, "receive", payload)

	sendResponse(w, "ok", http.StatusOK)
}
