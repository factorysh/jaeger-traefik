package main

import (
	"net/http"
	"os"

	"github.com/factorysh/jaeger-lite/server"
	"github.com/onrik/logrus/filename"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {
	filenameHook := filename.NewHook()
	log.AddHook(filenameHook)
	log.SetLevel(log.DebugLevel)

	http.Handle("/metrics", promhttp.Handler())
	adminListen := os.Getenv("ADMIN_LISTEN")
	if adminListen == "" {
		adminListen = "127.0.0.1:8080"
	}
	log.WithField("listen", adminListen).Info("Listening admin")
	go http.ListenAndServe(adminListen, nil)

	s, err := server.New()
	if err != nil {
		panic(err)
	}
	s.Serve()
}
