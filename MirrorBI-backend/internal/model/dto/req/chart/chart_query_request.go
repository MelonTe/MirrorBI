package chart

import "mrbi/internal/common"

type ChartQueryRequest struct {
	common.PageRequest
	UserID    uint64 `json:"userId,string" swaggertype:"string"` // 用户Id
	Name      string `json:"name"`                               // 图表名称
	Goal      string `json:"goal"`                               // 目标
	ChartData string `json:"chartData"`                          // 图表数据
	ChartType string `json:"chartType"`                          // 图表类型
}
