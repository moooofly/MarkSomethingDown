

# Inside the Erlang VM - with focus on SMP 

----------

## Introduction

在 Erlang 中支持 SMP (对称多处理器) 的历史可以追溯到 1997-1998 左右，自 Pekka Hedqvist 的硕士论文开始，其导师为 Tony Rogvall (Ericsson Computer Science Lab)；

The implementation was run on a Compaq with 4 Pentium Pro 200 Mhz CPU’s
(an impressive machine in those days) and showed a great potential for scalability
with additional processors but suffered from bad IO performance.

The work with SMP did not continue at that time since it was so easy to increase
performance by just upgrading the HW to the newest processor. There simply
was no business case for it at the time.

The SMP work was restarted at 2005 and now as part of the ordinary
development. The work was driven by the Erlang development team at Ericsson
with participation and contributions from Tony Rogvall (then at Synapse) and the
HiPE group at Uppsala University. 

策略如下（现在仍然是这个策略）：
- 首先，”make it work”
- 其次，”measure” 并找到瓶颈点
- 最后，通过移除瓶颈点进行 ”optimize”

带有 SMP 支持的第一个稳定运行时 release 版本为 2006 年  5 月发布的 OTP R11B ；

This ended the first cycle of the strategy and a new iteration with “measure”,
“optimize” and “make it work” started. Read more about it in the next pages. 

## How it works

### Erlang VM with no SMP support

不带 SMP 支持的 Erlang VM 只会在主线程中运行一个 scheduler ；scheduler 从 run queue 中选取可运行的 Erlang 进程和 IO 任务进行执行，并且不需要锁定任何数据结构，因为只有一个线程进行数据访问；

![Erlang (non SMP) VM today](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Erlang (non SMP) VM today.png  "Erlang (non SMP) VM today")

### Erlang VM with SMP support (in R11B and R12B)

带 SMP 支持的 Erlang VM 能够启动 1 - 1024 个 scheduler ，每一个 scheduler 都运行于一个线程之中；

全部 scheduler 都会从同一个 common run queue 中选取可运行 Erlang 进程和 IO 任务；在支持 SMP 的 VM 中，所有共享数据结构都会被锁保护，而 run queue 是通过锁保护共享数据结构的其中一个例子；

![Erlang SMP VM today](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Erlang SMP VM today.png "Erlang SMP VM today")

#### First release for use in Products, March 2007

Measurements from a real telecom product showed a 1.7 speed improvement between a single and a dual core system.

It should be noted that it took only about a week to port the telecom system to a new OTP release with SMP support, to a new Linux distribution and to a new incompatible CPU architecture, the Erlang code was not even recompiled.

It took a little longer to get the telecom system in product status, a few minor changes was needed in the Erlang code because Erlang processes now can run truly parallel which changes the timing and ordering of events which the old application code did not count for.

The performance improvements achieved on a dual core processor for a real telecom system where encouraging and after that several other telecom systems have also taken benefit from the SMP support in Erlang.


#### SMP in R12B

从 OTP R12B 开始，如果操作系统发现自身具有超过 1 个 CPU（或核心），则 VM 默认就会启动SMP ，并创建与 CPU 或核心数目相同的 scheduler ；

你可以在 erl 命令的第一行输出中看到如下信息

例如
```shell
Erlang (BEAM) emulator version 5.6.4 [source] [smp:4] .....
```
上面的 `[smp:4]` 表明正在运行支持 SMP 的 VM ，并且启动了 4 个  scheduler ；

默认行为可以通过 "-smp [enable|disable|auto]" 进行覆盖；`auto` 为默认值；
若想设置启动 scheduler 的具体数量，需要设置 -smp 为 enable 或 auto ，并使用 "+S Number" 选项，其中 Number 值为 scheduler 的数量（1..1024）； 

> ⚠️ 运行超过 CPU 或 CPU 核数的 scheduler 通常不会有任何额外的收益；

> ⚠️ 在一些操作系统上，单个进程可以使用的 CPU 或核心数量可以通过命令进行限制；例如，在Linux 上，命令 "taskset" 就可用于此目的；Erlang VM 当前只能检测到可用 CPU 或核心数量，而不会将 "taskset" 设置的 mask 值考虑在内；

