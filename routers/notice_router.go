package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterNoticeRouter(r *gin.Engine) {
	NoticeGroup := r.Group("/api/notice")
	{
		NoticeGroup.GET("/list", controller.GetNoticeList)
		NoticeGroup.GET("/:id", controller.GetNotice)
		NoticeGroup.POST("/create", controller.CreateNotice)
		NoticeGroup.POST("/edit/:id", controller.UpdateNotice)
		NoticeGroup.POST("/delete/:id", controller.DeleteNotice)
	}
}
