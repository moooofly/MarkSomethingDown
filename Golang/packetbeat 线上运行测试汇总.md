# packetbeat 线上运行测试

标签（空格分隔）： packetbeat

---

## redis 测试

### 测试中数据流

`packetbeat` 实时抓取 `bond0` 上的 `redis` 协议数据包，并进行 request-response 匹配，最终将封装成 json 结构的匹配信息写入本地文件；

基于 `python` 脚本对上述文件中的内容进行 topN 计算，得到 request-response 延迟最大的一组数据；


### 测试结果

测试命令（抓取 10 个 port 上的 `redis` 通信）

```
[root@xg-bigkey-rediscluster-1 packageCaptureAnalysis]# LD_LIBRARY_PATH=. ./packageCaptureAnalysis -c ./packetbeat.yml -E packetbeat.protocols.redis.ports=7101,7102,7103,7104,7105,7106,7107,7108,7109,7110
^C
[root@xg-bigkey-rediscluster-1 packageCaptureAnalysis]#
```

`top` 输出（24 核，pb 运行平均占用一个核）

```
top - 17:46:52 up 311 days,  3:31,  2 users,  load average: 0.71, 0.35, 0.28
Tasks: 369 total,   2 running, 367 sleeping,   0 stopped,   0 zombie
%Cpu0  :  8.8 us,  2.4 sy,  0.0 ni, 88.2 id,  0.0 wa,  0.0 hi,  0.7 si,  0.0 st
%Cpu1  : 12.1 us,  1.7 sy,  0.0 ni, 85.9 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu2  :  8.8 us,  1.7 sy,  0.0 ni, 89.5 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu3  :  9.1 us,  2.0 sy,  0.0 ni, 88.6 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu4  : 10.4 us,  2.3 sy,  0.0 ni, 86.9 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu5  :  8.1 us,  2.0 sy,  0.0 ni, 89.6 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu6  :  1.7 us,  0.3 sy,  0.0 ni, 98.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu7  :  4.7 us,  0.3 sy,  0.0 ni, 95.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu8  :  3.7 us,  0.7 sy,  0.0 ni, 95.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu9  :  3.7 us,  0.7 sy,  0.0 ni, 95.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu10 :  2.7 us,  0.3 sy,  0.0 ni, 97.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu11 :  1.0 us,  0.0 sy,  0.0 ni, 99.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu12 :  6.4 us,  1.4 sy,  0.0 ni, 91.9 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu13 :  6.8 us,  1.0 sy,  0.0 ni, 91.9 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu14 :  4.0 us,  1.7 sy,  0.0 ni, 94.0 id,  0.0 wa,  0.0 hi,  0.3 si,  0.0 st
%Cpu15 :  2.7 us,  0.7 sy,  0.0 ni, 96.6 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu16 :  5.7 us,  2.0 sy,  0.0 ni, 92.3 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu17 :  7.4 us,  1.7 sy,  0.0 ni, 90.3 id,  0.0 wa,  0.0 hi,  0.7 si,  0.0 st
%Cpu18 :  2.0 us,  0.3 sy,  0.0 ni, 97.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu19 :  4.3 us,  0.0 sy,  0.0 ni, 95.7 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu20 :  0.3 us,  0.3 sy,  0.0 ni, 99.3 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu21 :  1.7 us,  0.3 sy,  0.0 ni, 98.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu22 :  1.3 us,  0.3 sy,  0.0 ni, 98.3 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
%Cpu23 :  3.3 us,  0.7 sy,  0.0 ni, 96.0 id,  0.0 wa,  0.0 hi,  0.0 si,  0.0 st
KiB Mem : 98437376 total, 28629740 free, 60484752 used,  9322880 buff/cache
KiB Swap: 16383996 total, 16383996 free,        0 used. 37403828 avail Mem

  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
10766 root      20   0 1995004 829012   7188 S 127.0  0.8   2:02.10 packageCaptureA
14969 redis     20   0 5816220 5.424g   1596 S  13.0  5.8   7619:27 redis-server
13151 redis     20   0  183832  50288   1556 S   8.3  0.1   1332:16 redis-server
 3509 root      20   0   28756  18136   5352 R   2.7  0.0 652:00.40 esm-agent
14974 redis     20   0 5795352 5.402g   1580 S   1.0  5.8   3592:05 redis-server
14972 redis     20   0 6348312 5.293g   1592 S   0.7  5.6   5296:49 redis-server
   34 root      20   0       0      0      0 S   0.3  0.0 424:06.44 rcu_sched
 1781 root      20   0  353848  40768   5148 S   0.3  0.0  47:16.38 corvus_web
10717 root      20   0  146376   2348   1432 R   0.3  0.0   0:00.31 top
    1 root      20   0   68048  30556   2416 S   0.0  0.0   8:16.97 systemd
    2 root      20   0       0      0      0 S   0.0  0.0   0:02.30 kthreadd
```

pb 运行大约 7 分钟，保存到文件中的分析结果占用大约 1.8G ；

```
[root@xg-bigkey-rediscluster-1 logs]# ps aux|grep packet|grep -v grep ;ll -h
root     10766  118  1.7 2878708 1714756 pts/0 Sl+  17:45   8:00 ./packageCaptureAnalysis -c ./packetbeat.yml -E packetbeat.protocols.redis.ports=7101,7102,7103,7104,7105,7106,7107,7108,7109,7110
total 1.8G
-rw-r--r-- 1 root root  67M Feb 21 17:51 packetbeat
-rw-r--r-- 1 root root 101M Feb 21 17:51 packetbeat.1
-rw-r--r-- 1 root root 101M Feb 21 17:48 packetbeat.10
-rw-r--r-- 1 root root 101M Feb 21 17:47 packetbeat.11
-rw-r--r-- 1 root root 101M Feb 21 17:47 packetbeat.12
-rw-r--r-- 1 root root 101M Feb 21 17:46 packetbeat.13
-rw-r--r-- 1 root root 101M Feb 21 17:46 packetbeat.14
-rw-r--r-- 1 root root 101M Feb 21 17:46 packetbeat.15
-rw-r--r-- 1 root root 101M Feb 21 17:45 packetbeat.16
-rw-r--r-- 1 root root 101M Feb 21 17:45 packetbeat.17
-rw-r--r-- 1 root root 2.6K Feb 21 17:51 packetbeat.18
-rw-r--r-- 1 root root  661 Feb 21 17:45 packetbeat.19
-rw-r--r-- 1 root root 101M Feb 21 17:51 packetbeat.2
-rw-r--r-- 1 root root 101M Feb 21 17:50 packetbeat.3
-rw-r--r-- 1 root root 101M Feb 21 17:50 packetbeat.4
-rw-r--r-- 1 root root 101M Feb 21 17:50 packetbeat.5
-rw-r--r-- 1 root root 101M Feb 21 17:49 packetbeat.6
-rw-r--r-- 1 root root 101M Feb 21 17:49 packetbeat.7
-rw-r--r-- 1 root root 101M Feb 21 17:48 packetbeat.8
-rw-r--r-- 1 root root 101M Feb 21 17:48 packetbeat.9
[root@xg-bigkey-rediscluster-1 logs]# du -shx .
1.8G	.
```

