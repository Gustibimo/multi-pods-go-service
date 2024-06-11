package actors

import (
	"bom-import-xls/internal/domain"
	"bom-import-xls/internal/shared"
	"errors"
	"github.com/tealeg/xlsx"
	"github.com/xuri/excelize/v2"
	"log"
)

func ReadComponents(file *xlsx.File) (domain.BoMComponents, error) {
	components := make(domain.BoMComponents)

	if len(file.Sheets) == 0 {
		return nil, errors.New("no sheets found")
	}

	sheet := file.Sheets[0]
	if len(sheet.Rows) == 0 {
		return nil, errors.New("no rows found in the sheet")
	}

	for i, row := range sheet.Rows {
		// Skip the header row
		if i == 0 {
			continue
		}

		// Ensure the row has enough columns
		if len(row.Cells) < 7 {
			log.Println("Skipping invalid row:", row)
			continue
		}

		// Read cell values and trim whitespace
		bomName := row.Cells[0].String()
		bomCode := row.Cells[1].String()

		// Directly create and assign the struct
		component := domain.RawMaterials{
			BoMName:       bomName,
			BoMCode:       bomCode,
			ComponentName: row.Cells[2].String(),
			Description:   row.Cells[3].String(),
			ProductID:     row.Cells[4].String(),
			SKU:           row.Cells[5].String(),
			Qty:           row.Cells[6].String(),
		}
		// Check if the BOM code exists, if not, initialize the slice
		if _, ok := components[bomCode]; !ok {
			components[bomCode] = []domain.RawMaterials{component}
		} else {
			components[bomCode] = append(components[bomCode], component)
		}
	}

	return components, nil
}

func ParseComponents(file *excelize.File) (domain.BoMComponents, error) {
	components := make(domain.BoMComponents)

	if len(file.GetSheetList()) == 0 {
		return nil, errors.New("no sheets found")
	}

	sheetName := file.GetSheetName(0)
	rows, err := file.Rows(sheetName)
	if err != nil {
		return nil, errors.New("no rows found in the sheet")
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			return nil, err
		}
		// Ensure the row has enough columns
		if len(row) < 7 {
			//fmt.Println("Skipping invalid row:", row)
			continue
		}

		if shared.RowContainsHeader(row) {
			continue
		}

		bomName := row[0]
		bomCode := row[1]
		component := domain.RawMaterials{
			BoMName:       bomName,
			BoMCode:       bomCode,
			ComponentName: row[2],
			Description:   row[3],
			ProductID:     row[4],
			SKU:           row[5],
			Qty:           row[6],
		}

		// Check if the BOM code exists, if not, initialize the slice
		if _, ok := components[bomCode]; !ok {
			components[bomCode] = []domain.RawMaterials{component}
		} else {
			// Append to existing slice
			components[bomCode] = append(components[bomCode], component)
		}
	}

	return components, nil
}
