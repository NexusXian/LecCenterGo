package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"net/url"
	"time"
)

func ExportDayData(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "缺少日期参数"})
		return
	}

	t, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "日期格式错误，请使用 YYYY/MM/DD 格式"})
		return
	}

	formattedDate := t.Format("2006-01-02")

	var records []models.Attendance
	DB := dao.DB

	var dateFilterCondition string
	dbName := DB.Dialector.Name() // 获取数据库类型
	switch dbName {
	case "sqlite":
		dateFilterCondition = fmt.Sprintf("strftime('%%Y-%%m-%%d', created_at) = '%s'", formattedDate)
	case "mysql":
		dateFilterCondition = fmt.Sprintf("DATE_FORMAT(created_at, '%%Y-%%m-%%d') = '%s'", formattedDate)
	case "postgres":
		dateFilterCondition = fmt.Sprintf("TO_CHAR(created_at, 'YYYY-MM-DD') = '%s'", formattedDate)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": "不支持的数据库类型，无法进行日期过滤"})
		return
	}

	// 查询指定日期的所有打卡记录，并按创建时间升序排序
	if err := DB.Where(dateFilterCondition).Order("created_at ASC").Find(&records).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": fmt.Sprintf("获取导出数据失败: %v", err)})
		return
	}

	if len(records) == 0 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "message": "该日暂无打卡数据可导出"})
		return
	}

	// --- 开始生成 Excel 文件 ---
	f := excelize.NewFile() // 创建一个新的 Excel 文件

	// 设置工作表名称
	sheetName := "打卡记录"
	// 删除默认的 Sheet1，并创建新的工作表，或者直接使用 index 0 的工作表
	// 为了清晰，我们这里仍然创建新的，并设置为活跃
	index, err := f.NewSheet(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": fmt.Sprintf("创建 Excel 工作表失败: %v", err)})
		return
	}
	err = f.DeleteSheet("Sheet1")
	if err != nil {
		return
	} // 删除默认创建的 Sheet1

	// 设置表头
	// **注意：这里不再包含 "ID"**
	headers := []string{"序号", "用户名", "邮箱", "打卡时间", "打卡信息"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1) // 例如 A1, B1...
		err := f.SetCellValue(sheetName, cell, header)
		if err != nil {
			return
		}
	}

	// 写入数据
	for i, record := range records {
		rowIndex := i + 2 // 数据从第二行开始写入 (第一行是表头)

		// **注意：这里不再写入 record.ID**
		err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowIndex), i+1)
		if err != nil {
			return
		} // 序号 (自定义的行号)
		err = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowIndex), record.Username)
		if err != nil {
			return
		} // 用户名
		err = f.SetCellValue(sheetName, fmt.Sprintf("C%d", rowIndex), record.Email)
		if err != nil {
			return
		} // 邮箱
		err = f.SetCellValue(sheetName, fmt.Sprintf("D%d", rowIndex), record.CreatedAt.Format("2006-01-02 15:04:05"))
		if err != nil {
			return
		} // 打卡时间
		err = f.SetCellValue(sheetName, fmt.Sprintf("E%d", rowIndex), record.CheckInfo)
		if err != nil {
			return
		} // 打卡信息
	}

	// 设置活动工作表
	f.SetActiveSheet(index)

	// --- 设置 HTTP 响应头，通知浏览器下载 Excel 文件 ---
	// 1. 原始文件名（包含中文）
	fileName := fmt.Sprintf("打卡记录_%s.xlsx", formattedDate)
	// 2. 编码中文文件名（关键：处理非ASCII字符）
	encodedFileName := url.QueryEscape(fileName)

	// 3. 设置响应头（按标准格式传递）
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	// 核心：使用filename*=UTF-8''格式传递编码后的文件名
	c.Header("Content-Disposition", fmt.Sprintf(
		"attachment; filename=\"%s\"; filename*=UTF-8''%s",
		fileName,        // 兼容旧浏览器（可选）
		encodedFileName, // 现代浏览器优先解析，确保中文正常显示
	))
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache") // 避免浏览器缓存

	// 将 Excel 文件写入 Gin 的响应体
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": fmt.Sprintf("写入 Excel 文件到响应失败: %v", err)})
		return
	}
}
