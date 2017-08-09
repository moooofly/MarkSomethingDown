# Memory

## Virtual Memory

**Virtual memory 使用 disk 作为 RAM 的扩展**，因此可用内存的有效大小相应的变大了；kernel 会将当前不再使用的内存块的内容写到 hard disk 上，以便 memory 能够被用于其他目的；

当（上述内存块中的）原始内容再次被需要时，将被再次读回到内存中；而上述行为对于用户来说是完全透明到；运行在 Linux 下到程序只会看到更多的可用内存，而不会发现内存的哪部分位于 disk 上；当然，针对 hard disk 的读写会比使用真正的内存更慢（大概慢 1000 倍），因此程序跑的不会像使用纯内存时那么快；**用作 virtual memory 的那部分 hard disk 被称作 `swap space`** ；

## Virtual Memory Pages

Virtual memory 按 pages 进行划分；在 X86 架构上，每一个 virtual memory page 为 4KB 大小；当 kernel 将 memory 写入 disk 或从 disk 读出时，均按照 pages 进行操作；kernel 会将 memory pages 根据实际情况写入 **swap device** 或**文件系统**；

## Kernel Memory Paging

Memory `paging` 是一种常规活动，不要和 memory `swapping` 搞混淆；Memory `paging` 是指将 memory 按照**一定的时间间隔**、**同步（synching）**到 disk 的过程；只要时间运行长了，应用就会逐渐消耗光所有的 memory ；在某些时间点上，kernel 必须**扫描**（scan）memory 并**回收**（reclaim）不再使用的 pages ，以便其能够分配给其它应用使用； 

## The Page Frame Reclaim Algorithm (PFRA)

PFRA 负责**释放内存** ；PFRA 根据 page 类型选择要释放的 memory pages ；Page 类型包括：

- **Unreclaimable** – locked, kernel, reserved pages
- **Swappable** – anonymous memory pages
- **Syncable** – pages backed by a disk file
- **Discardable** – static pages, discarded pages

除了 “unreclaimable” 之外的其它类型 pages 均可以通过 PFRA 进行回收；
PFRA 中存在两个主要的 functions ；即 `kswapd` kernel 线程和 “`Low On Memory Reclaiming`” function ；

## kswapd

> The `kswapd` daemon is responsible for ensuring that memory stays free. It monitors the **pages_high** and **pages_low** watermarks in the kernel. If the amount of free memory is below **pages_low**, the `kswapd` process starts a scan to attempt to **free 32 pages at a time**. It repeats this process until the amount of free memory is above the **pages_high** watermark. 

`kswapd` 负责

- （用途）确保 memory 有空余（后面可知是确保 free list 有余量）；
- （周期）监控 kernel 中的高低两个水位，以便触发释放行为；
- （数量）每次尝试释放 32 个 pages ；

> The `kswapd` thread performs the following actions:
>
> - If the page is **unmodified**, it places the page on the **free list**.
> - If the page is **modified** and backed by a filesystem, it writes the contents of the page to **disk**.
> - If the page is **modified** and not backed up by any filesystem (**anonymous**), it writes the contents of the page to the **swap device**. 

`kswapd` 根据 page 的状态和特点，采取如下行动（回收对象）：

- 如果 page 未被修改，则将其放回 **free list** 中；
- 如果 page 被修改了，并且对应了文件系统中的文件，则将该 page 的内容写入 **disk** ；
- 如果 page 被修改了，但其并未对应文件系统中的文件（即 anonymous page），则将该 page 的内容写入 **swap device** ；

上述内容的隐藏知识点：

- 前两种情况中的 page 和 file buffer cache 有关；
- 后两种情况中的 page 均为 **dirty page** ；
- `kswapd` 既和 swap device 打交道，又和 disk 打交道；


> To maintain the **free list**, `kswapd` **steals** memory from the read/write buffers (`buff`) and assigns it to the **free list**. This is evident in the gradual decrease of the buffer cache (`buff`).
>
> The `kswapd` process then writes **dirty pages** to the **swap device** (`so`). This is evident in the fact that the amount of virtual memory utilized gradually increases (`swpd`). 

这段话说明：`kswapd` 会从 **Buffer Cache** 中**偷**内存，以便保证 free list 中的余量；

## Kernel Paging with pdflush 

> The `pdflush` daemon is responsible for **synchronizing** any pages associated with a file on a filesystem back to **disk**. In other words, when a file is modified in memory, the `pdflush` daemon writes it back to disk. 

**`pdflush` 负责将与 file 相关的 page 同步到 disk 上**；

