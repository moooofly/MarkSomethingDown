# How to achieve 100w Concurrent TCP connections

## 相关参考

- [Linux sysctl命令介绍](https://www.ifshow.com/linux-sysctl-command-introduced/)
- [The Road to 2 Million Websocket Connections in Phoenix](http://phoenixframework.org/blog/the-road-to-2-million-websocket-connections)
- [100万并发连接服务器笔记](http://www.blogjava.net/yongboy/category/54842.html)
- [kernel nf_conntrack: table full, dropping packet 解决办法](https://blog.yorkgu.me/2012/02/09/kernel-nf_conntrack-table-full-dropping-packet/)
- [Netfilter conntrack performance tweaking, v0.8](https://blog.yorkgu.me/wp-content/uploads/2012/02/netfilter_conntrack_perf-0.8.txt)
- [解决 nf_conntrack: table full, dropping packet 的几种思路](http://jaseywang.me/2012/08/16/%E8%A7%A3%E5%86%B3-nf_conntrack-table-full-dropping-packet-%E7%9A%84%E5%87%A0%E7%A7%8D%E6%80%9D%E8%B7%AF/)
- [通过 modprobe 彻底禁用 netfilter](http://jaseywang.me/2012/11/18/%E9%80%9A%E8%BF%87-modprobe-%E5%BD%BB%E5%BA%95%E7%A6%81%E7%94%A8-netfilter/)
- [nf_conntrack: table full, dropping packet. 终结篇](http://www.cnblogs.com/higkoo/articles/iptables_tunning_for_conntrack.html)
- [如何在Linux中启动/停止和启用/禁用FirewallD和Iptables防火墙](https://www.howtoing.com/start-stop-disable-enable-firewalld-iptables-firewall/)
- [Go语言TCP Socket编程](http://tonybai.com/2015/11/17/tcp-programming-in-golang/?url_type=39&object_type=webpage&pos=1)
- [扛住100亿次请求？我们来试一试](https://github.com/xiaojiaqi/10billionhongbaos/wiki/%E6%89%9B%E4%BD%8F100%E4%BA%BF%E6%AC%A1%E8%AF%B7%E6%B1%82%EF%BC%9F%E6%88%91%E4%BB%AC%E6%9D%A5%E8%AF%95%E4%B8%80%E8%AF%95)
- [A Million WebSockets and Go](https://medium.freecodecamp.org/million-websockets-and-go-cc58418460bb)
- [10M Concurrent Websockets](https://goroutines.com/10m)
- [如何实现支持数亿用户的长连消息系统](https://mp.weixin.qq.com/s/PrsttFAHqtUOSTxAHM6QZg)
- [10M Concurrent Websockets in Go](https://www.reddit.com/r/golang/comments/49q96w/10m_concurrent_websockets_in_go/)
- [10M Concurrent Websockets](https://news.ycombinator.com/item?id=11320023)
- [Golang适合高并发场景的原因分析](http://www.cnblogs.com/ghj1976/p/3763866.html)
- [A signaling system for 10M+ concurrent connections](https://blog.greta.io/a-signaling-system-for-10m-concurrent-connections-10d327fd6837)
- [聊一聊goroutine stack](https://zhuanlan.zhihu.com/p/28409657)



## 系统参数

- **`/proc/sys/fs/file-max`**

This file defines a **system-wide** limit on the **number of open files** for **all processes**. System calls that fail when encountering this limit fail with the error ENFILE. (See also `setrlimit(2)`, which can be used by a process to set the per-process limit, `RLIMIT_NOFILE`, on the number of files it may open.)  If you get lots of error messages in the kernel log about running out of file handles (look for  "VFS:  file-max limit &lt;number&gt; reached"), try increasing this value:

`echo 100000 > /proc/sys/fs/file-max`

Privileged processes (`CAP_SYS_ADMIN`) can override the `file-max` limit.

> 系统范围＋所有进程＋fd

- **`/proc/sys/fs/file-nr`**

This (**read-only**) file contains three numbers: the number of **allocated** file handles (i.e., the number of files presently **opened**); the number of **free** file handles; and the **maximum** number of file handles (i.e., the same value as `/proc/sys/fs/file-max`).  If the number of allocated file handles is close to the  maximum, you should consider increasing the maximum.  **Before Linux 2.6**, the kernel allocated file handles dynamically, but it didn't free them again.  Instead the free file handles were kept in a list for reallocation;  the "free file handles" value  indicates the size of that list.  A large number of free file handles indicates that there was a past peak in the usage of open file handles.  **Since Linux 2.6**, the kernel does deallocate freed file handles, and the "free file handles" value is always zero.


- **`/proc/sys/fs/nr_open`** (since Linux 2.6.25)

(man)
This file imposes ceiling on the value to which the `RLIMIT_NOFILE` resource limit can be raised (see `getrlimit(2)`).  This ceiling is enforced for both **unprivileged** and **privileged** process.  The default value in this file is 1048576.  (Before Linux 2.6.25, the ceiling for `RLIMIT_NOFILE` was hard-coded to the same value.)

(kernel doc)
This denotes the maximum number of file-handles **a process** can allocate. Default value is 1024*1024 (1048576) which should be enough for most machines. Actual limit depends on `RLIMIT_NOFILE` resource limit.

> 进程级别＋fd


- **`nproc`**

```
$ vi /etc/security/limits.conf
# /etc/security/limits.conf
#
#Each line describes a limit for a user in the form:
...
#        - nproc - max number of processes
```

在 `/etc/security/limits.conf` 中指定的 nofile 值的上限由内核参数 nr_open 决定；

> 用户级别＋fd


参考：

- [Documentation/sysctl/fs.txt](https://www.kernel.org/doc/Documentation/sysctl/fs.txt)


----------


> TODO

```
sysctl -w fs.nr_open=20000500
sysctl -w net.core.rmem_max=16384
sysctl -w net.core.wmem_max=16384
net.netfilter.nf_conntrack_tcp_timeout_established = 1200
net.ipv4.tcp_max_syn_backlog = 262144
net.core.netdev_max_backlog = 262144
net.core.somaxconn = 262144
```


> 辅助命令

```
grep conntrack /proc/slabinfo
cat /proc/sys/net/netfilter/nf_conntrack_count
ss -s
cat /proc/net/sockstat|grep -E "TCP:|sockets:"
cat /proc/sys/fs/file-nr
free -m
netstat -n | awk '/^tcp/ {++S[$NF]} END {for(a in S) print a, S[a]}'
```


----------


## 问题梳理

### "TCP: too many orphaned sockets"

在通过 `Ctrl+C` 停止服务的时候，会看到 `dmesg -w` 刷一波如下信息

```
[105172.135420] net_ratelimit: 981 callbacks suppressed
[105172.135427] TCP: too many orphaned sockets
[105172.135452] TCP: too many orphaned sockets
[105172.135466] TCP: too many orphaned sockets
[105172.135477] TCP: too many orphaned sockets
[105172.135485] TCP: too many orphaned sockets
[105172.135495] TCP: too many orphaned sockets
[105172.135506] TCP: too many orphaned sockets
[105172.135519] TCP: too many orphaned sockets
[105172.135527] TCP: too many orphaned sockets
[105172.135537] TCP: too many orphaned sockets
[105177.139974] net_ratelimit: 51638 callbacks suppressed
[105177.139978] TCP: too many orphaned sockets
[105177.139991] TCP: too many orphaned sockets
[105177.140002] TCP: too many orphaned sockets
[105177.140011] TCP: too many orphaned sockets
[105177.140021] TCP: too many orphaned sockets
[105177.140030] TCP: too many orphaned sockets
[105177.140040] TCP: too many orphaned sockets
[105177.140051] TCP: too many orphaned sockets
[105177.140061] TCP: too many orphaned sockets
[105177.140070] TCP: too many orphaned sockets
[105182.143899] net_ratelimit: 44560 callbacks suppressed
[105182.143904] TCP: too many orphaned sockets
[105182.143920] TCP: too many orphaned sockets
[105182.143931] TCP: too many orphaned sockets
...
```

**原因分析**：

当 tcp_mem 的 high 设置被达到后导致（允许所有 tcp sockets 用于排队缓冲数据报的页面量，当内存占用超过此值，系统拒绝分配 socket）；


### "nf_conntrack: table full, dropping packet"

当 server 侧连接数达到 262138 后，无法再建立新的连接，此时输出

```
[35470.975678] nf_conntrack: table full, dropping packet
[35704.617801] nf_conntrack: table full, dropping packet
[35706.621871] nf_conntrack: table full, dropping packet
[35710.633918] nf_conntrack: table full, dropping packet
[35718.649976] nf_conntrack: table full, dropping packet
```


原因分析：

系统设置为 server 端监听 10000-11000 共 1000 个端口；client 端对每个端口发起 50k 个连接，共计可以建立 1000 * 50k =~ 50M 连接；

但实际测试发现，server 侧连接数达到 262138 后，无法再建立新的连接；此时通过其他机器 ping 该机器也会时通时不通，使用 ssh 登录的情况也一样；

结论就是用来跟踪连接的哈希表满了，默认值为 262144（可以查看 `/proc/sys/net/netfilter/nf_conntrack_max`）；


服务器侧：

```
server elapsed=38s connected=243565 failed=0
server elapsed=39s connected=250994 failed=0
server elapsed=40s connected=257323 failed=0
server elapsed=41s connected=262107 failed=0
server elapsed=42s connected=262138 failed=0
server elapsed=43s connected=262138 failed=0
server elapsed=44s connected=262138 failed=0
...
server elapsed=3430s connected=262139 failed=0
```

通过 `dmesg` 可以看到(也可以看 `/var/log/kern.log`)

```
[26636.691175] nf_conntrack: table full, dropping packet
[26656.152537] nf_conntrack: table full, dropping packet
```

客户端侧：

```
client elapsed=31s pending=530 connected=221207 failed=0
client elapsed=32s pending=350 connected=226977 failed=0
client elapsed=33s pending=585 connected=232204 failed=0
client elapsed=34s pending=1024 connected=238480 failed=0
client elapsed=35s pending=481 connected=244613 failed=0
latency 0.32ms
client elapsed=36s pending=526 connected=251707 failed=0
client elapsed=37s pending=174 connected=258567 failed=0
client elapsed=38s pending=1024 connected=262106 failed=0
client elapsed=39s pending=1024 connected=262137 failed=0
client elapsed=40s pending=1024 connected=262137 failed=0
latency 0.53ms
client elapsed=41s pending=1024 connected=262137 failed=0
client elapsed=42s pending=1024 connected=262137 failed=0
...
client elapsed=3419s pending=1024 connected=262138 failed=26624
```

通过 `dmesg` 可以同样看到

```
[26705.484427] nf_conntrack: table full, dropping packet
[26706.853399] nf_conntrack: table full, dropping packet
```

## 扩展

- 每个 TCP 连接大致占用多少内存；
- 100w 连接时需要多少 goroutine ；
- 使用 Erlang 实现的效果；
- 使用 C 实现的效果；
- Golang 运行时调优；


## 案例

360 消息推送：

- 16 台机器，标配：24 个硬件线程，64GB 内存 ；
- Linux Kernel 2.6.32 x86_64 ；
- **单机 80 万并发连接**，load 0.2~0.4，CPU 总使用率 7%~10%，**内存占用20GB (res)** ；
- 目前接入的产品约1280万在线用户 
- **2分钟一次GC**，停顿2秒 (1.0.3 的 GC 不给力，直接升级到 tip，再次吃螃蟹) ；
- 15亿个心跳包/天，占大多数；


京东云消息推送系统：

- 单机并发tcp连接数峰值118w 
- 内存占用23G(Res) 
- Load 0.7左右 
- 心跳包 4k/s 
- gc时间 2-3.x s

