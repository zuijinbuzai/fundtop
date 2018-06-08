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
	config 			*Config
)

func Work() {
	arrayData, _ := getAllFundData()
	if arrayData == nil {
		return
	}

	outputHeader()
	markRank(arrayData)
	SortList(&arrayData)

	count := 0
	for _, v := range arrayData {
		if FilterResult(v) && !FilterOwned(v) {
			if count < config.Base.Limit {
				outputOneLine(v)
				count++
			}
		}
	}
	//fmt.Println("----------------------------------------------------------------------------------------------------------------------")
	//for i := len(arrayData) - 1; i >= len(arrayData) - config.Base.Limit && i >= 0; i-- {
	//	v := arrayData[i];
	//	if FilterResult(v) && !FilterOwned(v) {
	//		if count < config.Base.Limit {
	//			outputOneLine(v)
	//			count++
	//		}
	//	}
	//}

	fmt.Println("----------------------------------------------------------------------------------------------------------------------")
	for _, v := range arrayData {
		if FilterOwned(v) {
			outputOneLine(v)
		}
	}
	msg := fmt.Sprintf("----------------------------------------------------" +
		" 持有%d只", len(config.OwnedFundMap))
	fmt.Println(msg)
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
	fmt.Printf("\rDownloadAllFund finish, %d支 use = %f s,  %s\n", dlCount, time.Now().Sub(t).Seconds(), time.Now().Format("2006-01-02 15:04:05"))

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
			if *dlCount % 100 == 0 {
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
	first := true
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
				//2018-06-01
				timeLayout := "2006-01-02"
				loc, _ := time.LoadLocation("Local")
				theTime, _ := time.ParseInLocation(timeLayout, text2, loc)
				delt := time.Now().Unix() - theTime.Unix()
				//10天没更新
				if first && delt > 60 * 60 * 24 * 10 {
					//fmt.Println(code, theTime.Year(), theTime.Month(), theTime.Day())
					return nil
				}
				first = false;
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