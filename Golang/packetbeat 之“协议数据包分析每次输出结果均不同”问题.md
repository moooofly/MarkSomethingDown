# packetbeat 之“协议数据包分析每次输出结果均不同”问题

![packetbeat](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/packetbeat.png "packetbeat")

## 问题描述

通过 `packetbeat` 可执行程序进行 redis 协议数据包分析的输出结果每次都有所不同；

测试命令为

```
# ./packetbeat -c ./packetbeat.yml -e -I redis_xg-bjdev-rediscluster-2_prot-7101_20161222110723_20161222110733.pcap -E packetbeat.protocols.redis.ports=7101 -t
```

输出结果如下

```
... // 第一次
2017/01/12 04:18:56.634130 logp.go:245: INFO Total non-zero values:  libbeat.publisher.published_events=13859 redis.unmatched_responses=23 tcp.dropped_because_of_gaps=4
2017/01/12 04:18:56.634143 logp.go:246: INFO Uptime: 1.210802009s
... // 第二次
2017/01/12 04:21:10.543460 logp.go:245: INFO Total non-zero values:  redis.unmatched_responses=23 libbeat.publisher.published_events=14619 tcp.dropped_because_of_gaps=4
2017/01/12 04:21:10.543478 logp.go:246: INFO Uptime: 1.216717998s
... // 第三次
2017/01/12 04:22:49.583149 logp.go:245: INFO Total non-zero values:  libbeat.publisher.published_events=15006 tcp.dropped_because_of_gaps=4 redis.unmatched_responses=23
2017/01/12 04:22:49.583160 logp.go:246: INFO Uptime: 1.628717709s
```

从上面的输出中可以看到：

- 抓包文件中共有 **72257** 个数据包；
- 由 `packetbeat` 统计出的 "published_events" 三次分别为：**13859**、**14619** 和 **15006** ；

能够确认的是：
 
- published_events 的数值和 `logs/` 目录下生成文件中的数据行数是一致的；
- sniffer.go 计算得到的数据包总数和通过 `capinfos` 命令计算得到的数据包数量是一致的；

## 源码分析

针对如下每次发生变化的输出日志，进行代码反查：

```
2017/01/12 04:18:56.634130 logp.go:245: INFO Total non-zero values:  libbeat.publisher.published_events=13859 redis.unmatched_responses=23 tcp.dropped_because_of_gaps=4
2017/01/12 04:18:56.634143 logp.go:246: INFO Uptime: 1.210802009s
```

在 `logp.go` 中

```golang
func LogTotalExpvars(cfg *Logging) {
	if cfg.Metrics.Enabled != nil && *cfg.Metrics.Enabled == false {
		return
	}
	vals := map[string]int64{}
	prevVals := map[string]int64{}
	// 将注册到 expvar 中的全部 Int 类型内容保存到 vals 中
	snapshotExpvars(vals)
	// 构建“从开始运行到结束运行”的整个时间段内
	// 所有 Int 类型 expvar 变量的 delta 差值字符串
	metrics := buildMetricsOutput(prevVals, vals)
	// 输出“问题”打印
	Info("Total non-zero values: %s", metrics)
	// 输出“从开始运行到结束运行”的时间长度
	Info("Uptime: %s", time.Now().Sub(startTime))
}
```

> 小结：输出结果正确的体现了“从开始到结束”的差值计算；

在 `beat.go` 中

```golang
func (b *Beat) launch(bt Creator) error {
    ...
    // 标识 packetbeat 开始运行
	logp.Info("%s start running.", b.Name)
	// 标识 packetbeat 结束运行
	defer logp.Info("%s stopped.", b.Name)
	// 在结束运行之前，输出当前基于 expvar 记录
	// 的 metrics 变化值
	defer logp.LogTotalExpvars(&b.Config.Logging)

	return beater.Run(b)
}
```

> 小结：在运行结束前，将基于 `expvar` 维护的全局计数值进行计算输出；

推演：如果输出过程没有问题，那么只能是计算过程出了问题；

在 `client.go` 中

