# golang 之 expvar 使用

标签（空格分隔）： golang expvar

---

## 源码说明

在 `expvar.go` 中有如下说明：

> expvar 包提供了一种标准化接口用于**公共变量**，例如针对 server 中的操作计数器；
> expvar 以 JSON 格式通过 HTTP 的 `/debug/vars` 来暴露这些变量；
> 针对这些公共变量的 set 或 modify 操作具有原子性；
>
> 除了会添加 HTTP handler 以外，该包还会注册如下变量：
> 
> - `cmdline`   os.Args
> - `memstats`  runtime.Memstats
>
> 有些时候导入该包的目的仅仅是为了“副作用”，即注册其提供的 HTTP handler 和上述变量；
> 可以按照如下方式导入该包：
> `import _ "expvar"`
>

产生“副作用”的代码：

```golang
// Do calls f for each exported variable.
// The global variable map is locked during the iteration,
// but existing entries may be concurrently updated.
func Do(f func(KeyValue)) {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, k := range varKeys {
		f(KeyValue{k, vars[k]})
	}
}

// 针对 /debug/vars 的回调函数
func expvarHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(w, "{\n")
	first := true
	Do(func(kv KeyValue) {
		if !first {
			fmt.Fprintf(w, ",\n")
		}
		first = false
		fmt.Fprintf(w, "%q: %s", kv.Key, kv.Value)
	})
	fmt.Fprintf(w, "\n}\n")
}

// 获取命令参数
func cmdline() interface{} {
	return os.Args
}

// 获取内存统计信息
func memstats() interface{} {
	stats := new(runtime.MemStats)
	runtime.ReadMemStats(stats)
	return *stats
}

// 包 expvar 默认提供的东东
func init() {
	http.HandleFunc("/debug/vars", expvarHandler)
	Publish("cmdline", Func(cmdline))
	Publish("memstats", Func(memstats))
}
```

一个简单的使用示例：

```golang
package main

import "expvar"
import "net"
import "fmt"
import "net/http"

var (
    counts = expvar.NewMap("counters")
)

func init() {
    counts.Add("a", 10)
    counts.Add("b", 10)
}

func main() {
    sock, err := net.Listen("tcp", "localhost:8123")
    if err != nil {
        panic("sock error")
    }
    go func() {
        fmt.Println("HTTP now available at port 8123")
        http.Serve(sock, nil)
    }()
    fmt.Println("hello")
    select {}
}
```

启动示例 server

```shell
➜  Golang ./expvar_usage_1
hello
HTTP now available at port 8123

```

获取由 expvar 提供的内容

