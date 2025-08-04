package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterFileLoadRouter(router *gin.Engine) {
	LoadFile := router.Group("/api/file")
	{
		LoadFile.POST("/avatar", controller.UploadAvatars)
	}
}
