package cmd

import (
	"net/http"
	"os"

	"github.com/factorysh/jaeger-traefik/conf"
	"github.com/factorysh/jaeger-traefik/server"
	"github.com/getsentry/raven-go"
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

		dsn := os.Getenv("SENTRY_DSN")
		if dsn != "" {
			raven.SetDSN(dsn)
		}

		cfgPath := os.Getenv("CONFIG")
		if cfgPath == "" {
			cfgPath = "/etc/jaeger-traefik.yml"
		}

		cfg, err := conf.Read(cfgPath)
		if err != nil {
			return err
		}
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		if dsn != "" {
			http.HandleFunc("/", raven.RecoveryHandler(mux.ServeHTTP))
		} else {
			http.Handle("/", mux)
		}
		log.WithField("listen", cfg.ListenAdmin).Info("Listening admin")
		go http.ListenAndServe(cfg.ListenAdmin, nil)

		s, err := server.NewServer(cfg)
		if err != nil {
			return err
		}
		s.Serve()
		return nil
	},
}