pb 输出的统计结果

```
2017-02-21T17:51:53+08:00 INFO Input finish. Processed 12383127 packets. Have a nice day!
2017-02-21T17:51:53+08:00 INFO Total non-zero values:  libbeat.publisher.published_events=2956771 libbeat.publisher.messages_in_worker_queues=116 redis.unmatched_responses=8763 tcp.dropped_because_of_gaps=2735
2017-02-21T17:51:53+08:00 INFO Uptime: 6m49.384982569s
```

基于 python 脚本进行 topN 分析（耗费大约 2 分钟）

```
[root@xg-bigkey-rediscluster-1 packageCaptureAnalysis]# time python redis_analysis.py -p logs_bak -f packetbeat,packetbeat.1,packetbeat.2,packetbeat.3,packetbeat.4,packetbeat.5,packetbeat.6,packetbeat.7,packetbeat.8,packetbeat.9,packetbeat.10,packetbeat.11,packetbeat.12,packetbeat.13,packetbeat.14,packetbeat.15,packetbeat.16,packetbeat.17 -t 10


@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

total transactions : 2956656
total failure nums : 2
failure rate       : 0.000068%

@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

responsetime(143457699 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-02-21T09:47:47.122Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":66,"bytes_out":47,"client_ip":"10.0.27.40","client_port":7102,"client_proc":"","client_server":"","direction":"out","ip":"10.0.28.25","method":"EXPIRE","port":52180,"proc":"","query":"EXPIRE app:hotfood:187147299:foodclick 172800","redis":{"return_value":"[REPLCONF, ACK, 2311648430835]"},"resource":"app:hotfood:187147299:foodclick","responsetime":143457699,"server":"","status":"OK","type":"redis"}


responsetime(143444975 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-02-21T09:47:47.152Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":82,"bytes_out":47,"client_ip":"10.0.28.25","client_port":7103,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.40","method":"HSET","port":31182,"proc":"","query":"HSET app:hotfood:159897675:shopclick 1414481 1487670442","redis":{"return_value":"[REPLCONF, ACK, 2305174778527]"},"resource":"app:hotfood:159897675:shopclick","responsetime":143444975,"server":"","status":"OK","type":"redis"}


responsetime(143294293 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-02-21T09:47:47.116Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":83,"bytes_out":47,"client_ip":"10.0.28.25","client_port":7107,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.40","method":"HSET","port":36024,"proc":"","query":"HSET app:hotfood:40827992:search update_time 1487670467","redis":{"return_value":"[REPLCONF, ACK, 2298184951759]"},"resource":"app:hotfood:40827992:search","responsetime":143294293,"server":"","status":"OK","type":"redis"}


responsetime(143231599 microseconds)	==>    No.<4>
----
{"@timestamp":"2017-02-21T09:47:47.235Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":93,"bytes_out":47,"client_ip":"10.0.28.25","client_port":7109,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.40","method":"EXPIRE","port":47650,"proc":"","query":"EXPIRE app:hotfood:2C537BAE-CE37-45E6-A4D5-EE014DB35B4C:foodclick 172800","redis":{"return_value":"[REPLCONF, ACK, 2296028618540]"},"resource":"app:hotfood:2C537BAE-CE37-45E6-A4D5-EE014DB35B4C:foodclick","responsetime":143231599,"server":"","status":"OK","type":"redis"}


responsetime(143186975 microseconds)	==>    No.<5>
----
{"@timestamp":"2017-02-21T09:47:47.280Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":114,"bytes_out":47,"client_ip":"10.0.27.40","client_port":7107,"client_proc":"","client_server":"","direction":"out","ip":"10.0.29.23","method":"HSET","port":60426,"proc":"","query":"HSET app:hotfood:5994AE71-AE42-4B38-8BA1-8313E341AB27:shopclick update_time 1487670467","redis":{"return_value":"[REPLCONF, ACK, 2298803422475]"},"resource":"app:hotfood:5994AE71-AE42-4B38-8BA1-8313E341AB27:shopclick","responsetime":143186975,"server":"","status":"OK","type":"redis"}


responsetime(143162442 microseconds)	==>    No.<6>
----
{"@timestamp":"2017-02-21T09:47:47.246Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":111,"bytes_out":47,"client_ip":"10.0.27.40","client_port":7105,"client_proc":"","client_server":"","direction":"out","ip":"10.0.29.23","method":"HSET","port":35231,"proc":"","query":"HSET app:hotfood:2B679B89-70E6-4F0F-A6A0-989B04376669:shopclick 150024477 1487670441","redis":{"return_value":"[REPLCONF, ACK, 2290388973817]"},"resource":"app:hotfood:2B679B89-70E6-4F0F-A6A0-989B04376669:shopclick","responsetime":143162442,"server":"","status":"OK","type":"redis"}


responsetime(143042767 microseconds)	==>    No.<7>
----
{"@timestamp":"2017-02-21T09:47:47.250Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":93,"bytes_out":47,"client_ip":"10.0.27.40","client_port":7109,"client_proc":"","client_server":"","direction":"out","ip":"10.0.29.23","method":"EXPIRE","port":57062,"proc":"","query":"EXPIRE app:hotfood:55B530EF-99D6-4E8A-AA68-150D3A5D5644:shopclick 172800","redis":{"return_value":"[REPLCONF, ACK, 2302142943147]"},"resource":"app:hotfood:55B530EF-99D6-4E8A-AA68-150D3A5D5644:shopclick","responsetime":143042767,"server":"","status":"OK","type":"redis"}


responsetime(142726281 microseconds)	==>    No.<8>
----
{"@timestamp":"2017-02-21T09:47:47.138Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":93,"bytes_out":47,"client_ip":"10.0.27.40","client_port":7101,"client_proc":"","client_server":"","direction":"out","ip":"10.0.29.23","method":"EXPIRE","port":29811,"proc":"","query":"EXPIRE app:hotfood:969a380e-925d-3cd2-b634-4af9cbdcf215:foodclick 172800","redis":{"return_value":"[REPLCONF, ACK, 2323417606067]"},"resource":"app:hotfood:969a380e-925d-3cd2-b634-4af9cbdcf215:foodclick","responsetime":142726281,"server":"","status":"OK","type":"redis"}


responsetime(142483891 microseconds)	==>    No.<9>
----
{"@timestamp":"2017-02-21T09:47:47.227Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":108,"bytes_out":47,"client_ip":"10.0.28.25","client_port":7105,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.40","method":"HSET","port":38744,"proc":"","query":"HSET app:hotfood:59dcd09d-83ed-343c-9149-cea86176830a:shopclick 800050 1487670439","redis":{"return_value":"[REPLCONF, ACK, 2304936939673]"},"resource":"app:hotfood:59dcd09d-83ed-343c-9149-cea86176830a:shopclick","responsetime":142483891,"server":"","status":"OK","type":"redis"}


responsetime(142475195 microseconds)	==>    No.<10>
----
{"@timestamp":"2017-02-21T09:47:47.234Z","beat":{"hostname":"xg-bigkey-rediscluster-1","name":"xg-bigkey-rediscluster-1","version":"6.0.0-alpha1"},"bytes_in":66,"bytes_out":47,"client_ip":"10.0.29.23","client_port":7104,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.40","method":"EXPIRE","port":24908,"proc":"","query":"EXPIRE app:hotfood:183753799:shopclick 172800","redis":{"return_value":"[REPLCONF, ACK, 2314157704993]"},"resource":"app:hotfood:183753799:shopclick","responsetime":142475195,"server":"","status":"OK","type":"redis"}



real	1m56.783s
user	1m55.540s
sys	0m1.333s
[root@xg-bigkey-rediscluster-1 packageCaptureAnalysis]#
```

