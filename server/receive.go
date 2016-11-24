package server

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
)

type PrimusPayload struct {
	Channel string              `json:"channel"`
	Method  string              `json:"string"`
	Query   string              `json:"query"`
	Body    string              `json:"body"`
	Header  map[string][]string `json:"header"`
}

func (ps *PrimusServer) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	ps.AccessLog.WithFields(logrus.Fields{
		"method":     r.Method,
		"remoteAddr": r.RemoteAddr,
		"uri":        r.RequestURI,
		"path":       r.URL.Path,
	}).Infoln("Requests")

	channel := strings.TrimLeft(r.URL.Path, "/receive/")

	if channel == "" {
		ps.sendResponse(w, "invalid url", http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	reqBody := bufio.NewWriter(&buf)

	if _, err := io.Copy(reqBody, r.Body); err != nil {
		ps.sendResponse(w, "failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	payload := &PrimusPayload{
		Channel: channel,
		Method:  r.Method,
		Body:    buf.String(),
		Header:  r.Header,
		Query:   r.URL.RawQuery,
	}
	ps.SocketIO.BroadcastTo(channel, "receive", payload)

	ps.sendResponse(w, "ok", http.StatusOK)
}
