
# 内存使用

官网地址：[这里](http://www.rabbitmq.com/memory-use.html)

RabbitMQ 能够报告自身的内存使用情况，以便用户获知系统的哪部分正在使用内存；

需要注意的是，所有的指标值都是基于底层 Erlang VM 返回的大概值；但是，你仍旧应该认为其足够准确可用；

你可以通过 `rabbitmqctl status` 命令获取内存使用情况报告；或者通过 management 插件 Web UI 中的 node 详情页进行查看；内存使用情况被划分为下面几类（相互直接没有覆盖）：

Connection 相关内存占用
这部分内存由外部程序向 RabbitMQ 创建的 connection 和 channel 构成；在使能了某些插件的情况下，一些由 RabbitMQ 向外部业务发起的 connection 和 channel 占用的内存也将计算在内；若启用了 SSL ，那么也包括这部分使用的内存；

Queue 相关内存占用
由每条 queue 进程占用的内存构成；需要注意的是，queue 会在达到设置的内存压力阈值时，将其持有的消息内容 swap 到磁盘上；消息的 body 不被计算在这里，而是会被算进 Binary 内存占用统计中；

Plugin 相关内存占用
被插件所占用的内存；这部分统计值不包括 Erlang client 所使用的部分（已被统计到 Connection 占用内存中），也不包括 management 插件数据库所使用的部分（该内存单独统计）；
注意：RabbitMQ 会针对协议插件，如 STOMP 和 MQTT 统计每条 connection 使用的内存；

其他进程占用的内存
Memory belonging to processes not counted above, and memory assigned to "processes" by the Erlang VM, but not to any one process. Memory that has recently been garbage collected can show up here briefly.


Mnesia
Mnesia keeps an in-memory copy of all its data (even on disc nodes). Typically this will only be large when there are a large number of queues, exchanges, bindings, users or virtual hosts.


Message store index
The default message store implementation keeps an in-memory index of all messages, including those paged out to disc.

Management database
The management database (if the management plugin is loaded). In a cluster, this will only be present on one node.

Other ETS tables
Other in-memory tables besides the three sets above. Note that due to a bug in current versions of the Erlang runtime, some memory will be counted under this heading for all tables, including the three sets above.

Binaries
Memory used by shared binary data in the Erlang VM. In-memory message bodies show up here.

Code
Memory used by code. Should be fairly constant.

Atoms
Memory used by atoms. Should be fairly constant.

Other system memory
Other memory used by Erlang. One contributor to this value is the number of available file descriptors.