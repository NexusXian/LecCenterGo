package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterExportRoutes(r *gin.Engine) {
	ExportGroup := r.Group("/api")
	{
		ExportGroup.GET("/export", controller.ExportDayData)
	}

}
