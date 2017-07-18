# 基于 vmstat 进行系统分析

## 区分应用类型

不同类型的系统用途也不尽相同，要找到性能瓶颈，则首先需要知道系统跑的究竟是什么应用，以及应用本身有些什么特点；

区分应用类型很重要，通常可分为：

- **IO 相关**：IO 相关的应用通常用来处理大量数据，需要大量内存和存储，频繁 IO 操作读写数据，而对 CPU 的要求则较少，大部分时候 CPU 都在等待硬盘，比如，数据库服务器、文件服务器等。
- **CPU 相关**：CPU 相关的应用需要使用大量 CPU，比如高并发的 web/mail 服务器、图像/视频处理、科学计算等都可被视作 CPU 相关的应用。

## vmstat 输出参数

在 man vmstat 中有

```
FIELD DESCRIPTION FOR VM MODE
   Procs
       r: The number of runnable processes (running or waiting for run time).
       b: The number of processes in uninterruptible sleep.

   Memory
       swpd: the amount of virtual memory used.
       free: the amount of idle memory.
       buff: the amount of memory used as buffers.
       cache: the amount of memory used as cache.
       inact: the amount of inactive memory.  (-a option)
       active: the amount of active memory.  (-a option)

   Swap
       si: Amount of memory swapped in from disk (/s).
       so: Amount of memory swapped to disk (/s).

   IO
       bi: Blocks received from a block device (blocks/s).
       bo: Blocks sent to a block device (blocks/s).

   System
       in: The number of interrupts per second, including the clock.
       cs: The number of context switches per second.

   CPU
       These are percentages of total CPU time.
       us: Time spent running non-kernel code.  (user time, including nice time)
       sy: Time spent running kernel code.  (system time)
       id: Time spent idle.  Prior to Linux 2.5.41, this includes IO-wait time.
       wa: Time spent waiting for IO.  Prior to Linux 2.5.41, included in idle.
       st: Time stolen from a virtual machine.  Prior to Linux 2.6.11, unknown.
```

