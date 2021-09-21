package cmd

import (
	"log"

	"github.com/EwanValentine/capuchin/conf"
	"github.com/EwanValentine/capuchin/grpc"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(grpcServer)
}

var grpcServer = &cobra.Command{
	Use: "grpc",
	Run: func(cmd *cobra.Command, args []string) {
		c := conf.Load()
		if err := grpc.NewServer(c); err != nil {
			log.Fatal(err)
		}
	},
}
