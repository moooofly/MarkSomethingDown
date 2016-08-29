


# Lazy Queues

Since RabbitMQ 3.6.0 the broker has the concept of Lazy Queues: these are queues that try to keep as many messages as possible on disk, and only load them in RAM when requested by consumers, therefore the lazy denomination.

One of the main goals of lazy queues is to be able to support very long queues (many millions of messages). These queues can arise when consumers are unable to fetch messages from queues for long periods of times. This can happen for various reasons and use cases: because consumers are offline; because they have crashed, or they have been taken down for maintenance; and so on.

By default, queues keep an in-memory cache of messages that's filled up as messages are published into RabbitMQ. The idea of this cache is to be able to deliver messages to consumers as fast as possible (note that persistent messages are written to disk as they enter the broker and kept in this cache at the same time). Whenever the broker considers it needs to free up memory, messages from this cache will be paged out to disk. Paging messages to disk takes time and block the queue process, making it unable to receive new messages while it's paging. Even if on recent RabbitMQ versions we have improved the paging algorithm, the situation is still not ideal for use cases where you have many millions on messages in the queue that might need to be paged out.

Lazy queues help here by eliminating this cache and only loading messages in memory when requested by consumers. Lazy queues will send every message that arrives to the queue right away to the file system, completely eliminating the in-memory cache mentioned before. This has the consequence of heavily reducing the amount of RAM consumed by a queue and also eliminates the need for paging. While this will increase I/O usage, it is the same behaviour as when publishing persistent messages.


# Using Lazy Queues

Queues can be made to work in default mode or lazy mode either by specifying the mode via queue.declare arguments, or by applying a policy in the server. In the case where both policy and queue arguments specify a queue mode, then the queue argument has priority over the policy value. This means that if a queue mode is set via a declare argument, it can only be changed by deleting the queue, and re-declaring it later with a different argument.

## Configuration using arguments

The queue mode can be set by supplying the x-queue-mode queue declaration argument with a string specifying the desired mode. Valid modes are "default" and "lazy". If no mode is specified during declare, then "default" is assumed. The default mode is the behaviour already present in pre 3.6.0 versions of the broker, so there are no breaking changes in this regard.

This example in Java declares a queue with the queue mode set to "lazy":

```
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-queue-mode", "lazy");
channel.queueDeclare("myqueue", false, false, false, args);
```

## Configuration using policy

To specify a queue mode using a policy, add the key queue-length to a policy definition. For example:

|                  |                              |                  |
 ----------------- | ---------------------------- | ------------------
| rabbitmqctl | rabbitmqctl set_policy Lazy "^lazy-queue$" '{"queue-mode":"lazy"}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Lazy "^lazy-queue$" "{""queue-mode"":""lazy""}" --apply-to queues |

This ensures the queue called lazy-queue will work in the lazy mode.

Policies can also be defined using the management plugin, see the policy documentation for more details.

## Changing queue modes

If you specified the queue mode via a policy, then you can change it at run time without the need of deleting the queue and re-declaring it with a different mode. If you want the previous lazy-queue to start working like a default queue, then you can do so by issuing the following command:

|                  |                              |                  |
 ----------------- | ---------------------------- | ------------------
| rabbitmqctl | rabbitmqctl set_policy Lazy "^lazy-queue$" '{"queue-mode":"default"}' --apply-to queues|
| rabbitmqctl (Windows) | rabbitmqctl set_policy Lazy "^lazy-queue$" "{""queue-mode"":""default""}" --apply-to queues |


# Performance Considerations for Lazy Queues

## Disk Utilization
As stated above, lazy queues will send every message to disk right as they enter the queue. This will increase I/O opts, but keep in mind that this is the same behaviour as when persistent messages are delivered to queues. Note that even if you publish transient messages, they will still be sent to disk when using lazy queues. With default queues transient messages are only sent to disk if paging requires it.

## RAM Utilization
lazy queues use much less memory than default queues. While it's hard to give numbers that make sense for every use case, here's what we found: we tried publishing 10 million messages into a queue, with no consumers online. The message body size was 1000 bytes. default queues required 1.2GB of RAM, while lazy queues only used 1.5MB of RAM.

For a default queue, it took 801 seconds to send 10MM messages, with an average sending rate of 12469 msg/s. To publish the same amount of messages into a lazy queue, the time required was 421 seconds, with an average sending rate of 23653 msg/s. The difference can be explained by the fact that from time to time, the default queue had to page messages to disk. Once we activated a consumer, the lazy queue had a RAM consumption of approximately 40MB while it was delivering messages. The message receiving rate average was 13938 msg/s for one active consumer.

You can reproduce the test with our Java library by running:

```shell
./runjava.sh com.rabbitmq.examples.PerfTest -e test -u test_queue \
-f persistent -s 1000 -x1 -y0 -C10000000
```
Note that this was a very simplistic test. Please make sure to run your own benchmarks.

Don't forget to change the queue mode between benchmarks runs.


## Converting between queue modes

If we need to convert a default queue into a lazy one, then we will suffer the same performance impact as when a queue needs to page messages to disk. When we convert a queue into a lazy one, first it will page messages to disk and then it will start accepting publishes, acks, and other commands.

When a queue goes from the lazy mode to the default one, it will perform the same process as when a queue is recovered after a server restart. A batch of 16384 messages will be loaded in the cache mentioned above.