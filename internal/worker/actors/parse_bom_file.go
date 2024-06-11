package actors

import (
	"bom-import-xls/internal/domain"
	"bom-import-xls/kafkaclient"
	"context"
	"fmt"
	"github.com/xuri/excelize/v2"
	"strings"
)

func ParseBomFile(ctx context.Context, message string) (domain.BomMap, error) {

	fmt.Println("Consumer trace id: ", ctx.Value("trace_id"))
	_, cancel := context.WithCancel(context.Background())
	//time.Sleep(10000000000)
	defer cancel()

	fileName := strings.Trim(message, "\"")

	f, err := excelize.OpenFile(fileName)
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
	if err != nil {
		return nil, err
	}
	bomMap := MergeComponents(rm, cost)

	// send file-parsing-success or file-parsing-failed to kafka
	kafkaclient.Publish(ctx, bomMap, "file-parsing-succeed")

	return bomMap, nil
}
