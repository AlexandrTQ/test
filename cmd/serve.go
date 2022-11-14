package cmd

import (
	"TransactionServer/server"
	"log"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "execute server",

	Run: func(cmd *cobra.Command, args []string) {
		if err := server.StartDb(); err != nil {
			log.Fatal(err)
		}
		if err := server.StartServer(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
