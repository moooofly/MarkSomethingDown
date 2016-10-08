


# [RabbitMQ Sharding Plugin](https://github.com/rabbitmq/rabbitmq-sharding)

该插件为 RabbitMQ 引入了 sharded queues 概念；Sharding 机制通过 exchanges 实现，也就是说，消息通过一个定义为 sharded 的 exchange 被分区保存到 "shard" queues 中；

幕后使用的机制表明：我们会定义一个用于跨 queues 分区或 shard 消息的 exchange ；而分区操作会为你自动完成，也就是说，一旦你定义了用于 _sharded_ 功能的 exchange，相应的 queues 就会被自动创建到 cluster 中的每一个节点上，之后消息会在这些 queue 中被 sharded ；

下图描述了从 publisher 和 consumer 的角度该插件是如何工作的：

![Sharding Overview](https://raw.githubusercontent.com/rabbitmq/rabbitmq-sharding/master/docs/sharded_queues.png)

正如你从图中所见，当 producers 发布一系列消息后，这些消息会被分区到不同的 queue 中，之后我们的 consumer 可以从这些 queue 中获取到消息；换句话说，如果你有一个由 3 个 queue 构成的分区，那么你将需要至少 3 个 consumer 才能获取到所需的全部消息；

## Auto-scaling

该插件的其中一项有趣的特性为，如果你添加更多节点到 RabbitMQ cluster 中，那么该插件将会在新节点上自动创建出更多的 shards；假如你有一个由 4 个 queue 构成的 shard 位于 `node a` 中，同时 `node b` 刚刚加入了 cluster ；那么该插件将自动创建出 4 个 queue 在 `node b` 中，并将这些 queue 加入到 shard 分区中；已经投递到消息 _将不会_ 被 rebalanced ，但是新到达的消息将会分区到新 queues 中；

## Partitioning Messages

RabbitMQ 中默认提供的 exchanges 以 "all or nothing" 的模式工作；也就是说，一个 routing key 会匹配上绑定到 exchange 上的一组 queues ，而 RabbitMQ 会路由消息到相应的所有 queue 中；因此，对于该插件的工作方式来说，我们需要路由消息到负责分区消息的特定 exchange 上，以便消息_至多_被路由到一个 queue 中；

该插件提供了一种新 exchange 类型 `"x-modulus-hash"`，其基于传统的 hash 技术跨 queue 进行消息分区；

类型为 `"x-modulus-hash"` 的 exchange 将会对 routing key 进行 hash 运算，之后再通过 `Hash mod N` 来选择消息被路由到的目标 queue ，其中 N 为绑定到该 exchange 上的 queue 的数目；**该 exchange 将会完全忽略掉用于绑定 queue 到 exchange 上的  binding key  的作用**；

你也可以使用其他的具有类似行为模式的 exchanges ，例如 _一致性 hash Exchange_ 或者 _随机 Exchange_；前者具有伴随 RabbitMQ 一起发布的优势；

如果 _只是需要进行消息分区功能_，而不需要此插件提供的自动 queue 创建功能，那么你可以仅使用 [一致性 hash Exchange](https://github.com/rabbitmq/rabbitmq-consistent-hash-exchange)；

## Consuming From a Sharded [Pseudo-]Queue

尽管该插件会创建一组 "shard" queues ，但背后的想法是那些 queues 共同表现为一个大的、逻辑 queue ，供你进行消息的 consume ；跨 shard 时的消息整体顺序未进行定义；

一例胜千言：我们假设你声明来 exchange _images_ 作为 sharded exchange ；之后 RabbitMQ 会在幕后创建出多个 "shard" queues ：

 * _shard: - nodename images 1_
 * _shard: - nodename images 2_
 * _shard: - nodename images 3_
 * _shard: - nodename images 4_.

为了从 sharded queue 上进行 consume，需要使用 `basic.consume` 方法注册一个 consumer 到 `"images"` pseudo-queue 上；RabbitMQ 会在幕后“偷偷的“将 consumer 附着到 shard 上；需要注意的是，在进行消费行为之前，**consumers 不可以声明一个与 sharded pseudo-queue 同名的 queue **；

TL;DR: 如果你拥有一个 shard 叫做 _images_，那么你就可以直接从名为 _images_ 的 queue 上进行消费；

How does it work? 该插件将会在 shard 上关闭持有_最少 consumers 数量的_queue ，前提是 queue 中的内容对于你所连接的 broker 来说是属于本地的；

**注意：在 RabbitMQ 更新 queue 内部关于 consumers 的统计信息和 clients 发送 `basic.consume` 命令之间，存在一个小的 race condition ；**问题的根源在于，如果你的 client 发出了许多间隔很短的 `basic.consume` 命令的话，可能发生插件以非均匀的方式将 consumers 分配到 queues 的情况；

## 安装 ##

### RabbitMQ 3.6.0 或之后的版本

从 RabbitMQ `3.6.0` 版本开始，该插件已经被包含到 RabbitMQ 发布包中来；

可以使用如下命令使能该插件：

```bash
rabbitmq-plugins enable rabbitmq_sharding
```

你可能还想要使能一致性 hash Exchange 插件；

### 针对早期的 RabbitMQ 版本

安装相应的 .ez 文件，下载地址为：[Community Plugins archive](http://www.rabbitmq.com/community-plugins/)；

然后运行如下命令：

```bash
rabbitmq-plugins enable rabbitmq_sharding
```

你可能还想要使能一致性 hash Exchange 插件；


## 用法

一旦该插件安装成功，你就可以定义 exchange 为 sharded 了，只要建立一套 policy 用于匹配 exchange 的名字；例如，如果我们有一个名为 `shard.images` 的 exchange ，我们就可以定义如下 policy 来对其 shard ：

```bash
$CTL set_policy images-shard "^shard.images$" '{"shards-per-node": 2, "routing-key": "1234"}'
```

这将在 cluster 中的每个 node 上创建出 `2` 个 sharded queues ，并将那些 queues 通过 `"1234"` 这个 routing key 进行绑定；

### About the routing-key policy definition

在上面的例子中，我们在定义 policy 时使用了 `1234` 作为 routing key ，这意味着底层用作 sharding 功能的 exchanges 将会使用 `1234` 这个 routing key 绑定 sharded queues 到该 exchange 上；

这也意味着，对于 direct 类型的 exchange ，只有使用 routing key `1234` 进行 publish 的消息才会被路由到 sharded queues 中；

如果你决定使用 fanout 类型的 exchange 用作 sharding ，那么 `1234` 这个 routing key 尽管在绑定操作中被使用，但仍旧会被 exchange 所忽略；如果你使用了 `"x-modulus-hash"` exchange，那么使用的 routing key 同样会被忽略掉；因此，取决于你所使用的 exchange ，在进行消息路由时，与 `routing-key` 相关的 policy 定义产生的效果会有所不同；

针对 `routing-key` 的 policy 定义是可选的；

## 插件构建方式

参考 [RabbitMQ Plugin Development Guide](http://www.rabbitmq.com/plugin-development.html) 中的说明获取RabbitMQ Public Umbrella ；

切换到 umbrella 文件夹，之后运行如下命令：

```bash
make up
cd ../rabbitmq-sharding
make
```

## Plugin Status

此时此刻，应该认为该插件仍旧处于 __experimental__ 状态，以便更好的从社区获取反馈；


## Extra information ##

关于该插件如何影响消息顺序，以及一些其他细节内容可以查看 [README.extra.md](https://github.com/rabbitmq/rabbitmq-sharding/blob/master/README.extra.md) ；


----------


# [Additional information](https://github.com/rabbitmq/rabbitmq-sharding/blob/master/README.extra.md)

在这里你可以找到一些额外信息：关于该插件如何工作，以及如此工作的理由；

## 我们为何需要这个插件？

RabbitMQ 中的 queues 会默认在其首次声明的节点上建立绑定关系；这意味着即使你创建了 RabbitMQ cluster，但从某种角度来说，所有消息通信流量还是会发往 queue 所位于的节点上；而该插件所解决的就是给你提供了一个消息发送的中心点，并提供了跨多节点的 __负载均衡__ 功能（将 queues 分散到同一个 cluster 中的其他节点上）；

这种方式的优势在于，你的 consumers 获取消息的 queue 相对其所连接的节点来说是本地的；换句话说，producers 不需要关心 exchange 后面有些什么；

针对 shard queues 需要进行维护，全部由插件本身 __自动完成__；如果你向 cluster 中添加了更多的节点，则该插件会 __在那些节点上自动创建相应的 queues__；

如果你从 cluster 中移除了节点，RabbitMQ 会负责将相应信息从 queues 绑定信息列表中移除；消息丢失是可能发生的，因为存在一种竞争条件：即节点失效的过程中，有消息到达了用于 shard 的 exchange 上；如果你无法承受消息丢失的问题，那么你需要使用 [publisher confirms](http://www.rabbitmq.com/confirms.html) 机制进行处理；

## Message Ordering

针对每个 sharded queue 来说，消息顺序是能够得到保证的；这意味着，一旦一条消息进入到了 queue 中，那么对于该 queue 以及该 queue 的相应 consumer 来说，消息顺序是确定的；

如果你需要全局顺序性保证，那么请使用 [mirrored queues](http://www.rabbitmq.com/ha.html)；

## What strategy is used for picking the queue name

当你发送了 `basic.consume` 命令后，该插件将会选择拥有 _最少数量 consumers_ 的 queue ；并且该 queue 对于你的 client 所连接 broker 而言是本地的；当然，本地 sharded queue 将会是属于选中 shard 的 queue 分组的一部分； 

## Intercepted Channel Behaviour

This plugin works with the new `channel interceptors`. An interceptor basically allows a plugin to modify parts of an AMQP method. For example in this plugin case, whenever a user sends a `basic.consume`, the plugin will map the queue name sent by the user to one of the sharded queues.

Also a plugin can decide that a certain AMQP method can't be performed on a queue that's managed by the plugin. In this case declaring a queue called `my_shard` doesn't make much sense when there's actually a sharded queue by that name. In this case the plugin will return a channel error to the user.

下面给出插件能够处理的 AMQP 方法，以及相应的处理方式：

- `'basic.consume', QueueName`: 插件会将 `QueueName` 作为 shard 名字，选择具有最少数量 consumers 的 sharded queue ；
- `'basic.get', QueueName`: 插件会将 `QueueName` 作为 shard 名字，选择具有最少数量 consumers 的 sharded queue ；
- `'queue.declare', QueueName`: 插件会以 `QueueName` 作为名字在 shard 中声明 queue，因此 `queue.declare_ok` 将返回该 queue 的统计信息；
- `'queue.bind', QueueName`: 由于没有实际的名为 `QueueName` 的 queue ，因此该方法将返回一个 channel 错误信息；
- `'queue.unbind', QueueName`: 由于没有实际的名为 `QueueName` 的 queue ，因此该方法将返回一个 channel 错误信息；
- `'queue.purge', QueueName`: 由于没有实际的名为 `QueueName` 的 queue ，因此该方法将返回一个 channel 错误信息；
- `'queue.delete', QueueName`: 由于没有实际的名为 `QueueName` 的 queue ，因此该方法将返回一个 channel 错误信息；