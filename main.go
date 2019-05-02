package main

import (
	"github.com/factorysh/jaeger-traefik/cmd"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"
)

func main() {
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)
	cmd.Execute()
}
