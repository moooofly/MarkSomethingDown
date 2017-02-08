
> 说明：RabbitMQ 文档的变化还是比较大的，四年前曾真对**确认机制**进行过[翻译](https://my.oschina.net/moooofly/blog/142095)，这两天旧事重提，发现又有了新内容可补充；
>
> 官网原文：[这里](http://www.rabbitmq.com/confirms.html)

# Consumer Acknowledgements and Publisher Confirms

## Introduction

> Systems that use a messaging broker such as RabbitMQ are by definition distributed. Since protocol methods (messages) sent are not guaranteed to reach the peer or be successfully processed by it, both publishers and consumers need a mechanism for delivery and processing confirmation. Several messaging protocols supported by RabbitMQ provide such features. This guide covers the features in AMQP 0-9-1 but the idea is largely the same in other protocols (STOMP, MQTT, et cetera).

使用诸如 RabbitMQ 这类消息中间件的系统，根据定义应该满足分布式要求；由于被发送的协议方法（消息）无法被保证到达指定 peer ，或到达指定 peer 后一定被成功处理，因此，无论是 publishers 还是 consumers 都需要一种机制来确保 delivery 和 processing 过程被确认；RabbitMQ 所提供的多种消息协议都提供了该特性；本文主要讲述 AMQP 0-9-1 中的特性，但相应内容在其他协议中也是一样的（STOMP, MQTT 等）

> Delivery processing acknowledgements from consumers to RabbitMQ are known as acknowledgements in AMQP 0-9-1 parlance; broker acknowledgements to publishers are a protocol extension called publisher confirms.

由 consumers 发送到 RabbitMQ 的 Delivery 处理确认，在 AMQP 0-9-1 术语中通常被称作 **acknowledgements**；而由 broker 发送给 publishers 的确认实际上是一种协议扩展，称作 **publisher confirms** ；

## (Consumer) Delivery Acknowledgements

> When RabbitMQ delivers a message to a consumer, it needs to know when to consider the message successfully sent. What kind of logic is optimal depends on the system. It is therefore primarily an application decision. In AMQP 0-9-1 it is made when a consumer is registered using the basic.consume method or a message is fetched on demand with the basic.get method.

当 RabbitMQ 将消息 deliver 给 consumer 时，是需要知道如何判定消息成功发送的；而哪种判定逻辑是合适的选择是取决于系统要求的，实际来讲，主要由业务使用决定；在 AMQP 0-9-1 中，当 consumer 使用 `basic.consume` 方法或 `basic.get` 方法进行“注册”时，实际上已经做出了业务选择了；

> If you prefer a more example-oriented and step-by-step material, consumer acknowledgements are also covered in RabbitMQ tutorial #2.

在 [RabbitMQ tutorial #2](http://www.rabbitmq.com/getstarted.html) 中由具体的例子可以参考；

### Delivery Identifiers: Delivery Tags

> Before we proceed to discuss other topics it is important to explain how deliveries are identified (and acknowledgements indicate their respective deliveries). When a consumer (subscription) is registered, messages will be delivered (pushed) by RabbitMQ using the basic.deliver method. The method carries a delivery tag, which uniquely identifies the delivery on a channel. Delivery tags are therefore scoped per channel.

首先，需要理解一个重要的概念：deliveries 是如何被确认的（以及 acknowledgements 表明 their respective deliveries）；当一个 consumer 进行注册后，消息会被 RabbitMQ 通过 `basic.deliver` 方法源源不断的推给（deliver）该 consumer ；该方法中会携带一个 **delivery tag** 字段，用于在当前 channel 上唯一确定当前 delivery ；即 **Delivery tags 的作用域为当前 channel** ；

> Delivery tags are monotonically growing positive integers and are presented as such by client libraries. Client library methods that acknowledge deliveries take a delivery tag as an argument.

Delivery tags 是单调递增的正整数，并且在 client 库中也是这么处理的；client 库中的方法在对 deliveries 进行确认时，同样会将 delivery tag 的内容作为参数；

### Acknowledgement Modes

> Depending on the acknowledgement mode used, RabbitMQ can consider a message to be successfully delivered either immediately after it is sent out (written to a TCP socket) or when an explicit ("manual") client acknowledgement is received. Manually sent acknowledgements can be positive or negative and use one of the following protocol methods:
> 
> - `basic.ack` is used for positive acknowledgements
> - `basic.nack` is used for negative acknowledgements (note: this is a RabbitMQ extension to AMQP 0-9-1)
> - `basic.reject` is used for negative acknowledgements but has one limitations compared to basic.nack

取决于使用了何种 acknowledgement 模式，RabbitMQ 可能会在消息一发送（写入 TCP socket）就认为成功 delivery ，也可能在收到显式（“手动”）client 确认时认为成功 delivery ；当采取手动确认时，则可以采用如下方法进行成功确认或失败确认；

- `basic.ack` 用于成功确认；
- `basic.nack` 用于失败确认（注意：其为 RabbitMQ 针对 AMQP 0-9-1 进行的扩展）；
- `basic.reject` 用于失败确认，但和 `basic.nack` 比具有一个限制条件；

> Positive acknowledgements simply instruct RabbitMQ to record a message as delivered. Negative acknowledgements with basic.reject have the same effect. The difference is primarily in the semantics: positive acknowledgements assume a message was successfully processed while their negative counterpart suggests that a delivery wasn't processed but still should be deleted.

`basic.ack` 仅用于通知 RabbitMQ 将对应消息记录为 delivery 成功；
`basic.reject` 具有和 `basic.ack` 相同的效果；
两者的差别仅在于语义：成功确认认为消息被成功处理；而失败确认则表明当前 delivery 未被处理，但仍需进行 delete 操作（针对消息）；

### Acknowledging Multiple Deliveries at Once

> Manual acknowledgements can be batched to reduce network traffic. This is done by setting the `multiple` field of acknowledgement methods (see above) to true. Note that `basic.reject` doesn't historically have the field and that's why `basic.nack` was introduced by RabbitMQ as a protocol extension.

手动确认支持**批处理**操作以减少网络通信量；可以通过在上述确认方法中设置 `multiple` 的值为 true 实现；需要注意的是，由于历史原因，`basic.reject` 方法中未提供 `multiple` 字段，这也是为何 RabbitMQ 要引入 `basic.nack` 的原因；

> When the `multiple` field is set to , RabbitMQ will acknowledge all outstanding delivery tags up to and including the tag specified in the acknowledgement. Like everything else related to acknowledgements, this is scoped per channel. For example, given that there are delivery tags 5, 6, 7, and 8 unacknowledged on channel Ch, when an acknowledgement frame arrives on that channel with delivery_tag set to 8 and `multiple` set to true, all tags from 5 to 8 will be acknowledged. If `multiple` was set to false, deliveries 5, 6, and 7 would still be unacknowledged.

当 `multiple` 域被设置后，RabbitMQ 将会根据 acknowledgement 消息中指定的 tag 值，直接对尚未确认的全部 delivery tags 进行确认；正如其他和 acknowledgement 相关的行为一样，该操作的作用域同样是 channel 级别的；例如，在 channel Ch 上存在 delivery tags 为 5, 6, 7 和 8 的消息尚未被确认；当 delivery_tag 为 8 ，并且 `multiple` 为 true 的 acknowledgement 消息出现在该 channel 上时，将直接将 5 到 8 的 tags 全部确认；若 `multiple` 设置为 false ，则 5 到 7 的 deliveries 则仍需被确认；

### Channel Prefetch Setting (QoS)

> Because messages are sent (pushed) to clients asynchronously, there is usually more than one message "in flight" on a channel at any given moment. In addition, manual acknowledgements from clients are also inherently asynchronous in nature. So there's a sliding window of delivery tags that are unacknowledged. Developers would often prefer to cap the size of this window to avoid the unbounded buffer problem on the consumer end. This is done by setting a "prefetch count" value using the `basic.qos` method. The value defines the max number of unacknowledged deliveries that are permitted on a channel. Once the number reaches the configured count, RabbitMQ will stop delivering more messages on the channel unless at least one of the outstanding ones is acknowledged.

由于消息是异步推给 client 的，因此在任意给定时刻，都可能存在不止一条处于 "in flight" 状态的消息；除此之外，client 侧的手动确认机制本质上同样是异步的；因此，总是会存在一个针对 delivery tags 的滑动窗口；开发者通常都需要对该窗口进行最大值限制，以免造成 consumer 侧占用无限大的缓存浪费；可以通过 `basic.qos` 方法设置 "prefetch count" 的值以达到此目的；该值定义的是允许出现在 channel 上的、未确认 deliveries 的最大数目；一旦达到该配置值，RabbitMQ 将会停止当前 channel 上的消息投递，除非至少一条消息被确认；

> For example, given that there are delivery tags 5, 6, 7, and 8 unacknowledged on channel Ch and channel Ch's prefetch count is set to 4, RabbitMQ will not push any more deliveries on Ch unless at least one of the outstanding deliveries is acknowledged. When an acknowledgement frame arrives on that channel with delivery_tag set to 8, RabbitMQ will notice and deliver one more message.

举例说明；

> It's worth reiterating that the flow of deliveries and manual client acknowledgements is entirely asynchronous. Therefore if prefetch value is changed while there already are deliveries in flight, a natural race condition arises and there can temporarily be more than prefetch count unacknowledged messages on a channel.

需要再次重申：deliveries 和 client 手动确认都是完全异步的；因此，如果 prefetch 的值在 deliveries 处于 "in flight" 状态时发生了变化，则会产生竞争条件，即可能临时性出现超过 prefetch 数量的未确认消息；

### Client Errors: Double Acking and Unknown Tags

> Should a client acknowledge the same delivery tag more than once, RabbitMQ will result a channel error such as `PRECONDITION_FAILED - unknown delivery tag 100`. The same channel exception will be thrown if an unknown delivery tag is used.

如果 client 针对相同 delivery tag 确认了不止一次，RabbitMQ 将会产生诸如 `PRECONDITION_FAILED - unknown delivery tag 100` 的 channel error ；类似地，如果遇到了未知 delivery tag ，则也会产生该错误信息；

## Publisher Confirms

> Using standard AMQP 0-9-1, the only way to guarantee that a message isn't lost is by using transactions -- make the channel transactional, publish the message, commit. In this case, transactions are unnecessarily heavyweight and decrease throughput by a factor of 250. To remedy this, a confirmation mechanism was introduced. It mimics the consumer acknowledgements mechanism already present in the protocol.

如果采用标准的 AMQP 0-9-1 协议，则唯一能够保证消息不会丢失的方式是利用**事务机制** -- 令 channel 处于 transactional 模式、向其 publish 消息、执行 commit 动作。在这种方式下，事务机制会带来大量的多余开销，并会导致吞吐量下降 **250%** 。为了补救事务带来的问题，引入了 confirmation 机制（即 Publisher Confirm）。其处理方式正是模仿的上述 consumer acknowledgements 机制；

> To enable confirms, a client sends the `confirm.select` method. Depending on whether `no-wait` was set or not, the broker may respond with a `confirm.select-ok`. Once the `confirm.select` method is used on a channel, it is said to be in confirm mode. A transactional channel cannot be put into confirm mode and once a channel is in confirm mode, it cannot be made transactional.

为了使能 confirm 机制，client 首先要发送 `confirm.select` 方法帧。取决于是否设置了 `no-wait` 属性，broker 会相应的判定是否以 `confirm.select-ok` 进行应答。一旦在 channel 上使用 `confirm.select` 方法，channel 就将处于 confirm 模式。处于 transactional 模式的 channel 不能再被设置成 confirm 模式，反之亦然。

> Once a channel is in confirm mode, both the broker and the client count messages (counting starts at 1 on the first `confirm.select`). The broker then confirms messages as it handles them by sending a `basic.ack` on the same channel. The `delivery-tag` field contains the sequence number of the confirmed message. The broker may also set the `multiple` field in `basic.ack` to indicate that all messages up to and including the one with the sequence number have been handled.

一旦 channel 处于 confirm 模式，broker 和 client （译者注：client 的计数自行实现）都将启动消息计数（以 `confirm.select` 为基础从 1 开始计数）。broker 会在处理完消息后（译者注：这里的说法会让人产生错误理解，何为处理完消息？后续还有涉及），在当前 channel 上通过发送 `basic.ack` 的方式对其（消息）进行 confirm 。`delivery-tag` 域的值标识了被 confirm 消息的序列号。broker 也可以通过设置 `basic.ack` 中的 `multiple` 域来表明到指定序列号为止的所有消息都已被 broker 正确的处理了。

> An example in Java that publishes a large number of messages to a channel in confirm mode and waits for the acknowledgements can be found here.

一个 Java 示例展现了 publish 大量消息到一个处于 confirm 模式的 channel 并等待获取 acknowledgement 的情况，示例在[这里](http://hg.rabbitmq.com/rabbitmq-java-client/file/default/test/src/com/rabbitmq/examples/ConfirmDontLoseMessages.java)。

### Negative Acknowledgment

> In exceptional cases when the broker is unable to handle messages successfully, instead of a `basic.ack`, the broker will send a `basic.nack`. In this context, fields of the `basic.nack` have the same meaning as the corresponding ones in `basic.ack` and the `requeue` field should be ignored. By nack'ing one or more messages, the broker indicates that it was unable to process the messages and refuses responsibility for them; at that point, the client may choose to **re-publish** the messages.

在异常情况发生时，broker 将无法成功处理相应的消息，此时 broker 将发送 `basic.nack` 来代替 `basic.ack` 。在这个情形下，`basic.nack` 中各域值的含义与  `basic.ack` 中相应各域含义是相同的，同时 `requeue` 域的值应该被忽略。 通过 nack 一条或多条消息 ， broker 表明自身无法对相应消息完成处理，并拒绝为这些消息的处理负责。在这种情况下，client 可以选择将消息 **re-publish** 。

> After a channel is put into confirm mode, all subsequently published messages will be confirmed or nack'd once. No guarantees are made as to how soon a message is confirmed. No message will be both confirmed and nack'd.

在 channel 被设置成 confirm 模式之后，所有被 publish 的后续消息都将被 confirm（即 ack） 或者被 nack 一次。但是**没有对消息被 confirm 的快慢做任何保证**，并且同一条消息不会既被 confirm 又被 nack 。

> `basic.nack` will only be delivered if an internal error occurs in the Erlang process responsible for a queue.

`basic.nack` 只会在负责 queue 功能的 **Erlang 进程发生内部错误时**被发送。

### When will messages be confirmed?

> For unroutable messages, the broker will issue a confirm once the exchange verifies a message won't route to any queue (returns an empty list of queues). If the message is also published as `mandatory`, the `basic.return` is sent to the client before `basic.ack`. The same is true for negative acknowledgements (`basic.nack`).

对于无法路由的消息，broker 会在确认了通过 exchange 无法将消息路由到任何 queue 后，发送回客户端 `basic.ack` 进行确认（其中包含空的 queue 列表）。如果客户端发送消息时使用了 `mandatory` 属性，则会发送回客户端 `basic.return` + `basic.ack` 信息。上述行为对于 `basic.nack` 是一样的；

> For routable messages, the `basic.ack` is sent when a message has been accepted by all the queues. For persistent messages routed to durable queues, this means persisting to disk. For mirrored queues, this means that all mirrors have accepted the message.

对于能够进行路由的消息，broker 会在消息被所有 queue “接受”后，发送回客户端 `basic.ack` 进行确认。**对于 persistent 消息路由到 durable queue 的情况**，意味着持久化到硬盘动作的完成。**对于镜像队列而言**，意味着所有镜像队列都“接受”了该消息之后。

### Ack Latency for Persistent Messages

> basic.ack for a persistent message routed to a durable queue will be sent after persisting the message to disk. The RabbitMQ message store persists messages to disk in batches after an interval (a few hundred milliseconds) to minimise the number of fsync(2) calls, or when a queue is idle. This means that under a constant load, latency for basic.ack can reach a few hundred milliseconds. To improve throughput, applications are strongly advised to process acknowledgements asynchronously (as a stream) or publish batches of messages and wait for outstanding confirms. The exact API for this varies between client libraries.

针对 persistent 消息路由到 durable queue 的情况，`basic.ack` 是在消息被持久化到 disk 之后发送的；RabbitMQ 的消息存储模块，或者按照一定的时间间隔（几百毫秒）批量持久化消息的（为了最小化 `fsync(2)` 的调用次数），或者当 queue 空闲时处理；这就意味着，在恒定的负载压力下，针对 `basic.ack` 的延迟可达几百毫秒；为了改进吞吐量，强烈建议应用要异步处理 acknowledgements 消息（按照流的方式处理），或者批量 publish 消息（消息打包），之后再等待 confirms ；需要注意的时，针对该功能的具体 API 在不同 client 库中的实现会有所不同；

### Confirms and Guaranteed Delivery

> The broker loses persistent messages if it crashes before said messages are written to disk. Under certain conditions, this causes the broker to behave in surprising ways.

如果在将消息写入磁盘前（瞬间） broker 发生异常，则 broker （理论上）会丢失持久化消息。在特定条件下，还会导致 broker 运行不正常。

> For instance, consider this scenario:
>
> - a client publishes a persistent message to a durable queue
> - a client consumes the message from the queue (noting that the message is persistent and the queue durable), but doesn't yet ack it,
> - the broker dies and is restarted, and
> - the client reconnects and starts consuming messages.

例如，考虑下述情景：

- 一个 client 将 persistent 消息发送到 durable queue 时；
- 一个 client 从 queue 中 consume 消息后（注意：要求消息具有 persistent 属性，queue 具有 durable 属性），尚未对其进行 ack 时；
- broker 异常重启时；
- client 重连后，重新 consume 消息时；

> At this point, the client could reasonably assume that the message will be delivered again. This is not the case: the restart has caused the broker to lose the message. In order to guarantee persistence, a client should use confirms. If the publisher's channel had been in confirm mode, the publisher would not have received an ack for the lost message (since the message hadn't been written to disk yet).

在上述情景中，client 可能理所当然的认为消息会被（broker）重新 deliver 。但这并非事实：重启（有可能）会令 broker 丢失消息。为了确保持久性，client 应该使用 confirm 机制。如果 publisher 使用的 channel 被设置为 confirm 模式，publisher 将不会收到已丢失消息的 ack（这是因为 consumer 没有对消息进行 ack ，同时该消息尚未被写入磁盘）。

## Limitations

### Maximum Delivery Tag

> Delivery tag is a 64 bit long value, and thus its maximum value is 9223372036854775807. Since delivery tags are scoped per channel, it is very unlikely that a publisher or consumer will run over this value in practice.

Delivery tag 的值以 64 bit 数字表示，因此其最大值为 9223372036854775807 ；由于 delivery tags 是按每个 channel 计算的，因此在实际使用过程中不太可能遇到越界问题；



