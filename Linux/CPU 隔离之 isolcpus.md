# CPU 隔离之 isolcpus

标签（空格分隔）： linux

---

## 内核配置中的 isolcpus 说明

> 原文地址：[这里](http://www.linuxtopia.org/online_books/linux_kernel/kernel_configuration/re46.html)

isolcpus — **用于将指定 CPUs 从 kernel scheduler 中隔离出来**；

设置方式为

```shelll
isolcpus= cpu_number [, cpu_number ,...]
```

通过指定 cpu_number 值，将指定 CPUs 从常规的 kernel SMP balancing 和 scheduler 算法中移除；将某个进程搬移到 "isolated" CPU 或从 "isolated" CPU 上搬移走是通过 CPU affinity 相关的系统调用完成；cpu_number 的起始值为 0 ，因此最大值为系统中的 CPUs 数目减去 1 ；

**该选项是进行 CPUs 隔离的推荐方式**；另外一种替代方案是，针对系统中的所有 task 手动设置 CPU mask ，但是容易导致一些问题，以及负载均衡器性能的次优（suboptimal）问题的产生；


----------


## [isolcpus, numactl and taskset](https://codywu2010.wordpress.com/2015/09/27/isolcpus-numactl-and-taskset/)

> 本文主要说明：
> 
> - isolcpus 的设置；
> - isolcpus 设置后对 numactl 和 taskset 的影响；

`isolcpus` is one of the kernel boot params that isolated certain cpus from kernel scheduling, which is especially useful if you want to dedicate some cpus for special tasks with least unwanted interruption (but cannot get to 0) in a multi-core system.

With the options set, by default all user processes are created with cpu affinity mask excluding the isolated cpus.

To check whether the kernel is booted with `isolcpus` set simply check `/proc/cmdline`, for example, for my ubuntu system, I have 4 cpus:

```
$lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                4
On-line CPU(s) list:   0-3
Thread(s) per core:    1
Core(s) per socket:    1
Socket(s):             4
NUMA node(s):          1
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 30
Stepping:              5
CPU MHz:               2925.979
BogoMIPS:              5851.95
Virtualization:        VT-x
Hypervisor vendor:     VMware
Virtualization type:   full
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              8192K
NUMA node0 CPU(s):     0-3
$uname -a
Linux ubuntu 3.16.0-44-generic #59~14.04.1-Ubuntu SMP Tue Jul 7 15:07:27 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
$cat /proc/cmdline
BOOT_IMAGE=/boot/vmlinuz-3.16.0-44-generic root=UUID=.. ro find_preseed=/preseed.cfg auto noprompt priority=critical locale=en_US isolcpus=3 quiet
```

Now all user processes should be scheduled free from cpu 3.

For example, we can check the current shell process, and our focus is on **Cpus_allowed** and **Cpus_allowed_list** which has specified only 0-2 cpus is allowed (the system only has 4 cpus):

```
$cat /proc/$$/cmdline
/bin/ksh93
$cat /proc/$$/status|tail -6
Cpus_allowed:   ffffffff,fffffff7
Cpus_allowed_list:      0-2,4-63
Mems_allowed:   00000000,00000001
Mems_allowed_list:      0
voluntary_ctxt_switches:        399
nonvoluntary_ctxt_switches:     146
```

Setting of `isolcpus` has one interesting side effect that `numactl` stopped working.
We can bind to any other cpus but binding to cpu 3 would fail immediately.

First we try to bind to cpu1 and cpu2 and we can see **Cpus_allowed_list** updated correctly:

```
$numactl --physcpubind=1 /bin/ksh -c "cat /proc/\$\$/status|grep Cpus_allowed"
Cpus_allowed:   00000000,00000002
Cpus_allowed_list:      1
$numactl --physcpubind=2 /bin/ksh -c "cat /proc/\$\$/status|grep Cpus_allowed"
Cpus_allowed:   00000000,00000004
Cpus_allowed_list:      2
```

But if we try to bind to cpu3 it will fail:

```
$numactl --physcpubind=3 /bin/ksh -c "cat /proc/\$\$/status|grep Cpus_allowed"
libnuma: Warning: cpu argument 3 is out of range
 
<3> is invalid
usage: numactl [--all | -a] [--interleave= | -i <nodes>] [--preferred= | -p <node>]
...
```

But we can still use `taskset` to bind it to cpu3:

```
$taskset -c 3 /bin/ksh -c "cat /proc/\$\$/status|grep Cpus_allowed"
Cpus_allowed:   00000000,00000008
Cpus_allowed_list:      3
```

And why?

So it turns out `numactl`‘s logic is a bit funny.
It first read back the cpu affinity mask and use that to check whether the specified cpu list is valid or not and only when that is valid will it apply the cpu affinity mask.

And `taskset`‘s logic is simpler in that it directly go and apply the affinity mask using syscall `sched_setaffinity`.

Probably from users’ point of view, `taskset`‘s way is better.


----------

## [如何通过isolcpus指定CPU只运行特定任务](http://blog.sina.com.cn/s/blog_508d2c500100h4po.html)

> 此文应该是 RedHat 的内部文档；介绍了如何设置和使用 `isolcpus` 功能；

### Introduction

In some situations it may be desirable to configure a server equipped with multiple processors so that a subset of the total available processors are reserved for use only by the applications specifically assigned to them, and have all other server processes -- including the servicing of hardware interrupts and kernel maintenance tasks such as periodic flushing of the kernel's routing cache -- handled by the remaining processors.  This can be useful for some applications which are highly timing-sensitive, or when it is desired otherwise that an application always have one or more processors available to it.
 
While **Red Hat Enterprise Linux** is not designed to serve as a real-time operating system, it is possible to effectively dedicate a subset of processors to servicing one or more applications, which may help improve  performance.
 
For applications requiring realtime operating system support, **Red Hat Enterprise Linux MRG** is recommended.  For more information on Red Hat Enterprise Linux MRG, consult your Red Hat sales representative or visit the Red Hat Enterprise Linux MRG product overview page at http://www.redhat.com/mrg/.
 
 
### Assumptions

The following assumptions are made, and are beyond the scope of this article.
 
- The system administrator is familiar with editing kernel command line parameters in `/boot/grub/grub.conf` and operating within **rescue mode** in the event of a problem, and will take appropriate steps to ensure the original configuration can be restored in the event that such is desired.
- The server is equipped with enough processors that reserving a subset of them will leave sufficient resources to handle the remaining load by other processes and kernel functions.
- The system administrator is aware of the needs and capabilities of the application, such as its ability to utilize multiple threads.
 
> **Note**: Throughout the remainder of this article, a hypothetical scenario will be used for purpose of example.  The example server has two quad-core processors for a total of 8 distinct CPU cores, and cores 5-8 will be isolated for exclusive use by timing-sensitive applications.
 
### Configuration Procedure

For **Red Hat Enterprise Linux releases 4 and 5**, the following manual procedure can be used to achieve this type of configuration.
 
- 1> **Determine which processors are to be reserved/isolated**
    It should be kept in mind that the **kernel uses zero-based numbering to identify the processors** (the first CPU has ID 0, the second has ID 1, etc.).  In our example, we are isolating cores 5-8 (CPU ID's 4-7).
 
- 2> **Add the appropriate `isolcpus` kernel boot command line parameter**
    Edit the kernel command line in the `/boot/grub/grub.conf` file for the desired kernel, listing each processor to be isolated separated by commas.
     
    For example, to isolate CPU cores 5-8 (processor ID numbers 4-7), add `isolcpus=4,5,6,7` to the kernel command line options.
     
    If the original kernel stanza in `/boot/grub/grub.conf` reads:
    
    ```
    title Red Hat Enterprise Linux AS (2.6.9-78.EL.smp)        root (hd0,0)        kernel /vmlinuz-2.6.9-78.ELsmp ro root=LABEL=/        initrd /initrd-2.6.9-78.ELsmp.img
    ```
    
    the modified stanza would read:
    
    ```
    title Red Hat Enterprise Linux AS (2.6.9-78.EL.smp)        root (hd0,0)        kernel /vmlinuz-2.6.9-78.ELsmp ro root=LABEL=/ isolcpus=4,5,6,7        initrd /initrd-2.6.9-78.ELsmp.img
    ```

- 3> **Disable the irqbalance service**
    The `irqbalance` service periodically re-distributes the servicing of **hardware interrupt request (IRQ)** signals among available processors to equally balance the load.  The simplest way to prevent the irqbalance daemon from distributing IRQ service tasks to the isolated processors is to disable the irqbalance service.
    
    ```
    chkconfig --level 12345 irqbalance off
    ```

- 4> **Make a List of Numeric Interrupts**
    Examine the output of `/proc/interrupts`, noting each numeric interrupt listed.  The interrupts listed are those which have been active since booting the system.  The **Symmetric Multi-Processing (SMP)** affinity will need to be set for each of these to prevent the isolated CPU's from being assigned to them by the kernel.
    
    Example output:
    
    ```
    [vincew@pharaoh ~]$ cat /proc/interrupts
               CPU0       CPU1       CPU2       CPU3       CPU4       CPU5       CPU6       CPU7
      0:         90          0          0          0          1          0          0          0   IO-APIC-edge      timer
      1:          3          2          1          3          1          1          1          2   IO-APIC-edge      i8042
      4:          0          0          0          1          0          1          0          1   IO-APIC-edge
      6:          0          0          1          0          1          0          0          0   IO-APIC-edge      floppy
      7:          0          0          0          0          0          0          0          0   IO-APIC-edge      parport0
      8:          0          0          1          0          0          0          0          0   IO-APIC-edge      rtc
      9:          0          0          0          0          0          0          0          0   IO-APIC-fasteoi   acpi
    12:          0          1          0          1          0          1          0          1   IO-APIC-edge      i8042
    14:          0          0          0          0          0          0          0          0   IO-APIC-edge      ata_piix
    15:          0          0          0          0          0          0          0          0   IO-APIC-edge      ata_piix
    16:     168930     168814     168690     168740     168737     168665     168801     168558   IO-APIC-fasteoi   firewire_ohci, nvidia
    17:         32         36         39         42         36         29         36         35   IO-APIC-fasteoi   HDA Intel
    18:      11423      11423      11531      11517      11551      11495      11379      11440   IO-APIC-fasteoi   eth0
    20:      87758      87742      87837      87780      87764      87794      87861      88025   IO-APIC-fasteoi   ehci_hcd:usb1, uhci_hcd:usb2
    21:          0          0          0          0          0          0          0          0   IO-APIC-fasteoi   uhci_hcd:usb3
    22:          0          0          0          0          0          0          0          0   IO-APIC-fasteoi   uhci_hcd:usb4
    23:       6528       6656       6574       6590       6583       6688       6595       6612   IO-APIC-fasteoi   uhci_hcd:usb5, ahci
    NMI:          0          0          0          0          0          0          0          0   Non-maskable interrupts
    LOC:    3832781    3410057    3135853    1628674    1521292    1545106    1187711    1247285   Local timer interrupts
    RES:      96836      80994      40640      31803      50957      35632      23327      18266   Rescheduling interrupts
    CAL:      17378      17232      17098      19941      22691      22167      22265      22218   function call interrupts
    TLB:      36208      36502      26155      23681      29545      23042      21444      17505   TLB shootdowns
    TRM:          0          0          0          0          0          0          0          0   Thermal event interrupts
    THR:          0          0          0          0          0          0          0          0   Threshold APIC interrupts
    SPU:          0          0          0          0          0          0          0          0   Spurious interrupts
    ERR:          0
    ```
    
    Note that the IRQ numbers seen will likely vary from server to server, and can also change on a given server if hardware is added or removed, interrupts previously unused become active, or if options affecting **IRQ routing** are changed in the server BIOS or on the kernel command line.
     
    Furthermore, **inactive interrupts** should also be considered. By listing the contents of the `/proc/irq` virtual directory, additional interrupt numbers which do not display in the output of `/proc/interrupts` can also be identified.  This would be useful in the case of a device which is normally inactive, but is used occasionally.
     
    Finally, note that **IRQ numbers 0 and 2 are special**, and should be omitted from the list of numeric interrupts to set an SMP affinity mask for.  The kernel will not allow these two to be changed anyway.
     
    So in our example, we will need to set an SMP affinity mask for active IRQ's 1, 4, 6-9, 12, 14, 15-18, and 20-23.  For the sake of example, we will also assume that IRQ numbers 3, 5, 7, 10, 11, 13, and 19 are also seen when listing the contents of the `/proc/irq` virtual directory, and have decided for simplicity to have processors 1-4 (ID's 0-3) handle all hardware interrupt service.
     
    Therefore all numeric interrupt numbers between 1 - 23 (omitting IRQ 0 and 2) will need to have their SMP affinity masks set.
 
- 5> **Determine the Correct IRQ SMP Affinity Mask**
    For simplicity we will work with decimal masks.  To determine the correct mask, the decimal values represented by the bit number representing each processor we wish to allow a hardware interrupt to execute on are added together.  The following table can be used:
    
    ```
    Zero-based CPU ID:      7       6       5       4       3       2       1       0 Decimal Value:        128      64      32      16       8       4       2       1
    ```
    
    For systems with a larger number of processors, the decimal value doubles with each subsequently higher processor number.  For example CPU ID 8 = 256, CPU ID 9 = 512, etc.
    
    In our example, we wish to allow all hardware interrupts be serviced by CPU's 1-4 (CPU ID's 0-3), so we add 1 + 2 + 4 + 8 to arrive at a decimal mask of 15.
 
