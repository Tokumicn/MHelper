package utils

import "errors"

const (
	DefaultRM2MHRatio = 16 // 1Yuan:16W 转换比例
)

var (
	rm2mhRatio int = DefaultRM2MHRatio
)

// 设置转换比例
func SetRatio(ratio int) {
	if ratio <= 0 {
		ratio = DefaultRM2MHRatio
	}
	rm2mhRatio = ratio
}

func MH2RM(mh float32) (float32, error) {

	if rm2mhRatio <= 0 {
		return 0, errors.New("rm2mhRatio invalid")
	}

	return mh / float32(rm2mhRatio), nil
}

func RM2MH(rm float32) (float32, error) {
	if rm2mhRatio <= 0 {
		return 0, errors.New("rm2mhRatio invalid")
	}

	return rm * float32(rm2mhRatio), nil
}
