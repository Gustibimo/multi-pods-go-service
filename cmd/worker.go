package cmd

import (
	"bom-import-xls/internal/worker"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "worker app",
	Run:   runWorker(),
}

func runWorker() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// Create a channel to listen for termination signals
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

		// Create a channel to indicate when the worker has finished
		done := make(chan struct{})

		// Start the worker
		fmt.Println("[worker app]")
		go func() {
			worker.Run()
			close(done)
		}()

		// Wait for termination signal or worker completion
		select {
		case <-stopChan:
			fmt.Println("Shutting down gracefully...")
			// Signal the worker to stop
			worker.Stop()
			worker.Wait()
		case <-done:
			fmt.Println("worker finished its task.")
		}

	}
}
