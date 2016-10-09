




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

----------

You will increase the throughput with a **larger prefetch count** AND at the same time **ACK multiple messages** (instead of sending ACK for each message) from your consumer.

But, of course, ACK with multiple flag on (http://www.rabbitmq.com/amqp-0-9-1-reference.html#basic.ack) requires extra logic on your consumer application (http://lists.rabbitmq.com/pipermail/rabbitmq-discuss/2013-August/029600.html). You will have to keep a list of delivery-tags of the messages delivered from the broker, their status (whether your application has handled them or not) and ACK every N-th delivery-tag (NDTAG) when all of the messages with delivery-tag less than or equal to NDTAG have been handled.



