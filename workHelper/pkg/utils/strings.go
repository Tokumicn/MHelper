package utils

import "errors"

// ArrGetWithCheck 检查数组长度是否允许提取字符串
func ArrGetWithCheck(strArr []string, index int) (string, error) {
	if len(strArr) >= index+1 {
		return strArr[index], nil
	}

	return "", errors.New("out of range")
}
