package initialization

import (
	"LecCenterGo/utils"
	"fmt"
)

func InitAvatar() {
	err := utils.CreateAvatarDirectory()
	if err != nil {
		fmt.Printf("创建头像目录失败: %v\n", err)
	}
	err = utils.CreateCheckDirectory()
	if err != nil {
		fmt.Printf("创建签到目录失败: %v\n", err)
	}
}
