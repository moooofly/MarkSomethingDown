# Context Switch

![Context_switch](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/Context_switch.png "Context_switch")

> In computing, a `context switch` is the process of **storing** and **restoring** the state (more specifically, the `execution context`) of a `process` or `thread` so that execution can be resumed from the same point at a later time. This enables multiple processes to share a single CPU and is an essential feature of a multitasking operating system.

context switch 一般是指 execution context switch ；作用对象为 `process` 或 `thread` ；

> The precise meaning of "context switch" varies significantly in usage, most often to mean "**thread switch** or **process switch**" or "**process switch only**", either of which may be referred to as a "`task switch`". More finely, one can distinguish **thread switch** (switching between two threads within a given process), **process switch** (switching between two processes), **mode switch** (domain crossing: switching between user mode and kernel mode within a given thread), **register switch**, a **stack frame switch**, and **address space switch** (memory map switch: changing virtual memory to physical memory map). The computational cost of context switches varies significantly depending on what precisely it entails, from little more than a subroutine call for light-weight user processes, to very expensive, though typically much less than that of saving or restoring a process image.

在广义上，context switch 可以指：

- **thread switch** (task switch)
- **process switch** (task switch)
- **mode switch** (switching between user mode and kernel mode within a given thread))
- **register switch**
- **stack frame switch**
- **address space switch** (memory map switch: changing virtual memory to physical memory map)

## Context switche 的 Cost

上述 switch 的成本有大有小；

> Context switches are usually computationally intensive, and much of the design of operating systems is to optimize the use of context switches. **Switching from one process to another** requires a certain amount of time for doing the administration – saving and loading registers and memory maps, updating various tables and lists, etc. What is actually involved in a context switch varies between these senses and between processors and operating systems. For example, in the Linux kernel, context switching involves switching registers, stack pointer, and program counter, but is independent of address space switching, though in a process switch an address space switch also happens. Further still, analogous context switching happens between user threads, notably green threads, and is often very light-weight, saving and restoring minimal context. In extreme cases, such as switching between goroutines in Go, a context switch is equivalent to a coroutine yield, which is only marginally more expensive than a subroutine call.

这段描述了 context switching 在 kernel 层次，process 层次、thread 层次、goroutines/coroutine/subroutine 层次上的成本差异；

## 何时发生 switch

存在三种可能触发 context switch 的情况；

### Multitasking

> Most commonly, within some scheduling scheme, one process must be switched out of the CPU so another process can run. This context switch can be **triggered by the process making itself unrunnable**, such as by waiting for an I/O or synchronization operation to complete. On a pre-emptive multitasking system, **the scheduler may also switch out processes which are still runnable**. To prevent other processes from being starved of CPU time, preemptive schedulers often configure a timer interrupt to fire when a process exceeds its time slice. This interrupt ensures that the scheduler will gain control to perform a context switch.

在多任务调度中：

- process 由于自身原因**进入 unrunnable 状态**（可能原因：**等待 I/O 完成**或**等待同步操作完成**）
- scheduler 在指定 processes 的调度时间片到达后，**强行切出**仍处于 runnable 状态的 processes ，以避免其它 processes 由于无法获得 CPU time 而被饿死；


### Interrupt handling

> **Modern architectures are interrupt driven**. This means that if the CPU requests data from a disk, for example, it does not need to busy-wait until the read is over; it can issue the request and continue with some other execution. When the read is over, the CPU can be interrupted and presented with the read. For interrupts, a program called an interrupt handler is installed, and it is the interrupt handler that handles the interrupt from the disk.
>
> **When an interrupt occurs, the hardware automatically switches a part of the context** (at least enough to allow the handler to return to the interrupted code). **The handler may save additional context**, depending on details of the particular hardware and software designs. Often only a minimal part of the context is changed in order to minimize the amount of time spent handling the interrupt. **The kernel does not spawn or schedule a special process to handle interrupts**, but instead the handler executes in the (often partial) context established at the beginning of interrupt handling. Once interrupt servicing is complete, the context in effect before the interrupt occurred is restored so that the interrupted process can resume execution in its proper state.

这里描述的是**和 Interrupt 处理相关**的 context switch ；


### User and kernel mode switching

> When a transition between user mode and kernel mode is required in an operating system, a context switch is not necessary; **a mode transition is not by itself a context switch**. However, depending on the operating system, a context switch may also take place at this time.

**mode transition** 并不等价于 context switch ；但 context switch 有可能发生在 mode transition 的过程之中；

## Steps

> In a switch, the state of process currently executing must be saved somehow, so that when it is rescheduled, this state can be restored.
>
> The `process state` includes all the **registers** that the process may be using, especially the **program counter**, plus any other **operating system specific data** that may be necessary. This is usually stored in a data structure called a **process control block** (`PCB`) or `switchframe`.

在 context switch 过程中涉及的内容：

- `process state` 包含了 process 所使用的所有 registers ；
- `process state` 保存在称作 **process control block** (`PCB`) 或 `switchframe` 的地方；

