package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(grpcServer)
}

var grpcServer = &cobra.Command{
	Use: "grpc",
	Run: func(cmd *cobra.Command, args []string) {
		// server := grpc.NewServer()
	},
}
