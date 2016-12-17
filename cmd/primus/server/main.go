package main

import (
	"os"

	"github.com/papix/primus/server"
)

func main() {
	os.Exit(server.New().Run())
}
