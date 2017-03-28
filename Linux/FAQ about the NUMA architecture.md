# [FAQ about the NUMA architecture](http://lse.sourceforge.net/numa/faq/)

> 本文解答了一些 NUMA 架构的 FAQ ；

## What does NUMA stand for?

NUMA 即 Non-Uniform Memory Access.

## OK, So what does Non-Uniform Memory Access really mean to me?

Non-Uniform Memory Access 表示在访问某些内存区域时会比访问另外一些区域耗时更长；这是因为不同内存区域位于物理上的不同总线上；要想获得更佳具象的描述，请参考下面关于 NUMA 架构实现的描述；另外，还可以参照针对 NUMA 架构的现实世界类比；对于不感知 NUMA 的程序来说，可能会造成运行效果变差；NUMA 还引入了 local 和 remote 内存概念；

## What is the difference between NUMA and SMP?

NUMA 架构被设计出来用于解决 SMP 架构扩展性受限问题；对于 SMP 来说，即 Symmetric Multi-Processing，所有的内存访问请求都被发送给相同的共享内存总线上；这种方式对于相对少量的 CPUs 来说没有什么问题，但是当你使用数十，甚至上百 CPUs 时，共享总线的问题就出现了，因为所有 CPUs 都会竞争共享内存总线的访问权；NUMA 通过限制可以出现在任何内存总线上的 CPUs 数量，以及通过高速互联（路径）互通的方式缓解了这个瓶颈问题；

## What is the difference between NUMA and ccNUMA?

在当前这个时间点上两者几乎没有什么差别；ccNUMA 表示 Cache-Coherent NUMA，但是 NUMA 和 ccNUMA 几乎可以当作是同义的，因为基于 non-cache coherent NUMA 机器的应用几乎不存在；因此，除非特别说明，我们在说 NUMA 的时候实际上就是指 ccNUMA ；

## What is a node?

描述 NUMA 时常常遇到的问题之一就是实现该技术的方式存在许多种；这就导致了针对什么是 node 存在大量的“定义”方式；一种从技术角度来说非常正确，但同时非常丑陋的定义方式为：**a region of memory in which every byte has the same distance from each CPU** ；而更通俗的定义为：**a block of memory and the CPUs, I/O, etc. physically on the same bus as the memory** ；在某些架构中，确实存在 memory, CPUs 和 I/O 并非全部位于相同物理总线上的情况，因此，第二种定义并不总是正确的；在许多情况下，越少的技术定义可能越显得充分，但越多的技术定义会显得越准确；

## What is meant by local and remote memory?

术语 local memory 和 remote memory 通常情况下是针对当前运行进程来说的；也就是说，local memory 通常被定义为：与当前运行进程的 CPU 同属一个 node 的 memory ；而不输于该 node 的 memory 则被定义为 remote ；

Local 和 remote memory 的定义也可用于针对其他东东，而非当前的运行进程；当谈及中断上下文时，技术上讲不存在当前执行进程的概念，但我们仍将处理该中断的 CPU 所属 node 上的 memory 称作 local memory ；同样的，你还可以从 disk 的角度使用 local 和 remote memory 的概念；例如，如果存在一个 disk 附属于（attatched to） node 1 ，正在执行 DMA 操作，那么其正在读取或写入的 memory 将被称作 remote ，如果该 memory 位于另外一个 node 上（即 node 0 上）；

## What do you mean by distance?

基于 NUMA 的架构中非常有必要引入一个系统组件（即 CPUs, memory, I/O 总线等）之间的 distance 概念；用于衡量 distance 的 metric 经常有所不同，但是 hops 是比较受欢迎的一种 metric ，另外还有 latency 和 bandwidth ；这些术语的含义和其出现在网络上下文中时是一样的（几乎可以认为 NUMA 机器与紧耦合集群没有什么太大差别）；因此，当用于描述 node 时，我们可能会将特定范围内的 memory 描述成距离 CPUs 0..3 和 SCSI Controller 0 的 distance 为 2 hops (busses)；换句话说，CPUs 0..3 和 SCSI Controller 均为相同 node 的一部分；

## Could you give a real-world analogy of the NUMA architecture to help understand all these terms?

想象一下你正在烘焙蛋糕：你手头有一组配料（=**memory pages**）用于满足配方要求（=**process**）；其中一些配料你可能已经保存在了小柜子中（=**local memory**），而另外一些配料你可能根本就没有，因此你不得不向邻居索取（=**remote memory**）；一般性想法就是保证尽可能多的配料保存在你自己的小柜子中，因为这样才能缩短你做蛋糕的耗时和劳动量；

还需要谨记：你的小柜子只能保存一定数量的配料（=**physical nodal memory**）；如果你尝试购买更多的东西，但却没有空间来保存，你将不得不拜托邻居帮你保存配料在其小柜子中，直到你需要时获取（=**local memory full, so allocate pages remotely**）；

A bit of a strange example, I'll admit, but I think it works. If you have a better analogy, I'm all ears! ;)

## Why should I use NUMA? What are the benefits of NUMA?

NUMA 带来的最大收益，正如上面所说，是可扩展性；将 SMP 扩展至可支持 8-12 CPUs 是极其困难的；在这个 CPU 数目下，memory 总线将面临严重的竞争问题；而 NUMA 正是用于减少竞争访问共享 memory 总线的 CPUs 数量的一种方式；具体的解决办法为：提供多条 memory 总线，然后搭配相对少量的 CPUs 到每一条总线上；事实上，存在很多其他方式可以构建大规模多处理器机器，但此处为 NUMA FAQ ，因此还是将其他方法如何实现的讨论留到其他 FAQs 上吧；

## What are the peculiarities of NUMA?

CPU 和/或 node caches 会导致 NUMA effects 问题；例如，位于特定 node 上的 CPUs 在访问同一 node 上的 memory 和 CPUs 时，将得到更高的带宽和/或更低的延迟；由于这个原因，您可能会在高竞争压力下看到诸如 **lock starvation** 现象；这是因为如果目标 node 上的 CPU x 请求了已被同一 node 上 CPU y 获取的 lock 时，该请求将会（倾向于）“击退”来自 remote CPU z 的请求；

## What are some alternatives to NUMA?

将 memory 进行拆分后（尽可能任意）分配给一组 CPUs 能够获得一些和真正 NUMA 类似的性能收益；类似这种构建方式已经很像常规的 NUMA 机器了，其中 local 和 remote memory 之间的界限是模糊的；因为全部 memory 实际上仍位于相同的总线之上；PowerPC 的 Regatta 系统就是这种实现的例子；

你也可以基于 clusters 达到一定程度类似 NUMA 的性能；一个 cluster 非常相似于一个 NUMA 机器，其中位于 cluster 中的每一台机器将成为我们虚拟 NUMA 机器中的 node ；唯一的真正差别在于节点间延迟；在集群环境中，节点相互之间的 latency 和 bandwidth 很可能更加糟糕；

## Could you give a brief description of the main NUMA architecture implementations?

当然！主要类型为 IBM NUMA-Q, Compaq Wildfire, 以及 SGI MIPS64 ；你可以点击[这里](NUMA%20之%20System%20Descriptions.md)查阅关于上述系统类型的详细描述和图标，以及用作比较目的的标准 SMP 系统模型；
