package main

import (
	"flag"
	"log"

	"github.com/papix/primus"
	"github.com/papix/primus/server"
)

func main() {
	version := flag.Bool("v", false, "primus version")
	confPath := flag.String("c", "", "configuration file for primus")
	flag.Parse()

	if *version {
		primus.PrintVersion()
		return
	}

	// Load conf
	if *confPath != "" {
		err := server.LoadConf(*confPath, &server.Conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	// set logger
	err := server.SetLogLevel(server.LogAccess, "info")
	if err != nil {
		log.Fatal(err)
	}

	err = server.SetLogLevel(server.LogError, server.Conf.Log.Level)
	if err != nil {
		log.Fatal(err)
	}
	err = server.SetLogOut(server.LogAccess, server.Conf.Log.AccessLog)
	if err != nil {
		log.Fatal(err)
	}
	err = server.SetLogOut(server.LogError, server.Conf.Log.ErrorLog)
	if err != nil {
		log.Fatal(err)
	}

	server.Run()
}
