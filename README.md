## 爬取基金数据，做一些简单分析

## 用法
1.releases目录下有exe可执行文件，直接下载。
2.在一个单独的目录下运行exe，运行后会生成
- fundtop.json(基金列表数据)不用管
- fundtop.yaml (配置文件) 程序如何执行有此文件指定
- fundtop.db (基金历史净值缓存数据)不用管

## 配置

~~~
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
  Watched: [军工, 白酒]

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
  Owned: [002379, 020009]
~~~

## 输出展示
![](https://github.com/zuijinbuzai/fundtop/blob/master/img/output.jpg)