```shell
➜  ~ curl http://localhost:8123/debug/vars
{
"cmdline": ["./expvar_usage_1"],
"counters": {"a": 10, "b": 10},
"memstats": {"Alloc":321496,"TotalAlloc":321496,"Sys":3084288,"Lookups":5,"Mallocs":4743,"Frees":98,"HeapAlloc":321496,"HeapSys":1769472,"HeapIdle":1040384,"HeapInuse":729088,"HeapReleased":0,"HeapObjects":4645,"StackInuse":327680,"StackSys":327680,"MSpanInuse":14080,"MSpanSys":16384,"MCacheInuse":4800,"MCacheSys":16384,"BuckHashSys":2357,"GCSys":131072,"OtherSys":820939,"NextGC":4194304,"LastGC":0,"PauseTotalNs":0,"PauseNs":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"PauseEnd":[0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0],"NumGC":0,"GCCPUFraction":0,"EnableGC":true,"DebugGC":false,"BySize":[{"Size":0,"Mallocs":0,"Frees":0},{"Size":8,"Mallocs":15,"Frees":0},{"Size":16,"Mallocs":333,"Frees":0},{"Size":32,"Mallocs":3926,"Frees":0},{"Size":48,"Mallocs":101,"Frees":0},{"Size":64,"Mallocs":30,"Frees":0},{"Size":80,"Mallocs":28,"Frees":0},{"Size":96,"Mallocs":8,"Frees":0},{"Size":112,"Mallocs":4,"Frees":0},{"Size":128,"Mallocs":7,"Frees":0},{"Size":144,"Mallocs":3,"Frees":0},{"Size":160,"Mallocs":16,"Frees":0},{"Size":176,"Mallocs":7,"Frees":0},{"Size":192,"Mallocs":3,"Frees":0},{"Size":208,"Mallocs":24,"Frees":0},{"Size":224,"Mallocs":0,"Frees":0},{"Size":240,"Mallocs":2,"Frees":0},{"Size":256,"Mallocs":9,"Frees":0},{"Size":288,"Mallocs":18,"Frees":0},{"Size":320,"Mallocs":3,"Frees":0},{"Size":352,"Mallocs":9,"Frees":0},{"Size":384,"Mallocs":0,"Frees":0},{"Size":416,"Mallocs":26,"Frees":0},{"Size":448,"Mallocs":0,"Frees":0},{"Size":480,"Mallocs":1,"Frees":0},{"Size":512,"Mallocs":0,"Frees":0},{"Size":576,"Mallocs":7,"Frees":0},{"Size":640,"Mallocs":2,"Frees":0},{"Size":704,"Mallocs":6,"Frees":0},{"Size":768,"Mallocs":0,"Frees":0},{"Size":896,"Mallocs":6,"Frees":0},{"Size":1024,"Mallocs":6,"Frees":0},{"Size":1152,"Mallocs":5,"Frees":0},{"Size":1280,"Mallocs":0,"Frees":0},{"Size":1408,"Mallocs":0,"Frees":0},{"Size":1536,"Mallocs":0,"Frees":0},{"Size":1664,"Mallocs":6,"Frees":0},{"Size":2048,"Mallocs":16,"Frees":0},{"Size":2304,"Mallocs":4,"Frees":0},{"Size":2560,"Mallocs":0,"Frees":0},{"Size":2816,"Mallocs":0,"Frees":0},{"Size":3072,"Mallocs":0,"Frees":0},{"Size":3328,"Mallocs":4,"Frees":0},{"Size":4096,"Mallocs":2,"Frees":0},{"Size":4608,"Mallocs":0,"Frees":0},{"Size":5376,"Mallocs":5,"Frees":0},{"Size":6144,"Mallocs":2,"Frees":0},{"Size":6400,"Mallocs":0,"Frees":0},{"Size":6656,"Mallocs":1,"Frees":0},{"Size":6912,"Mallocs":0,"Frees":0},{"Size":8192,"Mallocs":0,"Frees":0},{"Size":8448,"Mallocs":0,"Frees":0},{"Size":8704,"Mallocs":0,"Frees":0},{"Size":9472,"Mallocs":0,"Frees":0},{"Size":10496,"Mallocs":0,"Frees":0},{"Size":12288,"Mallocs":0,"Frees":0},{"Size":13568,"Mallocs":0,"Frees":0},{"Size":14080,"Mallocs":0,"Frees":0},{"Size":16384,"Mallocs":0,"Frees":0},{"Size":16640,"Mallocs":0,"Frees":0},{"Size":17664,"Mallocs":0,"Frees":0}]}
}
➜  ~
```

## [一个golang http包自带的绝佳示例](http://studygolang.com/articles/4105)

这个示例完整展现了 expvar 包中可以使用的各种数据类型，示例代码如下

