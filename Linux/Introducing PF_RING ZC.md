# Introducing PF_RING ZC (Zero Copy)

> 原文地址：[这里](http://www.ntop.org/pf_ring/introducing-pf_ring-zc-zero-copy/)

在将近 18 个月的开发后，我们高兴的宣布 [PF_RING ZC (Zero Copy)](http://www.ntop.org/products/packet-capture/pf_ring/pf_ring-zc-zero-copy/) 发布了；基于从 DNA 和 libzero 处学到的经验教训，我们决定从零开始重新设计一个全新的、提供有一致性保证的（consistent）、zero-copy API ，用以实现广受欢迎的网络模式；设计目标是能够提供一套简单 API 用以保证网络应用程序开发者能够获得线速性能（从 1 到 multi-10 Gbit）；我们已经帮你隐藏了全部内部和底层细节，以打造 developer-centric API 而非 network/hardware-centric API ；

对于熟悉 DNA 和 Libzero 到用户来说，ZC 和它们的主要差别在于：

- 我们统一了 in-kernel（之前的 PF_RING-aware drivers）和 kernel-bypass（之前的 DNA）驱动的实现；现在，你可以通过 “`-i eth0`”（in-kernel 处理模式）和 “`-i zc:eth0`”（kernel bypass）打开相同的设备；尤其是，你可以运行时决定使用哪种操作模式；
- 全部的 drivers memory 被映射为（mapped）[huge-pages](https://github.com/ntop/PF_RING/blob/dev/doc/README.hugepages.md) 以后取更好性能；
- 如果你通过 “`zc:ethX`” 打开设备，则所有操作都以 zero-copy 模式完成；你能通过简单测试看到具体效果（eth3 为一个 10 Gbit 接口，运行了 PF_RING-aware [ixgbe](https://github.com/ntop/PF_RING/tree/dev/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src) 驱动）；第一个命令发送可达 0.82 Gbit ，而第二个可达 10 Gbit；
```
# ./zsend -i eth3 -c 3 -g 0
# ./zsend -i zc:eth3 -c 3 -g 0
```
- ZC 具有 [KVM](https://en.wikipedia.org/wiki/Kernel-based_Virtual_Machine) 友好性，意味着你能够从运行在 VM 中的应用和线程上以 10 Gbit 线速发送/接收 packets ，而无需使用类似 `PCIe passthrough` 这种技术；需要注意的是，同一种应用可以不做人和修改的运行在 VM 和物理机上：我们真的真的让 VM 和 ZC 之间成为了好基友；
- 与其它网络框架类似，例如 [Click Modular Router](http://read.cs.ucla.edu/click/click) ，我们提供了诸如 queue, pool, worker 这类简单组件用于保证通过几行代码构建应用程序；
- API 已经被简化和重写过了；例如，仅使用 6 行代码你就能够创建一个 traffic aggregator 和 balancer （详见[示例程序](https://github.com/ntop/PF_RING/blob/dev/userland/examples_zc/README.examples)）；
- 当工作在 kernel-bypass 模式下时，我们允许你和 IP stack 进行交互，并向其发送 packets 或从其获取 packets ；这可以简化既需要满足线速要求，又（有时）需要和主机 IP stack 交互的应用程序的开发过程；
- 当和低速（low-speed）设备进行交互时（例如 WiFi 适配器），我们能够采用 one-copy-mode 方式进行操作，或者允许未被提供 accelerate drivers 的 NICs 同样从 ZC 中受益：当 packet 进入 ZC 时，仅存在一次拷贝行为，然后没有然后了；
- 针对 no-for-profit 群众和研究机构，我们免费提供产品供使用；针对商业客户，我们简化了 license 模型，因此仅需要一个 license（不再是之前的两个，一个针对 DNA ，一个针对 Libzero ，因此现在您可以享受更低的 license 费用）；
- 我们期望这项技术能够普及起来，因此你只会在使用 accelerated drivers 时需要 license ，而对于所有其它东东（例如 KVM 和 one-copy 支持），你都能够免费使用，无需 license ；

如果你想要测试 PF_RING ZC ，最好的方式是：阅读这个 [quick start guide](https://github.com/ntop/PF_RING/blob/dev/userland/examples_zc/README.quickstart) ；


## FAQ

Q. **PF_RING ZC 的性能如何？**
A. 线速，任意 packet 大小，multi 10 Gbit ；你可以使用各种 apps 进行测试，例如 `zcount` 和 `zsend` ；

Q. **是否仍旧对 DNA 和 Libzero 进行支持？**
A. 就目前而言，我们会继续支持 DNA/Libzero ，尽管未来会是 PF_RING ZC 的天下，因为后者提供了许多新特性和一致的 API ；

Q. **PF_RING ZC 是否会支持 legacy applications ，例如那些基于 libpcap 开发的应用？**
A. 会的，一切和其前任 DNA 一样，我们支持 pcap-based 应用程序，以及诸如基于 PF_RING DAQ 实现的 Snort 等其它 apps ；

Q. **我应该如何打开一个 PF_RING ZC/libpcap 中的 queue ？**
A. 你可以使用如下语法实现：“`zc:<clusterId>@<queueId>`” ；例如：`pcount -i zc:4@0`

Q. **What adapters feature native ZC drivers?**
A. 我们当前支持了 Intel 适配器（1 和 10 Gbit）使用 zero-copy 模式，以及所有基于 1-copy 模式的其它适配器；需要记住的是，一旦 packets 被搬移到了 ZC 的世界，你就能够将它们以 zero-copy 模式传递给任意数量的应用程序、线程和 VMs 了；重点是，你只需在入口处支付 "copy ticket" ；

Q. **How do you position ZC with respect to other technologies such as DPDK?**
A. DPDK 主要面向的是贴近硬件的应用开发者，即对 X86 架构细节很熟悉的人，能够并愿意使用十六进制的 PCI ID （例如 0d:01.0）调用网络接口；对于 PF_RING ZC 的使用者来说，设备调用可以通过其名字（例如 eth1）完成，PF_RING ZC 负责管理所有底层（low-level）细节信息，因此你就可以无阻碍的运行你已有的 pcap-based 应用；ZC 是一种开发者友好的技术；
