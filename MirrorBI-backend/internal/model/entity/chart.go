package entity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Chart struct {
	ID          uint64         `gorm:"primaryKey;comment:id" json:"id,string" swaggertype:"string"`
	Name        string         `gorm:"type:varchar(128);comment:图表名称" json:"name"`
	Goal        string         `gorm:"type:text;comment:分析目标" json:"goal"`
	ChartDataID uint64         `gorm:"comment:存放图表数据的表中元组ID" json:"chartDataId,string" swaggertype:"string"`
	ChartType   string         `gorm:"type:varchar(128);comment:图表类型" json:"chartType"`
	GenChart    string         `gorm:"type:text;comment:AI生成的图表数据" json:"genChart"`
	GenResult   string         `gorm:"type:text;comment:AI生成的分析结论" json:"genResult"`
	UserID      uint64         `gorm:"comment:创建用户 id" json:"userId,string" swaggertype:"string"`
	CreateTime  time.Time      `gorm:"autoCreateTime;comment:创建时间" json:"createTime"`
	UpdateTime  time.Time      `gorm:"autoUpdateTime;comment:更新时间" json:"updateTime"`
	IsDelete    gorm.DeletedAt `gorm:"comment:是否删除" swaggerignore:"true" json:"isDelete"`
}

type ChartDataJSON struct {
	ID   uint64         `json:"id,string" swaggertype:"string"`
	Data datatypes.JSON `gorm:"type:json;comment:二维表格数据" json:"data"`
}

// AutoMigrateChart 执行数据库迁移
func AutoMigrateChart(db *gorm.DB) {
	err := db.AutoMigrate(&Chart{})
	if err != nil {
		panic("⚠️ 图表信息表迁移失败: " + err.Error())
	}
	err = db.AutoMigrate(&ChartDataJSON{})
	if err != nil {
		panic("⚠️ 图表数据表迁移失败: " + err.Error())
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

// 钩子，使用sonyflake生成ID
func (c *ChartDataJSON) BeforeCreate(tx *gorm.DB) error {
	if c.ID == 0 {
		id, err := sf.NextID()
		if err != nil {
			return err
		}
		c.ID = id
	}
	return nil
}
