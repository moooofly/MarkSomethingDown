
PF_RING 由 Linux 内核模块和用户空间框架构成，允许包处理相关应用程序基于一致性 API 高速地处理收到的网络包； 

源码目录说明如下：

```
drivers                     PF_RING optimized drivers
drivers/                    PF_RING-aware/ZC drivers (suggested option)
userland/			        用户空间代码
userland/lib/			    用户空间库
userland/libpcap-XXX-ring/	支持 PF_RING 的 Libpcap
userland/examples/		    使用 PF_RING 的样例应用程序
userland/examples_zc/		使用 PF_RING ZC 样例应用程序
userland/snort/			    Snort DAQ module for using snort over PF_RING
```

PF_RING Drivers Models
----------------------

- **`PF_RING-aware/ZC drivers`**

这些 drivers 被称作 "aware" 的原因在于：针对 PF_RING 进行了优化；

存在两种使用方式：

**作为标准 drivers 使用（packets 仍会进入 Linux stack）**，例如： 

```
pfcount -i eth1
```

或者**按 zero copy 模式使用，以便在 RX 和 TX 过程中完全旁路（bypassing）掉 Linux stack ，以达到线速**；

只要一个应用程序以 zero copy (ZC) 模式打开了接口，那么该接口（尽管仍可通过 `ifconfig` 进行查看）将无法被用于标准网络功能（例如 ping 或 SSH）；这种情况开始于基于 ZC 模式打开接口之时，终止于接口关闭之时，之后立即可以重新作为标准 Linux 网络功能使用；

以 ZC 模式打开接口，需要使用 'zc:' 作为接口名前缀，例如：

```
pfcount -i zc:eth1
```

- **`ZC drivers`**

这些 drivers 被用于实现内核旁路功能，直到 PF_RING 5.x 和 ZC 被引入之前；这些 drivers 还会持续可用一段时间，但在 ZC 功能可用后会被废弃，因为后者提供了相同的性能和更高的灵活性；在将来 PF_RING 发布后，我们将不再维护这些 drivers ；

PF_RING from Virtual Machines (KVM)
-----------------------------------

在使用了 PF_RING ZC 功能后，即使在 KVM 虚拟机中，你同样可以成功获得加速功能，而无需使用诸如 PCIe bypass 这类技术；这意味着 VM 能够动态的按 ZC 模式打开网络接口，实现复杂的包处理拓扑，并保持在 10 Gbit 线速；关于如何在 KVM 中利用 PF_RING ZC 请参考 `userland/examples_zc/README.kvm` 中的内容；

在将来，我们还会将其移植到其它 hypervisors 上，但当前只面向 KVM ；

Compilation
-----------

首先安装一些基本的编译工具和库；

在 Ubuntu 系统中，可以执行

```
# apt-get install build-essential bison flex linux-headers-$(uname -r) libnuma-dev
```

之后进行编译

```
# make
```

Installation
------------

```
# sudo su
# cd kernel; make install
# cd ../userland/lib; make install
```

实际

```
[root@xg-esm-data-4 kernel]# make install
mkdir -p /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring
cp *.ko /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring
mkdir -p /usr/include/linux
cp linux/pf_ring.h /usr/include/linux
/sbin/depmod 3.10.0-229.el7.x86_64
[root@xg-esm-data-4 kernel]#
```


Testing PF_RING
---------------

在 `PF_RING/userland/examples` 中，已经包含了大量可用于测试 PF_RING 的示例应用程序，尤其

- `pfcount` 允许你进行 packets 捕获
- `pfsend` 可用于重放（或生成）traffic

你可以参考这些应用程序的源码文件，学习如何使用 PF_RING API ；

需要注意的是：在运行任何应用程序之前，你需要首先加载 `pf_ring.ko` 这个内核模块：

```
# sudo su
# insmod ./kernel/pf_ring.ko
```

Documentation
-------------

如果你需要一个更全面的文档，你可以参考保存在 `doc/` 目录下的手册内容；在该目录中，你能够找到 PF_RING 支持的、针对各种网络适配器的独立 README 文件；

可以通过 doxygen 生成 API 文档，只需运行

```
# make documentation
```

最终生成到 `doc/html` 中；
