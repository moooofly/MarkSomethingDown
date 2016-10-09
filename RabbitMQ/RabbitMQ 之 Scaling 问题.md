
# [Scaling RabbitMQ](http://jaredrobinson.com/Scaling_RabbitMQ/)

## What's Covered

- Basics
- Work queues, then and now
- Latency, then and now
- Speed $ Throughput
- Fail-over
- limits.conf
- Transient vs durable trade-offs
- Backpressure and RabbitMQ death
- Grouping messages
- Resources


## What's Not Covered

- Publish/Subscribe, Topics, RPC
- Alternate technologies
   - Apache Kafka
   - ZeroMQ


## Did You Know?

- Good book: **RabbitMQ in Action** 
- Practical advice, born from experience
- Discusses High-Availability Pitfalls & Solutions
    - especially with regard to durable queues

## Latency

- Delivery by carrier pigeon isn't good enough. Nor do we want delivery by cargo-ship 
- "**Latency is the root of all that is evil on the Internet**... or so the saying goes." -- Theo Schlossnagle
- In the beginning: Round-robin distribution & slow workers
- Database: Contention, Slow queries, Expensive queries, Non-indexed queries
- CPU contention and VMs on the same hypervisor
- Slow third-party services


## Latency & Prefetch

> Prefetch protects us against slow consumers


## Did You Know?

When publishing to an exchange that isn't bound to a queue, there are no errors from the client perspective.


## Throughput & Speed

- "Fast hands make light work"? 
- **Our max speed is ~10K messages/sec (w/o HiPE)**
- **HiPE**
    - 20-50% better performance
    - not available on Windows
    - Experimental. Disable if it segfaults.
    - not available with RHEL/CentOS erlang
- Speed of hardware
- Is it running in VM? CPU/RAM contention.
- Do too many consumers slow it down?


## Throughput Needs

- As of March, 2015:
    - 50 user-initiated commands/second
    - 1,500 inbound signals/second
    - 7,000 messages/second on our busiest queue
- Doubling this year
- Need headroom for spikes and a buffer for growth beyond this year
- For coming year:
    - 21,000 messages/second on our busiest queue
    - Even with HiPE, a single work queue isn't sufficient


## Throughput Via Sharding

"Many hands make light work"

