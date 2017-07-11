# CPU Utilization 和 CPU Load Average

## Total CPU Utilization

> CPU 利用率的正确翻译为 CPU Utilization 不是 CPU Usage ；

cpu utilization 即 cpu 利用率；

cpu 利用率的计算是基于 `/proc/stat` 文件中的内容得到的（通过两次抽样值进行计算，抽样间隔自选）；

cpu 利用率可以基于如下公式进行计算（取两个采样点，然后基于差值计算）：

```
cpu_utilization = [(user_2 + sys_2 + nice_2) - (user_1 + sys_1 + nice_1)] / (total_2 - total_1) * 100;
```

> 注意：上述计算方式其实存在问题（因为 `total = user + nice + system + idle + ioWait + irq + softIRQ + steal + guest + guestNice`）；在实际使用过程中，最简单且正确的计算方式是通过 idle 指标进行计算得到，即
> 
> ```
> %cpu_utilization = 1 - %idle
> %cpu_utilization = 1 - (idle2 - idle1) / (total_2 - total_1)
> ```

下面是一段网友给出的基于 Bash 采集 cpu 利用率的代码：

```
#!/bin/sh
##echo user nice system idle iowait irq softirq
CPULOG_1=$(cat /proc/stat | grep 'cpu ' | awk '{print $2" "$3" "$4" "$5" "$6" "$7" "$8}')
SYS_IDLE_1=$(echo $CPULOG_1 | awk '{print $4}')
Total_1=$(echo $CPULOG_1 | awk '{print $1+$2+$3+$4+$5+$6+$7}')

sleep 5

CPULOG_2=$(cat /proc/stat | grep 'cpu ' | awk '{print $2" "$3" "$4" "$5" "$6" "$7" "$8}')
SYS_IDLE_2=$(echo $CPULOG_2 | awk '{print $4}')
Total_2=$(echo $CPULOG_2 | awk '{print $1+$2+$3+$4+$5+$6+$7}') 

SYS_IDLE=`expr $SYS_IDLE_2 - $SYS_IDLE_1`

Total=`expr $Total_2 - $Total_1`
SYS_USAGE=`expr $SYS_IDLE/$Total*100 |bc -l`  ## 这里的变量名有问题，此处得到的应该是 %idle * 100

SYS_Rate=`expr 100-$SYS_USAGE |bc -l`         ## 此处得到的是 %cpu_utilization * 100

Disp_SYS_Rate=`expr "scale=3; $SYS_Rate/1" |bc`  ## 取小数点后 3 位
echo $Disp_SYS_Rate%
```

## per-process CPU Utilization

基于 `/proc/<pid>/stat` 文件中的内容计算；

指定进程占用的 CPU 时间（该值包括其包含的所有线程占用的 cpu 时间）计算公式如下：

```
processCpuTime = utime + stime + cutime + cstime
```

具体例子：

```
$cat /proc/10996/stat
10996 (docker) S 10970 10996 10970 34845 10996 1077960704 3124 0 1 0 18 27 0 0 20 0 7 0 252734195 222781440 2866 18446744073709551615 4194304 33785809 140732225023424 140732225022864 4652003 0 1945795109 0 2143420159 18446744073709551615 0 0 17 1 0 0 0 0 0 35882968 36246560 44216320 140732225031681 140732225031714 140732225031714 140732225032162 0
[#49#root@dockermonitor ~]$
```

在 man proc 中可以看到如下信息：

```
/proc/[pid]/stat

Status information about the process.
```

> **`utime %lu (14)`**
Amount of time that this process has been scheduled in **user mode**, measured in clock ticks (divide by `sysconf(_SC_CLK_TCK)`). This includes guest time, `guest_time` (time spent running a virtual CPU, see below), so that applications that are not aware of the guest time field do not lose that time from their calculations.

该任务在**用户态**运行的时间，单位为 jiffies ；

> **`stime %lu (15)`**
Amount of time that this process has been scheduled in **kernel mode**, measured in clock ticks (divide by `sysconf(_SC_CLK_TCK)`).

该任务在**核心态**运行的时间，单位为 jiffies ；

> **`cutime %ld (16)`**
Amount of time that this process's **waited-for children** have been scheduled in **user mode**, measured in clock ticks (divide by `sysconf(_SC_CLK_TCK)`). (See also `times(2)`.) This includes guest time, `cguest_time` (time spent running a virtual CPU, see below).