基于上述原因，可能会发生，并且实际已经发生了诸如“尽管 Erlang VM 运行了 4 个 scheduler ，但只有 2 个核心被使用“的情况；这是由于操作系统自身采取的限制导致的，因为其将 "taskset" 设置的 mask 考虑在内了；

在 Erlang VM 中，每一个 scheduler 都运行在一个操作系统线程中，并由操作系统自行决定这些线程是否在不同的核上被执行；通常情况下，操作系统的默认处理方式就很好，会令线程在执行过程中一直跑在同一个核心上；

Erlang 进程在不同时段内会被不同的 scheduler 所运行，因为只要某个 scheduler 空闲，其就会从同一个 common run-queue 中提取 Erlang 进程或 IO 任务进行调度；


## Performance and scalability

只启动一个 scheduler 的 SMP VM  的性能要略低于（10%）non SMP VM ；

这是因为 SMP VM 需要针对所有共享数据结构使用锁；但是只要没有锁冲突问题，由其导致的额外开销就不会很高（只有锁冲突才会花费大量时间）

This explains why it in some cases can be more efficient to run several SMP VM's
with one scheduler each instead on one SMP VM with several schedulers. Of course
the running of several VM's require that the application can run in many parallel tasks
which has no or very little communication with each other.

If a program scale well with the SMP VM over many cores depends very much on the
characteristics of the program, some programs scale linearly up to 8 and even 16
cores while other programs barely scale at all even on 2 cores.

This might sound bad, but in practice many real programs scale well on the number
of cores that are common on the market today, see below.

Real telecom products supporting a massive number if simultaneously ongoing
"calls" represented as one or several Erlang processes per core have shown very
good scalability on dual and quad core processors.

Note, that these products was written in the normal Erlang style long before the SMP
VM and multi core processors where available and they could benefit from the Erlang
SMP VM without changes and even without need to recompile the code. 


## Our strategy with SMP

早在最开始实现 SMP VM 的时候，我们就定下了如下策略：

```shell
"First make it work, then measure, then optimize".
```

We are still following this strategy consistently since the first stable working SMP VM
that we released in May 2006 (R11B).

Another important part of the strategy is to hide the problems and awareness of SMP
execution for the Erlang programmer. Erlang programs should be written as usual
using processes for parallel tasks, the utilization of CPUs and cores should be
handled by the Erlang VM. It must be easy and cost effective to utilize multicore and
SMP HW with Erlang this is one of our absolute strengths compared to other
programming languages.

There will be new BIF’s for SMP related stuff but we try to avoid that as much as
possible. We think it is preferable to configure SMP related things such as number of
cores to use, which cores to use on the OS level and as parameters to the Erlang
VM at startup.

The principle is that an Erlang program should run perfectly well on any system no
matter what number of cores or processors there are. 

## Next steps with SMP and Erlang

There are more known things to improve and we address them one by one taking the
one we think gives most performance per implementation effort first and so on.

We are now putting most focus on getting consistent better scaling on many cores
(more than 4).

The SMP implementation is continually improved in order to get better performance
and scalability. In each service release R12B-1, 2, 3, 4, 5 , ..., R13B-0, 1, …, R14B
etc. you will find new optimizations. 


### Some known bottlenecks

下面列出一些我们知道的、最重要的瓶颈点；可以肯定的是，还有更多的瓶颈存在着，等待我们一个一个去解决；在移除一个瓶颈后，马上会有新的瓶颈点出现，还可能导致已知的其他点的重要性发生变更；

#### The common run-queue

当 CPU 或 CPU 核数增多时，单独一个 common run-queue 将会成为主要瓶颈；

从 4 核开始该问题就会显现出来，但对于许多应用来说，4 核情况下仍能给出不错的性能表现；

我们正在实现每个 scheduler 一个 run-queue 的解决方案，并将此作为当前最重要的改进点；本文的后续内容会有说明；

#### Ets tables

Ets tables involves locking. Before R12B-4 there was 2 locks involved in every
access to an ets-table, but in R12B-4 the locking of the meta-table is optimized to
reduce the conflicts significantly (as mentioned earlier it is the conflicts that are
expensive).

If many Erlang processes access the same table there will be a lot of lock conflicts
causing bad performance especially if these processes spend a majority of their work
accessing ets-tables.

