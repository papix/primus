package main

import (
	"os"

	"github.com/papix/primus/client"
)

func main() {
	os.Exit(client.New(os.Stdout).Run())
}
