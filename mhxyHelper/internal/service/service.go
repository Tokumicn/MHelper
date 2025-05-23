package service

import (
	"context"
	"fmt"
	"mhxyHelper/internal/data/const_val"
	"mhxyHelper/internal/service/query/local_query"
	"mhxyHelper/internal/service/query/mhjl_query"
	"mhxyHelper/pkg/logger"
)

const (
	TypeUnknown   = "Unknown"
	TypeStuff     = "Stuff"     // 物品信息
	TypeAttribute = "Attribute" // 物品属性
	TypeMHJL      = "MHJL"      // 梦幻精灵
)

// QueryLocal 查询信息 本地数据库检索
func QueryLocal(inStr string) (int64, string, interface{}, error) {

	_, ok := const_val.QueryQNameMapStuff[inStr]
	if ok {
		total, stuffs, err := local_query.QueryStuff(inStr)
		if err != nil {
			logger.Log.Error("Query QueryStuff err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeStuff, stuffs, nil
	}

	_, ok = const_val.QueryNameMapStuff[inStr]
	if ok {
		total, stuffs, err := local_query.QueryStuff(inStr)
		if err != nil {
			logger.Log.Error("Query QueryStuff err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeStuff, stuffs, nil
	}

	_, ok = const_val.QueryQNameMapAttribute[inStr]
	if ok {
		total, attributes, err := local_query.QueryAttribute(inStr)
		if err != nil {
			logger.Log.Error("Query QueryAttribute err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeAttribute, attributes, nil
	}

	_, ok = const_val.QueryQNameMapAttribute[inStr]
	if ok {
		total, attributes, err := local_query.QueryAttribute(inStr)
		if err != nil {
			logger.Log.Error("Query QueryAttribute err:%v", err)
			return 0, "", nil, err
		}
		return total, TypeAttribute, attributes, nil
	}

	return 0, "", nil, fmt.Errorf("input query string[%s] no match ", inStr)
}

// QueryMHJL for  查询梦幻精灵
func QueryMHJL(ctx context.Context, query string) (string, error) {

	return mhjl_query.QueryMHJL(ctx, query)
}
