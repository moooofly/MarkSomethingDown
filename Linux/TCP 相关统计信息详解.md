# TCP 相关统计信息详解

我们知道 TCP 相关统计信息包含在如下文件中

- `/proc/net/netstat`
- `/proc/net/snmp`

可以通过 `netstat -s` 的输出信息进行确认；

```shell
root@vagrant-ubuntu-trusty:~] $ strace -e open netstat -s
...
open("/proc/meminfo", O_RDONLY|O_CLOEXEC) = 3
open("/proc/net/snmp", O_RDONLY)        = 3
...
Ip:
    23207 total packets received
    26 with invalid addresses
    0 forwarded
    0 incoming packets discarded
    23181 incoming packets delivered
    17146 requests sent out
    40 outgoing packets dropped
Icmp:
    83 ICMP messages received
    0 input ICMP message failed.
    ICMP input histogram:
        destination unreachable: 83
    80 ICMP messages sent
    0 ICMP messages failed
    ICMP output histogram:
        destination unreachable: 80
IcmpMsg:
        InType3: 83
        OutType3: 80
Tcp:
    2577 active connections openings
    3 passive connection openings
    2570 failed connection attempts
    0 connection resets received
    2 connections established
    21855 segments received
    16307 segments send out
    0 segments retransmited
    0 bad segments received.
    2574 resets sent
Udp:
    659 packets received
    80 packets to unknown port received.
    0 packet receive errors
    767 packets sent
    IgnoredMulti: 506
UdpLite:
open("/proc/net/netstat", O_RDONLY)     = 3
TcpExt:
    4 TCP sockets finished time wait in fast timer
    55 delayed acks sent
    7 delayed acks further delayed because of locked socket
    Quick ack mode was activated 4 times
    8 packets directly queued to recvmsg prequeue.
    3551 packet headers predicted
    85 acknowledgments not containing data payload received
    10436 predicted acknowledgments
    TCPRcvCoalesce: 462
    TCPOrigDataSent: 10445
    TCPHystartTrainDetect: 1
    TCPHystartTrainCwnd: 18
IpExt:
    InBcastPkts: 506
    InOctets: 3732899
    OutOctets: 1305530
    InBcastOctets: 44382
    InNoECTPkts: 24856
+++ exited with 0 +++
root@vagrant-ubuntu-trusty:~] $
```

可以看出

`/proc/net/netstat` 文件中包含：

- TcpExt
- IpExt

`/proc/net/snmp` 文件中包含：

- Ip
- Icmp
- IcmpMsg
- Tcp
- Udp
- UdpLite

本文主要讨论 TCP 和 IP 协议相关内容，并按照类别将上述文件内容进行了整合；

## 计数器分类