```golang
package main

import (
	"bytes"
	"expvar"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

// hello world, the web server
var helloRequests = expvar.NewInt("hello-requests")

func HelloServer(w http.ResponseWriter, req *http.Request) {
	helloRequests.Add(1)
	io.WriteString(w, "hello, world!\n")
}

// Simple counter server. POSTing to it will set the value.
type Counter struct {
	mu sync.Mutex // protects n
	n  int
}

// This makes Counter satisfy the expvar.Var interface, so we can export
// it directly.
func (ctr *Counter) String() string {
	ctr.mu.Lock()
	defer ctr.mu.Unlock()
	return fmt.Sprintf("%d", ctr.n)
}

func (ctr *Counter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctr.mu.Lock()
	defer ctr.mu.Unlock()
	switch req.Method {
	case "GET":
		ctr.n++
	case "POST":
		buf := new(bytes.Buffer)
		io.Copy(buf, req.Body)
		body := buf.String()
		if n, err := strconv.Atoi(body); err != nil {
			fmt.Fprintf(w, "bad POST: %v\nbody: [%v]\n", err, body)
		} else {
			ctr.n = n
			fmt.Fprint(w, "counter reset\n")
		}
	}
	fmt.Fprintf(w, "counter = %d\n", ctr.n)
}

// simple flag server
var booleanflag = flag.Bool("boolean", true, "another flag for testing")

func FlagServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Flags:\n")
	flag.VisitAll(func(f *flag.Flag) {
		if f.Value.String() != f.DefValue {
			fmt.Fprintf(w, "%s = %s [default = %s]\n", f.Name, f.Value.String(), f.DefValue)
		} else {
			fmt.Fprintf(w, "%s = %s\n", f.Name, f.Value.String())
		}
	})
}

// simple argument server
func ArgServer(w http.ResponseWriter, req *http.Request) {
	for _, s := range os.Args {
		fmt.Fprint(w, s, " ")
	}
}

// a channel (just for the fun of it)
type Chan chan int

func ChanCreate() Chan {
	c := make(Chan)
	go func(c Chan) {
		for x := 0; ; x++ {
			c <- x
		}
	}(c)
	return c
}

func (ch Chan) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, fmt.Sprintf("channel send #%d\n", <-ch))
}

// exec a program, redirecting output
func DateServer(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8")

	date, err := exec.Command("/bin/date").Output()
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	rw.Write(date)
}

func Logger(w http.ResponseWriter, req *http.Request) {
	log.Print(req.URL)
	http.Error(w, "oops", 404)
}

var webroot = flag.String("root", os.Getenv("HOME"), "web root directory")

func main() {
	flag.Parse()

	// The counter is published as a variable directly.
	ctr := new(Counter)
	expvar.Publish("counter", ctr)
	http.Handle("/counter", ctr)
	http.Handle("/", http.HandlerFunc(Logger))
	http.Handle("/go/", http.StripPrefix("/go/", http.FileServer(http.Dir(*webroot))))
	http.Handle("/chan", ChanCreate())
	http.HandleFunc("/flags", FlagServer)
	http.HandleFunc("/args", ArgServer)
	http.HandleFunc("/go/hello", HelloServer)
	http.HandleFunc("/date", DateServer)
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Panicln("ListenAndServe:", err)
	}
}
```

测试

```
➜  ~ curl http://localhost:12345/counter
counter = 1
➜  ~ curl http://localhost:12345/counter
counter = 2
➜  ~ curl http://localhost:12345/counter
counter = 3
➜  ~
➜  ~ curl -XPOST http://localhost:12345/counter -d '100'
counter reset
counter = 100
➜  ~ curl http://localhost:12345/counter
counter = 101
➜  ~ curl http://localhost:12345/counter
counter = 102
➜  ~ curl -XPOST http://localhost:12345/counter -d '100a'
bad POST: strconv.ParseInt: parsing "100a": invalid syntax
body: [100a]
counter = 102
➜  ~
➜  ~ curl http://localhost:12345/counter
counter = 103
➜  ~
➜  ~ curl http://localhost:12345/
oops
➜  ~
➜  ~ curl http://localhost:12345/go
<a href="/go/">Moved Permanently</a>.

➜  ~
➜  ~ curl -L http://localhost:12345/go
<pre>
<a href=".bash_history">.bash_history</a>
<a href=".bash_profile">.bash_profile</a>
<a href=".bash_sessions/">.bash_sessions/</a>
<a href="http_-L.pcap">http_-L.pcap</a>
<a href="http_-d.pcap">http_-d.pcap</a>
<a href="http_-e.pcap">http_-e.pcap</a>
<a href="http_-r.pcap">http_-r.pcap</a>
<a href="workspace/">workspace/</a>
</pre>
➜  ~
➜  ~ curl http://localhost:12345/chan
channel send #0
➜  ~ curl http://localhost:12345/chan
channel send #1
➜  ~ curl http://localhost:12345/chan
channel send #2
➜  ~
➜  ~ curl http://localhost:12345/flags
Flags:
boolean = true
root = /Users/sunfei
➜  ~
➜  ~ curl http://localhost:12345/args
./expvar_usage_2 %                                                                                                  ➜  ~
➜  ~ curl http://localhost:12345/go/hello
hello, world!
➜  ~
➜  ~ curl http://localhost:12345/date
2017年 1月12日 星期四 17时16分46秒 CST
➜  ~ curl http://localhost:12345/date
2017年 1月12日 星期四 17时16分51秒 CST
➜  ~ curl http://localhost:12345/date
2017年 1月12日 星期四 17时16分56秒 CST
➜  ~ curl http://localhost:12345/date
2017年 1月12日 星期四 17时16分58秒 CST
➜  ~
➜  ~ curl http://localhost:12345/xx
oops
➜  ~ curl http://localhost:12345/yy
oops
➜  ~
```

## packetbeat 中的使用

在运行 packetbeat 时，console 上会周期性输出如下日志

