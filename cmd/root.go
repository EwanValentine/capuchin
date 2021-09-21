package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "capuchin",
	Short: "Capuchin distributed query engine",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Capuchin...")
	},
}

// Execute command line interface
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
