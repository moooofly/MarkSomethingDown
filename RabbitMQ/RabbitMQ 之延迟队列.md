


# Lazy Queues

从 RabbitMQ 3.6.0 版本开始，broker 中增加了 **Lazy Queues** 这个东东：这种 queue 基于 disk 保存尽量多的消息，并且尽在 consumers 进行消息请求时才将对应的内容加载到 RAM 中；这也就是 lazy 的真谛；

lazy queues 存在的主要目标就是为了支持超长 queues 的存在（100w+ 消息）；这种 queues 通常出现在 consumers 在较长时间内都无法从 queue 中取走消息的场景；问题出现原因是多样的：比如 consumers 处于离线状态；consumers 发生 crash ；或者由于维护的需要将其取消了等等；

默认情况下，queues 会在内存中对 publish 到 RabbitMQ 的消息进行缓存；该缓存主要用于尽快向 consumers 投递消息（⚠️ 持久化消息会在进入 broker 后被写入磁盘，同时在缓存中保留一份）；无论何时，只要 broker 认为应该释放内存了，缓存中的消息就会被 page out 到 disk 上；将消息 Page 到 disk 到行为需要花费一定的时间，并阻塞 queue 进程本身，即在此过程中无法接收新消息；尽管我们在最新 RabbitMQ 版本中已经改进了 paging 算法，但在某些使用场景下仍旧不是很理想：例如 queue 存在 100w+ 消息需要 page out 到磁盘时；

Lazy queues 在这种场景下可以发挥作用：移除了缓存的使用，仅在 consumers 请求时加载消息到内存中；Lazy queues 会将每一条发送到当前 queue 的消息立刻写入文件系统，完全消除了前面提到的缓存的使用；好处就是极大的减少了 queue 所需的 RAM 量，同时消除了 paging 问题；尽管 lazy queue 会增加 I/O 使用，但充其量也就等同于处理持久化消息的情况；


# Using Lazy Queues

通过 `queue.declare` 或 policy 设置，Queues 可以通过在 `default` 模式或 `lazy` 模式下；如果出现同时在 policy 和 queue 声明参数中设置了 queue 模式的情况，则通过 queue 声名参数进行的设置具有更高的优先级；这就意味着，如果 queue 的模式是通过声明参数进行的指定，则只能通过先删除该 queue ，再重新声明的方式进行变更；

## Configuration using arguments

queue 模式可以通过在 queue 声明中国年指定 `x-queue-mode` 参数的方式进行设置；有效模式为 "default" 和  "lazy" ；如果在声明的时候没有指定模式，则使用 "default" 模式；`default` 模式就是在 3.6.0 版本以前一直在使用的模式，因此此变更不会导致兼容问题；

在下面的例子中，通过 Java 代码声明了模式为 "lazy" 的 queue ：

```
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-queue-mode", "lazy");
channel.queueDeclare("myqueue", false, false, false, args);
```

## Configuration using policy

若想通过 policy 方式设置 queue 模式，只需将 queue-length 作为 key 添加到 policy 定义中，例如：

|                  |                              |                  |
 ----------------- | ---------------------------- | ------------------
| rabbitmqctl | rabbitmqctl set_policy Lazy "^lazy-queue$" '{"queue-mode":"lazy"}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Lazy "^lazy-queue$" "{""queue-mode"":""lazy""}" --apply-to queues |

上述设置会令名为 lazy-queue 的 queue 工作在 lazy 模式下；

Policies 也可以通过 management 插件进行定义，详情参见 [policy]() 文档；

## Changing queue modes

如果你是通过 policy 设置的 queue 模式，那么你就能够在无需删除后重建的情况下，直接运行时修改；如果你想要另之前设置的 lazy-queue 变回 default 模式，则你可以通过如下命令进行操作：

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


----------


官网原文：[这里](http://www.rabbitmq.com/lazy-queues.html)