![](http://jaredrobinson.com/Scaling_RabbitMQ/Hands.jpg)


## Throughput: Are You My Solution?

- Are you my solution? 
    - rabbitmq-consistent-hash-exchange
    - rabbitmq-sharding
    - client-side, shared nothing
    - Federated Queues
- Try it, test it, benchmark it. Fail fast & move on.


## Consistent Hash Exchange

- Plugin ships with Rabbit!
- Gracefully handles the disappearance of a queue on a crashed rabbit node
- We determine link between the exchange and downstream queues (where queues live, how many there are, etc.)
- A single RabbitMQ management console shows the entire cluster!

## Consistent Hash Exchange Configuration

```shell
$ umask 0022    -- important if you're doing this as root
$ rabbitmq-plugins enable rabbitmq_consistent_hash_exchange
--
The following plugins have been enabled:
  rabbitmq_consistent_hash_exchange
  Plugin configuration has changed. Restart RabbitMQ for changes to take effect.
--
$ service rabbitmq-server restart
$ /usr/sbin/rabbitmq-plugins list -E
```

## Consistent Hash Exchange Configuration

- Declare the exchange as type "**x-consistent-hash**"
- Bind the exchange to downstream queues (you decide where queues live, how many there are, etc.)
- When binding the queue to the exchange, the routing key must be a number-as-a-string
- Insufficient: "2", Better: "101"
- Producer must vary the per-message routing key -- random, or an id
- Not round-robin -- not an even distribution
- Don't forget cluster port and firewall configuration


## Consistent Hash Exchange Caveats

- Unlike a simple work-queue, **requires differing routing keys** -- random, or use an id
- Definitely consistent, but **not an even distribution**
- **Has a chance of losing messages** when unbinding or deleting a queue
    - There's a moment of time between the determination of which queue a message should be sent to, and the time it's published
- **Clients must be configured to evenly connect to the queue shards**
- Try it, test it, benchmark it. Fail fast & move on.


## Rabbitmq-Sharding

- Built by the **Pivotal/RabbitMQ** folks
- Based on the consistent hash exchange
- **Auto-creates and binds** the queues to the exchange!
- **Transparently auto-connects** client to the shard with the least number of consumers!
- Doesn't ship with Rabbit
- https://github.com/rabbitmq/rabbitmq-sharding

> 从 3.6.0 开始，已经随 RabbitMQ 直接发布了；

## Rabbitmq-Sharding Configuration

Download from https://www.rabbitmq.com/community-plugins/v3.3.x/

```shell
cp rabbitmq_sharding-3.3.x.ez  \
  /usr/lib/rabbitmq/lib/rabbitmq_server-3.3.5/plugins/ 

rabbitmq-plugins enable rabbitmq_sharding
rabbitmq-plugins enable rabbitmq_consistent_hash_exchange

service rabbitmq-server restart

rabbitmqctl set_policy history-shard "^history" \
  '{"shards-per-node": 2, "routing-key": "1234"}' \
  --apply-to exchanges
```

## Rabbitmq-Sharding Client Usage

- Clients declare the exchange type as "**x-modulus-hash**"
- **Equal message distrubituion** using uniform random routing key
    - Use same routing key, and all messages will go to a single queue shard
- When nodes appear and disappear, the exchange automatically creates queue shards on the new node.
    - It re-creates missing shards from nodes where I deleted them.


## Rabbitmq-Sharding Caveats

- Auto-creates and binds the queues to the exchange!
    - **Must decide how many queue shards per node**. Want 3 queue shards on a four-node Rabbit cluster? Sorry.
- Transparently auto-connects client to the shard with the least number of consumers!
    - **Only on a per node basis, not per cluster**
- When one node fills up, producers block, even though another node may be ready to accept messages
- **Has a chance of losing messages** when unbinding or deleting a queue
- Try it, test it, benchmark it. Fail fast & move on.

> 关于 producer 被 block 的问题，应该是 cluster 本身特定导致；

## Did You Know?

- "Each connection uses a file descriptor on the server. Channels don't"
- **Publishing a large message on one channel will block a connection while it goes out**.
- "it's a good idea to separate publishing and consuming connections"

> ⚠️ 为什么在一个 channel 上发送一个大消息后，会阻塞整个 connection ？不明白！

## Client-Side Sharding

- Doesn't ship with Rabbit -- more work to implement
- Shared-nothing RabbitMQ nodes
    - Homegrown monitoring to combine metrics from multiple servers
- Implement once for each language we use: Python, Java
- Gives us fail-over!
- Must balance consumers among the shards

## Client-Side Sharding

Producers, Consumers, YAML
- Producers: distribute messages round-robin
- Producers: auto-detect node failure and reconnect after a period of time
- Consumers: Choosing a shard to read from didn't give us an even distribution
- Consumers: "sticky" -- no fail-over to a separate shard
    - need extra consumers to handle the failure of a single node
- Configuration is done via YAML
    - Every producer knows about every shard
    - Update and deploy using salt-cp

## Client-Side Sharding

Benefits
- If a node dies, we don't have to remove it from the cluster
- Durable queues work well in the face of node failure (shared nothing)
- Mix-and-match rabbit versions (atypical)


## Client-Side Sharding Caveats

- More work to implement and monitor
- **Producers still block when one of the nodes exerts backpressure (flow control)**

## Federated Queues

- Distribute "the same 'logical' queue... over many brokers."
- Configured via policy
- Mix-and-match RabbitMQ versions
- Performs "best when there is some degree of locality...."

## Federated Queues

Better together?
- Combine with rabbit sharding plugins or with client-side sharding
- Keep queues drained

## Did You Know?

Client can publish to an exchange that doesn't exist, without an error.

## Default Limits.Conf

Default (insufficient) limits.conf for RHEL/CentOS/Fedora

```shell
*  soft  nproc   1024 # threads/processes
*  soft  nofile  1024 # Number of open files
*  hard  nofile  4096 # Number of open files
```

## Limits.Conf For RabbitMQ 

Here's what I'm using for /etc/security/limits.conf

```shell
rabbitmq  soft  nproc   16384
rabbitmq  hard  nofile  16000
rabbitmq  soft  nofile  16000
```

## Did You Know?

- A client can't read from a queue that doesn't exist. It causes an error.
- Queues balloon in size when being fed, and no consumers are listening
    - **Use message TTL**
    - **Consider combining TTL with a dead letter exchange**
    - Dead-letter exchange can be configured by clients, or via RabbitMQ policy
    - `rabbitmqctl set_policy DLX ".*" '{"dead-letter-exchange":"my-dlx"}' --apply-to queues`



## Fail-Over

- What happens to your system when a rabbit node dies?
- Load balancer
- Client-side fail-over
- Consumer capacity
- Note: Durable queues can't be re-created on a separate node in the cluster


## Transient Vs Durable

- Opinion: **The best place to persist messages is at the origin**
- Origin can resend when not acknowledged (publisher confirms)
- Transient is faster than durable, although durable survives RabbitMQ restarts
- RabbitMQ is fastest when few or no messages are in the queue
- RabbitMQ stores transient messages to disk when the queue balloons in size


## Did You Know?

- A queue can be bound to multiple exchanges, or the same exchange, multiple times, with different binding keys.


## Backpressure And RabbitMQ Death

- When RabbitMQ gets busy (persisting messages to disk), it exerts backpressure on producers. Producers block on publish.
- What if we don't want a producer to block?
- RabbitMQ monitors disk space for persistent queues.
- Bug: Out-of-disk space crashes RabbitMQ when transient messages are being paged to disk
    - **Restarting doesn't clear/fix the disk space consumed**
    - Remove `/var/rabbitmq/mnesia/rabbit@YourHost/msg_store_transient/`
    - Consider using configuration setting "`{vm_memory_high_watermark_paging_ratio, 1.1}`"


## Grouping Messages

- It may be possible to achieve higher throughput by combining your messages into a single Rabbit message
- Increase in latency

> 这里的说法更合理：将消息分组后发送能够增加吞吐量，但也会增大延迟；

## Summary 1 Of 2

- Configure **prefetch**
- Try **HiPE**
- When **scaling horizontally**, there's more than one way to do it
- Plan for fail-over
- Configure limits.conf


## Summary 2 Of 2

- Use good tools -- RabbitMQ is a good tool
- It's available, cost-effective, practical
- Try solutions, test solutions. If one idea doesn't work, get feedback, and try something else. Fail fast.


## References
- [RabbitMQ in Action](https://www.manning.com/books/rabbitmq-in-action)
- [Mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users)
- stackoverflow on [Maximizing throughput with RabbitMQ](http://stackoverflow.com/questions/10030227/maximize-throughput-with-rabbitmq)
- RabbitMQ blog on [queuing theory: throughput, latency and bandwidth](https://www.rabbitmq.com/blog/2012/05/11/some-queuing-theory-throughput-latency-and-bandwidth/)
- [Building a Distributed Data Ingestion System with RabbitMQ](http://www.erlang-factory.com/euc2014/alvaro-videla) by Alvaro Videla, co-author of RabbitMQ in Action



