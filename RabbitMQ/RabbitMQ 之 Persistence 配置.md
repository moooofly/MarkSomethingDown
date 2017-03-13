
> 官网原文：[这里](https://www.rabbitmq.com/persistence-conf.html)

# Persistence Configuration

> The RabbitMQ persistence layer is intended to give good results in the majority of situations without configuration. However, some configuration is sometimes useful. This page explains how you can configure it. You are advised to read it all before taking any action.

RabbitMQ 针对“持久层（persistence layer）” 提供的默认配置可以保证在大多数情况下工作良好；然而在某些特定场景下，进行配置配置调整还是非常有必要的，本文针对这个问题进行展开；

## How persistence works

> First, some background: both persistent and transient messages can be written to disk. **`Persistent` messages will be written to disk as soon as they reach the queue**, while **`transient` messages will be written to disk only so that they can be evicted from memory while under memory pressure**. **`Persistent` messages are also kept in memory when possible and only evicted from memory under memory pressure**. The "persistence layer" refers to the mechanism used to store messages of both types to disk.

关键点：

- 持久化消息消息在到达 (durable) queue 后，会尽快被写入磁盘；
- 非持久化消息仅在触发内存压力后，需要从内存中驱逐时，才被写入磁盘；
- 持久化消息（在内存条件允许点情况下）同样会在内存中保留一份，并在触发内存压力后，才从内存中被驱逐；
- “**持久层**”指的是将上述两种类型的消息保存到磁盘时使用的机制；

> On this page we say "queue" to refer to an unmirrored queue or a queue master or a queue slave. Queue mirroring happens "above" persistence.

在本文中提及的 "queue" 概念。指的是非镜像 queue 或者镜像情况下的 queue master 或 queue slave ；因为 **queue 的镜像行为是发生在持久行为“之上”的**；

> The persistence layer has two components: the `queue index` and the `message store`. The queue index is responsible for maintaining knowledge about where a given message is in a queue, along with whether it has been delivered and acknowledged. There is therefore one queue index per queue.
>
> The message store is a key-value store for messages, shared among all queues in the server. Messages (the body, and any properties and / or headers) can either be stored directly in the queue index, or written to the message store. There are technically two message stores (one for transient and one for persistent messages) but they are usually considered together as "the message store".

关键：

- 持久层由 `queue index` 和 `message store` 两个组件构成；
- `queue index` 负责维护消息在 queue 中的位置信息，以及消息是否被投递（deliver）和确认（ack）；
- 每个 queue 都有一个对应的 queue index ；
- `message store` 是用于保存消息本身的 k/v 存储，并且在 server 层面共享于所有 queue ；
- 消息（消息体，以及任何属性值，和/或 headers 内容）或者被直接保存到 `queue index` 中，或者被保存到 `message store` 里；
- 从技术实现层面上讲，存在两个 `message store` 分别用于 transient 和 persistent 消息存储，但通常情况下我们两者当成一个统一的 "the message store" ；

## Memory costs

> Under memory pressure, the persistence layer tries to write as much out to disk as possible, and remove as much as possible from memory. There are some things however which must remain in memory:
> 
> - Each queue maintains some metadata for each **unacknowledged** message. The message itself can be removed from memory if its destination is the message store.
> - The `message store` needs an index. The default message store index uses a small amount of memory for every message in the store.

在内存压力下，持久层会将尽量多的消息写出到磁盘，以便从空出尽量多的内存占用；但是，仍然存在一些信息必须留存在内存之中：

- 每一个 queue 都会为每一条未确认消息维护一些元数据；如果消息的“目的地”是 `message store` ，那么消息本身是可以从内存占用中移除的；
- `message store` 本身也会为其内部保存的消息维护一个索引，即 `message store index`，后者会针对每条消息占用少量内存；

## Messages in the queue index

> There are advantages and disadvantages to writing messages to the queue index.
>
> **Advantages**:
> 
> - Messages can be written to disk in one operation rather than two; for tiny messages this can be a substantial gain.
> - Messages that are written to the queue index do not require an entry in the message store index and thus do not have a memory cost when paged out.

> **Disadvantages**:
> 
> - The queue index keeps blocks of a fixed number of records in memory; if non-tiny messages are written to the queue index then memory use can be substantial.
> - If a message is routed to multiple queues by an exchange, the message will need to be written to multiple queue indices. If such a message is written to the message store, only one copy needs to be written.
> - Unacknowledged messages whose destination is the queue index are always kept in memory.

将消息直接保存到 `queue index` 中既有优点，也有缺点；

优点为：

- 消息的落盘能够以一次操作完成，而非两次；对于小消息来说，收益明显；
- 仅写入 `queue index` 的消息不会再占用 `message store index` 中的 entry ，因此，及时后续触发了 page out 行为，也不会有“残留的”内存消耗；

缺点为：

- `queue index` 在内存中维护了具有固定 record 数目的内存块；如果写入到 `queue index` 中的消息比较大，那么内存占用将非常可观；
- 如果一条消息经由 exchange 被路由到多个 queue 中，那么该消息将需要被写入多个 `queue index` ；与之对比，如果一条消息被写入 `message store` ，那么只需要写一个副本；
- 如果未确认消息的“目的地”是 `queue index` ，那么该消息只会被保存在内存中；

> The intent is for very small messages to be stored in the queue index as an optimisation, and for all other messages to be written to the message store. This is controlled by the configuration item `queue_index_embed_msgs_below`. By default, messages with a serialised size of less than **4096** bytes (including properties and headers) are stored in the queue index.

主要想法就是：针对非常小的消息，将其保存在 `queue index` 中是一种非常好的优化手段；而对于其他消息，写入 `message store` 才是最佳实践；可以通过配置项 `queue_index_embed_msgs_below` 对消息大小判定做划分；默认情况下，若消息在序列化后小于 **4096** 字节（包括属性和 header 信息），则会被保存在 `queue index` 中；

> Each queue index needs to keep at least one segment file in memory when reading messages from disk. The segment file contains records for 16,384 messages. Therefore be cautious if increasing `queue_index_embed_msgs_below`; a small increase can lead to a large amount of memory used.

当从磁盘上加载（读取）消息时，每一个 `queue index` 都需要在内存中维护至少一个 segment 文件；每一个 segment 文件中都会包含针对 16,384 条消息的 record；因此，在增大 `queue_index_embed_msgs_below` 配置时需要格外注意，即使微小的调整都可能导致大量的内存占用；

# Accidentally limited persister performance

> It is possible for persistence to underperform because the persister is limited in the number of file handles or async threads it has to work with. In both cases this can happen when you have a large number of queues which need to access the disk simultaneously.

有些情况下，持久化行为变现平平也是可能的；因为 persister 本身也要受**文件句柄**或**异步线程**数目的限制；当你的系统中使用了非常多的 queue ，并且这些 queue 需要访问同时访问磁盘时，这两种限制可能就会被触发；

## Too few file handles

> The RabbitMQ server is typically limited in the number of file handles it can open (on Unix, anyway). Every running network connection requires one file handle, and the rest are available for queues to use. If there are more disk-accessing queues than file handles after network connections have been taken into account, then the disk-accessing queues will share the file handles among themselves; each gets to use a file handle for a while before it is taken back and given to another queue.

RabbitMQ server 可使用的文件句柄数目是受限的；每一个处于运行状态的网络链接占用一个文件句柄，其余的句柄可供所有的 queue 使用；如果实际情况中遇到 disk-accessing queue 的数量超过可用文件句柄的数量（去掉用于连接的部分），那么 disk-accessing queues 则会通过共享这些文件句柄的方式解决不足问题，即每一个 queue 会占用某个文件句柄一段时间，之后文件句柄会被收回，以便给其他 queue 使用；

> This prevents the server from crashing due to there being too many disk-accessing queues, but it can become expensive. The management plugin can show I/O statistics for each node in the cluster; as well as showing rates of `reads`, `writes`, `seeks` and so on it will also show a rate of `reopens` - the rate at which file handles are recycled in this way. **A busy server with too few file handles might be doing hundreds of reopens per second** - in which case its performance is likely to increase notably if given more file handles.

这种实现策略可以有效防止当存在过多 disk-accessing queues 时 server 自身发生崩溃；但是，这种策略是有代价的；从管理插件输出的内容中可以看到集群中每一个节点上的 I/O 统计信息，以及各种速率值：`reads`, `writes`, `seeks` 和 `reopens` ；其中 `reopens` 对应的就是由于上诉策略文件句柄被回收的速率；**一个仅允许使用少量文件句柄的 busy server 可能会导致 n*100+ reopen/s** ；在这种情况下，如果能分配更多的文件句柄使用，则会显著提高服务器性能；

## Too few async threads

> The Erlang virtual machine creates a pool of async threads to handle long-running file I/O operations. These are shared among all queues. Every active file I/O operation uses one async thread while it is occurring. Having too few async threads can therefore hurt performance.

Erlang VM 会创建一个异步线程池用于处理所有 long-running 文件 I/O 操作；该池同样共享于所有 queue ；每一个 active 的文件 I/O 操作都会占用一个异步线程；因此，**若异步线程数目设置的过少会对性能有损害**；

> Note that the situation with async threads is not exactly analogous to the situation with file handles. If a queue executes a number of I/O operations in sequence it will perform best if it holds onto a file handle for all the operations; otherwise we may flush and seek too much and use additional CPU orchestrating it. However, queues do not benefit from holding an async thread across a sequence of operations (in fact they cannot do so).

需要注意**异步线程**和**文件句柄**所面对的问题的不同；如果一个 queue 按顺序执行一组 I/O 操作，那么所有操作都基于同一个文件句柄完成时性能才最好，否则我们可能需要 flush 和 seek 多次，并额外占用 CPU 以完成编排（orchestrating）；与此相反，如果一个 queue 针对一系列操作独占一个异步线程则不会得到额外好处（事实上，queue 也做不到这样）；

> Therefore there should ideally be enough file handles for all the queues that are executing streams of I/O operations, and enough async threads for the number of simultaneous I/O operations your storage layer can plausibly execute.

因此，**需要提供足够多的文件句柄给所有 queue 使用，以满足流式 I/O 操作的需求；需要提供足够多的异步线程，以满足存储层完成“同时” I/O 操作的需求**；

> It's less obvious when a lack of async threads is causing performance problems. (It's also less likely in general; check for other things first!) Typical symptoms of too few async threads include the number of I/O operations per second dropping to zero (as reported by the management plugin) for brief periods when the server should be busy with persistence, while the reported time per I/O operation increases.

