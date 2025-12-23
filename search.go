package main

import (
	"github.com/gin-gonic/gin"
	"search-system/router"
)

func main() {
	// 创建一个默认的路由引擎
	r := gin.Default()
	router.InitRouter(r)

	r.Run(":8080")
}
