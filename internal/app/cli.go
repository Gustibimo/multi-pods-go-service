package app

import (
	"bom-import-xls/kafkaclient"
	"context"
)

func NewCli(args []string) {
	ctx := context.Context(context.Background())
	if len(args) == 0 {
		return
	}

	if args[0] == "publish" {
		fileName := "bom.xlsx"
		kafkaclient.Publish(ctx, fileName, "file-bom-import-parsing")
	}

	if args[0] == "test" {
		// start worker
		kafkaclient.Publish(ctx, "test-hello", "test-hello")
	}

}
