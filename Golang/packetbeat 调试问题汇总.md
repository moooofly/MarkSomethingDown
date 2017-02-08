# packetbeat 调试问题汇总

标签（空格分隔）： packetbeat 

---

> 本文用于记录 packetbeat 在调试运行中遇到的各种问题；

## #00 "protos.go:100: ERR Unknown protocol plugin: xxx" 错误

### 解决办法

需要在 `packetbeat.yml` 中将不需要的协议插件注释掉；

## #01 调整默认的 Outputs 配置（默认输出到 elasticsearch）

### 解决办法

调整为 `file` 和 `console` 输出；

## #02 "sniffer.go:365: WARN Time in pcap went backwards: 0" 警告

### 问题描述

实际测试过程中会输出以下几种数值：

```shell
➜  packetbeat git:(master) ✗ ./packetbeat -c ./packetbeat.yml -e -I redis_xxxx.pcap 2>&1 |grep "pcap went backwards"
...
2017/01/11 06:57:47.091832 sniffer.go:365: WARN Time in pcap went backwards: -2000
...
2017/01/11 06:57:47.189023 sniffer.go:365: WARN Time in pcap went backwards: -1000
...
2017/01/11 06:57:47.189238 sniffer.go:365: WARN Time in pcap went backwards: 0
```

### 源码分析

在 `sniffer.go` 中，有如下代码

```golang
func (sniffer *SnifferSetup) Run() error {
        ...
		if sniffer.config.File != "" {  // 如果是读取的 pcap 文件
			if lastPktTime != nil && !sniffer.config.TopSpeed { // TopSpeed 对应 -t 选项
		        // 计算前后两个数据包的时间戳之差
				sleep := ci.Timestamp.Sub(*lastPktTime)
				if sleep > 0 {
					time.Sleep(sleep)
				} else {
				    // 发现时间戳有“回退”
					logp.Warn("Time in pcap went backwards: %d", sleep)
				}
			}
			_lastPktTime := ci.Timestamp
			lastPktTime = &_lastPktTime
			// 若没有设置 -t 选项，则使用当前系统时间作为包的时间戳
			if !sniffer.config.TopSpeed {
				ci.Timestamp = time.Now() // overwrite what we get from the pcap
			}
		}
		...
}
```

### 问题原因

相邻两个数据包之间的时间差，是基于包捕获时的时间戳计算得到的；当存在大量 TCP 重传情况时（如下图），则会出现时间戳相同的情况；当有乱序发生时，则会出现时间戳差值为负的情况；

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Seconds%20since%20beginning%20of%20capture.png)

该问题与是否使用 `-t` 选项没有直接关系；

### 解决办法

可以通过设置 `-t` 参数绕过此问题，因为此时不会输出相应的打印内容；

### 补充试验

使用 `-t` 选项的情况：

```shell
...
2017/01/11 07:19:05.673469 logp.go:246: INFO Uptime: 523.788161ms
```

未使用 `-t` 选项的情况：

```shell
...
2017/01/11 07:20:27.991096 logp.go:246: INFO Uptime: 14.847278439s
```

经确认，使用 `-t` 选项的目的在于能够准确的按照包内容进行数据重放（包括时间延迟）；

### 其他

