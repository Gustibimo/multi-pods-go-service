package app

import (
	"bom-import-xls/internal/domain"
	"encoding/json"
	"fmt"
)

func SaveBoM(data string) {
	// Save BoM to database
	var bomMap domain.BomMap

	// unmarshal data string to bomMap
	err := json.Unmarshal([]byte(data), &bomMap)
	if err != nil {
		fmt.Println("Error unmarshalling data string to bomMap")
		return
	}

	for bomCode, bom := range bomMap {
		fmt.Printf("save Raw Materials with bomCode %v: %v\n", bomCode, bom.RawMaterials)
		if bom.CostAccount != nil {
			fmt.Printf("save Cost Account with bomCode %v: %v\n", bomCode, bom.CostAccount)
		}
	}
}
