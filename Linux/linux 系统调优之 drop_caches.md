# linux 系统调优之 drop_caches

查阅[资料](http://git.kernel.org/cgit/linux/kernel/git/torvalds/linux.git/tree/Documentation/sysctl/vm.txt?id=HEAD#n189)如下：

```shell
==============================================================

drop_caches

Writing to this will cause the kernel to drop clean caches, as well as
reclaimable slab objects like dentries and inodes.  Once dropped, their
memory becomes free.

To free pagecache:
	echo 1 > /proc/sys/vm/drop_caches
To free reclaimable slab objects (includes dentries and inodes):
	echo 2 > /proc/sys/vm/drop_caches
To free slab objects and pagecache:
	echo 3 > /proc/sys/vm/drop_caches

This is a non-destructive operation and will not free any dirty objects.
To increase the number of objects freed by this operation, the user may run
`sync' prior to writing to /proc/sys/vm/drop_caches.  This will minimize the
number of dirty objects on the system and create more candidates to be
dropped.

This file is not a means to control the growth of the various kernel caches
(inodes, dentries, pagecache, etc...)  These objects are automatically
reclaimed by the kernel when memory is needed elsewhere on the system.

Use of this file can cause performance problems.  Since it discards cached
objects, it may cost a significant amount of I/O and CPU to recreate the
dropped objects, especially if they were under heavy use.  Because of this,
use outside of a testing or debugging environment is not recommended.

You may see informational messages in your kernel log when this file is
used:

	cat (1234): drop_caches: 3

These are informational only.  They do not mean that anything is wrong
with your system.  To disable them, echo 4 (bit 3) into drop_caches.

==============================================================
```

关键点：

- Writing to drop_caches will cause the kernel to drop clean **caches**, as well as
reclaimable **slab objects** like **dentries** and **inodes**.
- This is a non-destructive operation and **will not free any dirty objects**.
- use outside of a testing or debugging environment is not recommended.


问题：

- slab objects 中的 dentries 和 inodes 对应什么？
- 通过 echo 修改 drop_caches 的值为 3 后，是否意味着一直处于 "**free slab objects and pagecache**" 状态？
- 什么时候需要设置 drop_caches 进行进行清理？如何清理？
- what's in the buffers and cache?
- 如何进行 swap 清理？


针对问题一：

> A **dentries** is a data structure that represents a directory. 
> An **inode** in your context is a data structure that represents a file. 
> 


针对问题二：

> It isn't sticky - you just write to the file to make it drop the caches and then it immediately starts caching again.
>
> Basically when you write to that file you aren't really changing a setting, you are issuing a command to the kernel. The kernel acts on that command (by dropping the caches) then carries on as before.
>
> The value you can read from /proc/sys/vm/drop_caches is whatever you put last, but it's not used anywhere, only the action of writing matters. The source code is in [fs/drop_caches.c](http://lxr.linux.no/linux+v3.0/fs/drop_caches.c).

针对问题三：

> A common case to "**manually flush**" those caches is purely for **benchmark** comparison: your first benchmark run may run with "empty" caches and so give poor results, while a second run will show much "better" results (due to the pre-warmed caches). By flushing your caches before any benchmark run, you're removing the "warmed" caches and so your benchmark runs are more "fair" to be compared with each other.
>
> 可以使用如下命令进行清理：
>> `free && sync && echo 3 > /proc/sys/vm/drop_caches && free`

针对问题四：

> Take a look at [linux-ftools](https://code.google.com/p/linux-ftools/) if you'd like to analyze the contents of the buffers & cache. 


针对问题五：

> `swapoff -a` -- disable swap
> `swapon -a` -- re-enable swap


参考：

- [Setting /proc/sys/vm/drop_caches to clear cache](http://unix.stackexchange.com/questions/17936/setting-proc-sys-vm-drop-caches-to-clear-cache)
- [How do you empty the buffers and cache on a Linux system?](http://unix.stackexchange.com/questions/87908/how-do-you-empty-the-buffers-and-cache-on-a-linux-system)
- [what are pagecache, dentries, inodes?](http://stackoverflow.com/questions/29870068/what-are-pagecache-dentries-inodes)





