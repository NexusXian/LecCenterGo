package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetAttendanceList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "9") // 前端固定为 9，表示每页显示 9 天

	DB := dao.DB // 使用你的数据库连接

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"code": -1, "message": "无效的页码参数"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 9 { // 限制 pageSize 为 1-9，确保每页显示的天数合理
		pageSize = 9
	}

	// --- 修正的核心逻辑：按天分页 ---

	// 1. 获取所有不重复的打卡日期，并按日期倒序排列
	// 使用 MySQL 的 DATE_FORMAT 函数来提取日期部分并去重
	var allDates []string

	// 构建获取去重日期的原生 SQL 查询
	// DATE_FORMAT(created_at, '%Y-%m-%d') 会将 datetime 转换为 'YYYY-MM-DD' 格式的字符串
	rawQuery := "SELECT DISTINCT DATE_FORMAT(created_at, '%Y-%m-%d') AS day_str FROM attendances ORDER BY day_str DESC"

	if err := DB.Raw(rawQuery).Scan(&allDates).Error; err != nil && err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": fmt.Sprintf("获取打卡日期列表失败: %v", err)})
		return
	}

	totalDays := int64(len(allDates)) // 总的打卡天数

	// 2. 对这些日期进行分页
	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize

	// 检查分页索引是否越界
	if startIndex >= int(totalDays) {
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "Success",
			"data":    []models.Attendance{}, // 当前页没有对应的日期，返回空数据
			"total":   totalDays,
		})
		return
	}
	if endIndex > int(totalDays) {
		endIndex = int(totalDays)
	}

	// 获取当前页需要查询的日期字符串列表
	paginatedDates := allDates[startIndex:endIndex]

	// 3. 根据分页后的日期列表，查询所有对应的打卡记录
	var attendances []models.Attendance

	// 构建 SQL WHERE 子句，用于筛选出属于这些日期的记录
	// 示例： "DATE_FORMAT(created_at, '%Y-%m-%d') IN ('2025-07-18', '2025-07-17')"
	dateFilterValues := make([]string, len(paginatedDates))
	for i, dateStr := range paginatedDates {
		dateFilterValues[i] = fmt.Sprintf("'%s'", dateStr) // 为每个日期字符串加上引号
	}

	// 使用 MySQL 的 DATE_FORMAT 进行过滤
	dateFilterCondition := fmt.Sprintf("DATE_FORMAT(created_at, '%%Y-%%m-%%d') IN (%s)", strings.Join(dateFilterValues, ","))

	// 查询这些日期下的所有打卡记录，并按打卡时间倒序排序（最新的在前面）
	if err := DB.Where(dateFilterCondition).Order("created_at DESC").Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": -1, "message": fmt.Sprintf("获取打卡记录失败: %v", err)})
		return
	}

	// 返回数据：data 包含当前页所有日期的打卡记录，total 是总的打卡天数
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "Success",
		"data":    attendances,
		"total":   totalDays, // 返回总的打卡天数，而非总的记录数
	})
}

func GetAttendanceByDate(c *gin.Context) {
	dateStr := c.Query("date") // Get the 'date' query parameter
	DB := dao.DB               // Assuming 'dao.DB' provides your GORM database instance

	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "缺少日期参数 (date)"})
		return
	}

	// Parse the date string into a time.Time object
	// The format string "2006/01/02" corresponds to "YYYY/MM/DD"
	targetDate, err := time.Parse("2006/01/02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "日期格式无效，请使用 YYYY/MM/DD 格式"})
		return
	}

	// For querying records within a specific day, we need the start and end of that day.
	// Start of the day: targetDate at 00:00:00
	startOfDay := targetDate
	// End of the day: targetDate at 23:59:59.999...
	endOfDay := targetDate.Add(24 * time.Hour).Add(-time.Nanosecond)

	var attendances []models.Attendance
	var total int64

	// Count total records for the specific date
	if err := DB.Model(&models.Attendance{}).
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取总记录数失败"})
		return
	}

	// Fetch attendance records for the specific date, ordered by creation time
	if err := DB.
		Where("created_at >= ? AND created_at <= ?", startOfDay, endOfDay).
		Order("created_at ASC"). // Order by ASC to show earliest records first for a day
		Find(&attendances).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "获取打卡记录失败"})
		return
	}

	// Return data
	c.JSON(http.StatusOK, gin.H{
		"data":  attendances,
		"total": total,
		"date":  dateStr, // Include the requested date in the response for clarity
	})
}