The locking is on table-level not on record level. An obvious solution is to introduce
more fine granular locking. 

Note! that this will have impact on Mnesia as well since Mnesia is a heavy user of
ets-tables. 

#### Message passing

当许多进程同时发送消息到相同的接收进程时，将会存在大量的锁冲突；可以通过减少需要访问锁的待处理工作量优化该问题；

#### A process can block the scheduler

一旦某个进程在阻塞等待获取访问某个 ets 表的锁，整个 scheduler 将会被阻塞住，什么也做不了；直到锁被成功获取后，进程才会继续执行；

上述情况可以通过引入“进程级别锁“进行改进，即如果某个进程在阻塞等待获取锁，则会被调度出 scheduler ，之后 scheduler 会从 run queue 中提取下一个进程进行调度；

我们已经实现并测量了这种解决方案，结论是该方案可以在独立 run queue 可用时被引入；我们仍旧需要确认该方案在某些特殊情况下是否会导致性能下降；

### Separate run-queues per scheduler

针对 SMP 的下一个 Erlang 运行时系统的重大性能改进，就是将所有 scheduler 共享同一个 run queue 变更为每一个 scheduler 使用一个独立的 run queue ；该变化会极大的减少多核或多处理器系统中锁冲突的数量；从 4 核 开始，性能改进的效果已经体现在许多应用中了，并且在具有 8, 16 或者更多核的系统中，将会有更佳出色的表现；

![Erlang SMP VM next step](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Erlang SMP VM next step.png "Erlang SMP VM next step")

#### Migration logic

当每个 scheduler 都具有独立的 run queue 时，问题将从访问同一个 run queue 时的锁冲突，变成了迁移逻辑的实现效率和公平性问题；

目前我们已有的实现还需要大量的 benchmark 测试和精细调优，才能保证其工作在最优状态；粗略的说，其工作方式如下：

The maximum number of runable processes over all schedulers is measured
approximately 4 times per second. This value divided by number of schedulers is
then used to trigger migration of processes from one scheduler to another scheduler.

When a scheduler is about to schedule in a new process it will first check if its
number of runable processes is above the max value described above and if it is it
will migrate the process to another scheduler according to the migration path set up. 

There are also 2 other occasions in addition to the “schedule in” of a new process
when a process migration can occur:
1. If a scheduler gets out of jobs it will steal jobs from other schedulers.
2. Underloaded schedulers will also steal jobs from heavily overloaded schedulers in their migration paths.

Below follows some measurements that show early indications of the improvements
the system with separate run-queues per scheduler and the migration logic described
above will give.

The graph below shows the results from running the benchmark “big bang” with 1, 2,
4, 8 schedulers on both the current system with one single run-queue and on the
next to come system with multiple run-queues one per scheduler.

The benchmark spawns 1000 processes which all sends a ‘ping’ message to all
other processes and answer with a ‘pong’ message for all ‘ping’ it receives.

The “fat” lines in the graph shows the multiple run-queue case and as can be seen
the improvement is significant. 

![Number of schedulers](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Number of schedulers.png "Number of schedulers")

## Frequently Asked Questions

### Is there any difference in the .beam file depending on if it should run in a SMP or non SMP system?

只要目标模块没有在使能了 HiPE 的系统中启动 “native” 选项进行编译，.beam 文件就是相同的，并且可以在 SMP 和 non-SMP 系统中通用；

### Can an Erlang process be locked to a specific processor core?

程序员无法将 Erlang 进程锁定到特定的处理器上执行，并且这种实现是有意为之；在未来的  release 版本中，可能允许将某个 scheduler 锁定到某个特定的 core 上执行； 

### What is the relation between asynch threads and SMP?

异步线程池和 SMP 没有一毛钱关系；异步线程池只用于文件驱动器和用户自己实现的（使用该池）驱动器；文件驱动器使用该池可以在遇到大文件操作时避免整个 Erlang VM 被长时间锁住；对于 VM 来说，异步线程池的引入要远远早于 SMP ；并且异步线程池在  non SMP VM 中也工作良好；事实上，对于 non SMP 系统来说异步线程池更加重要，因为如果没有这个池，遇到大文件操作时，整个 VM 都会被阻塞住；


----------

原文地址：[这里](http://www.erlang.se/euc/08/euc_smp.pdf)

