


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

正如上面所提到的，lazy queues 会将每一条进入该 queue 的消息写入磁盘；这会导致  I/O 操作的增加，但是请记住，此时的行为和具有持久化属性的消息到来时的情况是一样的；注意，即使被 publish 的消息是 transient 消息，在使用 lazy queue 时一样会被写入磁盘；而对于处于 default 模式的 queues 来说，transient 消息只会在需要进行 page out 处理时才被写入磁盘；

## RAM Utilization

lazy queues 使用的内存量远远小于 default queues ；尽管难于针对每一种场景给出具体的数字，但可以给出我们得到的一些结论：我们尝试发送了 10 million 消息到 queue 中，并且该 queue 不存在消费者；消息体大小为 1000 字节；default queues 需要占用 1.2GB 的 RAM，而 lazy queues 仅需要 1.5MB 的 RAM ；

对于 default queue 来说，其需要花费 801 秒来发送 10MM 消息，平均发送速率为 12469 msg/s ；而发送同样数目的消息到 lazy queue 时，时间消耗为 421 秒，平均发送速率为 23653 msg/s ；上述差别可以解释为：对于 default queue 来说，会时不时的发生将消息 page out 到 disk 到情况；一旦我们激活了一个 consumer ，在消息转发时 lazy queue 就会相应的产生大概 40MB 的 RAM 占用；平均消息接收速率变为 13938 msg/s ；


可以通过如下的 java 库进行上述测试的重放：

```shell
./runjava.sh com.rabbitmq.examples.PerfTest -e test -u test_queue \
-f persistent -s 1000 -x1 -y0 -C10000000
```

> ⚠️ 上述测试仅为简单测试用例，请确保使用你自己的 benchmark 进行测试；

Don't forget to change the queue mode between benchmarks runs.


## Converting between queue modes

如果我们需要将一个 default queue 转变成 lazy queue ；我们将会面临相似的性能冲击：因为 queue 将需要将消息 page 到 disk 上；当我们将某个 queue 转成 lazy 模式时，首先会进行消息 page 到 disk 到操作，之后才会开始处理 publish, ack 和其他命令；

当某个 queue 从 lazy 模式切回 default 模式，对应的处理过程就如同服务被重启后 queue 的恢复过程；会按照 16384 条消息一组的方式，批量加载到之前提到的缓存中；


----------


官网原文：[这里](http://www.rabbitmq.com/lazy-queues.html)