


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

How does it work? The plugin will chose the queue from the shard with the _least amount of consumers_, provided the queue contents are local to the broker you are connected to.

**NOTE: there's a small race condition between RabbitMQ updating the
queue's internal stats about consumers and when clients issue
`basic.consume` commands.** The problem with this is that if your
client issue many `basic.consume` commands without too much time in
between, it might happen that the plugin assigns the consumers to
queues in an uneven way.

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

Once the plugin is installed you can define an exchange as sharded by
setting up a policy that matches the exchange name. For example if we
have the exchange called `shard.images`, we could define the following
policy to shard it:

```bash
$CTL set_policy images-shard "^shard.images$" '{"shards-per-node": 2, "routing-key": "1234"}'
```

This will create `2` sharded queues per node in the cluster, and will
bind those queues using the `"1234"` routing key.

### About the routing-key policy definition ###

In the example above we use the routing key `1234` when defining the
policy. This means that the underlying exchanges used for sharding
will bind the sharded queues to the exchange using the `1234` routing
key specified above. This means that for a direct exchange, _only
messages that are published with the routing key `1234` will be routed
to the sharded queues. If you decide to use a fanout exchange for
sharding, then the `1234` routing key, while used during binding, will
be ignored by the exchange. If you use the `"x-modulus-hash"`
exchange, then the routing key will be ignored as well. So depending
on the exchange you use, will be the effect the `routing-key` policy
definition has while routing messages.

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

Here you can find some extra information about how the plugin works
and the reasons for it.

## Why do we need this plugin? ##

RabbitMQ queues are bound to the node where they were first
declared. This means that even if you create a cluster of RabbitMQ
brokers, at some point all message traffic will go to the node where
the queue lives. What this plugin does is to give you a centralized
place where to send your messages, plus __load balancing__ across many
nodes, by adding queues to the other nodes in the cluster.

The advantage of this setup is that the queues from where your
consumers will get messages will be local to the node where they are
connected.  On the other hand, the producers don't need to care about
what's behind the exchange.

All the plumbing to __automatically maintain__ the shard queues is
done by the plugin. If you add more nodes to the cluster, then the
plugin will __automatically create queues in those nodes__.

If you remove nodes from the cluster then RabbitMQ will take care of
taking them out of the list of bound queues. Message loss can happen
in the case where a race occurs from a node going away and your
message arriving to the shard exchange. If you can't afford to lose a
message then you can use
[publisher confirms](http://www.rabbitmq.com/confirms.html) to prevent
message loss.

## Message Ordering ##

Message order is maintained per sharded queue, but not globally. This
means that once a message entered a queue, then for that queue and the
set of consumers attached to the queue, ordering will be preserved.

If you need global ordering then stick with
[mirrored queues](http://www.rabbitmq.com/ha.html).

## What strategy is used for picking the queue name ##

When you issue a `basic.consume`, the plugin will choose the queue
with the _least amount of consumers_.  The queue will be local to the
broker your client is connected to. Of course the local sharded queue
will be part of the set of queues that belong to the chosen shard.

## Intercepted Channel Behaviour ##

This plugin works with the new `channel interceptors`. An interceptor
basically allows a plugin to modify parts of an AMQP method. For
example in this plugin case, whenever a user sends a `basic.consume`,
the plugin will map the queue name sent by the user to one of the
sharded queues.

Also a plugin can decide that a certain AMQP method can't be performed
on a queue that's managed by the plugin. In this case declaring a queue
called `my_shard` doesn't make much sense when there's actually a
sharded queue by that name. In this case the plugin will return a
channel error to the user.

These are the AMQP methods intercepted by the plugin, and the
respective behaviour:

- `'basic.consume', QueueName`: The plugin will pick the sharded queue
  with the least amount of consumers from the `QueueName` shard.
- `'basic.get', QueueName`: The plugin will pick the sharded queue
  with the least amount of consumers from the `QueueName` shard.
- `'queue.declare', QueueName`: The plugin rewrites `QueueName` to be
  the first queue in the shard, so `queue.declare_ok` returns the stats
  for that queue.
- `'queue.bind', QueueName`: since there isn't an actual `QueueName`
  queue, this method returns a channel error.
- `'queue.unbind', QueueName`: since there isn't an actual `QueueName`
  queue, this method returns a channel error.
- `'queue.purge', QueueName`: since there isn't an actual `QueueName`
  queue, this method returns a channel error.
- `'queue.delete', QueueName`: since there isn't an actual `QueueName`
  queue, this method returns a channel error.