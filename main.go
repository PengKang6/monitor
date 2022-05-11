package main

import (
	"monitor/common"
	"monitor/config"
	"monitor/http"
	"monitor/influxdb"
	"monitor/routes"
	"monitor/timing"
)

func main() {
	//读取配置文件
	configs := config.LoadConfig("config/configFile")

	//初始化common数据
	common.InitCommon()

	// 初始化InfluxDB
	influxDB := influxdb.InitInfluxDB(configs.InfluxDBConfig)
	defer influxDB.CloseInfluxDB()

	//初始化httpClient
	http.InitHttpClient()

	//初始化定时任务
	//路由初始化以后再来开启定时任务？
	timing.InitTimerTask()
	defer timing.StopTimerTask()

	//初始化路由
	routes.InitRoute(configs.ServerConfig)
}
