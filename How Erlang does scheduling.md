

In this, I describe why Erlang is different from most other language runtimes. I also describe why it often forgoes throughput for lower latency.

在本文里，我将说明 Erlang 为什么会与其它语言的运行时有所不同，还将说明它为什么总是用放弃吞吐量来换取更小的延时。

TL;DR - Erlang is different from most other language runtimes in that it targets different values. This describes why it often seem to perform worse if you have few processes, but well if you have many.

摘要: Erlang与其它语言的差异源于，双方有着不同的目标价值。我们常常发现Erlang在进程不多时表现不佳，反而在进程很多时表现优异，原因也在于此。

From time to time the question of Erlang scheduling gets asked by different people. While this is an abridged version of the real thing, it can act as a way to describe how Erlang operates its processes. Do note that I am taking Erlang R15 as the base point here. If you are a reader from the future, things might have changed quite a lot—though it is usually fair to assume things only got better, in Erlang and other systems.
长久以来，关于Erlang调度的问题总是被人们问起。这篇文章虽然对真实的情况有所省略，但是大家仍然可以将它看作是针对这个问题的解答。需要留意的是，我在这里是基于Erlang R15来写这篇文章的。所以，如果你是一位来自未来的读者，那么文章里的情况可能会跟你目前的现状有所不同。不过，一般来说，未来总是会变得更好，无论是Erlang，还是其它的什么东东，不是么？

Toward the operating system, Erlang usually has a thread per core you have in the machine. Each of these threads runs what is known as a scheduler. This is to make sure all cores of the machine can potentially do work for the Erlang system. The cores may be bound to schedulers, through the +sbt flag, which means the schedulers will not "jump around" between cores. It only works on modern operating systems, so OSX can't do it, naturally. It means that the Erlang system knows about processor layout and associated affinities which is important due to caches, migration times and so on. Often the +sbt flag can speed up your system. And at times by quite a lot.
Erlang一般在计算机的每个核心上运行一个线程，这个线程就是调度者。这么做，是为了确保计算机的所有核心都能为Erlang工作。核心可以通过+sbt标识绑定到调度者上，这样一来，调度者就无法在多个核心间“跳来跳去”。这个标识只能用于现代操作系统上，所以OSX是不能用的。这也意味着，Erlang系统是了解处理器布局及相互关系的，这些信息对缓存、迁移时间等都非常重要。+sbt标识往往可以加速系统的运行，有时甚至可以大大地加速。

The +A flag defines a number of async threads for the async thread pool. This pool can be used by drivers to block an operation, such that the schedulers can still do useful work while one of the pool-threads are blocked. Most notably the thread pool is used by the file driver to speed up file I/O - but not network I/O.

+A标识是为异步线程池定义线程数。驱动程序可以利用这个池子来阻塞某个操作，而调度者仍然可以继续做有用的工作。需要注意的是，这个线程池可以被文件驱动利用，以加速文件I/O，但不能加速网络I/O。


While the above describes a rough layout towards the OS kernel, we still need to address the concept of an Erlang (userland) process. When you call spawn(fun worker/0) a new process is constructed, by allocating its process control block in userland. This usually amounts to some 600+ bytes and it varies from 32 to 64 bit architectures. Runnable processes are placed in the run-queue of a scheduler and will thus be run later when they get a time-slice.


目前为止，我们讨论的都是OS内核层面的东西，当然也需要介绍一下Erlang(用户域)进程。当调用spawn时，系统会构建一个新进程，具体来说，是在用户域(userland)为它分配一个进程控制块。它的大小一般为600+字节，32位和64位系统下有所不同。系统将可运行的进程放入调度者的运行队列，当它们获得时间片时即可运行。

Before diving into a single scheduler, I want to describe a little bit about how migration works. Every once in a while, processes are migrated between schedulers according to a quite intricate process. The aim of the heuristic is to balance load over multiple schedulers so all cores get utilized fully. But the algorithm also considers if there is enough work to warrant starting up new schedulers. If not, it is better to keep the scheduler turned off as this means the thread has nothing to do. And in turn this means the core can enter power save mode and get turned off. Yes, Erlang conserves power if possible. Schedulers can also work-steal if they are out of work. For the details of this, see [1].

在详细介绍调度者之前，我想要先说明一下迁移的原理。进程偶尔会在调度者之间进行迁移，这是一个非常复杂的过程。它的目的是为了平衡所有调度者的负荷，保证所有核心能被有效地利用。但算法同时也会考虑工作是否多到要启动新的调度者，如果没有，它会保持调度器仍然处于关闭状态，这也意味着相应的核心可以进入节电模式并保持关闭。没错，Erlang会尽量省电。调度者们也可以在无工可开的情况下从别人那里「偷」工作。详情可参看[1]。

