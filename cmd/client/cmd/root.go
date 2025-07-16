package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ServerURL string
)

var rootCmd = &cobra.Command{
	Use:   "todocli",
	Short: "A CLI client for the TODO gRPC service",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	serverHost := os.Getenv("SERVER_HOST")
	if serverHost == "" {
		serverHost = "localhost"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	ServerURL = fmt.Sprintf("http://%s:%s", serverHost, serverPort)
}
