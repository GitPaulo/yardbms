package main

import (
	"fmt"
	"os"
	"yardbms/repl"

	"github.com/spf13/cobra"
)

var storageType string

var rootCmd = &cobra.Command{
	Use:   "yardbms",
	Short: "A simple RDBMS",
	Long:  `yardbms is a simple relational database management system.`,
	Run: func(cmd *cobra.Command, args []string) {
		repl.Start(storageType)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&storageType, "storage", "s", "ram", "Type of storage to use: 'ram' or 'file'")
}
