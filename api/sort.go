package api

import (
	"github.com/bradfitz/slice"
	"github.com/zuijinbuzai/fundtop/api/types"
)

func SortByZf(data *[]*types.Fund) {
	mapData := *data
	slice.Sort(mapData, func(i, j int) bool {
		return mapData[i].DeltaSum > mapData[j].DeltaSum
	})
}

func SortList(data *[]*types.Fund) {
	mapData := *data
	slice.Sort(mapData, func(i, j int) bool {
		switch config.Base.Sort {
		case "Zfd2":
			return mapData[i].Delta2 > mapData[j].Delta2
		case "Zfd":
			return mapData[i].Delta > mapData[j].Delta
		case "Zf2":
			return mapData[i].DeltaSum2 > mapData[j].DeltaSum2
		case "Zf":
			return mapData[i].DeltaSum > mapData[j].DeltaSum
		case "Zfd2_":
			return mapData[i].Delta2 < mapData[j].Delta2
		case "Zfd_":
			return mapData[i].Delta < mapData[j].Delta
		case "Zf2_":
			return mapData[i].DeltaSum2 < mapData[j].DeltaSum2
		case "Zf_":
			return mapData[i].DeltaSum < mapData[j].DeltaSum

		case "0D_":
			return getPercent(mapData, i, 0) < getPercent(mapData, j, 0)
		case "0D":
			return getPercent(mapData, i, 0) > getPercent(mapData, j, 0)
		case "1D_":
			return getPercent(mapData, i, 1) < getPercent(mapData, j, 1)
		case "1D":
			return getPercent(mapData, i, 1) > getPercent(mapData, j, 1)
		default:
			return mapData[i].DeltaSum < mapData[j].DeltaSum
		}
	})
}

func getPercent(mapData []*types.Fund, i int, offset int) float64 {
	if len(mapData[i].FArray) <= 2 {
		return 0
	}
	return mapData[i].FArray[offset].Dwjz/mapData[i].FArray[offset+1].Dwjz * 100 - 100
}