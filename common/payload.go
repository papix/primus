package common

import "net/http"

type PrimusPayload struct {
	Channel string
	Method  string
	Query   string
	Body    []byte
	Header  http.Header
}
