package cmd

import (
	"log"

	"github.com/EwanValentine/capuchin/conf"
	"github.com/EwanValentine/capuchin/grpc"
	"github.com/EwanValentine/capuchin/http"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use: "start",
	Run: func(cmd *cobra.Command, args []string) {
		c := conf.Load()

		// Start gRPC server
		go func() {
			if err := grpc.NewServer(c); err != nil {
				log.Fatal(err)
			}
		}()

		// Start HTTP server (proxies gRPC server)
		if err := http.NewServer(c); err != nil {
			log.Fatal(err)
		}
	},
}
