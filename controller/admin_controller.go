package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"LecCenterGo/request"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// UpdateUser handles the update of an existing user's information.

func UpdateUser(c *gin.Context) {
	var payload request.UserUpdatePayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	var user models.User
	result := dao.DB.Where("email = ?", payload.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "User not found with the provided email."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "Database query error: " + result.Error.Error()})
		return
	}

	updates := make(map[string]interface{})

	if payload.Username != "" {
		updates["username"] = payload.Username
	}
	if payload.Grade != "" {
		updates["grade"] = payload.Grade
	}
	if payload.Major != "" {
		updates["major"] = payload.Major
	}
	if payload.Phone != "" {
		updates["phone"] = payload.Phone
	}
	if payload.Role != "" {
		updates["role"] = payload.Role
	}
	if payload.Direction != "" {
		updates["direction"] = payload.Direction
	}

	updateResult := dao.DB.Model(&user).Updates(updates)

	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "Failed to update user: " + updateResult.Error.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "User updated successfully!"})
}

// DeleteUser handles the deletion of a user.
func DeleteUser(c *gin.Context) {
	var payload request.UserDeletePayload
	// Bind JSON request body to payload struct
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "msg": err.Error()})
		return
	}

	// 1. Find the user by email
	var user models.User
	result := dao.DB.Where("email = ?", payload.Email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"code": http.StatusNotFound, "msg": "User not found with the provided email."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "Database query error: " + result.Error.Error()})
		return
	}

	// Prevent admin users from being deleted
	if user.Role == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "msg": "Cannot delete an admin user."})
		return
	}
	deleteResult := dao.DB.Delete(&user)

	if deleteResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "msg": "Failed to delete user: " + deleteResult.Error.Error()})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "msg": "User deleted successfully!"})
}
func GetUserListWithPagination(c *gin.Context) {
	// 1. 获取并验证分页参数
	page, err := parsePageParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 统一每页数量为5
	const pageSize = 5
	offset := (page - 1) * pageSize

	// 2. 查询用户总数
	var total int64
	if err := dao.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户总数失败"})
		return
	}

	// 计算总页数（确保类型一致性）
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)

	// 3. 处理无效页码请求
	if page > int(totalPages) && total > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求的页码超出范围"})
		return
	}

	// 4. 处理空结果集
	if total == 0 {
		returnSuccessResponse(c, []models.User{}, total, page, pageSize)
		return
	}

	// 5. 查询当前页数据（管理员优先，按创建时间排序）
	var userList []models.User
	if err := dao.DB.
		Order("CASE WHEN role = 'admin' THEN 0 ELSE 1 END, created_at ASC").
		Limit(pageSize).
		Offset(offset).
		Find(&userList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "获取用户列表失败"})
		return
	}

	// 6. 返回成功响应
	returnSuccessResponse(c, userList, total, page, pageSize)
}

// 解析并验证页码参数
func parsePageParam(c *gin.Context) (int, error) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 0, errors.New("无效的页码参数，必须为正整数")
	}
	return page, nil
}

// 统一返回成功响应
func returnSuccessResponse(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"list":     list,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}
