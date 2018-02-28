# Linux proc-pid-xxx introduction

[/proc/[pid]/auxv](#auxv)  
[/proc/[pid]/cmdline](#cmdline)  
[/proc/[pid]/comm](#comm)  
[/proc/[pid]/cwd](#cwd)  
[/proc/[pid]/environ](#environ)  
[/proc/[pid]/exe](#exe)  
[/proc/[pid]/fd](#fd)  
[/proc/[pid]/latency](#latency)  
[/proc/[pid]/limits](#limits)  
[/proc/[pid]/maps](#maps)  
[/proc/[pid]/root](#root)  
[/proc/[pid]/stack](#stack)  
[/proc/[pid]/statm](#statm)  
[/proc/[pid]/syscall](#syscall)  
[/proc/[pid]/wchan](#wchan)  

## auxv

`/proc/[pid]/auxv`包含传递给进程的`ELF`解释器信息，格式是每一项都是一个`unsigned long`长度的`ID`加上一个`unsigned long`长度的值。最后一项以连续的两个`0x00`开头。举例如下：  

    # hexdump -x /proc/2948/auxv
    0000000    0021    0000    0000    0000    0000    1a82    7ffd    0000
    0000010    0010    0000    0000    0000    dbf5    1fc9    0000    0000
    0000020    0006    0000    0000    0000    1000    0000    0000    0000
    0000030    0011    0000    0000    0000    0064    0000    0000    0000
    0000040    0003    0000    0000    0000    2040    4326    7f4a    0000
    0000050    0004    0000    0000    0000    0038    0000    0000    0000
    0000060    0005    0000    0000    0000    0009    0000    0000    0000
    0000070    0007    0000    0000    0000    f000    4303    7f4a    0000
    0000080    0008    0000    0000    0000    0000    0000    0000    0000
    0000090    0009    0000    0000    0000    8e67    4327    7f4a    0000
    00000a0    000b    0000    0000    0000    0000    0000    0000    0000
    00000b0    000c    0000    0000    0000    0000    0000    0000    0000
    00000c0    000d    0000    0000    0000    0000    0000    0000    0000
    00000d0    000e    0000    0000    0000    0000    0000    0000    0000
    00000e0    0017    0000    0000    0000    0000    0000    0000    0000
    00000f0    0019    0000    0000    0000    3de9    1a80    7ffd    0000
    0000100    001f    0000    0000    0000    4fe5    1a80    7ffd    0000
    0000110    000f    0000    0000    0000    3df9    1a80    7ffd    0000
    0000120    0000    0000    0000    0000    0000    0000    0000    0000
    0000130
解析这个文件可以参考这段[代码](http://www.wienand.org/junkcode/linux/read-auxv.c)。

## cmdline

`/proc/[pid]/cmdline`是一个只读文件，包含进程的完整命令行信息。如果这个进程是`zombie`进程，则这个文件没有任何内容。举例如下：    

    # ps -ef | grep 2948
    root       2948      1  0 Nov05 ?        00:00:04 /usr/sbin/libvirtd --listen

    # cat /proc/2948/cmdline
    /usr/sbin/libvirtd--listen

## comm

`/proc/[pid]/comm`包含进程的命令名。举例如下：  

    # cat /proc/2948/comm
    libvirtd

## cwd

`/proc/[pid]/cwd`是进程当前工作目录的符号链接。举例如下：  

    # ls -lt /proc/2948/cwd
    lrwxrwxrwx 1 root root 0 Nov  9 12:14 /proc/2948/cwd -> /


## environ

`/proc/[pid]/environ`显示进程的环境变量。举例如下：  

    # strings /proc/2948/environ
    LANG=POSIX
    LC_CTYPE=en_US.UTF-8
    PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
    NOTIFY_SOCKET=@/org/freedesktop/systemd1/notify
    LIBVIRTD_CONFIG=/etc/libvirt/libvirtd.conf
    LIBVIRTD_ARGS=--listen
    LIBVIRTD_NOFILES_LIMIT=2048

## exe

`/proc/[pid]/exe`为实际运行程序的符号链接。举例如下：  

    # ls -lt /proc/2948/exe
    lrwxrwxrwx 1 root root 0 Nov  5 13:04 /proc/2948/exe -> /usr/sbin/libvirtd

## fd

`/proc/[pid]/fd`是一个目录，包含进程打开文件的情况。举例如下：  

    # ls -lt /proc/3801/fd
    total 0
    lrwx------. 1 root root 64 Apr 18 16:51 0 -> socket:[37445]
    lrwx------. 1 root root 64 Apr 18 16:51 1 -> socket:[37446]
    lrwx------. 1 root root 64 Apr 18 16:51 10 -> socket:[31729]
    lrwx------. 1 root root 64 Apr 18 16:51 11 -> socket:[34562]
    lrwx------. 1 root root 64 Apr 18 16:51 12 -> socket:[39978]
    lrwx------. 1 root root 64 Apr 18 16:51 13 -> socket:[34574]
    lrwx------. 1 root root 64 Apr 18 16:51 14 -> socket:[39137]
    lrwx------. 1 root root 64 Apr 18 16:51 15 -> socket:[39208]
    lrwx------. 1 root root 64 Apr 18 16:51 16 -> socket:[39221]
    lrwx------. 1 root root 64 Apr 18 16:51 17 -> socket:[41080]
    lrwx------. 1 root root 64 Apr 18 16:51 18 -> socket:[40014]
    lrwx------. 1 root root 64 Apr 18 16:51 19 -> socket:[34617]
    lrwx------. 1 root root 64 Apr 18 16:51 20 -> socket:[34620]
    lrwx------. 1 root root 64 Apr 18 16:51 23 -> socket:[42357]
    lr-x------. 1 root root 64 Apr 18 16:51 3 -> /dev/urandom
    lrwx------. 1 root root 64 Apr 18 16:51 4 -> socket:[37468]
    lrwx------. 1 root root 64 Apr 18 16:51 5 -> socket:[37471]
    lrwx------. 1 root root 64 Apr 18 16:51 6 -> socket:[289532]
    lrwx------. 1 root root 64 Apr 18 16:51 7 -> socket:[31728]
    lrwx------. 1 root root 64 Apr 18 16:51 8 -> socket:[37450]
    lrwx------. 1 root root 64 Apr 18 16:51 9 -> socket:[37451]
    l-wx------. 1 root root 64 Apr 13 16:35 2 -> /root/.vnc/localhost.localdomain:1.log
目录中的每一项都是一个符号链接，指向打开的文件，数字则代表文件描述符。  

## latency

`/proc/[pid]/latency`显示哪些代码造成的延时比较大（使用这个`feature`，需要执行“`echo 1 > /proc/sys/kernel/latencytop`”）。举例如下：  

    # cat /proc/2948/latency
    Latency Top version : v0.1
    30667 10650491 4891 poll_schedule_timeout do_sys_poll SyS_poll system_call_fastpath 0x7f636573dc1d
    8 105 44 futex_wait_queue_me futex_wait do_futex SyS_futex system_call_fastpath 0x7f6365a167bc
每一行前三个数字分别是后面代码执行的次数，总共执行延迟时间（单位是微秒）和最长执行延迟时间（单位是微秒），后面则是代码完整的调用栈。

## limits

`/proc/[pid]/limits`显示当前进程的资源限制。举例如下：  

    # cat /proc/2948/limits
    Limit                     Soft Limit           Hard Limit           Units
    Max cpu time              unlimited            unlimited            seconds
    Max file size             unlimited            unlimited            bytes
    Max data size             unlimited            unlimited            bytes
    Max stack size            8388608              unlimited            bytes
    Max core file size        0                    unlimited            bytes
    Max resident set          unlimited            unlimited            bytes
    Max processes             6409                 6409                 processes
    Max open files            1024                 4096                 files
    Max locked memory         65536                65536                bytes
    Max address space         unlimited            unlimited            bytes
    Max file locks            unlimited            unlimited            locks
    Max pending signals       6409                 6409                 signals
    Max msgqueue size         819200               819200               bytes
    Max nice priority         0                    0
    Max realtime priority     0                    0
    Max realtime timeout      unlimited            unlimited            us
`Soft Limit`表示`kernel`设置给资源的值，`Hard Limit`表示`Soft Limit`的上限，而`Units`则为计量单元。

## maps

`/proc/[pid]/maps`显示进程的内存区域映射信息。举例如下：  

    # cat /proc/2948/maps
    ......
    address                   perms offset  dev   inode                      pathname
    7f4a2e2ad000-7f4a2e2ae000 rw-p 00006000 08:14 6505977                    /usr/lib64/sasl2/libsasldb.so.3.0.0
    7f4a2e2ae000-7f4a2e2af000 ---p 00000000 00:00 0
    7f4a2e2af000-7f4a2eaaf000 rw-p 00000000 00:00 0                          [stack:94671]
    7f4a2eaaf000-7f4a2eab0000 ---p 00000000 00:00 0
    7f4a2eab0000-7f4a2f2b0000 rw-p 00000000 00:00 0                          [stack:94670]
    ......
    7f4a434d0000-7f4a434d5000 rw-p 0006e000 08:14 4292988                    /usr/sbin/libvirtd
    7f4a4520a000-7f4a452f7000 rw-p 00000000 00:00 0                          [heap]
    7ffd1a7e4000-7ffd1a805000 rw-p 00000000 00:00 0                          [stack]
    7ffd1a820000-7ffd1a821000 r-xp 00000000 00:00 0                          [vdso]
    ffffffffff600000-ffffffffff601000 r-xp 00000000 00:00 0                  [vsyscall]

其中注意的一点是`[stack:<tid>]`是线程的堆栈信息，对应于`/proc/[pid]/task/[tid]/`路径。  

## root

`/proc/[pid]/root`是进程根目录的符号链接。举例如下： 

    # ls -lt /proc/2948/root
    lrwxrwxrwx 1 root root 0 Nov  9 12:14 /proc/2948/root -> /

## stack

`/proc/[pid]/stack`显示当前进程的内核调用栈信息，只有内核编译时打开了`CONFIG_STACKTRACE`编译选项，才会生成这个文件。举例如下：  

    # cat /proc/2948/stack
    [<ffffffff80168375>] poll_schedule_timeout+0x45/0x60
    [<ffffffff8016994d>] do_sys_poll+0x49d/0x550
    [<ffffffff80169abd>] SyS_poll+0x5d/0xf0
    [<ffffffff804c16e7>] system_call_fastpath+0x16/0x1b
    [<00007f4a41ff2c1d>] 0x7f4a41ff2c1d
    [<ffffffffffffffff>] 0xffffffffffffffff

## statm

`/proc/[pid]/statm`显示进程所占用内存大小的统计信息，包含七个值，度量单位是`page`（`page`大小可通过`getconf PAGESIZE`得到）。举例如下：  

    # cat /proc/2948/statm  
    72362 12945 4876 569 0 24665 0

各个值含义：   
    a）进程占用的总的内存；  
    b）进程当前时刻占用的物理内存；   
    c）同其它进程共享的内存；  
    d）进程的代码段；  
    e）共享库（从`2.6`版本起，这个值为`0`）；  
    f）进程的堆栈；  
    g）`dirty pages`（从`2.6`版本起，这个值为`0`）。  
    
## syscall

`/proc/[pid]/syscall`显示当前进程正在执行的系统调用。举例如下：  

    # cat /proc/2948/syscall
    7 0x7f4a452cbe70 0xb 0x1388 0xffffffffffdff000 0x7f4a4274a750 0x0 0x7ffd1a8033f0 0x7f4a41ff2c1d
    
第一个值是系统调用号（`7`代表`poll`），后面跟着`6`个系统调用的参数值（位于寄存器中），最后两个值依次是堆栈指针和指令计数器的值。如果当前进程虽然阻塞，但阻塞函数并不是系统调用，则系统调用号的值为`-1`，后面只有堆栈指针和指令计数器的值。如果进程没有阻塞，则这个文件只有一个“`running`”的字符串。

内核编译时打开了`CONFIG_HAVE_ARCH_TRACEHOOK`编译选项，才会生成这个文件。  

## wchan

`/proc/[pid]/wchan`显示当进程`sleep`时，`kernel`当前运行的函数。举例如下：  

    # cat /proc/2948/wchan
    kauditd_thread
