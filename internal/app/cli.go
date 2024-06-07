package app

import (
	"bom-import-xls/kafkaclient"
)

func NewCli(args []string) {
	if len(args) == 0 {
		return
	}

	if args[0] == "publish" {
		fileName := "bom.xlsx"
		kafkaclient.Publish(fileName, "file-bom-import-parsing")
	}

}