该任务等待其子进程（线程）在**用户态**运行的时间，单位为 jiffies ；

> **`cstime %ld (17)`**
Amount of time that this process's **waited-for children** have been scheduled in **kernel mode**, measured in clock ticks (divide by `sysconf(_SC_CLK_TCK)`).

该任务等待其子进程（线程）在**核心态**运行的时间，单位为 jiffies ；

## per-thread CPU Utilization

基于 `/proc/<pid>/task/<tid>/stat` 文件计算；

线程耗费的 CPU 时间：

```
threadCpuTime = utime + stime
```

## CPU Load Average

cpu load 即 cpu 负载；一般在讨论时会描述成 CPU load average ，即 CPU 平均负载，这样更为准确；

**系统 CPU 平均负载**被定义为在特定时间间隔内运行队列（run queue）中的平均进程数（应该说成可调度实体更为准确）。如果一个进程（或线程）满 足以下条件则其就会位于运行队列中：

- **没有在等待 I/O 操作的结果**；
- **没有主动进入等待状态（即没有调用 `wait`）**；
- **没有被停止（例如等待终止）**；

cpu 平均负载是从 `/proc/loadavg` 中读取的，该值为所有 cpu 的总体值，计算单个 cpu 的负载平均值时，需要除 cpu 的数量；

常规情况

```
vagrant@vagrant-ubuntu-trusty:~$ cat /proc/loadavg
0.00 0.01 0.05 1/81 1954
```

每个值的含义依次为：

- **lavg_1** (0.00)：1 分钟平均负载；
- **lavg_5** (0.01)：5 分钟平均负载；
- **lavg_15** (0.05)：15 分钟平均负载；
- **nr_running** (1)：在采样时刻，可运行的任务数目，与 `/proc/stat` 的 `procs_running` 表示相同意思；
- **nr_threads** (81)：在采样时刻，系统中存在的任务数目（不包括运行已经结束的任务）；
- **last_pid** (1954)：最近被创建出来的进程 pid 值，包括轻量级进程，即线程；


> The first three fields in this file are load average figures giving the number of jobs in the `run queue` (**state R**) or waiting for disk I/O (**state D**) averaged over 1, 5, and 15 minutes. They are the same as the load average numbers given by `uptime(1)` and other programs. 
> 
> The fourth field consists of two numbers separated by a slash (/). The first of these is the number of currently **runnable** kernel scheduling entities (**processes**, **threads**). The value after the slash is the number of kernel scheduling entities that currently **exist** on the system. 
> 
> The fifth field is the PID of the process that was most recently created on the system.


之前遇到的一台由于内核 bug 而导致计数值溢出的情况：

```
[root@xg-pcd-commodity-service-8 ~]# cat /proc/loadavg
4294967293.21 4294967293.36 4294967293.39 4294967294/740 4272
```

## Q&A

- **如何判断系统是否已经 overload**

> “**有多少核心即为有多少负荷**”法则： 在多核处理中，你的系统均值不应该高于处理器核心的总数量。

对一般的系统来说，可以根据 CPU 数量去判断。如果平均负载始终在 1.2 以下，而你的机器有 2 颗 CPU ，那么基本不会出现 CPU 不够用的情况。一般结论：**Load average 平均要小于 CPU 的数量**；

- **低 CPU Utilization 的情况下是否会有高 Load Average 情况产生**

首先需要理解**占有时间 (occupy)**和**使用时间 (active use)**的区别：可以简单的认为使用时间为 `total - idle` 得到的时间；而占用时间为 `total` ；

当分配时间片以后，是否使用完全取决于使用者，因此完全可能出现低 CPU 利用率、高 Load Average 的情况。由此来看，**仅仅从 CPU 利用率来判断 CPU 是否处于一种超负荷的工作状态还是不够的**，必须结合 load average 来全局的看 CPU 的使用情况和申请情况。

- **Load average 与容量规划（Capacity Planning）**

一般是会根据 15 分钟的 load average 为基准进行考虑。

如果 1 分钟平均值出现大于 **CPU 个数 * 核数**的情况，还不用担心；   
如果 5 分钟平均也是，那就要警惕了；   
如果 15 分钟平均也是这样，就要分析哪里出问题了，防范于未然；   


