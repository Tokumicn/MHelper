package buildQA

import (
	"context"
	"fmt"
	"log/slog"
	"mhxyHelper/pkg/ocr_parser"
	strings_pipeline "mhxyHelper/pkg/string_pipeline"
	"mhxyHelper/pkg/utils"
	"os"
)

func BuildFromOCR() {
	//dir, _ := os.Getwd()
	//fmt.Println("[My-OCR] shell work dir: ", dir)

	var (
		file       *os.File
		err        error
		dirEntries []os.DirEntry
		ocrRes     []ocr_parser.OCRResult
		table      [][]ocr_parser.TableCell
	)

	ctx := context.TODO()

	//exePath, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//binDir := filepath.Dir(exePath)
	//
	//fmt.Println("[My-OCR] bin work dir: ", binDir)
	//dirPath := fmt.Sprintf("%s/ocr_images/", binDir)

	ocr_parser.SetThresholds(50, 40) // 提取到配置中

	imageDir := "./ocr_images/"
	dirEntries, err = os.ReadDir(imageDir)
	if err != nil {
		slog.ErrorContext(ctx, "read dir err: ", err.Error())
		return
	}

	for _, entry := range dirEntries {
		// 最终待写入文件结果
		var writeArr []string

		// 文件夹跳过不处理
		if entry.IsDir() {
			continue
		}

		// 检查文件后缀名是图片
		if !utils.IsImageExt(entry) {
			continue
		}

		// 将相对路径转换为绝对路径
		fileName := entry.Name()
		fullFileName := utils.ConvRelative2FullPath(imageDir, fileName)

		// 打开文件
		file, err = os.Open(fullFileName)
		if err != nil {
			slog.ErrorContext(ctx, "open file err: ", err.Error())
			return
		}

		// 发送OCR请求
		ocrRes, err = ocr_parser.PostOCR(ctx, fileName, file)
		if err != nil {
			slog.ErrorContext(ctx, "post ocr err: ", err.Error())
			return
		}

		//ocrResBytes, err := json.Marshal(ocrRes)
		//if err != nil {
		//	slog.ErrorContext(ctx, "json marshal err: ", err.Error())
		//}

		//fmt.Println("===================================================")
		//fmt.Println(string(ocrResBytes))
		//fmt.Println("================================================")

		//err = json.Unmarshal([]byte(rawJson), &ocrRes)
		//if err != nil {
		//	slog.ErrorContext(ctx, "json unmarshal err: ", err.Error())
		//	return
		//}

		// 解析OCR结果
		table, err = ocr_parser.ParseOCRToTableWithFilter(ocrRes, map[string]struct{}{
			"单价": {},
		}, []string{"等级"})

		//fmt.Println("===================================================")
		//ocr_parser.PrintOCRTable(table)
		//fmt.Println("===================================================")
		// 将结果映射到实体

		for _, row := range table {
			for _, cell := range row {
				writeArr = append(writeArr, cell.Text)
			}
		}

		// 输出实体
		err = utils.WriteLinesToFile(writeArr, "output.txt")
		if err != nil {
			slog.ErrorContext(ctx, "write file err: ", err.Error())
			return
		}

		// 初始化流水线
		pipeline := &strings_pipeline.Pipeline{}
		pipeline.AddStep(strings_pipeline.Deduplicate)
		pipeline.AddStep(strings_pipeline.NormalizeNames)

		// 解析输入数据
		items := strings_pipeline.ParseInput(writeArr)

		// 运行流水线处理
		result := pipeline.Run(items)

		// 输出结果
		for _, product := range result {
			fmt.Printf("{ Name: %q, Prices: %v }\n", product.Name, product.Prices)
		}
	}

	// 关闭文件
	if file != nil {
		err = file.Close()
		if err != nil {
			slog.ErrorContext(ctx, "close file err: ", err.Error())
			return
		}
	}
}
