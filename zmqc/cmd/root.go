package cmd

import (
	"fmt"
	"os"

	"github.com/Opsi/zmqc/zmqc/logger"
	"github.com/spf13/cobra"
)

var (
	host           string
	port           uint
	logLevel       string
	logFile        string
	showTimestamps bool
	showLogLevel   bool
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
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		logger.InitLogger(logLevel, logFile, showTimestamps, showLogLevel)
	})
	rootCmd.PersistentFlags().UintVarP(&port, "port", "p", 5555, "Port to connect to")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (debug, info, fatal, off)")
	rootCmd.PersistentFlags().StringVar(&logFile, "log-file", "", "Log file to write to (stdout if empty)")
	rootCmd.PersistentFlags().BoolVar(&showTimestamps, "show-ts", false, "Show timestamps in log messages")
	rootCmd.PersistentFlags().BoolVar(&showLogLevel, "show-level", false, "Show log level in log messages")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(subCmd)
	rootCmd.AddCommand(pubCmd)
}
