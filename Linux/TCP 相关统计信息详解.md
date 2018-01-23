# TCP 统计信息详解

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


------


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
| Reset 相关 | EstabResets |
| 内存 prune | PruneCalled <br> RcvPruned <br> OfoPruned <br> TCPMemoryPressures |
| PAWS 相关 | PAWSPassive <br> PAWSActive <br> PAWSEstab |
| Listen 相关 | ListenOverflows <br> ListenDrops |
| Undo 相关 | TCPFullUndo <br> TCPPartialUndo <br> TCPDSACKUndo <br> TCPLossUndo |
| 快速路径与慢速路径 | TCPHPHits <br> TCPHPHitsToUser <br> TCPPureAcks <br> TCPHPAcks |

### 常量

这些常量是 Linux 3.10 中的默认值，仅在升级了内核版本时才需要关心一下这些值的变化。

| 名称 | 含义 |
| --- | --- |
| RtoAlgorithm | The algorithm used to determine the timeout value used for retransmitting unacknowledged octets. <br><br> 用于计算 RTO 的算法，默认为 1 ，RTO 算法与 RFC2698 一致 |
| RtoMin | The minimum value permitted by a TCP implementation for the retransmission timeout, measured in milliseconds. More refined semantics for objects of this type depend upon the algorithm used to determine the retransmission timeout. In particular, when the timeout algorithm is ``rsre '' (3), an object of this type has the semantics of the LBOUND quantity described in RFC 793. <br><br> 限定 RTO 的最小值，默认值为 1/5HZ，即 **200ms** |
| RtoMax | The maximum value permitted by a TCP implementation for the retransmission timeout, measured in milliseconds. More refined semantics for objects of this type depend upon the algorithm used to determine the retransmission timeout. In particular, when the timeout algorithm is ``rsre'' (3), an object of this type has the semantics of the UBOUND quantity described in RFC 793. <br><br> 限定 RTO 的最大值，默认值为 120HZ，即 **120s** |
| MaxConn | The limit on the total number of TCP connections the entity can support. In entities where the maximum number of connections is dynamic, this object should contain the value -1. <br><br> TCP 连接数量的上限，协议栈本身并不会限制 TCP 连接总数，默认值为 -1 |

### 建链统计

这些统计值中，只有 **CurrEstab** 反应的是系统**当前状态**，而其他值则是反应的**历史状态**；同时需要注意的是，**这些计数器将处于 `ESTABLISHED` 和 `CLOSE-WAIT` 状态的连接都算进当前连接数**。

可以这么理解：这两个状态都认为 `local => remote peer` 方向的连接未被关闭；

| 名称 | 含义 |
| --- | --- |
| ActiveOpens | `<num>` active connections openings<br><br>The number of times TCP connections have made a direct transition to the `SYN-SENT` state from the `CLOSED` state. <br><br> 主动建链次数，对应 `CLOSED` => `SYN-SENT` 次数； <br> 在 `tcp_connect()` 函数中计数； <br> 相当于 SYN 包的发送次数（但不包含重传次数） |
| PassiveOpens | `<num>` passive connection openings<br><br>The number of times TCP connections have made a direct transition to the SYN-RCVD state from the `LISTEN` state. <br><br> 被动建链次数，RFC 原意对应 `LISTEN` => `SYN-RECV` 次数，但 Linux 实现选择在三次握手成功后才加 1 （即在建立 tcp_sock 结构体后） |
| AttemptFails | `<num>` failed connection attempts<br><br>The number of times TCP connections have made a direct transition to the `CLOSED` state from either the `SYN-SENT` state or the SYN-RCVD state, plus the number of times TCP connections have made a direct transition to the `LISTEN` state from the SYN-RCVD state. <br><br> 建链失败次数，即如下三项之和 <br> a) `SYN-SENT` => `CLOSED` 次数 <br> b) `SYN-RECV` => `CLOSED` 次数 <br> c) `SYN-RECV` => `LISTEN` 次数 <br><br> 回 `CLOSED` 部分在 `tcp_done()` 函数中计数 <br> 回 `LISTEN` 部分在 `tcp_check_req()` 中计数 |
| EstabResets | `<num>` connection resets received<br><br>The number of times TCP connections have made a direct transition to the `CLOSED` state from either the `ESTABLISHED` state or the `CLOSE-WAIT` state. <br><br> 连接被 RST 次数，即如下两项之和 <br> a) `ESTABLISHED` => `CLOSED` 次数 <br> b) `CLOSE-WAIT` => `CLOSED` 次 <br><br> 在 `tcp_set_state()` 函数中，如果之前的状态是 TCP_CLOSE_WAIT 或 TCP_ESTABLISHED 就加 1 |
| CurrEstab | `<num>` connections `ESTABLISHED`<br><br>The number of TCP connections for which the current state is either `ESTABLISHED` or `CLOSE-WAIT`. <br><br> 处于 `ESTABLISHED` 和 `CLOSE-WAIT` 状态的 TCP 连接数 <br> 在 `tcp_set_state()` 中进行处理 <br> 实现体现的是进入 `ESTABLISHED` 之后，进入 `CLOSED` 之前的 TCP 连接数 |

