package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/papix/primus"
)

type PrimusResponse struct {
	Message string `json:"message"`
}

func (ps *PrimusServer) sendResponse(w http.ResponseWriter, msg string, code int) {
	var (
		respBody   []byte
		respPrimus PrimusResponse
	)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Server", fmt.Sprintf("primus %s", primus.Version))

	respPrimus.Message = msg
	respBody, err := json.Marshal(respPrimus)

	if err != nil {
		msg := "Response-body could not be created"
		http.Error(w, msg, http.StatusInternalServerError)
		ps.Errorln(msg)
		return
	}

	switch code {
	case http.StatusOK:
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	default:
		w.WriteHeader(code)
		w.Write(respBody)
		ps.Errorln("status code:", code)
	}
}