> The `pdflush` daemon starts synchronizing **dirty pages** back to the filesystem when 10% of the pages in memory are dirty. This is due to a kernel tuning parameter called `vm.dirty_background_ratio`.

`pdflush` 的同步行为受内核参数 `vm.dirty_background_ratio` 的影响；

> The `pdflush` daemon works independently of the PFRA under most circumstances. When the kernel invokes the LMR algorithm, the LMR specifically forces `pdflush` to flush dirty pages in addition to other page freeing routines. 

大多数情况下，`pdflush` 独立于 PFRA 工作；kernel 触发 LMR 算法，LMR 会通过 `pdflush` 将 dirty pages 刷出，同时也会调用其它 page 释放程序；

> Under intense memory pressure in the **2.4 kernel**, the system would experience `swap thrashing`. This would occur when the PFRA would **steal** a page that an active process was trying to use. As a result, the process would have to reclaim that page only for it to be stolen again, creating a thrashing condition. This was **fixed in kernel 2.6** with the “**Swap Token**”, which prevents the PFRA from constantly stealing the same page from a process. 

在 2.4 内核中存在 `swap thrashing` 问题，在 2.6 内核中进行了修复；


## Conclusion

Virtual memory 性能监控由以下行为构成：

- **系统中 `major page faults` 越少，则 response times 将越小**，因为系统更多利用了 memory caches 而不是 disk caches ；
- **`free` memory 量比较低是一个好的现象**，其表明 caches 被高效使用了，除非存在持续向 swap device 和 disk 进行写入的情况；
- 若系统报告**在 swap device 上存在任何持续性读写活动**，则表示**当前系统存在 memory shortage 问题**；


## Types of Memory Pages

在 Linux kernel 中存在三种类型的 **memory pages** ；具体描述如下：

