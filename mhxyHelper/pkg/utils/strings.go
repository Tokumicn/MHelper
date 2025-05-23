package utils

import (
	"errors"
	"strings"
)

// ArrGetWithCheck 检查数组长度是否允许提取字符串
func ArrGetWithCheck(strArr []string, index int) (string, error) {
	if len(strArr) >= index+1 {
		return strArr[index], nil
	}

	return "", errors.New("out of range")
}

// IsMultiline 判断字符串是否包含多行
func IsMultiline(s string) bool {
	// 检查字符串中是否包含换行符
	return strings.Contains(s, "\n") || strings.Contains(s, "\r\n")
}
