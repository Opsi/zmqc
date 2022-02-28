package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Zmqc",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print("Zmqc v0.1") // TODO
	},
}
