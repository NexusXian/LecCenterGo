package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRecordRouter(r *gin.Engine) {
	r.GET("/record/list", controller.GetAttendanceList)
	r.GET("/record/attendance/day", controller.GetAttendanceByDate)
}
