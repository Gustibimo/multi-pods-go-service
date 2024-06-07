package cmd

import (
	"bom-import-xls/internal/app"
	"fmt"
	"github.com/spf13/cobra"
)

// cli command
var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "cli app",
	RunE:  runCli,
}

func runCli(cmd *cobra.Command, args []string) error {
	fmt.Println("[cli app]")
	err := cmd.Flags().Parse(args)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		return nil
	}

	app.NewCli(args)
	return nil
}
