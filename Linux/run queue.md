# run queue

> In modern computers many processes run at once. **Active processes** are placed in an array called a `run queue`, or `runqueue`. The run queue may contain **priority** values for each process, which will be used by the **scheduler** to determine which process to run next. To ensure each program has a fair share of resources, each one is run for some **time period** (`quantum`) before it is paused and placed back into the run queue. When a program is stopped to let another run, the program with the highest priority in the run queue is then allowed to execute.
>
> Processes are also removed from the run queue when they ask to sleep, are waiting on a resource to become available, or have been terminated.

关键：

- **Active processes** 被放在称作 `run queue` 的 array 中；
- run queue 中可能会包含针对每个 process 的 **priority** 值；
- Processes 被**从 run queue 中移出**的情况：
    - **Processes 主动要求 sleep 时** ；
    - **Processes 等待指定 resource 变成可用状态时**；
    - **Processes 被终止运行时**；

> In the Linux operating system (prior to kernel 2.6.23), **each CPU in the system is given a run queue**, which maintains both an **active** and **expired** array of processes. Each array contains **140** (one for each **priority level**) pointers to doubly linked lists, which in turn reference all processes with the given priority. The scheduler selects the next process from the active array with highest priority. When a process' quantum expires, it is placed into the expired array with some priority. When the active array contains no more processes, the scheduler swaps the active and expired arrays, hence the name `O(1) scheduler`.

关键：

- **每一个 CPU 都有一个 run queue** ；
- 每一个 run queue 都由 **active array** 和 **expired array** 构成；
- 每一个 array 中包含 **140** 个指针，分别指向对应不同 **priority** 的双向链表；
- scheduler 从 active array 中获取 process 进行执行，时间片用完之后，放回 expired array ；
- scheduler 通过交换 active array 和 expired array 实现了所谓的 `O(1) scheduler` ；


> In UNIX or Linux, the `sar` command is used to check the run queue. The `vmstat` UNIX or Linux command can also be used to determine the number of processes that are **queued to run** or **waiting to run**. These appear in the 'r' column.


----------

## 只言片语

Each process using or waiting for CPU (the ready queue or run queue) increments the load number by 1. Each process that terminates decrements it by 1. Most UNIX systems count only processes in the running (on CPU) or runnable (waiting for CPU) states. However, Linux also includes processes in uninterruptible sleep states (usually waiting for disk activity), which can lead to markedly different results if many processes remain blocked in I/O due to a busy or stalled I/O system.


----------



## 参考

- [Run_queue from Wikipedia](https://en.wikipedia.org/wiki/Run_queue)



