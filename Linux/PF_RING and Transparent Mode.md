# [PF_RING and Transparent Mode](http://www.ntop.org/pf_ring/pf_ring-and-transparent-mode/)

PF_RING 被设计用以增强 packet 捕获性能；这意味着需要在 RX path 上进行加速，而常常被使用的一种方式就是减少从适配器到用户态 packet 所经过的“旅途”；这可以通过允许驱动程序将 packet 直接从 NIC 推送到 PF_RING 实现，即不再经过通常的 kernel path ；基于这个原因，PF_RING 引入了一个名为 “transparent mode” 的选项，其目的就是用于调节（tune）packets 从 NIC 搬移到 PF_RING 的方式；该选项（可以在基于 `insmod` 添加 PF_RING 模块时指定）可以设置如下三种值：

- **`insmod pf_ring.ko transparent_mode=0`**
默认值，意味着 packets 会通过标准内核机制（NAPI）发送到 PF_RING 中；在该设置下，packets 会被发送给 PF_RING 和所有其它内核组件（不确定这里理解的是否正确）；所有 NIC drivers 都支持该模式；

- **`insmod pf_ring.ko transparent_mode=1`**
在该模式下，packets 是由 NIC driver 直接发送给 PF_RING 的，但 packets 同样会被传递给（propagated）其它内核组件；在该模式下，packet 捕获能被加速的原因是，packets 从 NIC driver 中被拷贝出去时无需经过常规的 kernel path ；需要注意的是，为了使能该模式，你必须使用支持 PF_RING 的 NIC driver ；可用的 PF_RING-enabled drivers 已放在了 `drivers/` 目录中；

- **`insmod pf_ring.ko transparent_mode=2`**
在该模式下，packets 是由 NIC driver 直接发送给 PF_RING 的，但不会发给所有其它内核组件（因为会拖慢 packet 捕获过程）；需要注意的是：
    - 为了使能该模式，你必须使用支持 PF_RING 的 NIC driver ；
    - Packets 在被投递给（delivered）PF_RING 后不会再被发送到 kernel ；这意味着你将无法从基于 PF_RING-aware drivers 的 NICs 获取到活动性信息（connectivity）；
    - 该模式是最快的方式，因为 packets 会以最快的速度拷贝给 PF_RING ，并在得到处理后立即丢弃；

------

# [PF_RING and transparent_mode](http://www.ntop.org/pf_ring/pf_ring-and-transparent_mode/)

许多 PF_RING 用户都知道为了不对 Linux kernel 打补丁，从 PF_RING 4.x 开始 packets 是通过 NAPI 机制进行接收的；这意味着 packet 的 journey 与标准 Linux 中的路径一样；因此，相对 vanilla Linux 的性能改进也是很微小的（< 5%），尽管 PF_RING 允许完成很多标准 AF_PACKET 之外的事；

为了大幅改进的性能，PF_RING 支持了一个名为 `transparent_mode` 的参数，可在进行内核模块加载时使用，如 `insmod pf_ring.ko transparent_mode=X` ，其中 X 的值可以为 0, 1 或 2 ；具体含义会随着你所 hook 的 PF_RING-based 应用（例如 `pfcount` ）的、接口 NIC driver 的不同而有所变化；因此，很可能在一个接口上你使用的是标准驱动，而在另外的接口上你使用的是 PF_RING-aware 驱动；当前，所有支持的 PF_RING-aware 驱动都保存在 `PF_RING/drivers` 目录下，除了 `TNAPI` ；需要注意的是，针对 DNA 驱动的情况，由于 kernel 被完全旁路掉了，故 `transparent_mode` 参数将不起作用；


Mode | Standard driver | PF_RING-aware driver | Packet Capture Acceleration
---|---|---|---
0 | Packets are received through Linux NAPI | Packets are received through Linux NAPI | Same as Vanilla Linux
1 | Packets are received through Linux NAPI | Packets are passed to NAPI (for sending them to PF_RING-unaware applications) and copied directly to PF_RING for PF_RING-aware applications (i.e. PF_RING does not need NAPI for receiving packets) | Limited
2 | The driver sends packets only to PF_RING so PF_RING-unaware applications do not see any packet | The driver copies packets directly to PF_RING only (i.e. NAPI does not receive any packet) | Extreme


