package repository

import (
	"errors"
	"gorm.io/gorm"
	"mrbi/internal/model/entity"
	"mrbi/pkg/db"
)

// 数据库操作层
type ChartRepository struct {
	db *gorm.DB
}

func NewChartRepository() *ChartRepository {
	return &ChartRepository{db.LoadDB()}
}

// 开启事务
func (r *ChartRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *ChartRepository) AddChart(tx *gorm.DB, chart *entity.Chart) (uint64, error) {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Create(chart).Error; err != nil {
		return 0, err
	}
	return chart.ID, nil
}

// 根据Id和用户Id删除图表，不存在则返回错误
func (r *ChartRepository) DeleteChart(tx *gorm.DB, chartId uint64) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Where("id = ?", chartId).Delete(&entity.Chart{}).Error; err != nil {
		return err
	}
	return nil
}

// 根据图标的Id查找图表，不存在则返回nil
func (r *ChartRepository) FindById(tx *gorm.DB, chartId uint64) (*entity.Chart, error) {
	if tx == nil {
		tx = r.db
	}
	var chart entity.Chart
	if err := tx.Where("id = ?", chartId).First(&chart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //无记录
		}
		return nil, err //数据库查询异常
	}
	return &chart, nil
}

// 根据Map的字段来更新Chart
func (r *ChartRepository) UpdateChartByMap(tx *gorm.DB, chartId uint64, updateMap map[string]interface{}) error {
	if tx == nil {
		tx = r.db
	}
	if err := tx.Model(&entity.Chart{}).Where("id = ?", chartId).Updates(updateMap).Error; err != nil {
		return err
	}
	return nil
}
