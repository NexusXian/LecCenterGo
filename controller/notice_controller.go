package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func GetNoticeList(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "21")

	page, err1 := strconv.Atoi(pageStr)
	pageSize, err2 := strconv.Atoi(pageSizeStr)
	if err1 != nil || err2 != nil || page <= 0 || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page or pageSize"})
		return
	}

	offset := (page - 1) * pageSize

	var notices []models.Notice
	var total int64

	db := dao.DB // 假设你有全局 DB 实例
	db.Model(&models.Notice{}).Count(&total)

	// 修改排序逻辑：
	// 1. 首先按 is_important 降序排序 (true/1 会排在 false/0 前面)
	// 2. 然后按 updated_at 降序排序 (最新更新的排在前面)
	// 3. 最后按 created_at 降序排序 (最新创建的排在前面)
	err := db.Order("is_important desc").
		Order("updated_at desc"). // Added this line for sorting by updated_at
		Order("created_at desc").
		Limit(pageSize).
		Offset(offset).
		Find(&notices).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "数据库查询失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
		"data":     notices,
	})
}

// 获取单个通知详情
func GetNotice(c *gin.Context) {
	// 获取路径参数中的ID
	noticeID := c.Param("id")

	var notice models.Notice

	// 查询通知
	result := dao.DB.Where("notice_id = ?", noticeID).First(&notice)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "通知不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": "获取通知详情失败",
				"error":   result.Error.Error(),
			})
		}
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取通知详情成功",
		"data":    notice,
	})
}

func CreateNotice(c *gin.Context) {
	var notice models.Notice
	db := dao.DB
	// 绑定请求JSON数据
	if err := c.ShouldBindJSON(&notice); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必要字段
	if notice.Title == "" || notice.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title and Content are required"})
		return
	}

	// 保存到数据库
	result := db.Create(&notice)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create notice"})
		return
	}

	c.JSON(http.StatusCreated, notice)
}

func UpdateNotice(c *gin.Context) {
	db := dao.DB
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notice ID"})
		return
	}

	var notice models.Notice
	result := db.First(&notice, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Notice not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notice"})
		}
		return
	}

	// 绑定更新数据
	var updateData map[string]interface{}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	result = db.Model(&notice).Updates(updateData)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notice"})
		return
	}

	// 获取更新后的通知
	db.First(&notice, id)

	c.JSON(http.StatusOK, notice)
}

// DeleteNotice 处理删除通告的请求
func DeleteNotice(c *gin.Context) {
	db := dao.DB // 获取数据库实例

	// 尝试将URL参数中的ID转换为整数
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// 如果ID无效，返回400状态码和自定义错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    40001, // 自定义错误码：无效的通告ID
			"message": "Invalid notice ID",
			"data":    nil,
		})
		return
	}

	var notice models.Notice // 定义一个Notice模型变量
	// 查找指定ID的通告
	result := db.First(&notice, id)
	if result.Error != nil {
		// 如果通告未找到
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound, // 自定义错误码：通告未找到
				"message": "Notice not found",
				"data":    nil,
			})
		} else {
			// 其他内部服务器错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError, // 自定义错误码：获取通告失败
				"message": "Failed to fetch notice",
				"data":    nil,
			})
		}
		return
	}

	// 尝试从数据库中删除找到的通告
	result = db.Delete(&notice)
	if result.Error != nil {
		// 如果删除失败，返回500状态码和自定义错误信息
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500, // 自定义错误码：删除通告失败
			"message": "Failed to delete notice",
			"data":    nil,
		})
		return
	}

	// 删除成功，返回200状态码和成功信息
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK, // 0 通常表示成功
		"message": "Notice deleted successfully",
		"data":    nil,
	})
}
