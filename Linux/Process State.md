# Process state

![Process States](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Process_states.png "Process States")

> In a **multitasking** computer system, processes may occupy a variety of states. **These distinct states may not be recognized as such by the operating system kernel**. However, they are a useful abstraction for the understanding of processes.

关键：

- 在 multitasking 计算机系统中才会说 processes 具有多种状态；
- 这些定义明确的状态并非操作系统本身的视角，而是为了方便进行 processes 理解的抽象；


## Primary process states

The following typical process states are possible on computer systems of all kinds. In most of these states, processes are "stored" on **main memory**.

### Created

> (Also called **New**) When a process is first created, it occupies the "`created`" or "`new`" state. **In this state, the process awaits admission to the "`ready`" state**. Admission will be **approved** or **delayed** by a `long-term scheduler`, or `admission scheduler`. Typically in most desktop computer systems, this admission will be **approved automatically**. However, for real-time operating systems this admission may be **delayed**. In a realtime system, admitting too many processes to the "ready" state may lead to **oversaturation** and **[overcontention](https://en.wikipedia.org/wiki/Bus_contention)** of the system's resources, leading to an inability to meet process deadlines.

关键：

- 当 process 被创建时，则处于 "**created**" 或 "**new**" state ；
- 由 "**created**" 或 "**new**" 到 "**ready**" state 的切换，是由 long-term scheduler (admission scheduler) 控制的；并且根据实际的系统类型决定 admission 的发生是 approved automatically 还是 delayed ；

### Ready

> A "`ready`" or "`waiting`" process has been loaded into **main memory** and is awaiting execution on a CPU (to be **[context switched](https://en.wikipedia.org/wiki/Context_switch)** onto the CPU by the `dispatcher`, or `short-term scheduler`). There may be many "ready" processes at any one point of the system's execution—for example, in a one-processor system, only one process can be executing at any one time, and all other "concurrently executing" processes will be waiting for execution.
>
> A `ready queue` or `run queue` is used in [computer scheduling](https://en.wikipedia.org/wiki/Scheduling_(computing)). Modern computers are capable of running many different programs or processes at the same time. However, the CPU is only capable of handling one process at a time. **Processes that are ready for the CPU are kept in a queue for "ready" processes**. Other processes that are waiting for an event to occur, such as loading information from a hard drive or waiting on an internet connection, are not in the **ready queue**.

关键：

- 处于 "`ready`" 或 "`waiting`" 状态的 process 已经被载入到了主存中，等待获取 CPU 时间片执行；
- `ready queue` 中保存的是 ready for the CPU 的 process ，而不是 waiting for an event to occur 的 process ；

### Running

> A process moves into the `running` state when it is chosen for execution. The process's instructions are executed by one of the CPUs (or cores) of the system. There is at most one running process per CPU or core. A process can run in either of the two modes, namely **kernel mode** or **user mode**.

process 被选中执行时进入 `running` state ；

> **Kernel mode**
>
> - Processes in kernel mode can access both: kernel and user addresses.
> - Kernel mode allows **unrestricted access** to hardware including execution of privileged instructions.
> - Various instructions (such as **I/O** instructions and **halt** instructions) are **privileged** and can be executed only in kernel mode.
> - A [system call](https://en.wikipedia.org/wiki/System_call) from a user program leads to a switch to kernel mode.

> **User mode**
> 
> - Processes in user mode can access their own instructions and data but not kernel instructions and data (or those of other processes).
> - **When the computer system is executing on behalf of a user application, the system is in user mode. However, when a user application requests a service from the operating system (via a system call), the system must transition from user to kernel mode to fulfill the request.**
> - User mode avoids various catastrophic failures:
>    - There is an **isolated virtual address space** for each process in user mode.
>    - User mode ensures **isolated execution** of each process so that it does not affect other processes as such.
>    - **No direct access** to any hardware device is allowed.

两种模式的特点；

### Blocked

> A process transitions to a [`blocked`](https://en.wikipedia.org/wiki/Blocking_(computing)) state when it cannot carry on without an external change in state or event occurring. For example, a process may block on a call to an I/O device such as a printer, if the printer is not available. **Processes also commonly block** when they require user input, or require access to a critical section which must be executed atomically. Such critical sections are protected using a synchronization object such as a **semaphore** or **mutex**.

关键：

- 若在没有外部变更或事件触发的情况下 process 无法再继续执行，则进入 `blocked` state ；
- Processes 通常会 block 的原因有：
    - 等待用户输入；
    - 请求对 critical section 对访问权；

### Terminated

> A process may be terminated, either from the "running" state by **completing its execution** or by **explicitly being killed**. In either of these cases, the process moves to the "`terminated`" state. The underlying program is no longer executing, but the process remains in the **process table** as a **zombie process** until its parent process calls the `wait` system call to read its exit status, at which point the process is removed from the process table, finally ending the process's lifetime. If the parent fails to call wait, this continues to consume the process table entry (concretely the process identifier or PID), and causes a resource leak.

关键：

- 当 process 完成指令执行时，或者被 kill 掉时，则进入 `terminated` state ；
- 处于 `terminated` state 的 process 仍然位于进程表中，直到其父进程将其回收；


## Additional process states

> Two additional states are available for processes in systems that support **virtual memory**. In both of these states, processes are "stored" on **secondary memory** (typically a hard disk).

在支持虚拟内存的系统中，还存在两种额外的 state ；

### Swapped out and waiting

> (Also called `suspended and waiting`.) In systems that support virtual memory, a process may be `swapped out`, that is, **removed** from main memory and **placed** on external storage by the scheduler. From here the process may be `swapped back` into the `waiting` state.

处于 waiting state 的 process 可能会被 swapped out ；

### Swapped out and blocked

> (Also called `suspended and blocked`.) Processes that are `blocked` may also be `swapped out`. In this event the process is both swapped out and blocked, and may be `swapped back` in again under the same circumstances as a swapped out and waiting process (although in this case, the process will move to the blocked state, and may still be waiting for a resource to become available).

处于 blocked state 的 process 可能会被 swapped out ；

----------

## PROCESS STATE CODES

Here are the different values that the s, stat and state output specifiers (header "STAT" or "S") will display to describe the state of a process:

- D --> uninterruptible sleep (usually IO)
- R --> **running** or **runnable** (on **run queue**)
- S --> interruptible sleep (waiting for an event to complete)
- T --> stopped, either by a job control signal or because it is being traced.
- W --> **paging** (not valid since the 2.6.xx kernel)
- X --> dead (should never be seen)
- Z --> defunct ("zombie") process, terminated but not reaped by its parent.

For BSD formats and when the stat keyword is used, additional characters may be displayed:

- < --> **high-priority** (not nice to other users)
- N --> **low-priority** (nice to other users)
- L --> has pages locked into memory (for real-time and custom IO)
- s --> is a **session leader**
- l --> is **multi-threaded** (using CLONE_THREAD, like NPTL pthreads do)
- +    is in the **foreground** process group.

----------

| 状态 | 状态解释 | 在 linux 内核中的编码 | 在 top 中的显示 |
| --- | --- | --- | --- |
| TASK_RUNNING | 运行中 or 就绪 | 0 | R |
| TASK_INTERRUPTIBLE | 睡眠状态，表示进程被阻塞，但能响应信号，直到被唤醒 | 1 | S |
| TASK_UNINTERRUPTIBLE | 不可中断等待状态，无法响应信号，该进程等待一个事件的发生或某种系统资源，通常是等待 IO ，如磁盘 IO 、网络 IO ，或其它外设 IO ；如果进程正在等待的 IO 在较长的时间内都没有响应，可能是外设本身出了故障，也可能是挂载的远程文件系统已经不可访问了 | 2 | D |
| __TASK_STOPPED | 暂停状态，进程收到暂停信号后停止运行 | 4 | T |
| __TASK_TRACED | 跟踪状态，类似于 __TASK_STOPPED ，都代表进程暂停下来 | 8 | t |
| EXIT_ZOMBIE | 僵尸状态，此时进程不能被调度，但是 PCB 未被释放 | 16 | Z |
| EXIT_DEAD | 死亡状体，表示一个已终止的进程，其 PCB 已被释放 | 32 | X |

----------


参考：

- [Process States](https://en.wikipedia.org/wiki/Process_state)