- 6> **Set the IRQ SMP Affinity Mask**
    To set the affinity mask, the calculated mask needs to be echoed into `/proc/irq/(NUMBER)/smp_affinity`, where `NUMBER` is the IRQ number to be set.  For example:
    
    ```
    # echo "15" > /proc/irq/7/smp_affinity
    ```
    
    This should be repeated for each numeric interrupt number.
     
    Correct behavior can be verified by monitoring the output of `/proc/interrupts`.  It should be noticed with repeated viewing of `/proc/interrupts` output that IRQ service counts do not increment on processors excluded by the SMP affinity mask.
     
    **Note** that with some hardware platforms, interrupts can still be serviced by CPU's excluded by the defined affinity mask, but this should be recognized as a **hardware** or **BIOS/firmware** implementation problem.
 
- 7> **Make IRQ SMP Affinity Mask Settings Persistent**
    To make the customized IRQ SMP affinity settings persistent across server reboots, the settings will need to be placed in a file which executes at each server boot, such as `/etc/rc.d/rc.local`.  While these settings are kernel tunable parameters, they exist outside of `/proc/sys` and therefore cannot be added to `/etc/sysctl.conf`.
     
    It is normal to see the isolated processors initially handle some hardware interrupts after a server boot, since the kernel may direct some interrupt work to the isolated processors until the script file containing the custom IRQ SMP affinity settings is executed.
 
