
> 原文地址：[RabbitMQ进程结构分析与性能调优](https://www.qcloud.com/community/article/135)

# RabbitMQ 进程

- `tcp_acceptor` 进程接收客户端连接，并创建 rabbit_reader、rabbit_writer、rabbit_channel 进程；
- `rabbit_reader` 负责接收客户端连接，解析 AMQP 帧；
- `rabbit_writer` 负责向客户端返回数据；
- `rabbit_channel` 负责解析 AMQP 方法，对消息进行路由，然后发给相应队列进程；
- `rabbit_amqqueue_process` 即队列进程，在 RabbitMQ 启动（恢复 durable 类型队列）或创建队列时创建；
- `rabbit_msg_store` 是负责消息持久化的进程；


# RabbitMQ 流控

- 针对内存和磁盘使用量设置阈值；
- 在内部实现流控（Flow Control）机制来确保自身的稳定性；

Erlang 进程之间并不共享内存（binary 类型除外），而是通过消息传递来通信，每个进程都有自己的进程邮箱。Erlang 默认没有对进程邮箱大小设限制，所以当有大量消息持续发往某个进程时，会导致该进程邮箱过大，最终内存溢出并崩溃。

在 RabbitMQ 中，**如果生产者持续高速发送，而消费者消费速度较低时**，在没有流控的情况下，会导致内部进程邮箱大小迅速增大，进而达到 RabbitMQ 的整体内存阈值限制，阻塞生产者（得益于这种阻塞机制，RabbitMQ 本身并不会崩溃），与此同时，RabbitMQ 会进行 page 操作，将内存中的数据持久化到磁盘中。

为了解决该问题，RabbitMQ 使用了一种基于信用证（Credit）的流控机制。即每个消息处理进程都具有一个信用组 `{InitialCredit，MoreCreditAfter}`，默认值为 {200, 50} 。

当消息发送者进程 A 向接收者进程 B 发消息时：

- 对于发送者 A ，每发送一条消息，Credit 数量减 1，直到为 0 后被 block 住；
- 对于接收者 B ，每接收 MoreCreditAfter 条消息，会向 A 发送一条消息，给予 A MoreCreditAfter 个 Credit ；
- 当 A 的 Credit > 0 时，A 可以继续向 B 发送消息；

> 可以看出：基于信用证的流控，消息发送进程的发送速度会被限制在消息处理进程的处理速度内；


# amqqueue 进程与 Paging

消息的存储和队列功能是在 amqqueue 进程中实现的；为了高效处理入队和出队的消息、避免不必要的磁盘 I/O ，amqqueue 进程为消息设计了 4 种状态和 5 个内部队列。

4 种状态包括：

- **alpha**，消息的内容和索引都在内存中；
- **beta**，消息的内容在磁盘，索引在内存；
- **gamma**，消息的内容在磁盘，索引在磁盘和内存中都有；
- **delta**，消息的内容和索引都在磁盘；

对于持久化消息，RabbitMQ 先将消息的内容和索引保存在磁盘中，然后才处于上面的某种状态（即只可能处于 alpha、gamma、delta 三种状态之一）。

5 个内部队列包括：

- **q1**、**q2**、**delta**、**q3**、**q4** ；
- q1 和 q4 队列中只有 alpha 状态的消息；
- q2 和 q3 包含 beta 和 gamma 状态的消息；
- delta 队列是消息按序存盘后的一种逻辑队列，只有 delta 状态的消息。

所以 delta 队列并不在内存中，其他 4 个队列则是由 erlang queue 模块实现。

内部队列消息传递顺序：

```
Q1 => Q2 => delta => Q3 => Q4
```

消息从 q1 入队，q4 出队，在内部队列中传递的过程一般是经 q1 顺序到 q4 。实际执行并非必然如此：**开始时，所有队列都为空，消息直接进入 q4（没有消息堆积时）；内存紧张时，将 q4 队尾部分消息转入 q3 ，进而再由 q3 转入 delta ，此时新来的消息将存入 q1（有消息堆积时）**。

**Paging 就是在内存紧张时触发的**，paging 将大量 alpha 状态的消息转换为 beta 和 gamma ；如果内存依然紧张，继续将 beta 和 gamma 状态转换为 delta 状态。Paging 是一个持续过程，涉及到大量消息的多种状态转换，所以 **Paging 的开销较大，严重影响系统性能**。

## 问题分析

在生产者、消费者均正常情况下，RabbitMQ 压测性能非常稳定，保持在一个恒定的速度。当消费者异常或不消费时，RabbitMQ 则表现极不稳定。

大致意思如下：

- exchange 和队列都是持久化的，消息也是持久化的、固定为 1K ，并且无消费者；
- **在达到内存 paging 阈值后，生产速率会降低，并持续较长时间（的低值）**；
- （通过某种方法获得）内存使用情况表明，在内存中的消息数目只有 18M 内容，（说明）其他消息已经 page 到磁盘中，然而进程内存仍占用 2G ；
- （从 Erlang 自身的内存使用统计数据上看）Erlang 内存使用表明，Queues 占用了 2G ，Binaries 占用了 2.1G ；该情况说明**在消息从内存 page 到磁盘后（即从 q2、q3 队列转到 delta 后），系统中产生了大量的垃圾（garbage），而 Erlang VM 没有进行及时的垃圾回收（GC）**，这导致 RabbitMQ 错误的计算了内存使用量，并持续调用 paging 流程，直到 Erlang VM 隐式垃圾回收；


# RabbitMQ 的内存管理

- RabbitMQ 的内存使用量是在 memory_monitor 进程内进行周期性计算的（负责统计计算）；
- amqqueue 进程会周期性（从 memory_monitor 中）拉取内存使用量；**当内存达到 paging 阈值时，会触发 amqqueue 进程进行 paging；当 paging 发生后，amqqueue 进程每收到一条新消息都会对内部队列进行 page（每次 page 都会计算出一定数目的消息存盘）**；

可行的优化方案是：

- 在 amqqueue 进程将大部分消息 paging 到磁盘后，显式调用 GC ;
- 将 memory_monitor 周期设为 0.5s ，将 amqqueue 的拉取周期设为 1s ，这样就能够达到秒级恢复；
- 去掉（当内存达到 paging 阈值后）对每条消息执行 paging 的操作，转为使用 amqqueue 周期性拉取内存使用量来触发 paging 的方式，这样能够更快将消息 paging 到磁盘，而且保持这个周期内生产速度不下降；

相关讨论：

- Issues: [Improve reduce_memory_usage performance when persisting message to the queue index](https://github.com/rabbitmq/rabbitmq-server/issues/289)
- rabbitmq-users: [3.5.5 server performance with pub and no sub](https://groups.google.com/forum/#!msg/rabbitmq-users/vj_9YGUfDgg/-_fx2BkqAwAJ)
- PR: [forces GC after queue paging](https://github.com/rabbitmq/rabbitmq-server/pull/339)

在读取 Message rates 图时，若看到生产速度有明显的下降，则需要判定（通过**流控分析**）：

- 是否发生了 paging ；
- 消息链路阻塞到了哪里；
- 是否发生了 GC ；

可能的情况：**链路阻塞在 amqqueue 进程**；

- 若发现节点内存使用下降了，则说明该节点执行了 GC ；Erlang 中的 GC 是按**进程级别**的**标记-清扫**模式，会将当前进程暂停，直至 GC 结束（详见[这里](http://erlang.org/faq/academic.html#idp33134160)和[这里](http://prog21.dadgum.com/16.html)）；
- 由于在 RabbitMQ 中，一个队列只对应一个 amqqueue 进程，该进程又会处理大量的消息，产生大量的垃圾，就会导致该进程 GC 较慢，进而基于流控阻塞上游更长时间；

另外，针对 [gen_server 行为模式](http://erlang.org/doc/man/gen_server.html)：发现 amqqueue 进程的 gen_server 模型在正常的逻辑中调用了 hibernate ，而该操作可能导致两次不必要的 GC ；因此，优化掉 hibernate 对系统稳定性有一些帮助；详见下面的说明

> The `gen_server` process can go into **hibernation** (see `erlang:hibernate/3`) if a callback function specifies 'hibernate' instead of a time-out value. This can be useful if the server is expected to be idle for a long time. However, use this feature with care, as `hibernation` implies at least two **garbage collections** (when hibernating and shortly after waking up) and is not something you want to do between each call to a busy server.

针对流控的可能优化方案：用多个 amqqueue 进程来实现一个队列（的功能），这样可以降低 rabbit_channel 被单个 amqqueue 进程阻塞住的概率；同时在单队列的场景下也能更好利用多核的特性。不过该方案对 RabbitMQ 现有的架构改动很大，难度也很大。

> 问题：流控分析如何做？

垃圾回收命令：

```erlang
%% 针对当前进程
garbage_collect()

%% 针对全部进程
[garbage_collect(Pid) || Pid <- processes()]
```

# 参数调优

## Erlang 参数调优

- [Erlang VM Tuning](https://docs.basho.com/riak/kv/2.1.4/using/performance/erlang/)

## RabbitMQ 参数调优

- [RabbitMQ Performance Measurements, part 1](http://www.rabbitmq.com/blog/2012/04/17/rabbitmq-performance-measurements-part-1/)
- [RabbitMQ Performance Measurements, part 2](http://www.rabbitmq.com/blog/2012/04/25/rabbitmq-performance-measurements-part-2/)

具体参数：

- IO_THREAD_POOL_SIZE：CPU 大于或等于 16 核时，将 Erlang 异步线程池数目设为 100 左右，提高文件 IO 性能。

- hipe_compile：开启 Erlang HiPE 编译选项（相当于 Erlang 的 jit 技术），能够[提高性能 20%-50%](https://www.cloudamqp.com/blog/2014-03-31-rabbitmq-hipe.html) 。在 Erlang R17 后 HiPE 已经相当稳定，RabbitMQ 官方也建议开启此选项。

- queue_index_embed_msgs_below：RabbitMQ 3.5 版本引入了将[小消息直接存入队列索引（queue_index）](http://www.rabbitmq.com/persistence-conf.html#index-embedding)的优化，消息持久化直接在 amqqueue 进程中处理，不再通过 msg_store 进程。由于消息在 5 个内部队列中是有序的，所以不再需要额外的位置索引（msg_store_index）。该优化提高了系统性能 10% 左右。

- vm_memory_high_watermark：用于配置内存阈值，建议小于 0.5 ，因为 **Erlang GC 在最坏情况下会消耗一倍的内存**。

- vm_memory_high_watermark_paging_ratio：用于配置 paging 阈值，该值为 1 时，直接触发内存满阈值，阻塞生产者。

- queue_index_max_journal_entries：journal 文件是 queue_index 为避免过多磁盘寻址添加的一层缓冲（内存文件）。对于生产消费正常的情况，消息生产和消费的记录在 journal 文件中一致，则不用再保存；对于无消费者情况，该文件增加了一次多余的 IO 操作。