小结：

> 由于 `packetbeat` 源码中没有针对抓取到的 `redis` 协议包区分来自 Client-Server 侧，还是来自 master-slave 侧；因此上述分析数据存在一点问题（后续修复），此处结果仅做演示使用；

### 监控输出

- cpu

![运行 packetbeat 分析 redis 时的 cpu 资源使用情况-1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20redis%20%E6%97%B6%E7%9A%84%20cpu%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-1.png)

![运行 packetbeat 分析 redis 时的 cpu 资源使用情况-2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20redis%20%E6%97%B6%E7%9A%84%20cpu%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-2.png)

- disk

![运行 packetbeat 分析 redis 时的 disk 资源使用情况-1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20redis%20%E6%97%B6%E7%9A%84%20disk%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-1.png)

![运行 packetbeat 分析 redis 时的 disk 资源使用情况-2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20redis%20%E6%97%B6%E7%9A%84%20disk%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-2.png)



----------


## mysql 测试

### 测试中数据流

`packetbeat` 实时抓取 `bond0` 上的 `mysql` 协议数据包，并进行 request-response 匹配，最终将封装后的匹配信息写入本地文件；

基于 `python` 脚本对上述文件中的内容进行 topN 计算，得到 request-response 延迟最大的一组数据；


### 测试结果（针对 slave mysql）

测试命令（抓去 3306 上的 slave mysql 通信 7 分钟）

```
[root@xg-restaurant-slave-2 packageCaptureAnalysis]# LD_LIBRARY_PATH=. ./packageCaptureAnalysis -c ./packetbeat.yml
^C
[root@xg-restaurant-slave-2 packageCaptureAnalysis]#
```

top 输出（32 核，pb 运行平均占用 6% 左右）

```
top - 18:25:17 up 89 days,  3:04,  3 users,  load average: 2.23, 2.12, 2.10
Tasks: 574 total,   1 running, 573 sleeping,   0 stopped,   0 zombie
Cpu(s):  0.7%us,  0.2%sy,  0.0%ni, 99.0%id,  0.1%wa,  0.0%hi,  0.0%si,  0.0%st
Mem:  132045660k total, 117768888k used, 14276772k free,   457824k buffers
Swap: 16383996k total,        0k used, 16383996k free, 33700564k cached

   PID USER      PR  NI  VIRT  RES  SHR S %CPU %MEM    TIME+  COMMAND
 80408 mysql     20   0 80.5g  76g 9536 S 18.9 61.0   5596:11 mysqld
 68675 root      20   0 1538m  80m 9168 S  5.6  0.1   0:13.91 packageCaptureA
 55336 root      20   0 32064  19m 5136 S  3.0  0.0 242:36.90 esm-agent
   203 root      20   0     0    0    0 S  0.3  0.0  28:38.63 kblockd/1
  3442 root      20   0     0    0    0 S  0.3  0.0 138:14.52 jbd2/sdb-8
 68268 root      20   0 15432 1632  924 R  0.3  0.0   0:01.14 top
     1 root      20   0 19360 1552 1224 S  0.0  0.0   0:32.46 init
     2 root      20   0     0    0    0 S  0.0  0.0   0:00.00 kthreadd
     3 root      RT   0     0    0    0 S  0.0  0.0   0:02.11 migration/0
```

pb 输出的统计结果

```
2017-02-22T18:55:06+08:00 INFO Input finish. Processed 103462 packets. Have a nice day!
2017-02-22T18:55:06+08:00 INFO Total non-zero values:  tcp.dropped_because_of_gaps=31 libbeat.publisher.published_events=14924 mysql.unmatched_responses=15496
2017-02-22T18:55:06+08:00 INFO Uptime: 4m28.81701879s
```

基于 python 脚本进行 topN 分析（耗费大约 0.34 秒）

