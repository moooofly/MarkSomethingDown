# Quick Start Guide

> 原文地址：[这里](https://github.com/ntop/PF_RING/blob/dev/userland/examples_zc/README.quickstart)

- 加载 PF_RING 内核模块

```shell
 $ cd PF_RING/kernel; make
 $ insmod pf_ring.ko
```

- 加载 PF_RING ZC-aware 驱动，用于 zero-copy Direct NIC Access ，以及任何其它工作在 1-copy 模式（zc 设备模拟）下的驱动；PF_RING ZC-aware 驱动位于 `drivers/PF_RING_aware` 目录下（带有 '`-zc`' 后缀）

```
 $ cd PF_RING/drivers/PF_RING_aware/intel/ixgbe/ixgbe-3.18.7-zc/src; make
 $ ./load_driver.sh
```

- 使能并挂在足够的 hugepages ，因为内存分配时需要（详情阅读 `PF_RING/README.hugepages`）；注意：hugepages 可以通过用来加载驱动程序的 `load_driver.sh` 脚本自动初始化；

```
 $ echo 1024 > /sys/kernel/mm/hugepages/hugepages-2048kB/nr_hugepages
 $ mkdir /dev/hugepages
 $ mount -t hugetlbfs nodev /dev/hugepages
```

- 运行众多 PF_RING ZC 示例应用程序之一；请为每一个应用程序使用一个唯一的 cluster id（除非你使用的时 multi-process 应用）；

```
 $ cd PF_RING/userland; make
 $ cd examples_zc
 $ ./zcount -i zc:eth1 -c 1
```

如果你的 eth1 和 eth2 时 cross connected 的，你可以执行

```
 $ ./zcount -i zc:eth1 -c 1
 $ ./zsend -i zc:eth2 -c 2
```

详情参考应用程序的 help (-h) 描述；
