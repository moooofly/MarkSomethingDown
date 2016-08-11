


# Blocked Connection Notifications

在某些情况下，若 broker 触发资源阈值限制（内存或磁盘）时，客户端能够接收到连接阻塞通知信息，将会非常有必要；

这里我们会介绍一种 AMQP 协议扩展，基于该扩展 broker 就能够在发生连接阻塞时，发送 `connection.blocked` 方法给客户端；在阻塞解除时，发送 `connection.unblocked` ；

为了接收（并正确处理）上述通知（方法），客户端必须在自身的 `client-properties` 能力集中表明支持 `connection.blocked` ；官方支持的客户端已经默认支持了上述方法，并提供了注册的 handler 方法，用于处理 connection.blocked 和 connection.unblocked ；

# Using Blocked Connection Notifications with Java Client

lve

# Using Blocked Connection Notifications with .NET Client

lve

# When Notifications are Sent

`connection.blocked` 方法的发送时机为 RabbitMQ 首次触发资源阈值时；例如，当 RabbitMQ 节点检测到当前 RAM 已经达到阈值时，则会发送 `connection.blocked` 方法给所有已连接的、支持该特性的、作为 Producer 的客户端；如果在连接重回非阻塞状态前，当前节点又触发了磁盘空闲空间不足告警，则不会发送额外的 `connection.blocked` ；

`connection.unblocked` 方法的发送时机为所有资源告警都被恢复时，之后连接将恢复到完全非阻塞状态；


----------

官网原文：[这里](http://www.rabbitmq.com/connection-blocked.html)

