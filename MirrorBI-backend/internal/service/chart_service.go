package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/panjf2000/ants/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"io"
	"log"
	"mime/multipart"
	"mrbi/internal/api/siliconflow"
	"mrbi/internal/common"
	"mrbi/internal/consts"
	"mrbi/internal/ecode"
	reqChart "mrbi/internal/model/dto/req/chart"
	resChart "mrbi/internal/model/dto/res/chart"
	"mrbi/internal/model/entity"
	"mrbi/internal/repository"
	"mrbi/internal/utils"
	"mrbi/pkg/db"
	"mrbi/pkg/mq"
	"mrbi/pkg/rds"
	"strconv"
	"strings"
	"sync"
)

type ChartService struct {
	ChartRepo *repository.ChartRepository
}

func NewChartService() *ChartService {
	return &ChartService{
		ChartRepo: repository.NewChartRepository(),
	}
}

// AI生成的协程池，最大支持4个AI任务并发，支持20个任务排队等待
var aiGenPool *ants.Pool
var aiGenPoolOnce sync.Once

// 后台协程，负责处理AI任务
func ChartBackgroundService() {
	//获取协程池
	aiPool := GetAiGenPool()
	//连接一个消息队列
	ch := mq.GetChannel()
	//最后归还连接（大概不会触发）
	defer mq.ReleaseChannel(ch)
	//获取消息的通道
	msgs, err := ch.Consume(
		consts.MQQueueName,       // 队列名称
		consts.ChartConsumerName, // 消费者名称
		false,                    // 是否自动确认
		false,                    // 是否排他
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panic("Failed to register a consumer")
	}
	type ChartError struct {
		Id  uint64
		Err error
	}
	//错误消息处理管道
	errChan := make(chan ChartError, 24) //顶多24个任务
	//创建错误处理协程
	go func() {
		for chanerr := range errChan {
			//获取图表
			chart, err := NewChartService().ChartRepo.FindById(nil, chanerr.Id)
			if err != nil {
				log.Println("保存图表错误失败,数据库报错:", err)
				continue
			}
			if chart == nil {
				log.Println("图表不存在，直接跳过")
				continue
			}
			//修改chart状态为失败
			chart.Status = consts.ChartStatusFailed
			chart.ExecMessage = chanerr.Err.Error()
			updateMap := map[string]interface{}{
				"status":       chart.Status,
				"exec_message": chart.ExecMessage,
			}
			err = NewChartService().ChartRepo.UpdateChartByMap(nil, chart.ID, updateMap)
			if err != nil {
				log.Println("保存图表错误失败,数据库报错:", err)
				continue
			}
		}
	}()
	//创建获取消息的协程
	go func() {
		//提交任务，进行异步处理
		for d := range msgs {
			taskProcessErr := aiPool.Submit(func() {
				s := NewChartService()
				//从消息中，获取图表的ID
				chartId, _ := strconv.ParseUint(string(d.Body), 10, 64)
				//获取图表
				chart, err := NewChartService().ChartRepo.FindById(nil, chartId)
				if err != nil {
					//此处出现错误，必定是数据库操作错误,将消息重放回列队
					log.Println("获取图表失败,数据库报错:", err)
					d.Nack(false, true)
					return
				}
				if chart == nil {
					//图表不存在了，直接处理下一条消息。消息可以丢弃
					d.Ack(false)
					return
				}
				//需要校验图片是否已被处理，若处理了，跳到下一条消息。
				if chart.Status == consts.ChartStatusSucceed {
					d.Ack(false)
					return
				}
				//若消息正在被处理
				if chart.Status == consts.ChartStatusRunning {
					//说明本来要处理这个消息的协程，断开了和RabbitMQ的连接
					//这个消息现在变成了不确定的状态，它可能会被那个协程正确处理，也可能会出现业务错误，并且得不到回应。
					//如何处理？为了降低业务的复杂性，选择将这个Chart的状态设置为失败。
					//交给后台协程处理
					errChan <- ChartError{
						Id:  chart.ID,
						Err: errors.New("因系统原因，任务执行失败"),
					}
					d.Ack(false)
					return
				}
				//修改chart状态为执行中
				chart.Status = consts.ChartStatusRunning
				chart.ExecMessage = "正在执行"
				updateMap := map[string]interface{}{
					"status":       chart.Status,
					"exec_message": chart.ExecMessage,
				}
				err = s.ChartRepo.UpdateChartByMap(nil, chart.ID, updateMap)
				if err != nil {
					//更新失败，输出错误
					errChan <- ChartError{
						Id:  chart.ID,
						Err: errors.New("系统错误，更新图表失败"),
					}
					d.Ack(false)
					return
				}
				//开始处理任务
				//构造AI调用请求参数
				userRequirement := fmt.Sprintf("分析需求:%s", chart.Goal)
				if chart.ChartType != "" {
					userRequirement += fmt.Sprintf(",图表类型:%s", chart.ChartType)
				}
				//获取图标数据Data
				chartData, err := s.ChartRepo.FindChartDataByChartDataId(nil, chart.ChartDataID)
				if err != nil {
					//获取失败，输出错误
					errChan <- ChartError{
						Id:  chart.ID,
						Err: errors.New("系统错误，获取图表数据失败"),
					}
					d.Ack(false)
					return
				}
				data := string(chartData.Data)
				//调用API
				res, err := siliconflow.NewLLMChatReqeustNoContext(userRequirement, data)
				if err != nil {
					//调用失败，输出错误
					errChan <- ChartError{
						Id:  chart.ID,
						Err: errors.New("系统错误，调用AI服务失败"),
					}
					log.Println("调用AI服务失败:", err)
					d.Ack(false)
					return
				}
				//提取res中的数据
				genChart, genResult, err := s.GetGenResultAndChart(res.Choices[0].Message.Content)
				if err != nil {
					//提取失败，输出错误
					errChan <- ChartError{
						Id:  chart.ID,
						Err: errors.New("系统错误，提取AI结果失败"),
					}
					log.Println("调用AI服务失败:", err)
					d.Ack(false)
					return
				}
				//保存状态
				chart.GenChart = genChart
				chart.GenResult = genResult
				updateMap = map[string]interface{}{
					"status":       consts.ChartStatusSucceed,
					"exec_message": "执行成功",
					"gen_chart":    chart.GenChart,
					"gen_result":   chart.GenResult,
				}
				//存储数据库
				err = s.ChartRepo.UpdateChartByMap(nil, chart.ID, updateMap)
				if err != nil {
					//更新失败，返回错误
					errChan <- ChartError{
						Id:  chart.ID,
						Err: errors.New("系统错误，更新图表失败"),
					}
					d.Ack(false)
					return
				}
				//通知消息队列
				d.Ack(false)
			})
			//进行错误处理
			if taskProcessErr != nil {
				log.Println("任务处理失败:", taskProcessErr)
				//将消息重放回列队
				d.Nack(false, true)
				continue
			}
		}
	}()
}

