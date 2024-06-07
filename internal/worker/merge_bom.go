package worker

import (
	"bom-import-xls/internal/domain"
)

func MergeComponents(components domain.BoMComponents, costAccounts domain.BoMCostAccount) domain.BomMap {
	bomMap := make(domain.BomMap)
	for bomCode, compList := range components {
		bomMap[bomCode] = domain.BoM{
			//BomCode:      bomCode,
			RawMaterials: compList,
			CostAccount:  costAccounts[bomCode],
		}
	}

	// Include BoM entries that have cost accounts but no components
	for bomCode, costAccList := range costAccounts {
		if _, exists := bomMap[bomCode]; !exists {
			bomMap[bomCode] = domain.BoM{
				//BoMCode:      bomCode,
				RawMaterials: []domain.RawMaterials{},
				CostAccount:  costAccList,
			}
		}
	}

	return bomMap
}
