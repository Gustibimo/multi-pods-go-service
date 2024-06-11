package actors

import (
	"bom-import-xls/internal/domain"
	"bom-import-xls/internal/shared"
	_ "bom-import-xls/kafkaclient"
	"github.com/tealeg/xlsx"
	"log"
	"sync"
	"time"
)

type Result struct {
	Success any
	Error   error
}

func BomFileParsingActor() (domain.BomMap, error) {
	// Open the Excel file
	start := time.Now()
	file, err := xlsx.OpenFileWithRowLimit("bom.xlsx", xlsx.NoRowLimit)
	if err != nil {
		log.Fatal(err)
	}

	elapsed := time.Since(start)
	log.Println("Time elapsed merge: ", elapsed)

	shared.CleanExcelFile(file)

	var wg sync.WaitGroup
	resultCh := make(chan Result, 2)

	wg.Add(2)

	//compCh := make(chan domain.BoMComponents)

	go func() {
		defer wg.Done()
		components, err := ReadComponents(file)
		if err != nil {
			resultCh <- Result{
				Error: err,
			}
			return
		}
		resultCh <- Result{
			Success: components,
		}
	}()

	//costCh := make(chan domain.BoMCostAccount)
	// read the cost account
	go func() {
		defer wg.Done()
		cost, err := ReadCostAccounts(file)
		if err != nil {
			resultCh <- Result{
				Error: err,
			}
			return
		}
		resultCh <- Result{
			Success: cost,
		}
	}()

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	//rawMaterials := <-compCh
	//costAccounts := <-costCh

	var rawMaterials domain.BoMComponents
	var costAccounts domain.BoMCostAccount

	for i := 0; i < 2; i++ {
		result := <-resultCh
		if result.Error != nil {
			return nil, result.Error
		}
		switch res := result.Success.(type) {
		case domain.BoMComponents:
			rawMaterials = res
		case domain.BoMCostAccount:
			costAccounts = res
		}
	}

	bomMap := MergeComponents(rawMaterials, costAccounts)

	//for bomCode, comp := range bomMap {
	//	// publish bom
	//	fmt.Println(bomCode, comp)
	//}

	return bomMap, nil
}