```
[root@xg-restaurant-slave-2 packageCaptureAnalysis]# python redis_analysis.py -p logs -f packetbeat -t 10


@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

total transactions : 14924
total failure nums : 0
failure rate       : 0.000000%

@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

responsetime(53 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-02-22T10:54:00.018Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":378,"bytes_out":9873,"client_ip":"10.0.13.33","client_port":47175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":10,"num_rows":101},"path":"eleme_restaurant.t_ord_order_serial","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^1047058104948068138|1487760840018\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select\n         \n        id,eleme_order_id,tp_order_id,new_status,extra,is_manage,create_time,create_by,modify_time,modify_by\n     \n        from t_ord_order_serial where is_manage = 0 and create_time \u003e= date_sub(curdate(),interval 1 day) ORDER BY new_status desc,id asc limit 800;","responsetime":53,"server":"","status":"OK","type":"mysql"}


responsetime(15 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-02-22T10:55:00.059Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":230,"bytes_out":63,"client_ip":"10.0.13.33","client_port":44645,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":1,"num_rows":1},"path":".","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^9112177956226558331|1487760900058\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select count(1) from  t_bms_push_backup \n    where \n     \n    is_delete=0\n    and status\u002616\u003c\u003e16\n    and status\u00268\u003c\u003e8 \n    and level=2","responsetime":15,"server":"","status":"OK","type":"mysql"}


responsetime(13 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-02-22T10:55:00.078Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":231,"bytes_out":63,"client_ip":"10.0.13.33","client_port":47175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":1,"num_rows":1},"path":".","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^-1815965779687898817|1487760900078\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select count(1) from  t_bms_push_backup \n    where \n     \n    is_delete=0\n    and status\u002616\u003c\u003e16\n    and status\u00268\u003c\u003e8 \n    and level=1","responsetime":13,"server":"","status":"OK","type":"mysql"}


responsetime(7 microseconds)	==>    No.<4>
----
{"@timestamp":"2017-02-22T10:55:00.018Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":272,"bytes_out":63,"client_ip":"10.0.13.33","client_port":47175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":1,"num_rows":1},"path":".","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^946622612971717496|1487760900017\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select count(1) from  t_bms_push_point \n    where \n     \n    is_delete=0\n    and type=2 \n    and status\u002616\u003c\u003e16 \n    and status\u00268\u003c\u003e8 \n    and status\u002632\u003c\u003e32 \n    and time\u003c=now()","responsetime":7,"server":"","status":"OK","type":"mysql"}


responsetime(6 microseconds)	==>    No.<5>
----
{"@timestamp":"2017-02-22T10:54:01.550Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":365,"bytes_out":1636,"client_ip":"10.0.13.28","client_port":54065,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":1},"path":"eleme_restaurant.t_restaurant","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^7120619295620167600|1487760841549\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select\n\t\t \n\t\tid, eleme_restaurant_id, restaurant_oid, store_id, store_code,\n\t\tstore_name, city_id,\n\t\tcity_name, name,is_valid, is_delete, brand_id, brand_name, create_time, create_by,\n\t\tmodify_time,\n\t\tmodify_by\n\t \n\t\tfrom t_restaurant\n\t\twhere restaurant_oid = 73506650","responsetime":6,"server":"","status":"OK","type":"mysql"}


responsetime(6 microseconds)	==>    No.<6>
----
{"@timestamp":"2017-02-22T10:54:00.750Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":364,"bytes_out":1636,"client_ip":"10.0.13.28","client_port":49326,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":1},"path":"eleme_restaurant.t_restaurant","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^-1004743063856515343|1487760840749\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select\n\t\t \n\t\tid, eleme_restaurant_id, restaurant_oid, store_id, store_code,\n\t\tstore_name, city_id,\n\t\tcity_name, name,is_valid, is_delete, brand_id, brand_name, create_time, create_by,\n\t\tmodify_time,\n\t\tmodify_by\n\t \n\t\tfrom t_restaurant\n\t\twhere restaurant_oid = 927281","responsetime":6,"server":"","status":"OK","type":"mysql"}


responsetime(6 microseconds)	==>    No.<7>
----
{"@timestamp":"2017-02-22T10:53:57.653Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":866,"bytes_out":2729,"client_ip":"10.0.13.28","client_port":49326,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":35,"num_rows":0},"path":"eleme_restaurant.tm","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^7859161504596717036|1487760837649\u0026rpcid=1.2\u0026appid=me.ele.zs.erp:E */ SELECT  \n    id, code, name, standard_name, category_level_1, category_level_2, material_group_id, short_name, unit_group_id,unit_group,\n    transport_unit_id, purchase_type, is_valid,is_ever_valid, min_unit_id,\n    min_unit, transport_unit, shelf_life,type, is_delete, create_time, create_by, modify_time,\n    modify_by, supply_chain_prop, lead_time, moq,material_length,material_width,material_height,\n    material_netweight,material_grossweight,storage_type,material_volume,storage_conversion\n   \n    FROM t_material tm\n    WHERE  EXISTS\n    (select 1 from t_store_stock si where si.store_id =2549 and  tm.id =si.material_id)\n    and\n    (instr(tm.code,'麻辣f') OR instr(tm.short_name,'麻辣f') OR instr(tm.name,'麻辣f'))\n    and tm.is_delete=0\n     limit 0,20","responsetime":6,"server":"","status":"OK","type":"mysql"}


responsetime(6 microseconds)	==>    No.<8>
----
{"@timestamp":"2017-02-22T10:54:02.810Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":364,"bytes_out":1664,"client_ip":"10.0.13.33","client_port":47175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":1},"path":"eleme_restaurant.t_restaurant","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^5702849354046426014|1487760842809\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select\n\t\t \n\t\tid, eleme_restaurant_id, restaurant_oid, store_id, store_code,\n\t\tstore_name, city_id,\n\t\tcity_name, name,is_valid, is_delete, brand_id, brand_name, create_time, create_by,\n\t\tmodify_time,\n\t\tmodify_by\n\t \n\t\tfrom t_restaurant\n\t\twhere restaurant_oid = 2166476","responsetime":6,"server":"","status":"OK","type":"mysql"}


responsetime(6 microseconds)	==>    No.<9>
----
{"@timestamp":"2017-02-22T10:54:02.744Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":364,"bytes_out":1651,"client_ip":"10.0.13.33","client_port":47175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":1},"path":"eleme_restaurant.t_restaurant","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^8689796431643186159|1487760842743\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select\n\t\t \n\t\tid, eleme_restaurant_id, restaurant_oid, store_id, store_code,\n\t\tstore_name, city_id,\n\t\tcity_name, name,is_valid, is_delete, brand_id, brand_name, create_time, create_by,\n\t\tmodify_time,\n\t\tmodify_by\n\t \n\t\tfrom t_restaurant\n\t\twhere restaurant_oid = 1881514","responsetime":6,"server":"","status":"OK","type":"mysql"}


responsetime(6 microseconds)	==>    No.<10>
----
{"@timestamp":"2017-02-22T10:54:01.625Z","beat":{"hostname":"xg-restaurant-slave-2","name":"xg-restaurant-slave-2","version":"6.0.0-alpha1"},"bytes_in":364,"bytes_out":1658,"client_ip":"10.0.13.33","client_port":47175,"client_proc":"","client_server":"","direction":"in","ip":"10.0.10.78","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":1},"path":"eleme_restaurant.t_restaurant","port":3306,"proc":"","query":"/* E:rid=me.ele.zs.erp^^4371748204424072675|1487760841624\u0026rpcid=1.1\u0026appid=me.ele.zs.erp:E */ select\n\t\t \n\t\tid, eleme_restaurant_id, restaurant_oid, store_id, store_code,\n\t\tstore_name, city_id,\n\t\tcity_name, name,is_valid, is_delete, brand_id, brand_name, create_time, create_by,\n\t\tmodify_time,\n\t\tmodify_by\n\t \n\t\tfrom t_restaurant\n\t\twhere restaurant_oid = 2326165","responsetime":6,"server":"","status":"OK","type":"mysql"}


[root@xg-restaurant-slave-2 packageCaptureAnalysis]#
```