- **Read Pages** – These are pages of data read in via **disk** ([`MPF`](https://en.wikipedia.org/wiki/Memory-prediction_framework)) that are **read only** and **backed on disk**. These pages **exist in the `Buffer Cache`** and include **static files**, **binaries**, and **libraries** that do not change. The Kernel will continue to **page these into memory** as it needs them. If memory becomes short, the kernel will "**steal**" these pages and put them back on the `free list` causing an application to have to `MPF` to bring them back in.
- **Dirty Pages** – These are pages of data that have been **modified by the kernel** while in memory. These pages need to be **synced back to disk** at some point using the `pdflush` daemon. In the event of a memory shortage, `kswapd` (along with `pdflush`) will write these pages to disk in order to make more room in memory.
- **Anonymous Pages** – These are pages of data that do belong to a process, but do not have any file or backing store associated with them. They can't be synchronized back to disk. In the event of a memory shortage, `kswapd` writes these to the **swap device** as temporary storage until more RAM is free ("`swapping`" pages). 


----------


## kswapd v.s. pdflush

首先，它们存在的目的不同，`kswapd` 的作用是**管理内存**，`pdflush` 的作用是**同步内存和磁盘**；

数据写入磁盘前可能会被缓存在内存中，通常在满足如下三个原因的情况下，这些缓存才真正写入磁盘：

- **用户要求**缓存马上写入磁盘（`sync()` or `fsync()`）；
- **缓存过多**，超过一定阀值，需要写入磁盘；
- **内存吃紧**，需要将缓存写入磁盘以腾出地方；

上述原因使得 `kswapd` 和 `pdflush` 有交叉的地方，因此很多人混淆了它们。 

它们相同的地方都是**定期被唤醒**，都是以**守护进程（内核进程）**的形式存在；**`kswapd` 试图保证内存永远都是可满足用户要求的**，为了实现这种承诺，它必须采取一定的策略；**`pdflush` 试图保证内存和磁盘的数据是同步的**，不会因为缓存的原因使内存和磁盘的数据不同步，从而造成数据丢失或者损坏，为了实现这种承诺，它同样也要采取一定的策略。

**那么它们之间的交叉点在何处呢？**比如，在用户所要求的内存不能被满足，或者空闲内存的数量已经低于某一个值的时候，`kswapd` 被唤醒，它必须为用户的要求提供服务，因此试图换出一部分正在使用的内存，使之成为空闲内存以供用户使用；这时，**磁盘缓存 (disk caches 或直接描述为 file buffer cache 更容易理解)** 也是正在被使用的内存；因此，`kswapd` 需要将它们换出，这里的换出和匿名页面被换到**交换分区 (swap device)** 是一样的概念，**将磁盘缓存换到哪里呢？**当然哪里来哪里去了。**linux 不区分匿名页面对应的交换分区和真实文件的磁盘缓存对应的磁盘文件分区**，实际上在将匿名页面写到交换分区的时候，也是按照写文件的形式进行的，读源代码的时候就会发现有一个 `address_space_operations` 结构体，里面的 `readpage` 和 `writepage` 就是读写页面的回调函数，**linux 的这个实现方式表明，写匿名页面和写 ext2 的缓存页面没有本质的区别**，仅仅换一下那几个 `address_space_operations` 里面的回调函数就行。因此 `kswapd` 也会将磁盘缓存回写到磁盘，和 `pdflush` 所作的工作一样，这就是它们交叉的地方，当然如果 `kswapd` 已经将页面写入了磁盘，就会清除掉页面的脏标志；这样，在 `pdflush` 扫描脏页的时候就不会二次**回写**了。 

既然 `kswapd` 和 `pdflush` 有联系，那么**联系它们的纽带是什么？**当然是内核中的 LRU 链表了；本来需要通过 `pdflush` 写入磁盘的页面，也许会通过 `kswapd` 写入，**如何让 `kswapd` 看到 `pdflush` 负责的页面呢？**实际上 linux 并没有刻意关注这个事情，内核那么复杂，如果这么细致的考虑问题谁都会发疯的。因此 linux 采用了更加宏伟的方式，就是将事情抽象，不再操心什么回写啊，内存释放之类的细节，而是抽象出了**内存管理**和**缓存管理**这些个模块，然后模块和模块之间建立一个耦合点，也可以理解成一个接口，这个东西就是 LRU 链表；**linux 规定，凡是想纳入内存管理范畴的内存物理页面都要加入 LRU 链表**，而 `kswapd` 就是内存管理的执行者，它操作的正是这个链表，这样它就不需要别的什么了，只需要告诉大家，你想让我管理，别让我去找你，你自己加入 LRU 链表吧，就这样而已。缓存管理模块当然想加入内存管理，因此所有的磁盘缓存页面都在加入缓存的同时加入了 LRU 链表，这样**缓存管理的执行者 `pdflush`** 和**内存管理的执行者 `kswapd`** 就不需要直接交互商量事情了，一个 LRU 链表解除了它们的耦合。 

linux 中到处都体现了这样的思想，它看似一个嗷嗷宏大的内核，实际上是高度模块化的，你不要觉得内核中有些东西好像杂糅在一起而被搞的焦头烂额，实际上仔细看看代码就会发现，它们之间的耦合点也就是一些很简单的结构，比如我前面文章提到的 `list_head` 或者 `kobject` 等等。不信的话再看看上面的 `kswapd` 和 `pdflush` ，**如果你想让内存加入缓存管理，那么就设置它为脏**（设置一个标志），并加入一棵 radix 树（**本质上 radix 树和链表没有区别**，都是一个连接数据结构，早期的内核版本中的缓存就是链表结构连接起来的）；**如果你想让内存页面加入内存管理，那么就加入 LRU 链表**，如果你想两个都加入呢？很简单，设为脏，加入 radix 树，再加入 LRU 链表，之后就不用管了，`kswapd` 和 `pdflush` 会各司其职的，前者查 LRU ，后者查 radix 和标志位，它们之间唯一需要交互的就是，一方做完工作后要让另一方看到；这实际上不是它们之间的交互，而仅仅是一项工作的收尾工作，或者说是汇报工作，你做完一件事总得有点效果吧。


----------

this chaper will cover the **page replacement daemon** `kswapd`, how it is implemented and what it's responsibilities are.

...

The second situation is where there is a single process with many file backed resident pages in the `inactive_list` that are being written to frequently. Processes and `kswapd` may go into a loop of constantly “laundering” these pages and placing them at the top of the `inactive_list` without freeing anything. In this case, few pages are moved from the `active_list` to `inactive_list` as the ratio between the two lists sizes remains not change significantly.

...

During system startup, a kernel thread called `kswapd` is started from `kswapd_init()` which continuously executes the function `kswapd()` in `mm/vmscan.c` which usually sleeps. **This daemon is responsible for reclaiming pages when memory is running low**. Historically, `kswapd` used to wake up every 10 seconds but now it is only woken by the **physical page allocator** when the `pages_low` number of free pages in a zone is reached.

It is this daemon that performs most of the tasks needed to maintain the **page cache** correctly, shrink **slab caches** and swap out processes if necessary. Unlike swapout daemons such, as Solaris, which are woken up with increasing frequency as there is memory pressure, `kswapd` keeps freeing pages until the `pages_high` watermark is reached. Under extreme memory pressure, processes will do the work of `kswapd` synchronously by calling `balance_classzone()` which calls `try_to_free_pages_zone()`. As shown in Figure 10.6, it is at try_to_free_pages_zone() where the **physical page allocator** synchonously performs the same task as `kswapd` when the zone is under heavy pressure.

When `kswapd` is woken up, it performs the following:

- Calls `kswapd_can_sleep()` which cycles through all zones checking the `need_balance` field in the `struct zone_t`. If any of them are set, it can not sleep;
- If it cannot sleep, it is removed from the `kswapd_wait` wait queue;
- Calls the functions `kswapd_balance()`, which cycles through all zones. It will free pages in a zone with `try_to_free_pages_zone()` if `need_balance` is set and will keep freeing until the `pages_high` watermark is reached;
- The task queue for `tq_disk` is run so that pages queued will be written out;
- Add `kswapd` back to the `kswapd_wait` queue and go back to the first step.


As stated in Section 2.6, there is now a `kswapd` for every memory node in the system. These daemons are still started from `kswapd()` and they all execute the same code except their work is confined to their local node. The main changes to the implementation of `kswapd` are related to the `kswapd-per-node` change.

The basic operation of `kswapd` remains the same. Once woken, it calls `balance_pgdat()` for the `pgdat` it is responsible for. `balance_pgdat()` has two modes of operation. When called with `nr_pages == 0`, it will continually try to free pages from each zone in the local pgdat until `pages_high` is reached. When `nr_pages` is specified, it will try and free either `nr_pages` or `MAX_CLUSTER_MAX * 8`, whichever is the smaller number of pages.

----------

## 杂七杂八

通过虚拟内存（Virtual Memory）能够把计算机内存空间扩展到硬盘，即物理内存（RAM）和硬盘的一部分空间（SWAP）组合在一起，共同作为虚拟内存为计算机提供一个连贯的虚拟内存空间；**好处**是拥有的内存 “变多了”，可以运行更多、更大的程序；**坏处**是把部分硬盘当内存用，会导致整体性能受到影响，因为硬盘读写速度要比内存慢几个数量级，并且在发生 SWAP 交换时会增加系统的负担。
 
在操作系统里，**虚拟内存被分成页**，在 x86 系统上，每个页大小是 4KB 。 Linux 内核读写虚拟内存是以“页”为单位进行操作的；把内存转移到硬盘交换空间（SWAP）和从硬盘交换空间读取到内存的时候都是按页来读写的。

物理内存（RAM）和 SWAP 的这种交换过程称为**页面交换（Paging）**；值得注意的是 paging 和 swapping 是两个完全不同的概念，国内很多参考书把这两个概念混为一谈，`swapping` 也翻译成交换，**在操作系统里是指把某程序完全交换到硬盘以腾出内存给新程序使用**，和 `paging` **只交换程序的部分（页面）**是两个不同的概念。纯粹的 swapping 在现代操作系统中已经很难看到了，因为把整个程序交换到硬盘的办法既耗时又费力而且没必要，现代操作系统基本都是 paging 或者 paging/swapping 混合，swapping 最初是在 Unix system V 上实现的。

**`kswapd` 是内核回收内存的线程，即便不使能（创建）SWAP 分区，`kswapd` 还是要跑，因为还有 file buffer cache 相关的内存需要进行回收**；

系统每过一定时间就会唤醒 `kswapd` 看看内存是否紧张，如果不紧张，则睡眠；

内存不足的时候，`kswapd` 和 `pdflush` 共同负责把数据写回硬盘并释放内存。

`kswapd` 的回收目标对象：

- **anonymous pages**
即匿名页，属于某个进程但是又和任何文件无关联，不能被同步到硬盘上，内存不足的时候由 kswapd 负责将它们写到交换分区并释放内存；

- **file buffer cache**
即和文件相关的 pages ；

针对上述两种对象的回收倾向性，可以通过调整 `vm.swappiness` 进行控制。swappiness 默认为 60 ，即更倾向回收 file buffer cache ；

Linux uses `kswapd` for virtual memory management such that pages that have been recently accessed are kept in memory and less active pages are paged out to disk.

the `kswapd` process regularly decreases the ages of unreferenced pages, and at the end they are paged out (moved out) to disk.


在内核 2.6.32 之前有个著名内核线程 `pdflush` ，用于 flush dirty page ，2.6.32 里面把 `pdflush` 也 flush 掉了（[Flushing out pdflush](https://lwn.net/Articles/326552/))，变成了 per-BDI (backing device info) 的 writeback flusher 线程（每个逻辑设备一个 flusher 线程），解决了 `pdflush` 线程可能阻塞在某个 block device 的问题，从而提高了性能；


在程序写入大量的脏页面，超过 `dirty_ratio` 的比例时候，由于 2.6.18 内核的页面回写是针对全局的，`pdflush` 工作的时候，会进入 `blk_congestion_wait`，导致所以 buffer IO 写被阻塞。阻塞时间在 16us-10ms 不等，不少是 ms 级别的，导致非常严重的性能下降。

删除大文件的时候性能也会下降严重，跟这个有关系吗？**删除大文件的时候系统会把这个文件的相关页面全部回写**，就是同一问题。

linux 2.6 中 pagecache 回写磁盘的 `pdflush` 线程的动态调整算法比较有意思，有一个最大和最小的默认值，如果所有已存在的线程持续运行超过 1s，就动态创建一个 `pdflush` 线程，如果一个 `pdflush` 线程睡眠时间超过 1s，就终止这个线程，但最大和最小线程不能超过指定的默认值。



----------


## active memory v.s. inactive memory

```
       -a, --active
              Display active and inactive memory, given a 2.5.41 kernel or better.
...
   Memory
       ...
       inact: the amount of inactive memory.  (-a option)
       active: the amount of active memory.  (-a option)
```

```
root@vagrant-ubuntu-trusty:~] $ vmstat -a 1 5
procs -----------memory---------- ---swap-- -----io---- -system-- ------cpu-----
 r  b   swpd   free  inact active   si   so    bi    bo   in   cs us sy id wa st
 0  0    736   7768 349936  94252    0    0     7     9  147   46  0  1 99  0  0
 0  0    736   7756 349920  94256    0    0     0     0  149  301  0  0 100  0  0
 0  0    736   7756 349920  94256    0    0     0     0  151  298  0  1 99  0  0
 0  0    736   7756 349920  94256    0    0     0     0  151  297  0  0 100  0  0
 1  0    736   7524 349920  94268    0    0     0     0  200  501  0  4 96  0  0
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $ cat /proc/meminfo |grep -i "active"
Active:            94332 kB
Inactive:         349920 kB
Active(anon):      29884 kB
Inactive(anon):    43156 kB
Active(file):      64448 kB
Inactive(file):   306764 kB
root@vagrant-ubuntu-trusty:~] $
```

the **page replacement policy** is frequently said to be a `Least Recently Used (LRU)`-based algorithm but this is not strictly speaking true as the lists are not strictly maintained in LRU order. The LRU in Linux consists of two lists called the `active_list` and `inactive_list`. The objective is for the `active_list` to contain the working set of all processes and the `inactive_list` to contain reclaim canditates. As all reclaimable pages are contained in just two lists and pages belonging to any process may be reclaimed, rather than just those belonging to a faulting process, the replacement policy is a global one.

the LRU lists consist of two lists called `active_list` and `inactive_list`. They are declared in `mm/page_alloc.c` and are protected by the `pagemap_lru_lock` spinlock. They, broadly speaking, store the “**hot**” and “**cold**” pages respectively, or in other words, the `active_list` contains all the working sets in the system and `inactive_list` contains reclaim canditates.

When caches are being shrunk, pages are moved from the `active_list` to the `inactive_list` by the function `refill_inactive()`. It takes as a parameter the number of pages to move, which is calculated in `shrink_caches()` as a ratio depending on `nr_pages`, the number of pages in `active_list` and the number of pages in `inactive_list`. The number of pages to move is calculated as

```
nr_pages * nr_active_pages / ((nr_inactive_pages + 1) * 2) 
```

This keeps the `active_list` about two thirds the size of the `inactive_list` and the number of pages to move is determined as a ratio based on how many pages we desire to swap out (`nr_pages`).

Pages are taken from the end of the `active_list`. If the `PG_referenced` flag is set, it is cleared and the page is put back at top of the `active_list` as it has been recently used and is still “hot”. This is sometimes referred to as rotating the list. If the flag is cleared, it is moved to the `inactive_list` and the `PG_referenced` flag set so that it will be quickly promoted to the `active_list` if necessary.


----------

## 参考

- [Linux System and Network Performance Monitoring](http://www.ufsdump.org/papers/oscon2009-linux-monitoring.pdf)
- [Page Frame Reclamation](https://www.kernel.org/doc/gorman/html/understand/understand013.html)
- [kswapd和pdflush](http://blog.csdn.net/dog250/article/details/5303269)
- [Linux内存点滴 用户进程内存空间](http://www.cnblogs.com/muahao/p/5974594.html)
- [Linux 性能监控](http://blog.csdn.net/tianlesoftware/article/details/6198780)

