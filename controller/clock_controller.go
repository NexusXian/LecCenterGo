package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 常量定义
const Dir = "images/clock" // 签到图片保存目录

// 1. 提交签到（含图片上传）
// （无需修改，表单上传本身通过FormData传递参数）
func CombinedCheckinHandler(c *gin.Context) {
	// 获取用户信息（表单上传仍通过PostForm获取）
	userEmail := c.PostForm("email")
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少用户邮箱"})
		return
	}

	// 查找用户
	var user models.User
	if result := dao.DB.Where("email = ?", userEmail).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询用户失败"})
		return
	}

	// 处理上传的图片
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传签到图片"})
		return
	}

	// 验证文件大小和类型
	if file.Size > 10*1024*1024 { // 10MB限制
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件过大，最大10MB"})
		return
	}
	allowedExtensions := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true,
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效文件类型"})
		return
	}

	// 生成新文件名并保存
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(Dir, newFilename)
	if err := os.MkdirAll(Dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建目录失败"})
		return
	}
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	publicPath := "/" + filePath

	// 开始事务
	err = dao.DB.Transaction(func(tx *gorm.DB) error {
		// 创建签到记录
		attendance := models.Attendance{
			Username:  user.Username,
			Email:     userEmail,
			CreatedAt: time.Now(),
			Img:       publicPath,
			CheckInfo: c.PostForm("check_info"), // 签到备注
		}

		if err := tx.Create(&attendance).Error; err != nil {
			return err
		}

		// 更新用户签到次数
		user.Count++
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		_ = os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签到处理失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "签到成功",
		"code":    http.StatusOK,
		"img":     publicPath,
	})
}

// 2. 获取签到记录列表（修改为JSON接收email）
func GetRecordList(ctx *gin.Context) {
	// 定义接收JSON参数的结构体
	type ListRequest struct {
		Email    string `json:"email" binding:"required"` // 必传
		Page     int    `json:"page" binding:"omitempty,min=1"`
		PageSize int    `json:"pageSize" binding:"omitempty,min=1,max=20"`
	}

	var req ListRequest
	// 从JSON请求体解析参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	// 分页参数默认值
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 || pageSize > 20 {
		pageSize = 12
	}
	offset := (page - 1) * pageSize

	// 查询该用户的所有签到记录（按时间倒序）
	var records []models.Attendance
	var total int64

	// 先查总数
	if err := dao.DB.Model(&models.Attendance{}).
		Where("email = ?", req.Email).
		Count(&total).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "查询记录失败",
		})
		return
	}

	// 再查当前页数据
	if err := dao.DB.
		Where("email = ?", req.Email).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "查询记录失败",
		})
		return
	}

	// 返回分页结果
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": records,
		"pagination": gin.H{
			"page":     page,
			"pageSize": pageSize,
			"total":    total,
			"pages":    (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// 3. 获取单条签到记录详情（修改为JSON接收email）
func GetRecordDetail(ctx *gin.Context) {
	// 定义接收JSON参数的结构体
	type DetailRequest struct {
		Email string `json:"email" binding:"required"` // 必传
	}

	var req DetailRequest
	// 从JSON请求体解析参数
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	// 获取记录ID（仍从URL路径获取）
	recordID := ctx.Param("id")
	if recordID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "记录ID缺失",
		})
		return
	}

	// 查询单条记录
	var record models.Attendance
	if err := dao.DB.First(&record, recordID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "记录不存在",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "查询记录失败",
		})
		return
	}

	// 权限验证：只能查看自己的记录
	if record.Email != req.Email {
		ctx.JSON(http.StatusForbidden, gin.H{
			"code":    http.StatusForbidden,
			"message": "没有权限查看该记录",
		})
		return
	}

	// 返回详情
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": record,
	})
}
