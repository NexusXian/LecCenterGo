package routers

import (
	"LecCenterGo/controller"
	"github.com/gin-gonic/gin"
)

func RegisterUserRouter(r *gin.Engine) {
	UserGroup := r.Group("/api/user")
	{
		UserGroup.POST("/login", controller.LoginHandler)
		UserGroup.POST("/register", controller.RegisterHandler)
		UpdateGroup := UserGroup.Group("/update")
		{
			UpdateGroup.POST("/password", controller.UpdatePasswordHandler)
		}
	}
}
