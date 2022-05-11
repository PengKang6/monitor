package common

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"monitor/config"
)

type Weight struct {
}

var downstreamInstances = config.GetConfig().ServerConfig.DownstreamInstances
var bucket = config.GetConfig().InfluxDBConfig.Bucket
var weightsMap map[string]float32
var weights []float32

func LoadBalance(interfaceName string) string {
	//minCallDownstreamInstance := downstreamInstances[0]
	//minCallNum := math.MaxInt64

	//暂时先通过同步的方式来进行数据读取，后续可以基于协程异步实现
	//协程池并行调用实现。
	//for _, downstreamInstance := range downstreamInstances {
	//	var query []string
	//	query = append(query, Bucket(bucket))
	//	query = append(query, TimeRange("-5m", ""))
	//
	//	//var filters1 = [][3]string{{"_measurement", "==", "\"interact\""}, {"_filed", "==", "downstreamInstance"}, {"_filed", "==", "downstreamInstance"}}
	//	//filters = append(filters, [3]string{"_measurement", "==", "\"request\""})
	//	//filters = append(filters, [3]string{"_filed", "==", "downstreamInstance"})
	//	//filters = append(filters, [3]string{"_filed", "==", "downstreamInstance"})
	//
	//	//最终实现方式
	//	//`from(bucket: "test")
	//	//  |> range(start: -1h)
	//	//  |> filter(fn: (r) => r["_measurement"] == "stat")
	//	//  |> filter(fn: (r) => r["_field"] == "avg")
	//	//  |> filter(fn: (r) => r["unit"] == "temperature")
	//	//  |> count()`
	//
	//	//query = append(query, utils.Filter(filters))
	//	fmt.Println(downstreamInstance)
	//}

	return WeightedRandomIndex()
}

//基于当前的权重值数组返回对应的下游实例
func WeightedRandomIndex() string {
	if len(weights) == 1 {
		return downstreamInstances[0]
	}
	var sum float32 = 0.0
	for _, w := range weights {
		sum += w
	}
	r := rand.Float32() * sum
	var t float32 = 0.0
	for i, w := range weights {
		t += w
		if t > r {
			return downstreamInstances[i]
		}
	}
	return downstreamInstances[len(weights)-1]
}

//初始化相关数据
func InitCommon() {
	weightsMap = make(map[string]float32, len(downstreamInstances))
	weights = make([]float32, len(downstreamInstances))

	for i, downstreamInstance := range downstreamInstances {
		weights[i] = float32(1.0 / len(downstreamInstances))
		weightsMap[downstreamInstance] = weights[i]
	}
}

//拦截请求
func Interceptor(context *gin.Context) bool {
	for _, ip := range config.GetConfig().UpstreamInstances {
		if ip == context.ClientIP() {
			return true
		}
	}

	return false
}
