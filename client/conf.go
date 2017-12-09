package client

import (
	"github.com/BurntSushi/toml"
)

type ConfToml struct {
	Server ServerSection  `toml:"server"`
	Route  []RouteSection `toml:"route"`
}

type ServerSection struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	SSL  bool   `toml:"ssl"`
}

type RouteSection struct {
	Channel string `toml:"channel"`
	Host    string `toml:"host"`
	Port    int    `toml:"port"`
	Path    string `toml:"path"`
	SSL     bool   `toml:"ssl"`
}

func BuildDefaultConf() ConfToml {
	var conf ConfToml

	// Server
	conf.Server.Host = "localhost"
	conf.Server.Port = 14300
	conf.Server.SSL = false

	return conf
}

func LoadConf(confPath string, conf *ConfToml) error {
	_, err := toml.DecodeFile(confPath, conf)
	if err != nil {
		return err
	}
	return nil
}

func (conf *ConfToml) FetchRouteByChannel(channel string) RouteSection {
	for _, r := range conf.Route {
		if r.Channel == channel {
			return r
		}
	}
	return RouteSection{}
}