### 数据包统计

这些统计值反应的也是历史状态，独立的来看意义并不大。一般可统计一段时间内的变化，关注以下几个指标

-  **（发送）TCP 分段重传占比**：`ΔRetransSegs / ΔOutSegs` ；该值越小越好，如果超过 20% 则应该引起注意（这个值根据实际情况而定）；
-  **（发送）RST 分段占比**：`ΔOutRsts / ΔOutSegs` ；该值越小越好，一般应该在 1% 以内；
-  **（接收）错误分段占比**：`ΔInErrs / ΔInSegs` ；该值越小越好，一般应该在 1% 以内，同时由 checksum 导致的问题包应该更低；

| 名称 | 含义 |
| --- | --- |
| InSegs | `<num>` segments received<br><br>The total number of segments received, including those received in error. This count includes segments received on currently `ESTABLISHED` connections. <br><br> 所有收到的 TCP 分段，即使是个错误分段 <br><br> 在 `tcp_v4_rcv()` 和 `tcp_v6_rcv()` 中计数 |
| OutSegs | `<num>` segments send out<br><br>The total number of segments sent, including those on current connections but **excluding those containing only retransmitted octets**. <br><br> 所有发送出去的 TCP 分段，包括 <br><br> a) 新数据包 <br> b) 重传数据包 <br> c) SYN 包 <br> d) SYN,ACK 包 <br> e) RST 包 <br><br> 不包括那些只包含重传字节的分段 <br><br> `tcp_v4_send_reset()` 中统计 RST 包 <br> `tcp_v4_send_ack()` 中统计 `SYN-RECV` 和 `TIME-WAIT` 状态下发送的 ACK 包 <br> `tcp_v6_send_response()` 中统计 ipv6 相应数据 <br> `tcp_make_synack()` 中统计发送的 SYN,ACK 包 <br> `tcp_transmit_skb()` 中统计所有的其他包 |
| RetransSegs | `<num>` segments retransmited<br><br>The total number of segments retransmitted - that is, the number of TCP segments transmitted containing one or more previously transmitted octets. <br><br> 所有重传出去的 TCP 分段 <br><br> `tcp_v4_rtx_synack()` 和 `tcp_v6_rtx_synack()` 中统计重传的 SYN,ACK 包 <br> `tcp_retransmit_skb()` 中统计其他重传包 |
| InErrs | `<num>` bad segments received<br><br>The total number of segments received in error (for example, bad TCP checksums). <br><br> 所有收到的有问题的 TCP 分段数量，比如 checksum 有问题 <br><br> `tcp_validate_incoming()` 中统计 seq 有问题的包 <br> `tcp_rcv_established()`、`tcp_v4_do_rcv()`、`tcp_v4_rcv()`、`tcp_v6_do_rcv()`、`tcp_v6_rcv()` 中根据 checksum 来判断出错误分段 |
| OutRsts | `<num>` resets sent<br><br> The number of TCP segments sent containing the RST flag. <br><br> 发送的带 RST 标记的 TCP 分段数量 <br><br> 在 `tcp_v4_send_reset()`、`tcp_send_active_reset()`、`tcp_v6_send_response()` 中统计 |
| InCsumErrors | 收到的 checksum 有问题的数据包数量 <br><br> 属于 3.10 相对于 2.6.32 新增的内容，算是细化 InErrs 统计，InErrs 中应该只有*小部分*属于该类型 |
| EmbryonicRsts | number of RSTs received for embryonic SYN_RECV sockets <br> 在 `SYN-RECV` 状态收到带 RST/SYN 标记的包个数 |

### Syncookies 相关

syncookies 一般不会被触发，只有在 `tcp_max_syn_backlog` 队列被占满时才会被触发；因此 SyncookiesSent 和 SyncookiesRecv 一般情况下应该是 0 。但是 SyncookiesFailed 的值和 syncookies 机制是否被触发没有直接关系，因此可能不为 0 ，原因在于：SyncookiesFailed 值的计算方式为：当一个处于 `LISTEN` 状态的 socket 收到一个不带 SYN 标记的数据包时，就会调用 `cookie_v4_check()` 尝试验证 cookie 信息。而如果验证失败，值加 1 。

