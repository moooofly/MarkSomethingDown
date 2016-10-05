



# [Queue Length Limit](http://www.rabbitmq.com/maxlength.html)

The maximum length of a queue can be limited to a set number of messages, or a set number of bytes (the total of all message body lengths, ignoring message properties and any overheads), or both.

For any given queue, the maximum length (of either type) can be defined by clients using the queue's arguments, or in the server using [policies](http://www.rabbitmq.com/parameters.html#policies). In the case where both policy and arguments specify a maximum length, the minimum of the two values is applied.

In all cases the number of ready messages is used; unacknowledged messages do not count towards the limit. The fields messages_ready and message_bytes_ready from rabbitmqctl list_queues and the management API show the values that would be limited.

Messages will be dropped or dead-lettered from the front of the queue to make room for new messages once the limit is reached.

## Configuration using arguments

Maximum number of messages can be set by supplying the x-max-length queue declaration argument with a non-negative integer value.

Maximum length in bytes can be set by supplying the x-max-length-bytes queue declaration argument with a non-negative integer value.

If both arguments are set then both will apply; whichever limit is hit first will be enforced.

This example in Java declares a queue with a maximum length of 10 messages:

```java
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-max-length", 10);
channel.queueDeclare("myqueue", false, false, false, args);
```

## Configuration using policy

To specify a maximum length using policy, add the key max-length and / or max-length-bytes to a policy definition. For example:

| | |
-------- | ---
| rabbitmqctl | rabbitmqctl set_policy Ten "^one-meg$" '{"max-length-bytes":1000000}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Ten "^one-meg$" "{""max-length-bytes"":1000000}" --apply-to queues |

	

This ensures the queue called one-meg can contain no more than 1MB of message bodies.

Policies can also be defined using the management plugin, see the [policy documentation](http://www.rabbitmq.com/parameters.html#policies) for more details.

