# 折腾 PF_RING 测试

## PF_RING ZC Examples

> 原文地址：[这里](https://github.com/ntop/PF_RING/blob/dev/userland/examples_zc/README.examples)

该文件中包含了一些用于介绍 PF_RING ZC 用法的示例； 

`sysdig` 相关示例请跳转到本文最后；

这里假定你已经加载了某一种支持 ZC 的 PF_RING-aware 驱动程序，以单 RX queue 模式；同样假设存在两个 ports ，即 eth1 和 eth2 ，并且两者直接联通（基于 loopback）；

```shell
PF_RING/kernel# insmod ./pf_ring.ko
PF_RING/drivers/PF_RING_aware/intel/ixgbe/ixgbe-XXX-zc/src# ./load_driver.sh
```

我么使用 single-RX queues 来演示 ZC 的 packet 分发（distribution）能力，并且在如下全部测试中未使用 Intel 的 RSS 功能；

启动 **traffic generator** 并令其按照如下方式一直运行：

```shell
PF_RING/userland/examples_zc# ./zsend -i zc:eth1 -c 1 -g 0 -b 8 -l 60
```

或者（更好的方式为）使用一种具备跨核心均衡 traffic 能力的 traffic generator（例如具有地址 rotating 特点的 packets）；

为了测试最糟可能情景下的代码情况，请发送 60+4 子节的 packets ；

---

- **Example 1. Hash incoming packets and read them on 2 threads**

```shell
PF_RING/userland/examples_zc# ./zbalance -i zc:eth4 -c 4 -m 0 -r 1 -g 2:3
```

> 单进程多线程；   
> 不同线程分别绑定到不同 CPU core 上；   
> Balancer thread 以轮询方式将 packets 分发给 Consumer threads ；   


---

- **Example 2. Hash incoming packets and read them on 2 processes**

```shell
PF_RING/userland/examples_zc# ./zbalance_ipc -i zc:eth2 -c 99 -n 2 -m 0 -g 1

PF_RING/userland/examples_zc# ./zcount_ipc -c 99 -i 0 -g 2
PF_RING/userland/examples_zc# ./zcount_ipc -c 99 -i 1 -g 3
```

> 多进程；   
> Balancer process 以轮询方式将 packets 分发给 2 个 Consumer processes ；   
> Balancer process (zbalance_ipc) 会创建设备 zc:99@0 和 zc:99@1 ；   
> Consumer processes (zcount_ipc) 从指定的 sw queue 上 consume packets ；   
> zbalance_ipc 要先于 zcount_ipc 启动；   


---

- **Example 3. Hash incoming packets and read them on 2 legacy pcap-based applications**

```shell
PF_RING/userland/examples_zc# ./zbalance_ipc -i zc:eth2 -c 99 -n 2 -m 0 -g 1

PF_RING/userland/tcpdump-4.1.1# ./tcpdump -i zc:99@0
PF_RING/userland/tcpdump-4.1.1# ./tcpdump -i zc:99@1
```

> 多进程；   
> Balancer process 以轮询方式将 packets 分发给 2 个 Consumer processes ；   
> Balancer process (zbalance_ipc) 会创建设备 zc:99@0 和 zc:99@1 ；   
> Consumer processes (tcpdump) 从设备 zc:99@0 和 zc:99@1 上 consume packets ；   
> zbalance_ipc 要先于 tcpdump 启动；   

---

- **Example 4. Enqueue incoming packets to a pipeline with 2 threads**

```shell
PF_RING/userland/examples_zc# ./zpipeline -i zc:eth2 -c 99 -g 2:3 
```

> 未看出该程序有何效果；   
> 每个 thread 对应一个 pipeline stage ，thread 数量由 -g 参数控制；   

---

- **Example 5. Enqueue incoming packets to a queue, on another process forward packets from the queue to another queue, send packets from the second queue to an egress interface** 

```shell
PF_RING/userland/examples_zc# ./zpipeline_ipc -i zc:eth2;0 -o zc:eth3;1 -n 2 -c 99 -r 1 -t 2

PF_RING/userland/examples_zc# ./zbounce_ipc -c 99 -i 0 -o 1 -g 3
```

> 调用上面命令时，需要使用 "\" 对 ";" 进行转义；

（注意：`zbounce_ipc` 应用能够运行在 VM 上，用于 multiple VMs 的 pipeline 能够在创建时分配更多的 queues）


### PF_RING ZC Sysdig Examples

我们假设你已经

1. 已经加载了 `pf_rig.ko` 模块
2. 安装了 [sysdig](http://sysdig.org) ，并且同样已经加载了 sysdig 内核模块；你可以通过如下命令进行确认：

```shell
   # lsmod  |grep sysdig
   sysdig_probe          205269  0 
```

由于 PF_RING ZC 需要 huge-pages 才能工作，你需要确保其已按如下方式已配置好：

```shell
# echo 1024 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages
# mkdir /dev/hugepages
# mount -t hugetlbfs nodev /dev/hugepages
```

需要注意的是：

1. 取决于你的安装方式，你可能需要增加 huge pages 的数量；
2. 如果你打算使用 PF_RING ZC 用于网络相关行为，上面已讨论过的 `load_driver.sh` 脚本已经帮你正确设置了 huge pages ；

---

- **Example 1. Hash incoming sysdig events and read them on 2 threads balancing them per PID**

```shell
PF_RING/userland/examples_zc# ./zbalance -i sysdig -c 4 -m 0 -r 1 -g 2:3
```

---

- **Example 2. Hash incoming packets and read them on 2 processes**

```shell
PF_RING/userland/examples_zc# ./zbalance_ipc -i sysdig -c 99 -n 2 -m 0 -g 1

PF_RING/userland/examples_zc# ./zcount_ipc -c 99 -i 0 -g 2 -s
PF_RING/userland/examples_zc# ./zcount_ipc -c 99 -i 1 -g 3 -s
```

---

- **Example 3. Hash incoming packets and read them on 2 non-ZC applications**

```shell
PF_RING/userland/examples_zc# ./zbalance_ipc -i zc:eth2 -c 99 -n 2 -m 0 -g 1

PF_RING/userland/examples# ./pfcount -i zc:99@0 -v 1 -q
PF_RING/userland/examples# ./pfcount -i zc:99@1 -v 1 -q
```

---

- **Example 4. Enqueue incoming sysdig events to a pipeline with 2 threads**

```shell
PF_RING/userland/examples_zc# ./zpipeline -i sysdig -c 99 -g 2:3 
```

---

- **Example 5. Enqueue incoming sysdig events to a queue, on another process forward packets from the queue to another queue, send packets from the second queue to an egress interface (perhaps we should first encapsulate the events into a ethernet frame for best results)**

```
PF_RING/userland/examples_zc# ./zpipeline_ipc -i sysdig;0 -o zc:eth3;1 -n 2 -c 99 -r 1 -t 2

PF_RING/userland/examples_zc# ./zbounce_ipc -c 99 -i 0 -o 1 -g 3
```

(note that the zbounce_ipc application can run on a VM, a pipeline with multiple VMs can be created allocating more queues)


----------


## zsend

```
[root@wg-esm-hc-1 examples_zc]# ./zsend -h
zsend - (C) 2014-17 ntop.org
Using PFRING_ZC v.6.5.0.170303
A traffic generator able to replay synthetic udp packets or hex from standard input.
Usage:    zsend -i <device> -c <cluster id>
                [-h] [-g <core id>] [-r <rate>] [-p <pps>] [-l <len>] [-n <num>]
                [-b <num>] [-N <num>] [-S <core id>] [-P <core id>]
                [-z] [-a] [-Q <sock>] [-f <.pcap file>] [-m <MAC>] [-o <num>]

-h              Print this help
-i <device>     Device name (optional: do not specify a device to create a cluster with a sw queue)
-c <cluster id> Cluster id
-f <.pcap file> Send packets as read from a pcap file
-m <dst MAC>    Reforge destination MAC (format AA:BB:CC:DD:EE:FF)
-o <num>        Offset for generated IPs (-b) or packets in pcap (-f)
-g <core id>    Bind this app to a core
-p <pps>        Rate (packets/s)
-r <Gbps rate>  Rate to send (example -r 2.5 sends 2.5 Gbit/sec, -r -1 pcap capture rate)
-l <len>        Packet len (bytes)
-n <num>        Number of packets
-b <num>        Number of different IPs
-N <num>        Simulate a producer for n2disk multi-thread (<num> threads)
-S <core id>    Append timestamp to packets, bind time-pulse thread to a core
-P <core id>    Use a time-pulse thread to control transmission rate, bind the thread to a core
-z              Use burst API
-a              Active packet wait
-Q <sock>       Enable VM support to attach a consumer from a VM (<sock> is a QEMU monitor sockets)
[root@wg-esm-hc-1 examples_zc]#
```

- non-ZC

```
[root@wg-esm-hc-1 examples_zc]# ./zsend -i eno49 -c 1 -g 0 -b 8 -l 60
Sending packets to eno49
Estimated CPU freq: 2589571000 Hz
=========================
Absolute Stats: 887'670 pkts - 74'564'280 bytes
=========================

=========================
Absolute Stats: 1'783'366 pkts - 149'802'744 bytes
Actual Stats: 895'596.58 pps - 0.60 Gbps [75238464 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 2'673'242 pkts - 224'552'328 bytes
Actual Stats: 889'760.33 pps - 0.60 Gbps [74749584 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 3'566'465 pkts - 299'583'060 bytes
Actual Stats: 893'147.08 pps - 0.60 Gbps [75030732 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 4'459'209 pkts - 374'573'556 bytes
Actual Stats: 892'677.94 pps - 0.60 Gbps [74990496 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 5'351'597 pkts - 449'534'148 bytes
Actual Stats: 892'330.89 pps - 0.60 Gbps [74960592 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 6'244'310 pkts - 524'522'040 bytes
Actual Stats: 892'654.08 pps - 0.60 Gbps [74987892 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 7'137'173 pkts - 599'522'532 bytes
Actual Stats: 892'804.07 pps - 0.60 Gbps [75000492 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 8'030'022 pkts - 674'521'848 bytes
Actual Stats: 892'793.64 pps - 0.60 Gbps [74999316 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 8'923'046 pkts - 749'535'864 bytes
Actual Stats: 892'967.74 pps - 0.60 Gbps [75014016 bytes / 1.0 sec]
=========================

^CLeaving...
=========================
Absolute Stats: 8'923'111 pkts - 749'541'324 bytes
Actual Stats: 890'410.95 pps - 0.60 Gbps [5460 bytes / 0.0 sec]
=========================

=========================
Absolute Stats: 8'923'111 pkts - 749'541'324 bytes
Actual Stats: 0.00 pps - 0.00 Gbps [0 bytes / 0.0 sec]
=========================

[root@wg-esm-hc-1 examples_zc]#
```

- ZC

```
[root@wg-esm-hc-1 examples_zc]# ./zsend -i zc:eno49 -c 1 -g 0 -b 8 -l 60
#########################################################################
# ERROR: You do not seem to have a valid PF_RING ZC license 6.5.0.170303 for eno49 [Intel 10 Gbit ixgbe 82599-based]
# ERROR: Please get one at http://shop.ntop.org/.
#########################################################################
# We're now working in demo mode with packet capture and
# transmission limited to 5 minutes
#########################################################################
#########################################################################
# ERROR: You do not seem to have a valid PF_RING ZC license 6.5.0.170303 for eno49 [Intel 10 Gbit ixgbe 82599-based]
# ERROR: Please get one at http://shop.ntop.org/.
#########################################################################
Sending packets to zc:eno49
Estimated CPU freq: 2593191000 Hz
=========================
Absolute Stats: 14'868'342 pkts - 1'248'940'728 bytes
=========================

=========================
Absolute Stats: 29'748'153 pkts - 2'498'844'852 bytes
Actual Stats: 14'877'876.87 pps - 10.00 Gbps [1249904124 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 44'630'645 pkts - 3'748'974'180 bytes
Actual Stats: 14'880'959.26 pps - 10.00 Gbps [1250129328 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 59'512'756 pkts - 4'999'071'504 bytes
Actual Stats: 14'880'875.88 pps - 10.00 Gbps [1250097324 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 74'393'940 pkts - 6'249'090'960 bytes
Actual Stats: 14'880'097.75 pps - 10.00 Gbps [1250019456 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 89'276'380 pkts - 7'499'215'920 bytes
Actual Stats: 14'881'383.42 pps - 10.00 Gbps [1250124960 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 104'157'658 pkts - 8'749'243'272 bytes
Actual Stats: 14'880'161.98 pps - 10.00 Gbps [1250027352 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 119'039'853 pkts - 9'999'347'652 bytes
Actual Stats: 14'881'049.15 pps - 10.00 Gbps [1250104380 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 133'921'996 pkts - 11'249'447'664 bytes
Actual Stats: 14'881'012.04 pps - 10.00 Gbps [1250100012 bytes / 1.0 sec]
=========================

=========================
Absolute Stats: 148'803'112 pkts - 12'499'461'408 bytes
Actual Stats: 14'880'029.75 pps - 10.00 Gbps [1250013744 bytes / 1.0 sec]
=========================

^CLeaving...
=========================
Absolute Stats: 153'943'848 pkts - 12'931'283'232 bytes
Actual Stats: 14'881'101.39 pps - 10.00 Gbps [431821824 bytes / 0.3 sec]
=========================

=========================
Absolute Stats: 153'943'849 pkts - 12'931'283'316 bytes
Actual Stats: 25'641.02 pps - 0.02 Gbps [84 bytes / 0.0 sec]
=========================

[root@wg-esm-hc-1 examples_zc]#
```

## zcount

```
[root@wg-esm-hc-1 examples_zc]# ./zcount -h
zcount - (C) 2014 ntop.org
Using PFRING_ZC v.6.5.0.170303
A simple packet counter application.

Usage:   zcount -i <device> -c <cluster id>
                [-h] [-g <core id>] [-R] [-H] [-S <core id>] [-v] [-a]

-h              Print this help
-i <device>     Device name
-c <cluster id> Cluster id
-g <core id>    Bind this app to a core
-a              Active packet wait
-R              Test hw filters adding a rule (Intel 82599)
-H              High stats refresh rate (workaround for drop counter on 1G Intel cards)
-S <core id>    Pulse-time thread for inter-packet time check
-C              Check license
-v              Verbose
[root@wg-esm-hc-1 examples_zc]#
```

- non-ZC

```
[root@wg-esm-hc-1 examples_zc]# ./zcount -i eno49 -c 0
=========================
Absolute Stats: 1 pkts (0 drops) - 363 bytes
=========================

=========================
Absolute Stats: 2 pkts (0 drops) - 506 bytes
Actual Stats: 1.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 2 pkts (0 drops) - 506 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 3 pkts (0 drops) - 649 bytes
Actual Stats: 1.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 3 pkts (0 drops) - 649 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

^CLeaving...
=========================
Absolute Stats: 3 pkts (0 drops) - 649 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 3 pkts (0 drops) - 649 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

[root@wg-esm-hc-1 examples_zc]#
```

- ZC

```
[root@wg-esm-hc-1 examples_zc]# ./zcount -i zc:eno49 -c 0
#########################################################################
# ERROR: You do not seem to have a valid PF_RING ZC license 6.5.0.170303 for eno49 [Intel 10 Gbit ixgbe 82599-based]
# ERROR: Please get one at http://shop.ntop.org/.
#########################################################################
# We're now working in demo mode with packet capture and
# transmission limited to 5 minutes
#########################################################################
#########################################################################
# ERROR: You do not seem to have a valid PF_RING ZC license 6.5.0.170303 for eno49 [Intel 10 Gbit ixgbe 82599-based]
# ERROR: Please get one at http://shop.ntop.org/.
#########################################################################
=========================
Absolute Stats: 1 pkts (0 drops) - 143 bytes
=========================

=========================
Absolute Stats: 1 pkts (0 drops) - 143 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 2 pkts (0 drops) - 286 bytes
Actual Stats: 1.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 2 pkts (0 drops) - 286 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 3 pkts (0 drops) - 429 bytes
Actual Stats: 1.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 3 pkts (0 drops) - 429 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

^CLeaving...
=========================
Absolute Stats: 3 pkts (0 drops) - 429 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

=========================
Absolute Stats: 3 pkts (0 drops) - 429 bytes
Actual Stats: 0.00 pps (0.00 drops) - 0.00 Gbps
=========================

[root@wg-esm-hc-1 examples_zc]#
```

## pfcount

```
[root@wg-esm-hc-1 examples]# ./pfcount -h
pfcount - (C) 2005-17 ntop.org

-h              Print this help
-i <device>     Device name. Use:
                - ethX@Y for channels
                - zc:ethX for ZC devices
                - sysdig for capturing sysdig events
-n <threads>      Number of polling threads (default 1)
-f <filter>       BPF filter
-e <direction>    0=RX+TX, 1=RX only, 2=TX only
-l <len>          Capture length
-g <core_id>      Bind this app to a core
-d <device>       Device on which incoming packets are copied
-w <watermark>    Watermark
-p <poll wait>    Poll wait (msec)
-b <cpu %>        CPU pergentage priority (0-99)
-a                Active packet wait
-N <num>          Read <num> packets and exit
-q                Force printing packets as sysdig events with -v
-m                Long packet header (with PF_RING extensions)
-r                Rehash RSS packets
-c <cluster id>   Cluster ID (kernel clustering)
-H <cluster hash> Cluster hash type (kernel clustering)
                   2 - src ip,           dst ip
                   3 - src ip, src port, dst ip, dst port
                   4 - src ip, src port, dst ip, dst port, proto (default)
                   0 - src ip, src port, dst ip, dst port, proto, vlan
                   5 - src ip, src port, dst ip, dst port, proto for TCP, src ip, dst ip otherwise
                   7 - tunneled src ip,           dst ip
                   8 - tunneled src ip, src port, dst ip, dst port
                   9 - tunneled src ip, src port, dst ip, dst port, proto (default)
                   6 - tunneled src ip, src port, dst ip, dst port, proto, vlan
                  10 - tunneled src ip, src port, dst ip, dst port, proto for TCP, src ip, dst ip otherwise
                   1 - round-robin
-s              Enable hw timestamping
-S              Do not strip hw timestamps (if present)
-t              Touch payload (to force packet load on cache)
-M              Packet memcpy (to test memcpy speed)
-C <mode>       Work with the adapter in chunk mode (1=chunk API, 2=packet API)
-x <path>       File containing strings to search string (case sensitive) on payload.
-o <path>       Dump matching packets onto the specified pcap (need -x).
-u <1|2>        For each incoming packet add a drop rule (1=hash, 2=wildcard rule)
-v <mode>       Verbose [1: verbose, 2: very verbose (print packet payload)]
-z <mode>       Enabled hw timestamping/stripping. Currently the supported TS mode are:
                ixia	Timestamped packets by ixiacom.com hardware devices
-L              List all interafces and exit (use -v for more info)
[root@wg-esm-hc-1 examples]#
```

- non-ZC

```
[root@wg-esm-hc-1 examples]# ./pfcount -i eno49
Using PF_RING v.6.5.0
Capturing from eno49 [mac: 14:02:EC:82:58:CC][if_index: 19][speed: 10000Mb/s]
# Device RX channels: 1
# Polling threads:    1
Dumping statistics on /proc/net/pf_ring/stats/5509-eno49.52
=========================
Absolute Stats: [4 pkts total][0 pkts dropped][0.0% dropped]
[4 pkts rcvd][336 bytes rcvd]
=========================

=========================
Absolute Stats: [10 pkts total][0 pkts dropped][0.0% dropped]
[10 pkts rcvd][923 bytes rcvd][10.00 pkt/sec][0.01 Mbit/sec]
=========================
Actual Stats: [6 pkts rcvd][1'000.07 ms][6.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [14 pkts total][0 pkts dropped][0.0% dropped]
[14 pkts rcvd][1'259 bytes rcvd][7.00 pkt/sec][0.01 Mbit/sec]
=========================
Actual Stats: [4 pkts rcvd][1'000.08 ms][4.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [20 pkts total][0 pkts dropped][0.0% dropped]
[20 pkts rcvd][1'822 bytes rcvd][6.67 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [6 pkts rcvd][1'000.09 ms][6.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [23 pkts total][0 pkts dropped][0.0% dropped]
[23 pkts rcvd][2'074 bytes rcvd][5.75 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [3 pkts rcvd][1'000.08 ms][3.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [27 pkts total][0 pkts dropped][0.0% dropped]
[27 pkts rcvd][2'469 bytes rcvd][5.40 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [4 pkts rcvd][1'000.08 ms][4.00 pps][0.00 Gbps]
=========================

^CLeaving...
=========================
Absolute Stats: [31 pkts total][0 pkts dropped][0.0% dropped]
[31 pkts rcvd][2'805 bytes rcvd][5.19 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [4 pkts rcvd][968.32 ms][4.13 pps][0.00 Gbps]
=========================

[root@wg-esm-hc-1 examples]#
```

- ZC

```
[root@wg-esm-hc-1 examples]# ./pfcount -i zc:eno49
#########################################################################
# ERROR: You do not seem to have a valid PF_RING ZC license 6.5.0.170303 for eno49 [Intel 10 Gbit ixgbe 82599-based]
# ERROR: Please get one at http://shop.ntop.org/.
#########################################################################
# We're now working in demo mode with packet capture and
# transmission limited to 5 minutes
#########################################################################
Using PF_RING v.6.5.0
Capturing from zc:eno49 [mac: 14:02:EC:82:58:CC][if_index: 19][speed: 10000Mb/s]
# Device RX channels: 1
# Polling threads:    1
Dumping statistics on /proc/net/pf_ring/stats/5922-eno49.55
=========================
Absolute Stats: [0 pkts total][0 pkts dropped][0.0% dropped]
[0 pkts rcvd][0 bytes rcvd]
=========================

=========================
Absolute Stats: [1 pkts total][0 pkts dropped][0.0% dropped]
[1 pkts rcvd][143 bytes rcvd][1.00 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [1 pkts rcvd][1'000.13 ms][1.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [1 pkts total][0 pkts dropped][0.0% dropped]
[1 pkts rcvd][143 bytes rcvd][0.50 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [0 pkts rcvd][1'000.13 ms][0.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [2 pkts total][0 pkts dropped][0.0% dropped]
[2 pkts rcvd][286 bytes rcvd][0.67 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [1 pkts rcvd][1'000.04 ms][1.00 pps][0.00 Gbps]
=========================

=========================
Absolute Stats: [3 pkts total][0 pkts dropped][0.0% dropped]
[3 pkts rcvd][434 bytes rcvd][0.75 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [1 pkts rcvd][1'000.05 ms][1.00 pps][0.00 Gbps]
=========================

^CLeaving...
=========================
Absolute Stats: [3 pkts total][0 pkts dropped][0.0% dropped]
[3 pkts rcvd][434 bytes rcvd][0.73 pkt/sec][0.00 Mbit/sec]
=========================
Actual Stats: [0 pkts rcvd][99.27 ms][0.00 pps][0.00 Gbps]
=========================

[root@wg-esm-hc-1 examples]#
```

