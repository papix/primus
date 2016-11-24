package server

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"

	"github.com/papix/primus/common"
	"github.com/pkg/errors"
)

type PrimusServer struct {
	Conf      ConfToml
	SocketIO  *gosocketio.Server
	Listner   net.Listener
	AccessLog *logrus.Logger
	ErrorLog  *logrus.Logger
}

func New() *PrimusServer {
	return &PrimusServer{
		Conf:     BuildDefaultConf(),
		SocketIO: gosocketio.NewServer(transport.GetDefaultWebsocketTransport()),
	}
}

func (ps *PrimusServer) Run() int {
	if err := ps.optionParse(); err != nil {
		unwrapped := common.TraceBack(err)
		if unwrapped != nil {
			ps.Errorln(unwrapped.Error())
			return 1
		}
		return 0
	}

	if err := ps.SetupLogger(); err != nil {
		// Stderr
		ps.Errorln(err.Error())
		return 1
	}
	// To write to log files
	defer ps.loggerClose()

	// Setup route
	ps.Mount()

	// Listen TCP Port
	if err := ps.ListenServer(); err != nil {
		ps.Infoln(err.Error())
		return 0
	}

	return 0
}

func (ps *PrimusServer) ListenServer() error {
	ps.AccessLog.WithField("port", ps.Conf.Server.Port).Infoln("Primus start")
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", ps.Conf.Server.Port),
		Handler: http.DefaultServeMux,
	}

	ln, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return err
	}
	ps.Listner = ln

	go ps.ListenSignal()

	return server.Serve(ps.Listner)
}

func (ps *PrimusServer) ListenSignal() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	switch <-signalCh {
	case syscall.SIGHUP:
		fallthrough
	case syscall.SIGINT:
		fallthrough
	case syscall.SIGTERM:
		fallthrough
	case syscall.SIGQUIT:
		ps.Infoln("Primus close")
		ps.Listner.Close()
		close(signalCh)
	}
}

func (ps *PrimusServer) optionParse() error {
	version := flag.Bool("v", false, "primus version")
	confPath := flag.String("c", "", "configuration file for primus")
	flag.Parse()

	if *version {
		os.Stdout.Write(common.VersionInfo())
		return common.MakeIgnore()
	}

	// Load conf
	if *confPath != "" {
		err := ps.loadConf(*confPath)
		if err != nil {
			return errors.Wrap(err, "Failed to load conf")
		}
	}

	return nil
}
