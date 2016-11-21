package client

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type PrimusPayload struct {
	Channel string              `json:"channel"`
	Method  string              `json:"string"`
	Query   string              `json:"query"`
	Body    string              `json:"body"`
	Header  map[string][]string `json:"header"`
}

func Run() {
	client, err := gosocketio.Dial(
		gosocketio.GetUrl(Conf.Server.Host, Conf.Server.Port, Conf.Server.SSL),
		transport.GetDefaultWebsocketTransport(),
	)
	defer client.Close()

	if err != nil {
		log.Fatal(err)
	}

	for c, _ := range Conf.Route {
		client.Emit("join", c)
	}

	client.On("receive", func(c *gosocketio.Channel, payload PrimusPayload) {
		url := Conf.Route[payload.Channel]
		if payload.Query != "" {
			url = fmt.Sprintf("%s?%s", url, payload.Query)
		}

		body := strings.NewReader(payload.Body)
		req, _ := http.NewRequest(payload.Method, url, body)
		for h, v := range payload.Header {
			req.Header.Set(h, v[0])
		}

		client := new(http.Client)
		_, err := client.Do(req)

		if err != nil {
			log.Print(err)
		}
	})

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM)

	exitCh := make(chan int)
	go func() {
		for {
			s := <-signalCh
			switch s {
			case syscall.SIGTERM:
				exitCh <- 0
			default:
				exitCh <- 1
			}
		}
	}()

	code := <-exitCh
	os.Exit(code)
}
