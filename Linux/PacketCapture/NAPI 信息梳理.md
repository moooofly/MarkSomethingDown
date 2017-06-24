# NAPI 信息梳理

----------

NAPI 的核心在于：在一个繁忙网络，每次有网络数据包到达时，不需要都引发中断，因为高频率的中断可能会影响系统的整体效率；

假象一个场景，我们使用标准的 100Mbps 网卡，可能实际达到的接收速率为 **80Mbps**，而此时数据包平均长度为 1500Bytes，则每秒产生的中断数目为：

```
80M bits/s / (8 Bits/Byte * 1500 Byte) = 6667 中断/s
```

每秒 6667 个中断，对于系统是个很大的压力，此时其实可以转为使用轮询 (polling) 方式来处理，而不是中断；但轮询在网络流量较小的时候没有效率，因此低流量时，基于中断的方式进行处理比较合适，这就是 NAPI 出现的原因：在低流量时候使用中断方式接收数据包，而在高流量时候则使用轮询方式接收数据包。

现在内核中 NIC 基本上已经全部支持 NAPI 功能，由前面的叙述可知，NAPI 适合处理高速率数据包的处理，而带来的好处则是：

- **中断缓和 (Interrupt mitigation)**，由上面的例子可以看到，在高流量下，网卡产生的中断可能达到每秒几千次，而如果每次中断都需要系统来处理，是一个很大的压力，而 NAPI 在进行轮询时是禁止了网卡的接收中断的，这样会减小系统处理中断的压力；
- **数据包节流 (Packet throttling)**，NAPI 之前的 Linux NIC 驱动总是在接收到数据包之后产生一个 IRQ ，接着在中断服务例程里将这个 skb 加入本地的 softnet ，然后触发本地 NET_RX_SOFTIRQ 软中断后续处理。如果包速过高，因为 **IRQ 的优先级高于 SoftIRQ** ，导致系统的大部分资源都在响应中断，但 softnet 的队列大小有限，接收到的超额数据包也只能丢掉，所以这时这个模型是在用宝贵的系统资源做无用功。而 NAPI 则在这样的情况下，直接把包丢掉，不会继续将需要丢掉的数据包扔给内核去处理，这样，网卡将需要丢掉的数据包尽可能的早丢弃掉，内核将不可见需要丢掉的数据包，这样也减少了内核的压力。


----------


NAPI (“New API”) is an extension to the device driver packet processing framework, which is designed to improve the performance of high-speed networking. NAPI works through:

- **Interrupt mitigation** 
**High-speed networking** can create thousands of interrupts per second, all of which tell the system something it already knew: it has lots of packets to process. NAPI allows drivers to run with (some) interrupts disabled during times of high traffic, with a corresponding decrease in system load.
- **Packet throttling** 
When the system is overwhelmed and must drop packets, it's better if those packets are disposed of before much effort goes into processing them. **NAPI-compliant drivers** can often cause packets to be dropped in the network adaptor itself, before the kernel sees them at all.

New drivers should use NAPI if the hardware can support it. However, NAPI additions to the kernel do not break backward compatibility and drivers may still process completions directly in interrupt context if necessary.

## NAPI Driver design

- Packets should not be passed to netif_rx(); instead, use: `int netif_receive_skb(struct sk_buff *skb);`.
- The budget parameter places a limit on the amount of work the driver may do.
- The poll() function must return the amount of work done.
- If and only if the return value is less than the budget, your driver must reenable interrupts and turn off polling.

## Hardware Architecture

NAPI, however, requires the following features to be available:

- **DMA ring or enough RAM to store packets in software devices.**
- **Ability to turn off interrupts or maybe events that send packets up the stack.**

NAPI processes packet events in what is known as napi→poll() method. Typically, only packet receive events are processed in napi→poll(). The rest of the events MAY be processed by the regular interrupt handler to reduce processing latency (justified also because there are not that many of them).

Note, however, NAPI does not enforce that napi→poll() only processes receive events. Tests with the tulip driver indicated slightly increased latency if all of the interrupt handler is moved to napi→poll(). Also MII/PHY handling gets a little trickier.

## Advantages

NAPI provides an “inherent mitigation” which is bound by system capacity.

## Disadvantages

### Latency

In some cases, NAPI may introduce additional software IRQ latency.

### IRQ masking

On some devices, changing the IRQ mask may be a slow operation, or require additional locking. This overhead may negate any performance benefits observed with NAPI


----------


## NAPI

NAPI (New API) 是 Linux 新的网卡数据处理 API ，据说是由于找不到更好的名字，在 2.5 之后引入。NAPI 是 Linux 上采用的一种提高网络处理效率的技术；

