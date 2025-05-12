package utils

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

//使用excelize库操作excel文件

// dst为excel文件目标路径
// 返回解析出的CSV内容和错误（CSV未落盘）
func ExcelToCSV(dst string) (string, error) {
	//打开excel文件
	f, err := excelize.OpenFile(dst)
	if err != nil {
		return "", err
	}
	//获取sheet名称
	sheetName := f.GetSheetName(0)
	//获取sheet数据
	// 1. 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return "", fmt.Errorf("读取 sheet %s 失败: %w", sheetName, err)
	}

	// 2. 找到首个非空行作为表头
	var header []string
	var headerIdx int
	for i, row := range rows {
		nonEmpty := false
		for _, cell := range row {
			if strings.TrimSpace(cell) != "" {
				nonEmpty = true
				break
			}
		}
		if nonEmpty {
			header = row
			headerIdx = i
			break
		}
	}
	if header == nil {
		return "", fmt.Errorf("未找到有效表头（所有行均为空）")
	}

	// 3. 确定需要保留的列索引（表头不为空的列）
	var validCols []int      //存储表头有效列下标
	var cleanHeader []string //存储表头有效列的值
	for j, h := range header {
		title := strings.TrimSpace(h)
		if title != "" {
			validCols = append(validCols, j)
			cleanHeader = append(cleanHeader, title)
		}
	}
	if len(validCols) == 0 {
		return "", fmt.Errorf("表头无有效列")
	}

	// 4. 收集数据行，跳过空行
	var cleanData [][]string
	for _, row := range rows[headerIdx+1:] {
		var record []string
		emptyRow := true
		for _, col := range validCols {
			var cell string
			if col < len(row) {
				cell = strings.TrimSpace(row[col])
			}
			if cell != "" {
				emptyRow = false
			}
			record = append(record, cell)
		}
		if !emptyRow {
			cleanData = append(cleanData, record)
		}
	}

	// 5. 生成 CSV 文本
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	// 写入表头
	if err := writer.Write(cleanHeader); err != nil {
		return "", fmt.Errorf("写入 CSV 表头失败: %w", err)
	}
	// 写入所有数据行
	if err := writer.WriteAll(cleanData); err != nil {
		return "", fmt.Errorf("写入 CSV 数据失败: %w", err)
	}
	writer.Flush() //写入到缓冲区buf中
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV 写入过程中出错: %w", err)
	}
	//log.Println(buf.String())
	return buf.String(), nil
}
