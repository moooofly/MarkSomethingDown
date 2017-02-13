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

补充结论：经过深入研究发现，应答内容大并不是充分条件，导致数据解析出错和分包位置有关（有些情况分包数据并不会导致错误）；

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
> - 触发快速重传的机制？
> - 为何会丢包？

### 解决办法

xxx

### 其他

在官网论坛上的[讨论](https://discuss.elastic.co/t/packetbeat-err-failed-to-read-integer-reply-expected-digit/74352)；


## #07 packetbeat 在进行 request-response 关联（构建 transaction）时，在某些情况下是不正确的

### 问题描述

情况一：

```
responsetime(947201 microseconds)	==>    No.<1>
----
{"@timestamp":"2016-12-29T07:21:24.657Z","beat":{"hostname":"xg-mesos-39","name":"xg-mesos-39","version":"6.0.0-alpha1"},"bytes_in":14,"bytes_out":40,"client_ip":"10.0.242.43","client_port":7125,"client_proc":"","client_server":"","ip":"10.0.246.114","method":"PING","port":48877,"proc":"","query":"PING","redis":{"return_value":"[REPLCONF, ACK, 5372098]"},"resource":"","responsetime":947201,"server":"","status":"OK","type":"redis"}
```

情况二：

发送请求后，没有收到对应的应答响应，一段时间后收到其他响应信息，结果被当成了对应的应答，此时的关联是错误的；

> 应该不存在这种情况：因为 transaction 的构建是基于 TCP 连接的；

### 问题原因

上述将 `PING` 和 `[REPLCONF, ACK, 5372098]` 进行了关联，而这两者的关联明显是不对的；

- **`PING`** 的使用

1. [**客户端-服务器**] 使用客户端向 Redis 服务器发送一个 PING ，如果服务器运作正常的话，会返回一个 PONG 。通常用于测试与服务器的连接是否仍然生效，或者用于测量延迟值；

2. [**Sentinel**] 在默认情况下，Sentinel 会以每秒一次的频率向所有与它创建了命令连接的实例（包括主服务器、从服务器、其他 Sentinel 在内）发送 PING 命令，并通过实例返回的 PING 命令回复（有效回复为 +PONG/-LOADING/-MASTERDOWN）来判断实例是否在线（主观下线状态检测）；

3. [**主从复制**] 当从服务器成为主服务器的客户端后，做的第一件事就是向主服务器发送一个 PING 命令；两个作用：a) 检查套接字的读写状态是否正常；b) 检查主服务器能否正常处理命令请求；只有从服务器在规定时间内读取到主服务器返回的 PONG 才算成功；

4. [**主从复制**] Slaves 以预定义的周期向 server 发送 PING；该周期通过 `repl_ping_slave_period` 选项进行配置，默认为 10 秒； 
The original replication protocol was vulnerable to network/Internet outages where the master detects the outage and closes the connection, but the slave does not. The slave thinks the connection is still open and the master simply has no updates to send (low traffic or no traffic). So the slave never disconnects and re-connects to restart the replication. I know this very well. I have some v2.0.x Redis instances that replicate across 3,000 miles and once or twice a month this problem occurs.
Adding PING to the replication protocol solved that. The slave now detects the connection problem when the PING replies stop coming from the master. The slave can close its end of the connection and re-connect again.

5. [**集群**] 集群里的每个节点默认每隔一秒钟就会从已知节点列表中随机选出五个节点，然后对这五个节点中最长时间没有发送过 PING 消息的节点发送 PING 消息，以此来检测被选中的节点是否在线；除此之外，如果节点 A 最后一次收到节点 B 发送的 PONG 消息的时间，距离当前时间已经超过了节点 A 的 cluster-node-timeout 选项设置时长的一半，那么节点 A 也会向节点 B 发送 PING 消息，这可以防止节点 A 因为长时间没有随机选中节点 B 作为 PING 消息的发送对象，而导致对节点 B 的信息更新滞后；

- **`[REPLCONF, ACK, <replication_offset>]`** 的使用
在命令传播阶段，从服务器默认会以每秒一次的频率，向主服务器发送该命令；该命令的作用为：a) 检测主从服务器的网络连接状态；b) 辅助实现 min-slaves 选项；c) 检测命令丢失；

具体抓包数据如下

```
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*3
$8
REPLCONF
$3
ACK
$11
81234009046
*1
$4
PING
*3
$8
REPLCONF
$3
ACK
$11
81234009060
```

可以看到，在一个 10s 的抓包周期中，出现了 10 个 `[REPLCONF, ACK, xxxx]` 和 1 个 `PING` ；可以确定，该 PING 为基于 `repl_ping_slave_period` 选项的包括 PING ；

### 解决办法

这种问题需要具体情况具体分析了，可能需要进行命令过滤处理；换一种说法：一般情况下我们不太会通过 PING 来确定网络延迟，因为常规的 Redis 命令交互会起到同样到作用；因此，理论上讲非客户端直接发起的 Redis 命令都可以过滤掉（当前想法）；

## #08 编译出的 packetbeat 可执行程序需要动态链接 libpcap.so 库（当前默认情况），而目标生产服务器上存在多种版本的操作系统，另外相应的 .so 版本也可能存在不一致问题

### 解决办法

目前采用直接将编译机环境中的 libpcap.so 和可执行程序捆绑打包提供的方式进行解决；在执行时通过 `LD_LIBRARY_PATH` 引用当前目录下的 .so 文件；





