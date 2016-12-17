package client

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/papix/primus/common"
	"github.com/pkg/errors"
	"github.com/uber-go/zap"
)

type PrimusClient struct {
	Conf   ConfToml
	Logger zap.Logger
}

func New(Out zap.WriteSyncer) *PrimusClient {
	return &PrimusClient{
		Conf: BuildDefaultConf(),
		Logger: zap.New(
			zap.NewTextEncoder(zap.TextTimeFormat(time.ANSIC)),
			zap.AddCaller(), // Add line number option
			zap.Output(Out),
		),
	}
}

func (pc *PrimusClient) optionParse() error {
	version := flag.Bool("v", false, "primus version")
	confPath := flag.String("c", "", "configuration file for primus")
	flag.Parse()

	if *version {
		os.Stdout.Write(common.VersionInfo())
		return common.MakeIgnore()
	}

	// Load conf
	if *confPath != "" {
		err := LoadConf(*confPath, &pc.Conf)
		if err != nil {
			return errors.Wrap(err, "Failed to load conf")
		}
	}

	return nil
}

func (pc *PrimusClient) Run() int {

	if err := pc.optionParse(); err != nil {
		unwrapped := common.TraceBack(err)
		if unwrapped != nil {
			pc.Logger.Error(unwrapped.Error())
			return 1
		}
		return 0
	}

	socket, err := gosocketio.Dial(
		gosocketio.GetUrl(pc.Conf.Server.Host, pc.Conf.Server.Port, pc.Conf.Server.SSL),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		pc.Logger.Error(err.Error())
		return 1
	}
	defer socket.Close()

	for _, c := range pc.Conf.Route {
		socket.Emit("join", c.Channel)
	}

	socket.On("receive", func(c *gosocketio.Channel, payload common.PrimusPayload) {
		var url string
		if payload.Query != "" {
			url = fmt.Sprintf("http://localhost:8080/?%s", payload.Query)
		} else {
			url = fmt.Sprintf("http://localhost:8080/")
		}

		// gzip decompression
		body, err := gzip.NewReader(bytes.NewBuffer(payload.Body))
		if err != nil {
			pc.Logger.Error(err.Error())
			return
		}
		defer body.Close()

		// create request
		req, err := http.NewRequest(payload.Method, url, body)
		if err != nil {
			pc.Logger.Error(err.Error())
			return
		}
		req.Header = payload.Header

		if _, err = http.DefaultClient.Do(req); err != nil {
			pc.Logger.Error(err.Error())
		}
	})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)

	exitCh := make(chan int)
	go func() {
		s := <-signalCh
		switch s {
		case syscall.SIGTERM:
			exitCh <- 0
		default:
			exitCh <- 1
		}
	}()

	code := <-exitCh
	close(exitCh)
	return code
}
