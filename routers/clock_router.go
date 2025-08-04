package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterClockRouter(router *gin.Engine) {
	ClockGroup := router.Group("/api")
	{
		ClockGroup.POST("/checkin", controller.CombinedCheckinHandler)
		ClockGroup.POST("/checkin/records", controller.GetRecordList)
		ClockGroup.POST("/checkin/records/:id", controller.GetRecordDetail)
	}
}
