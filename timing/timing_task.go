package timing

import (
	"monitor/config"
	"monitor/influxdb"
)

//构建定时器定时更新请求和耗时指标数据

var downstreamInstances = config.GetConfig().ServerConfig.DownstreamInstances

func TimerForWeights() {
	//构建两个[]chan来分别并行的查询数据
	chansReqNums := make([]chan int64, len(downstreamInstances))
	chansTimeCost := make([]chan int64, len(downstreamInstances))

	//分别初始化两个chan
	for i, downstreamInstance := range downstreamInstances {
		//todo 问题：协程数量的增加是否也会有一定的对性能影响
		chansReqNums[i] = make(chan int64)
		go influxdb.QueryReqNums(downstreamInstance, chansReqNums[i])
		chansTimeCost[i] = make(chan int64)
		go influxdb.QueryTimeCost(downstreamInstance, chansTimeCost[i])
	}

	//基于通道同步获取对应数据
	//构建一个临时的weights，然后再构建临界区来写入weights数据
	var weightsTemp = make([]float32, len(downstreamInstances))
	var weightsTempCount float32 = 0
	var reqNums = make([]int64, len(downstreamInstances))
	var reqNumsCount int64 = 0
	var timeCosts = make([]int64, len(downstreamInstances))
	var timeCostsCount int64 = 0
	for i := range downstreamInstances {
		reqNums[i] = <-chansReqNums[i]
		timeCosts[i] = <-chansTimeCost[i]
		reqNumsCount += reqNums[i]
		timeCostsCount += timeCosts[i]
	}
	for i := range downstreamInstances {
		weightsTemp[i] = (float32(reqNums[i])/float32(reqNumsCount))*0.5 + (float32(timeCosts[i])/float32(timeCostsCount))*0.5
		weightsTempCount += weightsTemp[i]
	}

}