| 类别 | counters |
| --- | --- |
| 常量 | RtoAlgorithm <br> RtoMin <br> RtoMax <br> MaxConn
| 建链统计 | ActiveOpens <br> PassiveOpens <br> AttemptFails <br> CurrEstab <br> EstabResets |
| 数据包统计 | InSegs <br> OutSegs <br> RetransSegs <br> InErrs <br> OutRsts <br> InCsumErrors <br> EmbryonicRsts |
| syncookies 相关 | SyncookiesSent <br> SyncookiesRecv <br> SyncookiesFailed |
| TIME_WAIT 相关 | TW <br> TWRecycled <br> TWKilled <br> TCPTimeWaitOverflow |
| RTO 相关 | TCPTimeouts <br> TCPSpuriousRTOs <br> TCPLossProbes <br> TCPLossProbeRecovery <br> TCPRenoRecoveryFail <br> TCPSackRecoveryFail <br> TCPRenoFailures <br> TCPSackFailures <br> TCPLossFailures |
| Retrans 相关 | TCPFastRetrans <br> TCPForwardRetrans <br> TCPSlowStartRetrans <br> TCPLostRetransmit <br> TCPRetransFail |
| FastOpen 相关 | TCPFastOpenActive <br> TCPFastOpenPassive <br> TCPFastOpenPassiveFail <br> TCPFastOpenListenOverflow <br> TCPFastOpenCookieReqd |
| MD5 相关 | TCPMD5NotFound <br> TCPMD5Unexpected |
| DelayedACK 相关 | DelayedACKs <br> DelayedACKLocked <br> DelayedACKLost <br> TCPSchedulerFailed |
| DSACK 相关 | TCPDSACKOldSent <br> TCPDSACKOfoSent <br> TCPDSACKRecv <br> TCPDSACKOfoRecv <br> TCPDSACKIgnoredOld <br> TCPDSACKIgnoredNoUndo |
| Reorder 相关 | TCPFACKReorder <br> TCPSACKReorder <br> TCPRenoReorder <br> TCPTSReorder |
| Recovery 相关 | TCPRenoRecovery <br> TCPSackRecovery <br> TCPRenoRecoveryFail <br> TCPSackRecoveryFail|
| Abort 相关 | TCPAbortOnData <br> TCPAbortOnClose <br> TCPAbortOnMemory <br> TCPAbortOnTimeout <br> TCPAbortOnLingerTCPAbortFailed |
| Reset 相关 | |
| 内存 prune | PruneCalled <br> RcvPruned <br> OfoPruned <br> TCPMemoryPressures |
| PAWS 相关 | PAWSPassive <br> PAWSActive <br> PAWSEstab |
| Listen 相关 | ListenOverflows <br> ListenDrops |
| Undo 相关 | TCPFullUndo <br> TCPPartialUndo <br> TCPDSACKUndo <br> TCPLossUndo |
| 快速路径与慢速路径 | TCPHPHits <br> TCPHPHitsToUser <br> TCPPureAcks <br> TCPHPAcks |

### 常量

这些常量是 Linux 3.10 中的默认值，仅在升级了内核版本时才需要关心一下这些值的变化。

| 名称 | 含义 |
| --- | --- |
| RtoAlgorithm | 用于计算 RTO 的算法，默认为 1 ，RTO 算法与 RFC2698 一致 |
| RtoMin | 	限定 RTO 的最小值，默认值为 1/5HZ，即 **200ms** |
| RtoMax | 限定 RTO 的最大值，默认值为 120HZ，即 **120s** |
| MaxConn | TCP 流数量的上限，协议栈本身并不会限制 TCP 连接总数，默认值为 -1 |

### 建链统计

这些统计值中，只有 CurrEstab 反应的是系统当前状态，而其他值则是反应的历史状态；同时需要注意的是，这些计数器将 ESTABLISHED 和 CLOSE-WAIT 状态都作为当前连接数。

可以这么理解：这两个状态都认为 local => peer 方向的连接未被关闭；

| 名称 | 含义 |
| --- | --- |
| ActiveOpens | 主动建链次数，对应 CLOSE => SYN-SENT 次数； <br> 在 `tcp_connect()` 函数中计数； <br> 相当于 SYN 包的发送次数（但不包含重传次数） |
| PassiveOpens | 被动建链次数，RFC 原意对应 LISTEN => SYN-RECV 次数，但 Linux 实现选择在三次握手成功后才加 1 （即在建立 tcp_sock 结构体后） |
| AttemptFails |  建链失败次数，即如下三项之和 <br> a) SYN-SENT => CLOSE 次数 <br> b) SYN-RECV => CLOSE 次数 <br> c) SYN-RECV => LISTEN 次数 <br><br> 回 CLOSE 部分在 `tcp_done()` 函数中计数 <br> 回 LISTEN 部分在 `tcp_check_req()` 中计数 |
| EstabResets | 连接被 reset 次数，即如下两项之和 <br> a) ESTABLISHED => CLOSE 次数 <br> b) CLOSE-WAIT => CLOSE 次 <br><br> 在 `tcp_set_state()` 函数中，如果之前的状态是TCP_CLOSE_WAIT 或 TCP_ESTABLISHED 就加 1 |
| CurrEstab | 处于 ESTABLISHED 和 CLOSE-WAIT 状态的 TCP 流数 <br> 在 `tcp_set_state()` 中进行处理 <br> 实现体现的是进入 ESTABLISHED 之后，进入 CLOSE 之前的 TCP 流数 |

