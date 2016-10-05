



# [Queue Length Limit](http://www.rabbitmq.com/maxlength.html)


queue 的最大长度可以基于**消息数量**进行限制，也可以基于**总字节数**进行限制（由所有消息 body 的长度之和决定，忽略消息属性占用的长度和任何其他额外开销），或者基于两种方式一起进行限制；

对于任何给定的 queue ，其最大长度（无论任何类型）都可以由 client 使用 queue 参数进行定义，或者在 server 侧通过 [policies](http://www.rabbitmq.com/parameters.html#policies) 进行定义；当通过两种方式同时设置了最大长度时时，取其中的最小值；

在任何情况下，ready 消息数总是被计算在内的；unacknowledged 消息不会被计算到限制值内；`rabbitmqctl list_queues` 输出项中的 messages_ready 和 message_bytes_ready ，以及相应的 management API 可以展示的限制值；

当限制值被达到时，消息会被从 queue 的前端被丢弃或者 dead-lettered，以便为新消息留出空间；

## Configuration using arguments

最大消息数目（非负整数）可以通过  `x-max-length` 参数在 queue 声明中进行设置；

最大消息字节数目（非负整数）可以通过  `x-max-length-bytes` 参数在 queue 声明中进行设置；

如果两种参数都设置了，则同时生效；先达到限定值的限制先起作用；

下面示例中的 Java 代码声明了最多允许保存 10 条消息的 queue ：

```java
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-max-length", 10);
channel.queueDeclare("myqueue", false, false, false, args);
```

## Configuration using policy

若想通过 policy 指定最大长度，可以将 key `max-length` 和/或 `max-length-bytes` 添加到 policy 定义中，例如：

| | |
-------- | ---
| rabbitmqctl | rabbitmqctl set_policy Ten "^one-meg$" '{"max-length-bytes":1000000}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Ten "^one-meg$" "{""max-length-bytes"":1000000}" --apply-to queues |

上述命令确保了名为 one-meg 的 queue 能够保存不超过 1MB 的消息 body 大小；

Policies 的设置还可以使用 management plugin 进行，详见 [policy documentation](http://www.rabbitmq.com/parameters.html#policies) ；

