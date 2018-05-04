package api

import (
	"zuijinbuzai/fundtop/api/types"
	"math"
	"strings"
)

func FilterResult(fi *types.Fund) bool {
	if math.Abs(fi.DeltaSum) >= config.FundConfig.Zhangfu {
		return true
	}
	if math.Abs(fi.DeltaSum2) >= config.FundConfig.Zhangfu2 {
		return true
	}
	return false
}

func FilterOwned(fi *types.Fund) bool {
	for _, v := range config.FundConfig.Owned {
		if fi.Code == v {
			return true
		}
	}
	return false
}

func FilterNeedWork(fund *types.Fund) bool {
	if fund.Type == "分级杠杆" || fund.Type == "固定收益" {
		return false
	}

	for _, v := range config.FundConfig.Owned {
		if fund.Code == v || fund.Name == v {
			return true
		}
	}
	for _, v := range config.FundConfig.Watched {
		if strings.Index(fund.Name, v) >= 0 {
			return true
		}
	}
	return false
}