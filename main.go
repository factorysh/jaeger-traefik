package main

import (
	"github.com/factorysh/jaeger-lite/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		panic(err)
	}
	s.Serve()
}