```golang
...
// Metrics that can retrieved through the expvar web interface.
// 用于计算 publish_events 值的 expvar 变量
var (
	publishedEvents = expvar.NewInt("libbeat.publisher.published_events")
)
...
func (c *client) PublishEvent(event common.MapStr, opts ...ClientOption) bool {
    // 向 event 中添加自定义字段内容
	c.annotateEvent(event)

    // 基于配置的 Processors 进行定制化 event 过滤
    // 由于我没有配置这个，因为不会有 event 被过滤掉
	publishEvent := c.filterEvent(event)
	if publishEvent == nil {
		return false
	}

    // 根据配置获取一种投递 event 的管道
	ctx, pipeline := c.getPipeline(opts)
	// 将 publish_events 统计变量 +1
	publishedEvents.Add(1)
	// 将 event 封装成 message 投递到管道中
	return pipeline.publish(message{
		client:  c,
		context: ctx,
		datum:   outputs.Data{Event: *publishEvent},
	})
}
...
func (c *client) PublishEvents(events []common.MapStr, opts ...ClientOption) bool {
	data := make([]outputs.Data, 0, len(events))
	// 针对 N 个 event 的循环处理
	for _, event := range events {
		c.annotateEvent(event)

		publishEvent := c.filterEvent(event)
		if publishEvent != nil {
			data = append(data, outputs.Data{Event: *publishEvent})
		}
	}

	ctx, pipeline := c.getPipeline(opts)
	if len(data) == 0 {
		logp.Debug("filter", "No events to publish")
		return true
	}

    // 将 publish_events 变量 +N
	publishedEvents.Add(int64(len(data)))
	return pipeline.publish(message{client: c, context: ctx, data: data})
}
...
```

> 小结：针对每个 event 都进行了 +1 操作；

那么谁调用了 `PublishEvent` 和 `PublishEvents` 呢？

在 `publish.go` 中

```golang
...
func (p *PacketbeatPublisher) onTransaction(event common.MapStr) {
    // 确认 event 的有效性，即特定字段校验
	if err := validateEvent(event); err != nil {
		logp.Warn("Dropping invalid event: %v", err)
		return
	}

    // 针对 event 中的地址信息进行统一化处理
	if !p.normalizeTransAddr(event) {
		return
	}

    // 将 event 发布到管道中
	p.client.PublishEvent(event)
}

func (p *PacketbeatPublisher) onFlow(events []common.MapStr) {
	pub := events[:0]
	// 循环处理 N 个 event
	for _, event := range events {
		if err := validateEvent(event); err != nil {
			logp.Warn("Dropping invalid event: %v", err)
			continue
		}

		if !p.addGeoIPToFlow(event) {
			continue
		}

		pub = append(pub, event)
	}

	p.client.PublishEvents(pub)
}
...
```

> 小结：上述代码没进行任何特别处理；

```golang
func (p *PacketbeatPublisher) Start() {
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-p.done:
				return
			// 从名为 trans 的 channel 获取一个 event
			case event := <-p.trans:
				p.onTransaction(event)
			}
		}
	}()

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		for {
			select {
			case <-p.done:
				return
			// 从名为 flows 的 channel 获取 N 个 event
			case events := <-p.flows:
				p.onFlow(events)
			}
		}
	}()
}
```

> 小结：这里似乎可能会出问题，因为只要是基于 channel 传递内容，就无法避免 buffer 长度的问题；

相应代码如下

```golang
type PacketbeatPublisher struct {
    ....
	trans chan common.MapStr
	flows chan []common.MapStr
}
...
func NewPublisher(
	pub publisher.Publisher,
	hwm, bulkHWM int,
	ignoreOutgoing bool,
) (*PacketbeatPublisher, error) {
    ...
	return &PacketbeatPublisher{
		pub:            pub,
		topo:           topo,
		geoLite:        topo.GeoLite(),
		ignoreOutgoing: ignoreOutgoing,
		client:         pub.Connect(),
		done:           make(chan struct{}),
		// trans channel 的 buffer 长度为 hwm
		trans:          make(chan common.MapStr, hwm),
		// flows channel 的 buffer 长度为 bulkHWM
		flows:          make(chan []common.MapStr, bulkHWM),
	}, nil
}
```

在 `packetbeat.go` 中

```golang
// init packetbeat components
func (pb *packetbeat) init(b *beat.Beat) error {
    ...
    // This is required as init Beat is called before the beat publisher is initialised
	b.Config.Shipper.InitShipperConfig()
	
    // hwm 即 QueueSize 的值；
    // bulkHWM 即 BulkQueueSize 的值；
	pb.pub, err = publish.NewPublisher(b.Publisher, *b.Config.Shipper.QueueSize, *b.Config.Shipper.BulkQueueSize, pb.config.IgnoreOutgoing)
	if err != nil {
		return fmt.Errorf("Initializing publisher failed: %v", err)
	}
	...
}
```

在 `publisher/publish.go` 中可以看到相应定义

