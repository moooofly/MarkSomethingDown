


本文用于梳理 run queue 相关资料；

----------

以下内容取自《[learn you some Erlang](http://learnyousomeerlang.com/content)》

# Concurrency Implementation

轻量级进程 ＋ 异步消息传递是 Erlang 的两大法宝，但是你知道 Erlang 的实现者是如何做到的么？

首先，不能相信操作系统处理进程的方式和能力；

操作系统有很多种处理进程的方式，但不同方式下性能差距很大；其中大部分（甚至有可能是全部）实现方式对于 Erlang 应用来说要么太慢，要么太重；采取在 VM 中进行进程处理的方式，Erlang 实现者就能够在可优化性和可靠性上进行把控；在当前实现下，每个 Erlang 进程会占用 300 个字的内存空间，并可在微秒时间粒度上被创建；而这种能力在大多数操作系统上当前都是无法做到的；

![VM-scheduler-runqueue](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/VM-scheduler-runqueue.png "VM-scheduler-runqueue")

为了处理目标程序中可能创建出来的所有 Erlang 进程，VM 会在每个核心上启动一个线程作为 scheduler 运行；而每个运行中的 scheduler 都具有一个 run queue ，Erlang 进程会被放入其中等待时间片可用时被调度；

当其中某个 scheduler 的 run queue 中被放入了过多的任务时，任务中的一部分会被迁移到其它 run queue 中；这意味着每一个 Erlang VM 都需要参与到全部工作任务的负载均衡处理当中，因此编程人员就不需要担心相应的问题了；VM 同样会进行一些其它方面的优化处理，例如限制消息被发送到过载进程的速率，以便能够对负载压力进行调节和分散；


# symmetric multiprocessing and you

使能或不使能 SMP 从结果上来看没多少差别；为了验证此结论，你可以通过执行 `erl -smp disable`  命令启动 Erlang VM ；若想确认你的 Erlang VM 是否支持 SMP ，可以采用不带任何参数启动 VM 的方式，之后查看输出的第一行信息；如果你看到了类似 `[smp:2:2]` 的文本，则说明支持 SMP ，并且是在两个核心上使用了两个 run queue ；如果没有看到类似的文本，则意味着不支持 SMP ；
 
文本 `[smp:2:2]` 意味着有两个核心可用，并启动了两个 schedulers（每一个 scheduler 都具有一个 run queue）；在 Erlang 的早期版本中，你可能发现多个 scheduler 使用唯一一个共享 run queue 的情况；但是从 R13B 开始，每个 scheduler 都拥有一个单独的 run queue 了，因此提供了更好的并发性能； 


----------



以下内容取自：《[About Erlang/OTP and Multi-core performance in particular](http://erlang-factory.com/upload/presentations/105/KennethLundin-ErlangFactory2009London-AboutErlangOTPandMulti-coreperformanceinparticular.pdf)》

# Erlang VM 的演变（SMP & run queue）

![Erlang (non SMP)](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Erlang (non SMP).png "Erlang (non SMP)")

![Erlang SMP VM (before R13)](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Erlang SMP VM (before R13).png "Erlang SMP VM (before R13)")

![Erlang SMP VM (R13)](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Erlang SMP VM (R13).png "Erlang SMP VM (R13)")


迁移逻辑：

- 尽量保证可运行进程最大数目在所有 scheduler 中相等；
- 负载均衡行为触发于 scheduler 达到了设定的 reduction 最大值时；
   - 针对全部 scheduler 的 run-queue 的最大长度统计信息进行收集；
   - 计算每个 run queue/prio 的平均限度，并建立迁移路径信息；
   - 从高于限度值的 scheduler 中取走任务，将任务分给低于限度值的 scheduler；
   - 迁移行为发生于当 scheduler 完成了某个任务的调度，并继续进行下一次调度，直到触发限制值或者新的负载进行操作发生时；
- 存在 work-stealing 机制，触发于某个 scheduler 的 run queue 为空时；
- Running on full load or not!
   - 如果全部 scheduler 都未处于满载状态，任务会被迁移到具有更低 id 的 scheduler 上，这会导致某些 scheduler 进入 inactive 状态；


----------



以下内容取自：《[Characterizing the Scalability of Erlang VM on Many-core Processors](http://kth.diva-portal.org/smash/get/diva2:392243/FULLTEXT01)》


# Erlang’s Concurrency primitives

`Spawn`, “`!`” (send) 和 `receive` 为 Erlang 并发编程的三大原语；这些原语允许一个进程创建新的进程，通过异步消息传递和其它进程进行通信；当创建一个进程时，节点名，模块名，函数名，以及函数参数都会被传入 BIF spawn() 中；`spawn` 调用成功后会返回进程 id ；消息可以通过 `Pid ! Message` 结构进行发送，其中 `Pid` 为进程 id ，`Message` 为任意合法的 Erlang 数据类型的值；receive 语句用于从进程的消息队列中获取消息，使用方式如下：

```erlang
receive
    Pattern1 when Guard1 −> expressions1 ; 
    Pattern2 when Guard2 −> expressions2 ; 
    Other −> expressionsother
after % optional clause 
    Timeout −> expressionstimeout
end
```

在上述语句中，`after` 关键字（实现 timeout 机制），`other` 字句和 guard 断言都是可选的；当一个进程的 receive 执行时，VM 会检查进程消息队列中的每一条消息以确认其是否匹配表达式中给出的某个模式；模式的匹配是按照顺序进行的；如果某个模式和相应的 guard 成功匹配，其后的表达式将被求值，之后其它的模式将不再继续进行匹配；当消息队列中没有任何消息时，或者没有任何消息能匹配模式时，当前进程将被挂起，并调度出 scheduler ；被挂起的进程需要等待新消息的到来以重新变回可运行状态，进而被重新放入当前进程关联的 scheduler 的 run queue 之中；. Then when the process is selected to execute, the new message is matched to the pat- terns in the receive statement again. It is possible that the new message doesn’t match any patterns, and the process is suspended once more. Sometimes, the last pattern other is set to match all messages, and if a message doesn’t match any previous patterns, the expressions following the last pattern will be executed and the message is removed from the message queue.
When there is an after clause and the process is suspended waiting for a message, it will be woken up after Timeout milliseconds if it doesn’t receive a matching message during that time and then the corresponding expressions are executed.


# Erlang Runtime System

Currently BEAM1 is the standard virtual machine for Erlang, originating from Turbo Erlang [16]. It is an efficient register-based abstract machine2. The first experimental implementation of SMP (parallel) VM occurred in 1998 as a result of a master degree project [17]. From 2006, the SMP VM is included in official releases.

The SMP Erlang VM is a multithreaded program. On Linux, it utilizes POSIX thread (Pthread) libraries. Threads in an OS process share a memory space. An Erlang scheduler is a thread that schedules and executes Erlang processes and ports. Thus it is both a scheduler and a worker. Scheduling and execution of processes and ports are interleaved. There is a separate run queue for each scheduler storing the runnable processes and ports associated with it. On many-core processors, the Erlang VM is usually configured with one scheduler per core or one scheduler per hardware thread if hardware multi-threading is supported.

The Erlang runtime system provides many features often associated with operating systems, for instance, memory management, process scheduling and networking. In the remainder of this chapter, we will introduce and analyze the different parts of the current SMP VM implementation (R13B04 as mentioned before) which are relevant to the scalability on many-core processors, including process structure, message passing, scheduling, synchronization and memory management.



# Message Passing

Message passing between two processes on the same node is implemented by copying the message residing on the heap of the sending process to the heap of the receiving process. In the SMP VM, when sending a message, if the receiving process is executing on another scheduler, its heap cannot accommodate the new message or another mes- sage is being copied to it by another process, the sending process allocates a temporary heap fragment for the receiving process to store the new message. The heap fragments of a process are merged into its private heap during garbage collection. After copying, a management data structure containing a pointer to the actual message is put at the end of the receiving process’ message queue. Then the receiving process is woken up and appended to a run queue if it is suspended. In the SMP VM, the message queue of a process actually consists of two queues. Other processes send messages to the end of its external or public queue. It is protected by locks to achieve mutual exclusion (see Section 3.4). A process usually works on its private queue when retrieving messages in order to reduce the overhead of lock acquisition. But if it can’t find a matching mes- sage in the private queue, the messages in the public queue are removed and appended to the private queue. After that these messages are matched. The public queue is not required in the sequential Erlang VM and there is only one queue.

If a process sends a message to itself, the message doesn’t need to be copied. Only a new management data structure with a pointer to it is allocated. The management data in the public queue of the process cannot contain pointers into its heap, since data in the public queue are not in the root set of garbage collection. As a result, the management data pointing to a message in the heap is put to the private queue which is a part of the root set, and otherwise the message would be lost during garbage collection. But before the management data pointing into the heap is appended, earlier management data in the public queues have to be merged into the private queue. The order in which the messages arrive is always maintained. Messages in the heap fragments are always reserved during garbage collection. The message queue of a process is a part of its PCB and not stored in the heap.

A process executing receive command checks its message queue for a message which matches one of the specified patterns. If there is a matching message, the cor- responding management data are removed from the queue, and related instructions are executed. If there is no matching message, the process is suspended. When it is woken up after receiving a new message and scheduled to run, the new message is examined against the patterns. If it is not matching, the process is suspended again.

Since messages are sent by copying, Erlang messages are expected to be small. This also applies to arguments passed to newly spawned processes. The arguments cannot be placed in a memory location that is shared by different processes. They are copied every time a process is spawned.

Message passing can affect the scalability of the Erlang VM on many-core proces- sors. First, on many-core systems access to the external message queue of a process has to be synchronized which introduces overhead. Second, the allocation and release of memory for messages and their management data also require synchronization. All the scheduler threads in a node acquire memory from a common memory space of an OS process which needs to be protected. A memory block for a message or a manage- ment data structure may be allocated from a memory pool whose memory can only be assigned by the sending scheduler. But if the message or management data structure is sent to a process on another scheduler, when the memory block is deallocated and put back to its original memory pool, synchronization is still required to prevent multiple schedulers from releasing memory blocks to the pool simultaneously. Third, if many processes can run in parallel, their messages can be sent in an order that is quite differ- ent from the order in which they are sent on the sequential Erlang VM. When messages arrive differently, the time spent on message matching can vary, which means the work- load can change. As a result, the number or frequency of message passing in an Erlang application has an influence on the scalability. It is also affected by how the messages are sent and received.


# Scheduling

There are four types of work that have to be scheduled, process, port, linked-in driver and system-level activity. System-level tasks include checking I/O activities such as user input on the Erlang terminal. Linked-in driver is another mechanism for integrat- ing external programs written in other languages into Erlang. While with normal port the external program is executed in a separate OS process, the external program written as a linked-in driver is executed as a thread in the OS process of an Erlang node. It also relies on a port to communicate with other Erlang processes. The following description of scheduler is focused on scheduling processes.

## Overview

Erlang schedulers are based on reduction counting as a method for measuring execution time. A reduction is roughly equivalent to a function call. Since each function call may take a different amount of time, the actual periods are not the same between different reductions. When a process is scheduled to run, it is assigned a number of reductions that it is allowed to execute (by default 2000 reductions in R13B04). The process can execute until it consumes all its reduction quantum or pauses to wait for a message. A process waiting for a message is rescheduled when a new message comes or a timer expires. Rescheduled or new processes are put to the end of corresponding run queues. Suspended (blocked) processes are not stored in the run queues.

There are four priorities for processes: maximum, high, normal and low. Each scheduler has one queue for the maximum priority and another queue for the high priority. Processes with the normal and low priority share the same queue. Thus in the run queue of a scheduler, there are three queues for processes. There is also a queue for ports. The queue for each process priority or port is called priority queue in the remainder of the report. In total, a scheduler’s run queue consists of four priority queues storing all the processes and ports that are runnable. The number of processes and ports in all priority queues of a run queue is regarded as run queue length. Processes in the same priority queue are executed in round-robin order. Round-robin is a scheduling algorithm that assigns equal time slice (here a number of reductions) to each process in circular order, and the processes have the same priority to execute.

A scheduler chooses processes in the queue with the maximum priority to execute until it is empty. Then it does the same for the queue with the high priority. When there are no processes with the maximum or high priority, the processes with the normal priority are executed. As low priority and normal priority processes are in the same queue, the priority is realized by skipping a low priority process for a number of times before executing it.

Another important task of schedulers is balancing workload on multiple processors or cores. Both work sharing and stealing [7] approaches are employed. In general, the workload is checked and shared periodically and relatively infrequently. During a period, work stealing is employed to further balance the workload. Every period one of the schedulers will check the load condition on all schedulers (or run queues). It determines the number of active schedulers for the next period based on the load of the current period. It also computes migration limit, which is the target number of processes or ports, for each priority queue of a scheduler based upon the system load and availability of the queue. Then it establishes migration paths indicating which priority queues should push work to other queues and which priority queues should pull work from other queues.

After the process and port migration relationships are settled, priority queues with less work will pull processes or port from their counterparts during their scheduling time slots, while priority queues with more work will push tasks to other queues. Scheduling time slots are interleaved with time slots (or slices) for executing processes, ports and other tasks. When a system is under loaded and some schedulers are inac- tive, the work is mainly pushed by inactive schedulers. Inactive schedulers will become standby after all their work is pushed out. But when a system in full load and all available schedulers are active, the work is mainly pulled by schedulers which have less workload.
If an active scheduler has no work left and it cannot pull work from another sched- uler any more, it tries to steal work from other schedulers. If the stealing is not suc- cessful and there are no system-level activities, the scheduler thread goes into waiting state. It is in the state of waiting for either system-level activities or normal work. In normal waiting state it spins on a variable for a while waiting to be woken by another scheduler. If no other scheduler wakes it up, the scheduler thread is blocked on a con- ditional variable (see Subsection 3.4.6). When a scheduler thread is blocked, it takes longer time to wake it up. A scheduler with high workload will wake up another wait- ing scheduler either spinning or blocked. The flowchart in Figure 3.3 shows the major parts of the scheduling algorithm in the SMP VM. The balance checking and work stealing are introduced in more details in the remainder of this section.

## Number of Schedulers


The load of an Erlang system (a node) is checked during a scheduling slot of an ar- bitrary scheduler when a counter in it reaches zero. The counter in each scheduler is decreased every time when a number of reductions are executed by processes or ports on that scheduler. The counter in the scheduler which checks balance is reset to a value (default value 2000*2000 in R13B04) after each check. As a result, the default period between two balance checks is the time spent in executing 2000*2000 reductions by the scheduler which does the balance checks. If a scheduler has executed 2000*2000 reductions and finds another scheduler is checking balance, it will skip the check, and its counter is set to the maximum value of the integer type in C. Thus in every period there is only one scheduler thread checking the load.
The number of scheduler threads can be configured when starting the Erlang VM. By default it is equal to the number of logical processors in the system. A core or hardware thread is a logical processor. There are also different options to bind these threads to the logical processors. User can also set only a part of the scheduler threads on-line or available when starting the Erlang VM, and by default all schedulers are available. The number of on-line schedulers can be changed at runtime. When running, some on-line schedulers may be put into inactive state according the workload in order to reduce power consumption. The number of active schedulers is set during balance checking. It can increase in the period between two consecutive balance checks if some inactive schedulers are woken up due to high workload. Some of the active schedulers may be out of work and in the waiting state.
As illustrated in Figure 3.4, the active run queues (or schedulers) are always the ones with the smallest indices starting from 0 (1 for schedulers), and the run queues which are not on-line have the largest indices. Off-line schedulers are suspended after initialization.
The objectives of balance check are to find out the number of active schedulers, establish process and port migration paths between different schedulers, and set the target process or port number for each priority queue. The first step of balance checking is to determine the number of active schedulers for the beginning of the next period based on the workload of the current period. Then if all the on-line schedulers should be active, migration paths and limits are determined to share workload between priority queues.





