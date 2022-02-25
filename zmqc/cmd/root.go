package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var rootCmd = &cobra.Command{
	Use:   "zmqc",
	Short: "zmqc is a command line tool for ZeroMQ",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Root command!") // TODO
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of zmqc",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Zmqc v0.1") // TODO
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