| 名称 | 含义 |
| --- | --- |
| SyncookiesSent | SYN cookies sent <br><br> An application wasn't able to accept a connection fast enough, so the kernel couldn't store an entry in the queue for this connection. Instead of dropping it, it sent a cookie to the client <br><br> 使用 syncookie 技术发送的 syn/ack 包个数 |
| SyncookiesRecv | SYN cookies received <br><br> After sending a cookie, it came back to us and passed the check. <br><br> 收到携带有效 syncookie 信息包个数 |
| SyncookiesFailed | Num of invalid SYN cookies received <br><br> After sending a cookie, it came back to us but looked invalid <br><br> 收到携带无效 syncookie 信息包个数 |

> 注：syncookies 机制是为应对 syn flood 攻击而被提出来的。

### TIME-WAIT 相关

`TIME-WAIT` 状态是 TCP 协议状态机中的重要一环，服务器设备一般都有非常多处于 `TIME-WAIT` 状态的 socket ，如果是在主要提供 HTTP 服务的设备上，TW 值应该接近 TcpPassiveOpens 值。

一般情况下，`sysctl_tcp_tw_reuse` 和 `sysctl_tcp_tw_recycle` 都是不推荐开启的。所以 TWKilled 和 TWRecycled 都应该是 0 。同时 TCPTimeWaitOverflow 也应该是 0 ，否则就意味着内存使用方面出了大问题。


| 名称 | 含义 |
| --- | --- |
| TW | number of TCP sockets finished time wait in **fast** timer <br> 经过正常时间（`TCP_TIMEWAIT_LEN`）结束 TW 状态的 socket 数量 |
| TWRecycled | number of time wait sockets recycled by time stamp <br> `TIME-WAIT` socket 被复用的次数；只有在 `sysctl_tcp_tw_reuse` 开启时，才可能加 1 |
| TWKilled | number of TCP sockets finished time wait in **slow** timer <br> 经过更短时间结束 TW 状态的 socket 数量；只有在 `net.ipv4.tcp_tw_recycle` 开启时，调度 TW timer 时才可能用更短的 timeout 值 |
| TCPTimeWaitOverflow | 如果没有内存分配 TIMEWAIT 结构体，则加 1 |

### RTO 相关

RTO 超时对 TCP 性能的影响是巨大的，因此关心 RTO 超时的次数也非常必要。

当然 3.10 中的 TLP 机制能够减少一定量的 TCPTimeouts 数，将其转换为快速重传。


| 名称 | 含义 |
| --- | --- |
| TCPTimeouts | a) 在 RTO timer 中，从 CWR/Open 状态下第一次超时的次数，其余状态不计入这个计数器；<br> b) SYN,ACK 的超时次数 |
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
| TCPSlowStartRetrans | 成功在 Loss 状态发送的、重传 SKB 数量，而且这里仅记录非 RTO 超时进入 Loss 状态下的重传数量；目前找到的一种非 RTO 进入 Loss 状态的情况就是：`tcp_check_sack_reneging()` 函数发现接收端违反(renege)了之前的 SACK 信息时，会进入 Loss 状态 |
| TCPRetransFail | 尝试 FastRetrans、ForwardRetrans、SlowStartRetrans 重传失败的次数 |

### FastOpen

**TCP FastOpen (`TFO`)** 技术是 Google 提出来减少三次握手开销的技术，核心原理就是在第一次建链时，由 server 计算出一个 cookies 发给 client ，之后 client 向 server 再次发起建链请求时，就可以携带该 cookies 信息以验明正身。如果 cookies 验证通过，则 server 可以不等三次握手的最后一个 ACK 包，就将 client 放在 SYN 包里面的数据传递给应用层。

在 3.10 内核中，`TFO` 由 `sysctl_tcp_fastopen` 开关控制，默认值为 0(关闭)。而且 `sysctl_tcp_fastopen` 目前也是推荐关闭的，因为网络中有些 middle box 会丢弃那些带有不认识 option 的 SYN 包；所以正常情况下，这些值也应该都是 0 ，当然如果收到过某些不怀好意的、带 TFO cookies 信息的 SYN 包，TCPFastOpenPassive 计数器就可能不为 0 。