由于异步线程数量不足导致的性能问题，其实很难被察觉（通常情况下，这种问题也不应该是最先被怀疑的；建议先排查其他原因）；过少异步线程导致问题的典型症状包括：当 server 忙于持久化时，大部分时间的每秒 I/O 操作数都为零（可以从管理插件上看到），而每次 I/O 操作的时间却增长了；

> The number of async threads is configured by the `+A` argument to the Erlang virtual machine as described [here](http://erlang.org/doc/man/erl.html#async_thread_pool_size), and is typically configured through the envirnment variable `RABBITMQ_SERVER_ERL_ARGS`. The default value is `+A 64`. It is likely to be a good idea to experiment with several different values before changing this.

异步线程数目可以通过 Erlang VM 的 [`+A`](http://erlang.org/doc/man/erl.html#async_thread_pool_size) 参数进行配置，通常情况下，可以基于 RabbitMQ 的环境变量 `RABBITMQ_SERVER_ERL_ARGS` 进行设置；默认值为 `+A 64` ，具体调整可以根据实际情况测试得到；

# Alternate message store index implementations

> As mentioned above, each message which is written to the message store uses a small amount of memory for its index entry. The message store index is pluggable in RabbitMQ, and other implementations are available as plugins which can remove this limitation. (The reason we do not ship any with the server is that they all use native code.) Note that such plugins typically make the message store run more slowly.

如前所述，写入 `message store` 的每一条消息都会占用少量内存用于其相应的 index entry ；而 `message store index` 在 RabbitMQ 的设计中是可插拔的，因此，若使用其他插件提供的相应实现确实能够去除上述内存占用的限制；而没有将这类插件随 server 一起发布的原因是：这些插件均使用了 native code 实现；另外，这些插件实现通常来讲都会令 `message store` 运行的更慢；




