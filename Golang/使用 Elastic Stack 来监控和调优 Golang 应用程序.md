# 使用 Elastic Stack 来监控和调优 Golang 应用程序

> 原文地址：[这里](https://my.oschina.net/u/569210/blog/852351)（部分内容有删减和调整）


一个 Golang 程序开发好了之后，势必要关心其运行情况，今天就介绍一下如何基于 Elastic Stack 分析 Golang 程序的内存使用情况，从而方便针对 Golang 程序做长期监控、性能调优，以及发现一些潜在的内存泄露等问题。

## Metricbeat 概述

Elastic Stack 其实是一个集合，今天主要使用 Elasticsearch、Metricbeat 和 Kibana ；

Metricbeat 是一个专门用来获取**服务器**或**应用服务**内部运行指标数据的收集程序，基于 Golang 写成，部署包才 10M 左右，对目标服务器的部署环境也没有依赖，内存资源占用和 CPU 开销也较小，目前除了可以监控服务器本身的资源使用情况外，还支持常见的应用服务器和服务，目前支持列表如下：

- Apache Module
- Couchbase Module
- Docker Module
- HAProxy Module
- Kafka Module
- MongoDB Module
- MySQL Module
- Nginx Module
- PostgreSQL Module
- Prometheus Module
- Redis Module
- System Module
- ZooKeeper Module

当然，你的应用有可能不在上述列表，不过没关系，Metricbeat 是可以扩展的，你可以很方便的实现一个扩展模块；本文接下来所使用的 Golang Module 就是我刚刚为 Metricbeat 添加的扩展模块，目前已经 merge 进入 Metricbeat 的 master 分支，预计会在 6.0 版本发布，想了解是如何扩展这个模块的可以查看[代码路径](https://github.com/elastic/beats/tree/master/metricbeat/module/golang)和 [PR 地址](https://github.com/elastic/beats/pull/3536)。

## 在 Kibana 上进行 Metricbeat 数据展示

上面这些描述可能还不够吸引人，我们可以 Kibana 上看一下 Metricbeat 基于 Golang 扩展模块收集数据的可视化展示；

![Golang GC example - 1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Golang%20GC%20example%20-%201.png "Golang GC example - 1")

上图解读：

- **Golang: Heap Summary** ：位于最上方的是**堆内存**摘要信息，可用于大致了解 Golang 的内存使用和 GC 情况；其中
    - **System Total Memory** 表示 Golang 程序从操作系统申请的内存，可以理解为进程所占的内存（注意不是进程对应的虚拟内存）；
    - **Bytes Allocated** 表示 Heap 目前分配的内存，也就是 Golang 里面直接可使用的内存；
    - **GC limit** 表示当 Golang 的 Heap 内存分配达到这个 limit 值之后就会开始执行 GC ，这个值会随着每次 GC 而变化；
    - **GC Cycles(count)** 则代表监控周期内的 GC 次数；
- **Golang: Heap**：中间左一，表示堆内存统计情况；
    - **Heap Total**
    - **Heap Inuse**：表示活跃对象大小；
    - **Heap Allocated**：表示正在用和没有用但还未被回收的对象的大小；
    - **Heap Idle**：表示已分配但空闲的内存；
- **Golang: System**：中间左二，表示进程内存统计情况；
    - **System Total**
    - **System Obtained**
    - **System Stack**
    - **System Released**
- **Golang: Objects**：中间左三，表示对象统计情况；
    - **Object Count(avg)**
    - **Allocation Rate**
- **Golang: GC durations** GC 时间；
    - **sum of GC Pause durations(ns)**
    - **Total GC Pause(ns) Rate**
    - **Max GC Pause(ns)**
    - **Avg GC Pause(ns)**
    - **GC Pause count**
- **Golang: GC count**：GC 次数；
    - **GC Count**
    - **GC Rate**
    - **CPU Fraction**：表示该进程 CPU 占用时间花在 GC 上面的百分比，值越大说明 GC 越频繁，浪费更多的时间在 GC 上面，上图虽然趋势陡峭，但是看范围在 0.41%~0.52% 之间，看起来还算可以，如果 GC 比率占到个位数甚至更多比例，那肯定需要进一步优化程序了。

> ⚠️ `mvavg` 是 movingaverage 的简写，即移动平均数，计算方式为用指定的窗口大小计算移动平均值；

有了这些信息，我们就能够知道 Golang 程序的内存使用、分配情况，以及 GC 的执行情况了；

假如要分析是否有内存泄露，看**内存使用**和**堆内存分配**的趋势是否平稳就可以了，另外，如果 **GC Limit** 和 **Bytes Allocated** 一直上升的话，那么肯定有内存泄露了；结合历史信息还能针对不同版本或提交进行 Golang 内存使用分析以及 GC 分析。

## 分析用的数据源

具体如下：首先启用 Golang 的 `expvar` 服务，方法很简单，只需要在 Golang 的程序通过 `import` 引入该包即可，它会自动注册到现有的 http 服务上，如果 Golang 没有启动 http 服务，也可以使用下面的方式启动即可；

> [expvar](https://golang.org/pkg/expvar/) 是 Golang 提供的一个用于暴露内部变量或统计信息的标准包。

```golang
package main

import (
     "fmt"
     "net/http"
     "expvar"
)

func metricsHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")

    first := true
    report := func(key string, value interface{}) {
        if !first {
            fmt.Fprintf(w, ",\n")
        }
        first = false
        if str, ok := value.(string); ok {
            fmt.Fprintf(w, "%q: %q", key, str)
        } else {
            fmt.Fprintf(w, "%q: %v", key, value)
        }
    }

    fmt.Fprintf(w, "{\n")
    expvar.Do(func(kv expvar.KeyValue) {
        report(kv.Key, kv.Value)
    })
    fmt.Fprintf(w, "\n}\n")
}

func main() {
   mux := http.NewServeMux()
   mux.HandleFunc("/debug/vars", metricsHandler)
   http.ListenAndServe("localhost:6060", mux)
}
```

上述代码注册的访问路径为 `/debug/vars`，编译之后启动，就可以通过 [`http://localhost:6060/debug/vars`](http://localhost:6060/debug/vars) 来访问 `expvar` 以 JSON 格式暴露出来内部变量；`expvar` 默认提供了 `runtime.Memstats` 信息，也就是上图分析的数据源，当然你还可以注册自己的变量；

![memstats](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/golang%20example%20memstats.png "memstats")

OK，现在我们的 Golang 程序已经启动了，并且通过 `expvar` 暴露出了运行时的内存使用情况，现在我们需要使用 Metricbeat 来获取这些信息并存进 Elasticsearch。

## Metricbeat 安装配置

Metricbeat 的安装可以基于[官方提供的安装包](https://www.elastic.co/downloads/beats/metricbeat)，或直接[在 github 上下载源码](https://github.com/elastic/beats)后编译安装；
启动 Metricbeat 前，需要修改配置文件 `metricbeat.yml` 的内容：

> 注意：由于目前获取 Golang 内存信息的 Metricbeat module 只合并到了 master 分支，因此只能在相应版本基础上测试该功能，即 6.0.0-alpha1 和 github 上 master 分支；

```
metricbeat.modules:
  - module: golang
     metricsets: ["heap"]
     enabled: true
     period: 10s
     hosts: ["localhost:6060"]
     heap.path: "/debug/vars"
output.elasticsearch:
  hosts: ["localhost:9200"]
```

上面的参数设置用于

- 启用 Golang 监控模块；
- 每 10 秒从配置路径位置（hosts + heap.path）获取一次 `expvar` 数据；
- 设置数据输出为本机的 Elasticsearch ；


## Metricbeat 数据处理

```
➜  metricbeat git:(master) ✗ ./metricbeat -e -v
2017/05/23 09:47:04.643162 beat.go:334: INFO Home path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/metricbeat] Config path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/metricbeat] Data path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/metricbeat/data] Logs path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/metricbeat/logs]
2017/05/23 09:47:04.643194 beat.go:359: INFO Beat metadata path: /Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/elastic/beats/metricbeat/data/meta.json
2017/05/23 09:47:04.643210 metrics.go:23: INFO Metrics logging every 30s
2017/05/23 09:47:04.643585 beat.go:341: INFO Beat UUID: a3cb6ebe-845d-40fb-b381-77a28667ecc6
2017/05/23 09:47:04.643607 beat.go:214: INFO Setup Beat: metricbeat; Version: 6.0.0-alpha2
2017/05/23 09:47:04.644775 client.go:131: INFO Elasticsearch url: http://localhost:9200
2017/05/23 09:47:04.645047 outputs.go:107: INFO Activated elasticsearch as output plugin.
2017/05/23 09:47:04.645602 publish.go:191: INFO Publisher name: sunfeideMacBook-Pro.local
2017/05/23 09:47:04.645955 async.go:63: INFO Flush Interval set to: 1s
2017/05/23 09:47:04.645970 async.go:64: INFO Max Bulk Size set to: 50
2017/05/23 09:47:04.646083 metricbeat.go:31: INFO Register [ModuleFactory:[docker, mongodb, mysql, postgresql, system], MetricSetFactory:[apache/status, audit/kernel, ceph/cluster_disk, ceph/cluster_health, ceph/monitor_health, ceph/pool_disk, couchbase/bucket, couchbase/cluster, couchbase/node, docker/container, docker/cpu, docker/diskio, docker/healthcheck, docker/image, docker/info, docker/memory, docker/network, dropwizard/collector, elasticsearch/node, elasticsearch/node_stats, golang/expvar, golang/heap, haproxy/info, haproxy/stat, http/json, jolokia/jmx, kafka/consumergroup, kafka/partition, kibana/status, kubernetes/container, kubernetes/node, kubernetes/pod, kubernetes/system, kubernetes/volume, memcached/stats, mongodb/dbstats, mongodb/status, mysql/status, nginx/stubstatus, php_fpm/pool, postgresql/activity, postgresql/bgwriter, postgresql/database, prometheus/collector, prometheus/stats, redis/info, redis/keyspace, system/core, system/cpu, system/diskio, system/filesystem, system/fsstat, system/load, system/memory, system/network, system/process, vsphere/datastore, vsphere/host, vsphere/virtualmachine, zookeeper/mntr]]
2017/05/23 09:47:04.647145 log.go:144: WARN EXPERIMENTAL: The golang heap metricset is experimental
2017/05/23 09:47:04.647247 beat.go:262: INFO metricbeat start running.
2017/05/23 09:47:05.655894 client.go:653: INFO Connected to Elasticsearch version 5.3.2
2017/05/23 09:47:05.662934 load.go:50: INFO Loading template for elasticsearch version: 5.3.2
2017/05/23 09:47:05.663249 beat.go:495: INFO ES template successfully loaded.
2017/05/23 09:47:34.647791 metrics.go:39: INFO Non-zero metrics in the last 30s: beat.memstats.gc_next=4473924 beat.memstats.memory_alloc=2874456 beat.memstats.memory_total=2874456 metricbeat.golang.heap.events=3 metricbeat.golang.heap.success=3 output.elasticsearch.events.acked=3 output.elasticsearch.publishEvents.call.count=3 output.elasticsearch.read.bytes=1507 output.elasticsearch.write.bytes=3164 output.events.acked=3 output.write.bytes=3164 publisher.events.count=3 publisher.queue.messages.count=3
2017/05/23 09:48:04.645437 metrics.go:39: INFO Non-zero metrics in the last 30s: beat.memstats.memory_alloc=309104 beat.memstats.memory_total=309104 metricbeat.golang.heap.events=3 metricbeat.golang.heap.success=3 output.elasticsearch.events.acked=3 output.elasticsearch.publishEvents.call.count=3 output.elasticsearch.read.bytes=1046 output.elasticsearch.write.bytes=2836 output.events.acked=3 output.write.bytes=2836 publisher.events.count=3 publisher.queue.messages.count=3
2017/05/23 09:48:34.644285 metrics.go:39: INFO Non-zero metrics in the last 30s: beat.memstats.memory_alloc=292592 beat.memstats.memory_total=292592 metricbeat.golang.heap.events=3 metricbeat.golang.heap.success=3 output.elasticsearch.events.acked=3 output.elasticsearch.publishEvents.call.count=3 output.elasticsearch.read.bytes=1043 output.elasticsearch.write.bytes=2842 output.events.acked=3 output.write.bytes=2842 publisher.events.count=3 publisher.queue.messages.count=3
...
```

从上述信息中可以看到

```
...
2017/05/23 09:47:04.646083 metricbeat.go:31: INFO Register [ModuleFactory:[..., golang/expvar, golang/heap, ...]]
2017/05/23 09:47:04.647145 log.go:144: WARN EXPERIMENTAL: The golang heap metricset is experimental
...
```

进一步确认了模块已经生效，但为实验性模块；


从上面的日志输出中可以看到，已经有数据发送到 Elasticsearch 中了（当然要确保在运行 metricbeat 前 Elasticsearch 和 Kibana 是可用状态）；你可以在 Kibana 中根据需要灵活地自定义可视化展示，推荐使用 `Timelion` 进行分析；当然，为了方便也可以直接导入官方提供的样例仪表板，即上面第一个图的效果。

关于如何导入样例仪表板请参照《[Loading Sample Kibana Dashboards](https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-sample-dashboards.html)》； 

除了可以监控默认提供的内存信息外，如果你还有一些内部业务指标想要暴露出来，也是可以通过 `expvar` 实现。一个简单示例如下：

```golang
var inerInt int64 = 1024
pubInt := expvar.NewInt("your_metric_key")
pubInt.Set(inerInt)
pubInt.Add(2)
```

另外，由于 Metricbeat 自身的内部实现也通过 `expvar` 暴露了很多内部运行时信息，所以，Metricbeat 完全可以自己监控自己；

首先，启动的时候需要通过指定 `pprof` 参数设置 pprof http server 地址：

```shell
./metricbeat -httpprof="127.0.0.1:6060" -e -v
```

之后就可以通过 [http://127.0.0.1:6060/debug/vars](http://127.0.0.1:6060/debug/vars) 访问到 metricbeat 的内部运行情况了；

![expvar httpprof output](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/expvar%20httpprof%20output.png "expvar httpprof output")

如上图所示，此时就能看到 output 模块中和 Elasticsearch 相关的 `expvar` 统计变量的情况，例如 `output.elasticsearch.events.acked` 表示发送到 Elasticsearch Ack 返回之后的消息。

现在我们进一步修改 Metricbeat 配置文件，为 Golang module 设置两个 metricset ，即可以理解为设置两个**监控指标类型**；现在需要增加一个新的 `expvar` 类型，即新增一个自定义指标，相应配置文件修改如下：

```
- module: golang
  metricsets: ["heap","expvar"]
  enabled: true
  period: 1s
  hosts: ["localhost:6060"]
  heap.path: "/debug/vars"
  expvar:
    namespace: "metricbeat"
    path: "/debug/vars"
```

参数 `namespace` 表示自定义指标属于一个命令空间，主要是为了方便管理；这里用于 Metricbeat 自身信息，所以 namespace 就是 metricbeat。

重启 Metricbeat 应该就能收到新的数据了；

假设我们关注的是 `output.elasticsearch.events.acked` 和
`output.elasticsearch.events.not_acked` 这两个指标，可以在 Kibana 里面简单**定义一个曲线图**就能看到 Metricbeat 发往 Elasticsearch 消息的**成功和失败趋势**。

Timelion 表达式：

```
.es("metricbeat*",metric="max:golang.metricbeat.output.elasticsearch.events.acked").derivative().label("Elasticsearch Success"),.es("metricbeat*",metric="max:golang.metricbeat.output.elasticsearch.events.not_acked").derivative().label("Elasticsearch Failed")
```

效果如下：

![Golang GC example - 2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Golang%20GC%20example%20-%202.png "Golang GC example - 2")

从上图可以看到，发往 Elasticsearch 的消息很稳定，没有出现丢消息的情况，同时关于 Metricbeat 的内存情况，我们打开导入的 Dashboard 查看：

![Golang GC example - 3](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Golang%20GC%20example%20-%203.png "Golang GC example - 3")



----------


下图是我自己跑测试得到的结果：

- 针对测试小程序

![Golang GC in Kibana_myself](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Golang%20GC%20in%20Kibana_myself.png "Golang GC in Kibana_myself")

- 针对 Metricbeat 自身

![Golang GC in Kibana_Metricbeat](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Golang%20GC%20in%20Kibana_Metricbeat.png "Golang GC in Kibana_Metricbeat")