| 名称 | 含义 |
| --- | --- |
| TCPFastOpenActive | number of successful outbound TFO connections <br><br> 主动发送的、带 TFO cookie 的 SYN 包个数 |
| TCPFastOpenActiveFail | number of SYN,ACK packets received that did not acknowledge data sent in the SYN packet and caused a retransmissions without SYN data. Note that the original SYN packet contained a cookie + data, this is not the number of connections to servers that didn’t support TFO <br><br> 基于 TFO 主动建链失败的次数 |
| TCPFastOpenPassive | number of successful inbound TFO connections <br><br> 收到带 TFO cookie 的 SYN 包个数 |
| TCPFastOpenPassiveFail | number of inbound SYN packets with TFO cookie that was invalid <br><br> 基于 TFO 被动建链，但由于 cookie 无效而失败的次数 |
| TCPFastOpenListenOverflow | number of inbound SYN packets that will have TFO disabled because the socket has exceeded the max queue length <br><br> TFO 请求数超过监听队列设置上限，则加 1 |
| TCPFastOpenCookieReqd | number of inbound SYN packets requesting TFO with TFO set but no cookie <br><br> 收到一个请求 TFO cookies 的 SYN 包时，则加 1 |

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
| DelayedACKs | number of delayed acks sent <br><br> We waited for another packet to send an ACK, but didn't see any, so a timer ended up sending a delayed ACK. <br><br> 调用 tcp_send_ack() 的次数，无论发送是否成功 <br><br> 触发点：tcp_delack_timer() |
| DelayedACKLocked | number of delayed acks further delayed because of locked socket <br><br> We wanted to send a delayed ACK but failed because the socket was locked. So the timer was reset. <br><br> delay ACK 定时器因为 user 已经锁住而无法发送 ACK 的次数 <br><br> 触发点：tcp_delack_timer() |
| DelayedACKLost | Quick ack mode was activated %u times <br><br> We sent a delayed and duplicated ACK because the remote peer retransmitted a packet, thinking that it didn't get to us. <br><br> a) 当输入包不在接收窗口内，或者 PAWS 失败后，计数器加 1 ；触发点：tcp_validate_incoming()->tcp_send_dupack() <br> b) 输入包的结束序列号 < RCV_NXT 时，加 1 ；触发点：tcp_data_queue() |
| TCPSchedulerFailed | 在 delay ACK 处理功能内，如果 prequeue 中仍有数据，计数器就加 1 <br> 加入到 prequeue ，本来是期待着 userspace（使用 tcp_recvmsg() 之类的系统调用）尽快处理之。若其中仍有数据，则可能隐含着 userspace 行为不佳 <br><br> 触发点：tcp_delack_timer() |

### DSACK

该类型计数器统计的是收/发 DSACK 信息次数。

DSACKOldSent + DSACKOfoSent 可以当做是发送出的 DSACK 信息的次数，而且概率上来讲 OldSent 应该占比更大。

同理，DSACKRecv 的数量也应该远多于 DSACKOfoRecv 的数量。

另外，DSACK 信息的发送是需要 `sysctl_tcp_dsack` 开启的，如果发现 sent 两个计数器为零，则要检查一下了。

一般还是建议开启 dsack 选项；

| 名称 | 含义 |
| --- | --- |
| TCPDSACKOldSent | 如果收到的重复数据包序号比 rcv_nxt（接收端想收到的下一个序号）小，则增加 oldsent |
| TCPDSACKOfoSent | 如果收到的重复数据包序号比 rcv_nxt 大，则是一个乱序的重复数据包，增加 ofosent |
| TCPDSACKRecv | 收到的 old dsack 信息次数，判断 old 的方法：dsack 序号小于 ACK 号 |
| TCPDSACKOfoRecv | 收到的 Ofo dsack 信息次数 |
| TCPDSACKIgnoredOld | We got a duplicate SACK while retransmitting so we discarded it. <br><br> 当一个 dsack block 被判定为无效，且设置过 undo_marker ，则加 1 |
| TCPDSACKIgnoredNoUndo | We got a duplicate SACK and discarded it. <br><br> 当一个 dsack block 被判定为无效，且未设置 undo_marker ，则加 1 |

### Reorder

当发现了需要更新某条 TCP 连接的 reordering 值(乱序值)时，以下计数器可能被使用到。

不过下面四个计数器为互斥关系，最少见的应该是 TCPRenoReorder ，毕竟 SACK 已经被广泛部署使用了。

| 名称 | 含义 |
| --- | --- |
| TCPFACKReorder | We detected re-ordering using FACK <br> Forward ACK, the highest sequence number known to have been received by the peer when using SACK. FACK is used during congestion control |
| TCPSACKReorder | We detected re-ordering using SACK |
| TCPRenoReorder | We detected re-ordering using fast retransmit |
| TCPTSReorder | We detected re-ordering using the timestamp option |

