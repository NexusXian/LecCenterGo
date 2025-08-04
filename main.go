package main

import (
	"LecCenterGo/initialization"
	"LecCenterGo/routers"
)

func init() {
	initialization.InitViper()
	initialization.InitMySQL()
	initialization.InitAvatar()
}

func main() {
	r := routers.Router()
	r.Static("/images", "./images")
	err := r.Run(":8080")
	if err != nil {
		return
	}
}
