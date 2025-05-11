package entity

import (
	"gorm.io/gorm"
	"time"
)

type Chart struct {
	ID         uint64         `gorm:"primaryKey;comment:id" json:"id,string" swaggertype:"string"`
	Goal       string         `gorm:"type:text;comment:分析目标"`
	ChartData  string         `gorm:"type:text;comment:图表数据"`
	ChartType  string         `gorm:"type:varchar(128);comment:图表类型"`
	GenChart   string         `gorm:"type:text;comment:AI生成的图表数据"`
	GenResult  string         `gorm:"type:text;comment:AI生成的分析结论"`
	UserID     uint64         `gorm:"comment:创建用户 id"`
	CreateTime time.Time      `gorm:"autoCreateTime;comment:创建时间"`
	UpdateTime time.Time      `gorm:"autoUpdateTime;comment:更新时间"`
	IsDelete   gorm.DeletedAt `gorm:"comment:是否删除" swaggerignore:"true"`
}

// AutoMigrateChart 执行数据库迁移
func AutoMigrateChart(db *gorm.DB) {
	err := db.AutoMigrate(&Chart{})
	if err != nil {
		panic("⚠️ 图表信息表迁移失败: " + err.Error())
	}
}

// 钩子，使用sonyflake生成ID
func (c *Chart) BeforeCreate(tx *gorm.DB) error {
	if c.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		c.ID = id
	}
	return nil
}