在官方论坛上的[讨论](https://discuss.elastic.co/t/packetbeat-something-wired-with-warn-time-in-pcap-went-backwards/72142)；

## #03 在 packetbeat.yml 中移除（通过 '#' 注释掉） [Transaction protocols] 中的内容时出现协议端口错乱问题

> 这个问题比较二，浪费了我好长时间排查问题出在了哪里；

### 问题原因

在注释掉 "packetbeat.protocols.xxx" 内容后，还必须要同时注释掉其下的子配置选项！！！

## #04 移除 flow 相关的 report 配置以减少信息干扰

> 需要确认 flow report 到底是干什么的

## #05 基于 packetbeat 可执行程序分析 redis 协议数据包时每次输出结果都不同

详细说明参见《[packetbeat 之“协议数据包分析每次输出结果均不同”问题](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/packetbeat%20%E4%B9%8B%E2%80%9C%E5%8D%8F%E8%AE%AE%E6%95%B0%E6%8D%AE%E5%8C%85%E5%88%86%E6%9E%90%E6%AF%8F%E6%AC%A1%E8%BE%93%E5%87%BA%E7%BB%93%E6%9E%9C%E5%9D%87%E4%B8%8D%E5%90%8C%E2%80%9D%E9%97%AE%E9%A2%98.md)》

## #06 "redis_parse.go:306: ERR Failed to read integer reply: Expected digit"

### 问题描述

错误信息如下（增加了打印内容）：

```shell
➜  packetbeat git:(master) ✗ ./packetbeat -c ./packetbeat.yml -e -I redis_xg-bjdev-rediscluster-1_prot-7101_20161222110711_20161222110721.pcap -E packetbeat.protocols.redis.ports=7101 -t
2017/01/11 09:42:37.662148 logp.go:219: INFO Metrics logging every 30s
2017/01/11 09:42:37.661865 beat.go:267: INFO Home path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/moooofly/beats/packetbeat] Config path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/moooofly/beats/packetbeat] Data path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/moooofly/beats/packetbeat/data] Logs path: [/Users/sunfei/workspace/GIT/IdeaProjects/src/github.com/moooofly/beats/packetbeat/logs]
2017/01/11 09:42:37.662756 beat.go:177: INFO Setup Beat: packetbeat; Version: 6.0.0-alpha1
2017/01/11 09:42:37.663073 file.go:45: INFO File output path set to: ./logs
2017/01/11 09:42:37.663096 file.go:46: INFO File output base filename set to: packetbeat
2017/01/11 09:42:37.663105 file.go:49: INFO Rotate every bytes set to: 10240000
2017/01/11 09:42:37.663110 file.go:53: INFO Number of files set to: 7
2017/01/11 09:42:37.663140 outputs.go:106: INFO Activated file as output plugin.
2017/01/11 09:42:37.663409 publish.go:291: INFO Publisher name: sunfeideMacBook-Pro.local
2017/01/11 09:42:37.663797 async.go:63: INFO Flush Interval set to: -1s
2017/01/11 09:42:37.663817 async.go:64: INFO Max Bulk Size set to: -1
2017/01/11 09:42:37.664065 procs.go:79: INFO Process matching disabled
2017/01/11 09:42:37.664290 protos.go:89: INFO registered protocol plugin: amqp
2017/01/11 09:42:37.664309 protos.go:89: INFO registered protocol plugin: http
2017/01/11 09:42:37.664315 protos.go:89: INFO registered protocol plugin: mysql
2017/01/11 09:42:37.664320 protos.go:89: INFO registered protocol plugin: redis
2017/01/11 09:42:37.665784 beat.go:207: INFO packetbeat start running.
2017/01/11 09:47:45.218365 redis_parse.go:306: ERR Failed to read integer reply: Expected digit
value should be Int, but we get ->
"{\"id\":153,\"email\":\"xiaojiao.xie@ele.me\",\"work_code\":\"E000027\",\"mobile\":18603720925,\"name\":\"谢小佼\",\"walle_id\":77287,\"status\":6,\"pinyin_name\":\"xxj\",\"sex\":1,\"security_level\":60,\"certificate_type\":0,\"certificate_number\":\"420683198908113733\",\"created_at\":1431550029000,\"updated_at\":1449228237000,\"nchr_id\":\"0001A910000000002EQP\"}}"

2017/01/11 09:42:38.430211 sniffer.go:384: INFO Input finish. Processed 40644 packets. Have a nice day!
2017/01/11 09:42:38.430657 util.go:48: INFO flows worker loop stopped
2017/01/11 09:42:38.430709 logp.go:245: INFO Total non-zero values:  libbeat.publisher.published_events=8080 tcp.dropped_because_of_gaps=15 redis.unmatched_responses=15
2017/01/11 09:42:38.430722 logp.go:246: INFO Uptime: 909.957024ms
2017/01/11 09:42:38.430728 beat.go:211: INFO packetbeat stopped.
➜  packetbeat git:(master) ✗
```

能够看到错误信息为

```shell
2017/01/11 09:47:45.218365 redis_parse.go:306: ERR Failed to read integer reply: Expected digit
value should be Int, but we get ->
"{\"id\":153,\"email\":\"xiaojiao.xie@ele.me\",\"work_code\":\"E000027\",\"mobile\":18603720925,\"name\":\"谢小佼\",\"walle_id\":77287,\"status\":6,\"pinyin_name\":\"xxj\",\"sex\":1,\"security_level\":60,\"certificate_type\":0,\"certificate_number\":\"420683198908113733\",\"created_at\":1431550029000,\"updated_at\":1449228237000,\"nchr_id\":\"0001A910000000002EQP\"}}"
```

### 问题原因

当调用 `HMGET` 时同时查询多个数据，应答包含的数据量比较大时，需要分包进行回复；若在应答回复的过程中，出现丢包，则会导致数据解析出错；

![HMGET 的应答数据分包回复遇到丢包问题](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/HMGET%20%E7%9A%84%E5%BA%94%E7%AD%94%E6%95%B0%E6%8D%AE%E5%88%86%E5%8C%85%E5%9B%9E%E5%A4%8D%E9%81%87%E5%88%B0%E4%B8%A2%E5%8C%85%E9%97%AE%E9%A2%98.png "HMGET 的应答数据分包回复遇到丢包问题")

具体分析如下：

**37642 号包**：发起 `HMGET` 请求，查询 key 为 hr-e0acd6e0-4c21-4917-a676-c4fd8094f2aa:user 下的 8 个 filed 的值；

```
*10
$5
HMGET
$44
hr-e0acd6e0-4c21-4917-a676-c4fd8094f2aa:user
$5
18370
$5
52708
$1
0
$5
18370
$4
1117
$2
13
$3
153
$3
147
```


**37642 号包**：提示 *[TCP Previous segment not captured]* 信息，参考 `HMGET` 请求的内容，结合当前包的数据内容可以看出，该包为 `HMGET` 应答的最后两个 field 的部分内容（并且截断位置也应该是无规律的）；

```
:{"id":153,"email":"xiaojiao.xie@ele.me","work_code":"E000027","mobile":18603720925,"name":".........","walle_id":77287,"status":6,"pinyin_name":"xxj","sex":1,"security_level":60,"certificate_type":0,"certificate_number":"420683198908113733","created_at":1431550029000,"updated_at":1449228237000,"nchr_id":"0001A910000000002EQP"}}
$519
{"userId":147,"userBuList":[3175],"tagsList":[],"userBuRoleDto":[{"id":81524,"bu_id":3175,"bu_name":"............BU","role_id":859,"role_name":".........","user_id":147,"user_name":"......"}],"user":{"id":147,"email":"xin.jin@ele.me","work_code":"E000029","mobile":18607175626,"name":"......","walle_id":56063,"status":6,"pinyin_name":"jx","sex":1,"security_level":70,"certificate_type":0,"certificate_number":"420106198511242510","created_at":1431550030000,"updated_at":1449228115000,"nchr_id":"0001A910000000002ERE"}}
```

数据包细节展示如下：

```
0000   3a 7b 22 69 64 22 3a 31 35 33 2c 22 65 6d 61 69  :{"id":153,"emai
0010   6c 22 3a 22 78 69 61 6f 6a 69 61 6f 2e 78 69 65  l":"xiaojiao.xie
0020   40 65 6c 65 2e 6d 65 22 2c 22 77 6f 72 6b 5f 63  @ele.me","work_c
0030   6f 64 65 22 3a 22 45 30 30 30 30 32 37 22 2c 22  ode":"E000027","
0040   6d 6f 62 69 6c 65 22 3a 31 38 36 30 33 37 32 30  mobile":18603720
0050   39 32 35 2c 22 6e 61 6d 65 22 3a 22 e8 b0 a2 e5  925,"name":"....
0060   b0 8f e4 bd bc 22 2c 22 77 61 6c 6c 65 5f 69 64  .....","walle_id
0070   22 3a 37 37 32 38 37 2c 22 73 74 61 74 75 73 22  ":77287,"status"
0080   3a 36 2c 22 70 69 6e 79 69 6e 5f 6e 61 6d 65 22  :6,"pinyin_name"
0090   3a 22 78 78 6a 22 2c 22 73 65 78 22 3a 31 2c 22  :"xxj","sex":1,"
00a0   73 65 63 75 72 69 74 79 5f 6c 65 76 65 6c 22 3a  security_level":
00b0   36 30 2c 22 63 65 72 74 69 66 69 63 61 74 65 5f  60,"certificate_
00c0   74 79 70 65 22 3a 30 2c 22 63 65 72 74 69 66 69  type":0,"certifi
00d0   63 61 74 65 5f 6e 75 6d 62 65 72 22 3a 22 34 32  cate_number":"42
00e0   30 36 38 33 31 39 38 39 30 38 31 31 33 37 33 33  0683198908113733
00f0   22 2c 22 63 72 65 61 74 65 64 5f 61 74 22 3a 31  ","created_at":1
0100   34 33 31 35 35 30 30 32 39 30 30 30 2c 22 75 70  431550029000,"up
0110   64 61 74 65 64 5f 61 74 22 3a 31 34 34 39 32 32  dated_at":144922
0120   38 32 33 37 30 30 30 2c 22 6e 63 68 72 5f 69 64  8237000,"nchr_id
0130   22 3a 22 30 30 30 31 41 39 31 30 30 30 30 30 30  ":"0001A91000000
0140   30 30 30 32 45 51 50 22 7d 7d 0d 0a 24 35 31 39  0002EQP"}}..$519
0150   0d 0a 7b 22 75 73 65 72 49 64 22 3a 31 34 37 2c  ..{"userId":147,
0160   22 75 73 65 72 42 75 4c 69 73 74 22 3a 5b 33 31  "userBuList":[31
0170   37 35 5d 2c 22 74 61 67 73 4c 69 73 74 22 3a 5b  75],"tagsList":[
0180   5d 2c 22 75 73 65 72 42 75 52 6f 6c 65 44 74 6f  ],"userBuRoleDto
0190   22 3a 5b 7b 22 69 64 22 3a 38 31 35 32 34 2c 22  ":[{"id":81524,"
01a0   62 75 5f 69 64 22 3a 33 31 37 35 2c 22 62 75 5f  bu_id":3175,"bu_
01b0   6e 61 6d 65 22 3a 22 e4 ba a4 e6 98 93 e5 b9 b3  name":".........
01c0   e5 8f b0 42 55 22 2c 22 72 6f 6c 65 5f 69 64 22  ...BU","role_id"
01d0   3a 38 35 39 2c 22 72 6f 6c 65 5f 6e 61 6d 65 22  :859,"role_name"
01e0   3a 22 e5 89 af e6 80 bb e8 a3 81 22 2c 22 75 73  :".........","us
01f0   65 72 5f 69 64 22 3a 31 34 37 2c 22 75 73 65 72  er_id":147,"user
0200   5f 6e 61 6d 65 22 3a 22 e9 87 91 e9 91 ab 22 7d  _name":"......"}
0210   5d 2c 22 75 73 65 72 22 3a 7b 22 69 64 22 3a 31  ],"user":{"id":1
0220   34 37 2c 22 65 6d 61 69 6c 22 3a 22 78 69 6e 2e  47,"email":"xin.
0230   6a 69 6e 40 65 6c 65 2e 6d 65 22 2c 22 77 6f 72  jin@ele.me","wor
0240   6b 5f 63 6f 64 65 22 3a 22 45 30 30 30 30 32 39  k_code":"E000029
0250   22 2c 22 6d 6f 62 69 6c 65 22 3a 31 38 36 30 37  ","mobile":18607
0260   31 37 35 36 32 36 2c 22 6e 61 6d 65 22 3a 22 e9  175626,"name":".
0270   87 91 e9 91 ab 22 2c 22 77 61 6c 6c 65 5f 69 64  .....","walle_id
0280   22 3a 35 36 30 36 33 2c 22 73 74 61 74 75 73 22  ":56063,"status"
0290   3a 36 2c 22 70 69 6e 79 69 6e 5f 6e 61 6d 65 22  :6,"pinyin_name"
02a0   3a 22 6a 78 22 2c 22 73 65 78 22 3a 31 2c 22 73  :"jx","sex":1,"s
02b0   65 63 75 72 69 74 79 5f 6c 65 76 65 6c 22 3a 37  ecurity_level":7
02c0   30 2c 22 63 65 72 74 69 66 69 63 61 74 65 5f 74  0,"certificate_t
02d0   79 70 65 22 3a 30 2c 22 63 65 72 74 69 66 69 63  ype":0,"certific
02e0   61 74 65 5f 6e 75 6d 62 65 72 22 3a 22 34 32 30  ate_number":"420
02f0   31 30 36 31 39 38 35 31 31 32 34 32 35 31 30 22  106198511242510"
0300   2c 22 63 72 65 61 74 65 64 5f 61 74 22 3a 31 34  ,"created_at":14
0310   33 31 35 35 30 30 33 30 30 30 30 2c 22 75 70 64  31550030000,"upd
0320   61 74 65 64 5f 61 74 22 3a 31 34 34 39 32 32 38  ated_at":1449228
0330   31 31 35 30 30 30 2c 22 6e 63 68 72 5f 69 64 22  115000,"nchr_id"
0340   3a 22 30 30 30 31 41 39 31 30 30 30 30 30 30 30  :"0001A910000000
0350   30 30 32 45 52 45 22 7d 7d 0d 0a                 002ERE"}}..
```

**37648 号包**：提示 *[TCP Fast Retransmission]* 信息，结合数据包内容可以看出，重传数据即为之前判定丢失的数据；

```
*8
$532
{"userId":18370,"userBuList":[4594],"tagsList":[],"userBuRoleDto":[{"id":120993,"bu_id":4594,"bu_name":"..................","role_id":924,"role_name":"......","user_id":18370,"user_name":"......"}],"user":{"id":18370,"email":"hui.yaobj@ele.me","work_code":"E019529","mobile":13784728281,"name":"......","walle_id":23156752,"status":6,"pinyin_name":"yh","sex":1,"security_level":20,"certificate_type":0,"certificate_number":"130984199001023033","created_at":1438657893000,"updated_at":1449228019000,"nchr_id":"0001A910000000013PM5"}}
$536
{"userId":52708,"userBuList":[4704],"tagsList":[],"userBuRoleDto":[{"id":109320,"bu_id":4704,"bu_name":".....................","role_id":867,"role_name":"......","user_id":52708,"user_name":"Neil"}],"user":{"id":52708,"email":"liangang.qu@ele.me","work_code":"E055033","mobile":18637136367,"name":"Neil","walle_id":123225472,"status":6,"pinyin_name":"Neil","sex":1,"security_level":20,"certificate_type":0,"certificate_number":"370181198109040350","created_at":1476093008000,"updated_at":1477916894000,"nchr_id":"0001B61000000010WI0Q"}}
$-1
$532
{"userId":18370,"userBuList":[4594],"tagsList":[],"userBuRoleDto":[{"id":120993,"bu_id":4594,"bu_name":"..................","role_id":924,"role_name":"......","user_id":18370,"user_name":"......"}],"user":{"id":18370,"email":"hui.yaobj@ele.me","work_code":"E019529","mobile":13784728281,"name":"......","walle_id":23156752,"status":6,"pinyin_name":"yh","sex":1,"security_level":20,"certificate_type":0,"certificate_number":"130984199001023033","created_at":1438657893000,"updated_at":1449228019000,"nchr_id":"0001A910000000013PM5"}}
$533
{"userId":1117,"userBuList":[4596],"tagsList":[],"userBuRoleDto":[{"id":89811,"bu_id":4596,"bu_name":"..................","role_id":923,"role_name":"......","user_id":1117,"user_name":"........."}],"user":{"id":1117,"email":"xuhao.hu@ele.me","work_code":"E001000","mobile":18562612164,"name":".........","walle_id":1682514,"status":6,"pinyin_name":"hxh","sex":1,"security_level":20,"certificate_type":0,"certificate_number":"640102199404171817","created_at":1431550294000,"updated_at":1463403636000,"nchr_id":"0001A910000000002GCC"}}
$512
{"userId":13,"userBuList":[104],"tagsList":[],"userBuRoleDto":[{"id":769,"bu_id":104,"bu_name":".........","role_id":865,"role_name":"CEO","user_id":13,"user_name":"........."}],"user":{"id":13,"email":"mark.zhang@ele.me","work_code":"E000001","mobile":13482200180,"name":".........","walle_id":869105,"status":6,"pinyin_name":"zxh","sex":1,"security_level":90,"certificate_type":0,"certificate_number":"310103198504094033","created_at":1431550018000,"updated_at":1449228115000,"nchr_id":"0001A910000000002EGU"}}
$526
{"userId":153,"userBuList":[3174],"tagsList":[],"userBuRoleDto":[{"id":53306,"bu_id":3174,"bu_name":"............","role_id":922,"role_name":"......","user_id":153,"user_name":"........."}],"user"
```


> 遗留问题：
>
> - 抓到的数据包显示：每个报文都会被重传一次，原因何在？
> - 触发快速重传的机制
> - 为何会丢包

### 解决办法

xxx


## #07 packetbeat 在进行 request-response 关联（构建 transaction）时，在某些情况下是不正确的

### 问题描述

情况一：

```
responsetime(947201 microseconds)	==>    No.<1>
----
{"@timestamp":"2016-12-29T07:21:24.657Z","beat":{"hostname":"xg-mesos-39","name":"xg-mesos-39","version":"6.0.0-alpha1"},"bytes_in":14,"bytes_out":40,"client_ip":"10.0.242.43","client_port":7125,"client_proc":"","client_server":"","ip":"10.0.246.114","method":"PING","port":48877,"proc":"","query":"PING","redis":{"return_value":"[REPLCONF, ACK, 5372098]"},"resource":"","responsetime":947201,"server":"","status":"OK","type":"redis"}
```

情况二：

（目前未遇到，但推断可能）发送请求后，没有收到响应的应答，一段时间后收到其他回复信息，结果被当成了应答，此时的关联是错误的；


### 问题原因

上例中将 `PING` 和 `[REPLCONF, ACK, 5372098]` 进行了关联，而这种关联是不对的；

### 解决办法

这种问题需要具体情况具体分析了，可能需要进行命令过滤处理；

## #08 编译出的 packetbeat 可执行程序需要动态链接 libpcap.so 库（当前默认情况），而目标生产服务器上存在多种版本的操作系统，另外相应的 .so 版本也可能存在不一致问题

### 解决办法

目前采用直接将编译机环境中的 libpcap.so 和可执行程序捆绑打包提供的方式进行解决；在执行时通过 `LD_LIBRARY_PATH` 引用当前目录下的 .so 文件；





