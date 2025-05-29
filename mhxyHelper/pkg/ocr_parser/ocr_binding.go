package ocr_parser

import "unicode"

func isNumber(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

// 将相邻的数字和字符串合并为一个元素
// 例如：
// [羊 鬼谷子 99999 鬼谷子 99999 多 鬼谷子 99999 超级金柳露 288888]
//
//	[{鬼谷子 99999} {鬼谷子 99999} {鬼谷子 99999} {超级金柳露 288888}]
func MergeElements(arr []string) []string {
	var result []string
	i := 0
	for i < len(arr) {
		if i+1 < len(arr) && isNumber(arr[i+1]) {
			merged := arr[i] + " " + arr[i+1]
			result = append(result, merged)
			i += 2
		} else {
			i++
		}
	}
	return result
}
