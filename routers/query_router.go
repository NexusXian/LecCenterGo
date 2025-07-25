package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterQueryRouter(r *gin.Engine) {
	QueryGroup := r.Group("/query")
	{
		QueryGroup.POST("/user", controller.GetUserList)
	}
}
