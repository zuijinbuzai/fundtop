package api

import (
	"time"
	"fmt"
	"strings"
	"sync"
	"encoding/json"
	"github.com/zuijinbuzai/fundtop/api/types"
	"github.com/zuijinbuzai/fundtop/utils/fundutils"
	"github.com/zuijinbuzai/fundtop/utils/httputils"
	"github.com/zuijinbuzai/fundtop/utils"
)

/**
	https://github.com/weibycn/fund

	获取50项数据
	http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=160220&page=1&per=50
 */
var (
	fundDetailURL	= "http://fund.eastmoney.com/f10/F10DataApi.aspx?type=lsjz&code=%06s&page=1&per=50"
	fundGzURL 		= "http://fundgz.1234567.com.cn/js/%s.js"

	lock  *sync.Mutex
)

var (
	config 		*Config
)

func Work() {
	arrayData, _ := getAllFundData()
	if arrayData == nil {
		return
	}

	fmt.Print("编号   ", utils.FormatName("名称"), " ")
	for i := 0; i < config.Base.ShowDays; i++ {
		fmt.Print(fmt.Sprintf("%dD     ", i))
	}
	fmt.Print("| ", "涨跌       ", "涨跌2  ")
	fmt.Println()

	markRank(arrayData)
	SortList(&arrayData)

	for _, v := range arrayData {
		if FilterResult(v) && !FilterOwned(v) {
			outputOneLine(v)
		}
	}
	fmt.Println("---------------------------------------------------------------------------------------------------持有", len(config.FundConfig.Owned), "支")
	for _, v := range arrayData {
		if FilterOwned(v) {
			outputOneLine(v)
		}
	}
	fmt.Println("---------------------------------------------------------------------------------------------------持有", len(config.FundConfig.Owned), "支")
}

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

func outputOneLine(v *types.Fund) {
	dayText := getDDays(v)
	txt := fmt.Sprintf("%d天(%s)", v.Delta, utils.GetAbsText(v.DeltaSum))
	txt2 := fmt.Sprintf("%02d天(%s)", v.Delta2, utils.GetAbsText(v.DeltaSum2))
	fmt.Println(v.Code, utils.FormatName(v.Name), dayText, txt, txt2, v.Rank)
}

func getDDays(fi *types.Fund) string {
	if len(fi.FArray) < config.Base.ShowDays {
		return "no data          "
	}
	text := ""
	for i := 0; i < config.Base.ShowDays; i++ {
		val1 := (fi.FArray[i].Dwjz/fi.FArray[i+1].Dwjz-1) * 100
		text += utils.GetAbsText(val1)
		text += "  "
	}
	text += "|"
	return text
}

func getAllFundData() ([]*types.Fund, error) {
	var err error
	config, err = LoadConfig()
	if err != nil {
		fmt.Println("加载配置失败", err.Error())
		return  nil, nil
	}
	list, err := fundutils.GetAllFundList()
	if err != nil {
		fmt.Println("读取基金列表失败")
		return nil, nil
	}

	DBOpen()
	defer DBClose()

	dh := &downloadHeap{
		newDownloads:	make(chan struct{}, 1),
	}
	result := []*types.Fund {}
	for k, v := range list {
		fd := &types.Fund{}
		if fd == nil {
			continue
		}
		fd.Code = k
		fd.Name = v[2]
		fd.Type = v[3]
		if FilterNeedWork(fd) {
			dh.managedPush(fd)
			result = append(result, fd)
		}
	}

	lock = &sync.Mutex{}
	num := config.Base.ThreadNum
	dlCount := 0

	t := time.Now()
	wg := new(sync.WaitGroup)
	for i := 0; i < num; i++ {
		wg.Add(1)
		go threadDownFund(dh, wg, &dlCount)
	}
	wg.Wait()
	DBPutArray(result)
	fmt.Printf("\rDownloadAllFund finish, %d支 use = %f s\n", dlCount, time.Now().Sub(t).Seconds())

	return result, nil
}

func getOneFund(code string) *types.FundItem {
	data, err := httputils.HttpGet2(fmt.Sprintf(fundGzURL, code))
	if err != nil || len(data) < 20 {
		return nil
	}
	fi := &types.FundItem{}
	//text := string(data[8:len(data)-2])
	//fmt.Println(text)
	err = json.Unmarshal(data[8:len(data)-2], fi)
	if err != nil {
		//fmt.Println(err.Error())
		return nil
	}
	return fi
}

func threadDownFund(dh *downloadHeap, wg *sync.WaitGroup, dlCount *int) {
	for {
		fd := dh.managedPop()
		if fd != nil {
			fi := getOneFund(fd.Code)
			if fi == nil {
				continue
			}
			fd.FArray = append(fd.FArray, fi)

			arr := DBGet(fd.Code)
			if arr != nil {
				fd.FArray = append(fd.FArray, *arr...)
			} else {
				fd.FArray = append(fd.FArray, parseHtmlFundDetail(fd.Code)...)
			}

			Analyze(fd)
			lock.Lock()

			*dlCount++
			if *dlCount % 50 == 0 {
				fmt.Printf("\rdlCount=%d", *dlCount)
			}
			lock.Unlock()
		} else {
			break
		}
	}
	wg.Done()
}

func parseHtmlFundDetail(code string) ([]*types.FundItem) {
	url := fmt.Sprintf(fundDetailURL, code)
	data, err := httputils.HttpGet2(url)
	if err != nil {
		return nil
	}
	text := string(data)

	pos := strings.Index(text, "<tbody>")
	if pos == -1 {
		return nil
	}
	text = text[pos + len("<tbody>"):]
	fundArr := []*types.FundItem{}

	//解析table
	for {
		pos = strings.Index(text, "<tr>")
		if pos == -1 {
			break
		}
		text = text[pos + len("<tr>"):]

		fi := &types.FundItem{}
		for i := 0; i < 4; i++ {
			if strings.Index(text, "<td") < 0 {
				break
			}
			pos = strings.Index(text, ">")
			text = text[pos + len(">"):]

			pos = strings.Index(text, "</td>")
			if pos < 4 {
				break
			}
			text2 := text[:pos]
			text = text[pos + len("</td>"):]

			switch i {
			case 0:
				//fi.Time = text2
			case 1:
				fi.Dwjz = utils.ParseFloat(text2)
			case 2:
				//fi.Ljjz = fundutils.ParseFloat(text2)
			case 3:
				//fi.Rzzl = fundutils.ParseFloat(text2)
			}
		}
		fundArr = append(fundArr, fi)
	}
	return fundArr
}