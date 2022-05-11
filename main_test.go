package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"monitor/config"
	"monitor/utils"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_config(t *testing.T) {
	configs := config.LoadConfig("config/configFile")
	print(json.Marshal(configs))
}

func Test_idea(t *testing.T) {

	print(from("hello"))
}

func from(bucket string) string {
	return "from(bucket:\"" + bucket + "\")"
}

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func Test_json(t *testing.T) {
	//var s []string
	//s = append(s, "hello")
	//fmt.Println(strings.Join(s, " and "))
	//s = append(s, "world")
	//fmt.Println(strings.Join(s, " and "))

	b := strings.Contains("ssdf", "ef")
	fmt.Println(b)

	//user := User{}
	//fmt.Println(user.Age)
	//fmt.Println(user.Name)
	//s, _ := jsoniter.MarshalToString([]int64{2, 3, 4})
	//r := strconv.FormatInt(int64(12), 10)
	//print(s)
	//print(r)
}

func Test_ticker(t *testing.T) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("hello")
	}
}

//func Test_timing(t *testing.T) {
//	tick := utils.NewMyTick(1, Task)
//	tick.Start()
//	defer tick.Stop()
//}

func Task() {
	fmt.Println("hello")
}

func Test_intn(t *testing.T) {

	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		fmt.Println(rand.Float32())
	}
}

func Test_rand(t *testing.T) {
	weights := []float32{0.25, 0.25, 0.25, 0.25}
	var results = make([]int, len(weights))

	rand.Seed(time.Now().Unix())

	for i := 0; i < 1000; i++ {
		results[Weighted(weights)]++
	}
	fmt.Printf("%v\n", results)

}

func Weighted(weights []float32) int {
	if len(weights) == 1 {
		return 0
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
			return i
		}
	}
	return len(weights) - 1
}

func Test_str(t *testing.T) {
	start := time.Now().UnixNano() / 1e6
	fmt.Println(start)
	time.Sleep(2 * time.Millisecond)
	fmt.Println(time.Now().UnixNano() / 1e6)
	fmt.Println(time.Now().UnixNano()/1e6 - start)

}

func Test_routine(t *testing.T) {
	chs := make([]chan int, 3)
	//必须初始化[]chan的所有元素
	for i := range chs {
		chs[i] = make(chan int)
	}

	nums := make([]int, 3)
	for i := 0; i < 3; i++ {
		go hello(i, chs[i])
	}
	fmt.Println("hello")
	for i := 0; i < 3; i++ {
		nums[i] = <-chs[i]
		fmt.Println(i, "val", nums[i])
	}
	fmt.Println(nums)

	//nums := make([]int, 3)
	//for i:=0; i<3; i++ {
	//	go world(i, nums)
	//}
	//fmt.Println("hello")
	//time.Sleep(5*time.Second)
	//for i:=0; i<3; i++ {
	//	fmt.Println(i, "val", nums[i])
	//}
	//fmt.Println(nums)

	//ch := make(chan int)
	//go hello(1, ch)
	//if res, ok := <-ch; ok {
	//	fmt.Println(res)
	//}

}

func world(i int, ints []int) {
	ints[i] = i * i
}

func hello(i int, ch chan<- int) {
	ch <- i * i

}

func Test_convert(t *testing.T) {
	var a interface{} = 10
	//fmt.Println(int(a))
	s, ok := a.(int)
	if ok {
		fmt.Println("int", s)
	}
	s2, ok := a.(float32)
	if ok {
		fmt.Println("float32", s2)
	}
}

func Test_addstr(t *testing.T) {
	downstreamInstance := "127.0.0.1:8090"
	queryString := `from(bucket: "test")
	 |> range(start: -15s)
	 |> filter(fn: (r) => r["_measurement"] == "interact")
	 |> filter(fn: (r) => r["downstreamInstance"] == "` + downstreamInstance + `")
	 |> count()`
	fmt.Println(queryString)
}

func Test_float(t *testing.T) {
	a := 1 * 0.5
	fmt.Println(a, reflect.TypeOf(a))
}

func TestRedis(t *testing.T) {
	rdb := utils.InitRedisClient()
	ctx := context.Background()
	rdb.Set(ctx, "name", "pk", time.Hour*1)
	val, _ := rdb.Get(ctx, "name").Result()
	fmt.Println(val)
}

func TestDefer(t *testing.T) {
	fmt.Println("pk")
	defs()
	fmt.Println("?")
	defer func() {
		fmt.Println("sha")
	}()
}

func defs() {
	fmt.Println("hello")
	defer func() {
		fmt.Println("world")
	}()
}
