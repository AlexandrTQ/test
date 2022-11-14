package cmd

import (
	"TransactionServer/database"
	"TransactionServer/server"
	"log"

	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "init database",

	Run: func(cmd *cobra.Command, args []string) {
		if err := server.StartDb(); err != nil {
			log.Fatal(err)
		}
		if err := database.GetDatabase().InitSchema("public"); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}