```golang
type ShipperConfig struct {
    ...
	// internal publisher queue sizes
	QueueSize     *int `config:"queue_size"`
	BulkQueueSize *int `config:"bulk_queue_size"`
    ...
}
...
// 默认值
const (
	DefaultQueueSize     = 1000
	DefaultBulkQueueSize = 0
)
...
// 初始化函数
func (config *ShipperConfig) InitShipperConfig() {

	// TODO: replace by ucfg
	// QueueSize 的值在 packetbeat.yml 中定义
	if config.QueueSize == nil || *config.QueueSize <= 0 {
		queueSize := DefaultQueueSize
		config.QueueSize = &queueSize
	}

	if config.BulkQueueSize == nil || *config.BulkQueueSize < 0 {
		bulkQueueSize := DefaultBulkQueueSize
		config.BulkQueueSize = &bulkQueueSize
	}
}
```

在搞清楚了配置位置后，就剩下最后一个问题，trans 和 flows channel 中的消息来自哪里；

```golang
func (p *PacketbeatPublisher) PublishTransaction(event common.MapStr) bool {
	select {
	case p.trans <- event:
		return true
	default:
		// drop event if queue is full
		// 这个注释说明很关键：如果 queue 满了，event 会被丢弃
		return false
	}
}

func (p *PacketbeatPublisher) PublishFlows(event []common.MapStr) bool {
	select {
	case p.flows <- event:
		return true
	case <-p.done:
		// drop event, if worker has been stopped
		return false
	}
}
```

在 `redis.go` 中

```golang
// 将 request 和 response 进行关联
func (redis *redisPlugin) correlate(conn *redisConnectionData) {
	// drop responses with missing requests
	if conn.requests.empty() {
		for !conn.responses.empty() {
			debugf("Response from unknown transaction. Ignoring")
			unmatchedResponses.Add(1)
			conn.responses.pop()
		}
		return
	}

	// merge requests with responses into transactions
	for !conn.responses.empty() && !conn.requests.empty() {
		requ := conn.requests.pop()
		resp := conn.responses.pop()

		if redis.results != nil {
		    // 构建 transaction 消息内容（JSON 格式）
			event := redis.newTransaction(requ, resp)
			// 将 event 发送到 trans channel 中
			redis.results.PublishTransaction(event)
		}
	}
}
```

> 小结：从上面的 `PublishTransaction` 和 `PublishFlows` 函数实现中能够看到：导致输出结果每次不同的原因就是由于存在 event 被 drop 掉的问题；

至此，消息处理流程梳理完毕：

- redis 模块协议解析完成后，进行 request-response 配对，构成 transaction 后（即 event）发布到 trans 或 flows channel 中；
- event 从 channel 中获取出来，经过一系列判定和封装（构建为 message），再发送到 pipeline 中（试验发现使用的是 async pipeline）；
- 发送到 pipeline 中的 message 会以广播的方式发送给该 pipeline 下关联的每一个 worker ；在代码层面 messageWorker 实现了 worker 这个 interface ，因此 message 实际是被发送到了 messageWorker 结构上名为 queue 和 bulkQueue 的 channel 中；
- 作为 goroutine 运行的 messageWorker 不断从上述 channel 取走 message 并触发初始化时注册到 handler 上的 onMessage 回调函数；经确认，初始化过程中注册到 handler 上的为 outputWorker ；
- outputWorker 在拿到 message 后，会根据在 packetbeat.yml 中配置的 output 进行处理；在实际配置中，我只配置了 file 这个 output ，因此最终内容会写入磁盘文件（可配置的 output 包括：console, file, Redis, Kafka, logstash 和 Elasticsearch）；


## 问题原因

- 考虑到性能原因，官方默认配置 `queue_size` 为 **1000** ；需要注意的是，该值对应了 `packetbeat` 内部多种 channel 的 buffer 长度；如果你要处理的 pcap 文件中包数量非常多，则需要根据实际情况调大该值（否则，event 会丢失，进而导致写入文件中的 transaction 丢失）；
- 另外，还需要给 packetbeat 预留出足够长的包分析时间，否则可能出现尚未完成全部包的分析，就进入退出过程的情况（packetbeat stops running, because all packets have been send to the protcol analyzers. Not all events might be published yet）；

## 解决办法

- 将 packetbeat.yml 中的 `queue_size` 值按需要调大；
- 使用 `-waitstop=N` 增加 `packetbeat` 等待文件分析结束的时间；

## 其他

- 在官方论坛上的[讨论](https://discuss.elastic.co/t/why-packetbeat-generates-result-of-libbeat-publisher-published-events-changing-every-time/71699)；
- 给官方提的 [issue](https://github.com/elastic/beats/issues/3404)；

## 遗留问题

- tcp.dropped_because_of_gaps 含义；
- redis.unmatched_responses 含义；