> 关于 partial ack 的完整内容可参考 [RFC6582](https://tools.ietf.org/html/rfc6582) ；

### Recovery 相关

该类型计数器统计的是进入快速重传阶段的总次数及失败次数；失败次数是指先进入了 recovery 阶段，然后又 RTO 超时了。Fast Recovery 没有成功。

首先由于 SACK 选项已经大面积使用，RenoRecovery 的次数应该远小于 SackRecovery 的次数；另外，fail 的次数应该比例较小才比较理想；

| 名称 | 含义 |
| --- | --- |
| TCPRenoRecovery | A packet was lost and we recovered after a fast retransmit |
| TCPSackRecovery | A packet was lost and we recovered by using selective acknowledgements |
| TCPRenoRecoveryFail | 先进入 Recovery 阶段，然后又 RTO 的次数，对端不支持 SACK 选项 |
| TCPSackRecoveryFail | 先进入 Recovery 阶段，然后又 RTO 的次数，对端支持 SACK 选项 |

### Abort

abort 本身是一种很严重的问题，因此有必要关心这些计数器；

后三个计数器如果不为 0 ，则往往意味着系统发生了较为严重的问题，需要格外注意；

| 名称 | 含义 |
| --- | --- |
| TCPAbortOnSyn | We received an unexpected SYN so we sent a RST to the peer |
| TCPAbortOnData | We were in FIN_WAIT1 yet we received a data packet with a sequence number that's beyond the last one for this connection, so we RST'ed. <br><br> 如果在 FIN_WAIT_1 和 FIN_WAIT_2 状态下收到后续数据，或 TCP_LINGER2 设置小于 0 ，则计数器加 1 |
| TCPAbortOnClose | We received data but the user has `CLOSED` the socket, so we have no wait of handing it to them, so we RST'ed. <br><br> 如果调用 `tcp_close()` 关闭 socket 时，recv buffer 中还有数据，则加 1 ，此时会主动发送一个 RST 包给对端 |
| TCPAbortOnMemory | This is Really Bad. It happens when there are too many orphaned sockets (not attached a FD) and the kernel has to drop a connection. Sometimes it will send a RST to the peer, sometimes it wont. <br><br> 如果 orphan socket 数量或者 `tcp_memory_allocated` 超过上限，则加 1 ；一般值为 0 |
| TCPAbortOnTimeout | The connection timed out really hard. <br><br> 因各种计时器 (RTO/PTO/keepalive) 的重传次数超过上限，而关闭连接时，计数器加 1 |
| TCPAbortOnLinger | We killed a socket that was `CLOSED` by the application and lingered around for long enough. <br><br> `tcp_close()`中，因 tp->linger2 被设置小于 0 ，导致 FIN_WAIT_2 立即切换到 `CLOSED` 状态的次数；一般值为 0 |
| TCPAbortFailed | We tried to send a RST, probably during one of the TCPABort* situations above, but we failed e.g. because we couldn't allocate enough memory (very bad). <br><br> 如果在准备发送 RST 时，分配 SKB 或者发送 SKB 失败，则加 1 ；一般值为 0 |

### Reset 相关

| 名称 | 含义 |
| --- | --- |
| EstabResets | 连接被 RST 次数，即如下两项之和 <br><br> a) `ESTABLISHED` => `CLOSED` 次数 <br> b) `CLOSE-WAIT` => `CLOSED` 次 <br><br> 在 `tcp_set_state()` 函数中，如果之前的状态是 TCP_CLOSE_WAIT 或 TCP_ESTABLISHED 就加 1 |

### 内存 Prune

当 rcv_buf 不足时，可能需要 prune ofo queue ，这种情况就会导致 PruneCalled 计数器增加；但一般都应该通过 collapse 节省内存就可以了，并不需要真正 prune 掉被 SACK 的数据。所以 OfoPruned 和更严重的 RcvPruned 都应该计数为 0 。

| 名称 | 含义 |
| --- | --- |
| PruneCalled | packets **pruned** from **receive queue** because of socket buffer overrun <br> 慢速路径中，如果不能将数据直接复制到 user space ，需要加入到 sk_receive_queue 前，会检查 receiver side memory 是否允许，如果 rcv_buf 不足就可能 prune ofo queue 。此时计数器加 1 |
| RcvPruned | _obsolete: 2.2.0 doesn't do that anymore_ <br> packets **pruned** from **receive queue** <br><br> If the kernel is really really desperate and cannot give more memory to this socket even after dropping the ofo queue, it will simply discard the packet it received. This is Really Bad. <br><br> 慢速路径中，如果不能将数据直接复制到 user space ，需要加入到 sk_receive_queue 前，会检查 receiver side memory 是否允许，如果 rcv_buf 不足就可能 prune receive queue ，如果 prune 失败了，此计数器加 1 |
| OfoPruned | packets **dropped** from **out-of-order queue** because of socket buffer overrun <br><br> When a socket is using too much memory (rmem), the kernel will first discard any out-of-order packet that has been queued (with SACK). <br><br> 慢速路径中，如果不能将数据直接复制到 user space ，需要加入到 sk_receive_queue 前，会检查 receiver side memory 是否允许，如果 rcv_buf 不足就可能 prune ofo queue 。此时计数器加 1 |
| TCPMemoryPressures | Number of times a socket was put in "memory pressure" due to a non fatal memory allocation failure (reduces the send buffer size etc). <br><br> tcp_enter_memory_pressure() 在从“非压力状态”切换到“有压力状态”时计数器加 1 ；<br><br> 触发点：<br> a) tcp_sendmsg() <br> b) tcp_sendpage() <br> c) tcp_fragment() <br> d) tso_fragment() <br> e) tcp_mtu_probe() <br> f) tcp_data_queue() |


