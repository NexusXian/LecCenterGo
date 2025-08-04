package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterAdminRouter(r *gin.Engine) {
	AdminGroup := r.Group("/api/admin")
	{
		AdminGroup.GET("/info", controller.GetUserListWithPagination)
		AdminGroup.POST("/updateUser", controller.UpdateUser)
		AdminGroup.DELETE("/deleteUser", controller.DeleteUser)
		AdminGroup.GET("/users", controller.GetUserListWithPagination)
	}
}
