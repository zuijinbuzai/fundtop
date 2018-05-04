package api

import (
	"github.com/zuijinbuzai/fundtop/api/types"
	"math"
)

func Analyze(fi *types.Fund) {
	AnalyzeHighAndLow(fi)
	AnalyzeSum(fi)
}

func AnalyzeHighAndLow(fi *types.Fund) {
	init := false

	for i := 0; i < len(fi.FArray); i++  {
		v := fi.FArray[i]
		if !init {
			fi.Low = v.Dwjz
			init = true
		}
		if v.Dwjz > fi.High {
			fi.High = v.Dwjz
		}
		if v.Dwjz < fi.Low {
			fi.Low = v.Dwjz
		}
	}
}

func AnalyzeGetSumBegin(fi *types.Fund) int {
	for i := 1; i < len(fi.FArray) - 1; i++  {
		compare := fi.FArray[i].Dwjz - fi.FArray[i+1].Dwjz
		if compare != 0{
			return i
		}
	}
	return len(fi.FArray) - 1
}

func AnalyzeSum(fi *types.Fund) {
	beg := AnalyzeGetSumBegin(fi)
	if beg == len(fi.FArray) - 1 {
		return
	}
	compare := fi.FArray[beg].Dwjz - fi.FArray[beg + 1].Dwjz

	i := beg
	for ; i < len(fi.FArray) - 1; i++ {
		dv := (fi.FArray[i].Dwjz - fi.FArray[i+1].Dwjz)/fi.FArray[i+1].Dwjz * 100
		if dv * compare < 0 {
			break
		}
	}
	fi.Delta = i - 1
	fi.DeltaSum = (fi.FArray[1].Dwjz - fi.FArray[i].Dwjz) * 100 / fi.FArray[i].Dwjz

	goForward := 0
	i = beg
	for ; i < len(fi.FArray) - 1; i++  {
		dv := (fi.FArray[i].Dwjz - fi.FArray[i+1].Dwjz)/fi.FArray[i+1].Dwjz * 100
		if dv * compare >= 0 {
		} else {
			//翻过小山
			ret := AnalyzeSkip(fi, &i, compare, config.FundConfig.Lhd)
			if ret {
				continue
			}

			if goForward >= 1 {
				break
			}

			//翻过大山
			ret = AnalyzeSkip(fi, &i, compare, config.FundConfig.Lhd2)
			if ret {
				goForward++
				continue
			}
			break
		}
	}

	fi.Delta2 = i - 1
	fi.DeltaSum2 = (fi.FArray[1].Dwjz - fi.FArray[i].Dwjz) * 100 / fi.FArray[i].Dwjz
}

func AnalyzeSkip(fi *types.Fund, i *int, compare float64, cv float64) bool {
	for k := *i+1; k < len(fi.FArray) - 2 && k < *i+5; k++ {
		dv := (fi.FArray[*i].Dwjz - fi.FArray[k].Dwjz)/fi.FArray[k].Dwjz * 100
		if math.Abs(dv) <= cv && *i < len(fi.FArray) - 1 {
			dv2 := (fi.FArray[*i].Dwjz - fi.FArray[k+1].Dwjz)/fi.FArray[k+1].Dwjz * 100
			if dv2 * compare >= 0 {
				*i = k
				return true
			}
		}
	}
	return false
}