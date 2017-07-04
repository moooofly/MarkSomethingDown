# Linux TCP 队列相关参数的总结

> 原文地址：[这里](https://yq.aliyun.com/articles/4252)

## 可能丢包原因

针对**收包**：

- 当一个新的数据包到达，NIC 调用 DMA engine 通过 Ring Buffer 将数据包放置到内核内存区（这里描述为“通过”似乎不太好）。Ring Buffer 的大小固定，它不包含实际的数据包，而是包含了指向 sk_buff 的描述符。当 Ring Buffer 满的时候，新来的数据包将给丢弃。可以通过 `ifconfig` 命令的 **RX dropped** 数值进行查看；
- 数据包没到 Ring Buffer 就被网卡物理层给丢弃了，而 CPU 无法及时的处理中断是造成 Ring Buffer 满的原因之一，例如中断分配的不均匀。可以通过 `ifconfig` 的 **RX overruns** 数值进行查看；
- netdev_max_backlog 是内核从 NIC 收包后，交由协议栈（如IP、TCP）处理之前的缓冲队列。每个 CPU 核都有一个 backlog 队列，与 Ring Buffer 同理，当接收包的速率大于内核协议栈处理的速率时，CPU 的 backlog 队列不断增长，当达到设定的 netdev_max_backlog 值时，数据包将被丢弃。

针对**发包**：

- sendBuffer：取决于三个内核参数配置情况（详见下面说明）；
- QDisc：QDisc（queueing discipline ）位于 IP 层和网卡的 Ring Buffer 之间。Ring Buffer 是一个简单的 FIFO 队列，这种设计使网卡的驱动层保持简单和快速。而 QDisc 实现了流量管理的高级功能，包括流量分类，优先级和流量整形。可以使用 tc 命令配置 QDisc 。QDisc 的队列长度由 txqueuelen 设置，和接收数据包的队列长度由内核参数 `net.core.netdev_max_backlog` 控制所不同，txqueuelen 是和网卡关联的，可以通过 `ifconfig` 查看和调整大小；
- RingBuffer：`ethtool -g eth0` 输出中的 TX 项就是 RingBuffer 的传输队列大小；

> 

## recvBuffer 计算

BDP 的含义是任意时刻处于途中尚未确认的最大数据量。为了达到最大的吞吐量，recvBuffer 的设置应该大于 BDP ，即 `recvBuffer >= bandwidth * RTT`。
Linux 在 2.6.17 以后增加了 recvBuffer 自动调节机制，recvbuffer 的实际大小会自动在最小值和最大值之间浮动，以期找到性能和资源的平衡点，因此**大多数情况下不建议将 recvbuffer 手工设置成固定值**。

- 当 `net.ipv4.tcp_moderate_rcvbuf` 设置为 1 时，缓冲的自动调节机制生效，随后 recvbuffer 根据实际情况在最大值和最小值之间动态调节；在缓冲的动态调优机制开启的情况下，建议将 `net.ipv4.tcp_rmem` 的最大值设置为 BDP ；
- 当 `net.ipv4.tcp_moderate_rcvbuf` 被设置为 0 ，或者设置了 socket 选项 `SO_RCVBUF` 时，缓冲的动态调节机制被关闭。在缓冲动态调节机制关闭的情况下，建议把 `net.ipv4.tcp_rmem` 的缺省值设置为 BDP ；

每个 TCP 连接的 recvBuffer 由下面的 3 元数组指定：

- `net.core.rmem_default`
- `net.ipv4.tcp_rmem`
- `net.core.rmem_max`

recvbuffer 的缺省值由 `net.core.rmem_default` 设置，但如果设置了 `net.ipv4.tcp_rmem` ，该缺省值则被覆盖。

```
root@vagrant-ubuntu-trusty:~] $ sysctl -a|grep rmem
net.core.rmem_default = 212992
net.core.rmem_max = 212992
net.ipv4.tcp_rmem = 4096	87380	3827616
root@vagrant-ubuntu-trusty:~] $
```

- **`net.core.rmem_default`**

This sets the **default** OS receive buffer size for **all types of connections**.

- **`net.core.rmem_max`**

This sets the **max** OS receive buffer size for **all types of connections**.

- **`net.ipv4.tcp_rmem`**

    TCP Autotuning setting. 
    
    - The **first** value tells the kernel the **minimum** receive buffer for each TCP connection, and this buffer is always allocated to a TCP socket, even under high pressure on the system.
    - The **second** value specified tells the kernel the **default** receive buffer allocated for each TCP socket. This value overrides the `/proc/sys/net/core/rmem_default` value used by other protocols.
    - The **third** and last value specified in this variable specifies the **maximum** receive buffer that can be allocated for a TCP socket.

注意：这里还有一个细节，上述缓冲除了会保存接收的数据本身，还需要一部分空间保存 socket 数据结构等额外信息。因此上面讨论的 **recvbuffer 最佳值仅仅等于 BDP 是不够的**，还需要考虑保存 socket 等额外信息的开销。Linux 根据参数 `net.ipv4.tcp_adv_win_scale` 计算额外开销的大小：

```
recvBuffer/2^tcp_adv_win_scale
```

如果 `net.ipv4.tcp_adv_win_scale` 的值为 1 ，则二分之一的缓冲空间用来做额外开销，如果为 2 的话，则四分之一缓冲空间用来做额外开销。因此 recvbuffer 的最佳值应该设置为：

```
recvBuffer/(1 - 1/2^tcp_adv_win_scale)
```

> 上面公式计算时一般认为 recvBuffer 即 BDP ；


## sendBuffer 计算

同 recvBuffer 类似，和 sendBuffer 有关的参数如下：

- `net.core.wmem_default`
- `net.core.wmem_max`
- `net.ipv4.tcp_wmem`

发送端缓冲的自动调节机制很早就已经实现，并且是**无条件开启**，没有参数去设置。如果指定了 `net.ipv4.tcp_wmem` ，则 `net.core.wmem_default` 被 `net.ipv4.tcp_wmem` 的覆盖。sendBuffer 在 `net.ipv4.tcp_wmem` 的最小值和最大值之间自动调节。如果调用 `setsockopt()` 设置了 socket 选项 `SO_SNDBUF` ，将关闭发送端缓冲的自动调节机制，`net.ipv4.tcp_wmem` 将被忽略，`SO_SNDBUF` 的最大值由 `net.core.wmem_max` 限制。

```
root@vagrant-ubuntu-trusty:~] $ sysctl -a|grep wmem
net.core.wmem_default = 212992
net.core.wmem_max = 212992
net.ipv4.tcp_wmem = 4096	16384	3827616
root@vagrant-ubuntu-trusty:~] $
```

- **`net.core.wmem_default`**

This sets the **default** OS send buffer size for **all types of connections**.

- **`net.core.wmem_max`**

This sets the **max** OS send buffer size for **all types of connections**.


- **`net.ipv4.tcp_wmem`**

    TCP Autotuning setting.
    
    This variable takes 3 different values which holds information on how much TCP sendbuffer memory space each TCP socket has to use. Every TCP socket has this much buffer space to use before the buffer is filled up. Each of the three values are used under different conditions. 
    
    - The **first** value in this variable tells the **minimum** TCP send buffer space available for a single TCP socket.
    - The **second** value in the variable tells us the **default** buffer space allowed for a single TCP socket to use.
    - The **third** value tells the kernel the **maximum** TCP send buffer space.


## TCP Segmentation 和 Checksum Offloading

操作系统可以把一些 TCP/IP 的功能转交给网卡去完成，特别是 Segmentation 和 checksum 的计算，因为这样可以节省 CPU 资源，并且由硬件代替 OS 执行这些操作会带来性能的提升。一般以太网的 MTU 为 1500 bytes ，假设应用要发送数据包的大小为 7300 bytes ，则

```
MTU 1500 字节 - IP 头部 20 字节 - TCP 头部 20 字节 = 有效负载 1460 字节
```

因此，7300 字节需要拆分成 5 个 segment 进行发送；

![未使能 TCP segmentation offload 功能](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E6%9C%AA%E4%BD%BF%E8%83%BD%20TCP%20segmentation%20offload%20%E5%8A%9F%E8%83%BD.png)

Segmentation（分片）操作可以由操作系统移交给网卡完成，虽然最终线路上仍然是传输 5 个包，但这样节省了 CPU 资源并带来性能的提升：

![使能 TCP segmentation offload 功能](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E4%BD%BF%E8%83%BD%20TCP%20segmentation%20offload%20%E5%8A%9F%E8%83%BD.png)

可以使用 `ethtool -k eth0` 查看网卡当前的 offloading 情况：

```
root@vagrant-ubuntu-trusty:~] $ ethtool -k eth0
Features for eth0:
rx-checksumming: off                    -- a
tx-checksumming: on                     -- b
	tx-checksum-ipv4: off [fixed]
	tx-checksum-ip-generic: on
	tx-checksum-ipv6: off [fixed]
	tx-checksum-fcoe-crc: off [fixed]
	tx-checksum-sctp: off [fixed]
scatter-gather: on
	tx-scatter-gather: on
	tx-scatter-gather-fraglist: off [fixed]
tcp-segmentation-offload: on                  -- 1
	tx-tcp-segmentation: on
	tx-tcp-ecn-segmentation: off [fixed]
	tx-tcp6-segmentation: off [fixed]
udp-fragmentation-offload: off [fixed]        -- 2
generic-segmentation-offload: on              -- 3
generic-receive-offload: on                   -- 4
large-receive-offload: off [fixed]
rx-vlan-offload: on
tx-vlan-offload: on [fixed]
ntuple-filters: off [fixed]
receive-hashing: off [fixed]
highdma: off [fixed]
rx-vlan-filter: on [fixed]
vlan-challenged: off [fixed]
tx-lockless: off [fixed]
netns-local: off [fixed]
tx-gso-robust: off [fixed]
tx-fcoe-segmentation: off [fixed]
tx-gre-segmentation: off [fixed]
tx-ipip-segmentation: off [fixed]
tx-sit-segmentation: off [fixed]
tx-udp_tnl-segmentation: off [fixed]
fcoe-mtu: off [fixed]
tx-nocache-copy: off
loopback: off [fixed]
rx-fcs: off
rx-all: off
tx-vlan-stag-hw-insert: off [fixed]
rx-vlan-stag-hw-parse: off [fixed]
rx-vlan-stag-filter: off [fixed]
l2-fwd-offload: off [fixed]
busy-poll: off [fixed]
root@vagrant-ubuntu-trusty:~] $
```

如果想设置网卡的 offloading 开关，可以使用 `ethtool -K` 命令，例如下面的命令关闭了 tcp segmentation offload 功能：

```
ethtool -K eth0 tso off
```


----------

其它参考：

- [How To: Network / TCP / UDP Tuning](https://wwwx.cs.unc.edu/~sparkst/howto/network_tuning.php)