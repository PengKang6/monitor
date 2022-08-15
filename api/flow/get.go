package api

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"monitor/common"
	httpUtil "monitor/http"
	"monitor/influxdb"
	"monitor/model"
	"net/http"
	"time"
	"unsafe"
)

func Get(context *gin.Context) {
	//鉴权
	if !common.Interceptor(context) {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "ClientIP Unauthorized"})
		return
	}

	interfaceName := context.Param("realUrl")

	//负载均衡算法获取此次调用下游实例
	downstreamInstance := common.LoadBalance(interfaceName)

	//起始时间/ns
	start := time.Now()

	//调用下游
	var reqData []byte
	var err error
	reqData, err = ioutil.ReadAll(context.Request.Body)
	if err != nil {
		reqData = make([]byte, 0)
	}
	//todo
	url := ""
	resp, err := httpUtil.HttpClient.HttpPost(url, reqData)

	//记录本次调用流量数据
	interactDataPoint := &model.DataPoint{}
	interactDataPoint.Measurement = "interact"
	interactTags := make(map[string]string)
	interactTags["downstreamInstance"] = downstreamInstance
	interactTags["interface"] = interfaceName
	interactTags["clientIP"] = context.ClientIP()
	interactDataPoint.Tags = interactTags
	interactFields := make(map[string]interface{})
	interactFields["statusCode"] = resp.StatusCode
	interactFields["timeCost"] = time.Since(start).Milliseconds()
	interactDataPoint.Fields = interactFields
	interactDataPoint.Timestamp = time.Now()
	go influxdb.InfluxDBClient.WriteInfluxDB(interactDataPoint)

	//响应透传数据
	respData, err := httpUtil.ParseResponse(resp)
	context.JSON(resp.StatusCode, *(*gin.H)(unsafe.Pointer(respData)))
	defer resp.Body.Close()

	return
}
