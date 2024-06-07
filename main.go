package main

import (
	"bom-import-xls/cmd"
)

func main() {

	cmd.Execute()

	//fileName := "bom.xlsx"
	// Open the Excel file
	//file, err := xlsx.OpenFile("bom.xlsx")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//kafkaclient.Publish(fileName, "file-bom-import-parsing")

	// time to parse the file
	//start := time.Now()
	//consumer, err := worker.BomFileParsingActor()
	//if err != nil {
	//	fmt.Println("Error parsing file: ", err)
	//}
	//fmt.Println(consumer)

	//r, err := worker.ParseBomFile()
	//if err != nil {
	//	return
	//}
	//fmt.Println(r)
	//
	//elapsed := time.Since(start)
	//fmt.Println("Time elapsed: ", elapsed)

}
