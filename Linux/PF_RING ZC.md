# PF_RING ZC

> 原文地址：[这里](https://github.com/ntop/PF_RING/blob/dev/doc/README.ZC.md)

PF_RING 既可以工作在标准 NIC drivers 之上，也可以工作在特定的 drivers 之上；这一点对于 PF_RING 内核模块来说也是一样的，但对于不同的 drivers 来说，一些功能和性能表现会有所不同；

对于那些需要最大化 packet 捕获速度，同时希望在拷贝 packets 到 host 时不耗费任何（0%）CPU 利用率的用户来说（即不使用 NAPI polling 机制），可以使用 ZC drivers（或者称作新一代 DNA）；在 zero-copy 模式下，其允许直接从网络接口上进行 packets 读取，既绕过了 Linux kernel ，也绕过了 PF_RING 模块；

在 ZC 中，RX 和 TX 操作均支持；由于 kernel 被绕过了，一部分 PF_RING 功能将会缺失，其中包括在 kernel 层面中支持的 packet filtering 功能（即 BPF 和 PF_RING filters）；

在使用了 PF_RING ZC 后，你能够在任意 packet size 下获得 1/10G 的线速，能够创建 inter-process 和 inter-VM 的 clusters（PF_RING ZC 不只是 driver ，其还提供了简单但强有力的 API）；其可以被当作 `DNA/LibZero` 的后继者，但在汲取了过去几年的经验教训后，早已完成 API 的单一和一致保证；

用于进行测试的样例应用程序在 `userland/examples_zc` 中可以找到；

位于 `PF_RING/drivers/` 中的这些 drivers 都是支持 PF_RING ZC 库的标准 drivers ；其均能被用作标准内核 drivers ，或用在 zero-copy 
kernel-bypass 模式中（使用 PF_RING ZC 库），只需添加前缀 "zc:" 到接口名字上； 

一旦进行了安装，这些 drivers 就会如同标准 Linux drivers 一样允许进行常规的网络操作（例如 `ping` 或 `SSH`）；如果你是以 zero copy 方式指定 "zc:" 前缀打开的对应设备，则该设备对于标准网络操作将不可用，因为其是基于绕过 kernel 的方式，按照 zero-copy 方式被访问的，正如以前 DNA 的行为方式；一旦访问该设备的应用被关闭，标准网络功能将重新被激活；处于 ZC 模式的接口能够提供与 DNA 类似的性能；

例如：

```shell
pfcount -i zc:eth0
```

如果你省略了 'zc:' ，则对应的是基于 PF_RING 模式（无 ZC）打开当前设备；

## Supported Cards

为了能够利用 ZC 功能，你需要一个带有 ZC 支持的 PF_RING aware driver ，可以通过是否具有 '-zc' 后缀来辨识；当前存在三种 driver families 可用：

1 Gbit

- e1000e (RX and TX)
- igb    (RX and TX)

10 Gbit

- ixgbe (RX and TX)

10/40 Gbit

- i40e (RX and TX)

10/40/100 Gbit

- fm10k (RX and TX)

这些 drivers 可在 `drivers/` 中找到；

需要注意的是：

* PF_RING 内核模块必须在 ZC driver 之前被加载；
* 为了确保正确配置目标设备，强烈建议使用配套 drivers 的 `load_driver.sh` 脚本（可以基于该脚本进一步精细化调优）；
* ZC drivers 需要 hugepages 功能，在 `load_driver.sh` 脚本中给出了针对 hugepages 的配置处理；更多相关信息详见 **README.hugepages** 中的说明；

加载 PF_RING 和 ixgbe-ZC driver 的示例：

```shell
cd <PF_RING PATH>/kernel
insmod pf_ring.ko
cd PF_RING/drivers/intel/ixgbe/ixgbe-X.X.X-zc/src
make
./load_driver.sh
```

## ZC API

PF_RING ZC (Zero Copy) 是一套灵活的 packet 处理框架，允许你获得 1/10 Gbit 的线速 packet 处理能力（支持 RX 和 TX），无论 packet 的大小；

其实现了 zero-copy 操作，可用在 inter-process 和 inter-VM (KVM) 通信中；其可以被当做 `DNA/LibZero` 的后继者，提供了单一和一致的 API 并籍此实现了简单的构件块（queue, worker 和 pool），以便在 threads, applications 和 virtual machines 中使用；

下面的例子展示了如何基于 6 行代码创建一个 aggregator+balancer 应用程序：

```c
zc = pfring_zc_create_cluster(ID, MTU, MAX_BUFFERS, NULL);
for (i = 0; i < num_devices; i++)
  inzq[i] = pfring_zc_open_device(zc, devices[i], rx_only);
for (i = 0; i < num_slaves; i++)
  outzq[i] = pfring_zc_create_queue(zc, QUEUE_LEN);
zw = pfring_zc_run_balancer(inzq, outzq, num_devices, num_slaves, NULL, NULL, !wait_for_packet, core_id);
```

PF_RING ZC 允许你以在 KVM 
Virtual Machine 中以 zero-copy 方式进行 RX 和 TX 的 packets 转发（forward），而无需使用诸如 PCIe passthrough 这类技术；得益于 ZC devices 在 VMs 中的动态创建，你能够在 VM 中基于 zero-copy 方式捕获/发送 traffic ，而无需对 KVM 代码打 patch ，或在 ZC 设备被创建后才启动 KVM；本质上，你能够在 KVM 中达到 10 Gbit 线速能力，使用的是在物理主机中同样的命令，无需改变任何一行代码；

在 PF_RING ZC 中，即使面对的是 non-PF_RING aware drivers ，你依然能够使用 zero-copy 框架；这意味着你能够分发（dispatch）、处理（process）、发起（originate），和注入（inject）packets 到 zero-copy 框架中，即便其不是从 ZC devices 中发起的（originated）；

一旦 packet 被拷贝到（one-copy）ZC 世界中，从此该 packet 在整个生命周期中将一直按 zero-copy 方式被处理；例如在 `zbalance_ipc` 示例应用中，其以 1-copy 模式从 non-PF_RING aware 设备中读取 packet（例如 WiFI-device 或者 Broadcom NIC），之后在 ZC 内将 packet 以 zero-copy 操作进行发送；
