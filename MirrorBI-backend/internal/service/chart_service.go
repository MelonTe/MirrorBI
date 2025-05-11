package service

import (
	"mrbi/internal/common"
	"mrbi/internal/consts"
	"mrbi/internal/ecode"
	reqChart "mrbi/internal/model/dto/req/chart"
	resChart "mrbi/internal/model/dto/res/chart"
	"mrbi/internal/model/entity"
	"mrbi/internal/repository"
	"mrbi/pkg/db"

	"gorm.io/gorm"
)

type ChartService struct {
	ChartRepo *repository.ChartRepository
}

func NewChartService() *ChartService {
	return &ChartService{
		ChartRepo: repository.NewChartRepository(),
	}
}

func (s *ChartService) AddChart(goal string, chartData string, chartType string, userId uint64) (uint64, *ecode.ErrorWithCode) {
	// 1.校验
	if chartData == "" || chartType == "" {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数为空")
	}
	if len(chartData) < 4 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图表数据过短")
	}
	if len(chartType) < 4 {
		return 0, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图表类型过短")
	}

	// 2.添加图表
	chart := &entity.Chart{
		UserID:    userId,
		ChartData: chartData,
		ChartType: chartType,
		Goal:      goal,
	}
	id, err := s.ChartRepo.AddChart(nil, chart)
	if err != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	return id, nil
}

func (s *ChartService) DeleteChart(chartId uint64, loginUser *entity.User) (bool, *ecode.ErrorWithCode) {
	// 1.校验
	if chartId == 0 {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数错误")
	}
	// 2.校验图表是否存在
	chart, err := s.ChartRepo.FindById(nil, chartId)
	if err != nil {
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	if chart == nil {
		return false, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图表不存在")
	}
	if chart.UserID != loginUser.ID && loginUser.UserRole != consts.ADMIN_ROLE {
		return false, ecode.GetErrWithDetail(ecode.FORBIDDEN_ERROR, "没有权限删除该图表")
	}
	// 3.删除图表
	err = s.ChartRepo.DeleteChart(nil, chartId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图表不存在")
		}
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	return true, nil
}

func (s *ChartService) GetChart(id uint64, loginUser *entity.User) (*entity.Chart, *ecode.ErrorWithCode) {
	//获取图表
	chart, err := s.ChartRepo.FindById(nil, id)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	if chart == nil {
		return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图表不存在")
	}
	return chart, nil
}

func (s *ChartService) GetQueryWrapper(db *gorm.DB, req *reqChart.ChartQueryRequest) (*gorm.DB, *ecode.ErrorWithCode) {
	query := db.Session(&gorm.Session{})
	if req == nil {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数为空")
	}
	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.ChartType != "" {
		query = query.Where("chart_type = ?", req.ChartType)
	}
	if req.ChartData != "" {
		query = query.Where("chart_data LIKE ?", "%"+req.ChartData+"%")
	}
	if req.Goal != "" {
		query = query.Where("goal = ?", req.Goal)
	}
	return query, nil
}

// 批量获取图表
func (s *ChartService) ListChartsByPage(req *reqChart.ChartQueryRequest) (*resChart.ListChartResponse, *ecode.ErrorWithCode) {
	//参数校验
	if req.Current <= 0 || req.PageSize <= 0 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "页数或者页面大小不能小于0")
	}
	//获取查询对象
	query, err := s.GetQueryWrapper(db.LoadDB(), req)
	if err != nil {
		return nil, err
	}
	//查询总数
	var total int64
	query.Model(&entity.Chart{}).Count(&total)
	to := int(total)
	//分页查询
	var Charts []entity.Chart
	//重置query
	query, _ = s.GetQueryWrapper(db.LoadDB(), req)
	query = query.Offset((req.Current - 1) * req.PageSize).Limit(req.PageSize)
	query.Find(&Charts)
	p := (to + req.PageSize - 1) / req.PageSize
	//获取VO对象
	list := &resChart.ListChartResponse{
		Records: Charts,
		PageResponse: common.PageResponse{
			Total:   to,
			Size:    req.PageSize,
			Pages:   p,
			Current: req.Current,
		},
	}
	return list, nil
}

func (s *ChartService) EditChart(request *reqChart.ChartEditRequest, loginUser *entity.User) (bool, *ecode.ErrorWithCode) {
	//参数校验
	if request.ID == 0 {
		return false, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数错误")
	}
	//获取图表
	chart, err := s.ChartRepo.FindById(nil, request.ID)
	if err != nil {
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	if chart == nil {
		return false, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图表不存在")
	}
	if chart.UserID != loginUser.ID && loginUser.UserRole != consts.ADMIN_ROLE {
		return false, ecode.GetErrWithDetail(ecode.FORBIDDEN_ERROR, "没有权限编辑该图表")
	}
	//更新图表
	updateMap := map[string]interface{}{
		"chart_data": request.ChartData,
		"chart_type": request.ChartType,
		"goal":       request.Goal,
	}
	if err := s.ChartRepo.UpdateChartByMap(nil, request.ID, updateMap); err != nil {
		return false, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	return true, nil
}
