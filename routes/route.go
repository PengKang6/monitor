package routes

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"monitor/api"
	"monitor/config"
)

func buildRoute(r *gin.Engine) {
	//调用下游服务的路由接口
	r.GET("/get/*realUrl", api.Get)
	r.POST("/post/*realUrl", api.Post)

	//看板接口

	//管理接口
}

func InitRoute(config config.ServerConfig) {
	r := gin.Default()

	// 注入路由接口
	buildRoute(r)

	// 跨域配置
	r.Use(cors.Default())

	// 基于对应端口启动服务
	if err := r.Run(":" + config.Port); err != nil {
		fmt.Println("startup service failed.")
	}
}