### 数据包统计


这些统计值也是历史值，独立的来看意义并不大。一般可统计一段时间内的变化，关注以下几个指标

-  **TCP 层重传率**：`ΔRetransSegs / ΔOutSegs` ；该值越小越好，如果超过 20% 则应该引起注意（这个值根据实际情况而定）；
-  **Reset 发送频率**：`ΔOutRsts / ΔOutSegs` ；该值越小越好，一般应该在 1% 以内；
-  **错误包占比**：`ΔInErrs / ΔInSegs` ；该值越小越好，一般应该在 1% 以内，同时由 checksum 导致的问题包应该更低；

| 名称 | 含义 |
| --- | --- |
| InSegs | 所有收到的 TCP 包，即使是个错误包 <br> 	在 `tcp_v4_rcv()` 和 `tcp_v6_rcv()` 中计数 |
| OutSegs | 所有发送出去的 TCP 包，包括 <br><br> a) 新数据包 <br> b) 重传数据包 <br> c) syn 包 <br> d) synack 包 <br> e) reset 包 <br><br> `tcp_v4_send_reset()` 中统计 reset 包 <br> `tcp_v4_send_ack()` 中统计 SYN-RECV 和 TIME-WAIT 状态下发送的 ACK 包 <br> `tcp_v6_send_response()` 中统计 ipv6 相应数据 <br> `tcp_make_synack()` 中统计发送的 SYNACK 包 <br> `tcp_transmit_skb()` 中统计所有的其他包 |
| RetransSegs | 所有重传出去的 TCP 包 <br><br> `tcp_v4_rtx_synack()` 和 `tcp_v6_rtx_synack()` 中统计重传的 SYNACK 包 <br> `tcp_retransmit_skb()` 中统计其他重传包 |
| InErrs | 	所有收到的有问题的 TCP 包数量，比如 checksum 有问题 <br><br> `tcp_validate_incoming()` 中统计 seq 有问题的包 <br> `tcp_rcv_established()`、`tcp_v4_do_rcv()`、`tcp_v4_rcv()`、`tcp_v6_do_rcv()`、`tcp_v6_rcv()` 中根据 checksum 来判断出错误包 |
| OutRsts | 发送的带 RST 标记的 TCP 包数量 <br><br> 在 `tcp_v4_send_reset()`、`tcp_send_active_reset()`、`tcp_v6_send_response()` 中统计 |
| InCsumErrors | 收到的 checksum 有问题的数据包数量 <br><br> 属于 3.10 相对于 2.6.32 新增的内容，算是细化 InErrs 统计，InErrs 中应该只有*小部分*属于该类型 |
| EmbryonicRsts | 在 SYN-RECV 状态收到带 RST/SYN 标记的包个数 |

### Syncookies 相关

syncookies 一般不会被触发，只有在 `tcp_max_syn_backlog` 队列被占满时才会被触发；因此 SyncookiesSent 和 SyncookiesRecv 一般情况下应该是 0 。但是 SyncookiesFailed 的值和 syncookies 机制是否被触发没有直接关系，因此可能不为 0 ，原因在于：SyncookiesFailed 值的计算方式为：当一个处于 LISTEN 状态的 socket 收到一个不带 SYN 标记的数据包时，就会调用 `cookie_v4_check()` 尝试验证 cookie 信息。而如果验证失败，值加 1 。

| 名称 | 含义 |
| --- | --- |
| SyncookiesSent | 使用 syncookie 技术发送的 syn/ack 包个数 |
| SyncookiesRecv | 收到携带有效 syncookie 信息包个数 |
| SyncookiesFailed | 收到携带无效 syncookie 信息包个数 |

