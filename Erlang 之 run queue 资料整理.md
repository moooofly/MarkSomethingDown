


本文记录和 run queue 相关的资料；

----------

以下内容取自《[learn you some Erlang](http://learnyousomeerlang.com/content)》

# Concurrency Implementation

轻量级进程 ＋ 异步消息传递是 Erlang 的两大法宝，但是你知道 Erlang 的实现者是如何做到的么？

First of all, the operating system can’t be trusted to handle the processes. 
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