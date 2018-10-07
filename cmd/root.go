package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/factorysh/jaeger-lite/server"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose")
	cobra.OnInitialize(func() {
		if Verbose {
			log.SetLevel(log.DebugLevel)
		}
		log.WithField("level", log.GetLevel()).Info("Log level")
	})
}

var rootCmd = &cobra.Command{
	Use:   "jaeger-lite",
	Short: "Jaeger daemon with continuous consolidation",
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
