# [Introducing PF_RING DNA (Direct NIC Access)](http://www.ntop.org/pf_ring/introducing-pf_ring-dna-direct-nic-access/)

在此隆重宣布 PF_RING DNA (Direct NIC Access) 已经可用了，在 Linux 下，和基于 PF_RING (non-DNA) 进行 packet 捕获比较发现，其能够极大的提高性能（高达 80%）；

![PF_RING DNA](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/PF_RING%20DNA%20\(Direct%20NIC%20Access\).png "PF_RING DNA")


PF_RING 通过 Linux NAPI 从 NICs 进行 packets 的 polling ，这意味着 NAPI 会从 NIC 中拷贝 packets 到 PF_RING 的环形缓冲区（circular buffer），之后用户态应用程序再从该环形缓冲区中进行 packets 读取（polling）；在这种情况下，存在两次 poll 操作：应用程序和 NAPI 各一次；这导致了 CPU cycles 被这些 polling 所耗费；**优点**就是 PF_RING 能够将 incoming packets 同时分发到多个环形缓冲区中（即对应多个应用程序）；

**PF_RING DNA (Direct NIC Access)** 是一种**映射 NIC memory 和 registers 到用户态**的方案，因此从 NIC 拷贝 packets 到 DMA ring 中的操作是由 NIC **NPU (Network Process Unit)** 完成的，非 NAPI 完成；这必然带来了更好的性能表现，因为 CPU cycles 仅被用于 consuming packets ，未被浪费在从网络适配器上进行 packets 搬移动作上；而**缺点**就是一次只能有一个应用程序可以打开 DMA ring ，或者换句话说，用户态应用程序需要相互“沟通”以便进行 packets 的分发（distribute）；

简单来说，**如果您喜欢灵活性，那么你应该使用 PF_RING ，如果你想要追求纯粹的速度快，那么 PF_RING DNA 是解决方案**；需要注意的是，在 DNA 模式下，NAPI polling 不会发生，因此诸如 reflection 和 packet filtering 等 PF_RING 特性不再被支持；
