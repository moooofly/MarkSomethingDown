# Ring Buffer

## [Linux TCP队列相关参数的总结](https://yq.aliyun.com/articles/4252)

**Ring Buffer** 位于 NIC 和 IP 层之间（准确的说位于 NIC driver 中），是一个典型的 FIFO 环形队列。Ring Buffer 没有包含数据本身，而是包含了指向 `sk_buff`（socket kernel buffers）的描述符。
可以使用 `ethtool -g eth0` 查看当前 **Ring Buffer** 的设置：

```
root@vagrant-ubuntu-trusty:~] $ ethtool -g eth0
Ring parameters for eth0:
Pre-set maximums:
RX:		4096
RX Mini:	0
RX Jumbo:	0
TX:		4096
Current hardware settings:
RX:		256
RX Mini:	0
RX Jumbo:	0
TX:		256

root@vagrant-ubuntu-trusty:~] $
```

可以通过 `ifconfig` 观察接收和传输队列的运行状况：

```
eth0      Link encap:Ethernet  HWaddr 08:00:27:4a:c4:2f
          inet addr:10.0.2.15  Bcast:10.0.2.255  Mask:255.255.255.0
          inet6 addr: fe80::a00:27ff:fe4a:c42f/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:2567 errors:0 dropped:0 overruns:0 frame:0
          TX packets:1640 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:208893 (208.8 KB)  TX bytes:198997 (198.9 KB)
```

其中

- RX errors：收包总的错误数；
- `RX dropped`：表示数据包**已经进入了 Ring Buffer** ，但是**由于内存不够等系统原因，导致在拷贝到内存的过程中被丢弃**。
- `RX overruns`：**表示数据包没到 Ring Buffer 就被网卡物理层给丢弃了**，而 CPU 无法及时的处理中断是造成 Ring Buffer 满的原因之一，例如中断分配的不均匀。

当 dropped 数量持续增加，建议增大 Ring Buffer ，使用 `ethtool -G` 进行设置。

### 数据包的接收

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E6%95%B0%E6%8D%AE%E5%8C%85%E6%8E%A5%E6%94%B6%E8%B7%AF%E5%BE%84.jpeg)

从下往上经过了三层：网卡驱动、系统内核空间，最后到用户态空间的应用。Linux 内核使用 sk_buff 数据结构描述一个数据包。当一个新的数据包到达，NIC 调用 DMA engine ，通过 **Ring Buffer** 将数据包放置到内核内存区。**Ring Buffer** 的大小固定，它不包含实际的数据包，而是包含了指向 sk_buff 的描述符。当 **Ring Buffer** 满的时候，新来的数据包将给丢弃。一旦数据包被成功接收，NIC 发起中断，由内核的中断处理程序将数据包传递给 IP 层。经过 IP 层的处理，数据包被放入队列等待 TCP 层处理。每个数据包经过 TCP 层一系列复杂的步骤，更新 TCP 状态机，最终到达 recvBuffer ，等待被应用接收处理。有一点需要注意，数据包到达 recvBuffer ，TCP 就会回 ACK 确认，即 TCP 的 ACK 表示数据包已经被操作系统内核收到，但并不确保应用层一定收到数据（例如这个时候系统 crash），因此一般建议应用协议层也要设计自己的确认机制。

### 数据包的发送

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E6%95%B0%E6%8D%AE%E5%8C%85%E5%8F%91%E9%80%81%E8%B7%AF%E5%BE%84.jpeg)

和接收数据的路径相反，数据包的发送从上往下也经过了三层：用户态空间的应用、系统内核空间、最后到网卡驱动。应用先将数据写入 TCP sendbuffer ，TCP 层将 sendbuffer 中的数据构建成数据包转交给 IP 层。IP 层会将待发送的数据包放入队列 QDisc 。数据包成功放入 QDisc 后，指向数据包的描述符 sk_buff 被放入 **Ring Buffer** 输出队列，随后网卡驱动调用 DMA engine 将数据发送到网络链路上。


----------

## [How can I increase the ring buffer size of my NIC?](https://superuser.com/questions/284677/how-can-i-increase-the-ring-buffer-size-of-my-nic)

The NIC ring buffer maximum size is determined by how much memory is available on the NIC. Typically you do not adjust this setting, this is very much a system administration task and an advanced one at that. 4MB is quite large for a NIC ring buffer. Intel NICs tend to cap at this amount. Broadcom NICs tend to cap at less than one quarter that amount, 1020KB. It is extremely unlikely, unless you have a 10GigE NIC, that you can go above 4096KB in the NIC's internal ring buffer. But we would need the exact model to know for sure as it is a hardware limitation.