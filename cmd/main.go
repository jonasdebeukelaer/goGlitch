package main

import (
	"github.com/jonasdebeukelaer/goGlitch/api/server"
)

const (
	port = ":8080"
)

func main() {
	server.Serve(port)
}
