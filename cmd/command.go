package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "root",
	Short: "cli app",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		fmt.Println("cli app")
	},
}

func init() {
	rootCmd.AddCommand(workerCmd)
	rootCmd.AddCommand(cliCmd)
	rootCmd.AddCommand(httpServerCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
