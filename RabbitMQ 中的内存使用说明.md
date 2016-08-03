
# 内存使用

RabbitMQ 能够报告自身的内存使用情况，以便用户获知系统的哪部分正在使用内存；

需要注意的是，所有的指标值都是基于底层 Erlang VM 返回的大概值；但是你仍旧应该认为其足够准确；

你可以通过 `rabbitmqctl status` 命令获取内存使用情况报告；或者通过 management 插件 Web UI 中的 node 详情页进行查看；内存使用情况被划分为下面几类（相互直接没有覆盖）：

## Connection 相关内存占用
这部分内存由外部程序向 RabbitMQ 创建的 connection 和 channel 构成；在使能了某些插件的情况下，一些由 RabbitMQ 向外部业务发起的 connection 和 channel 占用的内存也将计算在内；若启用了 SSL ，那么也包括这部分使用的内存；

## Queue 相关内存占用
由每条 queue 进程占用的内存构成；需要注意的是，queue 会在达到设置的内存压力阈值时，将其持有的消息内容 swap 到磁盘上；消息的 body 不被计算在这里，而是会被算进 Binary 内存占用统计中；

## Plugin 相关内存占用
被插件所占用的内存；这部分统计值不包括 Erlang client 所使用的部分（已被统计到 Connection 占用内存中），也不包括 management 插件数据库所使用的部分（该内存单独统计）；
注意：RabbitMQ 会针对协议插件，如 STOMP 和 MQTT 统计每条 connection 使用的内存；

## 其他进程占用的内存
Memory belonging to processes not counted above, and memory assigned to "processes" by the Erlang VM, but not to any one process. Memory that has recently been garbage collected can show up here briefly.

## Mnesia
mnesia 数据库会维护所有数据的一份内存拷贝（即使当前节点为磁盘节点）；典型情况下，只有当存在大量 queues, exchanges, bindings, users 或 virtual hosts 时才会占用很大内存；

## Message store index
默认的消息存储实现会在内存中维护一份针对所有消息的索引，包括那些已经被 page out 到磁盘上的消息的索引；

## Management 插件的统计数据库占用的内存
在使能了 management 插件的情况下，其使用的统计数据库占用的内存；
在 cluster 中，该数据库仅会出现某一个节点上；

## 其他 ETS 表
除了上面三种集合之外的其他内存表占用的内存；
需要注意的是，由于当前 Erlang 运行时版本中的 bug ，某些内存占用会被统计到这里，但同时也统计到其他表自身内存占用中；

## Binary 数据
在 Erlang VM 中由共享 binary 数据使用的内存；In-memory message bodies show up here.

## Code
Memory used by code. Should be fairly constant.
代码本身占用的内存；该值基本不会有啥变化；

## Atoms
Memory used by atoms. Should be fairly constant.
由 atom 占用的内存；该值基本不会有啥变化；

## 其他系统内存
被 Erlang 自身使用的其他内存；其中的一种就是可用文件描述符数目；


----------


官网地址：[这里](http://www.rabbitmq.com/memory-use.html)