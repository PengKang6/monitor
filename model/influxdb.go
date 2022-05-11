package model

import "time"

type DataPoint struct {
	Measurement string                 `json:"measurement"`
	Tags        map[string]string      `json:"tags"`
	Fields      map[string]interface{} `json:"fields"`
	Timestamp   time.Time              `json:"timestamp"`
}

type QueryData struct {
	Bucket string      `json:"bucket"`
	Start  string      `json:"start"`
	End    string      `json:"end"`
	Fields [][3]string `json:"fields"`
}
