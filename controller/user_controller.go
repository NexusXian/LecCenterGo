package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"LecCenterGo/request"
	"LecCenterGo/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
)

func LoginHandler(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
	}
	var user *models.User
	user, err := models.GetUserByEmailOrPhone(req.Account)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "账号或密码错误",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "数据库查询用户出错",
			})
		}
	}
	if !utils.ComparePassword(user.Password, req.Password) {
		c.JSON(http.StatusOK, gin.H{
			"error": "账号或密码错误！",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"code":    http.StatusOK,
		"data": gin.H{
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"avatar":   user.Avatar,
			"grade":    user.Grade,
			"major":    user.Major,
			"role":     user.Role,
		},
	})
}

func RegisterHandler(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	_, err := models.GetUserByEmail(req.Email)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该邮箱已被注册"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询邮箱失败: " + err.Error()})
		return
	}

	_, err = models.GetUserByPhone(req.Phone)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "该手机号已被注册"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询手机号失败: " + err.Error()})
		return
	}

	var user models.User
	user.Username = req.Username
	user.Phone = req.Phone
	user.Email = req.Email
	user.Grade = req.Grade
	user.Major = req.Major
	user.Direction = req.Direction
	user.Password = req.Password
	user.CreatedAt = time.Now()

	if err := models.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "用户注册失败: ",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"code":    http.StatusOK,
	})
}

// UpdatePasswordHandler 处理用户密码更新请求
func UpdatePasswordHandler(c *gin.Context) {
	var req request.UpdatePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数无效，请检查！", "code": http.StatusBadRequest})
		return
	}

	user, err := models.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在！", "code": http.StatusNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器内部错误，请稍后再试！", "code": http.StatusInternalServerError})
		}
		return
	}

	if user.Password != req.OldPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "旧密码不正确！", "code": http.StatusUnauthorized})
		return
	}

	if req.OldPassword == req.NewPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "新密码不能与旧密码相同！", "code": http.StatusBadRequest})
		return
	}

	user.Password = req.NewPassword
	user.UpdatedAt = time.Now()

	if err := dao.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新密码失败，请稍后再试！", "code": http.StatusInternalServerError})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功！",
		"code":    http.StatusOK,
	})
}
