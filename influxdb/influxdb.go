package influxdb

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"monitor/config"
	"monitor/model"
)

var InfluxDBClient InfluxDB

type InfluxDB interface {
	WriteInfluxDB(data *model.DataPoint)
	ReadInfluxDB(queryString string) *api.QueryTableResult
	CloseInfluxDB()
}

type InfluxDBImpl struct {
	Client   influxdb2.Client
	WriteAPI api.WriteAPI
	QueryAPI api.QueryAPI
}

func InitInfluxDB(configInfluxDB config.InfluxDBConfig) InfluxDB {

	bucket := configInfluxDB.Bucket
	org := configInfluxDB.Org
	// 此token是基于monitor字符生成的，在8086面板找到的，踩坑。。
	token := configInfluxDB.Token
	// Store the URL of your InfluxDB instance
	url := configInfluxDB.Url

	client := influxdb2.NewClient(url, token)
	InfluxDBClient = &InfluxDBImpl{
		Client:   client,
		WriteAPI: client.WriteAPI(org, bucket),
		QueryAPI: client.QueryAPI(org),
	}

	return InfluxDBClient
}

func (idb *InfluxDBImpl) WriteInfluxDB(data *model.DataPoint) {
	//p := influxdb2.NewPoint("stat",
	//	map[string]string{"unit": "temperature"},
	//	map[string]interface{}{"avg": 27.3, "max": 48},
	//	time.Now())
	p := influxdb2.NewPoint(
		data.Measurement,
		data.Tags,
		data.Fields,
		data.Timestamp,
	)
	idb.WriteAPI.WritePoint(p)
	idb.WriteAPI.Flush()
}

//func ReadInfluxDB(queryData *model.QueryData) *api.QueryTableResult {
func (idb *InfluxDBImpl) ReadInfluxDB(queryString string) *api.QueryTableResult {
	//result, err := queryAPI.Query(context.Background(), `from(bucket:"test")
	//|> range(start: -1h)
	//|> filter(fn: (r) => r._measurement == "stat")`)

	//queryString := `from(bucket: "example-bucket")
	//|> range(start: -1h)
	//|> filter(fn: (r) => r._measurement == "example-measurement" and r.tag == "example-tag")
	//|> filter(fn: (r) => r._field == "example-field")
	//`

	//queryString := `from(bucket:"my-bucket") |> range(start: duration(params.start)) |> filter(fn: (r) => r._measurement == "stat") |> filter(fn: (r) => r._field == params.field) |> filter(fn: (r) => r._value > params.value)`
	result, err := idb.QueryAPI.Query(context.Background(), queryString)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			//fmt.Printf("value: %v\n", result.Record().Value())
			fmt.Println(result.Record())
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	return result

	/*
		// 2.8.0才能用，而2.8.0依赖go1.17
		parameters := struct {
				Start string  `json:"start"`
				Field string  `json:"field"`
				Value float64 `json:"value"`
			}{
				"-1h",
				"temperature",
				25,
			}
			// Query with parameters
			query := `from(bucket:"my-bucket")
						|> range(start: duration(params.start))
						|> filter(fn: (r) => r._measurement == "stat")
						|> filter(fn: (r) => r._field == params.field)
						|> filter(fn: (r) => r._value > params.value)`

			// Get result
			result, err := queryAPI.QueryWithParams(context.Background(), query, parameters)
	*/
}

func (idb *InfluxDBImpl) CloseInfluxDB() {
	idb.Client.Close()
}

//func Bucket(bucket string) string {
//	return "from(bucket:\"" + bucket + "\")"
//}
//
//func TimeRange(start string, stop string) string {
//	if stop != "" {
//		return "range(start: " + start + ", stop: " + stop + ")"
//	}
//	return "range(start: " + start + ")"
//}
//
//func Filter(items [][3]string) string {
//	var res []string
//	for _, v := range items {
//		res = append(res, "r."+v[0]+v[1]+v[2])
//	}
//	return "filter(fn: (r) =>" + strings.Join(res, " and ") + ")"
//}
