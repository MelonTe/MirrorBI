package chart

type ChartAddRequest struct {
	Goal      string `json:"goal"`      // 目标
	ChartData string `json:"chartData"` // 图表数据
	ChartType string `json:"chartType"` // 图表类型
}
