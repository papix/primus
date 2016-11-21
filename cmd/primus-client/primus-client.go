package main

import (
	"flag"
	"log"

	"github.com/papix/primus"
	"github.com/papix/primus/client"
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
		err := client.LoadConf(*confPath, &client.Conf)
		if err != nil {
			log.Fatal(err)
		}
	}

	client.Run()
}