- **针对 Load average 的常见误解**

> 系统 load average 高一定是性能有问题；

真相：Load average 高也许是因为在进行 CPU 密集型的计算；

> 系统 Load average 高一定是 CPU 性能问题或数量不够；

真相：Load average 高只是代表运行队列中累积了过多的调度实体（或称为任务）。但队列中的任务实体可能是耗 CPU 的，也可能是耗 I/O 或者其它因素的；因此，不能单纯认为 CPU 性能不足或数量不够；

> 系统长期 Load average 高，首先增加 CPU ；

真相：Load average 只是表象，不是实质。增加 CPU 个别情况下会临时看到 Load average 下降，但治标不治本；原因同上；


- 总 CPU 时间如何计算？要不要算 steal 和 guest 的时间，以及 guest_nice 时间？

```
totalCpuTime = user + nice + system + idle + iowait + irq + softirq + steal + guest + guest_nice
```

要，在纯粹的一台物理机上（即其上未跑其它 guest OS ，自身也未作为 guest OS 被虚拟机调度器管理），**steal**/**guest**/**guest_nice** 值应该都为 0 ；除此之外，上述值就应该不为 0 ；


----------


## 获取 CPU 时间占用情况的各种命令

### top 输出

```
%Cpu(s):  0.0 us,  0.0 sy,  0.0 ni, 99.7 id,  0.3 wa,  0.0 hi,  0.0 si,  0.0 st
```

### sar 输出

Collect, report, or save system activity information.

```
$sar -u ALL 10 3
Linux 3.10.0-229.11.1.el7.x86_64 (dockermonitor) 06/29/2016 _x86_64_ (2 CPU)

04:03:25 PM CPU %usr %nice %sys %iowait %steal %irq %soft %guest %gnice   %idle
04:03:35 PM    all 5.07   0.00   1.93       0.30     0.15 0.00  0.05    0.00     0.00   92.50
04:03:45 PM    all 2.53   0.00   1.47       0.35     0.10 0.00  0.00    0.00     0.00   95.55
04:03:55 PM    all 2.57   0.00   1.46       0.35     0.05 0.00  0.05    0.00    0.00   95.51
Average:          all 3.39   0.00   1.62       0.34     0.10 0.00   0.03    0.00    0.00   94.52
[#32#root@dockermonitor ~]$
```

> 注意：上述指标中 `%steal` 不为 0 ，说明当前 OS 是在虚拟机调度器的管理下运行的，且存在其它 OS 也被虚拟机调度器管理；Stolen time, which is the time spent in other operating systems when running in a virtualized environment.

sar 参数说明：

```
-u [ ALL ]

**Report CPU utilization**. The `ALL` keyword indicates that all the CPU fields should be displayed. The report may show the following fields:

%user
Percentage of CPU utilization that occurred while executing at the user level (application).
Note that this field includes time spent running virtual processors.

%usr
Percentage of CPU utilization that occurred while executing at the user level (application).
Note that this field does NOT include time spent running virtual processors.

%nice
Percentage of CPU utilization that occurred while executing at the user level with nice priority.

%system
Percentage of CPU utilization that occurred while executing at the system level (kernel).
Note that this field includes time spent servicing hardware and software interrupts.

%sys
Percentage of CPU utilization that occurred while executing at the system level (kernel).
Note that this field does NOT include time spent servicing hardware or software interrupts.

%iowait
Percentage of time that the CPU or CPUs were idle during which the system had an outstanding disk I/O request.

%steal
Percentage of time spent in involuntary wait by the virtual CPU or CPUs while the hypervisor was servicing another virtual processor.

%irq
Percentage of time spent by the CPU or CPUs to service hardware interrupts.

%soft
Percentage of time spent by the CPU or CPUs to service software interrupts.

%guest
Percentage of time spent by the CPU or CPUs to run a virtual processor.

%gnice
Percentage of time spent by the CPU or CPUs to run a niced guest.

%idle
Percentage of time that the CPU or CPUs were idle and the system did not have an outstanding disk I/O request.

Note: On SMP machines a processor that does not have any activity at all (0.00 for every field) is a disabled (offline) processor.
```

### mpstat 输出

Report processors related statistics.

