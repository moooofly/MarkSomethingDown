

当业务消息量大到一定程度时，RabbitMQ 的 Web 管理页面中 Overview 标签中会出现如下告警信息：

![RabbitMQ managerment 插件告警](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/managerment%20statistics%20database%20%E8%BF%87%E8%BD%BD%E9%97%AE%E9%A2%98.png "RabbitMQ managerment 插件告警")



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