小结：

> 从输出结果上看，能够获取到 mysql 的 query 语句内容（其中的注释字段可以用于和 etrace 打通），能够获取到 query 的源 ip 和 port（应该对应的是 DAL 地址）；

### 监控输出

- cpu

![运行 packetbeat 分析 mysql slave 时的 cpu 资源使用情况-1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20slave%20%E6%97%B6%E7%9A%84%20cpu%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-1.png)

![运行 packetbeat 分析 mysql slave 时的 cpu 资源使用情况-2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20slave%20%E6%97%B6%E7%9A%84%20cpu%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-2.png)

- disk

![运行 packetbeat 分析 mysql slave 时的 disk 资源使用情况-1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20slave%20%E6%97%B6%E7%9A%84%20disk%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-1.png)

![运行 packetbeat 分析 mysql slave 时的 disk 资源使用情况-2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20slave%20%E6%97%B6%E7%9A%84%20disk%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-2.png)



### 测试结果（针对 master mysql）


测试命令（抓去 3306 上的 master mysql 通信 8 分钟）

```
[root@xg-breakfast-master-1 packageCaptureAnalysis]# LD_LIBRARY_PATH=. ./packageCaptureAnalysis -c ./packetbeat.yml
^C
[root@xg-breakfast-master-1 packageCaptureAnalysis]#
```

top 输出（32 核，pb 运行平均占用 35% 左右）

```
top - 15:57:46 up 311 days,  1:43,  2 users,  load average: 0.67, 1.30, 1.39
Tasks: 537 total,   1 running, 536 sleeping,   0 stopped,   0 zombie
Cpu(s):  2.3%us,  0.3%sy,  0.0%ni, 97.4%id,  0.0%wa,  0.0%hi,  0.0%si,  0.0%st
Mem:  132042812k total, 131302132k used,   740680k free,   419752k buffers
Swap: 16383996k total,        0k used, 16383996k free, 37583944k cached

   PID USER      PR  NI  VIRT  RES  SHR S %CPU %MEM    TIME+  COMMAND
116392 root      20   0 1806m  83m 9152 S 41.5  0.1   1:03.66 packageCaptureA
102241 mysql     20   0 97.2g  82g 7100 S 38.8 65.6  60173:25 mysqld
 45417 root      20   0 33088  21m 3576 S  3.3  0.0 780:47.62 esm-agent
 58785 root      20   0     0    0    0 S  1.0  0.0   3899:03 shn_comp_wqa
 58787 root      20   0     0    0    0 S  0.3  0.0 124:32.66 shn_wqa
 59167 root      20   0     0    0    0 S  0.3  0.0   1567:36 jbd2/dfa-8
116278 root      20   0 15404 1596  924 R  0.3  0.0   0:00.58 top
     1 root      20   0 19364 1296  976 S  0.0  0.0   1:06.83 init
     2 root      20   0     0    0    0 S  0.0  0.0   0:00.00 kthreadd
     3 root      RT   0     0    0    0 S  0.0  0.0   1:26.06 migration/0
     4 root      20   0     0    0    0 S  0.0  0.0   2:06.43 ksoftirqd/0
     5 root      RT   0     0    0    0 S  0.0  0.0   0:00.00 stopper/0
```

pb 运行大约 8 分钟，保存到文件中的分析结果占用大约 710M 左右 ；

```
[root@xg-breakfast-master-1 packageCaptureAnalysis]# ll logs/
total 727848
-rw-r--r-- 1 root root  11269081 Feb 23 16:03 packetbeat
-rw-r--r-- 1 root root 104857674 Feb 23 16:02 packetbeat.1
-rw-r--r-- 1 root root 104857671 Feb 23 16:01 packetbeat.2
-rw-r--r-- 1 root root 104857633 Feb 23 16:00 packetbeat.3
-rw-r--r-- 1 root root 104857741 Feb 23 15:59 packetbeat.4
-rw-r--r-- 1 root root 104858429 Feb 23 15:58 packetbeat.5
-rw-r--r-- 1 root root 104858246 Feb 23 15:57 packetbeat.6
-rw-r--r-- 1 root root 104857785 Feb 23 15:56 packetbeat.7
-rw-r--r-- 1 root root      8162 Feb 23 16:03 packetbeat.8
-rw-r--r-- 1 root root       396 Feb 23 15:55 packetbeat.9
[root@xg-breakfast-master-1 packageCaptureAnalysis]#
[root@xg-breakfast-master-1 packageCaptureAnalysis]# du -shx logs/
711M	logs/
[root@xg-breakfast-master-1 packageCaptureAnalysis]#
```

pb 输出的统计结果

```
2017-02-23T16:03:03+08:00 INFO Input finish. Processed 2731373 packets. Have a nice day!
2017-02-23T16:03:03+08:00 INFO Total non-zero values:  libbeat.publisher.published_events=902136 mysql.unmatched_responses=406316 tcp.dropped_because_of_gaps=130 mysql.unmatched_requests=3391
2017-02-23T16:03:03+08:00 INFO Uptime: 7m54.66815743s
```

基于 python 脚本进行 topN 分析（耗费大约 26 秒）

