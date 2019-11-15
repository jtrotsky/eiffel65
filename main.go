package main

import (
	"log"

	"github.com/jtrotsky/eiffel65/server"
)

func main() {
	log.Fatal(server.Start())
}
