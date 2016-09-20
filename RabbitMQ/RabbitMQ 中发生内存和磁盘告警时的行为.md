


# Memory and Disk Alarms

在两种场景下，RabbitMQ 会为了避免自身崩溃而停止从客户端连接上进行消息读取：
- 当**内存使用**超过配置阈值上限时；
- 当**磁盘空闲空间**低于配置阈值下限时；

两种场景下，RabbitMQ 都会临时性的阻塞所有连接，即停止从客户端（发布消息的）连接上进行读取行为；针对每条连接的心跳监控功能也会同时被去使能；通过 `rabbitmqctl` 和 management 插件看到的所有（被 Producer 使用的）网络连接状态，或者变为 `blocking`（表明当前并未在这些连接上发布消息，而此时可以继续操作）或者变为 `blocked`（表明已经有消息发布在这些连接上，而此时处于 pause 状态）

当 RabbitMQ 运行于 cluster 中时，内存和磁盘告警的影响都是 cluster 范围的；如果 cluster 中的某个节点达到了阈值限制，则所有节点都会触发连接的阻塞；

这里的实现意图是：只停止  producer 但允许 consumer 不受影响的继续运行；然而，因为 AMQP 协议允许多个 producers 和 consumers 在同一个 channel 上工作；也允许在同一个连接上的不同 channel 中工作，因此上面的逻辑实现并不完美；

在实际场景中，上述策略对于大多数应用来说都不会导致任何问题，因为此时只会看到截流导致的延迟；无论如何，若设计考量上允许，基于独立的 AMQP 连接处理 producing 或 consuming 是明智的选择；


> 相关阅读：
> - [内存阈值工作方式](https://github.com/moooofly/MarkSomethingDown/blob/master/RabbitMQ/RabbitMQ%20%E4%B8%AD%E7%9A%84%E5%86%85%E5%AD%98%E5%91%8A%E8%AD%A6%E9%97%AE%E9%A2%98.md)
> - [磁盘阈值工作方式](https://github.com/moooofly/MarkSomethingDown/blob/master/RabbitMQ/RabbitMQ%20%E4%B8%AD%E7%9A%84%E7%A3%81%E7%9B%98%E5%91%8A%E8%AD%A6%E9%97%AE%E9%A2%98.md)
> - [客户端如何确定自己被阻塞了](https://github.com/moooofly/MarkSomethingDown/blob/master/RabbitMQ/RabbitMQ%20%E4%B9%8B%E8%BF%9E%E6%8E%A5%E9%98%BB%E5%A1%9E%E9%80%9A%E7%9F%A5%E5%8A%9F%E8%83%BD.md)

## Related concepts

当 RabbitMQ 的 fd 使用量接近了操作系统允许的限制值时，则会直接拒绝掉客户端的连接建立；    
当客户端发布消息的速度快于 RabbitMQ 的处理能力时，则会进入[流控](http://www.rabbitmq.com/flow-control.html)处理过程；    

----------

官网原文：[这里](http://www.rabbitmq.com/alarms.html)

