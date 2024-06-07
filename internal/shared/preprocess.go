package shared

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
)

func CleanExcelFile(file *xlsx.File) {
	for _, sheet := range file.Sheets {
		var cleanedRows []*xlsx.Row
		i := len(sheet.Rows)
		fmt.Println(i)
		for _, row := range sheet.Rows {
			isEmpty := true
			for _, cell := range row.Cells {
				if strings.TrimSpace(cell.String()) != "" {
					isEmpty = false
					break
				}
			}
			if !isEmpty {
				cleanedRows = append(cleanedRows, row)
			}
		}
		sheet.Rows = cleanedRows
	}
}
