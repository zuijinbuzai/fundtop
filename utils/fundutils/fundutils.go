package fundutils

import (
	"io/ioutil"
	"encoding/json"
	"github.com/zuijinbuzai/fundtop/utils"
	"github.com/zuijinbuzai/fundtop/utils/httputils"
)

const (
	jsonpath = "fundtop.json"
)

var (
	fundListMap map[string] []string
)

func GetMaxMin(arr []float64) (max, min float64){
	init := false
	for _, v := range arr {
		if !init {
			max = v
			min = v
			init = true
			continue
		}
		if v > max {
			max = v
		}
		if v < min {
			min = v
		}
	}
	return max, min
}

//"000001","HXCZ","华夏成长","混合型","HUAXIACHENGZHANG"
func GetAllFundList() (map[string] []string, error) {
	if len(fundListMap) > 0 {
		return fundListMap, nil
	}

	var data []byte
	var err error
	if !utils.IsFileExist(jsonpath) {
		data, err = getFundListFromServer()
		if err != nil {
			return nil, err
		}
		ioutil.WriteFile(jsonpath, data, 0)
	} else {
		data, err = ioutil.ReadFile(jsonpath)
		if err != nil {
			return nil, err
		}
	}
	fundListMap = map[string] []string {}

	var fss [][]string
	err = json.Unmarshal(data, &fss)
	if err != nil {
		return nil, err
	}
	for _, v := range fss {
		if v[3] == "分级杠杆"{
			continue
		}
		//if v[1][len(v[1]) - 1] == 'C' {
			fundListMap[v[0]] = v
		//}
	}
	return fundListMap, err
}

func getFundListFromServer() ([]byte, error) {
	data, err := httputils.HttpGet2("http://fund.eastmoney.com/js/fundcode_search.js")
	if err != nil || len(data) < 1000 {
		return nil, err
	}
	return data[10:len(data)-1], nil
}