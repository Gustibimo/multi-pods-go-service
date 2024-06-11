package handlers

import (
	"bom-import-xls/internal/app"
	"bom-import-xls/internal/worker/actors"
	"context"
	"fmt"
)

func HandleParseBomFile(ctx context.Context, message string) {
	_, err := actors.ParseBomFile(ctx, message)
	if err != nil {
		panic(err)
	}

	//fmt.Println(file)
}

func HandleHello(ctx context.Context) {
	fmt.Println("Hello")
}

func HandleFileParseSucceed(ctx context.Context, message string) {
	fmt.Println("file-parse-succeed with trace id: ", ctx.Value("trace_id"))
	app.SaveBoM(message)
}

func HandleFileParseFailed(ctx context.Context) {
	fmt.Println("file-parse-failed")
}