- 8> **Assign Desired Application(s) to the Isolated Processors**
    Application developers can use the `sched_setaffinity` and `sched_getaffinity` system function calls directly within an application to set and retrieve the application's CPU affinity, respectively.
     
    Applications can also be set to execute only on a specific processor, or group of processors, using the `taskset` utility, which is packaged for Red Hat Enterprise Linux in the `schedutils` RPM.
     
    `taskset` can be used to **modify** the CPU affinity of running processes if the program's process id (pid) is known, and can also be used to **launch** a command or program using the CPU affinity specified on the command line.  `taskset` will accept CPU masks in the form of numeric CPU ID's, as well as in the form of hexadecimal notation bitmasks.  For simplicity the examples below will utilize numeric CPU ID's.
     
    Examples:
    
    **Application Already Running**
    
    To set an application which is already running, having pid number 12345, to execute on CPU cores 5-8 (CPU ID's 4-7):
    
    ```
    # taskset -p -c 4-7 12345
    ```
    
    or this would be equally valid syntax:
    
    ```
    # taskset -p -c 4,5,6,7 12345
    ```
    
    Start Application With Specified CPU Affinity
    
    ```
    # taskset -c 4-7 /opt/foobar1.3/my-special-program
    ```

    *NOTE: The CPU affinity of a child process created via the `fork` system function call is inherited from its parent process, so if the application launches several related processes via `fork`, they will also run on the isolated CPU's.*
 
### More Information

For more information, refer to the `taskset` and `sched_setaffinity` man pages.

----------


## [isolcpus 功能与使用介绍](http://blog.csdn.net/haitaoliang/article/details/22427045)

`isolcpus` 功能存在已久，从内核版本 `v2.6.11`（2005年）那时就已经支持了该功能。

`isolcpus` 功能主要用于在 **SMP 均衡调度算法**中将一个或多个 CPU 孤立出来。同时可通过亲缘性设置将目标进程置于该“孤立 CPU”中运行；这种方法是推荐的使用“孤立 CPU”的方式，与手动设置每个任务的亲缘性相比，后者会降低调度器的性能；

`isolcpus` 带来的好处是：**有效地提高了“孤立 cpu” 上任务运行的实时性**。该功能在保证“孤立 cpu”上任务运行的同时，减少了其他任务可以使用的 cpu 资源，所以需要在使用前对 cpu 资源进行规划；

`isolcpus` 功能使用步骤：

- 决定需要孤立多少和哪些 cpu ；
- 基于命令行参数指定“孤立 cpu” ；
- 禁止使用 `irqbalance` 中断均衡服务；
- 了解所有中断，进行中断亲缘性设计与设置；
- 决定运行在“孤立 cpu”上的任务；

> 关于 irqbalance 问题可以学习《[深度剖析告诉你irqbalance有用吗？](http://blog.yufeng.info/archives/2422)》


----------


## [【鹅厂网事】走进腾讯公网传输系统](http://chuansong.me/n/1163233)

> 部分内容有删减；

### 背景

互联网企业的基础架构不断地被其业务发展规模所挑战，由于业务应用场景复杂多变，经常会出现单个业务部署在多个 IDC 或者多个业务之间需要相互通信的情况，这样，就会引起**跨 IDC 的通信问题**；

一般而言，在广域网上运营商会提供**网络专线服务**来实现这种通信。其优点是安全性好，QoS 也可以得到保证。但其缺点也很明显，价格高昂，建设部署由于涉及到外部运营商，经常需要较长周期；超长距离的跨国专线由于涉及多个运营商，建设运营不可控因素多等。很容易想到的一个替代的方案是**利用公网资源来规避**这些问题，公网资源获取成本很低，可以实现 anywhere/anytime 的服务；部署方便，通过两端架设设备的方式能够可控的实现部署；跨国公网依赖于全球运营商现有的相对成熟稳定的因特网，可以保证绝大部分情况可用；

> 公网传输系统（代号 BOBCAT）就是在这样一种应用场景和需求下提出的解决方案。

### 行业解决方案

业界应对这种场景，有一种成熟的技术，叫做**虚拟专用网（Virtual Private Network，简称VPN）**。虚拟专用网是透过公用的网络来传送私有网络的信息。它利用已加密的**隧道协议（Tunneling Protocol）**来达到保密、发送端认证、消息准确性等私有信息的安全效果。这种技术可以使用不安全的网络来发送可靠、安全的信息。当然，也有些场景不需要使用到加密，例如 MPLS VPN ，这种场景使用安全信道进行传输，使用 VPN 的目的更多在于方便网络维护管理。没有加密的也可以称为 VPN，但利用不加密 VPN 传输的消息有被窃取的危险。

根据 VPN 所在 OSI 模型的层次，VPN 有多种协议标准：

- 基于**数据链路层**的 `PPTP`、`L2TP` ；
- 基于**网络层**的 `IPsec VPN` ；
- 基于**应用层**的 `SSL VPN` 等；

根据网络连接方式的不同，还可以分为：

- **Site to Site**
- **End to Site**

Site to Site 主要指的是网络和网络之间的 VPN 连接，而 End to Site 指的是终端到网络之间的 VPN 连接。SSL VPN 就是 End to Site 类型的 VPN ；

> 目前 BOBCAT 系统属于 Site to Site 类型的 VPN 。

目前主流设备厂商的现状：

- 都有支持 VPN 功能的设备；
- 一般而言，这些厂商的设备集成了诸如 anti-DDOS、ADS、IPS 等多种安全功能；
- 为适应不同应用场景这些设备还支持了多种 VPN 方式；
- （因此）造价较高；

而随着技术的演进，**利用多核 CPU 进行数据面转发**的软硬件也越来越成熟，使用廉价的标准服务器进行数据转发成为一个选择。

同时，互联网公司的业务具有不断演进的特点，在这个过程中会对公网传输系统提出各种定制化需求。这些定制化需求开发性较强，传统设备厂商反馈周期较长，甚至没有针对特定的小众需求进行开发的意愿，无法满足业务快速响应的需求。

> 在这种背景下，我们决定采用软件的方式来解决公网传输的问题，BOBCAT 公网传输系统也因此而生。

### 腾讯解决方案

BOBCAT 系统架构分为 4 个层次，分别为

- 硬件层
- 转发层
- 管理层
- 应用层

硬件层主要负责物理硬件的支持，包括多核 CPU 支持、千兆网卡、万兆网卡的支持。

BOBCAT 目前支持两种 CPU 架构：

- Tilera 的网格架构；
- Intel 的 x86 架构；

最初该系统选用了 Tilera 的一款 64 核众核 CPU ，该 CPU 使用独有的“矩阵”架构（iMesh），该架构具备良好的扩展性，可以很轻易地将核心数量扩展到 100 个。Tilera 同时配套提供一套 MDE（Multicore Development Environment，多核开发环境），利用这个 MDE 环境可以非常方便的进行多核应用开发，发挥多核 CPU 性能。该 CPU 每个物理核时钟频率较低，为 866MHz ，使得整个 CPU 功耗较低，有利于降低整个数据中心的 TCO 。随着 Intel 的 DPDK 的发布，x86 CPU 也具备了**处理高速网络数据包的能力**。DPDK 与 Tilera 的思路有许多相似的地方，例如 **CPU 亲缘性绑定**、**内存零拷贝**、**轮询处理机制**等，通过最大限度地释放 CPU 处理能力、降低内存拷贝开销，来实现线速转发。

> BOBCAT 目前主要使用的是 Intel 的 DPDK 技术；

**一般而言，普通应用程序很难将网卡的处理能力发挥到极致**，一般应用程序能够实现数万级别 QPS 的处理性能，已经属于非常高的处理能力了。但 BOBCAT 要求的处理性能是**万兆吞吐量**，对应着**百万级别的 PPS 处理能力**，使用传统手段很难满足其线速处理需求。

限制应用程序性能的原因主要有两个：

- 一是因为普通应用程序架设在操作系统之上，而网络数据报文从网卡进入操作系统，经由驱动至 TCP/IP 协议栈处理，需要进行至少两次拷贝，内存拷贝影响了应用程序的处理性能；
- 二是普通应用程序与系统上其他应用程序共享 CPU 资源，操作系统对每个应用程序都是公平调度的，会导致频繁的进程调度、上下文切换，并导致大量 Cache Miss ，CPU 性能被无谓的浪费在这些调度过程中，无法将 CPU 能力发挥到极致。

针对上述问题，业界产生了许多解决方法，一个行之有效的方法是通过 **CPU 的亲缘性绑定**和**内存零拷贝技术**来规避这些问题。Intel 的 DPDK 将这些方法进行了整合实现，可以通过使用 DPDK 来实现这些特性，达到百万级别 PPS 转发的目的。

> 问题：什么叫线速转发？

![DPDK 框架图](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/DPDK%20%E6%A1%86%E6%9E%B6%E5%9B%BE.jpeg "DPDK 框架图")

DPDK 开发套件框架如上图所示，DPDK 提供一个**环境抽象层（Environment Abstraction Layer, EAL）**，EAL 将用户态进程与底层网卡驱动的通信行为进行了封装，为应用程序提供了通用的报文处理接口，应用程序可以很方便地调用这些接口，对报文进行处理，屏蔽了报文收发细节。

针对上面提到的普通应用程序的性能瓶颈问题，DPDK 主要采用以下方式解决：

- 基于大页表的内存管理
　　DPDK 采用**基于大页表的内存管理机制**，解决报文从驱动到协议栈，再从协议栈到用户态空间的内存拷贝问题。大页表是一项非常成熟的技术，在 Linux 2.6 内核中已经进行了支持，主要是解决大内存时的 `TLB Miss` 问题。我们知道应用程序访问的内存空间都是操作系统提供的虚拟内存空间，在实际进行访问时，需要将虚拟内存空间进行转换。虚拟内存到物理内存的转换就用到了 TLB 。**操作系统传统页面大小为 4K** ，若应用程序使用 4G 内存，将占用 1M 条表项，而一般 CPU 的 TLB 条目非常稀缺，无法同时维护如此巨大的表项，最终会通过置换算法将过期的表项移出 TLB ，导致 `TLB Miss` 。大页表通过增加单个页面的大小，减少了所需页表条目数量，可以有效的降低 TLB Miss 。同时由于大页表在系统启动时就分配好了，在使用过程中不需要换入和换出，也可以有效的减少内存交换。
　　DPDK 在大页表的基础上实现了特有的内存管理机制，DPDK 提供的网卡驱动将收到的报文通过 `DMA` 的方式映射到大页表对应的内存中，并抽取报文对应的报文指针，在应用程序和驱动之间维护一套无锁的环形队列来管理这些报文指针。应用程序只需要调用相对应的接口，从队列中提取相应的报文进行处理，完成后，再将报文指针放到发送队列中，由驱动程序进行发送。值得注意的是，在处理报文的整个过程中，不需要对报文进行任何拷贝，所有操作都依赖于报文指针，可以从整个报文生命周期上，解决内存拷贝问题，实现内存零拷贝。
　　当然，DPDK 的这种方案也有其**局限性**。**由于所有的报文没有经过内核协议栈的处理，所以所有的协议行为都需要应用程序来处理**。对开发人员而言，就需要对各种 TCP/IP 协议的特性有较深入的了解。但有些应用层协议行为十分复杂，通过应用程序来进行协议处理会带来巨大的性能开销。因此，**DPDK 这种方案更适合协议行为简单的网络应用**。DPDK 为应对这种情况，也提供了一种利用内核协议栈处理复杂协议的方式，即 **`KNI` 虚拟网卡**。通过 KNI 技术，可以将已经到达用户态的报文，再次导入内核协议栈进行处理。但这种方式的处理性能不高，需要平衡使用。

- 基于 CPU 亲缘性的进程分配
　　DPDK 的 EAL 提供**基于 CPU 亲缘性的多核并发处理机制**。前面提到由于普通应用程序运行过程中，CPU 同时需要处理其他应用及操作系统的中断，导致出现频繁的进程调度、上下文切换和大量的 `Cache Miss` ，导致普通应用程序无法达到百万 PPS 级别甚至更高的处理能力。CPU 亲缘性绑定，可以很好的解决这个问题。
　　CPU 亲缘性一般分为**软亲缘性**和**硬亲缘性**，软亲合性可以使应用程序不在处理器之间频繁迁移，而硬亲缘性更为彻底，应用程序将只运行在指定的处理器之上。**Linux 系统天然支持软亲缘性，其进程调度策略会尽量把进程固定在同一个处理器上运行**。但有些时候，仅靠操作系统的调度是不够的。例如对于网络应用而言，其处理速度要求非常高，并且要尽力减少其处理时延。这个时候就需要 CPU 的硬亲缘性绑定了，在 Linux 2.6 版本内核中对此进行了支持。
　　**DPDK 集成了 CPU 硬亲缘性绑定技术，通过启动参数的设定，可以很方便将某个 CPU 核心分配给某个应用进程使用**。当然，此种做法无法解决进程调度带来的开销，例如同一个 CPU 核心如果分配给了不同应用进程使用，仍然会带来进程调度、上下文切换、Cache Miss 等问题。这时候只能通过**独占**来解决这个问题，Linux 系统也对这种情况进行了支持。Linux 内核启动参数支持 `isolcpus`，通过该参数，可以将某些 CPU 核心进行隔离，默认不使用这些核心进行进程调度。这时候再使用 DPDK 提供的硬亲缘性绑定技术，只要小心处理，可以实现一个应用进程独占 CPU 核心的效果，避免进程调度带来的上下文切换和 Cache Miss 。


----------


