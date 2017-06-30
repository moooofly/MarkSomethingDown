# CPU 隔离之 numactl

## NUMA 的 wiki 说明

原文地址：[Non-uniform memory access](https://en.wikipedia.org/wiki/Non-uniform_memory_access)

> 主要内容如下：
>
> - NUMA 是一种用于**多核处理**的计算机**内存设计架构**；
> - 能否从 NUMA 中获得收益**与工作负载的类型密切相关**；
> - NUMA 架构逻辑上源于**对 SMP 架构的扩展**；
> - 在一些场景中，NUMA 系统会与软硬件配合**进行数据的内存搬移**；

`Non-uniform memory access (NUMA)` is a computer **memory design** used in multiprocessing, where the memory access time depends on the memory location relative to the processor. Under NUMA, a processor can access its own local memory faster than non-local memory (memory local to another processor or memory shared between processors). **The benefits of NUMA are limited to particular workloads**, notably on servers where the data are often associated strongly with certain tasks or users.

NUMA architectures logically follow in scaling from `symmetric multiprocessing (SMP)` architectures.

Modern CPUs operate considerably faster than the main memory they use. In the early days of computing and data processing, the CPU generally ran slower than its own memory. The performance lines of processors and memory crossed in the 1960s with the advent of the first supercomputers. Since then, CPUs increasingly have found themselves "starved for data" and having to stall while waiting for data to arrive from memory. Many supercomputer designs of the 1980s and 1990s focused on **providing high-speed memory access** as opposed to faster processors, allowing the computers to work on large data sets at speeds other systems could not approach.

**Limiting the number of memory accesses** provided the key to extracting high performance from a modern computer. For commodity processors, this meant **installing an ever-increasing amount of high-speed cache memory and using increasingly sophisticated algorithms** to avoid cache misses. But the **dramatic increase in size of the operating systems and of the applications run on them** has generally overwhelmed these cache-processing improvements. Multi-processor systems without NUMA make the problem considerably worse. Now a system can starve several processors at the same time, notably because only one processor can access the computer's memory at a time.

**NUMA attempts to address this problem by providing separate memory for each processor, avoiding the performance hit when several processors attempt to address the same memory**. For problems involving spread data (common for servers and similar applications), NUMA can improve the performance over a single shared memory by a factor of roughly the number of processors (or separate memory banks). Another approach to addressing this problem, used mainly in **non-NUMA** systems, is the **multi-channel memory architecture**, in which a linear increase in the number of memory channels increases the memory access concurrency linearly.

Of course, not all data ends up confined to a single task, which means that more than one processor may require the same data. To handle these cases, NUMA systems include additional hardware or software to move data between memory banks. This operation slows the processors attached to those banks, so **the overall speed increase due to NUMA depends heavily on the nature of the running tasks**.


----------


## [FAQ about the NUMA architecture](http://lse.sourceforge.net/numa/faq/)

> 中文翻译：[这里](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/FAQ%20about%20the%20NUMA%20architecture.md)


----------


## [Linux 的 NUMA 技术](http://www.ibm.com/developerworks/cn/linux/l-numa/index.html)

> 本文写于 2004 年；部分内容有调整；

随着科学计算、事务处理对计算机性能要求的不断提高，**SMP（对称多处理器）**系统的应用越来越广泛，规模也越来越大，但由于传统的 SMP 系统中，所有处理器都共享系统总线，因此当处理器的数目增大时，系统总线的竞争冲突加大，系统总线将成为瓶颈，所以**目前 SMP 系统的 CPU 数目一般只有数十个**，可扩展能力受到极大限制。NUMA 技术有效结合了 SMP 系统易编程性和 **MPP（大规模并行）**系统易扩展性的特点，较好解决了 SMP 系统的可扩展性问题，已成为当今高性能服务器的主流体系结构之一。目前国外著名的服务器厂商都先后推出了基于 NUMA 架构的高性能服务器，如 HP 的 Superdome、SGI 的 Altix 3000、IBM 的 x440、NEC 的 TX7、AMD 的 Opteron 等。随着 Linux 在服务器平台上的表现越来越成熟，Linux 内核对 NUMA 架构的支持也越来越完善，特别是**从 2.5 开始**，Linux 在调度器、存储管理、用户级 API 等方面进行了大量的 NUMA 优化工作，目前这部分工作还在不断地改进，如新近推出的 2.6.7-RC1 内核中增加了 NUMA 调度器。

NUMA 系统是由多个节点通过高速互连网络连接而成的；

NUMA 系统的节点（node）通常是由一组 CPU 和本地内存组成，有的节点可能还有 I/O 子系统。由于每个节点都有自己的本地内存，因此全系统的内存在物理上是分布的，每个节点访问本地内存和访问其它节点的远地内存的延迟是不同的，为了减少非一致性访存对系统的影响，在硬件设计时应尽量降低远地内存访存延迟（如通过 Cache 一致性设计等），而操作系统也必须能感知硬件的拓扑结构，优化系统的访存。

目前 IA64 Linux 所支持的 NUMA 架构服务器的物理拓扑描述是通过 **ACPI（Advanced Configuration and Power Interface**）实现的。ACPI 是由 Compaq、Intel、Microsoft、Phoenix 和 Toshiba 联合制定的 BIOS 规范，它定义了一个非常广泛的配置和电源管理；ACPI 规范也已广泛应用于 IA-32 架构的至强服务器系统中。

针对 NUMA 系统的物理内存分布信息，Linux 是从系统 firmware 的 ACPI 表中获得的，最重要的是 **SRAT（System Resource Affinity Table）**和 **SLIT（System Locality Information Table）**表，其中 SRAT 包含两个结构：

- Processor Local APIC/SAPIC Affinity Structure：记录某个 CPU 的信息；
- Memory Affinity Structure：记录内存的信息；

SLIT 表则记录了各个节点之间的距离；

Linux 采用 Node、Zone 和页三级结构来描述物理内存的：

![Linux中Node、Zone和页的关系](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Linux%E4%B8%ADNode%E3%80%81Zone%E5%92%8C%E9%A1%B5%E7%9A%84%E5%85%B3%E7%B3%BB.gif "Linux中Node、Zone和页的关系")

NUMA 系统中，由于局部内存的访存延迟低于远地内存访存延迟，因此将进程分配到局部内存附近的处理器上可极大优化应用程序的性能。Linux 2.4 内核中的调度器由于只设计了一个运行队列，可扩展性较差，在 SMP 平台表现一直不理想。当运行的任务数较多时，多个 CPU 增加了系统资源的竞争，限制了负载的吞吐率。在 2.5 内核开发时，Ingo Molnar 写了一个多队列调度器，称为 **O(1)**，从 2.5.2 开始 O(1) 调度器已集成到 2.5 内核版本中。**O(1) 是多队列调度器**，每个处理器都有一条自己的运行队列，但由于 O(1) 调度器不能较好地感知 NUMA 系统中节点这层结构，从而不能保证在调度后该进程仍运行在同一个节点上，为此，Eirch Focht 开发了节点亲和的 NUMA 调度器，它是建立在 Ingo Molnar 的 O(1) 调度器基础上的，Eirch 将该调度器向后移植到 2.4.X 内核中，该调度器最初是为基于 IA64 的 NUMA 机器的 2.4 内核开发的，后来 Matt Dobson 将它移植到基于 X86 的 NUMA-Q 硬件上。

在每个任务创建时都会赋予一个 HOME 节点（所谓 **HOME 节点**，就是该任务获得最初内存分配的节点），它是当时创建该任务时全系统负载最轻的节点，由于目前 Linux 中不支持任务的内存从一个节点迁移到另一个节点，因此在该任务的生命期内 HOME 节点保持不变。一个任务最初的负载平衡工作（也就是选该任务的 HOME 节点）缺省情况下是由 `exec()` 系统调用完成的，也可以由 `fork()` 系统调用完成。在任务结构中的 `node_policy` 域决定了最初的负载平衡选择方式；

在节点内，该 NUMA 调度器如同 O(1) 调度器一样。在一个空闲处理器上的动态负载平衡是由每隔 **1ms** 的时钟中断触发的，它试图寻找一个高负载的处理器，并将该处理器上的任务迁移到空闲处理器上。在一个负载较重的节点，则每隔 **200ms** 触发一次。调度器只搜索本节点内的处理器，只有还没有运行的任务可以从 Cache 池中移动到其它空闲的处理器。

如果本节点的负载均衡已经非常好，则计算其它节点的负载情况。如果某个节点的负载超过本节点的 **25％** ，则选择该节点进行负载均衡。如果本地节点具有平均的负载，则延迟该节点的任务迁移；如果负载非常差，则延迟的时间非常短，延迟时间长短依赖于系统的拓扑结构。


----------


## [NUMA (Non-Uniform Memory Access): An Overview](https://queue.acm.org/detail.cfm?id=2513149)

> TODO


----------


## NUMA 和 SMP

> **NUMA 和 SMP 是两种 CPU 相关的硬件架构**。

现在的机器上都是有多个 CPU 和多个内存块的；而很久以前我们都是将内存块看成是一大块内存，所有 CPU 到这个共享内存的访问都是一样的，这就是之前普遍使用的 **SMP (Symmetric Multi-Processor) 模型**，即对称多处理器结构；在 SMP 架构里面，所有的 CPU 通过争用一个总线来访问所有内存，优点是资源共享，而缺点是总线争用激烈。随着 PC 服务器上的 CPU 数量变多（不仅仅是 CPU 核数），通过总线争用来访问共享内存的弊端慢慢越来越明显，于是 Intel 在 Nehalem CPU 上推出了 NUMA 架构，而 AMD 也推出了基于相同架构的 Opteron CPU ；**NUMA（Non-Uniform Memory Access）模型**就是这样的环境下引入，即非一致存储访问结构。

NUMA 最大的特点是引入了 **node** 和 **distance** 概念。对于 CPU 和内存这两种最宝贵的硬件资源，NUMA 用近乎严格的方式划分了所属的**资源组（node）**，而每个资源组内的 CPU 和内存是几乎相等的；资源组的数量取决于物理 CPU 的个数（现有的 PC server 大多数有两个物理 CPU ，每个 CPU 有 4 个核）；distance 这个概念是用来定义各个 node 之间调用资源开销的，为资源调度优化算法提供数据支持。

由此可知，SMP 访问内存的都是代价都是一样的，但是在 NUMA 架构下，本地内存的访问和非本地内存的访问代价是不一样的；

比如一台机器上有 2 个处理器 4 个内存块。我们将 1 个处理器和两个内存块合起来，称为一个 `NUMA node` ，这样这个机器就会有两个 NUMA node 。**在物理分布上，NUMA node 的处理器和内存块的物理距离更小，因此访问也更快**。比如这台机器会分左右两个处理器（cpu1, cpu2），在每个处理器两边放两个内存块(memory1.1, memory1.2, memory2.1, memory2.2)，这样 NUMA node1 的 cpu1 访问 memory1.1 和 memory1.2 就比访问 memory2.1 和 memory2.2 更快。所以，**在使用 NUMA 模式时，如果能尽量保证本 node 内的 CPU 只访问本 node 内的内存块，那这样的效率就是最高的**。

### NUMA 策略

- 每个进程（或线程）都会从父进程继承 NUMA 策略，并分配有一个优先（prefered） node ；如果 NUMA 策略允许的话，进程可以调用其他 node 上的资源；
- NUMA 的 CPU 分配策略有 `cpunodebind` 和 `physcpubind` ；**`cpunodebind` 规定进程运行在某几个 node 之上，而 `physcpubind` 可以更加精细地规定运行在哪些核上**；
- NUMA 的内存分配策略有 `localalloc`、`preferred`、`membind` 和 `interleave` 四种；
    - `localalloc` 规定进程仅从当前 node 上请求分配内存；
    - `preferred` 则比较宽松地指定了一个推荐的 node 来获取内存；如果被推荐的 node 上没有足够内存，进程可以尝试别的 node ；
    - `membind` 可以指定若干个 node ，进程只能从这些指定的 node 上请求分配内存；
    - `interleave` 规定进程从指定的若干个 node 上以 RR（Round Robin 轮询调度）算法交织地请求分配内存；

因为 NUMA 默认的内存分配策略是优先在进程所在 CPU 的本地内存中分配，会导致 CPU 节点之间内存分配不均衡，当某个 CPU 节点的内存不足时，会导致 swap 产生，而不是从远程节点分配内存。这就是所谓的 **swap insanity** 现象。

### NUMA 和 swap

可能大家已经发现了，NUMA 的内存分配策略对于不同进程（或线程）来说，并不是公平的。在现有的 Redhat Linux 中，`localalloc` 是默认的 NUMA 内存分配策略，这个配置选项（有可能）导致资源独占程序很容易将某个 node 的内存用尽。而当某个 node 的内存耗尽时，Linux 又刚好将这个 node 分配给了某个需要消耗大量内存的进程（或线程），swap 就妥妥地产生了，尽管此时可能还有很多 page cache 可以释放，甚至还有很多的 free 内存。

在运行程序的时候，使用 `numactl -m` 和 `--physcpubind` 就能指定将这个程序使用哪个 memory 以及运行在哪个 cpu 中。[玩转 cpu-topology](http://www.searchtb.com/2012/12/%E7%8E%A9%E8%BD%ACcpu-topology.html) 给了一个表格，给出了当程序只使用一个 node 资源和使用多个 node 资源的比较（差不多是 38s 与 28s 的差距）。所以限定程序在 numa node 中运行是有实际意义的。

但是话又说回来了，指定 numa 就一定好吗？这就涉及 **numa 陷阱问题**。《[SWAP 的罪与罚](https://huoding.com/2012/11/08/198)》文章就说到了一个 numa 的陷阱问题。现象是**当你的服务器还有内存的时候，发现它已经在开始使用 swap 了，甚至已经导致机器出现停滞的现象**。这个就有可能是由于 numa 的限制，如果一个进程限制它只能使用自己的 numa node 
上的内存，那么当自身 numa node 内存使用光之后，就不会去使用其他 numa node 的内存了，会开始使用 swap ；甚至更糟的情况，机器没有设置 swap 的时候，可能会直接死机！所以，你可以使用 `numactl --interleave=all` 来取消 numa node 的限制。

综上所述得出的结论就是：**根据具体业务决定 NUMA 的使用**。

- 如果你的程序是会占用大规模内存的，你大多应该选择关闭 numa node 的限制（或从硬件关闭 numa），因为这个时候你的程序很有几率会碰到 numa 陷阱。
- 如果你的程序并不占用大内存，而是要求更快的程序运行时间，你大多应该选择限制只访问本 numa node 的方法来进行处理。

### NUMA 的取舍与优化设置

- 在 OS 层 numa 关闭时，打开 BIOS 层的 numa 会影响性能，QPS 会下降 15-30% ；
- 在 BIOS 层的 numa 关闭时，无论 OS 层面的 numa 是否打开，都不会影响性能； 


----------


## [numactl 手册](https://linux.die.net/man/8/numactl)

> 以下为 numactl 的 man 手册中的部分内容；

`numactl` 是用于针对**进程**或**共享内存**进行 NUMA 策略控制的命令行工具；

`numactl` 用于将 processes 运行在特定的 NUMA 调度策略或内存布局（placement）策略上；该策略作用于当前 command 并可以被其下所有子进程所继承；除此之外，还可以针对共享内存或文件设置持久性（persistent）策略；

可设置的策略如下：

- **`--interleave=nodes, -i nodes`**

设置内存交织策略（memory interleave policy）；即内存分配采用在 nodes 上轮询的方式；当内存无法从当前交织目标（interleave target）上成功分配时，则转移到其他 nodes 上尝试获取；可以在 `--interleave`, `--membind` 和 `--cpunodebind` 选项上指定多个 nodes ；你也可以指定 "all" 以表明可以使用当前 cpuset 上的全部 nodes ；nodes 的指定方式可以是 N,N,N 或 N-N 或 N,N-N 或 N-N,N-N 等等；Relative nodes 的指定方式可以是 +N,N,N 或 +N-N 或 +N,N-N 等等；其中 `+` 表明 node 号是相对 process 当前 cpuset 的 allowed nodes 集合；而 !N-N 表示的是和 N-N 相反的含义，即除了 N-N 的所有 nodes ；如何需要和 `+` 一起使用，则可以写成 !+N-N ；

- **`--membind=nodes, -m nodes`**

仅从指定 nodes 上分配内存；当在指定的这些 nodes 上没有足够的可用内存时，分配可能会失败；nodes 的指定方式同上面；

- **`--cpunodebind=nodes, -N nodes`**

仅在指定 nodes 上的 CPUs 中执行 command ；需要注意的是，nodes 可能由多个 CPUs 构成；nodes 的指定方式同上面；

- **`--physcpubind=cpus, -C cpus`**

仅在指定的 cpus 上执行；可指定的 cpu 编号值可以和 /proc/cpuinfo 文件中的 processor 域相同，也可以使用相对当前 cpuset 的 relative cpus 值；你可以指定 "all" ，即表明可以使用当前 cpuset 中的全部 cpus ；物理 cpus 的执行方式为 N,N,N 或 N-N 或 N,N-N 或 N-N,N-N 等等；Relative nodes 的指定方式可以是 +N,N,N 或 +N-N 或 +N,N-N 等等；其中 `+` 表明 node 号是相对 process 当前 cpuset 的 allowed nodes 集合；而 !N-N 表示的是和 N-N 相反的含义，即除了 N-N 的所有 nodes ；如何需要和 `+` 一起使用，则可以写成 !+N-N ；

- **`--localalloc, -l`**

总是在当前 node 上分配内存；

- **`--preferred=node`**

优先在指定 node 上分配内存，但如果无法分配到所需内存，则会转移到其他 nodes 上进行分配；该选项仅指定一个单独的 node 号；Relative 表示法也可以使用；

- **`--show, -s`**

展示当前进程的 NUMA 策略设置；

- **`--hardware, -H`**

显示当前系统中可用 nodes 的详细清单；

注意事项：

- 要求内核对 NUMA 策略是感知的（aware）；
- Command 不要基于 `shell` 来执行；如果你需要在子进程中使用 shell 的元字符，请使用 `sh -c` 对命令进行包装；
- 过时的 `--cpubind` 选项是接受 node 号作为参数的，而不是 cpu 号；因此已经被新的 `--cpunodebind` 和 `--physcpubind` 选项而取代；

> 问题：
> - `sh -c` 的使用；
> - 如何确定 kernel 是 numa 策略感知的？


----------


## `numactl` 使用

```shell
[root@nl-cloud-k8s-4 ~]# yum install numactl -y
[root@nl-cloud-k8s-4 ~]# numactl --show
policy: default
preferred node: current
physcpubind: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
cpubind: 0 1
nodebind: 0 1
membind: 0 1
[root@nl-cloud-k8s-4 ~]# numactl --hardware
available: 2 nodes (0-1)
node 0 cpus: 0 2 4 6 8 10 12 14 16 18 20 22 -- 这里
node 0 size: 15969 MB
node 0 free: 9490 MB     -- 1
node 1 cpus: 1 3 5 7 9 11 13 15 17 19 21 23 -- 这里
node 1 size: 16125 MB
node 1 free: 10796 MB    -- 2
node distances:
node   0   1
  0:  10  20
  1:  20  10
[root@nl-cloud-k8s-4 ~]#
```

> 可以通过对比 1 和 2 的数值判定一些问题：如果差距比较大，说明存在内存使用不均衡，可能由于服务器硬件、系统设置不当，或没有关闭 NUMA 导致；

对比

```shell
[root@xg-bigkey-rediscluster-1 ~]# numactl --show
policy: default
preferred node: current
physcpubind: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23
cpubind: 0 1
nodebind: 0 1
membind: 0 1
[root@xg-bigkey-rediscluster-1 ~]#
[root@xg-bigkey-rediscluster-1 ~]# numactl --hardware
available: 2 nodes (0-1)
node 0 cpus: 0 1 2 3 4 5 12 13 14 15 16 17 -- 这里
node 0 size: 48757 MB
node 0 free: 18496 MB
node 1 cpus: 6 7 8 9 10 11 18 19 20 21 22 23 -- 这里
node 1 size: 49152 MB
node 1 free: 8629 MB
node distances:
node   0   1
  0:  10  21
  1:  21  10
[root@xg-bigkey-rediscluster-1 ~]#
```


----------


## 通过命令判断 BIOS 层是否开启 numa

- numa 为 enable 的情况

```shell
[root@nl-cloud-k8s-4 ~]# grep -i numa /var/log/dmesg
[    0.000000] mempolicy: Enabling automatic NUMA balancing. Configure with numa_balancing= or the kernel.numa_balancing sysctl
[    1.639727] pci_bus 0000:00: on NUMA node 0
[    1.642406] pci_bus 0000:40: on NUMA node 1
[    1.645181] pci_bus 0000:3f: on NUMA node 0
[    1.647692] pci_bus 0000:7f: on NUMA node 1
[root@nl-cloud-k8s-4 ~]#
```

- numa 为 disable 的情况

```
No NUMA configuration found
```

> 问题：如何在 BIOS 层关闭 numa ？


----------


## 通过命令判断 OS 层是否开启 numa

```
[root@nl-cloud-k8s-4 ~]# cat /proc/cmdline
BOOT_IMAGE=/vmlinuz-4.9.14-1.el7.centos.x86_64 root=UUID=c42731a1-ffb6-4aed-a1ec-0758d79c5bee ro crashkernel=auto rhgb quiet LANG=en_US.UTF-8
[root@nl-cloud-k8s-4 ~]# 
```

如果出现 `numa=off` ，则表明在操作系统启动的时候就把 numa 给关掉了；否则就没有没有关闭；


----------


## 通过 `lscpu` 命令查看机器的 NUMA 拓扑结构

```shell
[root@nl-cloud-k8s-4 ~]# lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                24
On-line CPU(s) list:   0-23
Thread(s) per core:    2
Core(s) per socket:    6
Socket(s):             2
NUMA node(s):          2
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 45
Model name:            Intel(R) Xeon(R) CPU E5-2420 0 @ 1.90GHz
Stepping:              7
CPU MHz:               2199.890
BogoMIPS:              3805.58
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              15360K
NUMA node0 CPU(s):     0,2,4,6,8,10,12,14,16,18,20,22  --1
NUMA node1 CPU(s):     1,3,5,7,9,11,13,15,17,19,21,23  --2
[root@nl-cloud-k8s-4 ~]#
```


----------


## 通过 `numastat` 判定是否需要对分配策略进行调整

```shell
[root@nl-cloud-k8s-4 ~]# numastat
                           node0           node1
numa_hit               144896196       140993568
numa_miss                  15248           21040
numa_foreign               15248           21040
interleave_hit             12996           13201
local_node             144891768       140988903
other_node                  4428            4665
[root@nl-cloud-k8s-4 ~]#
```

> 当发现某个 node 上的 numa_miss 数值比较高时，说明需要对分配策略进行调整。例如，可以将指定进程关联绑定到指定的 CPU 上，从而提高内存命中率。


----------


## 应用案例

### Redis 与 NUMA

在《[How fast is Redis?](https://redis.io/topics/benchmarks)》提到：

> 在支持多 CPU sockets 的 servers 上，Redis 的性能将取决于 NUMA 配置和进程（运行）位置；最直观的效果是 `redis-benchmark` 的输出结果看起来存在了不确定性（non-deterministic），因为 client 和 server 进程是随机分布在不同的核心上的；为了获取到确定的（deterministic）的结果，需要使用进程布局工具（process placement tools），例如 Linux 系统中的 `taskset` 或 `numactl` ；最高效的组合自然是能够将 client 和 server 分开跑在同一个 CPU 的不同核心上的情况，以便从 L3 cache 中获益；下面给出针对 3 种服务器 CPU（AMD Istanbul, Intel Nehalem EX, and Intel Westmere)基于 4 KB SET benchmark 得到一些测试结果；需要注意的是，该 benchmark 结果不用作 CPU models 间的比较（因此 CPUs 的准确 model 和 frequency 情况没有给出）；

![redis_benchmark_NUMA_chart](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/redis_benchmark_NUMA_chart.gif "redis_benchmark_NUMA_chart")


### MongoDB 与 NUMA

在 《[记一次MongoDB性能问题](https://huoding.com/2011/08/09/104)》中，作者遇到了如下警告信息：

```shell
WARNING: You are running on a NUMA machine. We suggest launching mongod like this to avoid performance problems: numactl –interleave=all mongod [other options]
```

问题解决过程中，遇到了以下几种问题：

- 数据导入的速度下降；
- 导入脚本（PHP）出现超时异常（卡在了 recvfrom 操作上）；
- 通过 `db.currentOp()` 查询 MongoDB 的当前操作，会发现几乎每个操作会消耗大量的时间；
- 运行 `mongostat` 的话，结果会显示很高的 locked 值；
- 发现每当出问题的时候，总有一个名叫 irqbalance 的进程 CPU 占用率居高不下，同时发现很多介绍 irqbalance 的文章中都提及了 NUMA ；

mongodb 官方给出的解决办法：

```
echo 0 > /proc/sys/vm/zone_reclaim_mode
numactl --interleave=all mongod [options]
```

至于 NUMA 的含义，简单点说，在有多个物理 CPU 的架构下，NUMA 把内存分为本地和远程，每个物理 CPU 都有属于自己的本地内存，访问本地内存速度快于访问远程内存，缺省情况下，每个物理 CPU 只能访问属于自己的本地内存。对于 MongoDB 这种需要大内存的服务来说就可能造成内存不足；


----------


在 《[MongoDB and NUMA Hardware](https://docs.mongodb.com/manual/administration/production-notes/#production-numa)》中有如下说明：

在支持 NUMA 的系统上运行 MongoDB 可能会引起许多运行（operational）问题，包括：周期性的性能变慢（slow performance），以及更高的系统进程使用率；

当运行 MongoDB servers 和 clients 在 NUMA 硬件上时，你应该配置内存交织策略，以便主机（host）能够按照 non-NUMA 形式运行；当部署在 Linux 机器上时，MongoDB（从 2.0 版本开始）会在启动时检测 NUMA 设置；如果当前 NUMA 配置情况可能会导致性能降级，MongoDB 会打印出警告信息；

另外，还可以参考如下文章：

- [The MySQL “swap insanity” problem and the effects of NUMA](http://jcole.us/blog/archives/2010/09/28/mysql-swap-insanity-and-the-numa-architecture/) ：在这篇文章中，作者描述了 NUMA 之于数据库的影响；该文章介绍了 NUMA 和其主要的目标，并解释了这些目标和生产环境数据库之间是如何不兼容的；尽管该博客文章着重阐明的是 NUMA 对 MySQL 对影响，但是实际上对于 MongoDB 来说，也存在类似问题；
- [NUMA: An Overview](https://queue.acm.org/detail.cfm?id=2513149).

那么，**如何在 Linux 上配置 NUMA 呢？**

当在 Linux 上运行 MongoDB 时，你应该通过如下命令中的一个，将  `sysctl` 设置中的 `zone reclaim` 去使能：

```
echo 0 | sudo tee /proc/sys/vm/zone_reclaim_mode
```
```
sudo sysctl -w vm.zone_reclaim_mode=0
```

之后，你应该使用 `numactl` 来启动你的 `mongod` 实例，当然也包括配置服务器，`mongos` 实例，以及任何 clients ；如果你的系统中找不到 `numactl` 命令，请参考系统相应文档进行安装；

下面给出了如何使用 `numactl` 启动 MongoDB 实例的命令演示：

```
numactl --interleave=all <path> <options>
```

`<path>` 用于指定待启动的目标程序；`<options>` 用于指定传递给目标程序的选项参数；

若想完全去使能 NUMA 特性，你必须执行上述两步操作；更多信息详见 [Documentation for /proc/sys/vm/*](https://www.kernel.org/doc/Documentation/sysctl/vm.txt)；

> 问题：`zone_reclaim_mode` 和 numa 的关系？


----------


### MySQL 与 NUMA

在《[MySQL单机多实例方案](http://www.hellodb.net/tag/numa)》中，作者提到：

NUMA 的内存分配策略有四种：

- **缺省(default)**：总是在本地节点分配（分配在当前进程运行的节点上）；
- **绑定(bind)**：强制分配到指定节点上；
- **交叉(interleave)**：在所有节点或者指定的节点上交织分配；
- **优先(preferred)**：在指定节点上分配，失败则在其他节点上分配。

因为 NUMA 默认的内存分配策略是优先在进程所在 CPU 的本地内存中分配，会导致 CPU 节点之间内存分配不均衡，当某个 CPU 节点的内存不足时，会导致 swap 产生，而不是从远程节点分配内存。这就是所谓的 `swap insanity` 现象。

> 有经验的系统管理员或 DBA 都知道：SWAP 导致的数据库性能 下降有多么坑爹，所以最简单的方法还是关闭掉 numa （默认）不允许跨 node 分配内存的特性吧；

MySQL 采用了线程模式，对于 NUMA 特性的支持并不好；

如果单机只运行一个 MySQL 实例，我们可以选择关闭 NUMA ，关闭的方法有三种：

1. 硬件层：在 BIOS 中设置关闭（最好的方式）；
2. OS 内核：启动时设置 numa=off （可以直接在 `/etc/grub.conf` 的 kernel 行最后添加）；
3. 可以用 `numactl` 命令将内存分配策略修改为 interleave（这是使用交织分配模式启动一个程序，也就是说程序可以随意跨节点用其他节点的内存，传说中这是效率最高的关闭 NUMA 特性的方法）；

如果单机运行多个 MySQL 实例，我们可以将 MySQL 绑定在不同的 CPU 节点上，并且采用绑定的内存分配策略，强制在本节点内分配内存，这样既可以充分利用硬件的 NUMA 特性，又避免了单实例 MySQL 对多核 CPU 利用率不高的问题。


----------


在《[The MySQL “swap insanity” problem and the effects of the NUMA architecture](https://blog.jcole.us/2010/09/28/mysql-swap-insanity-and-the-numa-architecture/)》和《[A brief update on NUMA and MySQL](https://blog.jcole.us/2012/04/16/a-brief-update-on-numa-and-mysql/)》中，作者提到

> TODO


----------


## 其他

在《[Peculiar Linux kernel performance problem on NUMA systems](http://docs.datastax.com/en/landing_page/doc/landing_page/troubleshooting/cassandra/zoneReclaimMode.html)》中提到： 

- 在 NUMA 系统中遇到的奇怪的 Linux 内核性能问题可能与 `zone_reclaim_mode` 有关；
- Linux 内核在 enabling/disabling 参数 `zone_reclaim_mode` 后的表现是不一致的；这可能导致奇怪的性能问题；
- 随机的大量 CPU spikes 可能导致 latency 和 throughput 的极大增加；
- 程序在什么都没有做的情况下莫名其妙的 hang 住；
- 有些症状（Symptoms）会突然出现和突然消失；
- 在重启过后，有些症状（symptoms）通常一段时间内不会出现；

为了确保 `zone_reclaim_mode` 被去使能，可以执行：

```shell
$ echo 0 > /proc/sys/vm/zone_reclaim_mode
```


----------


## 相关内核参数

### **`vm.zone_reclaim_mode`**

在内核[文档](https://www.kernel.org/doc/Documentation/sysctl/vm.txt)中有如下说明：

通过设置 `zone_reclaim_mode` 能够实现在当某个 zone 内存用光时，采用相对多少有些激进（aggressive）的内存回收策略；如果设置为 0 ，那么将不会有 zone reclaim 发生；内存分配的需求将会从系统中其他 zones / nodes 上进行分配以满足需要；

以下值可以按位进行 or 操作：

- 1	= Zone reclaim on
- 2	= Zone reclaim writes dirty pages out
- 4	= Zone reclaim swaps pages

`zone_reclaim_mode` 默认是 disabled 的；对于文件服务器或者能够从数据缓存（cache）中获益的工作负载类型来说，`zone_reclaim_mode` 最好保持 disabled 状态，因为 caching 带来的收益很可能会比数据本地性（data locality）的收益更重要；

zone reclaim 可以被设置为 enabled ，如果明确知道工作负载具有分区特性（partitioned），即每一个分区对应到了一个 `NUMA node` 上，并且访问 remote
内存能够导致可测量到性能损耗；（在这种设置下）页分配器（page allocator）将更容易回收那些可重用的 pages（那些当前尚未被使用的、可用作 page cache 的 pages），在分配光 node 上的所有 pages 之前；

允许 zone reclaim 进行 pages 写出，可以阻止需要写大量数据的进程弄脏其他 nodes 上的 pages ；在 zone 被填满的情况下，zone
reclaim 会将 dirty pages 写出，并且有效的将进程进行节流（throttle）；这可能会降低单进程的性能，因为其无法使用系统内存的全部以便对所有写出进行缓存（buffer），但是却能够实现对其他 nodes 上对内存进行保留的效果，因此运行在其他 nodes 上的其他进程的性能将不会受到影响；

允许常规的 `swap` 可以非常有效的将内存分配限制在本地 node 上，除非显式的通过内存策略或者 cpuset 配置进行规则覆盖；


### **`vm.swappiness`**

`vm.swappiness` 是操作系统控制**物理内存**交换出去的策略。它允许的值是一个百分比的值，最小为 0 ，最大运行 100 ，该值默认为 60 。vm.swappiness 设置为 0 表示尽量少（对 inactive 内存页）进行 swap ，100 表示尽量将 inactive 的内存页交换出去。

具体的说：**当内存基本用满的时候，系统会根据这个参数来判断是把内存中很少用到的 inactive 内存交换出去，还是释放数据的 cache** 。

- `cache` 内存中缓存着从磁盘读出来的数据，根据**程序的局部性原理**，这些数据有可能在接下来又要被读取；
- `inactive` 内存，顾名思义，就是那些被应用程序映射着，但是长时间不用的内存；

可以利用 `vmstat` 看到 inactive 的内存的数量：

```
root@vagrant-ubuntu-trusty:~# vmstat -an 1
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free  inact active   si   so    bi    bo   in   cs us sy id wa st
 0  0      0 174268 168056 117464    0    0    67     3  145  426  0  1 99  0  0
 0  0      0 174252 168072 117476    0    0     0     0  145  282  0  0 100  0  0
 0  0      0 174252 168072 117476    0    0     0     0  144  290  0  0 100  0  0
 0  0      0 174252 168072 117476    0    0     0     0  139  268  0  0 100  0  0
 0  0      0 174252 168072 117476    0    0     0     0  142  271  0  1 99  0  0
 0  0      0 174252 168072 117488    0    0     0     0  149  280  0  1 99  0  0
 0  0      0 174252 168072 117488    0    0     0    24  172  307  0  0 100  0  0
 0  0      0 174252 168072 117488    0    0     0     0  144  281  0  1 99  0  0
 0  0      0 174252 168072 117488    0    0     0     0  142  273  0  0 100  0  0
^C
root@vagrant-ubuntu-trusty:~#
```

通过 `/proc/meminfo` 可以看到更详细的信息：

```
root@vagrant-ubuntu-trusty:~# cat /proc/meminfo | grep -i inact
Inactive:         168072 kB
Inactive(anon):     4592 kB
Inactive(file):   163480 kB
root@vagrant-ubuntu-trusty:~#
```

>> 针对 inactive 内存进一步深入讨论；
>
> 在 Linux 中，内存可能处于三种状态：
>
> - free
> - active
> - inactive
>
> 众所周知，Linux 内核在内部维护了很多 LRU 列表用来管理内存，比如 LRU_INACTIVE_ANON, LRU_ACTIVE_ANON, LRU_INACTIVE_FILE, LRU_ACTIVE_FILE, LRU_UNEVICTABLE。
> 其中
>
> - **LRU_INACTIVE_ANON**, **LRU_ACTIVE_ANON** 用来管理匿名页；
> - **LRU_INACTIVE_FILE**, **LRU_ACTIVE_FILE** 用来管理 page caches 页缓存；
> 
> 系统内核会根据内存页的访问情况，不定时的将活跃 active 内存被移到 inactive 列表中，这些 inactive 的内存可以被交换到 swap 中去。
> 


----------


## 参考

- wiki: [Non-uniform memory access](https://en.wikipedia.org/wiki/Non-uniform_memory_access)
- [FAQ about the NUMA architecture](http://lse.sourceforge.net/numa/faq/)
- [NUMA (Non-Uniform Memory Access): An Overview](https://queue.acm.org/detail.cfm?id=2513149)
- [NUMA 的取舍与优化设置](http://www.cnblogs.com/wjoyxt/p/4804081.html)
- [Linux 的 NUMA 技术](http://www.ibm.com/developerworks/cn/linux/l-numa/index.html)
- [NUMA Best Practices for Dell PowerEdge 12th Generation Servers](http://en.community.dell.com/techcenter/extras/m/white_papers/20266946)
- [Peculiar Linux kernel performance problem on NUMA systems](http://docs.datastax.com/en/landing_page/doc/landing_page/troubleshooting/cassandra/zoneReclaimMode.html)
- [numactl 手册](https://linux.die.net/man/8/numactl)
- [The MySQL “swap insanity” problem and the effects of the NUMA architecture](https://blog.jcole.us/2010/09/28/mysql-swap-insanity-and-the-numa-architecture/)
- [A brief update on NUMA and MySQL](https://blog.jcole.us/2012/04/16/a-brief-update-on-numa-and-mysql/)
- [MySQL单机多实例方案](http://www.hellodb.net/tag/numa)
- [MongoDB and NUMA Hardware](https://docs.mongodb.com/manual/administration/production-notes/#production-numa)
- [记一次MongoDB性能问题](https://huoding.com/2011/08/09/104)
- [How fast is Redis?](https://redis.io/topics/benchmarks)
- [详解服务器内存带宽计算和使用情况测量](http://blog.yufeng.info/archives/1511)
- [玩转 cpu-topology](http://www.searchtb.com/2012/12/%E7%8E%A9%E8%BD%ACcpu-topology.html)
- kernel: [sysctl/vm.txt](https://www.kernel.org/doc/Documentation/sysctl/vm.txt)
- [SWAP 的罪与罚](https://huoding.com/2012/11/08/198)