> 注：syncookies 机制是为应对 syn flood 攻击而被提出来的。

### TIME-WAIT 相关

TIME-WAIT 状态是 TCP 协议状态机中的重要一环，服务器设备一般都有非常多处于 TIME-WAIT 状态的 socket ，如果是在主要提供 HTTP 服务的设备上，TW 值应该接近 TcpPassiveOpens 值。

一般情况下，`sysctl_tcp_tw_reuse` 和 `sysctl_tcp_tw_recycle` 都是不推荐开启的。所以 TWKilled 和 TWRecycled 都应该是 0 。同时 TCPTimeWaitOverflow 也应该是 0 ，否则就意味着内存使用方面出了大问题。


| 名称 | 含义 |
| --- | --- |
| TW | 经过正常时间（`TCP_TIMEWAIT_LEN`）结束 TW 状态的 socket 数量 |
| TWRecycled | TIMEWAIT socket 被复用的次数；只有在 `sysctl_tcp_tw_reuse` 开启时，才可能加 1 |
| TWKilled | 经过更短时间结束 TW 状态的 socket 数量；只有在 `net.ipv4.tcp_tw_recycle` 开启时，调度 TW timer 时才可能用更短的 timeout 值 |
| TCPTimeWaitOverflow | 如果没有内存分配 TIMEWAIT 结构体，则加 1 |

### RTO 相关

RTO 超时对 TCP 性能的影响是巨大的，因此关心 RTO 超时的次数也非常必要。

当然 3.10 中的 TLP 机制能够减少一定量的 TCPTimeouts 数，将其转换为快速重传。


| 名称 | 含义 |
| --- | --- |
| TCPTimeouts | RTO timer 第一次超时的次数，仅包含直接超时的情况 |
| TCPSpuriousRTOs | 通过 F-RTO 机制发现的虚假超时个数 |
| TCPLossProbes |  Probe Timeout(PTO) 导致发送 Tail Loss Probe (TLP) 包的次数 |
| TCPLossProbeRecovery | 丢失包刚好被 TLP 探测包修复的次数 |
| TCPRenoRecoveryFail | 先进入 Recovery 阶段，然后又 RTO 的次数，对端不支持 SACK 选项 |
| <br>Fail | 先进入 Recovery 阶段，然后又 RTO 的次数，对端支持 SACK 选项 |
| TCPRenoFailures | 先进 TCP_CA_Disorder 阶段，然后又 RTO 超时的次数，对端不支持 SACK 选项 |
| TCPSackFailures | 先进 TCP_CA_Disorder 阶段，然后又 RTO 超时的次数，对端支持 SACK 选项 |
| TCPLossFailures | 先进 TCP_CA_Loss 阶段，然后又 RTO 超时的次数 |

### Retrans 相关

这些计数器统计的重传包，都不是由于 RTO 超时导致的重传数量；
如果结合 RetransSegs 统计来看，如果这些非 RTO 导致的重传占比较大的话，也算是不幸中的万幸。
另外 LostRetransmit 的数量应该偏低比较好，重传包如果都大量被丢弃，则真的要注意了。

| 名称 | 含义 |
| --- | --- |
| TCPLostRetransmit | 丢失的重传 SBK 数量，没有 TSO 时，等于丢失的重传包数量 |
| TCPFastRetrans | 成功快速重传的 SKB 数量 |
| TCPForwardRetrans | 成功 ForwardRetrans 的 SKB 数量，ForwardRetrans 重传的序号高于 retransmit_high 的数据 |
| TCPSlowStartRetrans | 成功在 Loss 状态发送的重传 SKB 数量，而且这里仅记录非 RTO 超时进入 Loss 状态下的重传数量；目前找到的一种非 RTO 进入 Loss 状态的情况就是：`tcp_check_sack_reneging()` 函数发现接收端违反(renege)了之前的 SACK 信息时，会进入 Loss 状态 |
| TCPRetransFail | 尝试 FastRetrans、ForwardRetrans、SlowStartRetrans 重传失败的次数 |

