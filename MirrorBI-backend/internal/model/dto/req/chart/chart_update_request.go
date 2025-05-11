package chart

import "time"

type ChartUpdateRequest struct {
	ID         uint64    `json:"id,string" swaggertype:"string"` // 图表ID
	Goal       string    `json:"goal,omitempty"`                 // 目标
	ChartData  string    `json:"chartData,omitempty"`            // 图表数据
	ChartType  string    `json:"chartType,omitempty"`            // 图表类型
	GenChart   string    `json:"genChart,omitempty"`
	GenResult  string    `json:"genResult,omitempty" `
	CreateTime time.Time `json:"createTime,omitempty" `
	UpdateTime time.Time `json:"updateTime,omitempty" `
}
