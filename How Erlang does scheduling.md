

在本文里，我将说明为什么 Erlang 的运行时会与其它语言的运行时有所不同，还将说明为什么它经常会采取放弃吞吐量的方式以换取更小的延时。

摘要：Erlang 与其它大多数语言的运行时差异源于双方有着不同的目标价值。我们常常发现 Erlang 在进程不多时表现不佳，反而在进程很多时表现优异，原因也在于此。

长久以来，关于 Erlang 的调度方式问题总是被人们问起。尽管本文是真实情况的缩减版本，但仍可作为针对这个问题的解答。需要留意的是，本文是基于 Erlang R15 来写的。如果你是一位来自未来的读者，那么文章里的情况可能会与你面对的现状有所不同。不过一般来说，未来总是会变得更好，无论是 Erlang ，还是系统的其它方面；

从操作系统层面来说，Erlang 通常会在计算机的每个核心上运行一个线程；每一个线程上运行的就是所谓的 scheduler 。这么做的原因，就是为了确保计算机的所有核心都能为 Erlang 系统工作。核心可以通过 `+sbt` 标识与调度器进行绑定；这样一来就可以保证，调度器无法在多个核心间“跳来跳去”。这个标识只能用于现代操作系统上，所以 OSX 是不能用的。这也意味着，Erlang 系统是了解处理器布局及亲缘（affinity）关系的，这些信息对缓存、迁移时间等都非常重要。通常情况下，`+sbt` 标识可以加速系统的运行，并且有些时候，加速很明显。

`+A` 标识定义了异步线程池中的线程数目。该池可被驱动程序用来执行某些阻塞操作，而与此同时，调度器仍可继续完成有用的工作。需要注意的是，该线程池可被文件驱动程序用于加速文件I/O，但不能加速网络 I/O 。

目前为止，我们讨论的都是 OS 内核层面轮廓性的东西，此时还需要介绍一下 Erlang（用户态）进程的概念。当调用 `spawn(fun worker/0)` 时，系统会构建一个新 Erlang 进程，具体来说，会在用户态为其分配一个进程控制块（PCB）。其大小一般为 600+ 字节，32 位和 64 位系统下有所不同。系统将可运行（待调度）进程放入调度器的 run queue 中，并在当它们获得时间片时即可运行。

在深入了解当个调度器之前，我想要先说明一下迁移的原理。进程偶尔会在调度器之间进行迁移，这是一个非常复杂的处理过程。该启发式算法的目的就是为了平衡所有调度器的负荷，保证所有核心能被有效地利用。与此同时，该算法也会考虑工作是否多到要启动新的调度器，如果没有，最好保持调度器处于关闭状态，这将意味着相应的核心可以进入节电模式并保持关闭。没错，Erlang 会尽量省电。调度器还可以在无工可开的情况下从别人那里「偷取」工作。详情可参看[1]。

重要说明：在 R15 中，调度器的启停采取了“延迟“处理模式。意思是说，Erlang/OTP 认为调度器的启停很昂贵，只有它认为真正需要启停的时候才会去做这件事。比方说，现在某个调度器的情况是无工可开。但 Erlang 不会马上让它入睡，而是先让它继续跑一会（spin for a little while），看看有没有新工作进来。如果恰好有，那就可以很快开工了，这样可以做到非常小的延时。但另一方面，这也意味着我们不能用 top(1) 或其它 OS 内核工具来衡量系统的执行效率，而必须用 Erlang 系统的内部调用来确定。许多人之所以认为 R15 比 R14 差，就是因为这个原因。

每个调度器运行着两类作业：进程作业和端口作业。它们运行的时候可以指定优先级，类似于操作系统内核中的方式。我们可以将某个进程标识为高优先级、低优先级或其它优先级。进程作业执行的是进程。端口作业处理的是端口。需要知道的是，Erlang 中的“端口“是一种系统与外界通信的机制；文件、网络 socket、与其它程序交互的管道，都属于“端口“。程序员可以通过添加“端口驱动”的方式，支持新的端口类型，不过那需要写 C 代码。调度器还可以轮询网络 socket ，从中读取数据。


普通进程和端口进程都有 2000 个运行次数（reduction）限制，称作 "reduction budget"；系统中的任何操作都会消耗 reduction 值，这包括：loop 中的函数调用、内置函数（BIF）调用、针对进程堆的垃圾回收[n1]、针对 ETS 表的保存和读取，以及消息发送（接收者邮箱中的消息越多，发送的成本越高）；可以说 reduction 计数无处不在；Erlang 的正则库已经被改写成用 c 实现，但同样需要进行 reduction 消耗计算，所以如果你长时间的进行正则运算，则可能会根据 reduction 消耗情况，多次被其它进程抢占执行；对于 Port 的情况也一样，在 Port 上进行 IO 操作会耗费 reduction 值，发送（分布式）消息也会；系统花费了大量时间以确保任何类型的步进兜会产生相应的 reduction 消耗[n2]； 

