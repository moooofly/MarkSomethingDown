# Linux Cached Memory

We were running a [**BMC firmware** (If you don’t know already, find out what it is?)](https://en.wikipedia.org/wiki/Intelligent_Platform_Management_Interface#Baseboard_management_controller) on top of our custom linux kernel and one fine day we realized that we have developed so many applications that we forgot checking how much memory these applications were using and we were running low on available RAM.

The first thing we saw was, the cached memory in **/proc/meminfo** was too high.

We eventually found what cached was and how we can limit/free it at will. So, lets take a stroll down the memory lane…

## What is Cached

In linux there are two things that are part of the of **file page cache**. The kernel caches
- The **files** that your processes access
- The **RAM based filesystems**, such as **tmpfs**/**ramfs**, are also part of the cache.

### RAM File System Caching

The [kernel tmpfs documentation](https://www.kernel.org/doc/Documentation/filesystems/tmpfs.txt) notes that
```
Since tmpfs lives completely in the page cache and on swap, all tmpfs
pages currently in memory will show up as cached. It will not show up
as shared or something like that

tmpfs has three mount options for sizing:

size:      The limit of allocated bytes for this tmpfs instance. The 
           default is half of your physical RAM without swap. 
           
           **If you oversize your tmpfs instances the machine will deadlock
           since the OOM handler will not be able to free that memory.**
           
nr_blocks: The same as size, but in blocks of PAGE_CACHE_SIZE.
nr_inodes: The maximum number of inodes for this instance. The default
           is half of the number of your physical RAM pages, or (on a
           machine with highmem) the number of lowmem RAM pages,
           whichever is the lower.
```
To be more clear, **any file that you place in your tmpfs will be in the cache forever, till it is deleted.** So the downside of tmpfs is, if you oversize the tmpfs eventually the OOM killer will start killing your running processes to get more RAM space, as it cannot free the files in tmpfs from cache unless you have swap (in embedded linux, you don’t have it – which means you are stuck forever in RAM)

### File Caching

**Anytime a process opens a file, the kernel automatically caches the file** there by making subsequent I/O calls to the file faster as it is done in the memory instead of directly on the disk.

### Eviction from the cache

Kernel automatically reclaims the pagecache when other processes or kernel needs memory. However, if the file is in use then the kernel wouldn’t be able to free the cache. And if the whole system is running short of memory, then the Out-Of-Memory (OOM) killer is called.

After a file is closed, it may still be part of the cache in case the process again tries to open it. But this is in a reclaimable part of the cache, which means if more memory is necessary this closed file will be evicted from the memory.

**A file in RAM filesystem will be part of the cache. Unneeded pages will be swapped to swap. If you don’t have swap, then unless you delete the file it is part of cache.**

## Best practices for reducing the cache

### RAM based FileSystem

When you mount, you can restrict the maximum size of the **tmpfs**. Note that you cannot limit it using **ramfs**. Because, [as noted in the tmpfs documentation](https://www.kernel.org/doc/Documentation/filesystems/tmpfs.txt)
```shell
If you compare it to ramfs (which was the template to create tmpfs)
you gain swapping and limit checking. Another similar thing is the RAM
disk (/dev/ram*), which simulates a fixed size hard disk in physical
RAM, where you have to create an ordinary filesystem on top. Ramdisks
cannot swap and you do not have the possibility to resize them. 
```
So the bottom line is: **tmpfs is a better option than ramfs, as you can add swap and also limit the mount size.**
```shell
# mount -t tmpfs -o size=60M tmpfs /tmp -> Size limited to 60M
```

### Operating on a file

If you are operating on a file of substantial size, then **after the end of the file access you can use [posix_fadvise](https://linux.die.net/man/2/posix_fadvise) to evict this file from the cache immediately**. This is a much safer option, as you will only evict your file instead of the other files that are being cached.
```c
#include <fcntl.h>

int posix_fadvise(int fd, off_t offset, off_t len, int advice);

If you want to drop the file from cache, the `advice` should be:
  POSIX_FADV_DONTNEED : The specified data will not be accessed in the near future.
```

### Global Drop Cache

If you want to drop caches of the whole system then write
- `echo 1 > /proc/sys/vm/drop_caches` - Will free **pagecache**
- `echo 2 > /proc/sys/vm/drop_caches` - Will free **dentries** and **inodes**
- `echo 3 > /proc/sys/vm/drop_caches` - Will free **pagecache**, **dentries** and **inodes**

### Tuning /proc/sys/vm/

You can control the below values which would help you in freeing up the cache more often.

- `/proc/sys/vm/swappiness`

This controls **how agressively kernel will swap memory pages out**.
```shell
Higher Values = Faster Swapping
Lower Values = Slow to clean up Swap
Value of 0 = Wait to cleanup swap until the high water mark in a zone is reached
Default = 60
```

- `/proc/sys/vm/vfs_cache_pressure`

This controls **how the kernel will reclaim memory that is used for caching of directory and inode objects**. Note: This doesn’t affect the file caching pressure
```shell
Higher Value = Faster Cleanup of directory and inode objects  
Lower Value = Slow to clean up directory and inode objects  
Value of 0 = Never clean up directory and inode objects  
Default = 100 (this is a fair rate of cleanup)
```

- `/proc/sys/vm/min_free_kbytes`

This forces the VM to keep a minimum number of configured kilobytes free for the atomic operations within the kernel. *Note: This is not the free memory that your userspace processes can use.*

> Choose a value that would be needed by your kernel depending on your kernel configurations

- `/proc/sys/vm/dirty_background_ratio`

The percentage value of memory that can be filled with **dirty pages** before `pdflush` begins to write them.

> If you have huge memory, then the default percentage may be too high. You can lower it

- `/proc/sys/vm/dirty_expire_centisecs`

The hundreth of the second after which data will be considered to be expired from the page cache and will be written at the next opportunity.

> Reducing will clean up the page cache, but will trigger IO congestion

- `/proc/sys/vm/dirty_ratio`

The percentage value of memory that can be filled with **dirty pages** before the processes begin to write them.

> Reducing this will kick in `pdflush` when a process is writing out huge files. This would **block** the IO for the process momentarily

- `/proc/sys/vm/dirty_writeback_centisecs`

The hundredth of a second after which the `pdflush` wakes up to write data to disk.

> Can reduce the value to avoid data loss, however this will have an effect of IO congestion

## Caching Tools

### fincore

This is a nice utility, part of [linux-ftools](https://code.google.com/p/linux-ftools/). You have to give the file name as input, and it will **stats for the files that are in the cache now**. Hint: Running it on files that are part of tmpfs will show that these files are always in cache
```shell
# fincore --pages=false --summarize --only-cached <file_name>
```

```shell
fincore [options] files...

  --pages=false      Do not print pages
  --summarize        When comparing multiple files, print a summary report
  --only-cached      Only print stats for files that are actually in cache.

root@xxxxxx:/var/lib/mysql/blogindex# fincore --pages=false --summarize --only-cached * 
stats for CLUSTER_LOG_2010_05_21.MYI: file size=93840384 , total pages=22910 , cached pages=1 , cached size=4096, cached perc=0.004365 
stats for CLUSTER_LOG_2010_05_22.MYI: file size=417792 , total pages=102 , cached pages=1 , cached size=4096, cached perc=0.980392 
stats for CLUSTER_LOG_2010_05_23.MYI: file size=826368 , total pages=201 , cached pages=1 , cached size=4096, cached perc=0.497512 
stats for CLUSTER_LOG_2010_05_24.MYI: file size=192512 , total pages=47 , cached pages=1 , cached size=4096, cached perc=2.127660 
stats for CLUSTER_LOG_2010_06_03.MYI: file size=345088 , total pages=84 , cached pages=43 , cached size=176128, cached perc=51.190476 
stats for CLUSTER_LOG_2010_06_04.MYD: file size=1478552 , total pages=360 , cached pages=97 , cached size=397312, cached perc=26.944444 
stats for CLUSTER_LOG_2010_06_04.MYI: file size=205824 , total pages=50 , cached pages=29 , cached size=118784, cached perc=58.000000 
stats for COMMENT_CONTENT_2010_06_03.MYI: file size=100051968 , total pages=24426 , cached pages=10253 , cached size=41996288, cached perc=41.975764 
stats for COMMENT_CONTENT_2010_06_04.MYD: file size=716369644 , total pages=174894 , cached pages=79821 , cached size=326946816, cached perc=45.639645 
stats for COMMENT_CONTENT_2010_06_04.MYI: file size=56832000 , total pages=13875 , cached pages=5365 , cached size=21975040, cached perc=38.666667 
stats for FEED_CONTENT_2010_06_03.MYI: file size=1001518080 , total pages=244511 , cached pages=98975 , cached size=405401600, cached perc=40.478751 
stats for FEED_CONTENT_2010_06_04.MYD: file size=9206385684 , total pages=2247652 , cached pages=1018661 , cached size=4172435456, cached perc=45.321117 
stats for FEED_CONTENT_2010_06_04.MYI: file size=638005248 , total pages=155763 , cached pages=52912 , cached size=216727552, cached perc=33.969556 
stats for FEED_CONTENT_2010_06_04.frm: file size=9840 , total pages=2 , cached pages=3 , cached size=12288, cached perc=150.000000 
stats for PERMALINK_CONTENT_2010_06_03.MYI: file size=1035290624 , total pages=252756 , cached pages=108563 , cached size=444674048, cached perc=42.951700 
stats for PERMALINK_CONTENT_2010_06_04.MYD: file size=55619712720 , total pages=13579031 , cached pages=6590322 , cached size=26993958912, cached perc=48.533080 
stats for PERMALINK_CONTENT_2010_06_04.MYI: file size=659397632 , total pages=160985 , cached pages=54304 , cached size=222429184, cached perc=33.732335 
stats for PERMALINK_CONTENT_2010_06_04.frm: file size=10156 , total pages=2 , cached pages=3 , cached size=12288,
```

### vmtouch

If you want to **modify the files in cache**, then you can use [vmtouch](https://github.com/hoytech/vmtouch/blob/master/vmtouch.pod). With vmtouch you can

- Evict files from cache
- Keep a file always in cache, without giving it a possibility to be evicted

## Conclusion

Damn, be extra careful with what you keep in your cache memory. That determines whether you system works as you want it to or not, esp in an embedded linux world.

I will write few more blog entries on how we eventually managed to handle all the memory issues in our firmware.

If you have more tools or thoughts on handling cache better, let me know in the comments.




