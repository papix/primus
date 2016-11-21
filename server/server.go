package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

func Run() {
	SocketIO = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	SocketIO.On(gosocketio.OnConnection, func(c *gosocketio.Channel, args interface{}) {
		LogError.Debug("New client: ", c.Id())
	})

	SocketIO.On("join", func(c *gosocketio.Channel, channel string) {
		LogError.Debug("Join: ", channel)
		c.Join(channel)
	})

	http.Handle("/socket.io/", SocketIO)
	http.HandleFunc("/receive/", ReceiveHandler)

	// Listen TCP Port
	if _, err := strconv.Atoi(Conf.Server.Port); err == nil {
		http.ListenAndServe(":"+Conf.Server.Port, nil)
	}

	log.Fatalf("port parameter is invalid: %s", Conf.Server.Port)
}
