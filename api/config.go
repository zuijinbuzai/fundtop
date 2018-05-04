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
	Base 		Base		`yaml:"Base"`
	FundConfig 	FundConfig	`yaml:"FundConfig"`
}

type Base struct {
	ThreadNum 	int		`yaml:"ThreadNum"`
	ShowDays	int		`yaml:"ShowDays"`
	Sort		string	`yaml:"Sort"`
}

type FundConfig struct {
	Owned		[]string	`yaml:"Owned"`
	Watched 	[]string	`yaml:"Watched"`
	Zhangfu		float64		`yaml:"Zhangfu"`
	Zhangfu2	float64		`yaml:"Zhangfu2"`
	Lhd			float64		`yaml:"Lhd"`
	Lhd2		float64		`yaml:"Lhd2"`
}

func LoadConfig() (*Config, error) {
	_, err := os.Stat(configpath)
	if err != nil && os.IsNotExist(err)  {
		newConfig()
	}

	data, err := ioutil.ReadFile("fundconfig.yaml")
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func newConfig() {
	text := `
Base:
  #下载线程数
  ThreadNum: 50
  #显示天数
  ShowDays: 5

  # 默认降序，升序加_
  #涨幅天数 Zfd Zfd2_, 涨幅 Zf,Zf_,当日涨幅，前一天 0D 1D
  Sort: 0D

FundConfig:
  #关注的关键词
  Watched: [军工, 计算机, 白酒, 消费, 人工, 医药, 汽车, 分级]

  #严格涨跌
  Zhangfu: 3.0
  #广义涨跌
  Zhangfu2: 5.0

  #自动无限翻过小山点数
  Lhd: 0.5
  #只翻过一座大山点数
  Lhd2: 1

  #翻过几天的山
  Skip: 10

  #持有的基金
  Owned: [
    002379,
    020009,
    160212,
    000457,
    001618,
    001630,
    001475,
    001629,
    001878,
    161631,
    501016,
    160625,
    001986,
    519066,
    110022,
    151725,
    160222,
    166011,
    161725
  ]
`
  ioutil.WriteFile(configpath, []byte(text), 0)
}