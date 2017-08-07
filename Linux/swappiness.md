# swappiness

## [sysctl/vm.txt/swappiness](https://www.kernel.org/doc/Documentation/sysctl/vm.txt)

This control is used to define how aggressive the kernel will swap memory pages. Higher values will increase agressiveness, lower values decrease the amount of swap.  A value of 0 instructs the kernel not to initiate swap until the amount of **free and file-backed pages** is less than the high water mark in a zone.

The default value is 60.

## [Swappiness](https://en.wikipedia.org/wiki/Swappiness)

> Swappiness is a Linux kernel parameter that controls the relative weight given to **[swapping out](https://en.wikipedia.org/wiki/Virtual_memory#Address_space_swapping) of [runtime memory](https://en.wikipedia.org/wiki/Memory_footprint)**, as opposed to **dropping [pages](https://en.wikipedia.org/wiki/Page_(computer_memory)) from the system [page cache](https://en.wikipedia.org/wiki/Page_cache)**. Swappiness can be set to values between 0 and 100 inclusive. A low value causes the kernel to avoid swapping; a higher value causes the kernel to try to use swap space. The default value is 60; **setting it higher** will increase performance of "hot" processes at the cost of making a return to inactive "cold" ones take a long pause, while **setting it lower** (even 0) may decrease response latency. Systems with more than adequate RAM for any expected task may want to drastically lower the setting.

| Value | Strategy |
| -- | -- |
| `vm.swappiness = 0` | The kernel will **swap only** to avoid an **out of memory** condition, when free memory will be below `vm.min_free_kbytes` limit. See the "[VM Sysctl documentation](https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/Documentation/sysctl/vm.txt)". |
| `vm.swappiness = 1` | Kernel version 3.5 and over, as well as Red Hat kernel version 2.6.32-303 and over: **Minimum amount of swapping** without disabling it entirely. |
| `vm.swappiness = 10` | This value is sometimes recommended to improve performance when sufficient memory exists in a system. |
| `vm.swappiness = 60` | The default value. |
| `vm.swappiness = 100` | The kernel will **swap aggressively**. |

With kernel version 3.5 and over, as well as kernel version 2.6.32-303 and over, it is likely better to use 1 for cases where 0 used to be optimal.

To temporarily set the swappiness in Linux, write the desired value (e.g. `10`) to `/proc/sys/vm/swappiness` using the following command, running as root user:

```
# Set the swappiness value as root
echo 10 > /proc/sys/vm/swappiness

# Alternatively, run this 
sysctl -w vm.swappiness=10

# Verify the change
cat /proc/sys/vm/swappiness
10

# Alternatively, verify the change
sysctl vm.swappiness
vm.swappiness = 10
```

Permanent changes are made in `/etc/sysctl.conf` via the following configuration line (inserted, if not present):

```
vm.swappiness = 10
```


## [Redis的Linux系统优化](https://cachecloud.github.io/2017/02/16/Redis%E7%9A%84Linux%E7%B3%BB%E7%BB%9F%E4%BC%98%E5%8C%96/)

swap 对于操作系统来比较重要，当物理内存不足时，可以 swap out 一部分内存页，已解燃眉之急。但世界上没有免费午餐，swap 空间由硬盘提供，对于需要**高并发**、**高吞吐**的应用来说，磁盘 IO 通常会成为系统瓶颈。在 Linux 中，并不是要等到所有物理内存都使用完才会使用到 swap ，**系统参数 swppiness 会决定操作系统使用 swap 的倾向程度**。swappiness 的取值范围是 0~100 ，swappiness 的值越大，说明操作系统可能使用 swap 的概率越高，swappiness 值越低，表示操作系统更加倾向于使用物理内存。swap 的默认值是 60 ，了解这个值的含义后，有利于 Redis 的性能优化。下表对 swappiness 的重要值进行了说明。


| swapniess | 策略 |
| -- | -- |
| 0 | Linux 3.5 以及以上：宁愿 OOM killer 也不用 swap<br/> Linux 3.4 以及更早：宁愿 swap 也不要 OOM killer |
| 1 | Linux 3.5 以及以上：宁愿 swap 也不要 OOM killer |
| 60 | 默认值 |
| 100 | 操作系统会主动地使用 swap |

从下表中可以看出，swappiness 参数在 Linux 3.5 版本前后的表现并不完全相同，Redis 运维人员在设置这个值需要关注当前操作系统的内核版本。


设置方法：

```
# 系统重启后会失效
echo {bestvalue} > /proc/sys/vm/swappiness
# 系统重启后仍有效
echo vm.swappiness={bestvalue} >> /etc/sysctl.conf
```

### 监控 swap 的方法

```
free -m
vmstat -w 1 30
cat /proc/{pid}/smaps
```

最后一种方法的意义在于：能够查看指定进程的 swap 使用情况，即基于内存块镜像信息查看 swap 使用量，求和后得到总量；


### 针对 redis 的最佳实践

如果 Linux>3.5 则设置 `vm.swapniess=1` ，否则设置 `vm.swapniess=0` ，从而实现如下两个目标：

- 物理内存充足时候，使 Redis 足够快；
- 物理内存不足时候，避免 Redis 死掉（如果当前 Redis 为高可用，死掉比阻塞更好）；


> 补充说明：
>
> - 本文没有说清楚“swappiness 值越低，表示操作系统更加倾向于使用物理内存”中的物理内存具体是指那一部分内存；
> - 没有给出基于 sysctl 命令进行设置的方法和设置效果；
> - 基于 smaps 中的内存块镜像信息获取 swap 数据的方法价值在何处没看出来；

## [Redis latency problems troubleshooting](https://redis.io/topics/latency)

### Latency induced by swapping (operating system paging)

Linux (and many other modern operating systems) is able to **relocate** memory pages from the memory to the disk, and vice versa, in order to use the system memory efficiently.

If a Redis page is moved by the kernel from the memory to the **swap file**, when the data stored in this memory page is used by Redis (for example accessing a key stored into this memory page) the kernel will stop the Redis process in order to move the page back into the main memory. **This is a slow operation involving random I/Os** (compared to accessing a page that is already in memory) and will result into anomalous latency experienced by Redis clients.

The kernel relocates Redis memory pages on disk mainly because of three reasons:

- The system is **under memory pressure since the running processes are demanding more physical memory than the amount that is available**. The simplest instance of this problem is simply Redis using more memory than the one available.
- The Redis instance **data set**, or part of the data set, **is mostly completely idle** (never accessed by clients), so the kernel could swap idle memory pages on disk. This problem is very rare since even a moderately slow instance will touch all the memory pages often, forcing the kernel to retain all the pages in memory.
- **Some processes are generating massive read or write I/Os on the system**. Because files are generally cached, it tends to put pressure on the kernel to increase the filesystem cache, and therefore generate swapping activity. Please note it includes Redis RDB and/or AOF background threads which can produce large files.

Fortunately Linux offers good tools to investigate the problem, so the simplest thing to do is when latency due to swapping is suspected is just to check if this is the case.

The first thing to do is to checking the amount of Redis memory that is swapped on disk. In order to do so you need to obtain the Redis instance pid:

```
$ redis-cli info | grep process_id
process_id:5454
```

Now enter the `/proc` file system directory for this process:

```
$ cd /proc/5454
```

Here you'll find a file called `smaps` that **describes the memory layout** of the Redis process (assuming you are using Linux 2.6.16 or newer). This file contains very detailed information about our process memory maps, and one field called `Swap` is exactly what we are looking for. However there is not just a single swap field since the `smaps` file contains the different memory maps of our Redis process (The memory layout of a process is more complex than a simple linear array of pages).

Since we are interested in all the memory swapped by our process the first thing to do is to grep for the `Swap` field across all the file:

```
$ cat smaps | grep 'Swap:'
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                 12 kB
Swap:                156 kB
Swap:                  8 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  4 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  4 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  4 kB
Swap:                  4 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
Swap:                  0 kB
```

If everything is 0 kB, or if there are sporadic 4k entries, everything is perfectly normal. Actually in our example instance (the one of a real web site running Redis and serving hundreds of users every second) there are a few entries that show more swapped pages. To investigate if this is a serious problem or not we change our command in order to also print the size of the memory map:

```
$ cat smaps | egrep '^(Swap|Size)'
Size:                316 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                  8 kB
Swap:                  0 kB
Size:                 40 kB
Swap:                  0 kB
Size:                132 kB
Swap:                  0 kB
Size:             720896 kB
Swap:                 12 kB
Size:               4096 kB
Swap:                156 kB
Size:               4096 kB
Swap:                  8 kB
Size:               4096 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:               1272 kB
Swap:                  0 kB
Size:                  8 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                 16 kB
Swap:                  0 kB
Size:                 84 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                  8 kB
Swap:                  4 kB
Size:                  8 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  4 kB
Size:                144 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  4 kB
Size:                 12 kB
Swap:                  4 kB
Size:                108 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
Size:                272 kB
Swap:                  0 kB
Size:                  4 kB
Swap:                  0 kB
```

As you can see from the output, there is a map of 720896 kB (with just 12 kB swapped) and 156 kB more swapped in another map: basically a very small amount of our memory is swapped so this is not going to create any problem at all.

> 这里给出了一种观察方法；

If instead a non trivial amount of the process memory is swapped on disk your latency problems are likely related to swapping. If this is the case with your Redis instance you can further verify it using the `vmstat` command:

```
$ vmstat 1
procs -----------memory---------- ---swap-- -----io---- -system-- ----cpu----
 r  b   swpd   free   buff  cache   si   so    bi    bo   in   cs us sy id wa
 0  0   3980 697932 147180 1406456    0    0     2     2    2    0  4  4 91  0
 0  0   3980 697428 147180 1406580    0    0     0     0 19088 16104  9  6 84  0
 0  0   3980 697296 147180 1406616    0    0     0    28 18936 16193  7  6 87  0
 0  0   3980 697048 147180 1406640    0    0     0     0 18613 15987  6  6 88  0
 2  0   3980 696924 147180 1406656    0    0     0     0 18744 16299  6  5 88  0
 0  0   3980 697048 147180 1406688    0    0     0     4 18520 15974  6  6 88  0
^C
```

The interesting part of the output for our needs are the two columns **si** and **so**, that counts the amount of memory swapped from/to the swap file. If you see non zero counts in those two columns then there is swapping activity in your system.

Finally, the `iostat` command can be used to **check the global I/O activity** of the system.

```
$ iostat -xk 1
avg-cpu:  %user   %nice %system %iowait  %steal   %idle
          13.55    0.04    2.92    0.53    0.00   82.95

Device:         rrqm/s   wrqm/s     r/s     w/s    rkB/s    wkB/s avgrq-sz avgqu-sz   await  svctm  %util
sda               0.77     0.00    0.01    0.00     0.40     0.00    73.65     0.00    3.62   2.58   0.00
sdb               1.27     4.75    0.82    3.54    38.00    32.32    32.19     0.11   24.80   4.24   1.85
```

If your latency problem is due to Redis memory being swapped on disk you need to lower the memory pressure in your system, either adding more RAM if Redis is using more memory than the available, or avoiding running other memory hungry processes in the same system.


## [OOM relation to vm.swappiness=0 in new kernel](https://www.percona.com/blog/2014/04/28/oom-relation-vm-swappiness0-new-kernel/)

本文讲述了 `vm.swappiness=0` 含义变化后导致的 OOM 问题；其中给出了变更时的 commit 信息；


----------


## 其它

- [Linux Swap Space](http://www.linuxjournal.com/article/10678)