func GetAiGenPool() *ants.Pool {
	aiGenPoolOnce.Do(func() {
		var err error
		aiGenPool, err = ants.NewPool(4,
			ants.WithMaxBlockingTasks(20),
			ants.WithPreAlloc(true),
			ants.WithNonblocking(true))
		if err != nil {
			panic(fmt.Sprintf("创建协程池失败: %v", err))
		}
	})
	return aiGenPool
}

// 根据excel文件、生成图表名称、目标和类型，返回生成结果（异步执行）
func (s *ChartService) ChartGenAsyncByAi(excelFile *multipart.FileHeader, name, goal, chartType string, loginUser *entity.User) (uint64, *ecode.ErrorWithCode) {
	//先尝试获取令牌
	acquire, originErr := rds.GetAIRateLimiter().Allow(context.Background(), "AI-Chart-Gen", 1)
	if originErr != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "获取令牌失败")
	}
	if !acquire {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "请求过于频繁,请稍后再试")
	}
	//文件存放至本地，需要校验文件不能＞20MB
	dst, originErr := utils.SaveFileToLocal(excelFile)
	if originErr != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "文件保存失败")
	} else {
		defer utils.DeleteFile(dst)
	}
	data, originErr := utils.ExcelToCSV(dst)
	if originErr != nil {
		return 0, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, originErr.Error())
	}

	//构建图表入库,状态为等待执行
	chart, ERR := s.AddChart(goal, data, chartType, loginUser.ID, "", "", name, consts.ChartStatusWait, "")
	if ERR != nil {
		return 0, ERR
	}
	//发送给消息队列
	mq.GetChannelPool().PublishMessage([]byte(fmt.Sprintf("%d", chart.ID)))
	//提前返回任务ID
	return chart.ID, nil
}

