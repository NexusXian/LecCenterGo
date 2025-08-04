package dao

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() (err error) {
	// 从配置文件中获取数据库连接信息
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetInt("mysql.port")
	dbname := viper.GetString("mysql.dbname")
	charset := viper.GetString("mysql.charset")

	// 增加调试输出，检查配置是否正确读取
	fmt.Printf("数据库配置信息: \n")
	fmt.Printf("username: %s\n", username)
	fmt.Printf("host: %s\n", host)
	fmt.Printf("port: %d\n", port)
	fmt.Printf("dbname: %s\n", dbname)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		username, password, host, port, dbname, charset)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("数据库连接失败: %v，DSN: %s", err, dsn)
	}

	return nil
}
