# TCP 之 TSO GSO LSO

## TSO: TCP Segment Offload

Abbreviated as TSO, TCP segmentation offload is used to reduce the CPU overhead of TCP/IP on fast networks. TSO breaks down large groups of data sent over a network into smaller segments that pass through all the network elements between the source and destination. This type of offload relies on the network interface controller (NIC) to segment the data and then add the TCP, IP and data link layer protocol headers to each segment. The NIC must support TSO. TSO is also called large segment offload (LSO).


TSO 全称为 TCP Segment Offload ，简单的讲，就是靠网卡硬件来分段 TCP ，计算 checksum ，从而解放 CPU 周期。

我们知道通常以太网的 MTU 是 1500 ，除去 TCP/IP 的包头，TCP 的 MSS(Max Segment Size)大小是 1460 ；通常情况下协议栈会对超过 1460 的 TCP payload 进行 segmentation 以保证生成的 IP 包不超过 MTU 的大小；但是对于支持 TSO/GSO 的网卡而言，就没这个必要了；我们可以把最多 64K 大小的 TCP payload 直接往下传给协议栈，此时 IP 层也不会进行 segmentation ，而是会一直传给网卡驱动，支持 TSO/GSO 的网卡会自己生成 TCP/IP 包头和帧头，这样可以 offload 很多协议栈上的内存操作，checksum 计算等原本靠 CPU 来做的工作都移给了网卡；

而 GSO 可以看作是 TSO 的增强，但 GSO 不只针对 TCP ，而是对任意协议；其会尽可能把 segmentation 推后到交给网卡的那一刻，才会判断网卡是否支持 SG(scatter-gather) 和 GSO ；如果不支持，则在协议栈里做 segmentation ，如果支持，则把 payload 直接发给网卡；

可以通过如下命令查看 TSO 支持情况的信息

```shell
ethtool -k <interface>
```

目前很多网卡都支持 TSO，但很少有支持 UFO 的，而 GSO/GRO 和网卡无关，只是内核的特性。

## GSO: Generic Segmentation Offload


GRO 是在内核 2.6.29 之后合并进去的，简介可以看[这里](http://lwn.net/Articles/358910/)；

GRO 的作用：GRO 针对网络包进行**接收**处理，并且只针对 NAPI 类型的驱动；因此，对 GRO 的支持，不仅要求内核支持，还要求驱动也必须调用相应的接口；使用 `ethtool -K gro on` 命令进行设置时，如果报错，则说明网卡驱动本身就不支持 GRO 。

GRO 类似 TSO ，但 TSO 只支持**发送**数据包；若不支持 GRO ，那么小数据段会被一个个的送到协议栈；若支持 GRO ，那么就会在接收端执行一个相对于 TSO 的反向操作，即将 TSO 切好的数据包组合成大包后，再传递给协议栈。

GRO 的主要思想就是，基于一些数据域，组合一些类似数据包为一个大数据包（对应一个 skb），然后 feed 给协议栈，这里主要是利用 Scatter-gather IO 来合并数据包；

---

This series adds Generic Segmentation Offload (**GSO**) support to the Linux
networking stack.

Many people have observed that a lot of the savings in **TSO** come from
traversing the networking stack once rather than many times for each
**super-packet**.  These savings can be obtained without hardware support.
In fact, the concept can be applied to other protocols such as TCPv6,
UDP, or even DCCP.

The key to minimising the cost in implementing this is to **postpone the
segmentation as late as possible**.  In the ideal world, the segmentation
would occur inside each NIC driver where they would rip the super-packet
apart and either produce SG lists which are directly fed to the hardware,
or linearise each segment into pre-allocated memory to be fed to the NIC.
This would elminate segmented skb's altogether.

Unfortunately this requires modifying each and every NIC driver so it
would take quite some time.  A much easier solution is to perform the
segmentation just before the entry into the driver's xmit routine.  This
series of patches does this.

I've attached some numbers to demonstrate the savings brought on by
doing this.  **The best scenario** is obviously the case where the underlying
NIC supports SG.  This means that we simply have to manipulate the SG
entries and place them into individual skb's before passing them to the
driver.  The attached file lo-res shows this.

The test was performed through the loopback device which is a fairly good
approxmiation of an SG-capable NIC.

**GSO like TSO is only effective if the MTU is significantly less than the
maximum value of 64K**.  So only the case where the MTU was set to 1500 is
of interest.  There we can see that the throughput improved by 17.5%
(3061.05Mb/s => 3598.17Mb/s).  The actual saving in transmission cost is
in fact a lot more than that as the majority of the time here is spent on
the RX side which still has to deal with 1500-byte packets.

**The worst-case scenario** is where the NIC does not support SG and the user
uses `write(2)` which means that we have to copy the data twice.  The files
gso-off/gso-on provide data for this case (the test was carried out on
e100).  As you can see, the cost of the extra copy is mostly offset by the
reduction in the cost of going through the networking stack.

For now GSO is off by default but can be enabled through `ethtool`.  It is
conceivable that with enough optimisation GSO could be a win in most cases
and we could enable it by default.

However, even without enabling GSO explicitly it can still function on
bridged and forwarded packets.  As it is, passing TSO packets through a
bridge only works if all constiuents support TSO.  With GSO, it provides
a fallback so that we may enable TSO for a bridge even if some of its
constituents do not support TSO.

This provides massive savings for Xen as it uses a bridge-based architecture
and TSO/GSO produces a much larger effective MTU for internal traffic between
domains.

## LSO: Large segment offload

In computer networking, large segment offload (**LSO**) is a technique for increasing `outbound` throughput of high-bandwidth network connections by reducing CPU overhead. It works by queuing up large buffers and letting the network interface card (**NIC**) split them into separate packets. The technique is also called TCP segmentation offload (**TSO**) when applied to TCP, or generic segmentation offload (**GSO**).

The `inbound` counterpart of large segment offload is large receive offload (**LRO**).

When a system needs to send large chunks of data out over a computer network, the chunks first need breaking down into smaller segments that can pass through all the network elements like **routers** and **switches** between the source and destination computers. This process is referred to as **segmentation**. Often the TCP protocol in the host computer performs this segmentation. Offloading this work to the NIC is called *TCP segmentation offload* (**TSO**).

For example, a unit of 64kB (65,536 bytes) of data is usually segmented to 46 segments of 1448 bytes each before it is sent through the NIC and over the network. With some intelligence in the NIC, the host CPU can hand over the 64 KB of data to the NIC in a single transmit-request, the NIC can break that data down into smaller segments of 1448 bytes, add the TCP, IP, and data link layer protocol headers - according to a template provided by the host's TCP/IP stack - to each segment, and send the resulting frames over the network. This significantly reduces the work done by the CPU. As of 2014 many new NICs on the market support TSO.

Some network cards implement TSO generically enough that it can be used for offloading fragmentation of other transport layer protocols, or by doing **IP fragmentation** for protocols that don't support fragmentation by themselves, such as **UDP**.



参考：

- [Large segment offload](https://en.wikipedia.org/wiki/Large_segment_offload)
- [关于GRO/GSO/LRO/TSO等patch的分析和测试](http://blog.csdn.net/majieyue/article/details/7929398)
- [GSO: Generic Segmentation Offload]( http://lwn.net/Articles/188489/)
- [linux kernel 网络协议栈之GRO(Generic receive offload)](http://www.pagefault.info/?p=159)