### PAWS

| 名称 | 含义 |
| --- | --- |
| PAWSPassive | number of **passive** connections rejected because of time stamp <br> 三路握手最后一个 ACK 的 PAWS 检查失败次数 <br><br> 触发点：tcp_v4_conn_request() |
| PAWSActive | number of **active** connections rejected because of time stamp <br> 在发送 SYN 后，接收到 ACK ，但 PAWS 检查失败的次数 <br><br> 触发点：tcp_rcv_synsent_state_process() |
| DelayedACKLocked | number of packets rejects in `ESTABLISHED` connections because of timestamp <br> 输入包 PAWS 失败次数 <br><br> 触发点： <br> a) tcp_validate_incoming() <br> b) tcp_timewait_state_process() <br> c) tcp_check_req() |


### Listen 相关

| 名称 | 含义 |
| --- | --- |
| ListenOverflows | `<num>` times the `LISTEN` queue of a socket overflowed <br><br> We completed a 3WHS but couldn't put the socket on the accept queue, so we had to discard the connection. <br><br> 三路握手最后一步完全之后，Accept queue 队列超过上限时加 1 <br><br> 触发点：tcp_v4_syn_recv_sock() |
| ListenDrops | `<num>` of SYNs to `LISTEN` sockets dropped <br><br> We couldn't accept a connection because one of: we had no route to the destination, we failed to allocate a socket, we failed to allocate a new local port bind bucket. Note: this counter also include all the increments made to ListenOverflows <br><br> 任何原因导致的失败后加 1，包括：Accept queue 超限，创建新连接，继承端口失败等 <br><br> 触发点：tcp_v4_syn_recv_sock() |

### undo 相关

| 名称 | 含义 |
| --- | --- |
| TCPFullUndo | We detected some erroneous retransmits and undid our CWND reduction <br><br> Recovery 状态时，接收到全部的确认（snd_una >= high_seq）后且已经 undo 完成（undo_retrans == 0）的次数 <br><br> tcp_ack() -> tcp_fastretrans_alert() -> tcp_try_undo_recovery() |
| TCPPartialUndo | We detected some erroneous retransmits, a partial ACK arrived while we were fast retransmitting, so we were able to partially undo some of our CWND reduction <br><br> Recovery 状态时，接收到到部分确认（snd_una < high_seq）时但已经 undo 完成（undo_retrans == 0）的次数 <br><br> tcp_ack() -> tcp_fastretrans_alert() -> tcp_undo_partial() |
| TCPDSACKUndo | We detected some erroneous retransmits, a D-SACK arrived and ACK'ed all the retransmitted data, so we undid our CWND reduction <br><br> Disorder 状态下，undo 完成（undo_retrans == 0）的次数 <br><br> tcp_ack() -> tcp_fastretrans_alert() -> tcp_try_undo_dsack() |
| TCPLossUndo | We detected some erroneous retransmits, a partial ACK arrived, so we undid our CWND reduction <br><br> Loss 状态时，接收到到全部确认（snd_una >= high_seq）后且已经 undo 完成（undo_retrans == 0）的次数 <br><br> tcp_ack() -> tcp_fastretrans_alert() -> tcp_try_undo_loss() |

### 快速路径与慢速路径

| 名称 | 含义 |
| --- | --- |
| TCPHPHits | 如果有 skb 通过“快速路径”进入到 sk_receive_queue 上，计数器加 1 ；特别地，Pure ACK 以及直接复制到 user space 上的都不算在这个计数器上 <br><br> 触发点：tcp_rcv_established() |
| TCPHPHitsToUser | 如果有 skb 通过“快速路径”直接复制到 user space 上，计数器加 1 <br><br> 触发点：tcp_rcv_established() |
| TCPPureAcks | 接收“慢速路径”中的 pure ACK 数量 <br><br> 触发点：tcp_ack() |
| TCPHPAcks | 接收到包，进入“快速路径”时加 1 <br><br> 触发点：tcp_ack() |


