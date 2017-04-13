# [PF_RING and Transparent Mode](http://www.ntop.org/pf_ring/pf_ring-and-transparent-mode/)

PF_RING 被设计用以增强 packet 捕获性能；这意味着需要在 RX path 上进行加速，而常常被使用的一种方式就是减少从适配器到用户态 packet 所经过的“旅途”；这可以通过允许驱动程序将 packet 直接从 NIC 推送到 PF_RING 实现，即不再经过通常的 kernel path ；基于这个原因，PF_RING 引入了一个名为 “transparent mode” 的选项，其目的就是用于调节（tune）packets 从 NIC 搬移到 PF_RING 的方式；该选项（可以在基于 `insmod` 添加 PF_RING 模块时指定）可以设置如下三种值：

- **`insmod pf_ring.ko transparent_mode=0`**
默认情况，意味着 packets 会通过标准内核机制发送到 PF_RING 中；在该设置下，packets 会被发送给 PF_RING ，而不会发给所有其它内核组件；所有 NIC drivers 都支持该模式；

- **`insmod pf_ring.ko transparent_mode=1`**
在该模式下，packets 是由 NIC driver 直接发送给 PF_RING 的，但 packets 仍会被传递给（propagated）其它内核组件；在该模式下，packet 捕获能被加速的原因是，packets 从 NIC driver 中被拷贝出去时无需经过常规的 kernel path ；需要注意的是，为了使能该模式，你必须使用支持 PF_RING 的 NIC driver ；可用的 PF_RING-enabled drivers 已放在了 `drivers/` 目录中；

- **`insmod pf_ring.ko transparent_mode=2`**
在该模式下，packets 是由 NIC driver 直接发送给 PF_RING 的，而不会发给所有其它内核组件，因为会拖慢 packet 捕获过程；需要注意的是：
    - 为了使能该模式，你必须使用支持 PF_RING 的 NIC driver ；
    - Packets 在被投递给（delivered）PF_RING 后不会再被发送到 kernel ；这意味着你将无法从基于 PF_RING-aware drivers 的 NICs 获取到活动性信息（connectivity）；
    - 该模式是最快的方式，因为 packets 会以最快的速度拷贝给 PF_RING ，并在得到处理后立即丢弃；