In effect, this is what makes me say that Erlang is one of a few languages that actually does preemptive multitasking and gets soft-realtime right. Also it values low latency over raw throughput, which is not common in programming language runtimes.

由于以上的原因，我认为Erlang是真正执行抢占式多任务处理的语言之一，也是正确理解软实时概念的语言之一。同时，Erlang对时延看得比吞吐量更重，这在编程语言里也是不多见的。

To be precise, preemption[2] means that the scheduler can force a task off execution. Everything based on cooperation cannot do this: Python twisted, Node.js, LWT (Ocaml) and so on. But more interestingly, neither Go (golang.org) nor Haskell (GHC) is fully preemptive. Go only switches context on communication, so a tight loop can hog a core. GHC switches upon memory allocation (which admittedly is a very common occurrence in Haskell programs). The problem in these systems are that hogging a core for a while—one might imagine doing an array-operation in both languages—will affect the latency of the system.

更精确地来说，抢占，指的是调度器可以强制让某个任务停止执行。所有基于协作的语言和系统，包括Python、Node.js、LWT(Ocaml)等等，都无法做到这一点。更有趣的是，即使Go(golang.org)和Haskell(GHC)也不完全是抢占式的。Go只在通信时切换上下文，因此只需一个密集的循环即可独占某个核心。GHC则是在内存分配时切换(在Haskell程序中十分常见)。这些系统的问题在于，对核心的独占会影响整个系统的时延，大家可以想象一下在这些语言里执行数组操作的情形。

This leads to soft-realtime[3] which means that the system will degrade if we fail to meet a timing deadline. Say we have 100 processes on our run-queue. The first one is doing an array-operation which takes 50ms. Now, in Go or Haskell/GHC[n3] this means that tasks 2-100 will take at least 50ms. In Erlang, on the other hand, task 1 would get 2000 reductions, which is sub 1ms. Then it would be put in the back of the queue and tasks 2-100 would be allowed to run. Naturally this means that all tasks are given a fair share.

这就是软实时[3]，即当定时失败时系统性能会急剧下降。比如说，在运行队列里共有100个进程。第一个进程执行的是数组操作，需要花费50ms。那么，如果是Go和Haskell/GHC[n3]的话，第2-100个进程就至少要50ms后才能完成。而如果换成Erlang，第一个任务只会得到2000个reduction，相当于1ms。当1ms过去后，系统就会把它放到队列尾部，换后面的任务来运行。所有的任务都能公平地分到属于自己的那一份时间。

Erlang is meticously built around ensuring low-latency soft-realtime properties. The reduction count of 2000 is quite low and forces many small context switches. It is quite expensive to break up long-running BIFs so they can be preempted mid-computation. But this also ensures an Erlang system tend to degrade in a graceful manner when loaded with more work. It also means that for a company like Ericsson, where low latency matters, there is no other alternative out there. You can't magically take another throughput-oriented language and obtain low latency. You will have to work for it. And if low latency matters to you, then frankly not picking Erlang is in many cases an odd choice.

Erlang是为保证低时延和软实时而精心打造的。2000的递减量非常短，因此会导致大量的上下文切换。而且，将需要长时间运行的BIF打碎成中型运算，也非常昂贵。但这保证了Erlang系统在加入更多任务时，性能的下降是平稳的。在Ericsson这样需要低时延的公司里，Erlang系统是无可替代的。不能指望换成一种 面向吞吐量的语言，还能得到这么低的时延。你得自己去搞了。坦白地说，如果你确实需要低时延，那么不选择Erlang真的是很奇怪。

[1] "Characterizing the Scalability of Erlang VM on Many-core Processors" http://kth.diva-portal.org/smash/record.jsf?searchId=2&pid=diva2:392243
[2] http://en.wikipedia.org/wiki/Preemption_(computing)
[3] http://en.wikipedia.org/wiki/Real-time_computing

[n1] Process heaps are per-process so one process can't affect the GC time of other processes too much.
[n2] This section is also why one must beware of long-running NIFs. They do not per default preempt, nor do they bump the reduction counter. So they can introduce latency in your system.
[n3] Imagine a single core here, multicore sort of "absorbs" this problem up to core-count, but the problem still persists.

[1] 《多核心处理器中Erlang虚拟机可扩展性的特征》 http://kth.diva-portal.org/smash/record.jsf?searchId=2&pid=diva2:392243 
[2] http://en.wikipedia.org/wiki/Preemption_(computing) 
[3] http://en.wikipedia.org/wiki/Real-time_computing 

[n1] 进程堆是独立进程的，因此一个进程并不能影响其它进程的 GC时间。 
[n2] 本节也说明了为什么我们必须避免长时间运行的NIFS。他们不是每个默认抢占，也不是他们触发减数计数器。因此，他们会使系统延迟。 
[n3] 设想一个单核的情况，多核依靠着核心数量近似的消化”这些问题，但问题仍然存在。

(Smaller edits made to the document at Mon 14th Jan 2013)


----------

原文地址：[这里](http://jlouisramblings.blogspot.dk/2013/01/how-erlang-does-scheduling.html)

