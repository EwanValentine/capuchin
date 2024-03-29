package cmd

import (
	"log"

	"github.com/EwanValentine/capuchin/conf"
	"github.com/EwanValentine/capuchin/http"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use: "http",
	Run: func(cmd *cobra.Command, args []string) {
		c := conf.Load()
		if err := http.NewServer(c); err != nil {
			log.Fatal(err)
		}
	},
}
