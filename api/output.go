package api

import (
	"fmt"
	"github.com/zuijinbuzai/fundtop/api/types"
	"github.com/zuijinbuzai/fundtop/utils"
)

func markRank(mapData []*types.Fund) {
	SortByZf(&mapData)
	for i, v := range mapData {
		if i < 3 {
			v.Rank = "***"
		}
		if i >= len(mapData) - 3 {
			v.Rank = "xxx"
		}
	}
}

func outputHeader() {
	fmt.Print("编号   ", utils.FormatName("名称"), " ")
	for i := 0; i < config.Base.ShowDays; i++ {
		if i < config.Base.ShowDays - 1 {
			fmt.Print(fmt.Sprintf("%dD    ", i))
		} else {
			fmt.Print(fmt.Sprintf("%dD   ", i))
		}
	}
	fmt.Print("| ", "涨跌       ", "涨跌2       ", "| ", "高/低  ", "当前/低 ")
	if config.Base.ShowType {
		fmt.Print(" 类型")
	}
	fmt.Println()
}

func getDDays(fi *types.Fund) string {
	if len(fi.FArray) < config.Base.ShowDays {
		return "no data          "
	}
	text := ""
	for i := 0; i < config.Base.ShowDays; i++ {
		val1 := (fi.FArray[i].Dwjz/fi.FArray[i+1].Dwjz-1) * 100
		text += utils.GetAbsText(val1)
		if i < config.Base.ShowDays - 1 {
			text += " "
		}
	}
	text += "|"
	return text
}

func outputOneLine(v *types.Fund) {
	if len(v.FArray) < 7 {
		fmt.Println(v.Code, v.Name, "no data")
		return
	}

	dayText := getDDays(v)

	txt := fmt.Sprintf("%d天(%s)", v.Delta, utils.GetAbsText(v.DeltaSum))
	txt2 := fmt.Sprintf("%02d天(%s)", v.Delta2, utils.GetAbsText(v.DeltaSum2))

	txt3 := utils.GetAbsText2(v.High/v.Low*100 - 100)
	txt4 := ""
	if len(v.FArray) >= 2 {
		txt4 = utils.GetAbsText2((v.FArray[1].Dwjz/v.Low*100 - 100) / (v.High/v.Low*100 - 100) * 100)
	}
	fmt.Print(v.Code, " ", utils.FormatName(v.Name), " ", dayText, " ", txt, " ", txt2, " ", "|", " ", txt3, "  ", txt4)

	if config.Base.ShowType {
		//fmt.Print("", v.Type)
	}
	fmt.Println()
}
