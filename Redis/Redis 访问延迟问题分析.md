# Redis 延迟问题分析

## 怀疑和 corvus 或 redis 自身处理机制有关


> 怀疑 corvus 中是否存在 40ms 延迟相关的代码或处理逻辑；
> 
> 因为看到两种情况：
> 
> - corvus 收到 redis 的数据包后，40ms 后做出回应；
> - redis 收到 corvus 的数据包后，40ms 后做出回应；
>
> 然后，实际并不存在这类逻辑；


- 80ms 延迟日志

![80ms 延迟日志](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/corvus_log_80ms.png "80ms 延迟日志")


完整日志：[这里](https://eleme.slack.com/files/guangxing.huang/F6EDRRHUH/-.txt)

- corvus 侧观察到的 80ms 延迟

![corvus_pkg_80ms_1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/corvus_pkg_80ms_1.png "corvus_pkg_80ms_1")

![corvus_pkg_80ms_2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/corvus_pkg_80ms_2.png "corvus_pkg_80ms_2")

- redis 侧观察到的 80ms 延迟

![redis_pkg_80ms_1](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis_pkg_80ms_1.png "redis_pkg_80ms_1")

![redis_pkg_80ms_2](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis_pkg_80ms_2.png "redis_pkg_80ms_2")


----------

- corvus 侧观察到的 40ms 延迟

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/corvus_pkg_40ms_2.png)

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/corvus_pkg_40ms_1.png)

- redis 侧观察到的 40ms 延迟

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis_pkg_40ms_1.png)

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis_pkg_40ms_2.png)

## 怀疑和 pipeline 使用有关

> 怀疑：corvus 某些情况下按照 pipeline 方式将多条请求批量发送到 redis ，导致延迟变大；
> 实际：corvus 对于所有命令都是以 pipeline 模式处理的；对于客户端来说，是否使用 pipeline 则完全由其自己控制；corvus 不会使用事务，因为会导致问题太复杂；
>
> 虽然有些情况中，大延迟确认发生在 pipeline 交互过程，但另外一些 pipeline 交互过程又表现的很正常；

观察发现：每次出现延迟时，corvus 都会连续发送多个 request ，而后 redis 才进行 response 回复；

经常看到 corvus 发送 3~4 个 request 后，redis 才进行 response 回复；

行为大致如下：

```
corvus => (request) => redis
(~40ms)
corvus <= (TCP Ack) <= redis
(~35ms or ~15ms)
corvus => (request) => redis
(~25ms or ~35ms)
corvus => (request) => redis
(~2ms)
corvus <= (TCP Ack) <= redis
(~80ms or ~1ms)
corvus <= (response) <= redis
corvus <= (response) <= redis
corvus <= (response) <= redis
...
```


## 怀疑和网络设备有关

> 基于 redis 侧的 **Send-Q** 经常会有积压，而 Corvus 侧的 **Recv-Q** 中却没有任何积压的情况，怀疑 redis 到 corvus 路径上的网络设备存在性能问题；
>
> 经确认，中间网络设备有开启概率性丢弃 ping 包功能，因此基于 ping 的 mtr 输出报告（默认）可能并不能反映出真实问题；目前 mtr 同样支持基于 TCP 的探测；
>
> 结论：基于 mtr 的 TCP 模式进行探测后，输出数据并不能表明中间网络设备的问题（也可能我测试的不到位）；


![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis_latency_Send-Q_Recv-Q.jpeg)

从截图中可以看到：

- 针对一条 TCP 连接，redis 侧的 **Send-Q** 经常会有积压，而 Corvus 侧的 **Recv-Q** 中却没有任何积压；
- 结合 ping 和 mtr 输出，似乎 redis 到 corvus 路径上最后一跳存在一些问题；

