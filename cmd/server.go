package cmd

import (
	"bom-import-xls/cmd/routes"
	"github.com/labstack/echo/v4"
	"github.com/spf13/cobra"
)

var httpServerCmd = &cobra.Command{
	Use:   "http-server",
	Short: "http server",
	RunE:  runHttpServer,
}

func runHttpServer(cmd *cobra.Command, args []string) error {
	// Create a new Echo instance
	e := echo.New()

	routes.SetupRoutes(e)

	// Start the HTTP server on port 8080
	return e.Start(":8080")
}
