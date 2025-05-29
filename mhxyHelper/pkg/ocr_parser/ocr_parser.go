package ocr_parser

import (
	"sort"
	"strings"
)

// 分组参数配置
var (
	//yThreshold = 15 // 垂直方向分组阈值
	//xThreshold = 50 // 水平方向合并阈值

	yThreshold = 100 // 垂直方向分组阈值
	xThreshold = 400 // 水平方向合并阈值 // TODO: 最重要的调节参数，考虑如何自动计算这个参数
)

func SetThresholds(y, x int) {
	yThreshold = y
	xThreshold = x
}

// 表格单元格结构体
type TableCell struct {
	Text     string  // 单元格文本
	Position [][]int // 位置坐标
}

// ParseOCRToTable 解析OCR结果为表格
func ParseOCRToTable(results []OCRResult) ([][]TableCell, error) {

	// 步骤1: 按Y坐标分组
	rows := groupByY(results, yThreshold)

	// 步骤2: 对每行按X坐标排序
	table := make([][]TableCell, 0, len(rows))
	for _, row := range rows {
		sort.Slice(row, func(i, j int) bool {
			return getMinX(row[i].Position) < getMinX(row[j].Position)
		})

		// 步骤3: 合并同行相邻元素
		mergedRow := mergeAdjacentCells(row, xThreshold)
		table = append(table, mergedRow)
	}

	return table, nil
}

// ParseOCRToTableWithFilter 解析OCR结果为表格，支持过滤文本
func ParseOCRToTableWithFilter(results []OCRResult, ignoreTextMap map[string]struct{}, prefixIgnore []string) ([][]TableCell, error) {
	// 过滤掉无用信息
	results = filter4Result(results, ignoreTextMap, prefixIgnore)
	return ParseOCRToTable(results)
}

func filter4Result(results []OCRResult, ignoreTextMap map[string]struct{}, prefixIgnore []string) []OCRResult {
	res := make([]OCRResult, 0)
	for _, result := range results {
		if _, ok := ignoreTextMap[result.Text]; ok {
			continue
		}

		for _, prefix := range prefixIgnore {
			if !strings.HasPrefix(result.Text, prefix) {
				res = append(res, result)
			}
		}
	}
	return res
}

// 按Y坐标分组函数
func groupByY(results []OCRResult, threshold int) [][]OCRResult {
	sort.Slice(results, func(i, j int) bool {
		return getMinY(results[i].Position) < getMinY(results[j].Position)
	})

	var groups [][]OCRResult
	var currentGroup []OCRResult
	var lastY int

	for _, res := range results {
		currentY := getMinY(res.Position)

		if len(currentGroup) == 0 {
			currentGroup = append(currentGroup, res)
			lastY = currentY
			continue
		}

		if abs(currentY-lastY) <= threshold {
			currentGroup = append(currentGroup, res)
		} else {
			groups = append(groups, currentGroup)
			currentGroup = []OCRResult{res}
			lastY = currentY
		}
	}

	if len(currentGroup) > 0 {
		groups = append(groups, currentGroup)
	}
	return groups
}

// 合并同行相邻元素（处理可能的识别碎片）
func mergeAdjacentCells(row []OCRResult, threshold int) []TableCell {
	var merged []TableCell
	var current *TableCell

	for _, res := range row {
		minX := getMinX(res.Position)
		//maxX := getMaxX(res.Position)

		if current == nil {
			current = &TableCell{
				Text:     res.Text,
				Position: res.Position,
			}
			continue
		}

		lastMaxX := getMaxX(current.Position)
		if minX-lastMaxX < threshold {
			// 合并文本和位置
			current.Text += " " + res.Text
			current.Position = [][]int{
				current.Position[0], // 保留最左坐标
				res.Position[1],     // 新的右上角
				res.Position[2],     // 新的右下角
				current.Position[3], // 保留最左下方坐标
			}
		} else {
			merged = append(merged, *current)
			current = &TableCell{
				Text:     res.Text,
				Position: res.Position,
			}
		}
	}

	if current != nil {
		merged = append(merged, *current)
	}
	return merged
}

// 辅助函数：获取坐标边界值
func getMinY(pos [][]int) int {
	if len(pos) == 0 {
		return 0
	}
	min := pos[0][1]
	for _, p := range pos {
		if p[1] < min {
			min = p[1]
		}
	}
	return min
}

func getMinX(pos [][]int) int {
	if len(pos) == 0 {
		return 0
	}
	min := pos[0][0]
	for _, p := range pos {
		if p[0] < min {
			min = p[0]
		}
	}
	return min
}

func getMaxX(pos [][]int) int {
	if len(pos) == 0 {
		return 0
	}
	max := pos[0][0]
	for _, p := range pos {
		if p[0] > max {
			max = p[0]
		}
	}
	return max
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
