package ocr_parser

type OCRResult struct {
	Text     string  `json:"text"`
	Score    float64 `json:"score"`
	Position [][]int `json:"position"`
}

// 定义OCR结果结构体
type OCRResponse struct {
	StatusCode int         `json:"status_code"`
	Results    []OCRResult `json:"results"`
}