### SACK

| 名称 | 含义 |
| --- | --- |
| TCPSACKReneging | 如果 snd_una（输入 skb->ack）之后的具有最小开始序号 skb（即 sk_write_queue 中的第一个 skb）中有 TCPCB_SACKED_ACKED 标志，此时加 1 ，这说明接收者已经丢掉了之前它已经 SACK 过的数据 <br><br> 触发点：tcp_clean_rtx_queue() |
| TCPSACKDiscard | We got a completely invalid SACK block and discarded it. <br><br> 非法 SACK 块（不包括 D-SACK）计数，即 SACK 中的序号太旧 <br><br> 触发点：tcp_sacktag_write_queue() |
| TCPSackShifted | 在 tcp_sacktag_walk() 时，一个 SACK 可能会导致切割某 skb ，新切出来的 skb 放到被切的 skb 之后。根据 SACK 的观点，如果“旧的 skb”（变小了）能够与它之前的 skb 合并，本计数器，就加 1 。这个合并过程，叫作 shift <br><br> tcp_ack()->tcp_sacktag_write_queue()->tcp_sacktag_walk()->tcp_shift_skb_data()->tcp_shifted_data() |
| TCPSackMerged | 在上面介绍的 shift 过程中，如果发现分割之后的 skb 被它之前的 skb 完全“吃掉”，本计数器加 1 <br><br> tcp_ack()->tcp_sacktag_write_queue()->tcp_sacktag_walk()->tcp_shift_skb_data()->tcp_shifted_data() |
| TCPSackShiftFallback | 与上相反，如果不能 shift ，本计数器加 1 。原因可能如下：<br> a) 不支持GSO <br> b) prev skb 不完全是 paged 的 <br> c) SACK 的序号已经 ACK 过 <br> d) 等等 <br><br> tcp_ack()->tcp_sacktag_write_queue()->tcp_sacktag_walk()->tcp_shift_skb_data() |


### TCP Others

| 名称 | 含义 |
| --- | --- |
| OutOfWindowIcmps | number of ICMP packets dropped because they were **out-of-window** <br> 接收到 ICMP ，但由于 ICMP 中的 TCP 头序号不在接收窗口之内而导致被丢弃的数量（待确认） <br><br> 有两个可能情况：<br> 1) `LISTEN` 状态时，序号不等于 ISN ；<br> 2) 其他状态时，序号不在 SND_UNA..SND_NXT 之间 <br><br> 触发点：tcp_v4_err() |
| LockDroppedIcmps | number of ICMP packets dropped because **socket was locked** <br> 接收到 ICMP 报文，但由于 socket 被 user 锁住的原因而丢弃的数量 <br><br> 触发点：tcp_v4_err() |
| ArpFilter | 与 TCP 无关，接收到 ARP packet 时做一次输出路由查找（sip, tip），如果找到的路由项的 device 与输入 device 的不同，计数器加 1 <br><br> ArpFilter    arp_rcv() -> NETFILTER(ARP_IN) -> arp_process() |
| TCPPrequeued | tcp_recvmsg() 发现可以从 prequeue 接收到报文，计数器加 1（不是每个 skb 加 1）<br><br>    tcp_recvmsg() -> tcp_prequeue_process() |
| TCPDirectCopyFromBacklog | 如果有数据在 softirq 里面直接从 backlog queue 中复制到 userland memory 上，则计数器加 1 <br><br> 触发点：tcp_recvmsg() |
| TCPDirectCopyFromPrequeue | 如果有数据在这个 syscall 里直接从 prequeue 中复制到 userland memory 上，计数器加 1 <br><br> 触发点：tcp_recvmsg() |
| TCPPrequeueDropped | 如果因为内存不足（ucopy.memory < sk->rcv_buf）而加入到 prequeue 失败，重新由 backlog 处理，计数器加 1 <br><br> tcp_v4_rcv() -> tcp_prequeue() |
| TCPRcvCollapsed | 每当合并 sk_receive_queue(ofo_queue) 中的连续报文时，计数器加 1 <br><br> 触发点：<br> a) tcp_prune_queue() -> tcp_collapse() -> tcp_collapse_one() <br> b) tcp_prune_ofo_queue() -> tcp_collapse()  |
| TCPBacklogDrop | We received something but had to drop it because the socket's receive queue was full. <br><br> 如果 socket 被 user 锁住，后退一步，内核会把包加到 sk_backlog_queue ，但如果因为 sk_rcv_buf 不足的原因入队失败，计数器加 1 <br><br> tcp_v4_rcv() |
| TCPMinTTLDrop | 在接收到 TCP 报文或者 TCP 相关的 ICMP 报文时，检查 IP TTL ，如果小于 socket option 设置的一个阀值，就丢包。这个功能是 RFC5082 (The Generalized TTL Security Mechanism, GTSM) 规定的，使用 GTSM 的通信双方，都将 TTL 设置成最大值 255 ，双方假定了解之间的链路情况，这样可以通过检查最小 TTL 值隔离攻击 <br><br> tcp_v4_err() / tcp_v4_rcv() |
| TCPDeferAcceptDrop | 如果启用 TCP_DEFER_ACCEPT ，这个计数器统计被丢掉的“Pure ACK”的个数。TCP_DEFER_ACCEPT 允许 listener 只有在连接上有数据时才创建新的 socket ，以抵御 syn-flood 攻击 <br><br> tcp_check_req() |
| IPReversePathFilter | 反向路径过滤掉的 IP 分组数量：要么反向路由查找失败，要么是找到的输出接口与输入接口不同 <br><br> ip_rcv_finish() -> ip_route_input_noref() |
| TCPReqQFullDoCookies | syn_table 过载，进行 SYN cookie 的次数（取决于是否打开 `sysctl_tcp_syncookies` ）<br><br> tcp_rcv_state_process() -> tcp_v4_conn_request() -> tcp_syn_flood_action() |
| TCPReqQFullDrop | syn_table 过载，丢掉 SYN 的次数 <br><br> tcp_rcv_state_process() -> tcp_v4_conn_request() -> tcp_syn_flood_action() |


