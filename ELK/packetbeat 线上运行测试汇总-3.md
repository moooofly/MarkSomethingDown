# packetbeat 线上运行测试汇总-3

标签（空格分隔）： packetbeat

---

> 本文是《[packetbeat 线上运行测试汇总](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/packetbeat%20%E7%BA%BF%E4%B8%8A%E8%BF%90%E8%A1%8C%E6%B5%8B%E8%AF%95%E6%B1%87%E6%80%BB.md)》和《[packetbeat 线上运行测试汇总-2](https://github.com/moooofly/MarkSomethingDown/blob/master/Golang/packetbeat%20%E7%BA%BF%E4%B8%8A%E8%BF%90%E8%A1%8C%E6%B5%8B%E8%AF%95%E6%B1%87%E6%80%BB-2.md)》之后的最新版；

当前版本功能：

- 抓包分析整体流程**变更**为：基于 eoc 触发 packetbeat 针对指定 redis cluster 进行定时抓包分析 => 上传分析结果文件到 ceph => 基于 web 查看分析结果文件
- 能够**区分** redis 的 master 和 slave ，在抓包时直接排除掉 slave 相关机器；
- 已**排除** redis 的 master-slave 通信引起的分析干扰；
- 默认情况下，基于 taskset 直接限制 packetbeat 跑在目标机器**最后一个** cpu 核心上；


当前测试针对：

- 只抓取 10s 的数据包；
- 基于 taskset 将 packetbeat 限制在 CPU 22 和 CPU 23 上；


----------

> 以下测试内容针对线上一个比较繁忙的 redis 集群进行；

集群信息如下

```shell
[root@xg-minos-rediscluster-1 ~]# /opt/redis/redis_bin/redis-3.0.3.tar.gz/bin/redis-cli -p 7602
127.0.0.1:7602> cluster info
cluster_state:ok
cluster_slots_assigned:16384
cluster_slots_ok:16384
cluster_slots_pfail:0
cluster_slots_fail:0
cluster_known_nodes:182
cluster_size:91
cluster_current_epoch:302
cluster_my_epoch:264
cluster_stats_messages_sent:547670487
cluster_stats_messages_received:547665351
127.0.0.1:7602>
127.0.0.1:7602>
```

## 抓包分析系统执行过程


- 查看 redis cluster 信息，并确定 guldan_redis 集群中 redis 分布情况

![web 版 guldan_redis 集群分布](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/web%20%E7%89%88%20guldan_redis%20%E9%9B%86%E7%BE%A4%E5%88%86%E5%B8%83.png "web 版 guldan_redis 集群分布")

可以看到 redis 集群 `guldan_redis` 分布在 9 台物理上，通过 8 个 corvus 代理对外提供服务；


- 通过 eoc 脚本触发抓包分析（持续时间 10s）

命令行参数配置情况

![web 版命令行参数配置](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/web%20%E7%89%88%E5%91%BD%E4%BB%A4%E8%A1%8C%E5%8F%82%E6%95%B0%E9%85%8D%E7%BD%AE.png "web 版命令行参数配置")

eoc 脚本执行情况（[job 详情](http://eoc.elenet.me/tasks/bcf7b1e4-1520-11e7-b17c-525400b2101e)）

![web 版 eoc 执行](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/web%20%E7%89%88%20eoc%20%E6%89%A7%E8%A1%8C.png "web 版 eoc 执行")

- eoc 脚本在 `xg-minos-rediscluster-1` 主机上的执行日志

```
2017-03-30 16:13:22,141 INFO me.ele.opdev.eoc.agent.Executor - agent version: 0.3.0-6
2017-03-30 16:13:22,142 DEBUG me.ele.opdev.eoc.agent.Executor - start download job package
2017-03-30 16:13:22,142 DEBUG me.ele.opdev.eoc.agent.Executor - start unpack package: kaixing.wu/eoc-PackageCaptureAnalysis
2017-03-30 16:13:22,258 DEBUG me.ele.opdev.eoc.agent.Executor - complete unpack package: kaixing.wu/eoc-PackageCaptureAnalysis
2017-03-30 16:13:22,259 DEBUG me.ele.opdev.eoc.agent.Executor - start render /run.sh.tmpl
2017-03-30 16:13:22,261 DEBUG me.ele.opdev.eoc.agent.Executor - complete render /run.sh.tmpl
2017-03-30 16:13:22,261 INFO me.ele.opdev.eoc.agent.Executor - ====== START Job<kaixing.wu/eoc-PackageCaptureAnalysis> ======
Run command:taskset -c 22-23 timeout 10 ./packageCaptureAnalysis -c ./packetbeat.yml -t -E packetbeat.protocols.redis.ports=7602,7603,7604,7605,7606,7607,7608,7609,7610,7611,7612,7613,7614 -e
2017/03/30 08:13:22.509571 beat.go:267: INFO Home path: [/data/dump] Config path: [/data/dump] Data path: [/data/dump/data] Logs path: [/data/dump/logs]
2017/03/30 08:13:22.509605 beat.go:177: INFO Setup Beat: packetbeat; Version: 6.0.0-alpha1
2017/03/30 08:13:22.509665 file.go:45: INFO File output path set to: ./logs
2017/03/30 08:13:22.509676 file.go:46: INFO File output base filename set to: packetbeat
2017/03/30 08:13:22.509681 file.go:49: INFO Rotate every bytes set to: 102400000
2017/03/30 08:13:22.509686 file.go:53: INFO Number of files set to: 20
2017/03/30 08:13:22.509716 logp.go:219: INFO Metrics logging every 30s
2017/03/30 08:13:22.509722 outputs.go:106: INFO Activated file as output plugin.
2017/03/30 08:13:22.509850 publish.go:291: INFO Publisher name: xg-minos-rediscluster-1
2017/03/30 08:13:22.511204 async.go:63: INFO Flush Interval set to: -1s
2017/03/30 08:13:22.511216 async.go:64: INFO Max Bulk Size set to: -1
2017/03/30 08:13:22.511269 procs.go:79: INFO Process matching disabled
2017/03/30 08:13:22.512528 protos.go:89: INFO registered protocol plugin: mysql
2017/03/30 08:13:22.512542 protos.go:89: INFO registered protocol plugin: redis
2017/03/30 08:13:22.512546 protos.go:89: INFO registered protocol plugin: amqp
2017/03/30 08:13:22.512551 protos.go:89: INFO registered protocol plugin: http
2017/03/30 08:13:22.524495 beat.go:207: INFO packetbeat start running.
2017/03/30 08:13:32.328886 packetbeat.go:184: INFO Packetbeat send stop signal
2017/03/30 08:13:32.329314 sniffer.go:384: INFO Input finish. Processed 761054 packets. Have a nice day!
2017/03/30 08:13:32.329536 logp.go:245: INFO Total non-zero values: libbeat.publisher.messages_in_worker_queues=273 libbeat.publisher.published_events=108883 redis.unmatched_responses=231733 tcp.dropped_because_of_gaps=236556
2017/03/30 08:13:32.329556 logp.go:246: INFO Uptime: 10.007247117s
2017/03/30 08:13:32.329564 beat.go:211: INFO packetbeat stopped.
./packageCaptureAnalysis run finish.
Run command: /bin/python redis_analysis.py -p logs -f packetbeat
run redis_analysis.py finish.
start upload dump file to the ceph.
% Total % Received % Xferd Average Speed Time Time Time Current
Dload Upload Total Spent Left Speed

 0 0 0 0 0 0 0 0 --:--:-- --:--:-- --:--:-- 0
100 6063 0 0 100 6063 0 196k --:--:-- --:--:-- --:--:-- 197k
dump file upload ceph successfully.
2017-03-30 16:13:39,277 INFO me.ele.opdev.eoc.agent.Executor - ====== END Job<kaixing.wu/eoc-PackageCaptureAnalysis>: SUCCESS ======
```

可以看出，抓包分析的持续时间为 `16:13:22 ～ 16:13:32` ，在 10 秒内处理了 761054 个数据包，因此该机器上针对 redis 协议的 `pps` 为 **76105.4** ；其他机器类同；

- 查看自动上传到 ceph 上的分析文件

![web 版自动上传 ceph 的分析文件](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/web%20%E7%89%88%E8%87%AA%E5%8A%A8%E4%B8%8A%E4%BC%A0%20ceph%20%E7%9A%84%E5%88%86%E6%9E%90%E6%96%87%E4%BB%B6.png "web 版自动上传 ceph 的分析文件")

## 监控曲线

- [redis cluster 的 cpu 曲线（整体）](https://t.elenet.me/dashboard/dashboard/db/system-monitor-cpu?var-machine=xg-minos-rediscluster-1&var-machine=xg-minos-rediscluster-2&var-machine=xg-minos-rediscluster-3&var-machine=xg-public-rediscluster-10&var-machine=xg-public-rediscluster-11&var-machine=xg-public-rediscluster-12&var-machine=xg-public-rediscluster-43&var-machine=xg-public-rediscluster-8&var-machine=xg-public-rediscluster-9&from=1490861519254&to=1490861760608)

![guldan_redis_cpu_10s_1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/guldan_redis_cpu_10s_1.png "guldan_redis_cpu_10s_1")

![guldan_redis_cpu_10s_2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/guldan_redis_cpu_10s_2.png "guldan_redis_cpu_10s_2")


- [redis cluster 上 pb 占用的 CPU 22 和 23 曲线](https://t.elenet.me/dashboard/dashboard/db/system-monitor-per-cpu?from=1490861519254&to=1490861760608&var-machine=xg-minos-rediscluster-1&var-machine=xg-minos-rediscluster-2&var-machine=xg-minos-rediscluster-3&var-machine=xg-public-rediscluster-10&var-machine=xg-public-rediscluster-11&var-machine=xg-public-rediscluster-12&var-machine=xg-public-rediscluster-43&var-machine=xg-public-rediscluster-8&var-machine=xg-public-rediscluster-9&var-cpu=22&var-cpu=23)

![system-monitor-per-cpu_guldan_redis_redis_cpu2223](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/system-monitor-per-cpu_guldan_redis_redis_cpu2223.png "system-monitor-per-cpu_guldan_redis_redis_cpu2223")

> 由上图可知，被 pb 绑定的 CPU 22 和 23 最高跑到 **85%** 左右（用户态）；

![system-monitor-per-cpu_guldan_redis_redis_cpu2223_up](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/system-monitor-per-cpu_guldan_redis_redis_cpu2223_up.png "system-monitor-per-cpu_guldan_redis_redis_cpu2223_up")

> 由上图可知，在设置 pb 只运行 10s 的情况下，曲线的误差会比较大：
>
> - 基于 eoc 在不同机器上运行 pb 的起始时间会有所不同（根据日志记录的时间，时间差距可达 3s 左右）；
> - 由于 esm 的采样精度为 10s ，因此得到的曲线可能会有问题；


- [corvus proxy 的 cpu 曲线（整体）](https://t.elenet.me/dashboard/dashboard/db/system-monitor-cpu?var-machine=xg-redis-corvusproxy-41&var-machine=xg-redis-corvusproxy-43&var-machine=xg-redis-corvusproxy-44&var-machine=xg-redis-corvusproxy-45&var-machine=xg-redis-corvusproxy-46&var-machine=xg-redis-corvusproxy-48&var-machine=xg-redis-corvusproxy-49&var-machine=xg-redis-corvusproxy-50&from=1490861519254&to=1490861760608)

![guldan_redis_corvus_cpu_10s_1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/guldan_redis_corvus_cpu_10s_1.png "guldan_redis_corvus_cpu_10s_1")

![guldan_redis_corvus_cpu_10s_2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/guldan_redis_corvus_cpu_10s_2.png "guldan_redis_corvus_cpu_10s_2")

- [corvus proxy 上 pb 占用的 CPU 22 和 23 曲线](https://t.elenet.me/dashboard/dashboard/db/system-monitor-per-cpu?from=1490861519254&to=1490861760608&var-machine=xg-redis-corvusproxy-41&var-machine=xg-redis-corvusproxy-43&var-machine=xg-redis-corvusproxy-44&var-machine=xg-redis-corvusproxy-45&var-machine=xg-redis-corvusproxy-46&var-machine=xg-redis-corvusproxy-48&var-machine=xg-redis-corvusproxy-49&var-machine=xg-redis-corvusproxy-50&var-cpu=22&var-cpu=23)

![system-monitor-per-cpu_guldan_redis_corvus_cpu2223](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/system-monitor-per-cpu_guldan_redis_corvus_cpu2223.png "system-monitor-per-cpu_guldan_redis_corvus_cpu2223")

> 由上图可知，被 pb 绑定的 CPU 22 和 23 最高跑到 **60%** 左右（用户态）；

![system-monitor-per-cpu_guldan_redis_corvus_cpu2223_up](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/system-monitor-per-cpu_guldan_redis_corvus_cpu2223_up.png "system-monitor-per-cpu_guldan_redis_corvus_cpu2223_up")

> 由上图可知，在设置 pb 只运行 10s 的情况下，曲线的误差会比较大：

## packetbeat 日志分析

> 针对每台机器日志，只讨论耗时最高 top 3；
> 若 responsetime 出现数值相等情况，应该理解成精度不够 or pipeline 行为导致；

- `2core10s_xg-minos-rediscluster-1_@7602@7603@7604@7605@7606@7607@7608@7609@7610@7611@7612@7613@7614@_20170330161322_20170330161332.log`

```
total transactions : 108609

responsetime(1006002 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:22.562Z","beat":{"hostname":"xg-minos-rediscluster-1","name":"xg-minos-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":56,"bytes_out":4,"client_ip":"10.0.58.190","client_port":30354,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.58","method":"SMEMBERS","port":7613,"proc":"","query":"SMEMBERS food_restaurant:844768:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:844768:contents","responsetime":1006002,"server":"","status":"OK","type":"redis"}


responsetime(686024 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:22.993Z","beat":{"hostname":"xg-minos-rediscluster-1","name":"xg-minos-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":4,"client_ip":"10.0.58.161","client_port":40253,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.58","method":"SMEMBERS","port":7602,"proc":"","query":"SMEMBERS restaurant:803716:contents","redis":{"return_value":"[]"},"resource":"restaurant:803716:contents","responsetime":686024,"server":"","status":"OK","type":"redis"}


responsetime(685128 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:22.995Z","beat":{"hostname":"xg-minos-rediscluster-1","name":"xg-minos-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.58.161","client_port":40253,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.58","method":"SMEMBERS","port":7602,"proc":"","query":"SMEMBERS food_restaurant:2214850:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:2214850:contents","responsetime":685128,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.006002s > 686.024ms > 685.128ms
- 耗时命令：均为 SMEMBERS
- 端口关系：30354 <-> 7613(1), 40253 <-> 7602(2)


----------


- `2core10s_xg-minos-rediscluster-2_@7602@7603@7604@7605@7606@7607@7608@7609@7610@7611@7612@7613@7614@_20170330161322_20170330161332.log`

```
total transactions : 105123

responsetime(1255779 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:22.798Z","beat":{"hostname":"xg-minos-rediscluster-2","name":"xg-minos-rediscluster-2","version":"6.0.0-alpha1"},"bytes_in":56,"bytes_out":4,"client_ip":"10.0.58.181","client_port":28999,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.59","method":"SMEMBERS","port":7610,"proc":"","query":"SMEMBERS food_restaurant:793253:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:793253:contents","responsetime":1255779,"server":"","status":"OK","type":"redis"}


responsetime(1102751 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:23.127Z","beat":{"hostname":"xg-minos-rediscluster-2","name":"xg-minos-rediscluster-2","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":48,"client_ip":"10.0.58.178","client_port":27588,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.59","method":"SMEMBERS","port":7613,"proc":"","query":"SMEMBERS food_restaurant:1066587:contents","redis":{"return_value":"[2488, 3258, 203645, 203727]"},"resource":"food_restaurant:1066587:contents","responsetime":1102751,"server":"","status":"OK","type":"redis"}


responsetime(684747 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:22.995Z","beat":{"hostname":"xg-minos-rediscluster-2","name":"xg-minos-rediscluster-2","version":"6.0.0-alpha1"},"bytes_in":56,"bytes_out":4,"client_ip":"10.0.58.161","client_port":56251,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.59","method":"SMEMBERS","port":7604,"proc":"","query":"SMEMBERS food_restaurant:991350:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:991350:contents","responsetime":684747,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.255779s > 1.102751s > 684.747ms
- 耗时命令：均为 SMEMBERS
- 端口关系：28999 <-> 7610(1), 27588 <-> 7613(1), 56251 <-> 7604(1)


----------


- `2core10s_xg-minos-rediscluster-3_@7602@7603@7604@7605@7606@7607@7608@7609@7610@7611@7612@7613@7614@_20170330161322_20170330161332.log`

```
total transactions : 118859

responsetime(1261181 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:22.762Z","beat":{"hostname":"xg-minos-rediscluster-3","name":"xg-minos-rediscluster-3","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":4,"client_ip":"10.0.58.161","client_port":62575,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.60","method":"SMEMBERS","port":7612,"proc":"","query":"SMEMBERS restaurant:811412:contents","redis":{"return_value":"[]"},"resource":"restaurant:811412:contents","responsetime":1261181,"server":"","status":"OK","type":"redis"}


responsetime(972203 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:23.125Z","beat":{"hostname":"xg-minos-rediscluster-3","name":"xg-minos-rediscluster-3","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":28,"client_ip":"10.0.58.178","client_port":53910,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.60","method":"SMEMBERS","port":7614,"proc":"","query":"SMEMBERS restaurant:515319:contents","redis":{"return_value":"[203606, 203636]"},"resource":"restaurant:515319:contents","responsetime":972203,"server":"","status":"OK","type":"redis"}


responsetime(969277 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:23.127Z","beat":{"hostname":"xg-minos-rediscluster-3","name":"xg-minos-rediscluster-3","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.58.178","client_port":49978,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.60","method":"SMEMBERS","port":7606,"proc":"","query":"SMEMBERS food_restaurant:2087173:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:2087173:contents","responsetime":969277,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.261181s > 972.203ms > 969.277ms
- 耗时命令：均为 SMEMBERS
- 端口关系：62575 <-> 7612(1), 53910 <-> 7614(1), 49978 <-> 7606(1)


----------


- `2core10s_xg-public-rediscluster-10_@7301@7302@7303@7304@7305@7306@7307@7308@_20170330161322_20170330161332.log`

```
total transactions : 190065

responsetime(194672 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:22.584Z","beat":{"hostname":"xg-public-rediscluster-10","name":"xg-public-rediscluster-10","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":7,"client_ip":"10.0.58.161","client_port":36465,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.30","method":"SMEMBERS","port":7304,"proc":"","query":"SMEMBERS food_restaurant:1256442:contents","redis":{"return_value":"0"},"resource":"food_restaurant:1256442:contents","responsetime":194672,"server":"","status":"OK","type":"redis"}


responsetime(163000 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:22.587Z","beat":{"hostname":"xg-public-rediscluster-10","name":"xg-public-rediscluster-10","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.58.183","client_port":42227,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.30","method":"SMEMBERS","port":7303,"proc":"","query":"SMEMBERS restaurant:1382695:contents","redis":{"return_value":"[]"},"resource":"restaurant:1382695:contents","responsetime":163000,"server":"","status":"OK","type":"redis"}


responsetime(79868 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:22.576Z","beat":{"hostname":"xg-public-rediscluster-10","name":"xg-public-rediscluster-10","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":182,"client_ip":"10.0.58.178","client_port":44749,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.30","method":"SMEMBERS","port":7306,"proc":"","query":"SMEMBERS restaurant:1344536:contents","redis":{"return_value":"[844, 2330, 2492, 2532, 2555, 2823, 3258, 203567, 203578, 203636, 203727, 203881, 204138, 204578, 204749, 205181]"},"resource":"restaurant:1344536:contents","responsetime":79868,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：194.672ms > 163.000ms > 79.868ms
- 耗时命令：均为 SMEMBERS
- 端口关系：36465 <-> 7304(1), 42227 <-> 7303(1), 44749 <-> 7306(1) 


----------


- `2core10s_xg-public-rediscluster-11_@7301@7302@7303@7304@7305@7306@7307@7308@7313@7314@7315@7316@_20170330161321_20170330161331.log`

```
total transactions : 179496

responsetime(968233 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:23.128Z","beat":{"hostname":"xg-public-rediscluster-11","name":"xg-public-rediscluster-11","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":4,"client_ip":"10.0.58.178","client_port":43522,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.106","method":"SMEMBERS","port":7305,"proc":"","query":"SMEMBERS restaurant:862182:contents","redis":{"return_value":"[]"},"resource":"restaurant:862182:contents","responsetime":968233,"server":"","status":"OK","type":"redis"}


responsetime(564689 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:22.079Z","beat":{"hostname":"xg-public-rediscluster-11","name":"xg-public-rediscluster-11","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.58.169","client_port":30730,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.106","method":"SMEMBERS","port":7306,"proc":"","query":"SMEMBERS food_restaurant:1466980:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:1466980:contents","responsetime":564689,"server":"","status":"OK","type":"redis"}


responsetime(329355 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:21.912Z","beat":{"hostname":"xg-public-rediscluster-11","name":"xg-public-rediscluster-11","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":40,"client_ip":"10.0.58.185","client_port":22382,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.106","method":"SMEMBERS","port":7315,"proc":"","query":"SMEMBERS restaurant:2095966:contents","redis":{"return_value":"[203567, 203852, 203881]"},"resource":"restaurant:2095966:contents","responsetime":329355,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：968.233ms > 564.689ms > 329.355ms
- 耗时命令：均为 SMEMBERS
- 端口关系：43522 <-> 7305(1), 30730 <-> 7306(1), 22382 <-> 7315(1)


----------


- `2core10s_xg-public-rediscluster-12_@7301@7302@7303@7304@7305@7306@7307@7308@_20170330161321_20170330161331.log`

```
total transactions : 188930

responsetime(546281 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:21.973Z","beat":{"hostname":"xg-public-rediscluster-12","name":"xg-public-rediscluster-12","version":"6.0.0-alpha1"},"bytes_in":81,"bytes_out":4,"client_ip":"10.0.58.183","client_port":61175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.91","method":"GET","port":7308,"proc":"","query":"GET cms_device_admin_by_uuid:A075F0DA-2E82-4C8B-86A1-F565C4041225","redis":{"return_value":"[]"},"resource":"cms_device_admin_by_uuid:A075F0DA-2E82-4C8B-86A1-F565C4041225","responsetime":546281,"server":"","status":"OK","type":"redis"}


responsetime(438768 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:22.520Z","beat":{"hostname":"xg-public-rediscluster-12","name":"xg-public-rediscluster-12","version":"6.0.0-alpha1"},"bytes_in":54,"bytes_out":4,"client_ip":"10.0.58.183","client_port":23686,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.91","method":"SMEMBERS","port":7301,"proc":"","query":"SMEMBERS restaurant:150840400:contents","redis":{"return_value":"[]"},"resource":"restaurant:150840400:contents","responsetime":438768,"server":"","status":"OK","type":"redis"}


responsetime(187638 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:22.021Z","beat":{"hostname":"xg-public-rediscluster-12","name":"xg-public-rediscluster-12","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.58.178","client_port":36870,"client_proc":"","client_server":"","direction":"in","ip":"10.0.43.91","method":"SMEMBERS","port":7308,"proc":"","query":"SMEMBERS restaurant:1088085:contents","redis":{"return_value":"[]"},"resource":"restaurant:1088085:contents","responsetime":187638,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：546.281ms > 438.768ms > 187.638ms
- 耗时命令：GET, SMEMBERS
- 端口关系：61175 <-> 7308(1), 23686 <-> 7301(1), 36870 <-> 7308(1)


----------

- `2core10s_xg-public-rediscluster-43_@7139@7140@7141@7142@_20170330161321_20170330161331.log`

```
total transactions : 186820

responsetime(196142 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:22.124Z","beat":{"hostname":"xg-public-rediscluster-43","name":"xg-public-rediscluster-43","version":"6.0.0-alpha1"},"bytes_in":59,"bytes_out":4,"client_ip":"10.0.58.178","client_port":41916,"client_proc":"","client_server":"","direction":"in","ip":"10.0.28.24","method":"SMEMBERS","port":7140,"proc":"","query":"SMEMBERS food_restaurant:143198951:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:143198951:contents","responsetime":196142,"server":"","status":"OK","type":"redis"}


responsetime(6845 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:22.265Z","beat":{"hostname":"xg-public-rediscluster-43","name":"xg-public-rediscluster-43","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":104,"client_ip":"10.0.58.181","client_port":39317,"client_proc":"","client_server":"","direction":"in","ip":"10.0.28.24","method":"SMEMBERS","port":7139,"proc":"","query":"SMEMBERS food_restaurant:1817059:contents","redis":{"return_value":"[1242, 2330, 2532, 2555, 203819, 204126, 204410, 204647, 204731]"},"resource":"food_restaurant:1817059:contents","responsetime":6845,"server":"","status":"OK","type":"redis"}


responsetime(6840 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:22.265Z","beat":{"hostname":"xg-public-rediscluster-43","name":"xg-public-rediscluster-43","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.58.181","client_port":39317,"client_proc":"","client_server":"","direction":"in","ip":"10.0.28.24","method":"SMEMBERS","port":7139,"proc":"","query":"SMEMBERS restaurant:1973552:contents","redis":{"return_value":"[]"},"resource":"restaurant:1973552:contents","responsetime":6840,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：196.142ms > 6.845ms > 6.840ms
- 耗时命令：均为 SMEMBERS
- 端口关系：41916 <-> 7140(1), 39317 <-> 7139(2)


----------

- `2core10s_xg-public-rediscluster-8_@7301@7302@7303@7304@7305@7306@7307@7308@_20170330161321_20170330161331.log`

```
total transactions : 2985

responsetime(182 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:31.750Z","beat":{"hostname":"xg-public-rediscluster-8","name":"xg-public-rediscluster-8","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":28,"client_ip":"10.0.58.183","client_port":40187,"client_proc":"","client_server":"","direction":"in","ip":"10.0.44.33","method":"SMEMBERS","port":7308,"proc":"","query":"SMEMBERS restaurant:1159600:contents","redis":{"return_value":"[204315, 204559]"},"resource":"restaurant:1159600:contents","responsetime":182,"server":"","status":"OK","type":"redis"}


responsetime(137 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:31.813Z","beat":{"hostname":"xg-public-rediscluster-8","name":"xg-public-rediscluster-8","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":52,"client_ip":"10.0.58.181","client_port":61566,"client_proc":"","client_server":"","direction":"in","ip":"10.0.44.33","method":"SMEMBERS","port":7302,"proc":"","query":"SMEMBERS restaurant:351379:contents","redis":{"return_value":"[826, 985, 2330, 2532, 2555]"},"resource":"restaurant:351379:contents","responsetime":137,"server":"","status":"OK","type":"redis"}


responsetime(137 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:31.813Z","beat":{"hostname":"xg-public-rediscluster-8","name":"xg-public-rediscluster-8","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.58.181","client_port":61566,"client_proc":"","client_server":"","direction":"in","ip":"10.0.44.33","method":"SMEMBERS","port":7302,"proc":"","query":"SMEMBERS food_restaurant:1874109:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:1874109:contents","responsetime":137,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：0.182ms > 0.137ms = 0.137ms
- 耗时命令：均为 SMEMBERS
- 端口关系：40187 <-> 7308(1), 61566 <-> 7302(2)


----------


- `2core10s_xg-public-rediscluster-9_@7301@7302@7303@7304@7305@7306@7307@7308@7309@7310@7311@7312@_20170330161321_20170330161332.log`

```
total transactions : 184957

responsetime(650636 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:22.699Z","beat":{"hostname":"xg-public-rediscluster-9","name":"xg-public-rediscluster-9","version":"6.0.0-alpha1"},"bytes_in":59,"bytes_out":4,"client_ip":"10.0.58.181","client_port":20035,"client_proc":"","client_server":"","direction":"in","ip":"10.0.44.63","method":"SMEMBERS","port":7307,"proc":"","query":"SMEMBERS food_restaurant:144976377:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:144976377:contents","responsetime":650636,"server":"","status":"OK","type":"redis"}


responsetime(311608 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:22.383Z","beat":{"hostname":"xg-public-rediscluster-9","name":"xg-public-rediscluster-9","version":"6.0.0-alpha1"},"bytes_in":81,"bytes_out":14,"client_ip":"10.0.58.181","client_port":45872,"client_proc":"","client_server":"","direction":"in","ip":"10.0.44.63","method":"GET","port":7301,"proc":"","query":"GET cms_device_admin_by_uuid:82a57896-adf2-33b0-9221-0031b4f256fb","redis":{"return_value":"[1242]"},"resource":"cms_device_admin_by_uuid:82a57896-adf2-33b0-9221-0031b4f256fb","responsetime":311608,"server":"","status":"OK","type":"redis"}


responsetime(184217 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:22.350Z","beat":{"hostname":"xg-public-rediscluster-9","name":"xg-public-rediscluster-9","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":28,"client_ip":"10.0.58.169","client_port":36673,"client_proc":"","client_server":"","direction":"in","ip":"10.0.44.63","method":"SMEMBERS","port":7303,"proc":"","query":"SMEMBERS restaurant:1483896:contents","redis":{"return_value":"[204413, 204458]"},"resource":"restaurant:1483896:contents","responsetime":184217,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：650.636ms > 311.608ms > 184.217ms
- 耗时命令：GET, SMEMBERS
- 端口关系：20035 <-> 7307(1), 45872 <-> 7301(1), 36673 <-> 7303


----------


- `2core10s_xg-redis-corvusproxy-41_@8003@_20170330161322_20170330161332.log`

```
total transactions : 115756

responsetime(2929241 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:24.380Z","beat":{"hostname":"xg-redis-corvusproxy-41","name":"xg-redis-corvusproxy-41","version":"6.0.0-alpha1"},"bytes_in":378,"bytes_out":406,"client_ip":"10.0.45.115","client_port":46987,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.190","method":"EVAL","port":8003,"proc":"","query":"EVAL \n        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then\n            if ARGV[2] ~= '' then\n                redis.call('pexpire', KEYS[1], ARGV[2])\n            end\n            return 1\n        end\n        return 0\n     1 _lockguldan.thrift.handler.query:deal_schedules|v1:base.guldan|(True,)[] bee33118152011e7adb69ce374424c69 10000","redis":{"return_value":"�\u0002cdogpile.cache.api\nCachedValue\nq\u0001]q\u0002(J�\"\u0003\u0000Mk\u0003JV\"\u0003\u0000M\u0006\u0003J�\"\u0003\u0000J�\"\u0003\u0000J�\"\u0003\u0000JW\"\u0003\u0000J�\"\u0003\u0000Mm\u0003JX!\u0003\u0000J�\"\u0003\u0000JQ\"\u0003\u0000J�\"\u0003\u0000J\u0018!\u0003\u0000J�!\u0003\u0000K\u001bJ�\"\u0003\u0000K\u001dK\u001eK\u001fJ�!\u0003\u0000J!\"\u0003\u0000J�\"\u0003\u0000K#M$\u0008J%\"\u0003\u0000Jq\"\u0003\u0000J)\"\u0003\u0000J*\"\u0003\u0000M\u0016\u0001K\u0018J�\u001f\u0003\u0000J2\"\u0003\u0000K�K\u001cJ�!\u0003\u0000J� \u0003\u0000K$J \"\u0003\u0000K\u000bJE\"\u0003\u0000K JH\"\u0003\u0000JI\"\u0003\u0000JK\"\u0003\u0000K\"JN\"\u0003\u0000JO\"\u0003\u0000JP\"\u0003\u0000J�!\u0003\u0000JR\"\u0003\u0000JS\"\u0003\u0000JT\"\u0003\u0000JU\"\u0003\u0000J�!\u0003\u0000J�!\u0003\u0000JX\u001f\u0003\u0000JY\"\u0003\u0000K\u001aJ�!\u0003\u0000K�K�K�Jf\"\u0003\u0000Jg\"\u0003\u0000J'\"\u0003\u0000KlJm\"\u0003\u0000Jo\"\u0003\u0000J� \u0003\u0000J� \u0003\u0000Jt\"\u0003\u0000Jv\"\u0003\u0000K�K!J\"\u0003\u0000JY!\u0003\u0000J�!\u0003\u0000e}q\u0003(U\u0001vK\u0001U\u0002ctq\u0004GA�7/��s*u�Rq\u0005."},"resource":"\n        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then\n            if ARGV[2] ~= '' then\n                redis.call('pexpire', KEYS[1], ARGV[2])\n            end\n            return 1\n        end\n        return 0\n    ","responsetime":2929241,"server":"","status":"OK","type":"redis"}


responsetime(1406325 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:28.034Z","beat":{"hostname":"xg-redis-corvusproxy-41","name":"xg-redis-corvusproxy-41","version":"6.0.0-alpha1"},"bytes_in":378,"bytes_out":406,"client_ip":"10.0.45.115","client_port":46987,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.190","method":"EVAL","port":8003,"proc":"","query":"EVAL \n        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then\n            if ARGV[2] ~= '' then\n                redis.call('pexpire', KEYS[1], ARGV[2])\n            end\n            return 1\n        end\n        return 0\n     1 _lockguldan.thrift.handler.query:deal_schedules|v1:base.guldan|(True,)[] c110cd88152011e7adb69ce374424c69 10000","redis":{"return_value":"�\u0002cdogpile.cache.api\nCachedValue\nq\u0001]q\u0002(J�\"\u0003\u0000Mk\u0003JV\"\u0003\u0000M\u0006\u0003J�\"\u0003\u0000J�\"\u0003\u0000J�\"\u0003\u0000JW\"\u0003\u0000J�\"\u0003\u0000Mm\u0003JX!\u0003\u0000J�\"\u0003\u0000JQ\"\u0003\u0000J�\"\u0003\u0000J\u0018!\u0003\u0000J�!\u0003\u0000K\u001bJ�\"\u0003\u0000K\u001dK\u001eK\u001fJ�!\u0003\u0000J!\"\u0003\u0000J�\"\u0003\u0000K#M$\u0008J%\"\u0003\u0000Jq\"\u0003\u0000J)\"\u0003\u0000J*\"\u0003\u0000M\u0016\u0001K\u0018J�\u001f\u0003\u0000J2\"\u0003\u0000K�K\u001cJ�!\u0003\u0000J� \u0003\u0000K$J \"\u0003\u0000K\u000bJE\"\u0003\u0000K JH\"\u0003\u0000JI\"\u0003\u0000JK\"\u0003\u0000K\"JN\"\u0003\u0000JO\"\u0003\u0000JP\"\u0003\u0000J�!\u0003\u0000JR\"\u0003\u0000JS\"\u0003\u0000JT\"\u0003\u0000JU\"\u0003\u0000J�!\u0003\u0000J�!\u0003\u0000JX\u001f\u0003\u0000JY\"\u0003\u0000K\u001aJ�!\u0003\u0000K�K�K�Jf\"\u0003\u0000Jg\"\u0003\u0000J'\"\u0003\u0000KlJm\"\u0003\u0000Jo\"\u0003\u0000J� \u0003\u0000J� \u0003\u0000Jt\"\u0003\u0000Jv\"\u0003\u0000K�K!J\"\u0003\u0000JY!\u0003\u0000J�!\u0003\u0000e}q\u0003(U\u0001vK\u0001U\u0002ctq\u0004GA�7/�Z\u003c�u�Rq\u0005."},"resource":"\n        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then\n            if ARGV[2] ~= '' then\n                redis.call('pexpire', KEYS[1], ARGV[2])\n            end\n            return 1\n        end\n        return 0\n    ","responsetime":1406325,"server":"","status":"OK","type":"redis"}


responsetime(1237671 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:23.002Z","beat":{"hostname":"xg-redis-corvusproxy-41","name":"xg-redis-corvusproxy-41","version":"6.0.0-alpha1"},"bytes_in":378,"bytes_out":406,"client_ip":"10.0.45.115","client_port":46987,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.190","method":"EVAL","port":8003,"proc":"","query":"EVAL \n        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then\n            if ARGV[2] ~= '' then\n                redis.call('pexpire', KEYS[1], ARGV[2])\n            end\n            return 1\n        end\n        return 0\n     1 _lockguldan.thrift.handler.query:deal_schedules|v1:base.guldan|(True,)[] be11099a152011e7adb69ce374424c69 10000","redis":{"return_value":"�\u0002cdogpile.cache.api\nCachedValue\nq\u0001]q\u0002(J�\"\u0003\u0000Mk\u0003JV\"\u0003\u0000M\u0006\u0003J�\"\u0003\u0000J�\"\u0003\u0000J�\"\u0003\u0000JW\"\u0003\u0000J�\"\u0003\u0000Mm\u0003JX!\u0003\u0000J�\"\u0003\u0000JQ\"\u0003\u0000J�\"\u0003\u0000J\u0018!\u0003\u0000J�!\u0003\u0000K\u001bJ�\"\u0003\u0000K\u001dK\u001eK\u001fJ�!\u0003\u0000J!\"\u0003\u0000J�\"\u0003\u0000K#M$\u0008J%\"\u0003\u0000Jq\"\u0003\u0000J)\"\u0003\u0000J*\"\u0003\u0000M\u0016\u0001K\u0018J�\u001f\u0003\u0000J2\"\u0003\u0000K�K\u001cJ�!\u0003\u0000J� \u0003\u0000K$J \"\u0003\u0000K\u000bJE\"\u0003\u0000K JH\"\u0003\u0000JI\"\u0003\u0000JK\"\u0003\u0000K\"JN\"\u0003\u0000JO\"\u0003\u0000JP\"\u0003\u0000J�!\u0003\u0000JR\"\u0003\u0000JS\"\u0003\u0000JT\"\u0003\u0000JU\"\u0003\u0000J�!\u0003\u0000J�!\u0003\u0000JX\u001f\u0003\u0000JY\"\u0003\u0000K\u001aJ�!\u0003\u0000K�K�K�Jf\"\u0003\u0000Jg\"\u0003\u0000J'\"\u0003\u0000KlJm\"\u0003\u0000Jo\"\u0003\u0000J� \u0003\u0000J� \u0003\u0000Jt\"\u0003\u0000Jv\"\u0003\u0000K�K!J\"\u0003\u0000JY!\u0003\u0000J�!\u0003\u0000e}q\u0003(U\u0001vK\u0001U\u0002ctq\u0004GA�7/�\u000e��u�Rq\u0005."},"resource":"\n        if redis.call('setnx', KEYS[1], ARGV[1]) == 1 then\n            if ARGV[2] ~= '' then\n                redis.call('pexpire', KEYS[1], ARGV[2])\n            end\n            return 1\n        end\n        return 0\n    ","responsetime":1237671,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：2.929241s > 1.406325s > 1.237671s
- 耗时命令：均为 EVAL
- 端口关系：46987 <-> 8003(3)



----------


- `2core10s_xg-redis-corvusproxy-43_@8006@_20170330161322_20170330161332.log`

```
total transactions : 90259

responsetime(875145 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:28.868Z","beat":{"hostname":"xg-redis-corvusproxy-43","name":"xg-redis-corvusproxy-43","version":"6.0.0-alpha1"},"bytes_in":56,"bytes_out":4,"client_ip":"10.0.18.77","client_port":23985,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.188","method":"SMEMBERS","port":8006,"proc":"","query":"SMEMBERS food_restaurant:929064:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:929064:contents","responsetime":875145,"server":"","status":"OK","type":"redis"}


responsetime(875145 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:28.868Z","beat":{"hostname":"xg-redis-corvusproxy-43","name":"xg-redis-corvusproxy-43","version":"6.0.0-alpha1"},"bytes_in":51,"bytes_out":4,"client_ip":"10.0.18.77","client_port":23985,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.188","method":"SMEMBERS","port":8006,"proc":"","query":"SMEMBERS restaurant:115013:contents","redis":{"return_value":"[]"},"resource":"restaurant:115013:contents","responsetime":875145,"server":"","status":"OK","type":"redis"}


responsetime(875145 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:28.868Z","beat":{"hostname":"xg-redis-corvusproxy-43","name":"xg-redis-corvusproxy-43","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.18.77","client_port":23985,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.188","method":"SMEMBERS","port":8006,"proc":"","query":"SMEMBERS food_restaurant:1504085:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:1504085:contents","responsetime":875145,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：875.145ms = 875.145ms = 875.145ms
- 耗时命令：均为 SMEMBERS
- 端口信息：23985 <-> 8006(3)


----------


- `2core10s_xg-redis-corvusproxy-44_@8003@_20170330161322_20170330161332.log`

```
total transactions : 109667

responsetime(1306936 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:29.350Z","beat":{"hostname":"xg-redis-corvusproxy-44","name":"xg-redis-corvusproxy-44","version":"6.0.0-alpha1"},"bytes_in":81,"bytes_out":1612,"client_ip":"10.0.45.116","client_port":28670,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.185","method":"GET","port":8003,"proc":"","query":"GET cms_device_admin_by_uuid:9da5ca36-a8b1-3a08-814d-c35744a7bbdb","redis":{"return_value":"(lp1\nI1383241\naI1403966\naI1475838\naI150110871\naI151705378\naI1564276\naI1473651\naI1510009\naI1142103\naI2157574\naI1533827\naI1998038\naI1424302\naI1196703\naI2076133\naI1447094\naI152106990\naI151693436\naI1480709\naI1282781\naI1282783\naI1466886\naI152107343\naI1196134\naI2111115\naI1172229\naI1522337\naI1178219\naI2150535\naI1363054\naI1283201\naI1178211\naI1265271\naI1143363\naI1515548\naI142261599\naI151705745\naI1480135\naI1385687\naI152170626\naI152110105\naI2346648\naI1172562\naI1244334\naI2238238\naI152106729\naI1170044\naI1347533\naI1142089\naI1482685\naI1542092\naI152170713\naI150124479\naI2376215\naI1220796\naI2132171\naI151706123\naI1250463\naI150994996\naI2365893\naI151706450\naI151713044\naI2300949\naI1172162\naI2300572\naI2028962\naI151710262\naI2102630\naI1178233\naI151709548\naI1318552\naI882128\naI152105580\naI150055159\naI1171525\naI150154534\naI151706238\naI150081533\naI1196480\naI1303846\naI150986268\naI1271244\naI1318710\naI1510042\naI1465975\naI150011473\naI1371537\naI1810090\naI1813592\naI151653900\naI1500689\naI2154537\naI150995453\naI1244164\naI2159117\naI1397472\naI152107193\naI2107408\naI1538558\naI1480997\naI1507939\naI1405864\naI1510830\naI1295140\naI152162151\naI1809756\naI2082929\naI2351732\naI152162012\naI150109199\naI152168144\naI1218712\naI1515811\naI2056287\naI1403943\naI1515529\naI150082492\naI151002469\naI2023521\naI1984117\naI1481228\naI1884729\naI152110224\naI1318730\naI2271673\naI1478976\naI1527189\naI151002829\naI150846364\naI151693571\naI1529698\naI151709665\naI1992481\naI1146713\naI151690325\naI1831013\naI1931268\naI1202445\naI1539268\naI151002215\naI1870160\naI1515640\naI150154828\naI151693219\naI2111080\naI2029781\naI1206426\naI1368360\naI150850501\naI1327061\naI1810046\na."},"resource":"cms_device_admin_by_uuid:9da5ca36-a8b1-3a08-814d-c35744a7bbdb","responsetime":1306936,"server":"","status":"OK","type":"redis"}


responsetime(914311 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:27.972Z","beat":{"hostname":"xg-redis-corvusproxy-44","name":"xg-redis-corvusproxy-44","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":26,"client_ip":"10.0.12.122","client_port":30380,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.185","method":"SMEMBERS","port":8003,"proc":"","query":"SMEMBERS food_restaurant:1954746:contents","redis":{"return_value":"[3258, 203727]"},"resource":"food_restaurant:1954746:contents","responsetime":914311,"server":"","status":"OK","type":"redis"}


responsetime(914311 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:27.972Z","beat":{"hostname":"xg-redis-corvusproxy-44","name":"xg-redis-corvusproxy-44","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.12.122","client_port":30380,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.185","method":"SMEMBERS","port":8003,"proc":"","query":"SMEMBERS restaurant:1954746:contents","redis":{"return_value":"[]"},"resource":"restaurant:1954746:contents","responsetime":914311,"server":"","status":"OK","type":"redis"}
```


结论：

- 耗时时间：1.306936s > 914.311ms = 914.311ms
- 耗时命令：GET, SMEMBERS
- 端口关系：28670 <-> 8003(1), 30380 <-> 8003(2)


----------


- `2core10s_xg-redis-corvusproxy-45_@8005@_20170330161321_20170330161331.log`

```
total transactions : 104284

responsetime(763994 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:30.601Z","beat":{"hostname":"xg-redis-corvusproxy-45","name":"xg-redis-corvusproxy-45","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":5,"client_ip":"10.0.13.26","client_port":62060,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.183","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS restaurant:1971574:contents","redis":{"return_value":"nil"},"resource":"restaurant:1971574:contents","responsetime":763994,"server":"","status":"OK","type":"redis"}


responsetime(663541 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:25.951Z","beat":{"hostname":"xg-redis-corvusproxy-45","name":"xg-redis-corvusproxy-45","version":"6.0.0-alpha1"},"bytes_in":59,"bytes_out":251,"client_ip":"10.0.21.49","client_port":40495,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.183","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS food_restaurant:150851816:contents","redis":{"return_value":"(lp1\nccopy_reg\n_reconstructor\np2\n(cme_ele_base_search_range_thrift\nSourceTypeWallSearchResult__SingleResult\np3\nc__builtin__\nobject\np4\nNtRp5\n(dp6\nS'type'\np7\nVrange_walle\np8\nsS'id'\np9\nV30953\np10\nsbag2\n(g3\ng4\nNtRp11\n(dp12\ng7\nV2\nsg9\nV8343\np13\nsba."},"resource":"food_restaurant:150851816:contents","responsetime":663541,"server":"","status":"OK","type":"redis"}


responsetime(175609 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:29.850Z","beat":{"hostname":"xg-redis-corvusproxy-45","name":"xg-redis-corvusproxy-45","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":7,"client_ip":"10.0.21.41","client_port":48328,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.183","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS restaurant:2215398:contents","redis":{"return_value":"0"},"resource":"restaurant:2215398:contents","responsetime":175609,"server":"","status":"OK","type":"redis"}

```

结论：

- 耗时时间：763.994ms > 663.541ms > 175.609ms 
- 耗时命令：均为 SMEMBERS
- 端口关系：62060 <-> 8005(1), 40495 <-> 8005(1), 48328 <-> 8005(1)


----------


- `2core10s_xg-redis-corvusproxy-46_@8005@_20170330161322_20170330161332.log`


```
total transactions : 106601

responsetime(1580521 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:27.571Z","beat":{"hostname":"xg-redis-corvusproxy-46","name":"xg-redis-corvusproxy-46","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":252,"client_ip":"10.0.45.118","client_port":51972,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.181","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS food_restaurant:1440658:contents","redis":{"return_value":"(lp1\nccopy_reg\n_reconstructor\np2\n(cme_ele_base_search_range_thrift\nSourceTypeWallSearchResult__SingleResult\np3\nc__builtin__\nobject\np4\nNtRp5\n(dp6\nS'type'\np7\nV2\nsS'id'\np8\nV10165\np9\nsbag2\n(g3\ng4\nNtRp10\n(dp11\ng7\nVrange_walle\np12\nsg8\nV10165\np13\nsba."},"resource":"food_restaurant:1440658:contents","responsetime":1580521,"server":"","status":"OK","type":"redis"}


responsetime(1414711 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:27.571Z","beat":{"hostname":"xg-redis-corvusproxy-46","name":"xg-redis-corvusproxy-46","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":252,"client_ip":"10.0.45.118","client_port":51972,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.181","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS restaurant:1440658:contents","redis":{"return_value":"(lp1\nccopy_reg\n_reconstructor\np2\n(cme_ele_base_search_range_thrift\nSourceTypeWallSearchResult__SingleResult\np3\nc__builtin__\nobject\np4\nNtRp5\n(dp6\nS'type'\np7\nV2\nsS'id'\np8\nV13313\np9\nsbag2\n(g3\ng4\nNtRp10\n(dp11\ng7\nVrange_walle\np12\nsg8\nV13313\np13\nsba."},"resource":"restaurant:1440658:contents","responsetime":1414711,"server":"","status":"OK","type":"redis"}


responsetime(1371968 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:28.708Z","beat":{"hostname":"xg-redis-corvusproxy-46","name":"xg-redis-corvusproxy-46","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.21.41","client_port":43833,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.181","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS restaurant:2049049:contents","redis":{"return_value":"[]"},"resource":"restaurant:2049049:contents","responsetime":1371968,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.580521s > 1.414711s > 1.371968s 
- 耗时命令：均为 SMEMBERS
- 端口关系：51972 <-> 8005(2), 43833 <-> 8005(1)


----------


- `2core10s_xg-redis-corvusproxy-48_@8005@_20170330161322_20170330161332.log`

```
total transactions : 97125

responsetime(1954545 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:29.145Z","beat":{"hostname":"xg-redis-corvusproxy-48","name":"xg-redis-corvusproxy-48","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.21.33","client_port":56615,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.178","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS restaurant:1174160:contents","redis":{"return_value":"[]"},"resource":"restaurant:1174160:contents","responsetime":1954545,"server":"","status":"OK","type":"redis"}


responsetime(1954545 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:29.145Z","beat":{"hostname":"xg-redis-corvusproxy-48","name":"xg-redis-corvusproxy-48","version":"6.0.0-alpha1"},"bytes_in":59,"bytes_out":4,"client_ip":"10.0.21.33","client_port":56615,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.178","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS food_restaurant:150972360:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:150972360:contents","responsetime":1954545,"server":"","status":"OK","type":"redis"}


responsetime(1954545 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:29.145Z","beat":{"hostname":"xg-redis-corvusproxy-48","name":"xg-redis-corvusproxy-48","version":"6.0.0-alpha1"},"bytes_in":54,"bytes_out":4,"client_ip":"10.0.21.33","client_port":56615,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.178","method":"SMEMBERS","port":8005,"proc":"","query":"SMEMBERS restaurant:150972360:contents","redis":{"return_value":"[]"},"resource":"restaurant:150972360:contents","responsetime":1954545,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.954545s = 1.954545s = 1.954545s 
- 耗时命令：均为 SMEMBERS
- 端口关系：56615 <-> 8005(3)


----------



- `2core10s_xg-redis-corvusproxy-49_@8002@_20170330161322_20170330161332.log`

```
total transactions : 107602

responsetime(1678905 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:23.909Z","beat":{"hostname":"xg-redis-corvusproxy-49","name":"xg-redis-corvusproxy-49","version":"6.0.0-alpha1"},"bytes_in":54,"bytes_out":250,"client_ip":"10.0.45.115","client_port":62725,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.169","method":"SMEMBERS","port":8002,"proc":"","query":"SMEMBERS restaurant:151667784:contents","redis":{"return_value":"(lp1\nccopy_reg\n_reconstructor\np2\n(cme_ele_base_search_range_thrift\nSourceTypeWallSearchResult__SingleResult\np3\nc__builtin__\nobject\np4\nNtRp5\n(dp6\nS'type'\np7\nV2\nsS'id'\np8\nV1842\np9\nsbag2\n(g3\ng4\nNtRp10\n(dp11\ng7\nVrange_walle\np12\nsg8\nV1842\np13\nsba."},"resource":"restaurant:151667784:contents","responsetime":1678905,"server":"","status":"OK","type":"redis"}


responsetime(1066459 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:29.903Z","beat":{"hostname":"xg-redis-corvusproxy-49","name":"xg-redis-corvusproxy-49","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":250,"client_ip":"10.0.45.115","client_port":21233,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.169","method":"SMEMBERS","port":8002,"proc":"","query":"SMEMBERS restaurant:1382695:contents","redis":{"return_value":"(lp1\nccopy_reg\n_reconstructor\np2\n(cme_ele_base_search_range_thrift\nSourceTypeWallSearchResult__SingleResult\np3\nc__builtin__\nobject\np4\nNtRp5\n(dp6\nS'type'\np7\nV2\nsS'id'\np8\nV1930\np9\nsbag2\n(g3\ng4\nNtRp10\n(dp11\ng7\nVrange_walle\np12\nsg8\nV1930\np13\nsba."},"resource":"restaurant:1382695:contents","responsetime":1066459,"server":"","status":"OK","type":"redis"}


responsetime(1030138 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:26.698Z","beat":{"hostname":"xg-redis-corvusproxy-49","name":"xg-redis-corvusproxy-49","version":"6.0.0-alpha1"},"bytes_in":59,"bytes_out":4,"client_ip":"10.0.12.65","client_port":28946,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.169","method":"SMEMBERS","port":8002,"proc":"","query":"SMEMBERS food_restaurant:150854199:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:150854199:contents","responsetime":1030138,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.678905s > 1.066459s > 1.030138s 
- 耗时命令：均为 SMEMBERS
- 端口关系：62725 <-> 8002(1), 21233 <-> 8002(1), 28946 <-> 8002(1)


----------


- `2core10s_xg-redis-corvusproxy-50_@8006@_20170330161322_20170330161333.log`

```
total transactions : 87123

responsetime(1500756 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-03-30T08:13:27.388Z","beat":{"hostname":"xg-redis-corvusproxy-50","name":"xg-redis-corvusproxy-50","version":"6.0.0-alpha1"},"bytes_in":57,"bytes_out":4,"client_ip":"10.0.45.118","client_port":54967,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.161","method":"SMEMBERS","port":8006,"proc":"","query":"SMEMBERS food_restaurant:1502929:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:1502929:contents","responsetime":1500756,"server":"","status":"OK","type":"redis"}


responsetime(1500756 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-03-30T08:13:27.388Z","beat":{"hostname":"xg-redis-corvusproxy-50","name":"xg-redis-corvusproxy-50","version":"6.0.0-alpha1"},"bytes_in":52,"bytes_out":4,"client_ip":"10.0.45.118","client_port":54967,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.161","method":"SMEMBERS","port":8006,"proc":"","query":"SMEMBERS restaurant:1537842:contents","redis":{"return_value":"[]"},"resource":"restaurant:1537842:contents","responsetime":1500756,"server":"","status":"OK","type":"redis"}


responsetime(1500756 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-03-30T08:13:27.388Z","beat":{"hostname":"xg-redis-corvusproxy-50","name":"xg-redis-corvusproxy-50","version":"6.0.0-alpha1"},"bytes_in":56,"bytes_out":4,"client_ip":"10.0.45.118","client_port":54967,"client_proc":"","client_server":"","direction":"in","ip":"10.0.58.161","method":"SMEMBERS","port":8006,"proc":"","query":"SMEMBERS food_restaurant:484024:contents","redis":{"return_value":"[]"},"resource":"food_restaurant:484024:contents","responsetime":1500756,"server":"","status":"OK","type":"redis"}
```

结论：

- 耗时时间：1.500756s = 1.500756s = 1.500756s 
- 耗时命令：均为 SMEMBERS
- 端口关系：54967 <-> 8006(3)


----------


## 结论

- redis server 上最长 responsetime 为 1.26s ，对应的命令为 SMEMBERS；
- corvus proxy 上最长 responsetime 为 2.92s ，对应的命令为 EVAL ；


> 遗留问题：
> 
> - 上述结果不排除 packetbeat 进行 request-response 匹配时存在 bug ，进而导致数据“糟糕”；（需要通过其他方式确认数据的准确性）
> - 在进行该测试时，由于 corvus proxy 没有提供相应的慢查询日志供比对，因此是否存在数据中的延迟，以及延迟是否正常的合理情况，尚不清楚；
> - 通过 redis-cli 登录目标 redis server 查看 slow log 日志，全都是 `CLUSTER NODES` 命令；
> - 需要确认 packetbeat 如何饿到 packet 的 timestamp ，以及如何计算 responsetime 的；