```
[root@xg-breakfast-master-1 packageCaptureAnalysis]# time python redis_analysis.py -p logs -f packetbeat,packetbeat.1,packetbeat.2,packetbeat.3,packetbeat.4,packetbeat.5,packetbeat.6,packetbeat.7 -t 10


@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

total transactions : 902136
total failure nums : 1
failure rate       : 0.000111%

@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@

responsetime(716 microseconds)	==>    No.<1>
----
{"@timestamp":"2017-02-23T08:02:00.629Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":1122,"bytes_out":27873,"client_ip":"10.0.47.27","client_port":61565,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":8,"num_rows":962},"path":"eleme_breakfast.to1","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^1220906522393322527|1487836920629\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT\n        to1.tradeAreaId as tradeAreaId,\n        to1.buildingId as buildingId,\n        to1.takeaway_id as takeawayId,\n        1 as granularity,\n        0 as businessType,\n        to1.delivery_type as deliveryType,\n        t.order_dish_type as dishType,\n        COUNT(DISTINCT(t.user_id))AS newUserNum\n        FROM t_ord_order_item t\n        JOIN t_ord_order_process to1 ON to1.orderId=t.orderId\n        WHERE   to1.createTime \u003e= '2017-02-14 00:00:00'  \n        AND    to1.createTime \u003c '2017-02-15 00:00:00'  \n        AND to1.status IN(20, 25, 30)\n        AND\n         \n            (to1.order_type = 1 or to1.order_type=4) AND to1.biz_role=0\n         \n         \n         \n         \n         \n        AND t.orderItemId in (select MIN(to2.orderItemId) from t_ord_order_item to2 JOIN t_order_ext oe on to2.orderId = oe.order_id  WHERE oe.order_id=to1.orderId  AND oe.is_first=1)\n        group by  to1.tradeAreaId,to1.buildingId,to1.takeaway_id,to1.delivery_type,t.order_dish_type\n        limit 0,1500","responsetime":716,"server":"","status":"OK","type":"mysql"}


responsetime(711 microseconds)	==>    No.<2>
----
{"@timestamp":"2017-02-23T08:01:44.208Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":1122,"bytes_out":27873,"client_ip":"10.0.47.27","client_port":61565,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":8,"num_rows":962},"path":"eleme_breakfast.to1","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^5315853952453105678|1487836904137\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT\n        to1.tradeAreaId as tradeAreaId,\n        to1.buildingId as buildingId,\n        to1.takeaway_id as takeawayId,\n        0 as granularity,\n        0 as businessType,\n        to1.delivery_type as deliveryType,\n        t.order_dish_type as dishType,\n        COUNT(DISTINCT(t.user_id))AS newUserNum\n        FROM t_ord_order_item t\n        JOIN t_ord_order_process to1 ON to1.orderId=t.orderId\n        WHERE   to1.createTime \u003e= '2017-02-14 00:00:00'  \n        AND    to1.createTime \u003c '2017-02-15 00:00:00'  \n        AND to1.status IN(20, 25, 30)\n        AND\n         \n            (to1.order_type = 1 or to1.order_type=4) AND to1.biz_role=0\n         \n         \n         \n         \n         \n        AND t.orderItemId in (select MIN(to2.orderItemId) from t_ord_order_item to2 JOIN t_order_ext oe on to2.orderId = oe.order_id  WHERE oe.order_id=to1.orderId  AND oe.is_first=1)\n        group by  to1.tradeAreaId,to1.buildingId,to1.takeaway_id,to1.delivery_type,t.order_dish_type\n        limit 0,1500","responsetime":711,"server":"","status":"OK","type":"mysql"}


responsetime(707 microseconds)	==>    No.<3>
----
{"@timestamp":"2017-02-23T08:01:47.647Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":1123,"bytes_out":27873,"client_ip":"10.0.47.27","client_port":61565,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":8,"num_rows":962},"path":"eleme_breakfast.to1","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^-7292380918749661957|1487836907647\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT\n        to1.tradeAreaId as tradeAreaId,\n        to1.buildingId as buildingId,\n        to1.takeaway_id as takeawayId,\n        1 as granularity,\n        0 as businessType,\n        to1.delivery_type as deliveryType,\n        t.order_dish_type as dishType,\n        COUNT(DISTINCT(t.user_id))AS newUserNum\n        FROM t_ord_order_item t\n        JOIN t_ord_order_process to1 ON to1.orderId=t.orderId\n        WHERE   to1.createTime \u003e= '2017-02-14 00:00:00'  \n        AND    to1.createTime \u003c '2017-02-15 00:00:00'  \n        AND to1.status IN(20, 25, 30)\n        AND\n         \n            (to1.order_type = 1 or to1.order_type=4) AND to1.biz_role=0\n         \n         \n         \n         \n         \n        AND t.orderItemId in (select MIN(to2.orderItemId) from t_ord_order_item to2 JOIN t_order_ext oe on to2.orderId = oe.order_id  WHERE oe.order_id=to1.orderId  AND oe.is_first=1)\n        group by  to1.tradeAreaId,to1.buildingId,to1.takeaway_id,to1.delivery_type,t.order_dish_type\n        limit 0,1500","responsetime":707,"server":"","status":"OK","type":"mysql"}


responsetime(706 microseconds)	==>    No.<4>
----
{"@timestamp":"2017-02-23T08:01:51.058Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":1121,"bytes_out":27873,"client_ip":"10.0.47.27","client_port":61565,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":8,"num_rows":962},"path":"eleme_breakfast.to1","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^921811928358534775|1487836911058\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT\n        to1.tradeAreaId as tradeAreaId,\n        to1.buildingId as buildingId,\n        to1.takeaway_id as takeawayId,\n        0 as granularity,\n        0 as businessType,\n        to1.delivery_type as deliveryType,\n        t.order_dish_type as dishType,\n        COUNT(DISTINCT(t.user_id))AS newUserNum\n        FROM t_ord_order_item t\n        JOIN t_ord_order_process to1 ON to1.orderId=t.orderId\n        WHERE   to1.createTime \u003e= '2017-02-14 00:00:00'  \n        AND    to1.createTime \u003c '2017-02-15 00:00:00'  \n        AND to1.status IN(20, 25, 30)\n        AND\n         \n            (to1.order_type = 1 or to1.order_type=4) AND to1.biz_role=0\n         \n         \n         \n         \n         \n        AND t.orderItemId in (select MIN(to2.orderItemId) from t_ord_order_item to2 JOIN t_order_ext oe on to2.orderId = oe.order_id  WHERE oe.order_id=to1.orderId  AND oe.is_first=1)\n        group by  to1.tradeAreaId,to1.buildingId,to1.takeaway_id,to1.delivery_type,t.order_dish_type\n        limit 0,1500","responsetime":706,"server":"","status":"OK","type":"mysql"}


responsetime(702 microseconds)	==>    No.<5>
----
{"@timestamp":"2017-02-23T08:01:54.306Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":1123,"bytes_out":27873,"client_ip":"10.0.47.27","client_port":61565,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":8,"num_rows":962},"path":"eleme_breakfast.to1","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^-6158599016008372488|1487836914306\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT\n        to1.tradeAreaId as tradeAreaId,\n        to1.buildingId as buildingId,\n        to1.takeaway_id as takeawayId,\n        1 as granularity,\n        0 as businessType,\n        to1.delivery_type as deliveryType,\n        t.order_dish_type as dishType,\n        COUNT(DISTINCT(t.user_id))AS newUserNum\n        FROM t_ord_order_item t\n        JOIN t_ord_order_process to1 ON to1.orderId=t.orderId\n        WHERE   to1.createTime \u003e= '2017-02-14 00:00:00'  \n        AND    to1.createTime \u003c '2017-02-15 00:00:00'  \n        AND to1.status IN(20, 25, 30)\n        AND\n         \n            (to1.order_type = 1 or to1.order_type=4) AND to1.biz_role=0\n         \n         \n         \n         \n         \n        AND t.orderItemId in (select MIN(to2.orderItemId) from t_ord_order_item to2 JOIN t_order_ext oe on to2.orderId = oe.order_id  WHERE oe.order_id=to1.orderId  AND oe.is_first=1)\n        group by  to1.tradeAreaId,to1.buildingId,to1.takeaway_id,to1.delivery_type,t.order_dish_type\n        limit 0,1500","responsetime":702,"server":"","status":"OK","type":"mysql"}


responsetime(698 microseconds)	==>    No.<6>
----
{"@timestamp":"2017-02-23T08:01:57.515Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":1123,"bytes_out":27873,"client_ip":"10.0.47.27","client_port":61565,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":8,"num_rows":962},"path":"eleme_breakfast.to1","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^-2254936189944000241|1487836917515\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT\n        to1.tradeAreaId as tradeAreaId,\n        to1.buildingId as buildingId,\n        to1.takeaway_id as takeawayId,\n        0 as granularity,\n        0 as businessType,\n        to1.delivery_type as deliveryType,\n        t.order_dish_type as dishType,\n        COUNT(DISTINCT(t.user_id))AS newUserNum\n        FROM t_ord_order_item t\n        JOIN t_ord_order_process to1 ON to1.orderId=t.orderId\n        WHERE   to1.createTime \u003e= '2017-02-14 00:00:00'  \n        AND    to1.createTime \u003c '2017-02-15 00:00:00'  \n        AND to1.status IN(20, 25, 30)\n        AND\n         \n            (to1.order_type = 1 or to1.order_type=4) AND to1.biz_role=0\n         \n         \n         \n         \n         \n        AND t.orderItemId in (select MIN(to2.orderItemId) from t_ord_order_item to2 JOIN t_order_ext oe on to2.orderId = oe.order_id  WHERE oe.order_id=to1.orderId  AND oe.is_first=1)\n        group by  to1.tradeAreaId,to1.buildingId,to1.takeaway_id,to1.delivery_type,t.order_dish_type\n        limit 0,1500","responsetime":698,"server":"","status":"OK","type":"mysql"}


responsetime(634 microseconds)	==>    No.<7>
----
{"@timestamp":"2017-02-23T08:00:31.128Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":2814,"bytes_out":13502,"client_ip":"10.0.47.28","client_port":38226,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":14,"num_rows":200},"path":".t","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^-4099252426591907081|1487836831128\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT cityId,\n        cityName,\n        tradeAreaId,\n        buildingId,\n        takeawayId,\n        1 as granularity,\n        0 as businessType,\n        deliveryType,\n        dishType,\n         \n         \n             \n             \n                sum(deliveryNumOld) deliveryNumOld,\n                sum(retailAmountOld) retailAmountOld,\n                sum(purchaseAmountOld) purchaseAmountOld,\n                sum(saleAmountOld) saleAmountOld,\n                sum(paidAmountOld) paidAmountOld\n             \n         \n        FROM (\n            SELECT cityId,cityName,tradeAreaId,buildingId,takeawayId,deliveryType,dishType,\n             \n             \n                 \n                 \n                    COUNT(DISTINCT takeawayDate) deliveryNumOld,\n                    sum(retailAmountOld) retailAmountOld,\n                    sum(purchaseAmountOld) purchaseAmountOld,\n                    sum(saleAmountOld) saleAmountOld,\n                    sum(paidAmountOld) paidAmountOld\n                 \n             \n            FROM (\n                SELECT t.city_id as cityId,\n                t.city_name as cityName,\n                t.trade_area_id as tradeAreaId,\n                t.building_id as buildingId,\n                t.takeaway_id as takeawayId,\n                t1.delivery_type as deliveryType,\n                t.dish_type as dishType,\n                t.takeaway_date as takeawayDate,\n                t1.elemeOrderId as elemeOrderId,\n                 \n                 \n                     \n                     \n                        t.price*t.quantity as retailAmountOld,\n                        t.purchase_price*t.quantity as purchaseAmountOld,\n                        t.selling_amount as saleAmountOld,\n                        t.selling_amount-t.voucher_amount as paidAmountOld\n                     \n                 \n                FROM t_delivery_order_item t\n                JOIN t_ord_order_process t1 ON t.order_id = t1.orderId\n                WHERE t.takeaway_date = '2017-02-23 00:00:00' and t.takeaway_status in (20,30) AND\n                 \n                    (t1.order_type = 1 or t1.order_type=4) AND t1.biz_role=0\n                 \n                 \n                 \n                 \n                 \n                 \n                     \n                     \n                        AND NOT EXISTS (select 1 from t_order_ext oe WHERE oe.order_id=t1.orderId  AND oe.is_first=1)\n                     \n                 \n            )t  GROUP BY cityId,cityName,tradeAreaId,buildingId,takeawayId,deliveryType,dishType,elemeOrderId\n        ) t\n        GROUP BY cityId,cityName,tradeAreaId,buildingId,takeawayId,deliveryType,dishType\n        limit 1500,1500","responsetime":634,"server":"","status":"OK","type":"mysql"}


responsetime(606 microseconds)	==>    No.<8>
----
{"@timestamp":"2017-02-23T08:00:42.276Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":469,"bytes_out":2334,"client_ip":"10.0.13.95","client_port":38975,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"SELECT","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":7},"path":"ecs.order_assign","port":3306,"proc":"","query":"select  \n    id, order_pool_id, order_id, cs_id, rst_id, user_id, status, order_type, is_valid,\n    assign_reason, invalid_reason, active_at, accept_at, handle_time, create_time, create_by,\n    update_time\n   \n    from order_assign\n    where cs_id = 2011\n    and is_valid = 1\n     \n      and status = 0\n     \n     \n      and order_type = 1\n     \n     \n     \n     \n    and 'bind_master' = 'bind_master'\n     \n        order by create_time asc\n       \n    limit 0, 15","responsetime":606,"server":"","status":"OK","type":"mysql"}


responsetime(593 microseconds)	==>    No.<9>
----
{"@timestamp":"2017-02-23T07:57:47.206Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":468,"bytes_out":2327,"client_ip":"10.0.13.30","client_port":42406,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"SELECT","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":17,"num_rows":7},"path":"ecs.order_assign","port":3306,"proc":"","query":"select  \n    id, order_pool_id, order_id, cs_id, rst_id, user_id, status, order_type, is_valid,\n    assign_reason, invalid_reason, active_at, accept_at, handle_time, create_time, create_by,\n    update_time\n   \n    from order_assign\n    where cs_id = 474\n    and is_valid = 1\n     \n      and status = 0\n     \n     \n      and order_type = 1\n     \n     \n     \n     \n    and 'bind_master' = 'bind_master'\n     \n        order by create_time asc\n       \n    limit 0, 15","responsetime":593,"server":"","status":"OK","type":"mysql"}


responsetime(590 microseconds)	==>    No.<10>
----
{"@timestamp":"2017-02-23T08:00:26.587Z","beat":{"hostname":"xg-breakfast-master-1","name":"xg-breakfast-master-1","version":"6.0.0-alpha1"},"bytes_in":2810,"bytes_out":94892,"client_ip":"10.0.47.28","client_port":38226,"client_proc":"","client_server":"","direction":"in","ip":"10.0.27.244","method":"/*","mysql":{"affected_rows":0,"error_code":0,"error_message":"","insert_id":0,"iserror":false,"num_fields":14,"num_rows":1500},"path":".t","port":3306,"proc":"","query":"/* E:rid=me.ele.breakfast.backend^^6683072567797086865|1487836826587\u0026rpcid=1.1\u0026appid=me.ele.breakfast.backend:E */ SELECT cityId,\n        cityName,\n        tradeAreaId,\n        buildingId,\n        takeawayId,\n        1 as granularity,\n        0 as businessType,\n        deliveryType,\n        dishType,\n         \n         \n             \n             \n                sum(deliveryNumOld) deliveryNumOld,\n                sum(retailAmountOld) retailAmountOld,\n                sum(purchaseAmountOld) purchaseAmountOld,\n                sum(saleAmountOld) saleAmountOld,\n                sum(paidAmountOld) paidAmountOld\n             \n         \n        FROM (\n            SELECT cityId,cityName,tradeAreaId,buildingId,takeawayId,deliveryType,dishType,\n             \n             \n                 \n                 \n                    COUNT(DISTINCT takeawayDate) deliveryNumOld,\n                    sum(retailAmountOld) retailAmountOld,\n                    sum(purchaseAmountOld) purchaseAmountOld,\n                    sum(saleAmountOld) saleAmountOld,\n                    sum(paidAmountOld) paidAmountOld\n                 \n             \n            FROM (\n                SELECT t.city_id as cityId,\n                t.city_name as cityName,\n                t.trade_area_id as tradeAreaId,\n                t.building_id as buildingId,\n                t.takeaway_id as takeawayId,\n                t1.delivery_type as deliveryType,\n                t.dish_type as dishType,\n                t.takeaway_date as takeawayDate,\n                t1.elemeOrderId as elemeOrderId,\n                 \n                 \n                     \n                     \n                        t.price*t.quantity as retailAmountOld,\n                        t.purchase_price*t.quantity as purchaseAmountOld,\n                        t.selling_amount as saleAmountOld,\n                        t.selling_amount-t.voucher_amount as paidAmountOld\n                     \n                 \n                FROM t_delivery_order_item t\n                JOIN t_ord_order_process t1 ON t.order_id = t1.orderId\n                WHERE t.takeaway_date = '2017-02-23 00:00:00' and t.takeaway_status in (20,30) AND\n                 \n                    (t1.order_type = 1 or t1.order_type=4) AND t1.biz_role=0\n                 \n                 \n                 \n                 \n                 \n                 \n                     \n                     \n                        AND NOT EXISTS (select 1 from t_order_ext oe WHERE oe.order_id=t1.orderId  AND oe.is_first=1)\n                     \n                 \n            )t  GROUP BY cityId,cityName,tradeAreaId,buildingId,takeawayId,deliveryType,dishType,elemeOrderId\n        ) t\n        GROUP BY cityId,cityName,tradeAreaId,buildingId,takeawayId,deliveryType,dishType\n        limit 0,1500","responsetime":590,"server":"","status":"OK","type":"mysql"}



real	0m26.254s
user	0m25.738s
sys	0m0.514s
[root@xg-breakfast-master-1 packageCaptureAnalysis]#
```

小结：

> 从输出结果上看，能够获取到 mysql 的 query 语句内容（其中的注释字段可以用于和 etrace 打通），能够获取到 query 的源 ip 和 port（应该对应的是 DAL 地址）；

### 监控输出

- cpu

![运行 packetbeat 分析 mysql master 时的 cpu 资源使用情况-1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20master%20%E6%97%B6%E7%9A%84%20cpu%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-1.png)

![运行 packetbeat 分析 mysql master 时的 cpu 资源使用情况-2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20master%20%E6%97%B6%E7%9A%84%20cpu%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-2.png)

- disk

![运行 packetbeat 分析 mysql master 时的 disk 资源使用情况-1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20master%20%E6%97%B6%E7%9A%84%20disk%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-1.png)

![运行 packetbeat 分析 mysql master 时的 disk 资源使用情况-2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%BF%90%E8%A1%8C%20packetbeat%20%E5%88%86%E6%9E%90%20mysql%20master%20%E6%97%B6%E7%9A%84%20disk%20%E8%B5%84%E6%BA%90%E4%BD%BF%E7%94%A8%E6%83%85%E5%86%B5-2.png)