package server

import (
	"github.com/BurntSushi/toml"
)

type ConfToml struct {
	Server ServerSection `toml:"server"`
	Log    LogSection    `toml:"log"`
}

type ServerSection struct {
	Port int `toml:"port"`
}

type LogSection struct {
	AccessLog string `toml:"access_log"`
	ErrorLog  string `toml:"error_log"`
	Level     string `toml:"level"`
}

func BuildDefaultConf() ConfToml {
	var conf ConfToml

	// Server
	conf.Server.Port = 14300
	// 	conf.Log.Level = "error"
	conf.Log.Level = "debug"

	return conf
}

func (ps *PrimusServer) loadConf(confPath string) error {
	_, err := toml.DecodeFile(confPath, &ps.Conf)
	if err != nil {
		return err
	}
	return nil
}
