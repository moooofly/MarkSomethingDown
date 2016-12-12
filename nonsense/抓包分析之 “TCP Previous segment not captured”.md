# 抓包分析之 “TCP Previous segment not captured”

标签（空格分隔）： tcp 丢包 retransmission

---


前两天遇到了这样一个奇怪的现象：向某个服务发起 HTTP 请求，会概率性的出现丢包现象；

抓包截图如下：

![HTTP 请求正常情况](https://github.com/moooofly/ImageCache/blob/master/Pictures/HTTP%20%E8%AF%B7%E6%B1%82%E6%AD%A3%E5%B8%B8%E6%83%85%E5%86%B5.png "HTTP 请求正常情况")

> 在正常情况下，服务器在接收到 HTTP 请求后会迅速进行回应；并且 TCP 连接关闭是客户端侧发起的；

![HTTP 请求异常情况](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/HTTP%20%E8%AF%B7%E6%B1%82%E5%BC%82%E5%B8%B8%E6%83%85%E5%86%B5.png "HTTP 请求异常情况")

> 在异常情况下，能够看到服务器在 20s 后主动关闭了 TCP 连接，并且 wireshark 的专家信息告诉我们，存在没有抓到的数据包（TCP Previous segment not captured）；

需要注意的是：上述两个截图是通过 curl 连续不断进行测试时抓到的，问题复现概率非常高；因为服务器在云端，因此只能在本地客户端侧进行抓包分析；

针对异常情况的数据包，分析如下：

1. 196 号包说明服务器认为自己已经完成了服务请求，但客户端却没有主动关闭连接，于是在 20s 后只好主动将连接关闭；196 号包显示 "TCP Previous segment not captured" 并且包中 "Seq=1011" ，而我们抓到的上一个服务器侧发出包需要为 "Seq=1" ，说明其中有 1010 字节的包我们确实没抓到；
2. 197 号包说明客户端在收到 "Seq=1011" 后认为自己收到了一个乱序包，因此试图通过 "TCP Dup ACK 192#1" 和 "Ack=1" 试图让服务器重传缺失的 1010 子节；
3. 198 号包说明客户端发现 196 包携带了 "FIN" 标志，因此（在过了大约 6s 后）进行了四次握手的连接关闭；
4. 199 号包说明客户端由于某种原因还进行了一次针对 "FIN" 的 "TCP Spurious Retransmission" ；

从上面的分析可以看出，虽然异常包存在好几个，但最为关键的是 "TCP Previous segment not captured" ；因为问题是概率性出现，并且数据包大小也属于常规；另外据相关人员说，服务器侧没有 CPU 等异常，并且网卡流量也在正常范围内；因此，问题看起来就如同间歇式“抽风”一般；那么，一般什么情况会导致 "TCP Previous segment not captured" 的出现呢？

在抓包分析时，可能会看到上述信息情况有：

- 在 tcp session 的中间阶段进行抓包，导致部分信息的缺失（属于正常情况，无法避免）；
- 抓包处理速度无法满足数据包到来的速度（可以通过 capture filter 进行调整）；
- 交换机、路由器和防火墙等在某些情况下会导致上述问题；
- 杀毒软件、恶意软件监测程序等也可能导致上述问题；
- 过于老旧的 TCP 协议栈实现可能存在相关 bug ；


另外，查阅 Wireshark 官网，存在如下解释说明：

> **TCP Previous segment lost** - Occurs when a packet arrives with a sequence number greater than the "next expected sequence number" on that connection, indicating that one or more packets prior to the flagged packet did not arrive. This event is a good indicator of packet loss and will likely be accompanied by "TCP Retransmission" events. 


总之，目前尚无法确定问题出现的原因，但怀疑和防火墙等设备有关，后续有进展后再更新结论～～

----------

> 网上找到的通用分析说明

正如告警信息所描述的那样，这种情况常见于抓包行为起始于一个 tcp session 的中间阶段的情况；如果真实情况确实是由于 acks 丢失所致，那么你需要去相对于你的 host 的 upstream 去查看包是怎么丢失的；

非常可能出现的一种情况是，tshark 的抓包速度无法跟上数据包到来的速度，然后，额，理所当然的会丢掉一些东东（metrics）；在抓包停止的时候，会有相应的信息告诉你是否存在 "**kernel dropped packet**" 的情况，以及 drop 了多少的信息；默认情况下，tshark 会 disable 掉 dns 查询，而 tcpdump 默认则不会；如果你使用 tcpdump 进行抓包，则需要使用 "-n" 开关来禁止 dns 查询以减少额外开销；如果你面临 disk IO 不足问题，那么你可以将抓包内容写到诸如 /dev/shm 等内存区中；但这种用法需要当心，因为一旦你等抓包内容过大，将会导致机器进入 swapping 困境之中；

一种常见的情况是，你的应用使用了一些长时间运行的 tcp session ，在后续启动抓包行为时，就必然会丢失这类 tcp session 的一部分信息；

具体来说，一些比较典型的、我所见过的、会导致 duplicate/missing acks 的情况有：

- **Switches** - 通常情况下不太可能，但有些时候交换机处于 sick 状态时会出现这种问题；
- **Routers** - 比交换机出现问题的可能性要高，但也不是非常高；
- **Firewall** - 比路由器出现问题的可能要高，常见原因和资源耗尽有关（license, cpu, etc）；
- **Client side filtering software** - 反病毒软件，恶意软件探测软件等；

如果在 busy interface 上针对全部网络流量进行 unfiltered 模式的抓包，那么很大概率会在停止 tshark 时看到大量的数据包显示为 'dropped' ，因为这种情况下，tshark 需要针对全部抓包进行排序操作；一旦设置了合适的 capture filter 以减少针对非目标包的捕获，则出现 'dropped' 的概率大大降低；


----------

> 下图是网友给出的一种特殊情况：TCP 协议栈 bug 导致的序列号错误；

![tcp 连接关闭时由于协议栈bug导致的认为数据包丢失](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/tcp%20%E8%BF%9E%E6%8E%A5%E5%85%B3%E9%97%AD%E6%97%B6%E7%94%B1%E4%BA%8E%E5%8D%8F%E8%AE%AE%E6%A0%88bug%E5%AF%BC%E8%87%B4%E7%9A%84%E8%AE%A4%E4%B8%BA%E6%95%B0%E6%8D%AE%E5%8C%85%E4%B8%A2%E5%A4%B1.png "tcp 连接关闭时由于协议栈bug导致的认为数据包丢失")

包地址：[这里](https://www.cloudshark.org/captures/c256982bb42d)

从上述抓包中明显可以看出，来自 server 地址 82.117.201.86 的包明显没有遵循 TCP RFC 进行协议栈实现：在关闭 tcp 会话的时候，FIN 包应该携带下一个 expected 序列包，并且在该 FIN 包被 ACKed 时，相应的 ACK 包需要针对 FIN 所占用的一字节（ one phantom byte）进行确定；然而，上述包的实际情况是，server 发送的 FIN 包中带有的序列号已经跨过了一字节的长度；这种行为明显是错误的，因此客户端此时正确的发送了一个 duplicate ACK 以请求获取正确序列号的包；

因此，导致上述抓包的可能原因，要么确实存在一字节数据被发送，之后丢失的情况，要么就是 server 端的 TCP 协议栈存在 bug ；

"tcp previous segment not captured" 是 Wireshark 软件提供的专家信息，用于告知在当前包捕获中缺少了本应出现的某些包；该告警信息在之前被描述为 "tcp previous segment lost" ；该告警的基本含义为：或者是包真的没有出现过（丢包），或者是 Wireshark 的抓包动作不够快，以致没有将包到来的情况记录下来；


----------

参考文章：

- [Understanding [TCP ACKed unseen segment] [TCP Previous segment not captured]](http://stackoverflow.com/questions/18325522/understanding-tcp-acked-unseen-segment-tcp-previous-segment-not-captured)
- [TCP previous segment not captured, why?](https://ask.wireshark.org/questions/12943/tcp-previous-segment-not-captured-why)


