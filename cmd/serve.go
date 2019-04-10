package cmd

import (
	"net/http"
	"os"

	"github.com/factorysh/jaeger-lite/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve as a Jæger daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		http.Handle("/metrics", promhttp.Handler())
		adminListen := os.Getenv("ADMIN_LISTEN")
		if adminListen == "" {
			adminListen = "127.0.0.1:8080"
		}
		log.WithField("listen", adminListen).Info("Listening admin")
		go http.ListenAndServe(adminListen, nil)

		s, err := server.New()
		if err != nil {
			return err
		}
		s.Serve()
		return nil
	},
}
