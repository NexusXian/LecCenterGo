package main

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
	"LecCenterGo/routers"
	"LecCenterGo/utils"
	"fmt"
)

func init() {
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	err = dao.DB.AutoMigrate(&models.User{}, &models.Notice{}, &models.Attendance{})
	if err != nil {
		return
	}
	err = utils.CreateAvatarDirectory()
	if err != nil {
		fmt.Printf("创建头像目录失败: %v\n", err)
	}
	err = utils.CreateCheckDirectory()
	if err != nil {
		fmt.Printf("创建签到目录失败: %v\n", err)
	}
}

func main() {
	r := routers.Router()
	r.Static("/images", "./images")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
