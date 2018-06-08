package api

import (
	"github.com/zuijinbuzai/fundtop/api/types"
	"strings"
)

func FilterResult(fi *types.Fund) bool {
	//if math.Abs(fi.DeltaSum) >= config.FundConfig.Zhangfu {
	//	return true
	//}
	//if math.Abs(fi.DeltaSum2) >= config.FundConfig.Zhangfu2 {
	//	return true
	//}
	//if len(fi.FArray) < 7 {
	//	return false
	//}
	//
	//ok := false
	//if math.Abs(fi.FArray[0].Dwjz/fi.FArray[1].Dwjz*100-100) >= config.FilterResultConfig.D0_ &&
	//	math.Abs(fi.FArray[1].Dwjz/fi.FArray[2].Dwjz*100-100) >= config.FilterResultConfig.D1_ &&
	//		math.Abs(fi.FArray[2].Dwjz/fi.FArray[3].Dwjz*100-100) >= config.FilterResultConfig.D2_ &&
	//			math.Abs(fi.DeltaSum) >= config.FilterResultConfig.Zhangfu_ &&
	//				math.Abs(fi.DeltaSum2) >= config.FilterResultConfig.Zhangfu2_ &&
	//					math.Abs(fi.High/fi.Low*100 - 100) >= config.FilterResultConfig.High2Low_ {
	//		ok = true
	//}
	//if !ok {
	//	return false
	//}
	//if math.Abs(fi.FArray[0].Dwjz/fi.FArray[1].Dwjz*100-100) >= config.FilterResultConfig.D0 ||
	//	math.Abs(fi.FArray[1].Dwjz/fi.FArray[2].Dwjz*100-100) >= config.FilterResultConfig.D1 ||
	//	math.Abs(fi.FArray[2].Dwjz/fi.FArray[3].Dwjz*100-100) >= config.FilterResultConfig.D2 ||
	//	math.Abs(fi.DeltaSum) >= config.FilterResultConfig.Zhangfu ||
	//	math.Abs(fi.DeltaSum2) >= config.FilterResultConfig.Zhangfu2 ||
	//	math.Abs(fi.High/fi.Low*100 - 100) >= config.FilterResultConfig.High2Low {
	//	ok = true
	//}
	//if ok {
	//	return true
	//}
	//return false
	return true
}

func FilterOwned(fi *types.Fund) bool {
	return config.OwnedFundMap[fi.Code] != nil
}

func FilterNeedWork(fund *types.Fund) bool {
	if fund.Type == "分级杠杆" || fund.Type == "固定收益" {
		return false
	}

	for _, v := range config.FundConfig.Black {
		if fund.Code == v {
			return false
		}
	}

	if config.FilterResultConfig.JustC && fund.Name[len(fund.Name) - 1] != 'C' {
		return false;
	}

	if FilterOwned(fund) {
		return true
	}
	for _, v := range config.FundConfig.Watched {
		if strings.Index(fund.Name, v) >= 0 {
			return true
		}
	}
	return false
}