// 添加一个图表，会保存CSV数据
func (s *ChartService) AddChart(goal string, chartData string, chartType string, userId uint64, genChart, genResult string, name string, status string, execmessage string) (*entity.Chart, *ecode.ErrorWithCode) {
	// 1.校验
	if chartData == "" || chartType == "" {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数为空")
	}
	if len(chartData) < 4 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图表数据过短")
	}
	if len(chartType) < 4 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "图表类型过短")
	}

	// 2.添加图表，开启事务
	tx := s.ChartRepo.BeginTransaction()
	//数据转化为JSON
	jsonData, err := s.ConvertChartDataToJSON(chartData)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//构造图表对象
	chartJsonData := &entity.ChartDataJSON{
		Data: datatypes.JSON([]byte(jsonData)),
	}
	chart := &entity.Chart{
		UserID:      userId,
		ChartType:   chartType,
		Goal:        goal,
		GenChart:    genChart,
		GenResult:   genResult,
		Name:        name,
		Status:      status,
		ExecMessage: execmessage,
	}
	//先存储chartData
	chartDataId, originErr := s.ChartRepo.AddChartDataJSON(tx, chartJsonData)
	if originErr != nil {
		tx.Rollback()
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	//设置Chart的chartDataId
	chart.ChartDataID = chartDataId
	_, originErr = s.ChartRepo.AddChart(nil, chart)
	if originErr != nil {
		tx.Rollback()
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	//提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	return chart, nil
}

// 将ChartData转化为JSON数据
func (s *ChartService) ConvertChartDataToJSON(chartData string) (string, *ecode.ErrorWithCode) {
	// 1.校验
	if chartData == "" {
		return "", ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "表为空")
	}
	// 2.转换
	r := csv.NewReader(strings.NewReader(chartData))
	// 如果你的 CSV 里用的是其他分隔符，可选： r.Comma = '\t'

	// 1. 读表头
	headers, err := r.Read()
	if err == io.EOF {
		return "[]", nil
	}
	if err != nil {
		return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "读取表头错误")
	}

	// 2. 逐行读取
	var records []map[string]interface{}
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "读取数据行错误")
		}
		// 每行要和表头一样长，不足的补空字符串，多余的忽略
		m := make(map[string]interface{}, len(headers))
		for i, key := range headers {
			var v interface{} = ""
			if i < len(row) {
				cell := row[i]
				// 尝试转换为整数
				if n, err := strconv.ParseInt(cell, 10, 64); err == nil {
					v = n
				} else if f, err := strconv.ParseFloat(cell, 64); err == nil {
					v = f
				} else {
					v = cell
				}
			}
			m[key] = v
		}
		records = append(records, m)
	}

	// 3. 序列化成 JSON
	out, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return "", ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "序列化数据成 JSON 错误")
	}
	return string(out), nil
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
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
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
	if req.Status != "" && req.Status != consts.ChartStatusNotSucceed {
		query = query.Where("status = ?", req.Status)
	} else if req.Status == consts.ChartStatusNotSucceed {
		query = query.Where("status != ?", consts.ChartStatusSucceed)
	}
	return query, nil
}

