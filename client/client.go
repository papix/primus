package client

import (
	"bytes"
	"compress/gzip"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	for _, r := range pc.Conf.Route {
		socket.Emit("join", r.Channel)
	}

	socket.On("receive", func(c *gosocketio.Channel, payload common.PrimusPayload) {
		r := pc.Conf.FetchRouteByChannel(payload.Channel)

		var url string
		var scheme string
		var host string
		var port string
		var path string

		if r.SSL {
			scheme = "https"
		} else {
			scheme = "http"
		}

		if r.Host != "" {
			host = r.Host
		} else {
			host = "localhost"
		}

		if r.Port != 0 {
			port = fmt.Sprintf(":%d", r.Port)
		}

		if r.Path != "" {
			if strings.HasPrefix(r.Path, "/") {
				path = r.Path
			} else {
				path = fmt.Sprintf("/%s", r.Path)
			}
		}

		if payload.Query != "" {
			url = fmt.Sprintf("%s://%s%s%s?%s", scheme, host, port, path, payload.Query)
		} else {
			url = fmt.Sprintf("%s://%s%s%s", scheme, host, port, path)
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
