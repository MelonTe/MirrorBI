package chart

type ChartGenByAiRequest struct {
	Name      string `json:"name" `      // 图表名称
	Goal      string `json:"goal" `      // 分析目标
	ChartType string `json:"chartType" ` // 图表类型
}
