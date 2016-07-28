

当业务和 RabbitMQ 的消息交互量大到一定程度时，RabbitMQ 的 Web 管理页面 Overview 标签中会出现如下告警信息：

![RabbitMQ managerment 插件告警](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/managerment%20statistics%20database%20%E8%BF%87%E8%BD%BD%E9%97%AE%E9%A2%98.png "RabbitMQ managerment 插件告警")

（据说最严重的情况下，积压了几十万条消息）

上述内容提供了以下几点信息：
- managerment 插件通过一个名为 statistics 的数据库维护用于 web 页面展示的相关统计数据；
- managerment 插件使用了内部 queue 有序处理消息，随着 queue 中消息的增多，势必造成内存消耗的增加，统计信息的即时性变差，甚至可能对磁盘 I/O 造成影响（待确认）；
- 设置 rates_mode 选项参数的值为 node 可能有所改善；



在 overview.ejs 中

```ejs
<div class="updatable">
<% if (overview.statistics_db_event_queue > 1000) { %>
<p class="warning">
  The management statistics database currently has a queue
  of <b><%= overview.statistics_db_event_queue %></b> events to
  process. If this number keeps increasing, so will the memory used by
  the management plugin.

  <% if (overview.rates_mode != 'none') { %>
  You may find it useful to set the <code>rates_mode</code> config item
  to <code>none</code>.
  <% } %>
</p>
<% } %>
</div>
```






刚才反馈的 publish 等曲线掉底的问题，经过 @张斌 确认，结论如下：
1.在掉底曲线的时间段内，rabbitmq 的统计信息数据库（或队列）积压了几十万条统计信息；
2.同一时间段内，业务获取 channel 超时飙高；
3.张斌重启 rabbitmq 的 managerment 管理插件后（等于清空积压的统计信息），统计数据库从 xg-napos-rmq-1 节点随机迁移到 xg-napos-rmq-3 节点，此时发现，整体 qps 从原来的 2000 上升到 4000 左右；此时业务获取 channel 超时时间恢复正常；


所以，建议将 rabbitmq managerment 插件所使用的统计数据库部署到单独一个节点上，避免对业务造成影响；应该可以立刻取的改善；


之后我会深入研究下 rabbitmq managerment 插件的使用和调优姿势，看看内否进一步改进


$sudo rabbitmqctl eval 'application:stop(rabbitmq_management), application:start(rabbitmq_management).'
or
$sudo rabbitmqctl eval 'exit(erlang:whereis(rabbit_mgmt_db), please_terminate).'
1 Comment
帮忙看下这两条命令的区别在哪