// 批量获取图表
func (s *ChartService) ListChartsByPage(req *reqChart.ChartQueryRequest) (*resChart.ListChartResponse, *ecode.ErrorWithCode) {
	//参数校验
	if req.Current <= 0 || req.PageSize <= 0 {
		//自动设置，默认为第一页，10条
		if req.Current <= 0 {
			req.Current = 1
		}
		if req.PageSize <= 0 || req.PageSize > 20 {
			req.PageSize = 10
		}
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

// 根据excel文件、生成图表名称、目标和类型，返回生成结果（同步执行）
func (s *ChartService) ChartGenByAi(excelFile *multipart.FileHeader, name, goal, chartType string, loginUser *entity.User) (*resChart.ChartGenByAiResponse, *ecode.ErrorWithCode) {
	//先尝试获取令牌
	acquire, originErr := rds.GetAIRateLimiter().Allow(context.Background(), "AI-Chart-Gen", 1)
	if originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "获取令牌失败")
	}
	if !acquire {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "请求过于频繁,请稍后再试")
	}
	//文件存放至本地，需要校验文件不能＞20MB
	dst, originErr := utils.SaveFileToLocal(excelFile)
	if originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "文件保存失败")
	} else {
		defer utils.DeleteFile(dst)
	}
	data, originErr := utils.ExcelToCSV(dst)
	if originErr != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, originErr.Error())
	}

	//构造AI调用请求参数
	userRequirement := fmt.Sprintf("分析需求:%s", goal)
	if chartType != "" {
		userRequirement += fmt.Sprintf(",图表类型:%s", chartType)
	}

	//调用API
	res, err := siliconflow.NewLLMChatReqeustNoContext(userRequirement, data)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, err.Error())
	}

	//提取res中的数据
	genChart, genResult, err := s.GetGenResultAndChart(res.Choices[0].Message.Content)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, err.Error())
	}
	//构建图表入库
	chart, ERR := s.AddChart(goal, data, chartType, loginUser.ID, genChart, genResult, name, consts.ChartStatusSucceed, "")
	if ERR != nil {
		return nil, ERR
	}
	return &resChart.ChartGenByAiResponse{
		GenChart:  genChart,
		GenResult: genResult,
		ChartID:   chart.ID,
	}, nil
}

// 从AI生成的文本中提取图表数据和结果，若未能成功提取返回错误
func (s *ChartService) GetGenResultAndChart(text string) (optionPart string, analysisPart string, err error) {
	// 查找 "option"
	idx := strings.Index(text, "option")
	if idx == -1 {
		err = errors.New("'option' 关键字未找到")
		return
	}
	// 查找 '='
	eqIdx := strings.Index(text[idx:], "=")
	if eqIdx == -1 {
		err = errors.New("'=' 未在 'option' 之后找到")
		return
	}
	// 查找第一个 '{'
	braceStart := strings.Index(text[idx+eqIdx:], "{")
	if braceStart == -1 {
		err = errors.New("未找到 opening '{'")
		return
	}
	// 计算绝对开始位置
	start := idx + eqIdx + braceStart
	// 查找匹配的 '}'
	depth := 0
	end := -1
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				end = i
				break
			}
		}
	}
	if end == -1 {
		err = errors.New("未找到匹配的 closing '}'")
		return
	}
	optionPart = strings.TrimSpace(text[start : end+1])

	// 查找分析标题（兼容两种说法）
	var headerIdx int
	if i := strings.Index(text, "数据结论分析"); i != -1 {
		headerIdx = i
	} else if i := strings.Index(text, "数据分析结论"); i != -1 {
		headerIdx = i
	} else {
		err = errors.New("未找到分析标题")
		return
	}
	// 查找 ':' 之后的内容
	colonIdx := strings.Index(text[headerIdx:], ":")
	if colonIdx == -1 {
		err = errors.New("分析标题之后未找到 ':'")
		return
	}
	analysisPart = strings.TrimSpace(text[headerIdx+colonIdx+1:])
	return
}

// 根据Chart表的I获取图表数据
func (s *ChartService) GetChartDataById(chartId uint64, loginUser *entity.User) (*entity.ChartDataJSON, *ecode.ErrorWithCode) {
	// 1.校验
	if chartId == 0 {
		return nil, ecode.GetErrWithDetail(ecode.PARAMS_ERROR, "参数错误")
	}
	// 2.获取图表数据，校验是否是本人或者管理员获取
	chart, err := s.ChartRepo.FindById(nil, chartId)
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	if chart == nil {
		return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图表不存在")
	}
	if chart.UserID != loginUser.ID && loginUser.UserRole != consts.ADMIN_ROLE {
		return nil, ecode.GetErrWithDetail(ecode.FORBIDDEN_ERROR, "没有权限获取该图表数据")
	}
	//获取图表数据
	var chartData *entity.ChartDataJSON
	if chartId != 0 {
		chartData, err = s.ChartRepo.FindChartDataByChartDataId(nil, chart.ChartDataID)
	}
	if err != nil {
		return nil, ecode.GetErrWithDetail(ecode.SYSTEM_ERROR, "数据库操作错误")
	}
	if chartData == nil {
		return nil, ecode.GetErrWithDetail(ecode.NOT_FOUND_ERROR, "图表数据不存在")
	}
	return chartData, nil
}
