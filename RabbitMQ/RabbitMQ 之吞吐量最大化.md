


# [Maximize throughput with RabbitMQ](http://stackoverflow.com/questions/10030227/maximize-throughput-with-rabbitmq)

In our project, we want to use the RabbitMQ in "**Task Queues**" pattern to pass data.

On the **producer** side, we build a few TCP server(in node.js) to recv high concurrent data and send it to MQ without doing anything.

On the **consumer** side, we use JAVA client to get the task data from MQ, handle it and then ack.

So **the question** is: 
- **To get the maximum message passing throughput/performance?**( For example, 400,000 msg/second) 
- **How many queues is best?**
- **Does that more queue means better throughput/performance?**
- And is there anything else should I notice? Any known best practices guide for using RabbitMQ in such scenario?

Any comments are highly appreciated!!

----------

For **best performance** in RabbitMQ, follow the advice of its creators. From the [RabbitMQ blog](http://www.rabbitmq.com/blog/2011/09/24/sizing-your-rabbits/):

> **RabbitMQ's queues are fastest when they're empty**. When a queue is empty, and it has consumers ready to receive messages, then as soon as a message is received by the queue, it goes straight out to the consumer. In the case of a **persistent** message in a **durable** queue, yes, it will also go to disk, but that's done in an **asynchronous** manner and is **buffered heavily**. The main point is that very little book-keeping needs to be done, very few data structures are modified, and very little additional memory needs allocating.

If you really want to dig deep into the performance of RabbitMQ queues, this [other blog entry](http://www.rabbitmq.com/blog/2011/10/27/performance-of-queues-when-less-is-more/) of theirs goes into the data much further.

----------

According to a response I once got from the rabbitmq-discuss mailing group there are other things that you can try to **increase throughput** and **reduce latency**:

- **Use a larger prefetch count**. Small values hurt performance.
- A topic exchange is slower than a direct or a fanout exchange.
- **Make sure queues stay short**. Longer queues impose more processing overhead.
- If you care about latency and message rates then **use smaller messages**. **Use an efficient format** (e.g. avoid XML) or **compress the payload**.
- Experiment with **HiPE**, which helps performance.
- **Avoid transactions and persistence**. Also avoid publishing in immediate or mandatory mode. Avoid HA. Clustering can also impact performance.
- You will achieve better throughput on a multi-core system if you have multiple queues and consumers.
- Use at least v2.8.1, which introduces flow control. Make sure the memory and disk space alarms never trigger.
- Virtualisation can impose a small performance penalty.
- Tune your OS and network stack. Make sure you provide more than enough RAM. Provide fast cores and RAM.


一点个人看法：
- **prefetch** 应该尽量调大；
- topic 匹配性能问题已经在新版本中有所优化；
- queue 中留存的消息应该尽量少（理想很丰满，现实很残酷）；
- 有人说越小的消息越能降低延迟，提供消息速率；也有人说将小消息合并成大消息进行处理，对性能更好；
- 对消息进行高效的表达（更精简的数据格式 or 压缩），能够提升传输性能（业务处理复杂度略微提高）；
- HiPE 是什么鸟？仍处于实验性阶段？适用于哪些场景？
- 事务肯定要禁止使用，可靠性通过其他方案解决；
- immediate 已经被取消，存在其他等效方案；mandatory 适用于某些特定场景，是否使用应该视情况而定；
- 避免 HA ？扯什么蛋！cluster 也会对性能有一定影响？嗯，这是一句没有错误的废话；
- 在多核系统上跑 RabbitMQ 性能会比泡在单核上强（需要注意下 docker 等虚拟环境等使用情况）；
- Virtualisation 也会性能有一定影响，如果这里是指通过 management 插件的 Virtualisation ，那么性能问题已经被我遇到过了！

----------

You will increase the throughput with a **larger prefetch count** AND at the same time **ACK multiple messages** (instead of sending ACK for each message) from your consumer.

But, of course, ACK with multiple flag on (http://www.rabbitmq.com/amqp-0-9-1-reference.html#basic.ack) requires extra logic on your consumer application (http://lists.rabbitmq.com/pipermail/rabbitmq-discuss/2013-August/029600.html). You will have to keep a list of delivery-tags of the messages delivered from the broker, their status (whether your application has handled them or not) and ACK every N-th delivery-tag (NDTAG) when all of the messages with delivery-tag less than or equal to NDTAG have been handled.

更大的 prefetch 和 ack 多条消息的组合必然可以提高吞吐量，但会增加一定的业务处理复杂度；

