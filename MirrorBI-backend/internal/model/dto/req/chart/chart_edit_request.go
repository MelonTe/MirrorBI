package chart

type ChartEditRequest struct {
	ID        uint64 `json:"id,string" swaggertype:"string"` // 图表ID
	Goal      string `json:"goal"`                           // 目标
	ChartData string `json:"chartData"`                      // 图表数据
	ChartType string `json:"chartType"`                      // 图表类型
}