> The PCB might be stored on a **per-process stack in kernel memory** (as opposed to the **user-mode call stack**), or there may be some specific operating system defined data structure for this information. A handle to the PCB is added to a queue of processes that are ready to run, often called the `ready queue`.

关键：

- PCB 保存在 **kernel memory** 中的 per-process stack 里；
- 指向 PCB 的 handle 被添加到 `ready queue` 中，其中保存的是 ready to run 的 processes ；

> Since the operating system has effectively suspended the execution of one process, it can then switch context by choosing a process from the `ready queue` and restoring its `PCB`. In doing so, the program counter from the PCB is loaded, and thus execution can continue in the chosen process. Process and thread priority can influence which process is chosen from the ready queue (i.e., it may be a [priority queue](https://en.wikipedia.org/wiki/Priority_queue)).

**switch context 的底层行为**简单来说就是：**根据 priority 从 `ready queue` 选出一个调度实体，之后恢复其 PCB ，拿到 program counter 后恢复之前的执行；**

## Performance

> Context switching itself has a cost in performance, due to running the task scheduler, [TLB](https://en.wikipedia.org/wiki/Translation_lookaside_buffer) flushes, and indirectly due to sharing the CPU cache between multiple tasks. **Switching between threads of a single process can be faster** than between two separate processes, because **threads share the same virtual memory maps**, so a TLB flush is not necessary.

和 Context switching 性能相关的因素： 

- running the task scheduler
- TLB flushes
- sharing the CPU cache between multiple tasks


----------

> A `context switch` (also sometimes referred to as a `process switch` or a `task switch`) is the switching of the CPU (central processing unit) from one process or thread to another.

简单的讲，`context switch` 就是 `process switch` 或 `task switch` ；

> A **process** (also sometimes referred to as a **task**) is an executing (i.e., running) instance of a program. In Linux, **threads** are lightweight processes that can run in parallel and share an address space (i.e., a range of memory locations) and other resources with their parent processes (i.e., the processes that created them).

process 和 thread 的自身特点；

> **A `context` is the contents of a `CPU's registers` and `program counter` at any point in time.** A `register` is a small amount of very fast memory **inside** of a CPU (as opposed to the slower RAM `main memory` **outside** of the CPU) that is used to speed the execution of computer programs by providing quick access to commonly used values, generally those in the midst of a calculation. A `program counter` is a specialized register that indicates the position of the CPU in its instruction sequence and which holds either the address of the instruction being executed or the address of the next instruction to be executed, depending on the specific system.

关键：

- `context` 由 `CPU's registers` 和 `program counter` 构成（前者包含后者）；
- register 和 main memory 的速度差别；
- program counter 的用途；

> Context switching can be described in slightly more detail as the kernel (i.e., the core of the operating system) performing the following activities with regard to processes (including threads) on the CPU: (1) **suspending** the progression of one process and **storing** the CPU's state (i.e., the context) for that process somewhere in memory, (2) **retrieving** the context of the next process from memory and **restoring** it in the CPU's registers and (3) **returning to** the location indicated by the `program counter` (i.e., returning to the line of code at which the process was interrupted) in order to resume the process.

Context switching 在 kernel 层次的行为；

> A context switch is sometimes described as the kernel suspending execution of one process on the CPU and resuming execution of some other process that had previously been suspended. Although this wording can help clarify the concept, it can be confusing in itself because a process is, by definition, an executing instance of a program. Thus **the wording suspending progression of a process might be preferable.**

将 context switch 描述为挂起 process 的 progression 更为恰当；

## Context Switches and Mode Switches

> **Context switches can occur only in kernel mode.** `Kernel mode` is a privileged mode of the CPU in which only the kernel runs and which provides access to all memory locations and all other system resources. Other programs, including applications, initially operate in user mode, but they can run portions of the kernel code via system calls. A `system call` is a **request** in a Unix-like operating system by an active process (i.e., a process currently progressing in the CPU) **for a service performed by the kernel**, such as input/output (I/O) or process creation (i.e., creation of a new process). `I/O` can be defined as any movement of information to or from the combination of the CPU and main memory (i.e. RAM), that is, communication between this combination and the computer's users (e.g., via the keyboard or mouse), its storage devices (e.g., disk or tape drives), or other computers.

- **Context switches 只能在 kernel mode 中发生**；
- `Kernel mode` 是一种 CPU 的**特权模式**，在该模式下，只有 kernel 能够运行并提供对全部 memory locations 和 system resources 的访问；

> The existence of these two modes in Unix-like operating systems means that a similar, but simpler, operation is necessary when a system call causes the CPU to shift to kernel mode. This is referred to as a `mode switch` rather than a context switch, because it does not change the current process.
>
> Context switching is an essential feature of multitasking operating systems. A multitasking operating system is one in which multiple processes execute on a single CPU seemingly simultaneously and without interfering with each other. This illusion of concurrency is achieved by means of context switches that are occurring in rapid succession (tens or hundreds of times per second). These context switches occur as a result of processes **voluntarily relinquishing** their time in the CPU or as a result of the **scheduler making the switch** when a process has used up its CPU time slice.
> 
> A context switch can also occur as a result of a **hardware interrupt**, which is a signal from a hardware device (such as a keyboard, mouse, modem or **system clock**) to the kernel that an event (e.g., a key press, mouse movement or **arrival of data from a network connection**) has occurred.

context switches 发生于：

- processes 自愿放弃（对照上面内容）；
- scheduler 强制切出；
- 发生硬件中断（**system clock** 和**网络数据包**是主要部分）

> Intel 80386 and higher CPUs contain **hardware support for context switches**. However, most modern operating systems perform **software context switching**, which can be used on any CPU, rather than hardware context switching in an attempt to obtain improved performance. Software context switching was first implemented in Linux for Intel-compatible processors with the 2.4 kernel.
> 
> One major advantage claimed for software context switching is that, whereas the **hardware mechanism saves almost all of the CPU state, software can be more selective and save only that portion that actually needs to be saved and reloaded.** However, there is some question as to how important this really is in increasing the efficiency of context switching. Its advocates also claim that software context switching allows for the possibility of improving the switching code, thereby further enhancing efficiency, and that it permits better control over the validity of the data that is being loaded.

关键：

- context switches 可以通过硬件支持，也可以纯软件实现； 
- 硬件机制会保存几乎全部 CPU state 内容，而软件机制只保存必要的内容；
- 还有其它关于性能方面的争论；

## The Cost of Context Switching

> Context switching is generally **computationally intensive**. That is, it requires considerable processor time, which can be on the order of nanoseconds for each of the tens or hundreds of switches per second. Thus, **context switching represents a substantial cost to the system in terms of CPU time** and can, in fact, be the most costly operation on an operating system.

关键：

- **Context switching 本质上是计算密集的**；
- 每秒数以千计的 switches 会耗费 nanoseconds 级别的处理器时间；
- 从 CPU 时间耗费的角度来说，context switching 成本是很高的；

> Consequently, a major focus in the design of operating systems has been to avoid unnecessary context switching to the extent possible. However, this has not been easy to accomplish in practice. In fact, although the cost of context switching has been declining when measured in terms of the absolute amount of CPU time consumed, this appears to be due mainly to increases in CPU clock speeds rather than to improvements in the efficiency of context switching itself.
> 
> One of the many advantages claimed for Linux as compared with other operating systems, including some other Unix-like systems, is its extremely low cost of context switching and mode switching.

关键：

- 从操作系统设计角度改善 context switching 是可能的，但目前似乎没有从改善 context switching 效率上下手，而只是从增加 CPU clock 速度的角度进行“改善”；
- Linux 上的 context switching 和 mode switching 相比其它 Unix-like 系统成本要低一些；


## 参考

- [Context switch](https://en.wikipedia.org/wiki/Context_switch)
- [Context Switch Definition](http://www.linfo.org/context_switch.html)


----------

## 只言片语

CPU Utilization 需要结合 Load Average 和 Context Switch Rate 来看，因为有可能 CPU Utilization 高的原因是因为后两个指标高导致的；

Context Switch Rate 就是指 Process 或 Thread 的切换速率；如果切换过多，会让 CPU 忙于切换，也会导致影响吞吐量；

Context Switch 大体上由两个部分组成：

- **中断**：一次中断（Interrupt）会引起一次切换
- **进程/线程切换**：进程（线程）的创建、激活之类的也会引起一次切换；

另外，Context Switch 的值也和 TPS (Transaction Per Second) 相关；假设每次调用会引起 N 次 Context Switch ，则有

```
Context Switch Rate = Interrupt Rate + TPS * N
CSR = IR + TPS * N
```

其中

- **TPS * N 对应的就是进程/线程的切换**（待确认）；
- Interrupt Rate 为每秒设备中断数（主要为 **system clock** 和**网络数据包**）；

**内核的 system clock 频率**可以通过如下命令得到：

```
root@vagrant-ubuntu-trusty:~] $ cat /boot/config-`uname -r` | grep '^CONFIG_HZ='
CONFIG_HZ=250
```

或

```
[root@xg-minos-rediscluster-1 ~]# cat /boot/config-`uname -r` | grep '^CONFIG_HZ='
CONFIG_HZ=1000
```

则每秒系统时钟的中断数（timer interrupt）为

```
每秒时钟中断数 = cpu num * core num * CONFIG_HZ
```

和 kernel timer system 相关的一段话：

> **Timer Wheel**, **Jiffies** and **HZ** (or, the way it was)
>
>> The original kernel timer system (called the "timer wheel") was based on incrementing a kernel-internal value (jiffies) every **timer interrupt**. The timer interrupt becomes the default scheduling quamtum, and all other timers are based on jiffies. The **timer interrupt rate** (and **jiffy increment rate**) is defined by a compile-time constant called `HZ`. Different platforms use different values for HZ. Historically, the kernel used **100** as the value for HZ, yielding a jiffy interval of **10 ms**. With 2.4, the HZ value for i386 was changed to **1000**, yeilding a jiffy interval of **1 ms**. Recently (2.6.13) the kernel changed HZ for i386 to 250. (1000 was deemed too high).




