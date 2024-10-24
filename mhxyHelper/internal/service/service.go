package service

import (
	"fmt"
	"mhxyHelper/pkg/common"
	"mhxyHelper/pkg/logger"
)

const (
	TypeUnknown   = "Unknown"
	TypeStuff     = "Stuff"
	TypeAttribute = "Attribute"
)

// 查询信息
func Query(inStr string) (int64, string, interface{}, error) {

	_, ok := common.QueryQNameMapStuff[inStr]
	if ok {
		total, stuffs, err := QueryStuff(inStr)
		if err != nil {
			logger.Log.Error("Query QueryStuff err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeStuff, stuffs, nil
	}

	_, ok = common.QueryNameMapStuff[inStr]
	if ok {
		total, stuffs, err := QueryStuff(inStr)
		if err != nil {
			logger.Log.Error("Query QueryStuff err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeStuff, stuffs, nil
	}

	_, ok = common.QueryQNameMapAttribute[inStr]
	if ok {
		total, attributes, err := QueryAttribute(inStr)
		if err != nil {
			logger.Log.Error("Query QueryAttribute err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeAttribute, attributes, nil
	}

	_, ok = common.QueryQNameMapAttribute[inStr]
	if ok {
		total, attributes, err := QueryAttribute(inStr)
		if err != nil {
			logger.Log.Error("Query QueryAttribute err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeAttribute, attributes, nil
	}

	return 0, "", nil, fmt.Errorf("input query string[%s] no match ", inStr)
}
