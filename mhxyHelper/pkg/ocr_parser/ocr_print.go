package ocr_parser

import (
	"fmt"
	"github.com/gosuri/uitable"
	"strings"
)

func PrintOCRTable(table [][]TableCell) {
	t := uitable.New()
	for _, row := range table {
		tempColumns := make([]string, 0)
		for _, cell := range row {
			splits := strings.Split(cell.Text, " ")
			tempColumns = append(tempColumns, splits...)
		}
		t.AddRow(tempColumns)
	}
	fmt.Println(t.String())
}
