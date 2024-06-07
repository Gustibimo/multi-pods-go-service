package worker

import (
	"bom-import-xls/internal/domain"
	"bom-import-xls/kafkaclient"
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func ParseBomFile() (domain.BomMap, error) {

	ctx, cancel := context.WithCancel(context.Background())
	//fileName := "bom.xlsx"
	//time.Sleep(10000000000)
	defer cancel()

	consumer := kafkaclient.Consumer{Ready: make(chan bool)}

	message := kafkaclient.Consume(ctx, "file-bom-import-parsing", consumer)
	messageString := strings.Trim(string(message), "\"")

	f, err := excelize.OpenFile(messageString)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	defer func(f *excelize.File) {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(f)

	rm, err := ParseComponents(f)
	if err != nil {
		return nil, err
	}
	cost, err := ParseCostAccount(f)
	bomMap := MergeComponents(rm, cost)

	return bomMap, nil
}
