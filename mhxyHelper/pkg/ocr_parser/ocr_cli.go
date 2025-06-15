package ocr_parser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
)

func PostOCR(ctx context.Context, fileName string, file *os.File) ([]OCRResult, error) {
	ocrRawRes, err := ocrCall(ctx, fileName, file)
	if err != nil {
		return nil, err
	}

	resp := &OCRResponse{}
	err = json.Unmarshal(ocrRawRes, resp)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("ocr call err: %d", resp.StatusCode)
	}

	return resp.Results, nil
}

func ocrCall(ctx context.Context, fileName string, file *os.File) ([]byte, error) {
	url := "http://0.0.0.0:8501/ocr" // TODO: 提取配置
	method := http.MethodPost

	// 2. 准备请求体
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 创建表单文件字段
	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		panic("创建表单字段失败: " + err.Error())
	}

	// 将文件内容拷贝到表单字段
	_, err = io.Copy(part, file)
	if err != nil {
		panic("写入文件数据失败: " + err.Error())
	}

	// 关闭multipart writer（自动添加结尾boundary）
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(ctx context.Context) {
		err := res.Body.Close()
		if err != nil {
			slog.ErrorContext(ctx, "close res body err: ", err.Error())
		}
	}(ctx)

	ResponseBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ResponseBody, nil
}
