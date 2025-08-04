package initialization

import (
	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigName("config")   // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")     // 配置文件类型
	viper.AddConfigPath(".")        // 当前目录
	viper.AddConfigPath("./config") // 配置文件可能存在的其他路径

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