[其它说明](http://www.lazysystemadmin.com/2011/04/understanding-vmstat-output-explained.html)：

```
Proc: 
-------
r: How many processes are waiting for CPU time.
b: Wait Queue - Process which are waiting for I/O (disk, network, user 
    input,etc..) 


Memory: 
-----------
swpd: shows how many blocks are swapped out to disk (paged). Total Virtual  
       memory usage. 
            
Note: you can see the swap area configured in server using "cat proc/swaps"


free: Idle Memory 
buff: Memory used as buffers, like before/after I/O operations
cache: Memory used as cache by the Operating System


Swap: 
---------
si: How many blocks per second the operating system is swapping in. i.e 
    Memory swapped in from the disk (Read from swap area to Memory)
so: How many blocks per second the operating system is swaped Out. i.e 
     Memory swapped to the disk 

(Written to swap area and cleared from Memory)

In Ideal condition, We like to see si and so at 0 most of the time, and we definitely don’t like to see more than 10 blocks per second.


IO: 
------
bi: Blocks received from block device - Read (like a hard disk) 
bo: Blocks sent to a block device - Write


System: 
-------------
in: The number of interrupts per second, including the clock. 
cs: The number of context switches per second. 


CPU: 
--------
us: percentage of cpu used for running non-kernel code. (user time, including 
     nice time) 
sy: percentage of cpu used for running kernel code. (system time - network, IO 
     interrupts, etc) 
id: cpu idle time in percentage.
wa: percentage of time spent by cpu for waiting to IO.


If you used to monitor this data, you can understand how is your server doing during peak usage times. 

Note: the memory, swap, and I/O statistics are in blocks, not in bytes. In Linux, blocks are usually 1,024 bytes (1 KB).
```

## 案例分析

### CPU us 高

可能原因：应用程序再做大量用户态计算；

```
$ vmstat 1
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 4  0    140 3625096 334256 3266584  0    0     0    16 1054  470 100 0  0  0  0
 4  0    140 3625220 334264 3266576  0    0     0    12 1037  448 100 0  0  0  0
 4  0    140 3624468 334264 3266580  0    0     0   148 1160  632 100 0  0  0  0
 4  0    140 3624468 334264 3266580  0    0     0     0 1078  527 100 0  0  0  0
 4  0    140 3624712 334264 3266580  0    0     0    80 1053  501 100 0  0  0  0
```

输出特点：

- r 高 (4 core 跑满)
- us 高 (~100%)
- in > cs (2/1)
- si 和 so 为 0
- bi 和 bo 几乎为 0
- swpd/free/cache 几乎不变

总结：

- CPU 和内存资源没有问题；
- 存在少量硬盘输出；
- 用户态跑满 CPU ；
- 对系统资源有一定量需求（in）；


### CPU sy 高 + cs 正常

可能原因：单个应用程序（或调度实体）请求了（依赖）大量 kernel 提供的服务（资源）以满足用户态的计算；

```
$ vmstat 1
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 4  0    140 2915476 341288 3951700  0    0     0     0 1057  523 19 81  0  0  0
 4  0    140 2915724 341296 3951700  0    0     0     0 1048  546 19 81  0  0  0
 4  0    140 2915848 341296 3951700  0    0     0     0 1044  514 18 82  0  0  0
 4  0    140 2915848 341296 3951700  0    0     0    24 1044  564 20 80  0  0  0
 4  0    140 2915848 341296 3951700  0    0     0     0 1060  546 18 82  0  0  0
```

特点：

- r 高
- sy > us (8/2 ~100%)
- in > cs (2/1)
- si 和 so 为 0
- bi 和 bo 几乎为 0
- swpd/free/cache 几乎不变

总结：

- CPU 和内存资源没有问题；
- 对系统资源有一定量需求（in）；
- 某个进程可能一直在提供态霸占着 CPU（sy 非常高，而且 cs 较低）；


###  CPU sy 高 + cs 高

可能原因：过量的应用程序（或调度实体）请求了（依赖）大量 kernel 提供的服务（资源）以满足用户态的计算；

```
$ vmstat 1
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
14  0    140 2904316 341912 3952308  0    0     0   460 1106 9593 36 64  1  0  0
17  0    140 2903492 341912 3951780  0    0     0     0 1037 9614 35 65  1  0  0
20  0    140 2902016 341912 3952000  0    0     0     0 1046 9739 35 64  1  0  0
17  0    140 2903904 341912 3951888  0    0     0    76 1044 9879 37 63  0  0  0
16  0    140 2904580 341912 3952108  0    0     0     0 1055 9808 34 65  1  0  0
```

特点：

- r 超高
- sy > us (6/4 ~100%)
- in < cs (1/9)
- si 和 so 为 0
- bi 和 bo 几乎为 0
- swpd/free/cache 几乎不变

总结：

- CPU 资源不足（r 过高）；
- 内核忙于上下文切换（cs 比 in 要高太多）；
- 多应用程序调用了大量的系统调用（sy 高，us 低，高 cs）；


### CPU wa 高 + bo 高 + b 高

可能原因：多个应用程序（或调度实体）进行大量文件写出，导致 IO 瓶颈；

```
$ vmstat 1
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 0  4    140 1962724 335516 4852308  0    0   388 65024 1442  563  0  2 47 52  0
 0  4    140 1961816 335516 4853868  0    0   768 65536 1434  522  0  1 50 48  0
 0  4    140 1960788 335516 4855300  0    0   768 48640 1412  573  0  1 50 49  0
 0  4    140 1958528 335516 4857280  0    0  1024 65536 1415  521  0  1 41 57  0
 0  5    140 1957488 335516 4858884  0    0   768 81412 1504  609  0  2 50 49  0
```

特点：

- b 高
- wa 高 (~50% > 20%)
- in > cs (3/1)
- bo >> bi > 0
- si 和 so 为 0
- swpd/free/cache 几乎不变

总结：

- CPU 和内存资源没有问题；
- IO 遇到瓶颈 (b 高，wa 高)
- 应用进行了大量写出，少量读入（in 略高于 cs 合理）；


### CPU wa 高 + si/so 高 + bi/bo 高 + swpd 增 + buff/cache 减

可能问题：文件大量读写，遇到系统内存不足情况，导致刷脏页到 SWAP 分区

```
# vmstat 1
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 0  3 252696   2432    268   7148 3604 2368  3608  2372  288  288  0  0 21 78  1
 0  2 253484   2216    228   7104 5368 2976  5372  3036  930  519  0  0  0 100  0
 0  1 259252   2616    128   6148 19784 18712 19784 18712 3821 1853  0  1  3 95  1
 1  2 260008   2188    144   6824 11824 2584 12664  2584 1347 1174 14  0  0 86  0
 2  1 262140   2964    128   5852 24912 17304 24952 17304 4737 2341 86 10  0  0  4
```

特点：

- wa 高 (~80%)
- free 保持在一定值
- swpd 逐渐增大
- buff 逐渐减少
- so 和 si 均有数值
- bi 和 bo 均有数值

总结：

- **free 基本没什么显著变化**，**swapd 逐步增加**，说明 free 已经达到最小阈值，即  2.56MB = 256MB * 10% 左右；并触发了 swap 操作（vm.dirty_background_ratio = 10）；
- **buff 逐步减少**，说明系统知道内存不够了，kswapd 正在从 buff 那里借用部分内存（FIXME：为什么从这里借？）；
- **so 一直有数值**，说明 kswapd 正持续把脏页面（FIXME：确定是 dirty page 么？为什么不是 pdflush 在刷？）写到 swap 交换区，并且从 swpd 逐渐增加也能看出确实如此（kswapd 在进行可用内存扫描时会如下检查：如果页面被修改了，但不是被文件系统修改的，则把页面写到硬盘上的 swap 交换空间；这就是此处 swapd 持续增加的原因）；

> 问题：
> 
> - 为何不提及 si 和 cache 的数值变化；
> - in 和 cs 的变化说明什么？

### CPU wa 高 + cs 高 + bi/bo 有值 + so 高 + swpd 增 + buff/cache 增

可能原因：**RAM Bottleneck (swapping) Example**

```
[user@fedora8 ~]$ vmstat 1 5
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 3  1 244208  10312   1552  62636    4   23    98   249   44  304 28  3 68  1  0
 0  2 244920   6852   1844  67284    0  544  5248   544  236 1655  4  6  0 90  0
 1  2 256556   7468   1892  69356    0 3404  6048  3448  290 2604  5 12  0 83  0
 0  2 263832   8416   1952  71028    0 3788  2792  3788  140 2926 12 14  0 74  0
 0  3 274492   7704   1964  73064    0 4444  2812  5840  295 4201  8 22  0 69  0
[user@fedora8 ~]$
```

总结：

同时打开了很多 applications（包括 VirtualBox with a Windows guest system, among others）；几乎所有 memory 都被占用了；之后，再启动一个应用时（OpenOffice），则会导致 Linux kernel 进行 swap out，将 several memory pages 换出到硬盘上的 swap file 中，以便为 OpenOffice 获取更多可用 RAM ；将 memory pages 换出到 swap file 到行为可以通过 vmstat 的 so 列 (swap out - memory swapped to disk) 看到：


### CPU sy 高 + bo 有值 + cache 增 + free 减

实际原因：基于 kernel 获取随机数写入本地文件；

> For this, `/dev/urandom` will supply random numbers, which will be generated by the kernel. This will lead to an increased load on the CPU (**sy** – system time). At the same time, the vmstat executing in parallel will indicate that between 93% and 97% of the CPU time is being used for the execution of kernel code (for the generation of random numbers, in this case).

```
root@vagrant-ubuntu-trusty:~] $ dd if=/dev/urandom of=500MBfile bs=1M count=500
```

vmstat 输出

```
root@vagrant-ubuntu-trusty:~] $ free -m
             total       used       free     shared    buffers     cached
Mem:           489        306        183          5         44        153
-/+ buffers/cache:        108        380
Swap:         2509          0       2509
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $ vmstat 1 60
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 2  0      0 188260  45700 156876    0    0     2     0  145  294  0  0 99  0  0
 0  0      0 188248  45700 156876    0    0     0     0  152  292  0  1 99  0  0
 0  0      0 188248  45700 156876    0    0     0     0  148  286  0  0 100  0  0
 0  0      0 188248  45700 156876    0    0     0     0  151  288  0  1 99  0  0
 1  0      0 188248  45700 156876    0    0     0     0  148  289  0  1 99  0  0
 
 # dd 命令开始
 1  0      0 184112  45700 159976    0    0    60     0  227  310  0 29 71  0  0
 1  0      0 172544  45700 171240    0    0     0     0  406  281  0 100  0  0  0
 1  0      0 160832  45700 182520    0    0     0     0  412  292  0 100  0  0  0
 
 # 第一次触发大 bo 
 # free 的减少量 ~= cache 的增加量
 # sy 开始进入 100% 状态忙于生成随机数
 2  0      0 149228  45704 193820    0    0     4 30640  425  285  0 100  0  0  0
 1  0      0 138796  45704 204044    0    0     0  6240  437  310  0 100  0  0  0
 1  0      0 127228  45712 215320    0    0     0    44  421  296  0 100  0  0  0
 1  0      0 115660  45712 226584    0    0     0     0  429  289  0 100  0  0  0
 1  0      0 103968  45712 237840    0    0     0 12312  425  295  0 100  0  0  0
 1  0      0  92392  45712 249104    0    0     0 12288  417  292  1 99  0  0  0
 1  0      0  80892  45712 260360    0    0     0 37888  418  312  0 100  0  0  0
 1  0      0  69324  45720 271624    0    0     0    20  420  297  0 100  0  0  0
 1  0      0  57632  45720 282880    0    0     0     0  424  284  0 100  0  0  0
 1  0      0  46076  45720 294184    0    0     0 12288  419  303  0 100  0  0  0
 1  0      0  34512  45720 305416    0    0     0     0  430  308  0 100  0  0  0
 1  0      0  22832  45724 316680    0    0     0 12288  423  318  0 100  0  0  0
 1  0      0  12392  45732 326920    0    0     0 37924  429  317  0 100  0  0  0
 1  0      0   6664  45732 332400    0    0     0     0  430  300  0 100  0  0  0
 1  0      0   6264  45732 332540    0    0     0     0  417  324  0 100  0  0  0
 1  0      0   6868  45724 331704    0    0     0 12300  424  321  0 100  0  0  0
 1  0      0   6868  45728 331152    0    0     4 12288  423  339  0 100  0  0  0
 1  0      0   6204  45736 331784    0    0     0 37916  471  360  0 100  0  0  0
 1  0      0   6328  45736 331728    0    0     0     0  422  313  0 100  0  0  0
 1  0      0   6696  45736 331216    0    0     0     0  436  323  1 99  0  0  0
 1  0      0   6408  45736 331568    0    0     0     0  414  309  0 100  0  0  0
 1  0      0   6400  45736 331648    0    0     0 12292  439  344  0 100  0  0  0
 1  0      0   6224  45744 332140    0    0     0 12312  425  338  0 100  0  0  0
 1  0      0   6356  45748 332116    0    0     4 37888  452  348  0 100  0  0  0
 1  0      0   5836  45748 332556    0    0     0     0  412  315  0 100  0  0  0
 
 # 第一次触发 so
 # free 的数值在 5836~6808 之间（和系统参数设置的阈值有关）
 1  0     24   6808  45748 330936    0   24     0    24  437  330  0 100  0  0  0
 1  0     24   6368  45748 331936    0    0     0 12288  435  333  0 100  0  0  0
 1  0     24   6504  45756 331684    0    0     0 12316  450  346  0 100  0  0  0
 1  0     24   6420  45756 331120    0    0     0 37888  438  325  0 100  0  0  0
 1  0     24   6452  45756 331132    0    0     0     0  438  326  0 100  0  0  0
 1  0     24   6820  45756 330840    0    0     0     0  424  315  0 100  0  0  0
 1  0     24   6932  45756 330504    0    0     0     0  427  326  0 100  0  0  0
 
  # 第二次触发 so
 1  0     36   7228  45764 330084    0   12     0 12348  438  345  1 99  0  0  0
 1  0     36   6844  45768 330996    0    0     4 12300  430  325  0 100  0  0  0
 1  0     40   5724  45768 332036    0    4     0 37892  457  343  0 100  0  0  0
 1  0     48   6096  45768 331668    0    8     0     8  419  306  0 100  0  0  0
 1  0     48   6308  45768 331476    0    0     0     0  417  329  0 100  0  0  0
 1  0     48   5736  45776 332024    0    0     0 12324  444  339  0 100  0  0  0
 0  1     48   6208  45788 331176    0    0    12 10240  429  336  0 99  0  1  0
 1  0     48   6460  45792 331456    0    0     4  2060  432  321  1 97  0  2  0
 1  0     48   5680  45796 331988    0    0     4 37888  451  338  0 98  0  2  0
 1  0     48   6560  45796 331240    0    0     0     0  416  337  0 100  0  0  0
 1  0     52   6916  44480 331568    0    4     0     4  426  320  0 100  0  0  0
 1  0     56   7948  43324 332900    0    4   252 12340  337  459  0 52 43  5  0
 0  0     56   7948  43324 332900    0    0     0     0  149  300  0  0 100  0  0
 1  0     56   7948  43324 332900    0    0     0     0  148  286  0  0 100  0  0
 0  0     56   7948  43324 332900    0    0     0     0  152  307  0  1 99  0  0
 0  0     56   7948  43324 332900    0    0     0     0  143  285  0  0 100  0  0
 0  0     56   7884  43352 332940    0    0    20 25636  181  332  0  0 89 11  0
 0  0     56   7884  43352 332940    0    0     0     0  149  288  0  1 99  0  0
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  0     56   7884  43352 332940    0    0     0     0  148  291  0  0 100  0  0
 0  0     56   7884  43352 332940    0    0     0     0  150  303  0  0 100  0  0
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $ free -m
             total       used       free     shared    buffers     cached
Mem:           489        481          7          5         42        325
-/+ buffers/cache:        113        375
Swap:         2509          0       2509
root@vagrant-ubuntu-trusty:~] $
```

可以看到，通过 dd 基于 kernel 获取随机数写入本地文件时（500M 文件）：

- sy 跑满
- free 逐步下降
- cache 逐步上升
- 在 free 数值为 5700 左右时会触发 swap out 
- free 减少 176M ，buffers 减少 2M ，cached 增大 172M（说明写文件会导致 cache 的大量使用）

### CPU wa 高 + bo 有值 + buff 增

实际原因：read from `/dev/zero` and write a file, **High IO Write Load Example**

> In contrast with the previous example, `dd` will read from `/dev/zero` and write a file. The `oflag=dsync` will cause the data to be written immediately to the disk (and not merely stored in the **page cache**).

测试

```
[user@fedora9 ~]$ dd if=/dev/zero of=500MBfile bs=1M count=500 oflag=dsync

[user@fedora9 ~]$ vmstat 1 5
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 1  1      0  35628  14700 1239164    0    0  1740   652  117  601 11  4 66 20  0
 0  1      0  34852  14896 1239788    0    0     0 23096  300  573  3 16  0 81  0
 0  1      0  32780  15080 1241304    0    0     4 21000  344  526  1 13  0 86  0
 0  1      0  36512  15244 1237256    0    0     0 19952  276  394  1 12  0 87  0
 0  1      0  35688  15412 1237180    0    0     0 18904  285  465  1 13  0 86  0
[user@fedora9 ~]$ 
```


### CPU sy 高

实际原因：基于 kernel 获取随机数写入 `/dev/null` 文件

```
root@vagrant-ubuntu-trusty:~] $ dd if=/dev/urandom of=/dev/null bs=1M count=500
```

则只有 sy 跑满，其它指标基本无变化；

### CPU sy/wa 高 + buff 减 + cache 增

实际原因：读取本地文件写入 `/dev/null` 文件

```
root@vagrant-ubuntu-trusty:~] $ dd if=500MBfile of=/dev/null bs=1M count=500
```

vmstat 输出：

```
root@vagrant-ubuntu-trusty:~] $ free -m
             total       used       free     shared    buffers     cached
Mem:           489        479          9          5         11        354
-/+ buffers/cache:        113        375
Swap:         2509          0       2509
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $ vmstat 1 60
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 2  0     64   9436  11900 363276    0    0     2    10  146  294  0  1 99  0  0
 0  0     64   9436  11900 363276    0    0     0     0  147  295  0  0 100  0  0
 0  0     64   9448  11900 363276    0    0     0     0  151  297  0  0 100  0  0
 0  0     64   9448  11900 363276    0    0     0     0  149  293  0  1 99  0  0
 0  0     64   9448  11900 363276    0    0     0     0  148  295  1  0 99  0  0
 
 # 第一次触发 bi 
 # 读取 250M 数据
 # sy 在 20% 左右
 # free 减少 3264，buff 减少 3604，cache 增大 10964
 # in 和 cs 同比增大
 0  1     64   6184   8296 374240    0    0 250628     0 2043 3242  0 18 52 30  0
 
 # 第二次触发 bi 
 # 读取 250M 数据
 # sy 在 20% 左右
 # free 增大 1484，buff 减少 12，cache 增大 736
 # in 和 cs 同比增大
 0  0     64   7668   8284 374976    0    0 261644     0 2102 3401  0 21 45 34  0
 0  0     64   7668   8284 374976    0    0     0     0  153  304  0  1 99  0  0
 2  0     64   7668   8284 374976    0    0     0     0  149  283  1  0 99  0  0
 1  0     64   7668   8284 374976    0    0     0     0  150  294  0  0 100  0  0
 0  0     64   7668   8292 374968    0    0     0    12  159  323  0  0 100  0  0
 0  0     64   7668   8292 374976    0    0     0     0  147  288  0  1 99  0  0
 0  0     64   7668   8292 374976    0    0     0     0  149  302  0  1 99  0  0
 0  0     64   7668   8292 374976    0    0     0     0  145  290  0  0 100  0  0
(Ctrl+C)
root@vagrant-ubuntu-trusty:~] $ free -m
             total       used       free     shared    buffers     cached
Mem:           489        481          7          5          8        366
-/+ buffers/cache:        107        381
Swap:         2509          0       2509
root@vagrant-ubuntu-trusty:~] $
```

可以看到，通过 dd 读取本地文件写入 /dev/null 文件（500M 文件）：

- sy 保持在 20%（当发生 bi 时）
- free 有所下降
- cache 有所上升
- buff 有所下降
- free 减少 2M ，buffers 减少 3M ，cached 增大 12M（说明读文件不会大量使用 cache）

### CPU wa 高 + bi 有值

实际原因：A large file (such as an ISO file) will be read and written to `/dev/null` using dd. **High IO Read Load Example**

```
[user@fedora9 ~]$ dd if=bigfile.iso of=/dev/null bs=1M
[user@fedora9 ~]$ vmstat 1 5
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 3  1 465872  36132  82588 1018364    7   17    70   127  214  838 12  3 82  3  0
 0  1 465872  33796  82620 1021820    0    0 34592     0  357  781  6 10  0 84  0
 0  1 465872  36100  82656 1019660    0    0 34340     0  358  723  5  9  0 86  0
 0  1 465872  35744  82688 1020416    0    0 33312     0  345  892  8 11  0 81  0
 0  1 465872  35716  82572 1020948    0    0 34592     0  358  738  7  8  0 85  0
[user@fedora9 ~]$ 
```

### CPU wa 高 + cs 高 + bi/bo 有值

实际原因：**CPU Waiting for IO Example**

> In the following example, an `updatedb` process is already running. The `updatedb` utility is part of `mlocate`. It examines the entire file system and accordingly creates the database for the `locate` command (by means of which file searches can be performed very quickly). Because `updatedb` **reads all of the file names from the entire file system**, the CPU must wait to get data from the IO system (the hard disk). For that reason, `vmstat` running in parallel will display large values for wa (waiting for IO):

```
[user@fedora9 ~]$ vmstat 1 5
procs -----------memory---------- ---swap-- -----io---- --system-- -----cpu------
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa st
 2  1 403256 602848  17836 400356    5   15    50    50  207  861 13  3 83  1  0
 1  0 403256 601568  18892 400496    0    0  1048   364  337 1903  5  7  0 88  0
 0  1 403256 600816  19640 400568    0    0   748     0  259 1142  6  4  0 90  0
 0  1 403256 600300  20116 400800    0    0   476     0  196  630  8  5  0 87  0
 0  1 403256 599328  20792 400792    0    0   676     0  278 1401  7  5  0 88  0
[user@fedora9 ~]$
```


----------


参考：

- [Linux Performance Measurements using vmstat](https://www.thomas-krenn.com/en/wiki/Linux_Performance_Measurements_using_vmstat)