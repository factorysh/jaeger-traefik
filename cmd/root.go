package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Verbose bool

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Verbose")
	cobra.OnInitialize(func() {
		if Verbose {
			log.SetLevel(log.DebugLevel)
			log.WithField("level", log.GetLevel()).Info("Log level")
		}
	})
}

var rootCmd = &cobra.Command{
	Use:   "jaeger-lite",
	Short: "Jaeger daemon with continuous consolidation",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
