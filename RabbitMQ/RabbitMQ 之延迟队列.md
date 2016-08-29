


# Lazy Queues

从 RabbitMQ 3.6.0 版本开始，broker 中增加了 **Lazy Queues** 这个东东：这种 queue 基于 disk 保存尽量多的消息，并且尽在 consumers 进行消息请求时才将对应的内容加载到 RAM 中；这也就是 lazy 的真谛；

lazy queues 存在的主要目标就是为了支持超长 queues 的存在（100w+ 消息）；这种 queues 通常出现在 consumers 在较长时间内都无法从 queue 中取走消息的场景；问题出现原因是多样的：比如 consumers 处于离线状态；consumers 发生 crash ；或者由于维护的需要将其取消了等等；

默认情况下，queues 会在内存中对 publish 到 RabbitMQ 的消息进行缓存；该缓存主要用于尽快向 consumers 投递消息（⚠️ 持久化消息会在进入 broker 后被写入磁盘，同时在缓存中保留一份）；无论何时，只要 broker 认为应该释放内存了，缓存中的消息就会被 page out 到 disk 上；将消息 Page 到 disk 到行为需要花费一定的时间，并阻塞 queue 进程本身，即在此过程中无法接收新消息；尽管我们在最新 RabbitMQ 版本中已经改进了 paging 算法，但在某些使用场景下仍旧不是很理想：例如 queue 存在 100w+ 消息需要 page out 到磁盘时；

Lazy queues 在这种场景下可以发挥作用：移除了缓存的使用，仅在 consumers 请求时加载消息到内存中；Lazy queues 会将每一条发送到当前 queue 的消息立刻写入文件系统，完全消除了前面提到的缓存的使用；好处就是极大的减少了 queue 所需的 RAM 量，同时消除了 paging 问题；尽管 lazy queue 会增加 I/O 使用，但充其量也就等同于处理持久化消息的情况；


# Using Lazy Queues

Queues can be made to work in default mode or lazy mode either by specifying the mode via queue.declare arguments, or by applying a policy in the server. In the case where both policy and queue arguments specify a queue mode, then the queue argument has priority over the policy value. This means that if a queue mode is set via a declare argument, it can only be changed by deleting the queue, and re-declaring it later with a different argument.

## Configuration using arguments

The queue mode can be set by supplying the x-queue-mode queue declaration argument with a string specifying the desired mode. Valid modes are "default" and "lazy". If no mode is specified during declare, then "default" is assumed. The default mode is the behaviour already present in pre 3.6.0 versions of the broker, so there are no breaking changes in this regard.

This example in Java declares a queue with the queue mode set to "lazy":

```
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-queue-mode", "lazy");
channel.queueDeclare("myqueue", false, false, false, args);
```

## Configuration using policy

To specify a queue mode using a policy, add the key queue-length to a policy definition. For example:

|                  |                              |                  |
 ----------------- | ---------------------------- | ------------------
| rabbitmqctl | rabbitmqctl set_policy Lazy "^lazy-queue$" '{"queue-mode":"lazy"}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Lazy "^lazy-queue$" "{""queue-mode"":""lazy""}" --apply-to queues |

This ensures the queue called lazy-queue will work in the lazy mode.

Policies can also be defined using the management plugin, see the policy documentation for more details.

## Changing queue modes

If you specified the queue mode via a policy, then you can change it at run time without the need of deleting the queue and re-declaring it with a different mode. If you want the previous lazy-queue to start working like a default queue, then you can do so by issuing the following command:

|                  |                              |                  |
 ----------------- | ---------------------------- | ------------------
| rabbitmqctl | rabbitmqctl set_policy Lazy "^lazy-queue$" '{"queue-mode":"default"}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Lazy "^lazy-queue$" "{""queue-mode"":""default""}" --apply-to queues |


# Performance Considerations for Lazy Queues

## Disk Utilization
As stated above, lazy queues will send every message to disk right as they enter the queue. This will increase I/O opts, but keep in mind that this is the same behaviour as when persistent messages are delivered to queues. Note that even if you publish transient messages, they will still be sent to disk when using lazy queues. With default queues transient messages are only sent to disk if paging requires it.

## RAM Utilization
lazy queues use much less memory than default queues. While it's hard to give numbers that make sense for every use case, here's what we found: we tried publishing 10 million messages into a queue, with no consumers online. The message body size was 1000 bytes. default queues required 1.2GB of RAM, while lazy queues only used 1.5MB of RAM.

For a default queue, it took 801 seconds to send 10MM messages, with an average sending rate of 12469 msg/s. To publish the same amount of messages into a lazy queue, the time required was 421 seconds, with an average sending rate of 23653 msg/s. The difference can be explained by the fact that from time to time, the default queue had to page messages to disk. Once we activated a consumer, the lazy queue had a RAM consumption of approximately 40MB while it was delivering messages. The message receiving rate average was 13938 msg/s for one active consumer.

You can reproduce the test with our Java library by running:

```shell
./runjava.sh com.rabbitmq.examples.PerfTest -e test -u test_queue \
-f persistent -s 1000 -x1 -y0 -C10000000
```
Note that this was a very simplistic test. Please make sure to run your own benchmarks.

Don't forget to change the queue mode between benchmarks runs.


## Converting between queue modes

If we need to convert a default queue into a lazy one, then we will suffer the same performance impact as when a queue needs to page messages to disk. When we convert a queue into a lazy one, first it will page messages to disk and then it will start accepting publishes, acks, and other commands.

When a queue goes from the lazy mode to the default one, it will perform the same process as when a queue is recovered after a server restart. A batch of 16384 messages will be loaded in the cache mentioned above.