# HugePage

> 未完成

## redis 中的 Transparent Huge Pages 问题

> 以下内容取自[这里](https://cachecloud.github.io/2017/02/16/Redis%E7%9A%84Linux%E7%B3%BB%E7%BB%9F%E4%BC%98%E5%8C%96/)

Redis 在启动时可能会看到如下日志：

```
WARNING you have Transparent Huge Pages (THP) support enabled in your kernel. This will create latency and memory usage issues with Redis. To fix this issue run the command 'echo never > /sys/kernel/mm/transparent_hugepage/enabled' as root, and add it to your /etc/rc.local in order to retain the setting after a reboot. Redis must be restarted after THP is disabled.
```

从提示上看，Redis 建议

- 修改 **Transparent Huge Pages (THP)** 的相关配置为 never ，因为默认 THP 是开启的；
- 开启 THP 后会导致 Redis 的 latency 和 memory 使用问题；
- 在修改过后需要重启 Redis 实例；

背景信息：Linux kernel 在 2.6.38 内核增加了 Transparent Huge Pages (THP) 特性，以支持**大内存页 (2MB)** 分配，**默认开启**。在开启状态下，可以降低 `fork` 子进程的速度，但 `fork` 之后，由于每个内存页从原来 4KB 变为 2MB ，会大幅增加重写期间父进程内存消耗。同时每次写命令引起的复制内存页 (COW) 单位放大了 512 倍，会拖慢写操作的执行时间，导致大量**写操作慢查询**。

这就是为何简单的 `incr` 命令也会出现在慢查询中的原因。

禁用方法如下：

```
echo never >  /sys/kernel/mm/transparent_hugepage/enabled
```

而且为了使机器重启后 THP 配置依然生效，可以在 `/etc/rc.local` 中追加 `echo never > /sys/kernel/mm/transparent_hugepage/enabled` 。

在设置 THP 配置时需要注意：有些 Linux 的发行版本没有将 THP 放到 `/sys/kernel/mm/transparent_hugepage/enabled` 中，例如 Red Hat 6 以上的 THP 配置放到 `/sys/kernel/mm/redhat_transparent_hugepage/enabled` 中。而 Redis 源码中在检查 THP 时把 THP 位置写死了：

```c
FILE *fp = fopen("/sys/kernel/mm/transparent_hugepage/enabled","r");
if (!fp) return 0;
```

所以在发行版中，虽然没有 THP 的日志提示，但是依然存在 THP 所带来的问题。

```
echo never >  /sys/kernel/mm/redhat_transparent_hugepage/enabled
```


----------


## [Redis latency problems troubleshooting](https://redis.io/topics/latency)

- **latency** is the maximum delay between the time a client issues a command and the time the reply to the command is received by the client.
- Usually Redis processing time is extremely low, in the **sub microsecond range**.

### Latency induced by transparent huge pages

Unfortunately when a Linux kernel has transparent huge pages enabled, Redis incurs to a big latency penalty after the `fork` call is used in order to persist on disk. Huge pages are the cause of the following issue:

- `Fork` is called, **two processes with shared huge pages are created**.
- In a busy instance, **a few event loops runs will cause commands to target a few thousand of pages, causing the copy on write of almost the whole process memory.**
- This will result in big latency and big memory usage.

Make sure to **disable transparent huge pages** using the following command:

```
echo never > /sys/kernel/mm/transparent_hugepage/enabled
```

> 补充说明：能否将不使用持久化特性的 Redis 上出现延迟的情况归为上面第二种情况？



## [Huge Pages and Transparent Huge Pages](https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/6/html/Performance_Tuning_Guide/s-memory-transhuge.html)

> RHEL 6

Memory is managed in **blocks** known as `pages`. A page is **4096** bytes. **1MB** of memory is equal to 256 pages; **1GB** of memory is equal to 256,000 pages, etc. CPUs have a built-in `memory management unit` that contains a list of these pages, with each page referenced through a `page table entry`.

There are two ways to enable the system to manage large amounts of memory:

- **Increase the number of page table entries in the hardware memory management unit**
- **Increase the page size**

The first method is expensive, since the hardware memory management unit in a modern processor **only supports hundreds or thousands of page table entries**. Additionally, hardware and memory management algorithms that work well with thousands of pages (megabytes of memory) may have difficulty performing well with millions (or even billions) of pages. This results in performance issues: **when an application needs to use more memory pages than the memory management unit supports, the system falls back to slower, software-based memory management, which causes the entire system to run more slowly**.

Red Hat Enterprise Linux 6 implements the second method via the use of `huge pages`.

Simply put, huge pages are blocks of memory that come in **2MB** and **1GB** sizes. The page tables used by the 2MB pages are suitable for managing multiple gigabytes of memory, whereas the page tables of 1GB pages are best for scaling to terabytes of memory. 

**Huge pages can be difficult to manage manually**, and often require significant changes to code in order to be used effectively. As such, Red Hat Enterprise Linux 6 also implemented the use of `transparent huge pages (THP)`. **`THP` is an abstraction layer that automates most aspects of creating, managing, and using huge pages.**

THP hides much of the complexity in using huge pages from system administrators and developers. As the goal of THP is improving performance, its developers (both from the community and Red Hat) have tested and optimized THP across a wide range of systems, configurations, applications, and workloads. This allows the default settings of THP to improve the performance of most system configurations. However, **THP is not recommended for database workloads.**

**THP can currently only map anonymous memory regions such as heap and stack space.**

### Configure Huge Pages

Huge pages **require** contiguous areas of memory, so **allocating them at boot is the most reliable method** since memory has not yet become fragmented. To do so, add the following parameters to the kernel boot command line:

**Huge pages kernel options**

#### hugepages

Defines the **number of persistent huge pages** configured in the kernel at boot time. The default value is 0. It is only possible to allocate (or deallocate) huge pages if there are sufficient physically contiguous free pages in the system. **Pages reserved by this parameter cannot be used for other purposes.**

Default size huge pages can be dynamically allocated or deallocated by changing the value of the `/proc/sys/vm/nr_hugepages` file.

In a **NUMA** system, **huge pages** assigned with this parameter **are divided equally between nodes**. You can assign huge pages to specific nodes at runtime by changing the value of the node's `/sys/devices/system/node/node_id/hugepages/hugepages-1048576kB/nr_hugepages` file.

For more information, read the relevant kernel documentation, which is installed in `/usr/share/doc/kernel-doc-kernel_version/Documentation/vm/hugetlbpage.txt` by default. This documentation is available only if the `kernel-doc` package is installed.

#### hugepagesz

Defines **the size of persistent huge pages** configured in the kernel at boot time. Valid values are **2 MB** and **1 GB**. The default value is 2 MB.

#### default_hugepagesz

Defines **the default size of persistent huge pages** configured in the kernel at boot time. Valid values are **2 MB** and **1 GB**. The default value is 2 MB.



## [CONFIGURING TRANSPARENT HUGE PAGES](https://access.redhat.com/documentation/en-US/Red_Hat_Enterprise_Linux/7/html/Performance_Tuning_Guide/sect-Red_Hat_Enterprise_Linux-Performance_Tuning_Guide-Configuring_transparent_huge_pages.html)

## [Disable Transparent Huge Pages (THP)](https://docs.mongodb.com/manual/tutorial/transparent-huge-pages/)


## [Transparent Hugepage Support](https://www.kernel.org/doc/Documentation/vm/transhuge.txt)


----------


## /sys/kernel/mm/hugepages/ 目录解析

`/sys/kernel/mm/hugepages/` contains a number of subdirectories of the form `hugepages-<size>kB`, where `<size>` is the page size of the hugepages supported by the kernel/CPU combination.

Under these directories are a number of files:

- nr_hugepages
- nr_overcommit_hugepages
- free_hugepages
- surplus_hugepages
- resv_hugepages

See `Documentation/vm/hugetlbpage.txt` for details.

```
[root@wg-public-rediscluster-119: ~]# ll /sys/kernel/mm/hugepages/
total 0
drwxr-xr-x 2 root root 0 Aug  7 18:09 hugepages-1048576kB
drwxr-xr-x 2 root root 0 Aug  7 18:09 hugepages-2048kB
[root@wg-public-rediscluster-119: ~]#
[root@wg-public-rediscluster-119: ~]# ll /sys/kernel/mm/transparent_hugepage/
total 0
-rw-r--r-- 1 root root 4096 Aug  4 21:02 defrag
-rw-r--r-- 1 root root 4096 Aug  4 21:02 enabled
drwxr-xr-x 2 root root    0 Aug  7 18:10 khugepaged
-rw-r--r-- 1 root root 4096 Aug  7 18:10 use_zero_page
[root@wg-public-rediscluster-119: ~]#
[root@wg-public-rediscluster-119: ~]# cat /sys/kernel/mm/transparent_hugepage/defrag
always madvise [never]
[root@wg-public-rediscluster-119: ~]# cat /sys/kernel/mm/transparent_hugepage/enabled
always madvise [never]
[root@wg-public-rediscluster-119: ~]# cat /sys/kernel/mm/transparent_hugepage/use_zero_page
1
[root@wg-public-rediscluster-119: ~]#
[root@wg-public-rediscluster-119: ~]# ll /sys/kernel/mm/transparent_hugepage/khugepaged/
total 0
-rw-r--r-- 1 root root 4096 Aug  7 18:11 alloc_sleep_millisecs
-rw-r--r-- 1 root root 4096 Aug  7 18:11 defrag
-r--r--r-- 1 root root 4096 Aug  7 18:11 full_scans
-rw-r--r-- 1 root root 4096 Aug  7 18:11 max_ptes_none
-r--r--r-- 1 root root 4096 Aug  7 18:11 pages_collapsed
-rw-r--r-- 1 root root 4096 Aug  7 18:11 pages_to_scan
-rw-r--r-- 1 root root 4096 Aug  7 18:11 scan_sleep_millisecs
[root@wg-public-rediscluster-119: ~]#
```

## /proc/meminfo 中和 hugepages 相关的内容

```
[root@wg-public-rediscluster-119: ~]# cat /proc/meminfo |grep "Huge"
AnonHugePages:  21479424 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
[root@wg-public-rediscluster-119: ~]#
```

## sysctl 中和 hugepages 相关的内容

```
[root@wg-public-rediscluster-119: ~]# sysctl -a|grep "hugepages"
vm.hugepages_treat_as_movable = 0
vm.nr_hugepages = 0
vm.nr_hugepages_mempolicy = 0
vm.nr_overcommit_hugepages = 0
[root@wg-public-rediscluster-119: ~]#
```


## 其它

- [Disabling Defrag for Red Hat and CentOS Systems](https://my.vertica.com/docs/7.2.x/HTML/index.htm#Authoring/InstallationGuide/BeforeYouInstall/defrag.htm)
- [Enabling or Disabling Transparent Hugepages](https://my.vertica.com/docs/7.2.x/HTML/index.htm#Authoring/InstallationGuide/BeforeYouInstall/transparenthugepages.htm)
- [vm/hugetlbpage.txt](https://www.kernel.org/doc/Documentation/vm/hugetlbpage.txt)