package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	r.GET("/api/record/list", controller.GetAttendanceList)
	r.GET("/api/record/attendance/day", controller.GetAttendanceByDate)
}
