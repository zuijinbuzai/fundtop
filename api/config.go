package api

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	configpath = "fundtop.yaml"
)

type Config struct {
	Base				Base					`yaml:"Base"`
	FundConfig			FundConfig				`yaml:"FundConfig"`
	FilterResultConfig	FilterResultConfig		`yaml:"FilterResultConfig"`
	OwnedFundMap		map[string]*OwnedFund	`yaml:"OwnedFund"`
}

type Base struct {
	ThreadNum 	int		`yaml:"ThreadNum"`
	ShowDays	int		`yaml:"ShowDays"`
	Sort		string	`yaml:"Sort"`
	ShowType	bool	`yaml:"ShowType"`
	Limit		int		`yaml:"Limit"`
}

type FundConfig struct {
	//Owned		[]string	`yaml:"Owned"`
	Watched 	[]string	`yaml:"Watched"`
	Black		[]string	`yaml:"Black"`
	Zhangfu		float64		`yaml:"Zhangfu"`
	Zhangfu2	float64		`yaml:"Zhangfu2"`
	Lhd			float64		`yaml:"Lhd"`
	Lhd2		float64		`yaml:"Lhd2"`
}

type FilterResultConfig struct {
	//D0				float64		`yaml:"0D"`
	//D1				float64		`yaml:"1D"`
	//D2				float64		`yaml:"2D"`
	//Zhangfu			float64		`yaml:"Zhangfu"`
	//Zhangfu2		float64		`yaml:"Zhangfu2"`
	//High2Low		float64		`yaml:"Highh2Low"`
	//Curr2Low		float64		`yaml:"Curr2Low"`
	//
	//D0_				float64		`yaml:"0D_"`
	//D1_				float64		`yaml:"1D_"`
	//D2_				float64		`yaml:"2D_"`
	//Zhangfu_		float64		`yaml:"Zhangfu_"`
	//Zhangfu2_		float64		`yaml:"Zhangfu2_"`
	//High2Low_		float64		`yaml:"Highh2Low_"`
	//Curr2Low_		float64		`yaml:"Curr2Low_"`

	JustC			bool		`yaml:"JustC"`
}

type OwnedFund struct {
	Name			string	`yaml:"name"`
	Code			string	`yaml:"code"`
}

func LoadConfig() (*Config, error) {
	_, err := os.Stat(configpath)
	if err != nil && os.IsNotExist(err)  {
		newConfig()
	}

	data, err := ioutil.ReadFile(configpath)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	//for _, v := range config.OwnedFund {
	//	config.OwnedFundMap[v.Code] = v
	//}
	return config, nil
}

func newConfig() {
	text := `
Base:
  #下载线程数
  ThreadNum: 100
  #显示天数
  ShowDays: 8
  #展示基金类型
  ShowType: false

  #默认降序，升序加_
  #涨幅天数 Zfd Zfd2_, 涨幅 Zf,Zf_,当日涨幅，前一天 0D 1D
  Sort: 0D_
  #显示条数限制
  Limit: 15

FundConfig:
  #关注的关键词
  Watched: [铁, 天, 智能, 汽车, 军工, 医药, 白酒, 体育, 煤, 碳, 石油, 计算机]
  Black: [003942, 001483, 003603, 003604, 003975, 003974]

  #自动无限翻过小山点数
  Lhd: 0.5
  #只翻过一座大山点数
  Lhd2: 2
  #翻过几天的山
  Skip: 10

FilterResultConfig:
  #只显示带C的基金
  JustC: false

#持有基金
OwnedFund:
  001986: {name: 前海人工智能}
  002251: {name: 华夏军工}
  001630: {name: 天弘计算机C}
  001618: {name: 天弘电子C}
  000457: {name: 上投摩根}
  160212: {name: 国泰估值}
`
  ioutil.WriteFile(configpath, []byte(text), 0)
}