```shell
2017/01/16 02:40:21.957781 logp.go:230: INFO Non-zero metrics in the last 30s: redis.unmatched_responses=1 libbeat.publisher.published_events=7
```

对应代码如下

```golang
// logExpvars 函数会以 Info 日志级别记录 integer 类型的 expvars 值在
// 指定 interval 内的变化情况；对于每一个 expvar 来说，变化量 delta 
// 是从 interval 的开始进行记录和计算的；
func logExpvars(metricsCfg *LoggingMetricsConfig) {
    // 没有配置或配置为 false
	if metricsCfg.Enabled != nil && *metricsCfg.Enabled == false {
		Info("Metrics logging disabled")
		return
	}
	// 默认 interval 设置为 30s
	if metricsCfg.Period == nil {
		metricsCfg.Period = &defaultMetricsPeriod
	}
	Info("Metrics logging every %s", metricsCfg.Period)

	ticker := time.NewTicker(*metricsCfg.Period)
	prevVals := map[string]int64{}
	for {
		<-ticker.C
		vals := map[string]int64{}
		// 将注册到 expvar 中的全部 Int 类型内容保存到 vals 中
		snapshotExpvars(vals)
		// 构建 interval 时间内的所有 Int 类型 expvar 变量的 delta 差值字符串
		metrics := buildMetricsOutput(prevVals, vals)
		prevVals = vals
		if len(metrics) > 0 {
			Info("Non-zero metrics in the last %s:%s", metricsCfg.Period, metrics)
		} else {
		    // 指定 interval 内所有 Int 类型的 expvar 变量均无变化
			Info("No non-zero metrics in the last %s", metricsCfg.Period)
		}
	}
}
```

函数 `snapshotExpvars` 的实现如下

```golang
// snapshotExpvars 对定义的全部 expvars 进行递归处理；针对 integer 类型
// 的内容，会对其 name 和 value 执行 snapshot 操作，保存到一个单独的
// (flat) map 中
func snapshotExpvars(varsMap map[string]int64) {
    // 遍历当前已全局注册过的 expvar 变量
	expvar.Do(func(kv expvar.KeyValue) {
		switch kv.Value.(type) {
		// 如果是 Int 类型，直接保存到 map 中
		case *expvar.Int:
			varsMap[kv.Key], _ = strconv.ParseInt(kv.Value.String(), 10, 64)
		// 如果是 Map 类型，则递归处理其中的内容（其实是处理其中的 Int 类型内容）
		case *expvar.Map:
			snapshotMap(varsMap, kv.Key, kv.Value.(*expvar.Map))
		}
	})
}
```

函数 `buildMetricsOutput` 的实现如下

```golang
// buildMetricsOutput 负责计算 vals 和 prevVals 的差值 delta ；并构建
// 一个便于打印输出 delta 内容的字符串；
// 注：至少存在一个 expvar 变量的 delta 不为 0 才有非空字符串返回
func buildMetricsOutput(prevVals map[string]int64, vals map[string]int64) string {
	metrics := ""
	for k, v := range vals {
		delta := v - prevVals[k]
		if delta != 0 {
			metrics = fmt.Sprintf("%s %s=%d", metrics, k, delta)
		}
	}
	return metrics
}
```

在结束 packetbeat 运行时，console 上会输出如下日志

```shell
2017/01/16 10:41:44.727504 logp.go:245: INFO Total non-zero values:  libbeat.publisher.published_events=12915 redis.unmatched_responses=23 tcp.dropped_because_of_gaps=4
2017/01/16 10:41:44.727537 logp.go:246: INFO Uptime: 1.22785826s
```

对应的代码如下

```golang
func LogTotalExpvars(cfg *Logging) {
	if cfg.Metrics.Enabled != nil && *cfg.Metrics.Enabled == false {
		return
	}
	vals := map[string]int64{}
	prevVals := map[string]int64{}
	snapshotExpvars(vals)
	metrics := buildMetricsOutput(prevVals, vals)
	Info("Total non-zero values: %s", metrics)
	Info("Uptime: %s", time.Now().Sub(startTime))
}
```


## 其他

- [Some notes on Go's expvar package](https://utcc.utoronto.ca/~cks/space/blog/programming/GoExpvarNotes)

- [A surprise to watch out for with Go's expvar package (in expvar.Var)](https://utcc.utoronto.ca/~cks/space/blog/programming/GoExpvarVarGotcha)


