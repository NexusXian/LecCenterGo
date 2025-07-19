package controller

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"LecCenterGo/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserList(c *gin.Context) {
	var params response.QueryParams

	// 绑定 JSON 请求体
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Msg: "参数错误或JSON格式不正确: " + err.Error()})
		return
	}

	// 显式验证分页参数并设置默认值
	if params.Page < 1 {
		params.Page = 1 // 默认为第一页
	}
	if params.PageSize < 1 || params.PageSize > 100 { // 示例：限制 PageSize
		params.PageSize = 5 // 默认或一个合理的限制
	}

	query := dao.DB.Model(&models.User{})

	// 根据查询参数构建查询条件
	if params.Username != "" {
		query = query.Where("username LIKE ?", "%"+params.Username+"%")
	}

	if params.Major != "" {
		query = query.Where("major = ?", params.Major)
	}

	if params.Grade != "" {
		query = query.Where("grade = ?", params.Grade)
	}

	if params.Direction != "" {
		query = query.Where("direction = ?", params.Direction)
	}

	if params.Role != "" {
		query = query.Where("role = ?", params.Role)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Msg: "获取总数失败"})
		return
	}

	// 查询列表
	var users []models.User
	offset := (params.Page - 1) * params.PageSize

	if err := query.
		Select("*, CASE WHEN role = 'admin' THEN 0 ELSE 1 END as role_order").
		Order("role_order ASC").
		Order("created_at ASC").
		Offset(offset).
		Limit(params.PageSize).
		Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Msg: "获取数据失败"})
		return
	}

	// 返回结果
	c.JSON(http.StatusOK, response.Response{
		Code: 200,
		Msg:  "获取成功",
		Data: response.PageResult{
			List:  users,
			Total: total,
		},
	})
}
