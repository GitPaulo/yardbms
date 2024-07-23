package main

import (
	"fmt"
	"os"
	"yardbms/repl"

	"github.com/spf13/cobra"
)

var storageType string
var filePath string
var rootCmd = &cobra.Command{
	Use:   "yardbms",
	Short: "A simple RDBMS",
	Long:  `yardbms is a simple relational database management system.`,
	Run: func(_ *cobra.Command, args []string) {
		repl.Start(storageType, filePath)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&storageType, "storage", "s", "ram", "Type of storage to use: 'ram' or 'file'")
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "database.json", "File path for file storage (only used if storage type is 'file')")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
