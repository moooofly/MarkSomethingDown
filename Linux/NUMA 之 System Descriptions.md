# NUMA 之 System Descriptions

标签（空格分隔）： linux

---

## [System Descriptions](http://lse.sourceforge.net/numa/faq/system_descriptions.html)

本节内容将描述各种平台相关信息；其中一些平台当前可能是存在的，而另外一些可能是虚构的，或者说将来可能会被制造出来；本文的**目的**就是为了展示**系统拓扑的潜在的多样性**；本文中描述的系统包括：

- 典型的 `SMP` 系统
- Alpha Wildfire 系统
- IBM Numa-Q 系统
- SGI Mips64 系统
- 能够在单块芯片上使用多 CPU 技术的系统
- 将 CPUs 和 Memory 环状连接在一起的系统


### A typical SMP system

下图展示了一个典型的 SMP 系统设计，其使用了 Intel x86 处理器：

![sys_smp](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/sys_smp.gif "sys_smp")

典型的 SMP 系统包括多个 CPUs ；典型 CPU 会包含一个 L1 cache ；L2 cache 在典型情况下由 CPU 所管理，但 L2 cache 的内存对于 CPU 来说是外部的；该系统可能还具有 L3 cache ，针对它的管理对于 CPU 来说是外部的；L3 cache 可能会被多个 CPUs 所共享；系统中还会包含主存（main memory），并且主存中的内容可能出现在任一个 caches 中；针对主存和各级 caches 的一致性维护需要由硬件来保证；典型的 memory latencies 为：

- L1 cache hit:
- L2 cache hit:
- L3 cache hit:
- memory access:

系统中还会包含一条或多条 IO 总线，IO 控制器会关联到这些 IO 总线上，所有设备会关联到这些 IO 控制器上；

### Compaq / Alpha Wildfire

Currently searching for more information. Any help identifying this information or a volumteer to write this section would be greatly appreciated. Please send any information to Paul Dorwin (pdorwin@us.ibm.com)

### The IBM Numa-Q system

IBM Numa-Q 的系统设计如下图所描述：

![sys_numaq](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/sys_numaq.gif "sys_numaq")

系统中的每一个 node 都由一个简单的 4 处理器 SMP 系统构成；Numa-Q 使用的 Intel X86 CPU ；该 node 上的每一个 CPU 都包含 L1 和 L2 cache ；该 node 中还回包含一个 L3 cache 用于其上所有处理器共享；单个 node 上的本地内存最多支持 2 Gb ；该 node 同样包含一个（不知道多少大小的）远端 cache ，用于缓存来自远端 nodes 的数据；Numa-Q 系统中的 nodes 是通过 Lynxer 进行互联的，其中包含了 SCI (Scalable Coherent Interface) 接口；memory latencies 如下：

- L1 cache hit:
- L2 cache hit:
- L3 cache hit:
- local memory hit:
- remote cache hit:
- remote memory hit:

Linux 系统针对 Numa-Q 的移植（The Linux port to Numa-Q）要求修改基于 APIC 访问 CPUs 的方式；默认情况下，CPU 上的 APIC 地址是平坦的（flat）, 允许最多 8 个处理器同时使用系统总线；Numa-Q 使用 `luster` 模式，其中 8 比特可以进一步拆分；4 比特用来确认（区分）最多 16 个 nodes ，而另外 4 比特用于确认（区分）每个 node 上的 4 个 cpus ；而 Lynxer card 负责确认和路由跨 node 的访问；每一个 node 同样包含 2 个 PCI 总线；（在撰写此文时）第一个 PCI 总线包含 3 个 slots 而第二个包含 4 slots (VERIFY THIS) ；Linux was verified to boot consistently on 4 nodes containing a total of 16 processors. The Numa-Q work is hidden behind CONFIG_MULTIQUAD, and the patches are being tested. Work is underway to update the kernel to allow IO boards in all nodes. Work is also underway to port discontiguous memory to the Numa-Q platform.

### Slicon Graphics Mips64

Currently searching for more information. Any help identifying this information or a volumteer to write this section would be greatly appreciated. Please send any information to Paul Dorwin (pdorwin@us.ibm.com)


### CPU/Memory ring system

下图展示了一个理论上的系统：一个 CPU 连接到两个独立的 memory 节点上，并且每一个 memory 节点也会连接到两个 CPUs 上：

![sys_ring](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/sys_ring.gif "sys_ring")

正如你所看到的，这样就会产生一个 ring 配置的系统；

### Multiple CPU on a single chip

下图展示了理论上单块芯片包含多个 CPU 的情况：

![sys_mpct](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/sys_mpct.gif "sys_mpct")

### Conclusions

在这里所描述的所有系统均展示了相似的特点：均含有多个处理器和一些配套的 cache ；也都包含一些形式的 memory 和 IO ；处理器之间彼此被分开一段距离；然而针对他们的连接方式存在比较大的差别；为了重申这一点，它们均具有独一无二的拓扑结构；Therefore, it is important to present an in-kernel infrastructure which easily allows the architecture dependant code to create an in-kernel description of the system's topology. The section describing the in-kernel infrastructure will provide the details of the proposed solution. 下面是一组特性列表，可用于从系统拓扑中得到确认：

- How many **processors** are in the system.
- How many **memories** are in the system.
- How many **nodes** are in the system
- What is **encapsulated** in any given node.
- For any given processor:
    - What is the **distance to any other processor** in the system.
    - What is the **distance to any memory** in the system.
    - How many **cache levels** exist and how large is each cache.
- For any given memory:
    - What is the **start pfn** of the memory
    - What is the **size** of the memory
    - What **processors are directly connected** to the memory.