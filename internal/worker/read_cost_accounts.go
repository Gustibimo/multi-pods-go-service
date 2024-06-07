package worker

import (
	"bom-import-xls/internal/domain"
	"bom-import-xls/internal/shared"
	"errors"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"log"
)

func ReadCostAccounts(file *xlsx.File) (domain.BoMCostAccount, error) {
	costAccounts := make(map[string][]domain.CostAccount)
	sheet := file.Sheets[1]
	for i, row := range sheet.Rows {
		// Skip the header row
		if i == 0 {
			continue
		}

		// Ensure the row has enough columns
		if len(row.Cells) < 3 {
			log.Println("Skipping invalid row:", row)
			continue
		}

		// Read cell values and trim whitespace
		bomCode := row.Cells[0].String()
		name := row.Cells[1].String()
		amount := row.Cells[2].String()

		costAccounts[bomCode] = append(costAccounts[bomCode], domain.CostAccount{
			BoMCode: bomCode,
			Name:    name,
			Amount:  amount,
		})
	}

	return costAccounts, nil
}

func ParseCostAccount(file *excelize.File) (domain.BoMCostAccount, error) {
	costAccounts := make(domain.BoMCostAccount)

	if len(file.GetSheetList()) == 0 {
		return nil, errors.New("no sheet found")
	}

	sheetName := file.GetSheetName(1)
	rows, err := file.Rows(sheetName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, errors.New("no columns")
		}

		if len(row) < 3 {
			continue
		}

		// delete the header
		if shared.RowContainsHeader(row) {
			continue
		}

		// Read cell values and trim whitespace
		bomCode := row[0]
		name := row[1]
		amount := row[2]

		costAccounts[bomCode] = append(costAccounts[bomCode], domain.CostAccount{
			BoMCode: bomCode,
			Name:    name,
			Amount:  amount,
		})
	}
	return costAccounts, nil
}