IMPORTANT: In R15, schedulers are started and stopped in a "lagged" fashion. What this means is that Erlang/OTP recognizes that starting a scheduler or stopping one is rather expensive so it only does this if really needed. Suppose there is no work for a scheduler. Rather than immediately taking it to sleep, it will spin for a little while in the hope that work arrives soon. If work arrives, it can be handled immediately with low latency. On the other hand, this means you cannot use tools like top(1) or the OS kernel to measure how efficient your system is executing. You must use the internal calls in the Erlang system. Many people were incorrectly assuming that R15 was worse than R14 for exactly this reason.

重点: 在R15中，调度者的启停是有延迟的。意思是说，Erlang/OTP认为调度者的启停很昂贵，只有它认为真正需要启停的时候才会去做这件事。比方说，现在某个调度者的情况是无工可开。但Erlang不会马上让它入睡，而是先让它继续跑一会，看看有没有新工作进来。如果恰好有，那就可以很快开工了，这样可以做到非常小的延时。但另一方面，这也意味着我们不能用top(1)和OS内核来衡量系统的效率，而必须用Erlang系统的内部调用来确定。许多人之所以认为R15比R14差，就是因为这个原因。

Each scheduler runs two types of jobs: process jobs and port jobs. These are run with priorities like in an operating system kernel and is subject to the same worries and heuristics. You can flag processes to be high-priority, low-priority and so on. A process job executes a process for a little while. A port job considers ports. To the uninformed, a "port" in Erlang is a mechanism for communicating with the outside world. Files, network sockets, pipes to other programs are all ports. Programmers can add "port drivers" to the Erlang system in order to support new types of ports, but that does require writing C code. One scheduler will also run polling on network sockets to read in new data from those.

每个调度者运行着两类作业: 进程作业和端口作业。它们运行的时候也带优先级，类似于操作系统内核。我们可以将某个进程标识为高优先级、低优先级或其它优先级。进程作业执行的是进程。端口作业处理的是端口。Erlang中的“端口”是系统与外界通信的机制，文件、网络socket、通往其它程序的管道，都是端口。程序员可以通过加入“端口驱动”，以支持新类型的端口，不过那需要写C代码。调度者还可以轮询网络socket，从中读取数据。

Both processes and ports have a "reduction budget" of 2000 reductions. Any operation in the system costs reductions. This includes function calls in loops, calling built-in-functions (BIFs), garbage collecting heaps of that process[n1], storing/reading from ETS, sending messages (The size of the recipients mailbox counts, large mailboxes are more expensive to send to). This is quite pervasive, by the way. The Erlang regular expression library has been modified and instrumented even if it is written in C code. So when you have a long-running regular expression, you will be counted against it and preempted several times while it runs. Ports as well! Doing I/O on a port costs reductions, sending distributed messages has a cost, and so on. Much time has been spent to ensure that any kind of progress in the system has a reduction cost[n2].

普通进程和端口进程都有2000个运行次数限制,任何系统操作都会消耗这个限制, 包括循环调用,内置函数调用,垃圾回收,ets的读取,发送消息(接受者的邮箱信息越多,发送的成本越高).这真是无处不在. btw, erlang的正则库已经被改进,甚至用c写, 所以如果你长时间的进行正则运算, 对你不利并会占用消耗几个限制次数. port也是. 在ports上做io操作消耗限制,发送分布消息也是. 大量的限制消耗以保证任何进程都有一定的运行次数. 

In effect, this is what makes me say that Erlang is one of a few languages that actually does preemptive multitasking and gets soft-realtime right. Also it values low latency over raw throughput, which is not common in programming language runtimes.

由于以上的原因，我认为Erlang是真正执行抢占式多任务处理的语言之一，也是正确理解软实时概念的语言之一。同时，Erlang对时延看得比吞吐量更重，这在编程语言里也是不多见的。

To be precise, preemption[2] means that the scheduler can force a task off execution. Everything based on cooperation cannot do this: Python twisted, Node.js, LWT (Ocaml) and so on. But more interestingly, neither Go (golang.org) nor Haskell (GHC) is fully preemptive. Go only switches context on communication, so a tight loop can hog a core. GHC switches upon memory allocation (which admittedly is a very common occurrence in Haskell programs). The problem in these systems are that hogging a core for a while—one might imagine doing an array-operation in both languages—will affect the latency of the system.

更精确地来说，抢占，指的是调度者可以强制让某个任务停止执行。所有基于协作的语言和系统，包括Python、Node.js、LWT(Ocaml)等等，都无法做到这一点。更有趣的是，即使Go(golang.org)和Haskell(GHC)也不完全是抢占式的。Go只在通信时切换上下文，因此只需一个密集的循环即可独占某个核心。GHC则是在内存分配时切换(在Haskell程序中十分常见)。这些系统的问题在于，对核心的独占会影响整个系统的时延，大家可以想象一下在这些语言里执行数组操作的情形。

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