```
$mpstat -P ALL
Linux 3.10.0-229.11.1.el7.x86_64 (dockermonitor) 06/29/2016 _x86_64_ (2 CPU)

07:20:32 PM CPU %usr %nice %sys %iowait %irq %soft %steal %guest %gnice %idle
07:20:32 PM    all   1.95  0.00    1.19      0.27  0.00  0.08   0.08     0.00     0.00  96.43
07:20:32 PM     0   1.99  0.00    1.25     0.24  0.00  0.08   0.08     0.00     0.00  96.36
07:20:32 PM      1   1.91   0.00    1.14      0.30  0.00  0.08   0.08    0.00     0.00  96.49
[#53#root@dockermonitor ~]$
```

mpstat 中输出的 CPU 使用率报告：

```
-P { cpu [,...] | ON | ALL }

Indicate the processor number for which statistics are to be reported.  cpu is the processor number. Note that processor 0 is the first processor.  The ON keyword indicates that statistics are to be reported for every online processor, whereas the ALL keyword indicates that statistics are to be reported for all processors.


-u

Report CPU utilization. The following values are displayed:

CPU
Processor number. The keyword all indicates that statistics are calculated as averages among all processors.

%usr
Show the percentage of CPU utilization that occurred while executing at the user level (application).

%nice
Show the percentage of CPU utilization that occurred while executing at the user level with nice priority.

%sys
Show the percentage of CPU utilization that occurred while executing at the system level (kernel). Note that this does not include time spent servicing hardware and software interrupts.

%iowait
Show the percentage of time that the CPU or CPUs were idle during which the system had an outstanding disk I/O request.

%irq
Show the percentage of time spent by the CPU or CPUs to service hardware interrupts.

%soft
Show the percentage of time spent by the CPU or CPUs to service software interrupts.

%steal
Show the percentage of time spent in involuntary wait by the virtual CPU or CPUs while the hypervisor was servicing another virtual processor.

%guest
Show the percentage of time spent by the CPU or CPUs to run a virtual processor.

%gnice
Show the percentage of time spent by the CPU or CPUs to run a niced guest.

%idle
Show the percentage of time that the CPU or CPUs were idle and the system did not have an outstanding disk I/O request.

Note: On SMP machines a processor that does not have any activity at all is a disabled (offline) processor.
```


### iostat 输出

Report Central Processing Unit (CPU) statistics and input/output statistics for devices and partitions.

```
$iostat
Linux 3.10.0-229.11.1.el7.x86_64 (dockermonitor) 06/29/2016 _x86_64_ (2 CPU)

avg-cpu: %user %nice %system %iowait %steal    %idle
                   1.95   0.00        1.28       0.27    0.08   96.43

Device: tps kB_read/s kB_wrtn/s kB_read kB_wrtn
vda 2.66 1.37 11.09 4190325 33851948
vdb 43.63 4.49 567.16 13713889 1731436536
vdc 0.06 0.21 0.45 629764 1377676
dm-0 21.25 4.11 470.90 12544414 1437579039
dm-4 0.01 0.17 0.06 521344 190041
dm-5 0.07 0.66 1.04 2022838 3160694
dm-6 0.00 0.01 0.00 27485 3210
dm-7 0.00 0.00 0.00 9686 2464
dm-8 0.00 0.00 0.00 9467 2784
dm-9 0.01 0.08 0.27 243767 837085
dm-10 0.45 0.41 10.15 1262443 30974262
dm-1 2.46 0.32 81.93 971740 250114665
dm-2 2.46 0.32 81.93 968732 250115322
dm-3 2.46 0.29 81.93 878884 250116748
dm-11 1.79 0.37 50.01 1142915 152658304
dm-12 0.00 0.02 0.00 54430 2361

[#55#root@dockermonitor ~]$
```

iostat 中输出的 CPU 使用率报告：

```
CPU Utilization Report

The first report generated by the iostat command is the CPU Utilization Report. For multiprocessor systems, the CPU values are global averages among all processors.  The report has the following format:

%user
Show the percentage of CPU utilization that occurred while executing at the user level (application).

%nice
Show the percentage of CPU utilization that occurred while executing at the user level with nice priority.

%system
Show the percentage of CPU utilization that occurred while executing at the system level (kernel).

%iowait
Show  the  percentage of time that the CPU or CPUs were idle during which the system had an outstanding disk I/O request.

%steal
Show the percentage of time spent in involuntary wait by the virtual CPU or CPUs while the hypervisor was servicing another virtual processor.

%idle
Show the percentage of time that the CPU or CPUs were idle and the system did not have an outstanding disk I/O request.
```


