package main

import (
	"fmt"

	"github.com/factorysh/jaeger-lite/server"
)

func main() {
	s, err := server.New()
	if err != nil {
		panic(err)
	}
	fmt.Println(s)
}
