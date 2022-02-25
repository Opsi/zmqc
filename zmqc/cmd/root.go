package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	host string
	port uint
)

var rootCmd = &cobra.Command{
	Use:   "zmqc",
	Short: "Zmqc is a command line tool for ZeroMQ",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().UintVarP(&port, "port", "p", 5555, "Port to connect to")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(subCmd)
	rootCmd.AddCommand(pubCmd)
}