### 未找到

| 名称 | 含义 |
| --- | --- |
| TCPRcvCoalesce |  |
| TCPOFOQueue |  |
| TCPOFODrop |  |
| TCPOFOMerge |  |
| TCPChallengeACK |  |
| TCPSYNChallenge |  |
| TCPSpuriousRtxHostQueues |  |
| BusyPollRxPackets |  |
| TCPAutoCorking |  |
| TCPFromZeroWindowAdv |  |
| TCPToZeroWindowAdv |  |
| TCPWantZeroWindowAdv |  |
| TCPSynRetrans | number of SYN and SYN/ACK retransmits to break down retransmissions into SYN, fast-retransmits, timeout retransmits, etc. |
| TCPOrigDataSent | number of outgoing packets with original data (excluding retransmission but including data-in-SYN). This counter is different from TcpOutSegs because TcpOutSegs also tracks pure ACKs. TCPOrigDataSent is more useful to track the TCP retransmission rate. |
| TCPHystartTrainDetect |  |
| TCPHystartTrainCwnd |  |
| TCPHystartDelayDetect |  |
| TCPHystartDelayCwnd |  |


----------


## IpExt

> TODO

| 名称 | 含义 |
| --- | --- | 
| InNoRoutes |  |
| InTruncatedPkts |  |
| InMcastPkts |  |
| OutMcastPkts |  |
| InBcastPkts |  |
| OutBcastPkts |  |
| InOctets |  |
| OutOctets |  |
| InMcastOctets |  |
| OutMcastOctets |  |
| InBcastOctets |  |
| OutBcastOctets |  |
| InCsumErrors |  |
| InNoECTPkts |  |
| InECT1Pkts |  |
| InECT0Pkts |  |
| InCEPkts |  |


----------


## 参考资料

- [RFC 2012: SNMPv2 Management Information Base for the Transmission Control Protocol using SMIv2](https://tools.ietf.org/html/rfc2012)
- [TCP Fast Open: expediting web services](https://lwn.net/Articles/508865/)
- [TCP SNMP counters一](http://blog.chinaunix.net/uid-20043340-id-2984198.html)
- [TCP SNMP counters二](http://blog.chinaunix.net/uid-20043340-id-3016560.html)
- [TCP SNMP counters三](http://blog.chinaunix.net/uid-20043340-id-3017972.html)
- [netstat -s输出解析(一)](http://perthcharles.github.io/2015/11/09/wiki-rfc2012-snmp-proc/)
- [netstat -st输出解析(二)](http://perthcharles.github.io/2015/11/10/wiki-netstat-proc/)
- [How To Read `netstat -s` Output](http://cromwell-intl.com/networking/netstat-s.html)
- [net-tools](https://sourceforge.net/p/net-tools/code/ci/v1.60/tree/statistics.c) 源码
- [Investigating Linux Network Issues with netstat and nstat](https://perfchron.com/2015/12/26/investigating-linux-network-issues-with-netstat-and-nstat/)
- [TCPCollector](https://github.com/BrightcoveOS/Diamond/wiki/collectors-TCPCollector)


