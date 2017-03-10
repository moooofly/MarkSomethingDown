
# [Maximize throughput with RabbitMQ](http://stackoverflow.com/questions/10030227/maximize-throughput-with-rabbitmq)

在我们的项目中，会使用 RabbitMQ 的 "**Task Queues**" 模式来传输数据；

在 **producer** 侧，我们基于 node.js 构建了一些 TCP server 负责接收高并发数据，以及将数据无任何修改的转发给 MQ ；

在 **consumer** 侧，我们使用 JAVA client 从 MQ 中获取、处理、应答数据；

所以**问题**集中在：
- **如何获取最大的消息传输吞吐量？**(例如 400,000 msg/second) 
- **使用多少 queues 是最佳实践？**
- **是否 queue 的数据越多即意味着更好的吞吐量和性能？**
- 是否还有其他的点需要注意的？在该使用场景中是否有任何已知的最佳实践指导原则？

Any comments are highly appreciated!!

----------

为了获得 RabbitMQ 的**最佳性能**，可以遵循其作者给出的使用建议；详见 《[Sizing your Rabbits](http://www.rabbitmq.com/blog/2011/09/24/sizing-your-rabbits/)》:

> - **RabbitMQ 中的 queues 在空状态下是最快的**；    
> - 当一个 queue 处于空状态，并且其下有 consumers 处于准备接收消息的状态时，只要有消息到达该 queue ，就会直接转发给相应的 consumer ；    
> - 对于 **persistent** 消息发送到 **durable** queue 的情况来说，没什么好说的，消息确实会落到磁盘上，但该动作是**异步**完成的，并且相应的消息都会 **buffered heavily**；    
> - 关键点在于：在上述场景中，RabbitMQ 需要进行 book-keeping 的内容几乎可以忽略不计，需要修改的 data structures 也非常少量，并且仅需要分配非常少量的内存；

如果你想要深入理解和 RabbitMQ queues 性能相关的更多内容，可以看看[这篇博客文章](http://www.rabbitmq.com/blog/2011/10/27/performance-of-queues-when-less-is-more/)；

----------

According to a response I once got from the `rabbitmq-discuss` mailing group there are other things that you can try to **增加吞吐量** 和 **降低延迟**：

- **Use a larger prefetch count**. Small values hurt performance.
- A topic exchange is slower than a direct or a fanout exchange.
- **Make sure queues stay short**. Longer queues impose more processing overhead.
- If you care about latency and message rates then **use smaller messages**. **Use an efficient format** (e.g. avoid XML) or **compress the payload**.
- Experiment with **HiPE**, which helps performance.
- **Avoid transactions and persistence**. Also avoid publishing in `immediate` or `mandatory` mode. Avoid HA. Clustering can also impact performance.
- You will achieve better throughput on a multi-core system if you have multiple queues and consumers.
- Use at least `v2.8.1`, which introduces `flow control`. Make sure the memory and disk space alarms never trigger.
- Virtualisation can impose a small performance penalty.
- Tune your OS and network stack. Make sure you provide more than enough RAM. Provide fast cores and RAM.


一点个人看法：
- **prefetch** 应该尽量调大；
- topic 匹配性能问题已经在新版本中有所优化；
- queue 中留存的消息应该尽量少（理想很丰满，现实很残酷）；
- 消息大小越小越能降低延迟，提高消息速率；而将小消息合并成大消息进行处理，能够增大吞吐量；
- 对消息进行高效的表达（更精简的数据格式 or 压缩），能够提升传输性能（业务处理复杂度略微提高）；
- 可以试试 `HiPE` ，对性能提升有好处（在 **Erlang R17** 后 `HiPE` 已经相当稳定）；
- 事务肯定要禁止使用，可靠性通过其他方案解决；
- `immediate` 已经被取消，存在其他等效方案；`mandatory` 适用于某些特定场景，是否使用应该视情况而定；
- 避免 HA ？扯什么蛋！cluster 也会对性能有一定影响？嗯，这是一句没有错误的废话；
- 在多核系统上跑 RabbitMQ 性能会比跑在单核上强（需要注意下 docker 等虚拟环境等使用情况）；
- `flow control` 功能从 `v2.8.1` 开始才有；
- Virtualisation 也会对性能有一定影响，如果这里是指通过 management 插件实现的 Virtualisation ，那么性能问题已经被遇到了！

----------

You will increase the throughput with a **larger prefetch count** AND at the same time **ACK multiple messages** (instead of sending ACK for each message) from your consumer.

But, of course, ACK with `multiple` flag on ([here](http://www.rabbitmq.com/amqp-0-9-1-reference.html#basic.ack)) requires extra logic on your consumer application ([here](http://lists.rabbitmq.com/pipermail/rabbitmq-discuss/2013-August/029600.html)). You will have to keep a list of delivery-tags of the messages delivered from the broker, their status (whether your application has handled them or not) and ACK every N-th delivery-tag (NDTAG) when all of the messages with delivery-tag less than or equal to NDTAG have been handled.

更大的 prefetch 数值 ＋ 一次 ack 多条消息的组合必然可以提高吞吐量，但会增加一定的业务处理复杂度；

