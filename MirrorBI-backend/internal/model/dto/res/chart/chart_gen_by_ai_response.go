package chart

type ChartGenByAiResponse struct {
	ChartID   uint64 `json:"chartId,string" swaggertype:"string"` // 图表ID
	GenChart  string `json:"genChart"`                            // 生成的图表数据代码用于展示
	GenResult string `json:"genResult"`                           // 生成的图表结果
}
