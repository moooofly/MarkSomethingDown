


# [RabbitMQ Sharding Plugin](https://github.com/rabbitmq/rabbitmq-sharding)

该插件为 RabbitMQ 引入了 sharded queues 概念；
Sharding 处理由 exchanges 实现，也就是说，消息将被分区到 "shard" queues 中，通过一个被定义为 sharded 的 exchange 实现；

The machinery used behind the scenes implies defining an exchange that will partition, or shard messages across queues. The partitioning will be done automatically for you, i.e: once you define an exchange as _sharded_, then the supporting queues will be automatically created on every cluster node and messages will be sharded across them.

下图描述了从 publisher 和 consumer 的角度该插件是如何工作的：

![Sharding Overview](https://raw.githubusercontent.com/rabbitmq/rabbitmq-sharding/master/docs/sharded_queues.png)

正如你从图中所见，当 producers 发布一系列消息后，这些消息会被分区到不同的 queue 中，之后我们的 consumer 可以从这些 queue 中获取到消息；换句话说，如果你有一个由 3 个 queue 构成的分区，那么你将需要至少 3 个 consumer 才能获取到所需的全部消息；

## Auto-scaling

该插件的其中一项有趣的特性为，如果你添加更多节点到 RabbitMQ cluster 中，那么该插件将会在新节点上自动创建出更多的 shards；假如你有一个由 4 个 queue 构成的 shard 位于 `node a` 中，同时 `node b` 刚刚加入了 cluster ；那么该插件将自动创建出 4 个 queue 在 `node b` 中，并将这些 queue 加入到 shard 分区中；已经投递到消息 _将不会_ 被 rebalanced ，但是新到达的消息将会分区到新 queues 中；

## Partitioning Messages ##

RabbitMQ 中默认提供的 exchanges 以 "all or nothing" 的模式工作；也就是说，一个 routing key 会匹配上绑定到 exchange 上的一组 queues ，而 RabbitMQ 会路由消息 message 到相应的所有 queue 中；因此，对于该插件的工作方式来说，我们需要路由消息到能够分区消息的特定 exchange 上，以便消息_至多_被路由到一个 queue 中；

该插件提供了一种新 exchange 类型 `"x-modulus-hash"`， that will use
the traditional hashing technique applying to partition messages
across queues.

The `"x-modulus-hash"` exchange will hash the routing key used to
publish the message and then it will apply a `Hash mod N` to pick the
queue where to route the message, where N is the number of queues
bound to the exchange. **This exchange will completely ignore the
binding key used to bind the queue to the exchange**.

You could also use other exchanges that have similar behaviour like
the _Consistent Hash Exchange_ or the _Random Exchange_.  The first
one has the advantage of shipping directly with RabbitMQ.

If _just need message partitioning_ but not the automatic queue
creation provided by this plugin, then you can just use the
[Consistent Hash Exchange](https://github.com/rabbitmq/rabbitmq-consistent-hash-exchange).

## Consuming From a Sharded [Pseudo-]Queue ##

While the plugin creates a bunch of "shard" queues behind the scenes, the idea
is that those queues act like a big logical queue where you consume
messages from it. Total ordering of messages between shards is not defined.

An example should illustrate this better: let's say you declared the
exchange _images_ to be a sharded exchange. Then RabbitMQ creates
several "shard" queues behind the scenes:

 * _shard: - nodename images 1_
 * _shard: - nodename images 2_
 * _shard: - nodename images 3_
 * _shard: - nodename images 4_.

To consume from a sharded queue, register a consumer on the `"images"` pseudo-queue
using the `basic.consume` method. RabbitMQ will attach the consumer to a shard
behind the scenes. Note that **consumers must not declare a queue with the same
name as the sharded pseudo-queue prior to consuming**.

TL;DR: if you have a shard called _images_, then you can directly
consume from a queue called _images_.

How does it work? The plugin will chose the queue from the shard with
the _least amount of consumers_, provided the queue contents are local
to the broker you are connected to.

**NOTE: there's a small race condition between RabbitMQ updating the
queue's internal stats about consumers and when clients issue
`basic.consume` commands.** The problem with this is that if your
client issue many `basic.consume` commands without too much time in
between, it might happen that the plugin assigns the consumers to
queues in an uneven way.

## Installing ##

### RabbitMQ 3.6.0 or later

As of RabbitMQ `3.6.0` this plugin is included into the RabbitMQ distribution.

Enable it with the following command:

```bash
rabbitmq-plugins enable rabbitmq_sharding
```

You'd probably want to also enable the Consistent Hash Exchange
plugin, too.

### With Earlier Versions

Install the corresponding .ez files from our
[Community Plugins archive](http://www.rabbitmq.com/community-plugins/).

Then run the following command:

```bash
rabbitmq-plugins enable rabbitmq_sharding
```

You'd probably want to also enable the Consistent Hash Exchange
plugin, too.

## Usage ##

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

The `routing-key` policy definition is optional.

## Building the plugin ##

Get the RabbitMQ Public Umbrella ready as explained in the
[RabbitMQ Plugin Development Guide](http://www.rabbitmq.com/plugin-development.html).

Move to the umbrella folder an then run the following commands, to
fetch dependencies:

```bash
make up
cd ../rabbitmq-sharding
make
```

## Plugin Status ##

At the moment the plugin is __experimental__ in order to receive
feedback from the community.

## LICENSE ##

See the LICENSE file.

## Extra information ##

Some information about how the plugin affects message ordering and
some other details can be found in the file README.extra.md


----------


# Additional information #

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