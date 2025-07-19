package routers

import (
	"LecCenterGo/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORS())
	RegisterUserRouter(r)
	RegisterFileLoadRouter(r)
	RegisterAdminRouter(r)
	RegisterRankRouter(r)
	RegisterQueryRouter(r)
	RegisterNoticeRouter(r)
	RegisterClockRouter(r)
	RegisterRecordRouter(r)
	RegisterExportRoutes(r)
	return r
}
