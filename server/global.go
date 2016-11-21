package server

import (
	"github.com/Sirupsen/logrus"
	"github.com/graarh/golang-socketio"
)

var (
	Conf      ConfToml
	SocketIO  *gosocketio.Server
	LogAccess *logrus.Logger
	LogError  *logrus.Logger
)
