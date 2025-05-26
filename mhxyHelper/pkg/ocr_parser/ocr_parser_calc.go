package ocr_parser

import "sort"

// 自动计算阈值策略
func CalculateThresholds(results []OCRResult) (yThreshold, xThreshold int) {
	// 先按Y坐标排序所有结果
	sort.Slice(results, func(i, j int) bool {
		return getMinY(results[i].Position) < getMinY(results[j].Position)
	})

	// 计算垂直间距分布
	var yDeltas []int
	for i := 1; i < len(results); i++ {
		delta := getMinY(results[i].Position) - getMinY(results[i-1].Position)
		if delta > 0 { // 忽略重叠项
			yDeltas = append(yDeltas, delta)
		}
	}
	yThreshold = calculateMedian(yDeltas) * 2

	// 初步按Y分组计算水平阈值
	rows := groupByY(results, yThreshold)
	var xDeltas []int
	for _, row := range rows {
		// 按X坐标排序
		sort.Slice(row, func(i, j int) bool {
			return getMinX(row[i].Position) < getMinX(row[j].Position)
		})

		// 计算行内间距
		for i := 1; i < len(row); i++ {
			delta := getMinX(row[i].Position) - getMaxX(row[i-1].Position)
			if delta > 0 {
				xDeltas = append(xDeltas, delta)
			}
		}
	}
	xThreshold = calculateMedian(xDeltas) * 3

	// 设置最小阈值
	if yThreshold < 10 {
		yThreshold = 10
	}
	if xThreshold < 30 {
		xThreshold = 30
	}
	return
}

// 计算中位数（自动阈值核心）
func calculateMedian(deltas []int) int {
	if len(deltas) == 0 {
		return 0
	}

	sort.Ints(deltas)
	middle := len(deltas) / 2

	if len(deltas)%2 == 0 {
		return (deltas[middle-1] + deltas[middle]) / 2
	}
	return deltas[middle]
}

func CalculateCenter(position [][]int) (float64, float64) {
	if len(position) == 0 {
		return 0, 0
	}
	sumX := 0
	sumY := 0
	count := 0
	for _, point := range position {
		if len(point) < 2 {
			continue // 跳过无效的点
		}
		sumX += point[0]
		sumY += point[1]
		count++
	}
	if count == 0 {
		return 0, 0
	}
	n := float64(count)
	centerX := float64(sumX) / n
	centerY := float64(sumY) / n
	return centerX, centerY
}
