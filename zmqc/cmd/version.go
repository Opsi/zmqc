package cmd

import (
	"github.com/Opsi/zmqc/zmqc/logger"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Zmqc",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("Zmqc v0.1") // TODO
	},
}