### FastOpen

**TCP FastOpen (`TFO`)** 技术是 Google 提出来减少三次握手开销的技术，核心原理就是在第一次建链时，由 server 计算出一个 cookies 发给 client ，之后 client 向 server 再次发起建链请求时，就可以携带该 cookies 信息以验明正身。如果 cookies 验证通过，则 server 可以不等三次握手的最后一个 ACK 包，就将 client 放在 SYN 包里面的数据传递给应用层。

在 3.10 内核中，`TFO` 由 `sysctl_tcp_fastopen` 开关控制，默认值为 0(关闭)。而且 `sysctl_tcp_fastopen` 目前也是推荐关闭的，因为网络中有些 middle box 会丢弃那些带有不认识 option 的 SYN 包；所以正常情况下，这些值也应该都是 0 ，当然如果收到过某些不怀好意的、带 TFO cookies 信息的 SYN 包，TCPFastOpenPassive 计数器就可能不为 0 。

| 名称 | 含义 |
| --- | --- |
| TCPFastOpenActive | 主动发送的、带 TFO cookie 的 SYN 包个数 |
| TCPFastOpenActiveFail | 基于 TFO 主动建链失败次数 |
| TCPFastOpenPassive | 收到带 TFO cookie 的 SYN 包个数 |
| TCPFastOpenPassiveFail | 基于 TFO 被动建链失败次数 |
| TCPFastOpenListenOverflow | TFO 请求数超过 listener queue 设置上限，则加 1 |
| TCPFastOpenCookieReqd | 收到一个请求 TFO cookies 的 SYN 包时，则加 1 |

### MD5

