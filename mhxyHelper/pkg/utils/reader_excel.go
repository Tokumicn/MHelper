package utils

import (
	"bufio"
	"context"
	"github.com/xuri/excelize/v2"
	"mhxyHelper/pkg/logger"
	"os"
)

// 从 Excel 中读取Attribute信息
func ReadAttributeFromExcel(ctx context.Context) ([][]string, error) {

	res := make([][]string, 0)

	f, err := excelize.OpenFile("./config/信息整理.xlsx")
	if err != nil {
		return nil, err
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			logger.Log.ErrorContext(ctx, "ReadAttributeFromExcel defer close excel file err: %v", err)
		}
	}()

	tempRows, err := f.Rows("属性类")
	if err != nil {
		return nil, err
	}

	tempRows.Next() // 丢弃第一行

	for tempRows.Next() {
		clos, err := tempRows.Columns()
		if err != nil {
			return nil, err
		}

		res = append(res, clos)
	}
	return res, nil
}

// 从 CSV 文件中读取 Stuff 数据
func ReadCSVFromStuffData() ([]string, error) {
	res := make([]string, 0)
	f, err := os.Open("/Users/zhangrui/Workspace/goSpace/src/Tokumicn/theBookofChangesEveryDay/mhxyHelper/config/stuff_data.csv")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	r.ReadLine() // 丢弃第一行
	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		res = append(res, string(line))
	}
	return res, nil
}
