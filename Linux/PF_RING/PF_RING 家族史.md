# PF_RING 家族史

> 参考自：[这里](http://jasonzhuo.com/pfring-and-jnetpcap/)

## PF_RING 简介

PF_RING 是 Luca Deri 发明的、基于提高内核处理数据包效率，进而改善网络数据包捕获程序（如 libpcap 和在其上的 tcpdump 等）的东东。PF_RING 是一种新型的网络 socket ，它可以极大的改进包捕获的速度。

### 术语

- **NAPI**：NAPI 是 Linux 新的网卡数据处理 API ，NAPI 是一种综合了中断与轮询两种方式的技术；NAPI 在高负载的情况下可以产生更好的性能，它避免了为每个传入的数据帧都产生中断。详情可以参考[这里](http://blog.csdn.net/zhangskd/article/details/21627963)；
- **Zero copy (ZC)**：简单一点来说，零拷贝就是一种避免 CPU 将数据从一块存储拷贝到另外一块存储的技术；
- **NPU**：网络处理单元；
- **DMA**：即 Direct Memory Access ，直接内存存取；它允许不同速度的硬件装置来沟通，而不需要依赖于 CPU 的大量中断负载；
- **Linux 网络栈**：如下图所示，它简单地为用户空间的应用程序提供了一种访问内核网络子系统的方法；

![Linux 网络栈](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Linux%20%E7%BD%91%E7%BB%9C%E6%A0%88.gif "Linux 网络栈")


## libpcap 抓包原理

libpcap 的包捕获机制就是**在数据链路层加一个旁路处理**。当一个数据包到达网络接口时，libpcap 首先利用已经创建的 Socket 从链路层驱动程序中获得该数据包的拷贝，再通过 **Tap 函数**将数据包发给 **BPF 过滤器**。BPF 过滤器根据用户已经定义好的过滤规则对数据包进行逐一匹配，匹配成功则放入内核缓冲区（一次拷贝），并传递给用户缓冲区（又一次拷贝），匹配失败则直接丢弃。如果没有设置过滤规则，所有数据包都将放入内核缓冲区，并传递给用户层缓冲区。

![Libpcap 的包捕获机制](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Libpcap%20%E7%9A%84%E5%8C%85%E6%8D%95%E8%8E%B7%E6%9C%BA%E5%88%B6.png "Libpcap 的包捕获机制")

在高速复杂网络环境下 **libpcap 丢包的原因**主要有以下两个方面：

- Cpu 处于频繁中断状态，造成接收数据包效率低下；
- 数据包被多次拷贝，浪费了大量时间和资源。从网卡驱动到内核，再从内核到用户空间；


### 为啥用 PF_RING 呢

随着信息技术的发展，1 Gbit/s，10 Gbit/s 以及 100 Gbit/s 的网络会越来越普及，那么**零拷贝技术**也会变得越来越普及，这是因为网络链接的处理能力比 CPU 的处理能力的增长要快得多。高速网络环境下，CPU 就有可能需要花费几乎所有的时间去拷贝要传输的数据，而没有能力再去做别的事情，这就产生了性能瓶颈。

## PF_RING 驱动家族

### PF_RING DNA (Direct NIC Access)

对于那些希望在 CPU 利用率为 0%（拷贝包到主机）的情况下，想要最大化数据包捕获速度的用户来说，可以使用 **DNA (Direct NIC Access)** 驱动；它允许数据直接从网络接口上读取，它以**零拷贝**的方式**同时绕过 Linux 内核和 PF_RING 模块**。

左图解释：**Vanilla PF_RING** 从 NIC 上通过 Linux NAPI 获取数据包拷贝，拷贝到 PF_RING 的 环状缓存空间。然后用户空间的应用程序会从该环状缓存空间中读取数据包。从图中可以看出 Vanilla PF_RING 方式有**两次 polling 操作**：一次是从 NIC 到 PF_RING 环状缓存空间（Linux 内核里面），另外一次从 PF_RING 环状缓存空间到用户程序。

左图中的实现（即 Vanilla PF_RING）相对于传统方式来说，由于 Application 使用的是**基于 `mmap` 的 libpcap 版本**，会较标准版本的 libpcap 效率更高。**libpcap 标准版**是目前使用最多的、用于从内核拷贝数据包到用户层的库，而 libpcap-mmap 是 libpcap 的一个改进版本。**传统 libpcap** 使用固定大小的**存储缓冲器**和**保持缓冲器**来完成数据包从内核到用户层的传递，而 **libpcap-mmap** 设计了一个大小可以配置的循环缓冲器，允许用户程序和内核程序同时对该循环缓冲器的不同数据区域进行直接的读取。其次，PF_RING 使用的是 NAPI ，而不是（标准）传统 libpcap 中基于 **DMA** 方式（调用系统函数 `netif_rx()`）将数据包从网卡拷贝到内核缓存。

右图解释：在 DNA 模式下，会通过 NIC 的 NPU 拷贝数据包到**网卡上的缓存空间**。然后直接通过 DMA 方式到用户层，同时**绕过了 PF_RING 模块和 Linux 内核缓存**。官网已经证明了 DNA 方式能有更快的处理速率 Performance 。

![PF_RING DNA (Direct NIC Access)](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/PF_RING%20DNA%20\(Direct%20NIC%20Access\).png "PF_RING DNA (Direct NIC Access)")

让人多少有点遗憾的是，该方式并不是可以随意使用的；根据其 License 说明，提供免费下载但是以二进制形式提供测试版本的库（也就是说使用 5 分钟或者达到一定包的处理数量之后软件就停了），如果需要长期使用，需要购买解锁的代码。

```
ERROR: You do not seem to have a valid DNA license for eth0 [Intel 1 Gbit e1000e family]. We're now working in demo mode with packet capture and transmission limited to 0 day(s) 00:05:00
```

关于传统数据发送的整个过程，可以看《[Linux 中的零拷贝技术](https://www.ibm.com/developerworks/cn/linux/l-cn-zerocopy1/)》；

### PF_RING-aware drivers (ZC support)

根据官方手册介绍：

> An interface in ZC mode provides the same performance as DNA. 

PF_RING ZC 和 PF_RING DNA 实际上都是**绕过 Linux 内核和 PF_RING 模块**的方式，因此在这些模式下 Linux 内核将看不到任何数据包。

### PF_RING ZC

关于 PF_RING ZC，PF_RING DNA，PF_RING-aware drivers 之间的关系是有点乱的，不太容易分清楚。

> "It can be considered as the successor of DNA/LibZero that offers a single and consistent API implementing simple building blocks (queue, worker and pool) that can be used from threads, applications and virtual machines."
> 
> "PF_RING ZC comes with a new generation of PF_RING aware drivers."

### 相互关系

- PF_RING ZC 的 API 更友好；
- PF_RING ZC 模块可以看做是 DNA/LibZero 的后继者；
- PF_RING ZC 和 PF_RING DNA 都是**绕过 Linux 内核和 PF_RING 模块**的方式，在这些模式下 Linux 内核将看不到任何数据包；


## PF_RING 工作模式

PF_RING 有三种工作模式（对应 `transparent_mode` 的不同值）：

- 为 0 时，走的是 Linux 标准的 NAPI 包处理流程；
- 为 1 时，包既走 Linux 标准包处理流程，也 copy 给 PR_RING 一份；
- 为 2 时，驱动只将包拷贝给 PF_RING ，内核不会接收到这些包；

1 和 2 模式需要 PF_RING-aware 的网卡驱动支持。

另外：

- 默认为 `transparent=0` ，数据包通过标准的 linux 接口接收，任何驱动都可以使用该模式；
- `transparent=1`（用于 Vanilla PF_RING 和 PF_RING-aware 驱动程序），数据包分别拷贝到 PF_RING 和标准 linux 网络协议栈各一份；
- `transparent=2`（用于 PF_RING-aware 驱动程序），数据包仅拷贝到 PF_RING ，而不会拷贝到标准的 linux 网络协议栈（即 tcpdump 不会看到任何数据包）；

不要同时使用模式 1 和模式 2 到 Vanilla 驱动，否则将会抓到任何数据包。

## PF_RING 包过滤

Vanilla PF_RING 支持传统的 BPF 过滤器，由于 DNA 模式下，不再使用 NAPI Poll ，所以 PF_RING 的数据包过滤功能就不支持了；目前可以使用硬件层数据包过滤的只有 intel 的 82599 网卡支持。