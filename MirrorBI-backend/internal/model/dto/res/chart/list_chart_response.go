package chart

import (
	"mrbi/internal/common"
	"mrbi/internal/model/entity"
)

type ListChartResponse struct {
	common.PageResponse
	Records []entity.Chart `json:"records"` // 图表列表
}
