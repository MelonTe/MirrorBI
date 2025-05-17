package controller

import (
	"mrbi/internal/common"
	"mrbi/internal/consts"
	"mrbi/internal/ecode"
	reqChart "mrbi/internal/model/dto/req/chart"
	resChart "mrbi/internal/model/dto/res/chart"
	"mrbi/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func dumb() {
	//防止引用包错误
	_ = resChart.ListChartResponse{}
}

// 接口前缀为/Chart
// param	用空格分隔的参数。param name,param type,data type,is mandatory?,comment attribute(optional)
// 获取一个Chartservice单例
var sChart *service.ChartService = service.NewChartService()

// AddChart godoc
// @Summary      添加一个图表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		request body reqChart.ChartAddRequest true "用户添加申请参数"
// @Success      200  {object}  common.Response{data=string} "添加成功，返回添加图表的ID"
// @Failure      400  {object}  common.Response "添加失败，详情见响应中的code"
// @Router       /api/chart/add [POST]
func AddChart(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var cAdd reqChart.ChartAddRequest
	if err := c.ShouldBind(&cAdd); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, err := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	if chart, err := sChart.AddChart(cAdd.Goal, cAdd.ChartData, cAdd.ChartType, loginUser.ID, "", "", "", "succeed", ""); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, strconv.FormatUint(chart.ID, 10))
		return
	}
}

// DeleteChart godoc
// @Summary      添加一个图表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		request body common.DeleteRequest true "要删除的图表的ID"
// @Success      200  {object}  common.Response{data=bool} "删除成功"
// @Failure      400  {object}  common.Response "删除失败，详情见响应中的code"
// @Router       /api/chart/delete [POST]
func DeleteChart(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var cDel common.DeleteRequest
	if err := c.ShouldBind(&cDel); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, err := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	if suc, err := sChart.DeleteChart(cDel.Id, loginUser); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, suc)
		return
	}
}

// GetChartById godoc
// @Summary      获取一个图表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		id query string true "图表的ID"
// @Success      200  {object}  common.Response{data=entity.Chart} "获取成功，返回图表数据"
// @Failure      400  {object}  common.Response "获取失败，详情见响应中的code"
// @Router       /api/chart/get [GET]
func GetChartById(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		common.BaseResponse(c, nil, "参数缺失", ecode.PARAMS_ERROR)
		return
	}
	id_parse, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	if chart, err := sChart.GetChart(id_parse, loginUser); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, chart)
		return
	}
}

// ListChartByPage godoc
// @Summary      根据页数查询图表列表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		request body reqChart.ChartQueryRequest true "需要查询的页数、以及图表关键信息"
// @Success      200  {object}  common.Response{data=resChart.ListChartResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /api/chart/list/page [POST]
func ListChartByPage(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var cQuery reqChart.ChartQueryRequest
	if err := c.ShouldBind(&cQuery); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	if charts, err := sChart.ListChartsByPage(&cQuery); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, charts)
		return
	}
}

// ListMyChartByPage godoc
// @Summary      根据页数查询图表列表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		request body reqChart.ChartQueryRequest true "需要查询的页数、以及图表关键信息"
// @Success      200  {object}  common.Response{data=resChart.ListChartResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /api/chart/list/page/my [POST]
func ListMyChartByPage(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var cQuery reqChart.ChartQueryRequest
	if err := c.ShouldBind(&cQuery); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	cQuery.UserID = loginUser.ID
	cQuery.Status = "succeed"
	if charts, err := sChart.ListChartsByPage(&cQuery); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, charts)
		return
	}
}

// ListMyChartByPage godoc
// @Summary      根据页数查询图表列表，是未成功分析的列表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		request body reqChart.ChartQueryRequest true "需要查询的页数、以及图表关键信息"
// @Success      200  {object}  common.Response{data=resChart.ListChartResponse} "查询成功"
// @Failure      400  {object}  common.Response "更新失败，详情见响应中的code"
// @Router       /api/chart/list/page/my/no [POST]
func ListMyChartByPageNo(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var cQuery reqChart.ChartQueryRequest
	if err := c.ShouldBind(&cQuery); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, _ := sUser.GetLoginUser(c)
	cQuery.UserID = loginUser.ID
	cQuery.Status = consts.ChartStatusNotSucceed
	if charts, err := sChart.ListChartsByPage(&cQuery); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, charts)
		return
	}
}