**中断方式**的好处是响应及时，如果数据量较小，则不会占用太多的 CPU 时间，缺点是数据量大时，会产生过多中断，而每个中断都要消耗不少的 CPU 时间，从而导致效率反而不如轮询高。**轮询**方式与中断方式相反，它更适合处理大量数据，因为每次轮询不需要消耗过多的 CPU 时间；缺点是即使只接收很少数据或不接收数据时，也要占用 CPU 时间。

简单来说，NAPI 是综合**中断方式**与**轮询方式**的技术，数据量低时采用中断，数据量高时采用轮询。平时是中断方式，当有数据到达时，会触发中断处理函数执行，中断处理函数关闭中断，开始处理。如果此时有数据到达，则没必要再触发中断了，因为中断处理函数中会轮询处理数据，直到没有新数据时才打开中断。

很明显，数据量很低与很高时，NAPI 可以发挥中断与轮询方式的优点，性能较好。如果数据量不稳定，且说高不高说低不低，则 NAPI 则会在两种方式切换上消耗不少时间，效率反而较低一些。
 

### NAPI 和 non-NAPI 的区别

- 支持 NAPI 的网卡驱动必须提供轮询方法 `poll()` ；
- non-NAPI 使用的内核接口为 `netif_rx()` ，NAPI 使用的内核接口为 `napi_schedule()` ；
- non-NAPI 使用共享的 CPU 队列 `softnet_data->input_pkt_queue` ，NAPI 使用设备内存或者
设备驱动程序的环状缓冲区；

### NAPI 解决了什么问题
 
- 第一，它限制了中断的数量，一旦有中断过来就停掉中断改为轮询，这样就不会造成 cpu 被频繁中断；
- 第二，cpu 不会做无用功，就是所谓的无用的轮询，因为只有在中断来了才改为轮询，中断来了说明有事可做；

###  NAPI 存在的严重缺陷

- 对于上层的应用程序而言，由于系统不能在接收到每个数据包时都及时处理，因此随着传输速度的增加，累计的数据包将会耗费大量的内存，经过实验表明，在 Linux 平台上这个问题会比在 FreeBSD 上要严重一些；
- 另外一个问题是对于大的数据包处理比较困难，原因是大的数据包在传送到网络层所耗费的时间比短数据包长很多（即使是采用 DMA 方式），所以如前所述，NAPI 技术适用于针对高速率、短长度数据包处理。

### 使用 NAPI 的先决条件

驱动可以继续使用老的 2.4 内核的网络驱动程序接口，NAPI 的加入并不会导致向前兼容性的丧失，但是 NAPI 的使用至少要得到下面的保证：

- 要使用 **DMA 环形输入队列**（也就是 `ring_dma`，这个在 2.4 驱动中关于 Ethernet 的部分有详细的介绍），或者是有足够的内存空间缓存驱动获得的包。
- 在发送/接收数据包产生中断的时候，有能力关断 NIC 中断的事件处理，并且在关断 NIC 以后，并不影响数据包被接收到**网络设备的环形缓冲区**处理队列中。


### Linux 网卡驱动与 NAPI

目前 NAPI 技术已经在网卡驱动层和网络层得到了广泛的应用；

由于在遇到**突发快速小包传输**的时候，过去的实现机制会导致频繁中断，造成大量遗漏（造成丢包的一种情况）和竞争；因此就将传统的网卡收发机制也纳入到了 NAPI 框架；

Linux 网卡驱动中的 NAPI 方式设计十分巧妙，就是在第一个包到来的时候中断，然后关闭中断开始轮询，等某一次轮询完毕后发现没有数据了，那么内核默认此次数据已经传输完毕，短时间内不会再有数据了，那么停止轮询，重新开启中断；


----------

- NAPI 运行概览

![NAPI 运行概览](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/NAPI%20%E8%BF%90%E8%A1%8C%E6%A6%82%E8%A7%88.jpg "NAPI 运行概览")

- 分组到NIC后穿过内核到达网络层路径

![分组到NIC后穿过内核到达网络层路径](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E5%88%86%E7%BB%84%E5%88%B0NIC%E5%90%8E%E7%A9%BF%E8%BF%87%E5%86%85%E6%A0%B8%E5%88%B0%E8%BE%BE%E7%BD%91%E7%BB%9C%E5%B1%82%E8%B7%AF%E5%BE%84.jpg "分组到NIC后穿过内核到达网络层路径")


----------


## 参考资料


- [内核接收分组理解](http://www.cnblogs.com/lxgeek/p/4182029.html)
- [数据包接收系列 - NAPI 的原理和实现](http://blog.csdn.net/zhangskd/article/details/21627963)
- [Linux 内核 NAPI 机制分析](http://blog.csdn.net/joshua_yu/article/details/591041)
- [Linux 下网络性能优化方法简析](https://www.ibm.com/developerworks/cn/linux/l-cn-network-pt/index.html)
- [napi](https://wiki.linuxfoundation.org/networking/napi)