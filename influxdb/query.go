package influxdb

func QueryReqNums(downstreamInstance string, ch chan int64) {
	queryString := `from(bucket: "test")
	 |> range(start: -15s)
	 |> filter(fn: (r) => r["_measurement"] == "interact")
	 |> filter(fn: (r) => r["downstreamInstance"] == "` + downstreamInstance + `")
	 |> count()`
	result := InfluxDBClient.ReadInfluxDB(queryString)
	val, ok := result.Record().Value().(int64)
	if ok {
		ch <- val
	}
}

func QueryTimeCost(downstreamInstance string, ch chan int64) {
	queryString := `from(bucket: "test")
	 |> range(start: -15s)
	 |> filter(fn: (r) => r["_measurement"] == "interact")
	 |> filter(fn: (r) => r["downstreamInstance"] == "` + downstreamInstance + `")
	 |> filter(fn: (r) => r["_field"] == "timeCost")
	 |> mean()`
	result := InfluxDBClient.ReadInfluxDB(queryString)
	val, ok := result.Record().Value().(int64)
	if ok {
		ch <- val
	}
}
