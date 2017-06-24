# PF_RING

> 原文地址：[这里](http://www.ntop.org/products/packet-capture/pf_ring/)

## High-speed packet capture, filtering and analysis.

PF_RING™ 是一种新型网络 socket ，能够极大到提高 packet 捕获速度，其具有以下特性：

- 可用于 kernel 2.6.32 及之后的版本；
- 无需对 kernel 打 patch ：只需加载 kernel 模块；
- 采用（常用的）商用网络适配器就可提供 10 Gbit 的 [Hardware Packet Filtering](http://www.ntop.org/products/packet-capture/pf_ring/hardware-packet-filtering/) 能力；
- 支持用户空间 [ZC](http://www.ntop.org/products/packet-capture/pf_ring/pf_ring-zc-zero-copy/) (新一代 DNA, Direct NIC Access) drivers ；在通过 NIC NPU (Network Process Unit) 将 packets 推到（pushing to）用户空间时，或者从用户空间获取（getting from）时，即便处于极端的 packet 捕获/传输速度情况下，也不需要任何形式的 kernel 干预；使用  10Gbit ZC driver 后，对于任何大小的 packet ，你都可以按照 wire-speed 发送或接收；
- [PF_RING ZC](http://www.ntop.org/products/packet-capture/pf_ring/pf_ring-zc-zero-copy/) 库可用于以 zero-copy 方式跨 threads, applications, Virtual Machines 进行 packets 分布；
- 设备驱动无关（independent）；
- 支持 Myricom, Intel 及 Napatech 的网络适配器；
- Kernel-based 的 packet 捕获和抽样；
- 支持 Libpcap ，可无缝集成已存在 pcap-based 应用；
- 除 BPF 外，还可以指定数以百计的 header filters ；
- 支持内容检视（Content inspection）功能，因此只有满足 payload filter 规则的 packets 才能通过；
- 可以通过 PF_RING™ 插件支持高级 packet 解析和内容过滤；

更多 PF_RING™ 的内部实现可以查看[用户手册](http://www.ntop.org/support/documentation/documentation/)；


> - Napatech 是提供网络管理和安全应用数据交付解决方案的全球领导者。随着数据量增多、复杂性增大，组织必须对流经其网络的信息进行监控、编译和分析。产品使用专利技术，能高速捕获和处理大量数据，同时保证性能并实现实时可见性；
> - Myricom 是著名的网络设备供应商，是为垂直市场应用提供高性能、低延时优化以太网解决方案的行业领导者；Myricom 在低延时、高包速率、垂直应用优化方面做的很好；Myri-10G 网络适配器具有低延迟、高带宽及 CPU 资源占用低的特点，能提供终极性能的以太网及 Myricom 专业软件平台；

## Vanilla PF_RING™

PF_RING™ 是基于 Linux NAPI 从 NICs 中进行 packets 轮询；这意味着 NAPI 会从 NIC 中拷贝 packets 到 PF_RING™ 的 circular buffer ，之后用户态应用程序再从该 ring 中读取这些 packets ；在这种场景下，存在两次轮询动作，即应用和 NAPI 各一次，这就会导致一些 CPU 时钟周期在轮询中被耗费掉了；而这种实现方式的好处就是：PF_RING™ 能够将进入的 packets 同时分发到多个 rings 中（对应多个应用）；

![vanilla_pf_ring](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/vanilla_pf_ring.png "vanilla_pf_ring")

## PF_RING™ Modules

PF_RING™ 实现了模块化架构（modular architecture），因此它能够使用除 PF_RING™ 这个 kernel 模块外的额外组件；当前，可用的其他模块包括：

- ZC module.
详情查看 [ZC](http://www.ntop.org/products/packet-capture/pf_ring/pf_ring-zc-zero-copy/) 页面说明；
- Accolade module.
该模块在 PF_RING™ 中为 Accolade cards 增加了本地支持；
- Endace module.
该模块在 PF_RING™ 中为 Endace DAG cards 增加了本地支持；
- Exablaze module.
该模块在 PF_RING™ 中为 Exablaze cards 增加了本地支持；
- Fiberblaze module.
该模块在 PF_RING™ 中为 Fiberblaze cards 增加了本地支持；
- Myricom module.
该模块在 PF_RING™ 中为 Myricom 10 Gbit cards 增加了本地支持；
- Napatech module.
该模块在 PF_RING™ 中为 Napatech cards 增加了本地支持；
- Stack module.
该模块可用于向 linux 网络协议栈中 inject packets ；
- Timeline module.
该模块可通过使用 PF_RING™ API 从 `n2disk` dump set 中无缝提取 traffic ；
- Sysdig module.
该模块使用 sysdig 内核模块捕获系统事件；

![PF_RING](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/PF_RING-modules.jpeg "PF_RING")

## Who needs PF_RING™?

一般来讲，任何需要面对 pps 比较高（many）情况的都需要；术语 ‘many’ 的含义针对你所使用的、用于 traffic 分析的不同硬件含义有所不同；可能从 1,2GHz ARM 上的 80k pkt/sec 变化成 low-end 2,5GHz Xeon 上的 14M pkt/sec 及以上的速度； PF_RING™ 不仅为你提供了更快的 packets 捕获速度，还保证了 packets 捕获的高效（减少对 CPU cycles 的浪费）；