### /proc/stat 中的内容

```
~$ cat /proc/stat
cpu 851 0 1784 956996 318 176 0 0 0 0
cpu0 851 0 1784 956996 318 176 0 0 0 0
intr 142546 58 10 0 0 0 0 0 0 0 0 0 0 156 0 9444 0 4977 0 0 10027 3894 7365 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 878495
btime 1467168528
processes 2032
procs_running 1
procs_blocked 0
softirq 102927 0 65552 3845 15021 11176 0 3 0 92 7238
vagrant@vagrant-ubuntu-trusty:~$
```

> The very first "cpu" line aggregates the numbers in all of the other "cpuN" lines.

第一行的数值表示的是 CPU 总的使用情况，所以计算整体 CPU 使用率只需要使用这一行数据即可；

> These numbers identify the amount of time the CPU has spent performing different kinds of work. Time units are in `USER_HZ` or `Jiffies` (typically hundredths of a second).

`jiffies` 是内核中的一个全局变量，用来记录自系统启动以来产生的节拍数；在 linux 中，一个节拍大致可理解为操作系统进程调度的最小时间片，不同 linux 内核可能值有不同，通常在 1ms 到 10ms 之间；

The meanings of the columns are as follows, from left to right:

- **user**(1): normal processes executing in user mode. 从系统启动开始，累计到当前时刻，处于**用户态**的运行时间，不包含 nice 值为负进程；
- **nice**(2): niced processes executing in user mode. 从系统启动开始，累计到当前时刻，nice 值为负的进程所占用的 CPU 时间；
- **system**(3): processes executing in kernel mode. 从系统启动开始，累计到当前时刻，处于**核心态**的进程的运行时间；
- **idle**(4): twiddling thumbs. 从系统启动开始，累计到当前时刻，除 IO 等待时间以外的其它等待时间；
- **iowait**(5): waiting for I/O to complete. 从系统启动开始，累计到当前时刻，IO 等待时间 (since 2.5.41)；
- **irq**(6): servicing interrupts. 从系统启动开始，累计到当前时刻，硬中断时间 (since 2.6.0-test4)；
- **softirq**(7): servicing softirqs. 从系统启动开始，累计到当前时刻，软中断时间 (since 2.6.0-test4)；
- **steal**(8) (since Linux 2.6.11): Stolen time, which is the time spent in other operating systems when running in a virtualized environment. 当运行在虚拟化环境中，花费在其它 OS 中的时间（基于虚拟机监视器 hypervisor 的调度）；可以理解成由于虚拟机调度器将 cpu 时间用于其它 OS 了，故当前 OS 无法使用 CPU 的时间；
- **guest**(9) (since Linux 2.6.24): Time spent running a virtual CPU for guest operating systems under the control of the Linux kernel. 花费给当前 host 上运行的其它 guest OS 的时间；
- **guest_nice**(10) (since Linux 2.6.33): Time spent running a niced guest (virtual CPU for guest operating systems under the control of the Linux kernel).

The "**intr**" line gives counts of interrupts serviced since boot time, for each
of the possible system interrupts. The first column is the total of all interrupts serviced; each subsequent column is the total for that particular interrupt.

The "**ctxt**" line gives the total number of context switches across all CPUs

The "**btime**" line gives the time at which the system booted, in seconds since
the Unix epoch.

The "**processes**" line gives the number of **processes** and **threads** created, which includes (but is not limited to) those created by calls to the `fork()` and `clone()` system calls.

The "**procs_running**" line gives the number of processes currently **running** on CPUs.

The "**procs_blocked**" line gives the number of processes currently **blocked**, waiting for I/O to complete.


----------



参考:

- [CPU利用率和Load Average的区别](http://www.voidcn.com/blog/chenhaotong/article/p-5996294.html)
- [Understanding Linux CPU Load - when should you be worried?](http://blog.scoutapp.com/articles/2009/07/31/understanding-load-averages)
- [理解 Linux 的处理器负载均值（翻译）](https://www.gracecode.com/posts/2973.html)

