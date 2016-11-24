package server

import (
	"net/http"

	"github.com/graarh/golang-socketio"
)

func (ps *PrimusServer) Mount() {

	ps.SocketIO.On(gosocketio.OnConnection, func(c *gosocketio.Channel, args interface{}) {
		ps.AccessLog.WithField("id", c.Id()).Infoln("New client")
	})

	ps.SocketIO.On("join", func(c *gosocketio.Channel, channel string) {
		ps.AccessLog.WithField("channel", channel).Infoln("Join")
		c.Join(channel)
	})

	http.Handle("/socket.io/", ps.SocketIO)
	http.HandleFunc("/receive", ps.ReceiveHandler)
}
