package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/papix/primus/common"
)

func (ps *PrimusServer) ReceiveHandler(w http.ResponseWriter, r *http.Request) {
	channel := strings.TrimLeft(r.URL.Path, "/receive/")

	ps.AccessLog.WithFields(logrus.Fields{
		"method":     r.Method,
		"remoteAddr": r.RemoteAddr,
		"uri":        r.RequestURI,
		"channel":    channel,
	}).Infoln("Request")

	if channel == "" {
		ps.sendResponse(w, "invalid url", http.StatusInternalServerError)
		return
	}

	zipped, err := compress(r.Body)
	if err != nil {
		ps.sendResponse(w, "failed to read request body", http.StatusInternalServerError)
		return
	}
	payload := &common.PrimusPayload{
		Channel: channel,
		Method:  r.Method,
		Body:    zipped.Bytes(),
		Header:  r.Header,
		Query:   r.URL.RawQuery,
	}
	ps.SocketIO.BroadcastTo(channel, "receive", payload)

	ps.sendResponse(w, "ok", http.StatusOK)
}

func compress(body io.ReadCloser) (*bytes.Buffer, error) {

	buf := new(bytes.Buffer)
	zipped := new(bytes.Buffer)
	if _, err := buf.ReadFrom(body); err != nil {
		return nil, err
	}

	reqBody, _ := gzip.NewWriterLevel(zipped, gzip.BestCompression)
	reqBody.Write(buf.Bytes())
	reqBody.Close()

	return zipped, nil
}