> 以下内容取自：[Use of Recv-Q and Send-Q](https://stackoverflow.com/questions/36466744/use-of-recv-q-and-send-q)

Send-Q

- **Established**: The count of bytes not acknowledged by the remote host.
- **Listening**: Since Kernel 2.6.18 this column contains the maximum size of the syn backlog.
- a non-zero Send-Q, if the other side TCP implementation at the OS level have been stucked and stopped ACKnowleding the data.
- the Send-Q scenario depicted above may also be a sending side issue if the Linux TCP implementation was misbehaving and continued to send data after the TCP window went down to 0: the receiving side then has no more room for this data, so does not ACKnowledge.
- Send-Q issue may be cause not because of the receiver, but by some routing issue somewhere between the sender and the receiver. Some packets are "on the fly" between the 2 hosts, but not ACKnowledge yet. 
- Send-Q issue to be caused most of the time by some routing issue/network poor performances between the sending and receiving side. 


Recv-Q

- **Established**: The count of bytes not copied by the user program connected to this socket.
- **Listening**: Since Kernel 2.6.18 this column contains the current syn backlog.
- an increasing Recv-Q, up to some roof value, where the other side stop sending data because the window get down to 0, since the application does not read the data available on its socket, and these data stay buffered in the TCP implementation in the OS, not going to the stucked application.
- Recv-Q issue is definitly on a host: packets received, ACKnowledged, but not read from the application yet.


结论：

- If you have this stuck to 0, this just mean that your applications, on both side of the connection, and the network between them, are doing OK. Actual instant values may be different from 0, but in such a transient, fugitive manner that you don't get a chance to actually observe it.
- The "on the fly" state of packets should never be forgotten:
    - The packet may be on the network between the sender and the receiver,
    - (or received but ACK not send yet, see above)
    - or the ACK may be on the network between the receiver and the sender.
- It takes a RTT (round time trip) for a packet to be send and then ACKed.

> 整理了一篇[如何使用 MTR](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/MTR.md) 的文章；

## 怀疑和 Buddy memory allocation 相关

> 怀疑 redis 的延迟是由于伙伴系统内存使用导致；
>
> 需要了解 redis 的内存使用分布情况，并对照 buddyinfo 内容分析是否由于内存使用问题导致延迟；

查看 buddyinfo 信息

```
[root@wg-public-rediscluster-119: ~]# cat /proc/buddyinfo
Node 0, zone      DMA      0      1      1      1      1      1      1      0      1      1      3
Node 0, zone    DMA32    481    291    216   1121   1004    736    592    484    409      0      0
Node 0, zone   Normal   8394   7766   2167   1292   8641  22518  19533  16089  13125      0      0
Node 1, zone   Normal 109715  64348  43552  33212  24212  16771  11738   9263   7718   2389   1859
[root@wg-public-rediscluster-119: ~]#
```

从如下输出中可以对比得到上面信息的具体含义：

```
[root@wg-public-rediscluster-119: ~]# echo m > /proc/sysrq-trigger
[root@wg-public-rediscluster-119: ~]# dmesg
...
[11464077.355400] SysRq : Show Memory
[11464077.355439] Mem-Info:
[11464077.355454] active_anon:6758464 inactive_anon:69505 isolated_anon:0
 active_file:600715 inactive_file:252796 isolated_file:0
 unevictable:0 dirty:14 writeback:0 unstable:0
 slab_reclaimable:72415 slab_unreclaimable:44153
 mapped:16344 shmem:139560 pagetables:15797 bounce:0
 free:16328239 free_pcp:10107 free_cma:0
[11464077.355459] Node 0 DMA free:15864kB min:12kB low:12kB high:16kB active_anon:0kB inactive_anon:0kB active_file:0kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:15948kB managed:15864kB mlocked:0kB dirty:0kB writeback:0kB mapped:0kB shmem:0kB slab_reclaimable:0kB slab_unreclaimable:0kB kernel_stack:0kB pagetables:0kB unstable:0kB bounce:0kB free_pcp:0kB local_pcp:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11464077.355465] lowmem_reserve[]: 0 1662 47995 47995
[11464077.355469] Node 0 DMA32 free:1018864kB min:1552kB low:1940kB high:2328kB active_anon:481848kB inactive_anon:6508kB active_file:62540kB inactive_file:30652kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:1952512kB managed:1702652kB mlocked:0kB dirty:8kB writeback:0kB mapped:1972kB shmem:14516kB slab_reclaimable:9132kB slab_unreclaimable:5196kB kernel_stack:352kB pagetables:2800kB unstable:0kB bounce:0kB free_pcp:11020kB local_pcp:348kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11464077.355475] lowmem_reserve[]: 0 0 46332 46332
[11464077.355478] Node 0 Normal free:31492388kB min:43312kB low:54140kB high:64968kB active_anon:11524676kB inactive_anon:262008kB active_file:1445780kB inactive_file:675912kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:48234496kB managed:47444940kB mlocked:0kB dirty:32kB writeback:0kB mapped:53164kB shmem:518556kB slab_reclaimable:180492kB slab_unreclaimable:112024kB kernel_stack:5824kB pagetables:56588kB unstable:0kB bounce:0kB free_pcp:18396kB local_pcp:364kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11464077.355484] lowmem_reserve[]: 0 0 0 0
[11464077.355487] Node 1 Normal free:32785840kB min:45228kB low:56532kB high:67840kB active_anon:15027332kB inactive_anon:9504kB active_file:894540kB inactive_file:304620kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:50331648kB managed:49541852kB mlocked:0kB dirty:16kB writeback:0kB mapped:10240kB shmem:25168kB slab_reclaimable:100036kB slab_unreclaimable:59392kB kernel_stack:5136kB pagetables:3800kB unstable:0kB bounce:0kB free_pcp:11012kB local_pcp:92kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11464077.355493] lowmem_reserve[]: 0 0 0 0
[11464077.355495] Node 0 DMA: 0*4kB 1*8kB (U) 1*16kB (U) 1*32kB (U) 1*64kB (U) 1*128kB (U) 1*256kB (U) 0*512kB 1*1024kB (U) 1*2048kB (M) 3*4096kB (M) = 15864kB
[11464077.355507] Node 0 DMA32: 998*4kB (UEM) 363*8kB (UEM) 560*16kB (UEM) 1162*32kB (UEM) 963*64kB (UEM) 718*128kB (UEM) 583*256kB (EM) 477*512kB (UM) 409*1024kB (UEM) 0*2048kB 0*4096kB = 1018864kB
[11464077.355519] Node 0 Normal: 12617*4kB (UEM) 5822*8kB (UEM) 4573*16kB (UEM) 13057*32kB (UEM) 22675*64kB (UEM) 22261*128kB (UM) 19401*256kB (UEM) 16034*512kB (UEM) 13113*1024kB (EM) 0*2048kB 0*4096kB = 31492420kB
[11464077.355530] Node 1 Normal: 97202*4kB (UM) 57273*8kB (UEM) 39870*16kB (UEM) 31610*32kB (UEM) 24404*64kB (UEM) 18200*128kB (UEM) 13599*256kB (UM) 10903*512kB (UEM) 8946*1024kB (UM) 289*2048kB (UM) 1851*4096kB (UM) = 32785840kB
[11464077.355544] Node 0 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=1048576kB
[11464077.355546] Node 0 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=2048kB
[11464077.355547] Node 1 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=1048576kB
[11464077.355549] Node 1 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=2048kB
[11464077.355550] 993068 total pagecache pages
[11464077.355552] 0 pages in swap cache
[11464077.355554] Swap cache stats: add 0, delete 0, find 0/0
[11464077.355555] Free swap  = 0kB
[11464077.355556] Total swap = 0kB
[11464077.355558] 25133651 pages RAM
[11464077.355559] 0 pages HighMem/MovableOnly
[11464077.355560] 457324 pages reserved
[11811682.679933] SysRq : Show Memory
[11811682.679970] Mem-Info:
[11811682.679985] active_anon:6619350 inactive_anon:71550 isolated_anon:0
 active_file:599958 inactive_file:252639 isolated_file:0
 unevictable:0 dirty:11 writeback:0 unstable:0
 slab_reclaimable:72439 slab_unreclaimable:43504
 mapped:18277 shmem:141600 pagetables:15794 bounce:0
 free:16467322 free_pcp:9694 free_cma:0
[11811682.679990] Node 0 DMA free:15864kB min:12kB low:12kB high:16kB active_anon:0kB inactive_anon:0kB active_file:0kB inactive_file:0kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:15948kB managed:15864kB mlocked:0kB dirty:0kB writeback:0kB mapped:0kB shmem:0kB slab_reclaimable:0kB slab_unreclaimable:0kB kernel_stack:0kB pagetables:0kB unstable:0kB bounce:0kB free_pcp:0kB local_pcp:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11811682.679996] lowmem_reserve[]: 0 1662 47995 47995
[11811682.680000] Node 0 DMA32 free:1020532kB min:1552kB low:1940kB high:2328kB active_anon:485320kB inactive_anon:6508kB active_file:61980kB inactive_file:30644kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:1952512kB managed:1702652kB mlocked:0kB dirty:8kB writeback:0kB mapped:2004kB shmem:14516kB slab_reclaimable:9132kB slab_unreclaimable:5084kB kernel_stack:288kB pagetables:2800kB unstable:0kB bounce:0kB free_pcp:9400kB local_pcp:492kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11811682.680006] lowmem_reserve[]: 0 0 46332 46332
[11811682.680009] Node 0 Normal free:30265844kB min:43312kB low:54140kB high:64968kB active_anon:12743300kB inactive_anon:270196kB active_file:1442948kB inactive_file:673972kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:48234496kB managed:47444940kB mlocked:0kB dirty:36kB writeback:0kB mapped:60776kB shmem:526740kB slab_reclaimable:180252kB slab_unreclaimable:112932kB kernel_stack:5776kB pagetables:56996kB unstable:0kB bounce:0kB free_pcp:19408kB local_pcp:48kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11811682.680015] lowmem_reserve[]: 0 0 0 0
[11811682.680018] Node 1 Normal free:34567048kB min:45228kB low:56532kB high:67840kB active_anon:13248780kB inactive_anon:9496kB active_file:894904kB inactive_file:305940kB unevictable:0kB isolated(anon):0kB isolated(file):0kB present:50331648kB managed:49541852kB mlocked:0kB dirty:0kB writeback:0kB mapped:10328kB shmem:25144kB slab_reclaimable:100372kB slab_unreclaimable:56000kB kernel_stack:5072kB pagetables:3380kB unstable:0kB bounce:0kB free_pcp:9968kB local_pcp:0kB free_cma:0kB writeback_tmp:0kB pages_scanned:0 all_unreclaimable? no
[11811682.680023] lowmem_reserve[]: 0 0 0 0
[11811682.680026] Node 0 DMA: 0*4kB 1*8kB (U) 1*16kB (U) 1*32kB (U) 1*64kB (U) 1*128kB (U) 1*256kB (U) 0*512kB 1*1024kB (U) 1*2048kB (M) 3*4096kB (M) = 15864kB
[11811682.680038] Node 0 DMA32: 627*4kB (UEM) 263*8kB (UEM) 217*16kB (UEM) 1119*32kB (UEM) 1004*64kB (UEM) 736*128kB (UM) 592*256kB (UEM) 484*512kB (UM) 409*1024kB (UEM) 0*2048kB 0*4096kB = 1020532kB
[11811682.680049] Node 0 Normal: 11530*4kB (UEM) 5664*8kB (UEM) 1211*16kB (UEM) 1298*32kB (UEM) 8641*64kB (UEM) 22518*128kB (UM) 19533*256kB (UEM) 16089*512kB (UEM) 13125*1024kB (EM) 0*2048kB 0*4096kB = 30265688kB
[11811682.680061] Node 1 Normal: 109516*4kB (UEM) 64449*8kB (UM) 43551*16kB (UEM) 33209*32kB (UEM) 24209*64kB (UEM) 16772*128kB (UEM) 11739*256kB (UM) 9262*512kB (UM) 7718*1024kB (UM) 2389*2048kB (UM) 1859*4096kB (UEM) = 34567048kB
[11811682.680074] Node 0 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=1048576kB
[11811682.680076] Node 0 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=2048kB
[11811682.680077] Node 1 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=1048576kB
[11811682.680079] Node 1 hugepages_total=0 hugepages_free=0 hugepages_surp=0 hugepages_size=2048kB
[11811682.680080] 994195 total pagecache pages
[11811682.680082] 0 pages in swap cache
[11811682.680084] Swap cache stats: add 0, delete 0, find 0/0
[11811682.680085] Free swap  = 0kB
[11811682.680086] Total swap = 0kB
[11811682.680088] 25133651 pages RAM
[11811682.680089] 0 pages HighMem/MovableOnly
[11811682.680090] 457324 pages reserved
[root@wg-public-rediscluster-119: ~]#
```

> [Buddy memory allocation](https://en.wikipedia.org/wiki/Buddy_memory_allocation)


## 怀疑和 NUMA node 内存访问相关

> 怀疑 redis 的延迟由于跨 NUMA 节点使用内存导致；
>
> 观察访问情况是否均衡；

```
[root@wg-public-rediscluster-119: ~]# lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                40
On-line CPU(s) list:   0-39
Thread(s) per core:    2
Core(s) per socket:    10
Socket(s):             2
NUMA node(s):          2
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 79
Model name:            Intel(R) Xeon(R) CPU E5-2640 v4 @ 2.40GHz
Stepping:              1
CPU MHz:               2400.000
BogoMIPS:              4799.38
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              25600K
NUMA node0 CPU(s):     0-9,20-29
NUMA node1 CPU(s):     10-19,30-39
[root@wg-public-rediscluster-119: ~]#
[root@wg-public-rediscluster-119: ~]# numastat
                           node0           node1
numa_hit              6675249920      3536826510
numa_miss                      0           60799
numa_foreign               60799               0
interleave_hit             32201           32655
local_node            6675173321      3536620723
other_node                 76599          266586
[root@wg-public-rediscluster-119: ~]#
```

## Transparent Hugepage

最初没怀疑和这个参数有关的原因在于，我们的 redis 实例默认是不开启 rdb 和 aof 的，即纯内存应用，因此理论上是不会进行 fork 操作的；而如下的 redis 官方文档特意强调了 fork 行为才导致问题，于是将我们带偏了；事实上，问题的原因正是因为 Transparent Hugepage 开启才导致；

在《[Redis latency problems troubleshooting](https://redis.io/topics/latency)》中提到

> Transparent huge pages must be disabled from your kernel. Use `echo never > /sys/kernel/mm/transparent_hugepage/enabled` to disable them, and restart your Redis process.
>
> ### Latency induced by transparent huge pages
>
> Unfortunately when a Linux kernel has transparent huge pages enabled, Redis incurs to a big latency penalty after the `fork` call is used in order to persist on disk. Huge pages are the cause of the following issue:
> 
> - `Fork` is called, two processes with shared huge pages are created.
> - In a busy instance, a few event loops runs will cause commands to target a few thousand of pages, causing the copy on write of almost the whole process memory.
> - This will result in big latency and big memory usage.
>
> Make sure to disable transparent huge pages using the following command:
> 
> `echo never > /sys/kernel/mm/transparent_hugepage/enabled`

在《[When to turn off Transparent Huge Pages for redis](https://stackoverflow.com/questions/42591511/when-to-turn-off-transparent-huge-pages-for-redis)》中提到

> The problem lies in **how THP shifts memory around to try and keep or create contiguous pages**. Some applications can tolerate this, most databases cannot and it causes intermittent performance problems, some pretty bad. This is not unique to Redis by any means.
>
> For your application, especially if it is JAVA, set up real HugePages and leave the transparent variety out of it. If you do that just make sure you alocate memory correctly for the app and redis. Though I have to say, I probably would not recommend running both the app and redis on the same instance/server/vm.

在《[Transparent Huge Pages and Alternative Memory Allocators: A Cautionary Tale](https://blog.digitalocean.com/transparent-huge-pages-and-alternative-memory-allocators/)》中提到（主要讲述内存分配器和 THP 相互作用下的内存泄漏问题）

> Despite initially looking like a leak, the problem was actually an issue between an alternative memory allocator and transparent huge pages.
>
> disabling transparent huge pages requires manually echoing settings
> ```
> echo never > /sys/kernel/mm/transparent_hugepage/enabled
> echo never > /sys/kernel/mm/transparent_hugepage/defrag
> ```


在《[Often Overlooked Linux OS Tweaks](https://blog.couchbase.com/often-overlooked-linux-os-tweaks/)》中提到

> ### Disable Transparent Huge Pages (THP)
>
> Starting in Red Hat Enterprise Linux (RHEL) version 6, so this includes CentOS 6 and 7 too, a new default method of managing huge pages was implemented in the OS. Ubuntu has this setting as well starting in 12.02, so it will need this changed as well. **THP** combines smaller memory pages into Huge Pages without the running processes knowing. The idea is to **reduce the number of lookups on TLB required** and therefor increase performance. It brings in abstraction for automatation and management of huge pages basically.  Couchbase Engineering has determined that under some conditions, **Couchbase Server can be negatively impacted by severe page allocation delays when THP is enabled**. Couchbase therefore recommends that THP be disabled on all Couchbase Server nodes
>
>
> #### Confirm if the OS settings need to be disabled
> 
> Check the status of THP by issuing the following commands:
>
> ```
> cat /sys/kernel/mm/transparent_hugepage/enabled
> cat /sys/kernel/mm/transparent_hugepage/defrag
> ```
> On some Red Hat or Red Hat variants, you might have to do this:
>
> ```
> cat /sys/kernel/mm/redhat_transparent_hugepage/enabled
> cat /sys/kernel/mm/redhat_transparent_hugepage/defrag
> ```
>
> If in one or both files, the output looks like this, you need the below procedure:
> 
> ```
> [always] madvise never
> ```
>
> Disable THP by
> 
> ```
> echo 'never' > /sys/kernel/mm/transparent_hugepage/enabled
> echo 'never' > /sys/kernel/mm/transparent_hugepage/defrag
> ```
>
> the output should be like this:
>
> ```
> always madvise [never]
> ```
>
> THP is a great feature for some things, but causes problems with applications like Couchbase. It is not alone in this. If you go search the Internet for transparent huge pages, **there are multiple documented issues from other DB and application vendors about this**. Until something has been found to work with this, it is just best to turn THP off.

