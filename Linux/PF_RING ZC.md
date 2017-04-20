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

如果你省略了 'zc:' ，则对应的是基于 PF_RING 模式（no ZC）打开当前设备；

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
* ZC drivers 需要 hugepages 功能，在 `load_driver.sh` 脚本中给出了针对 hugepages 的配置处理；更多相关信息详见 **[README.hugepages](https://github.com/ntop/PF_RING/blob/dev/doc/README.hugepages.md)** 中的说明；

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
Virtual Machine 中以 zero-copy 方式进行 RX 和 TX 的 packets 转发（forward），而无需使用诸如 `PCIe passthrough` 这类技术；得益于 ZC devices 在 VMs 中的动态创建，你能够在 VM 中基于 zero-copy 方式捕获/发送 traffic ，而无需对 KVM 代码打 patch ，或在 ZC 设备被创建后才启动 KVM；本质上，你能够在 KVM 中达到 10 Gbit 线速能力，使用的是在物理主机中同样的命令，无需改变任何一行代码；

在 PF_RING ZC 中，即使面对的是 non-PF_RING aware drivers ，你依然能够使用 zero-copy 框架；这意味着你能够分发（dispatch）、处理（process）、发起（originate），和注入（inject）packets 到 zero-copy 框架中，即便其不是从 ZC devices 中发起的（originated）；

一旦 packet 被拷贝到（one-copy）ZC 世界中，从此该 packet 在整个生命周期中将一直按 zero-copy 方式被处理；例如在 `zbalance_ipc` 示例应用中，其以 1-copy 模式从 non-PF_RING aware 设备中读取 packet（例如 WiFI-device 或者 Broadcom NIC），之后在 ZC 内将 packet 以 zero-copy 操作进行发送；

------

# PF_RING ZC (Zero Copy)

> 原文地址：[这里](http://www.ntop.org/products/packet-capture/pf_ring/pf_ring-zc-zero-copy/)

## Multi-10 Gbit RX/TX Packet Processing from Hosts and Virtual Machines

PF_RING™ ZC (Zero Copy) 是一种灵活的 packet 处理框架，允许针对任意大小的 packet 达到 1/10 Gbit 线速包处理能力（ RX 和 TX）；其实现了 zero copy 操作，包括 inter-process 和 inter-VM (KVM) 通信模式；可以将其看作 DNA/LibZero 的后继者，基于过去几年的经验教训，提供了单一且一致的 API ；

其提供了一套干净灵活的 API ，实现了一组简单方便的构建模块（queue, worker 和 pool），可被用在线程、应用和 virtual machines 中；用以实现了 10 Gbit 线速 packet 处理能力；

## Simple And Clean API

PF_RING™ ZC 提供的简单 API 能够用来以寥寥几行代码创建出复杂的应用；下面的例子展示了如何基于 6 行代码创建出一个 aggregator+balancer 应用：

```c
 zc = pfring_zc_create_cluster(ID, MTU, MAX_BUFFERS, NULL);
 for (i = 0; i < num_devices; i++)
  inzq[i] = pfring_zc_open_device(zc, devices[i], rx_only);
 for (i = 0; i < num_slaves; i++)
   outzq[i] = pfring_zc_create_queue(zc, QUEUE_LEN);
 zw = pfring_zc_run_balancer(inzq, outzq, num_devices, num_slaves, NULL, NULL, !wait_for_packet, core_id);
```

关于 API 的更多信息，请参考 [documentation](http://www.ntop.org/pfring_api/pfring__zc_8h.html) 和相应的 [code examples](https://github.com/ntop/PF_RING/blob/dev/userland/examples_zc/README.examples) ；

## On-Demand Kernel Bypass with PF_RING Aware Drivers

PF_RING™ ZC 提供了新一代 PF_RING™ aware 驱动，能够使用 in-kernel 模式或 kernel bypass 模式；一旦成功安装，该驱动能够像标准 Linux 驱动一样完成常规网络工作（例如 ping 或 SSH）；当基于 PF_RING™ 完成上述功能时，会比 vanilla drivers 更快，因为其能够直接与 drivers 交互；如果你使用 PF_RING-aware 驱动以 zero copy 模式打开设备（例如 `pfcount -i zc:eth1`），该设备将变得对标准网络功能不可用，因为此时是以 zero-copy 模式通过 kernel bypass 进行访问对，正如同之前的 DNA 的行为；一旦应用完成（关闭）了设备的访问，标准网络活动又能正常使用了；

## Zero Copy Operations to Virtual Machines (KVM)

PF_RING™ ZC 允许你针对 KVM virtual machine 以 zero-copy 模式（针对 RX 和 TX）转发（forward）packets，而无需使用诸如 `PCIe passthrough` 这类技术；由于具备在 VMs 中动态创建 ZC 设备的能力，你能够从你的 VM 中以 zero-copy 模式捕获/发送 traffic 而无需为 KVM 代码打布丁，或者在创建 ZC 设备后启动 KVM ；重要的是，现在你已经能够在 KVM 中获得 10 Gbit 线速能力，但同样使用在物理主机中采用的命令，无需变更一行代码；

![ZC_IntraVM](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/ZC_IntraVM.png "ZC_IntraVM")

上图展示了基于 ZC 如何创建应用程序的 pipeline ，以 zero copy 方式跨 VMs 交流；重要的是，PF_RING™ ZC 从诞生的第一天就为云服务做好了准备（cloud-ready）；


## Zero Copy Operations

和前辈 LibZero 类似，基于 PF_RING™ ZC 你能够跨线程，应用和 VMs 进行 zero copy 操作；你能够在 zero-copy 模式下跨应用进行 **packets 均衡**；

![ZC_Balancing](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/ZC_Balancing.png "ZC_Balancing")

或实现 **packet fanout** 功能

![ZC_Fanout](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/ZC_Fanout.png "ZC_Fanout")

在 PF_RING™ ZC 中，所有事情都在 zero-copy ，线速能力下发生；

## Performance

和前辈 LibZero/DNA 类似，无论 packet 是来自物理主机，还是来自 KVM ，无论 packet 的大小，PF_RING™ ZC 都能保证达到 10 Gbit 线速能力；你可以使用 [demo applications](https://github.com/ntop/PF_RING/tree/dev/userland/examples_zc) 自行测试；

## Integrating Zero-Copy with One-Copy Devices

在 PF_RING™ ZC 中，你同样能够针对 non-PF_RING aware 驱动使用 zero-copy 框架；这意味着你能够分发（dispatch）、处理（process）、发起（originate），以及注入（inject）packets 到 zero-copy 框架中，尽管其并非起源于 ZC 设备；

![ZC_OneCopy](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/ZC_OneCopy.png "ZC_OneCopy")

一旦 packet 被拷贝进入（one-copy）到 ZC 世界，从此 packet 在其整个生命期内将按照 zero-copy 方式被处理；例如，[zbalance_ipc](https://github.com/ntop/PF_RING/blob/dev/userland/examples_zc/zbalance_ipc.c) 示例应用能够基于 1-copy 方式从 non-PF_RING aware 设备中读取 packet（例如，WiFI-device 或 Broadcom NIC），之后在 ZC 框架内部发送时则基于 zero-copy 操作；

## Kernel Bypass and IP Stack Packet Injection

与其它 kernel-bypass 技术相比，基于 PF_RING™ ZC 的你，能够在任意时刻决定哪些 packets 以 kernel-bypass 方式接收，再（重新）注入（inject）到标准 Linux IP stack 中；PF_RING 现已提供了一个称作 “stack” 的 IP [stack packet injection module](https://github.com/ntop/PF_RING/blob/dev/userland/examples/README.stackinjection)，允许你选择哪些通过 kernel-bypass 接收的 packets 需要再被注入到标准 IP stack 中；你所需要做的仅仅是打开设备 “`stack:ethX`” 并进行 packets 发送，以便将其推入 IP stack 中，就好像这些 packets 是从 ethX 上接收的一样；

## DAQ for Snort

[Snort](https://www.snort.org/) 用户同样能够受益于 PF_RING™ ZC 的速度能力（最受欢迎的 IDS/IPS 之一）；本地化（native）后的 PF_RING™ ZC DAQ (Snort Data AcQuisition) 库要比标准 [PF_RING™ DAQ](https://github.com/ntop/PF_RING/tree/dev/userland/snort/) 快 [20% 到 50%](http://www.ntop.org/wp-content/uploads/2012/09/Snort_over_DNA_Silicom_30_07_2012_1.pdf) 左右，并且其可在 IPS 和 IDS 模式下运行；

PF_RING™ ZC DAQ 属于 [PF_RING™](https://github.com/ntop/PF_RING/tree/dev/userland/snort/) 的一部分；


|  | e1000e | igb | ixgbe | i40e
---|---|---|---|---
Capture Rate (Line-Rate) | 1 Gbit/s | 1 Gbit/s | 10 Gbit/sec | 40 Gbit/sec
Supported Cards | Intel 8254x/8256x/8257x/8258x-based | Intel 82575/82576/82580/I350-based | Intel 82599/X540/X710-based | Intel XL710
Operating System | Linux (kernel 2.6.32 or better) | Linux (kernel 2.6.32 or better) | Linux (kernel 2.6.32 or better) | Linux (kernel 2.6.32 or better)
Traffic Reception | included | included | included | included
Traffic Injection | included | included | included | included
Hw packet filtering | | | Intel 82599-based only | 
Hw timestamping (nsec) | | Intel 82580/I350-based only | 