// EditChart godoc
// @Summary      编辑图表
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		request body reqChart.ChartEditRequest true "图表编辑信息"
// @Success      200  {object}  common.Response{data=bool} "编辑成功"
// @Failure      400  {object}  common.Response "编辑失败，详情见响应中的code"
// @Router       /api/chart/edit [POST]
func EditChart(c *gin.Context) {
	//使用shouldbind绑定参数，参数不可复用
	var cEdit reqChart.ChartEditRequest
	if err := c.ShouldBind(&cEdit); err != nil {
		common.BaseResponse(c, nil, "参数绑定错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, err := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	if success, err := sChart.EditChart(&cEdit, loginUser); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		common.Success(c, success)
		return
	}
}

// ChartGenByAi godoc
// @Summary      上传excel文件和目标信息，使用AI生成信息。
// @Tags         chart
// @Accept multipart/form-data
// @Produce      json
// @Param        file formData file true "excel文件"
// @Param   name      formData  string true  "图表名称"            example(人数趋势)
// @Param   goal      formData  string true  "分析目标"            example(了解用户增长)
// @Param   chartType formData  string true  "图表类型"  example(折线图)
// @Success      200  {object}  common.Response{data=resChart.ChartGenByAiResponse} "生成成功"
// @Failure      400  {object}  common.Response "生成失败，详情见响应中的code"
// @Router       /api/chart/gen/ai [POST]
func ChartGenByAi(c *gin.Context) {
	//绑定请求参数
	file, _ := c.FormFile("file")
	Name := c.PostForm("name")
	Goal := c.PostForm("goal")
	ChartType := c.PostForm("chartType")
	if file == nil {
		common.BaseResponse(c, nil, "文件不能为空", ecode.PARAMS_ERROR)
		return
	}
	//调用生成服务
	loginUser, _ := sUser.GetLoginUser(c)
	res, err := sChart.ChartGenByAi(file, Name, Goal, ChartType, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	common.Success(c, *res)
}

// GetChartData godoc
// @Summary      根据chart表的ID，获取上传的原始EXCEL的JSON格式数据
// @Tags         chart
// @Accept       json
// @Produce      json
// @Param		id query string true "图表的ID"
// @Success      200  {object}  common.Response{data=bool} "信息获取成功"
// @Failure      400  {object}  common.Response "获取失败，详情见响应中的code"
// @Router       /api/chart/data [GET]
func GetChartData(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		common.BaseResponse(c, nil, "参数缺失", ecode.PARAMS_ERROR)
		return
	}
	id_parse, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		common.BaseResponse(c, nil, "参数错误", ecode.PARAMS_ERROR)
		return
	}
	loginUser, ERR := sUser.GetLoginUser(c)
	if err != nil {
		common.BaseResponse(c, nil, ERR.Msg, ERR.Code)
		return
	}
	if chartData, err := sChart.GetChartDataById(id_parse, loginUser); err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	} else {
		data := chartData.Data
		common.Success(c, data)
		return
	}
}

// ChartGenAsyncByAi godoc
// @Summary      上传excel文件和目标信息，异步执行AI生成信息。
// @Tags         chart
// @Accept multipart/form-data
// @Produce      json
// @Param        file formData file true "excel文件"
// @Param   name      formData  string true  "图表名称"            example(人数趋势)
// @Param   goal      formData  string true  "分析目标"            example(了解用户增长)
// @Param   chartType formData  string true  "图表类型"  example(折线图)
// @Success      200  {object}  common.Response{data=string} "生成成功，返回图表的ID"
// @Failure      400  {object}  common.Response "生成失败，详情见响应中的code"
// @Router       /api/chart/gen/ai/async [POST]
func ChartGenAsyncByAi(c *gin.Context) {
	//绑定请求参数
	file, _ := c.FormFile("file")
	Name := c.PostForm("name")
	Goal := c.PostForm("goal")
	ChartType := c.PostForm("chartType")
	if file == nil {
		common.BaseResponse(c, nil, "文件不能为空", ecode.PARAMS_ERROR)
		return
	}
	//调用生成服务
	loginUser, _ := sUser.GetLoginUser(c)
	id, err := sChart.ChartGenAsyncByAi(file, Name, Goal, ChartType, loginUser)
	if err != nil {
		common.BaseResponse(c, nil, err.Msg, err.Code)
		return
	}
	idStr := strconv.FormatUint(id, 10)
	common.Success(c, idStr)
}
