package shared

func RowContainsHeader(row []string) bool {
	return row[0] == "Nama BoM*" || row[0] == "Kode BoM*"
}
