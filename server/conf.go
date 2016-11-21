package server

import (
	"github.com/BurntSushi/toml"
)

type ConfToml struct {
	Server ServerSection `toml:"server"`
	Log    LogSection    `toml:"log"`
}

type ServerSection struct {
	Port string `toml:"port"`
}

type LogSection struct {
	AccessLog string `toml:"access_log"`
	ErrorLog  string `toml:"error_log"`
	Level     string `toml:"level"`
}

func init() {
	Conf = BuildDefaultConf()
}

func BuildDefaultConf() ConfToml {
	var conf ConfToml

	// Server
	conf.Server.Port = "14300"

	// Log
	conf.Log.AccessLog = "stdout"
	conf.Log.ErrorLog = "stderr"
	conf.Log.Level = "error"

	return conf
}

func LoadConf(confPath string, conf *ConfToml) error {
	_, err := toml.DecodeFile(confPath, conf)
	if err != nil {
		return err
	}
	return nil
}
