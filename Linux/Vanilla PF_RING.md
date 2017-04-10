# Vanilla PF_RING

> 原文地址：[Vanilla PF_RING](https://raw.githubusercontent.com/ntop/PF_RING/dev/doc/README.vanilla.md)

------

Vanilla PF_RING 的构成：

1. 内核加速模块（the accelerated kernel module），在 low-level 层次上提供了将 packet 拷贝到 PF_RING rings 的功能；
2. 用户空间 PF_RING SDK ，为用户空间应用程序提供了针对 PF_RING 的透明支持；

PF_RING 实现了一种 socket 类型，基于该 socket 类型实现了用户空间应用程序直接与 PF_RING 内核模块进行“对话”的能力；

应用程序可以通过获取到的 PF_RING 句柄发起各种 API 调用；

该句柄可以绑定到：

1. 物理网络接口；
2. 某个 RX queue ，仅在 multi-queue 网络适配器上允许；
3. 名为 'any' 的虚拟接口，即所有系统接口上的收发包均可被获取；

正如上面所说，packets 读取发生在 memory ring 中（创建时分配内存）；而进入的 packets 由内核模块负责拷贝到该 ring 中，之后再被用户空间应用所读取；内存的分配和释放并非基于每个 packet 进行，一旦某个 packet 被从 ring 中读取走了，那么之前 ring 中用于保存该 packet 的空间将被用作安放后续其它的 packets ；这就意味着，如果应用想要维护一个 packet archive ，则必须自行保存读取到的所有 packets ，因为 PF_RING 不会为应用做保存；

## Packet Filtering

PF_RING 既支持 legacy **BPF** filters（即那些被 pcap-based 应用，如 tcpdump ，所支持的 filters 规则），也支持两种额外类型的 filters（取名为 **wildcard** 和 **precise** filters ，分别对应了一部分或全部 filter elements 被指定的情况），提供给开发者很多选择可能；

Filters 是在 PF_RING 模块内部起作用的，即位于 kernel 中；一些现代适配器，例如 Intel 82599-based 或 Silicom Redirector 的 NICs ，支持了基于硬件的 filters ，PF_RING 同样通过特定的 API 调用进行了支持（例如 `pfring_add_hw_rule`）；另外，PF_RING filters （除 hw filters 外）能够指定 action ，用于告知 PF_RING 内核模块在发现匹配 filter 的 packet 出现时，将执行何种 action ；

Actions 包括 

- pass/don’t pass the filter to the user space application
- stop evaluating the filter chain, or 
- reflect packet. 

在 PF_RING 中，`packet reflection` 指的是传输（不做修改）匹配 filter 的 packet 到指定网络接口的能力（此处不包括接收到该 packet 的接口）；整个 reflection 功能实现在 PF_RING 内核模块之中，而对于用户空间应用程序来说，唯一能够请求的活动就是设置 filter specification ，因此，是无法针对 packet 进行任何其它处理的；

## Packet Clustering

PF_RING 还能够进一步提升 packet capture 应用的性能，基于已实现的两种机制：**balancing** 和 **clustering** ；这些机制允许一些应用程序（愿意处理一部分划分出来的 packets）只处理整个 packet stream 中的一部分，而将所有其它包发送到 cluster 中的其他成员来处理；这意味着，打开了 PF_RING sockets 的不同应用程序均能够将自身绑定到指定的 cluster Id 上（通过 `pfring_set_cluster`），进而成为包处理“劳动力”，以便对数据包的一部分进行分析处理；

在 cluster sockets 中进行 packets 划分的方法是通过 cluster policy 进行指定的，可以指定为（默认值）**per-flow**（即所有的 packets 都归属于相同的五元组 `<proto, ip src/dst, port src/dst>`）或 **round-robin** ；

这意味着，如果你选择的是 per-flow balancing ，那么所有归属相同 flow 的 packets 都将被发送给同一个应用程序；若选择的是 round-robin ，那么所有的应用都将收到相同数量的 packets ，但是不保证归属于相同 queue 的 packets 一定被单独一个应用所接收；因此一方面来讲，per-flow balancing 允许保留（preserve）应用逻辑，因为在这种情况下，应用将只会接收所有 packets 的子集，但保证了 traffic 的一致性（consistent）；而另一方面，如果你恰好遇到了某条特定的 flow 汇聚了几乎全部的 traffic 情况，那么处理该 flow 的应用程序将会出现 over-flooded 问题，并且此时 traffic 分布也将严重不均衡；
