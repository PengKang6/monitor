package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"net/http"
	"strings"
)

var HttpClient Http

type Http interface {
	HttpGet(url string) (*http.Response, error)
	HttpPost(url string, data []byte) (*http.Response, error)
}

type HttpClientImpl struct {
	Client *http.Client
}

func InitHttpClient() {
	HttpClient = &HttpClientImpl{
		Client: &http.Client{},
	}
}

func (hc *HttpClientImpl) HttpGet(url string) (*http.Response, error) {
	// todo
	//var err error
	//req := &http.Request{}
	resp := &http.Response{}

	return resp, nil
}

func (hc *HttpClientImpl) HttpPost(url string, data []byte) (*http.Response, error) {
	req := &http.Request{}
	resp := &http.Response{}
	var err error

	//data, err := json.Marshal(&dataObj)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}
	//req, err := http.NewRequest("POST", "http://127.0.0.1:6789/post", strings.NewReader(string(data)))
	req, err = http.NewRequest("POST", url, strings.NewReader(string(data)))
	if err != nil {
		fmt.Println(err)
		return resp, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err = hc.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return resp, err
	}
	//bytes, err := ioutil.ReadAll(resp.Body)
	//if err != nil{
	//	fmt.Println(err)
	//	return nil
	//}
	//defer resp.Body.Close()
	return resp, nil
}

func ParseResponse(resp *http.Response) (*map[string]interface{}, error) {
	data := make(map[string]interface{})
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &data, err
	}
	err = json.Unmarshal(bytes, data)
	return &data, err
}
