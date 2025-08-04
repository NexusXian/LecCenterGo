package initialization

import (
	"LecCenterGo/dao"
	"LecCenterGo/models"
)

func InitMySQL() {
	err := dao.InitMySQL()
	if err != nil {
		panic(err)
	}
	err = dao.DB.AutoMigrate(&models.User{}, &models.Notice{}, &models.Attendance{})
	if err != nil {
		return
	}

}