TCP MD5 Signature 选项是为提高 BGP Session 的安全性而提出的，详见 [RFC 2385](https://tools.ietf.org/html/rfc2385) 。因此内核中是以编译选项，而不是 sysctl 接口来配置是否使用该功能的。如果内核编译时的 CONFIG_TCP_MD5SIG 选项未配置，则不会支持 TCPMD5Sig ，下面两个计数器也就只能是 0 ；

| 名称 | 含义 |
| --- | --- |
| TCPMD5NotFound | 希望收到带 MD5 选项的包，但是包里面没有 MD5 选项 |
| TCPMD5Unexpected | 不希望收到带 MD5 选项的包，但是包里面有 MD5 选项 |


### DelayedACK

DelayedACK 是内核中默认支持的，但即使使用 DelayedACKs ，每收到两个数据包也必须发送一个 ACK 。所以 DelayedACKs 可以估算为发送出去的 ACK 数量的一半。

同时 DelayedACKLocked 反映的是应用与内核争抢 socket 的次数，如果占 DelayedACKs 比例过大，可能就需要看看应用程序是否有问题了。

| 名称 | 含义 |
| --- | --- |
| DelayedACKs | 尝试发送 delayed ack 的次数，包括未成功发送的次数 |
| DelayedACKLocked | 由于应用锁住了 socket ，而无法发送（即未成功发送）delayed ack 的次数 |
| DelayedACKLost | TODO |
| TCPSchedulerFailed | 如果在 delayed ack 处理函数中发现 prequeue 还有数据，就加 1 ；数据放到 prequeue ，就是想 user 能尽快处理。如果任由数据，则可能 user 行为调度效果不好，这个值应该非常接近于零才正常 |

### DSACK

该类型计数器统计的是收/发 DSACK 信息次数。

DSACKOldSent + DSACKOfoSent 可以当做是发送出的 DSACK 信息的次数，而且概率上来讲 OldSent 应该占比更大。

同理，DSACKRecv 的数量也应该远多于 DSACKOfoRecv 的数量。

另外，DSACK 信息的发送是需要 `sysctl_tcp_dsack` 开启的，如果发现 sent 两个计数器为零，则要检查一下了。

一般还是建议开启 dsack 选项；

| 名称 | 含义 |
| --- | --- |
| TCPDSACKOldSent | 如果收到的重复数据包序号比rcv_nxt(接收端想收到的下一个序号)小，则增加oldsent |
| TCPDSACKOfoSent | 如果收到的重复数据包序号比rcv_nxt大，则是一个乱序的重复数据包，增加ofosent |
| TCPDSACKRecv | 收到的old dsack信息次数，判断old的方法：dsack序号小于ACK号 |
| TCPDSACKOfoRecv | 收到的Ofo dsack信息次数 |
| TCPDSACKIgnoredOld | 当一个dsack block被判定为无效，且设置过undo_marker，则加1 |
| TCPDSACKIgnoredNoUndo | 当一个dsack block被判定为无效，且未设置undo_marker，则加1 |

### Reorder

当发现了需要更新某条 TCP 流的 reordering 值(乱序值)时，以下计数器可能被使用到。

不过下面四个计数器为互斥关系，最少见的应该是 TCPRenoReorder ，毕竟 sack 已经被广泛部署使用了。


| 名称 | 含义 |
| --- | --- |
| TCPFACKReorder | 如果在需要更新时判断支持 FACK ，则加 1 |
| TCPSACKReorder | 如果仅支持 SACK ，则该计数器加 1 |
| TCPRenoReorder | 如果被不支持 SACK 的 dupack 确认后，需要更新 reorder 值，则加 1 |
| TCPTSReorder | 如果是被一个 partial ack 确认后需要更新 reorder 值，则加 1 |

> 关于 partial ack 的完整内容可参考 [RFC6582](https://tools.ietf.org/html/rfc6582) ；

### Recovery 相关

该类型计数器统计的是进入快速重传阶段的总次数及失败次数；失败次数是指先进入了 recovery 阶段，然后又 RTO 超时了。Fast Recovery 没有成功。

首先由于 SACK 选项已经大面积使用，RenoRecovery 的次数应该远小于 SackRecovery 的次数；另外，fail 的次数应该比例较小才比较理想；

| 名称 | 含义 |
| --- | --- |
| TCPRenoRecovery | 进入 Recovery 阶段的次数，对端不支持 SACK 选项 |
| TCPSackRecovery | 进入 Recovery 阶段的次数，对端支持 SACK 选项 |
| TCPRenoRecoveryFail | 先进入 Recovery 阶段，然后又 RTO 的次数，对端不支持 SACK 选项 |
| TCPSackRecoveryFail | 先进入 Recovery 阶段，然后又 RTO 的次数，对端支持 SACK 选项 |

### Abort

abort 本身是一种很严重的问题，因此有必要关心这些计数器；

后三个计数器如果不为 0 ，则往往意味着系统发生了较为严重的问题，需要格外注意；

| 名称 | 含义 |
| --- | --- |
| TCPAbortOnData | 如果在 FIN_WAIT_1 和 FIN_WAIT_2 状态下收到后续数据，或 TCP_LINGER2 设置小于 0 ，则计数器加 1 |
| TCPAbortOnClose | 如果调用 `tcp_close()` 关闭 socket 时，recv buffer 中还有数据，则加 1 ，此时会主动发送一个 reset 包给对端 |
| TCPAbortOnMemory | 如果 orphan socket 数量或者 `tcp_memory_allocated` 超过上限，则加 1 ；一般值为 0 |
| TCPAbortOnTimeout | 因各种计时器 (RTO/PTO/keepalive) 的重传次数超过上限，而关闭连接时，计数器加 1 |
| TCPAbortOnLinger | `tcp_close()`中，因 tp->linger2 被设置小于 0 ，导致 FIN_WAIT_2 立即切换到 CLOSE 状态的次数；一般值为 0 |
| TCPAbortFailed | 如果在准备发送 reset 时，分配 SKB 或者发送 SKB 失败，则加 1 ；一般值为 0 |

### reset 相关

### 内存 Prune

当 rcv_buf 不足时，可能需要 prune ofo queue ，这种情况就会导致 PruneCalled 计数器增加；但一般都应该通过 collapse 节省内存就可以了，并不需要真正 prune 掉被 SACK 的数据。所以 OfoPruned 和更严重的 RcvPruned 都应该计数为 0 。

| 名称 | 含义 |
| --- | --- |
| PruneCalled |  |
| RcvPruned |  |
| OfoPruned |  |
| TCPMemoryPressures |  |


### PAWS

| 名称 | 含义 |
| --- | --- |
| PAWSPassive |  |
| PAWSActive |  |
| PAWSEstab |  |


### Listen 相关

| 名称 | 含义 |
| --- | --- |
| ListenOverflows |  |
| ListenDrops |  |

### undo 相关

| 名称 | 含义 |
| --- | --- |
| TCPFullUndo |  |
| TCPPartialUndo |  |
| TCPDSACKUndo |  |
| TCPLossUndo |  |

### 快速路径与慢速路径

| 名称 | 含义 |
| --- | --- |
| TCPHPHits |  |
| TCPHPHitsToUser |  |
| TCPPureAcks |  |
| TCPHPAcks |  |


### 未找到

OutOfWindowIcmps 
LockDroppedIcmps 
ArpFilter 


### 未找到

TCPPrequeued 
TCPDirectCopyFromBacklog 
TCPDirectCopyFromPrequeue 
TCPPrequeueDropped 


### 未找到

TCPSACKReneging 


### 未找到

TCPRcvCollapsed 
TCPSACKDiscard 


### 未找到

TCPSackShifted 
TCPSackMerged 
TCPSackShiftFallback 
TCPBacklogDrop 
TCPMinTTLDrop 
TCPDeferAcceptDrop 
IPReversePathFilter 
 
TCPReqQFullDoCookies 
TCPReqQFullDrop 

TCPRcvCoalesce 
TCPOFOQueue 
TCPOFODrop 
TCPOFOMerge 
TCPChallengeACK 
TCPSYNChallenge 



### 未找到

TCPSpuriousRtxHostQueues 
BusyPollRxPackets 
TCPAutoCorking 
TCPFromZeroWindowAdv 
TCPToZeroWindowAdv 
TCPWantZeroWindowAdv 
TCPSynRetrans 
TCPOrigDataSent 
TCPHystartTrainDetect 
TCPHystartTrainCwnd 
TCPHystartDelayDetect 
TCPHystartDelayCwnd


----------


## IpExt

InNoRoutes 
InTruncatedPkts 
InMcastPkts 
OutMcastPkts 
InBcastPkts 
OutBcastPkts 
InOctets 
OutOctets 
InMcastOctets 
OutMcastOctets 
InBcastOctets 
OutBcastOctets 
InCsumErrors 
InNoECTPkts 
InECT1Pkts 
InECT0Pkts 
InCEPkts


----------


## 参考资料

- [RFC 2012: SNMPv2 Management Information Base for the Transmission Control Protocol using SMIv2](https://tools.ietf.org/html/rfc2012)
- [TCP Fast Open: expediting web services](https://lwn.net/Articles/508865/)
- [TCP SNMP counters一](http://blog.chinaunix.net/uid-20043340-id-2984198.html)
- [TCP SNMP counters二](http://blog.chinaunix.net/uid-20043340-id-3016560.html)
- [TCP SNMP counters三](http://blog.chinaunix.net/uid-20043340-id-3017972.html)
- [netstat -s输出解析(一)](http://perthcharles.github.io/2015/11/09/wiki-rfc2012-snmp-proc/)
- [netstat -st输出解析(二)](http://perthcharles.github.io/2015/11/10/wiki-netstat-proc/)


