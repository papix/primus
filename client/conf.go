package client

import "github.com/BurntSushi/toml"

type ConfToml struct {
	Server ServerSection     `toml:"server"`
	Route  map[string]string `toml:"route"`
}

type ServerSection struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	SSL  bool   `toml:"ssl"`
}

func init() {
	Conf = BuildDefaultConf()
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
