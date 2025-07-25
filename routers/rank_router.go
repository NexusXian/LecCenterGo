package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRankRouter(r *gin.Engine) {
	RankGroup := r.Group("/rank")
	{
		RankGroup.GET("/board", controller.GetTopFiveLeaderboard)
		RankGroup.GET("/table", controller.GetUserListWithPagination)
	}
}
