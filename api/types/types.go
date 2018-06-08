package types

type Fund struct {
	Code 		string
	Name		string
	Type		string

	Low			float64
	High		float64

	Delta		int
	DeltaSum	float64

	Delta2		int
	DeltaSum2	float64

	FArray		[]*FundItem

	Rank		string

	Earnings	float64
}

type FundItem struct {
	//净值日期
	//Time	string		`json:"gztime"`
	//单位净值
	Dwjz	float64		`json:"gsz,string"`
	//累计净值
	//Ljjz	float64
	//日增长率
	//Rzzl	float64
}