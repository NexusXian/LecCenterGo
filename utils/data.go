package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 成功响应
func ResponseSuccess(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: message,
		Data:    data,
	})
}

// 错误响应
func ResponseError(c *gin.Context, httpCode, code int, message string, err error) {
	if err != nil {
		c.JSON(httpCode, Response{
			Code:    code,
			Message: message,
			Data:    err.Error(),
		})
	} else {
		c.JSON(httpCode, Response{
			Code:    code,
			Message: message,
		})
	}
}
