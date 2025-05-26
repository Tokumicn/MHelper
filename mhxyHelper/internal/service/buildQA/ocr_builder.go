package buildQA

import (
	"context"
	"fmt"
	"log/slog"
	"mhxyHelper/pkg/ocr_parser"
	"mhxyHelper/pkg/utils"
	"os"
)

const (
	imageDir = "./images"
)

func BuildFromOCR() {
	//dir, _ := os.Getwd()
	//fmt.Println("[My-OCR] shell work dir: ", dir)

	//exePath, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//binDir := filepath.Dir(exePath)
	//
	//fmt.Println("[My-OCR] bin work dir: ", binDir)
	//dirPath := fmt.Sprintf("%s/images/", binDir)
	// 读取文件

	var (
		ctx        = context.TODO()
		file       *os.File
		err        error
		dirEntries []os.DirEntry
		ocrRes     []ocr_parser.OCRResult
		table      [][]ocr_parser.TableCell
	)

	dirEntries, err = os.ReadDir(imageDir)
	if err != nil {
		slog.ErrorContext(ctx, "read dir err: ", err.Error())
		return
	}

	for _, entry := range dirEntries {

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

		// 解析OCR结果
		table, err = ocr_parser.ParseOCRToTable(ocrRes)
		if err != nil {
			slog.ErrorContext(ctx, "parse ocr err: ", err.Error())
			return
		}

		ocr_parser.PrintOCRTable(table)
		fmt.Println("===================================================")
		fmt.Println("===================================================")
		fmt.Println("===================================================")

		// TODO 将结果映射到实体

		// TODO 输出实体
	}

	if file != nil {
		err = file.Close()
		if err != nil {
			slog.ErrorContext(ctx, "close file err: ", err.Error())
		}
	}

}
