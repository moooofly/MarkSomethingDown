

----------

1. [Gossip protocol in wikipedia](https://en.wikipedia.org/wiki/Gossip_protocol)
2. [Gossip Protocol in Serf](https://www.serf.io/docs/internals/gossip.html)
3. [Using Gossip Protocols For Failure Detection, Monitoring, Messaging And Other Good Things](http://highscalability.com/blog/2011/11/14/using-gossip-protocols-for-failure-detection-monitoring-mess.html)
4. [Gossip算法学习](http://blog.csdn.net/yfkiss/article/details/6943682#comments)
5. [Gossip算法](http://tianya23.blog.51cto.com/1081650/530743)


----------

# Gossip protocol

A gossip protocol [1] is a style of computer-to-computer [communication protocol](https://en.wikipedia.org/wiki/Communications_protocol) inspired by the form of gossip seen in social networks. Modern [distributed systems](https://en.wikipedia.org/wiki/Distributed_computing) often use gossip protocols to solve problems that might be difficult to solve in other ways, either because the underlying network has an inconvenient structure, is extremely large, or because gossip solutions are the most efficient ones available.

The term epidemic protocol is sometimes used as a synonym for a gossip protocol, because gossip spreads information in a manner similar to the spread of a virus in a biological community.

## Gossip communication

The concept of gossip communication can be illustrated by the analogy of office workers spreading rumors. Let's say each hour the office workers congregate around the water cooler. Each employee pairs off with another, chosen at random, and shares the latest gossip. At the start of the day, Alice starts a new rumor: she comments to Bob that she believes that Charlie dyes his mustache. At the next meeting, Bob tells Dave, while Alice repeats the idea to Eve. After each water cooler rendezvous, the number of individuals who have heard the rumor roughly doubles (though this doesn't account for gossiping twice to the same person; perhaps Alice tries to tell the story to Frank, only to find that Frank already heard it from Dave). Computer systems typically implement this type of protocol with a form of random "peer selection": with a given frequency, each machine picks another machine at random and shares any hot rumors.

The power of gossip lies in the robust spread of information. Even if Dave had trouble understanding Bob, he will probably run into someone else soon and can learn the news that way.

Expressing these ideas in more technical terms, a gossip protocol is one that satisfies the following conditions:
- The core of the protocol involves periodic, pairwise, inter-process interactions.
- The information exchanged during these interactions is of bounded size.
- When agents interact, the state of at least one agent changes to reflect the state of the other.
- Reliable communication is not assumed.
- The frequency of the interactions is low compared to typical message latencies so that the protocol costs are negligible.
- There is some form of randomness in the peer selection. Peers might be selected from the full set of nodes or from a smaller set of neighbors.
- Due to the replication there is an implicit redundancy of the delivered information.

## Gossip protocol types

It is useful to distinguish three prevailing styles of gossip protocol:[2]

- Dissemination protocols (or rumor-mongering protocols). These use gossip to spread information; they basically work by flooding agents in the network, but in a manner that produces bounded worst-case loads:
      - Event dissemination protocols use gossip to carry out multicasts. They report events, but the gossip occurs periodically and events don’t actually trigger the gossip. One concern here is the potentially high latency from when the event occurs until it is delivered.
      - Background data dissemination protocols continuously gossip about information associated with the participating nodes. Typically, propagation latency isn’t a concern, perhaps because the information in question changes slowly or there is no significant penalty for acting upon slightly stale data.
- [Anti-entropy protocols](https://en.wikipedia.org/wiki/Error_detection_and_correction) for repairing replicated data, which operate by comparing replicas and reconciling differences.
- Protocols that compute aggregates. These compute a network-wide aggregate by sampling information at the nodes in the network and combining the values to arrive at a system-wide value – the largest value for some measurement nodes are making, smallest, etc. The key requirement is that the aggregate must be computable by fixed-size pairwise information exchanges; these typically terminate after a number of rounds of information exchange logarithmic in the system size, by which time an all-to-all information flow pattern will have been established. As a side effect of aggregation, it is possible to solve other kinds of problems using gossip; for example, there are gossip protocols that can arrange the nodes in a gossip overlay into a list sorted by node-id (or some other attribute) in logarithmic time using aggregation-style exchanges of information. Similarly, there are gossip algorithms that arrange nodes into a tree and compute aggregates such as "sum" or "count" by gossiping in a pattern biased to match the tree structure.
Many protocols that predate the earliest use of the term "gossip" fall within this rather inclusive definition. For example, Internet [routing protocols](https://en.wikipedia.org/wiki/Routing_protocol) often use gossip-like information exchanges. A gossip substrate can be used to implement a standard routed network: nodes "gossip" about traditional point-to-point messages, effectively pushing traffic through the gossip layer. Bandwidth permitting, this implies that a gossip system can potentially support any classic protocol or implement any classical distributed service. However, such a broadly inclusive interpretation is rarely intended. More typically gossip protocols are those that specifically run in a regular, periodic, relatively lazy, symmetric and decentralized manner; the high degree of symmetry among nodes is particularly characteristic. Thus, while one could run a [2-phase commit protocol](https://en.wikipedia.org/wiki/Two-phase_commit_protocol) over a gossip substrate, doing so would be at odds with the spirit, if not the wording, of the definition.

Frequently, the most useful gossip protocols turn out to be those with exponentially rapid convergence towards a state that "emerges" with probability 1.0. A classic distributed computing problem, for example, involves building a tree whose inner nodes are the nodes in a network and whose edges represent links between computers (for routing, as a dissemination overlay, etc.). Not all tree-building protocols are gossip protocols (for example, spanning tree constructions in which a leader initiates a flood), but gossip offers a decentralized solution that is useful in many situations.

The term convergently consistent is sometimes used to describe protocols that achieve exponentially rapid spread of information. For this purpose, a protocol must propagate any new information to all nodes that will be affected by the information within time logarithmic in the size of the system (the "mixing time" must be logarithmic in system size).

## Examples

Suppose that we want to find the object that most closely matches some search pattern, within a network of unknown size, but where the computers are linked to one another and where each machine is running a small agent program that implements a gossip protocol.

To start the search, a user would ask the local agent to begin to gossip about the search string. (We're assuming that agents either start with a known list of peers, or retrieve this information from some kind of a shared store.)
Periodically, at some rate (let's say ten times per second, for simplicity), each agent picks some other agent at random, and gossips with it. Search strings known to A will now also be known to B, and vice versa. In the next "round" of gossip A and B will pick additional random peers, maybe C and D. This round-by-round doubling phenomenon makes the protocol very robust, even if some messages get lost, or some of the selected peers are the same or already know about the search string.
On receipt of a search string for the first time, each agent checks its local machine for matching documents.
Agents also gossip about the best match, to date. Thus, if A gossips with B, after the interaction, A will know of the best matches known to B, and vice versa. Best matches will "spread" through the network.
If the messages might get large (for example, if many searches are active all at the same time), a size limit should be introduced. Also, searches should "age out" of the network.

It follows that within logarithmic time in the size of the network (the number of agents), any new search string will have reached all agents. Within an additional delay of the same approximate length, every agent will learn where the best match can be found. In particular, the agent that started the search will have found the best match.

For example, in a network with 25,000 machines, we can find the best match after about 30 rounds of gossip: 15 to spread the search string and 15 more to discover the best match. A gossip exchange could occur as often as once every tenth of a second without imposing undue load, hence this form of network search could search a big data center in about 3 seconds.

In this scenario, searches might automatically age out of the network after, say, 10 seconds. By then, the initiator knows the answer and there is no point in further gossip about that search.

Gossip protocols have also been used for achieving and maintaining [distributed database](https://en.wikipedia.org/wiki/Distributed_database) [consistency](https://en.wikipedia.org/wiki/Consistency) or with other types of data in consistent states, counting the number of nodes in a network of unknown size, spreading news robustly, organizing nodes according to some structuring policy, building so-called [overlay networks](https://en.wikipedia.org/wiki/Overlay_network), computing aggregates, sorting the nodes in a network, electing leaders, etc.

## Epidemic algorithms

Gossip protocols can be used to propagate information in a manner rather similar to the way that a viral infection spreads in a biological population. Indeed, the mathematics of epidemics are often used to model the mathematics of gossip communication. The term epidemic algorithm is sometimes employed when describing a software system in which this kind of gossip-based information propagation is employed.


## Biased gossip

Above, a purely random peer-selection scheme for gossip was described: when agent A decides to run a gossip round, it picks some peer B uniformly and at random within the network as a whole (or launches a message on a random walk that will terminate at a random agent). More commonly, gossip algorithms are designed so that agents interact mostly with nearby agents, and only sometimes with agents that are far away (in terms of network delay). These biased gossip protocols need to ensure a sufficient degree of connectivity to avoid the risk of complete disconnection of one side of a network from the other, but if care is taken, can be faster and more efficient than protocols that are purely random. Moreover, as a purely practical question, it is much easier to maintain lists of peers in ways that might be somewhat biased.

## Code examples

There are two known libraries which implement a Gossip algorithm to discover nodes in a peer-to-peer network: [teknek-gossip](https://github.com/edwardcapriolo/gossip) works with UDP and is written in Java. [gossip-python](https://github.com/thomai/gossip-python) utilizes the TCP stack and it is possible to share data via the constructed network as well.

## See also

- Gossip protocols are just one class among many classes of networking protocols. See also [virtual synchrony](https://en.wikipedia.org/wiki/Virtual_synchrony), distributed [state machines](https://en.wikipedia.org/wiki/Finite-state_machine), [Paxos algorithm](https://en.wikipedia.org/wiki/Paxos_(computer_science)), database [transactions](https://en.wikipedia.org/wiki/Database_transaction). Each class contains tens or even hundreds of protocols, differing in their details and performance properties but similar at the level of the guarantees offered to users.
- Some Gossip protocols replace the random peer selection mechanism with a more deterministic scheme. For example, in the [NeighbourCast](http://www.actapress.com/PaperInfo.aspx?PaperID=31994&reason=500) Algorithm, instead of talking to random nodes, information is spread by talking only to neighbouring nodes. There are a number of algorithms that use similar ideas. A key requirement when designing such protocols is that the neighbor set trace out an [expander graph](https://en.wikipedia.org/wiki/Expander_graph).
- [Routing](https://en.wikipedia.org/wiki/Routing)
- [Tribler](https://en.wikipedia.org/wiki/Tribler), BitTorrent peer to peer client using gossip protocol.


## References

略


原文：[这里](https://en.wikipedia.org/wiki/Gossip_protocol)


----------


# Gossip Protocol

Serf uses a [gossip protocol](https://en.wikipedia.org/wiki/Gossip_protocol) to broadcast messages to the cluster. This page documents the details of this internal protocol. The gossip protocol is based on ["SWIM: Scalable Weakly-consistent Infection-style Process Group Membership Protocol"](https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf), with a few minor adaptations, mostly to increase propagation speed and convergence rate.


## SWIM Protocol Overview

Serf begins by joining an existing cluster or starting a new cluster. If starting a new cluster, additional nodes are expected to join it. New nodes in an existing cluster must be given the address of at least one existing member in order to join the cluster. The new member does a full state sync with the existing member over TCP and begins gossiping its existence to the cluster.

Gossip is done over UDP with a configurable but fixed fanout and interval. This ensures that network usage is constant with regards to number of nodes. Complete state exchanges with a random node are done periodically over TCP, but much less often than gossip messages. This increases the likelihood that the membership list converges properly since the full state is exchanged and merged. The interval between full state exchanges is configurable or can be disabled entirely.

Failure detection is done by periodic random probing using a configurable interval. If the node fails to ack within a reasonable time (typically some multiple of RTT), then an indirect probe is attempted. An indirect probe asks a configurable number of random nodes to probe the same node, in case there are network issues causing our own node to fail the probe. If both our probe and the indirect probes fail within a reasonable time, then the node is marked "suspicious" and this knowledge is gossiped to the cluster. A suspicious node is still considered a member of cluster. If the suspect member of the cluster does not dispute the suspicion within a configurable period of time, the node is finally considered dead, and this state is then gossiped to the cluster.

This is a brief and incomplete description of the protocol. For a better idea, please read the SWIM paper in its entirety, along with the Serf source code.

## SWIM Modifications

As mentioned earlier, the gossip protocol is based on SWIM but includes minor changes, mostly to increase propagation speed and convergence rates.

The changes from SWIM are noted here:

- Serf does a full state sync over TCP periodically. SWIM only propagates changes over gossip. While both are eventually consistent, Serf is able to more quickly reach convergence, as well as gracefully recover from network partitions.

- Serf has a dedicated gossip layer separate from the failure detection protocol. SWIM only piggybacks gossip messages on top of probe/ack messages. Serf uses piggybacking along with dedicated gossip messages. This feature lets you have a higher gossip rate (for example once per 200ms) and a slower failure detection rate (such as once per second), resulting in overall faster convergence rates and data propagation speeds.

- Serf keeps the state of dead nodes around for a set amount of time, so that when full syncs are requested, the requester also receives information about dead nodes. Because SWIM doesn't do full syncs, SWIM deletes dead node state immediately upon learning that the node is dead. This change again helps the cluster converge more quickly.

## Serf-Specific Messages

On top of the SWIM-based gossip layer, Serf sends some custom message types.

Serf makes heavy use of [Lamport clocks](https://en.wikipedia.org/wiki/Lamport_timestamps) to maintain some notion of message ordering despite being eventually consistent. Every message sent by Serf contains a Lamport clock time.

When a node gracefully leaves the cluster, Serf sends a leave intent through the gossip layer. Because the underlying gossip layer makes no differentiation between a node leaving the cluster and a node being detected as failed, this allows the higher level Serf layer to detect a failure versus a graceful leave.

When a node joins the cluster, Serf sends a join intent. The purpose of this intent is solely to attach a Lamport clock time to a join so that it can be ordered properly in case a leave comes out of order.

For custom events and queries, Serf sends either a user event, or user query message. This message contains a Lamport time, event name, and event payload. Because user events are sent along the gossip layer, which uses UDP, the payload and entire message framing must fit within a single UDP packet.


原文：[这里](https://www.serf.io/docs/internals/gossip.html)


----------


# Using Gossip Protocols For Failure Detection, Monitoring, Messaging And Other Good Things

When building a system on top of a set of wildly uncooperative and unruly computers you have knowledge problems: knowing when other nodes are dead; knowing when nodes become alive; getting information about other nodes so you can make local decisions, like knowing which node should handle a request based on a scheme for assigning nodes to a certain range of users; learning about new configuration data; agreeing on data values; and so on.

How do you solve these problems? 

A common centralized approach is to use a database and all nodes query it for information. Obvious availability and performance issues for large distributed clusters. Another approach is to use [Paxos](https://en.wikipedia.org/wiki/Paxos_(computer_science)), a protocol for solving consensus in a network to maintain strict consistency requirements for small groups of unreliable processes. Not practical when larger number of nodes are involved.

So what's the super cool decentralized way to bring order to large clusters?

[Gossip protocols](https://en.wikipedia.org/wiki/Gossip_protocol), which maintain relaxed consistency requirements amongst a very large group of nodes. A gossip protocol is simple in concept. Each nodes sends out some data to a set of other nodes. Data propagates through the system node by node like a virus. Eventually data propagates to every node in the system. It's a way for nodes to build a global map from limited local interactions.

As you might imagine there are all sorts of subtleties involved, but at its core it's a simple and robust system. A node only has to send to a subset of other nodes. That's it.

Cassandra, for example, uses what's called an [anti-entropy version](http://wiki.apache.org/cassandra/ArchitectureAntiEntropy) of the gossip protocol for repairing unread data using [Merkle Trees](https://en.wikipedia.org/wiki/Hash_tree). Riak uses a gossip protocol to share and communicate ring state and bucket properties around the cluster. 

For a detailed look at using gossip protocols take a look at GEMS: Gossip-Enabled Monitoring Service for Scalable Heterogeneous Distributed Systems by Rajagopal Subramaniyan, Pirabhu Raman, Alan George, and Matthew Radlinski. I really like this paper because of how marvelously well written and clear it is on how to use gossip protocols to detect node failures and load balance based on data sampled from other other nodes. Details are explained clearly and it dares to cover a variety of possibly useful topics.

From the abstract:

Gossip protocols have proven to be effective means by which failures can be detected in large, distributed systems in an asynchronous manner without the limitations associated with reliable multicasting for group communications. In this paper, we discuss the development and features of a Gossip-Enabled Monitoring Service (GEMS), a highly responsive and scalable resource monitoring service, to monitor health and performance information in heterogeneous distributed systems. GEMS has many novel and essential features such as detection of network partitions and dynamic insertion of new nodes into the service. Easily extensible, GEMS also incorporates facilities for distributing arbitrary system and application-specific data. We present experiments and analytical projections demonstrating scalability, fast response times and low resource utilization requirements, making GEMS a potent solution for resource monitoring in distributed computing.

Failure Detection

The failure detection part of the paper is good and makes sense. By combining the reachability data from a lot of different nodes you can quickly determine when a node is down. When a node is down, for example, there's no need to attempt to write to that node, saving queue space, CPU, and bandwidth.

In a distributed system you need at least two independent sources of information to mark a node down. It's not enough to simply say because your node can't contact another node that the other node is down. It's quite possible that your node is broken and the other node is fine. But if other nodes in the system also see that other node is dead then you can with some confidence conclude that that node is dead. Many complex hard to debug bugs are hidden here. How do you know what other nodes are seeing? Via a gossip protocol exchanging this kind of reachability data.

In embedded systems the backplane often has traces between nodes so a local system can get an independent source of confirmation that a given node is dead, or alive, or transitioning between the two states. If the datacenter is really the computer, it would be nice to see datacenters step up and implement higher level services like node liveness and time syncing so every application doesn't have to worry about these issues, again.

The paper covers the obvious issue of scaling as the number of nodes increases by dividing nodes into groups and introducing a hierarchy of layers at which node information is aggregated. They found running the gossip protocol used less than 60 Kbps of bandwidth and less than 2% of CPU for a system of 128 nodes.

One thing I would add is the communication subsystem can also contribute what it learns about reachability, we don't just have to rely on a gossip heartbeat. If the communication layer can't reach a node that fact can be noted in a reachability table. This keeps data as up to date as possible.

Using Gossip As A Form Of Messaging 

In addition to failure detection, the paper shows how to transmit node and subsystem properties between nodes. This is a great extension and is a far more robust mechanism than individual modules using TCP connections to exchange data and command and control. We want to abstract communication out of application level code and this type of approach accomplishes that.

It seems somewhat obvious that you would transmit node properties to other nodes. Stats like load average, free memory, etc. would allow a local node to decide where to send work, for example. If a node is idle send it work (as long as everyone doesn't send it work at the same time). This local decision making angle is the key to scale. There's no centralized controller. Local nodes make local decisions based on local data. This can scale as far as the gossip protocol can scale.

What goes to another level is that they use an architecture I've used on several products, sending subsystem information so software modules on a node can send information to other modules on other nodes. For example, queue depth for a module could be sent out so other modules could gauge the work load. Alarm information could be sent out so other entities know the status of modules they are dependent on. Key information like configuration changes can be passed on. Even requests and response can be sent through this mechanism. At an architecture level this allows the aggregation of updates (from all sources on a node) so they can be sent in big blocks through the system instead of small messages, which is a big win.

This approach can be combined with a publish/subscribe topic registration system to reduce useless communication between nodes.

Another advantage of this approach is data could flow directly into your monitoring system rather than having a completely separate monitoring subsystem bolted on.

In the meat world we are warned against gossiping, it's a sin, it can ruin lives, it can ruin your reputation, etc., but in software, gossiping is a powerful tool in your distributed toolbox. So go forth and gossip.


原文：[这里](http://highscalability.com/blog/2011/11/14/using-gossip-protocols-for-failure-detection-monitoring-mess.html)


----------


# Gossip算法学习

## 概述

gossip，顾名思义，类似于流言传播的概念，是一种可以按照自己的期望，自行选择与之交换信息的节点的通信方式；

> **gossip**, or **anti-entropy**,  is an attractive way of replicating state that does not have strong consistency requirements.


## 算法描述

假设有 `{p, q, ...}` 为协议参与者。 每个参与者都有一个关于自己信息的表；

用编程语言可以描述为： 
- 每一个参与者都有一个关于自己信息的表，即要维护一个 `InfoMap` 类型的 `localInfo` ，记 `InfoMap = Map<Key, (Value, Version)>` ；    
- 每一个参与者还要知道所有其他参与者的信息，即要维护一个 `globalMap` 类型的全局表，记 `globalMap = Map<participant, InfoMap>` ；    
- 每一个参与者负责更新自己的 `localInfo`， 并由 `gossip` 协议负责将更新的信息同步到整个网络上；    
- 每个节点和系统中的某些节点成为 peer （如果系统规模比较小，则和系统中所有其他节点成为 peer）；    

 `gossip` 中有三种不同的同步信息方法：
- **push-gossip** -> 最简单的情况下， 一个节点 p 向 q 发送整个 `globalMap` ；    
- **pull-gossip** -> p 向 q 发送 digest ，q 根据 digest 向 p 发送 p 过期的 `(key, (value, version))` 列表；    
- **push-pull-gossip** -> 与 pull-gossip 类似，只是多了一步，p 再将本地比 q 新的数据推送给 q ，q 更新本地信息；    

## 特点

- gossip 不要求一个节点知道所有其他节点，因此具有**去中心化**的特点，节点之间完全对等，不需要任何的中心节点；    
- gossip 算法又被称为**反熵（Anti-Entropy）**；熵是物理学上的一个概念，代表杂乱无章，而反熵就是在杂乱无章中寻求一致，这充分说明了 gossip 的特点；     
- 在一个有界网络中，每个节点都随机地与其他节点通信，经过一番杂乱无章的通信，最终所有节点的状态都会达成一致；每个节点可能知道所有其他节点，也可能仅知道几个邻居节点，只要这些节点可以通过网络连通，最终他们的状态都是一致的；    
- gossip 算法是一个**最终一致性算法**，其无法保证在某个时刻所有节点状态一致，但可以保证”最终“所有节点会一致，”最终“是一个现实中存在、但理论上无法证明的时间点；    

## 协调机制

- **协调机制**是针对：在每次 2 个节点通信时，如何交换数据能最快的达到一致性，也即消除两个节点的不一致性；    
- 协调机制所面临的最大问题：因为受限于网络负载，不可能每次都把一个节点上的（全部）数据发送给另外一个节点，也即每个 gossip 的消息大小都有上限；在有限空间上、高效的交换所有消息是协调机制要解决的主要问题；    

在文章 **“[Efficient Reconciliation and Flow Control for Anti-Entropy Protocols](http://www.cs.cornell.edu/home/rvr/papers/flowgossip.pdf)”** 中描述了两种同步机制：

1. **precise reconciliation**
> precise reconciliation 希望在每次通信周期内都非常准确地消除双方的不一致性，具体表现为相互发送对方需要更新的数据； 然而，因为每个节点都在**并发**与多个节点通信，所以理论上很难做到上述要求。precise reconciliation 需要针对每个数据项独立地维护各自的 version，并在每次交互时，把所有的 `(key,value,version)` 发送到目标节点进行比对，从而找出双方不同之处进而更新。但因为 Gossip 消息存在大小限制，因此每次选择发送哪些数据就成了问题。当然，可以随机选择一部分数据，也可确定性的选择数据。对确定性的选择而言，可以有**最老优先**（根据版本）和**最新优先**两种：最老优先会优先更新版本最新的数据，而最新更新正好相反，这样会造成老数据始终得不到机会更新，也即饥饿。

2. **Scuttlebutt Reconciliation**
> Scuttlebutt Reconciliation 与 precise reconciliation 不同之处是，Scuttlebutt Reconciliation 不是为每个数据都维护单独的版本号，而是为每个节点上的宿主数据维护统一的 version 。比如节点 P 会为 `(p1,p2,...)` 维护一个一致的全局 version ，相当于把所有的宿主数据看作一个整体，当与其他节点进行比较时，只需比较这些宿主数据的最高 version ；如果最高 version 相同，则说明这部分数据全部一致，否则再进行 precise reconciliation 。

## Merkle tree

信息同步无疑是 gossip 的核心，Merkle tree 是一个非常适合同步的数据结构；

简单来说，Merkle tree 就是一颗 hash 树；在这棵树中，叶子节点的值是一些 hash 值、非叶节点的值均是由其子节点值计算 hash 得来的；因此，一旦某个文件被修改，修改信息（事件）就会迅速传播到 hash 树的根。需要（在发现差异后进行）同步的系统只需要不断查询根节点的 hash 值是否有变化，一旦发现变化，顺着树状结构就能够在 logN 级别的时间内找到发生变化的内容，马上同步；

在 Dynamo 中，每个节点保存一个范围内的 key 值，不同节点间存在一定范围的重叠 key 值；在去熵操作中，考虑的仅仅是某两个节点间共有 key 值范围；Merkle tree 的叶子节点就是这个共有 key 值范围内每个 key 的 hash 值；通过叶子节点的 hash 值，自底向上便可以构建出一棵 Merkle tree ；Dynamo 首先比对 Merkle tree 根处的 hash 值，如果一致，则表示两者完全相同，否则将其子节点交换并继续比较的过程；

## 总结

Gossip 常见于大规模、无中心的网络系统，可以用于众多能接受“**最终一致性**”的领域：失败检测、路由同步、Pub/Sub、动态负载均衡。

## 参考文献

- [Efficient Reconciliation and Flow Control for Anti-Entropy Protocols](http://www.cs.cornell.edu/home/rvr/papers/flowgossip.pdf)
- [Gossip 算法](http://tianya23.blog.51cto.com/1081650/530743)


原文：[这里](http://blog.csdn.net/yfkiss/article/details/6943682#comments)


----------


# Gossip算法

在公司的某个关于 Pub/Sub 的项目中使用到了 Gossip 算法。 为方便使用，整理网上找的资料如下：

Gossip 算法因为 Cassandra 而名声大噪，Gossip 看似简单，但要真正弄清楚其本质远没看起来那么容易。为了寻求 Gossip 的本质，下面的内容主要参考 Gossip 的原始论文：[“Efficient Reconciliation and Flow Control for Anti-Entropy Protocols“](http)

## Gossip背景

Gossip 算法如其名，灵感来自办公室八卦，只要一个人八卦一下，在有限的时间内所有的人都会知道该八卦的信息，这种方式也与病毒传播类似，因此 Gossip 有众多的别名：“闲话算法”、“疫情传播算法”、“病毒感染算法”、“谣言传播算法”；

但 Gossip 并不是一个新东西，之前的**泛洪查找**、**路由算法**都归属于这个范畴，不同的是 Gossip 给这类算法提供了明确的语义、具体实施方法及收敛性证明；

## Gossip特点

Gossip 算法又被称为**反熵（Anti-Entropy）**；**熵**是物理学上的一个概念，代表杂乱无章，而**反熵**就是在杂乱无章中寻求一致；这充分说明了 **Gossip 的特点**：在一个有界网络中，每个节点都随机地与其他节点通信，经过一番杂乱无章的通信，最终所有节点的状态都会达成一致。每个节点可能知道所有其他节点，也可能仅知道几个邻居节点，只要这些节可以通过网络连通，最终他们的状态都是一致的，当然这也是疫情传播的特点。

要注意到的一点是，即使有的节点因宕机而重启，有新节点加入，但经过一段时间后，这些节点的状态也会与其他节点达成一致，也就是说，Gossip天然具有分布式容错的优点。

## Gossip本质

Gossip 是一个带冗余的容错算法，更进一步，Gossip 是一个**最终一致性算法**；虽然无法保证在某个时刻所有节点状态一致，但可以保证在”最终“所有节点一致，”最终“是一个现实中存在，但理论上无法证明的时间点。

因为 Gossip 不要求节点知道所有其他节点，因此又具有去中心化的特点，节点之间完全对等，不需要任何的中心节点。实际上 Gossip 可以用于众多能接受“最终一致性”的领域：失败检测、路由同步、Pub/Sub、动态负载均衡。

但 Gossip 的**缺点**也很明显，冗余通信会对网路带宽、CPU 资源造成很大的负载，而这些负载又受限于通信频率，该频率又影响着算法收敛的速度，后面我们会讲在各种场合下的优化方法。

## Gossip节点的通信方式及收敛性

根据原论文，两个节点（A, B）之间存在三种通信方式：
- **push**: A 节点将数据及对应的版本号 (key,value,version) 推送给 B 节点，B 节点更新 A 中比自己新的数据；
- **pull**: A 仅将数据 (key,version) 推送给 B，B 将本地比 A 新的数据 (key,value,version) 推送给 A，A 更新本地；
- **push/pull**: 与 pull 类似，只是多了一步，A 再将本地比 B 新的数据推送给 B，B 更新本地；

如果把两个节点数据同步一次定义为一个周期，则在一个周期内，push 需通信 1 次，pull 需 2 次，push/pull 则需 3 次；因此从效果上来讲，push/pull 最好，理论上一个周期内可以使两个节点完全一致。直观上也感觉，push/pull 的收敛速度是最快的。

假设每个节点通信周期都能选择（感染）一个新节点，则 Gossip 算法退化为一个二分查找过程，每个周期构成一个平衡二叉树，收敛速度为 O(n2)，对应的时间开销则为 O(logn) 。这也是 Gossip 理论上最优的收敛速度。但在实际情况中最优收敛速度是很难达到的，假设某个节点在第 i 个周期被感染的概率为 pi ，第 i+1 个周期被感染的概率为 pi+1 ，则 pull 的方式:

![gossip_pull_probability](http "gossip_pull_probability")

而 push 为：

![gossip_push_probability](http "gossip_push_probability")

显然 pull 的收敛速度大于 push，而每个节点在每个周期被感染的概率都是固定的 p (0<p<1)，因此 Gossip 算法是基于 p 的平方收敛，也成为概率收敛，这在众多的一致性算法中是非常独特的。

Gossip 的节点的工作方式又分两种：
- **Anti-Entropy（反熵）**：以固定的概率传播所有的数据；
- **Rumor-Mongering（谣言传播）**：仅传播新到达的数据；

Anti-Entropy 模式有完全的容错性，但有较大的网络、CPU 负载；    
Rumor-Mongering 模式有较小的网络、CPU 负载，但必须为数据定义”最新“的边界，并且难以保证完全容错，对失败重启且超过”最新“期限的节点，无法保证最终一致性，或需要引入额外的机制处理不一致性；    

我们后续着重讨论 Anti-Entropy 模式的优化。

## Anti-Entropy 的协调机制

协调机制是讨论在每次两个节点通信时，如何交换数据能达到最快的一致性，也即消除两个节点的不一致性。上面所讲的 push、pull 等是通信方式，协调是在通信方式下的数据交换机制。协调所面临的最大问题是，因为受限于网络负载，不可能每次都把一个节点上的数据发送给另外一个节点，也即**每个 Gossip 的消息大小都有上限**，在有限的空间上、有效率地交换所有的消息是协调要解决的主要问题。

在讨论之前先声明几个概念：

令 `N = {p,q,s,...}` 为需要 gossip 通信的 server 集合，有界大小；
令 `(p1,p2,...)` 是宿主在节点 p 上的数据，其中数据由 `(key,value,version)` 构成，q 的规则与 p 类似；
为了保证一致性，规定数据的 value 及 version 只有宿主节点才能修改，其他节点只能通过 Gossip 协议间接请求数据对应的宿主节点修改；

### 精确协调（Precise Reconciliation）

精确协调希望在每次通信周期内都非常准确地消除双方的不一致性，具体表现为相互发送对方需要更新的数据，因为每个节点都在并发与多个节点通信，理论上精确协调很难做到。精确协调需要给每个数据项独立地维护自己的 version，在每次交互是把所有的 (key,value,version) 发送到目标进行比对，从而找出双方不同之处从而更新。但因为 Gossip 消息存在大小限制，因此每次选择发送哪些数据就成了问题。当然可以随机选择一部分数据，也可确定性的选择数据。对确定性的选择而言，可以有最老优先（根据版本）和最新优先两种，最老优先会优先更新版本最新的数据，而最新更新正好相反，这样会造成老数据始终得不到机会更新，也即饥饿。

当然，开发这也可根据业务场景构造自己的选择算法，但始终都无法避免消息量过多的问题。

### 整体协调（Scuttlebutt Reconciliation）

整体协调与精确协调不同之处是，整体协调不是为每个数据都维护单独的版本号，而是为每个节点上的宿主数据维护统一的 version。比如节点 P 会为 (p1,p2,...) 维护一个一致的全局 version，相当于把所有的宿主数据看作一个整体，当与其他节点进行比较时，只需必须这些宿主数据的最高 version，如果最高 version 相同说明这部分数据全部一致，否则再进行精确协调。

整体协调对数据的选择也有两种方法：
- 广度优先：根据整体 version 大小排序，也称为公平选择；
- 深度优先：根据包含数据多少的排序，也称为非公平选择。因为后者更有实用价值，所以原论文更鼓励后者；

## Cassandra中的实现

经过验证，Cassandra 实现了基于整体协调的 push/push 模式，有几个组件：

三条消息分别对应 push/pull 的三个阶段：
- GossipDigitsMessage
- GossipDigitsAckMessage
- GossipDigitsAck2Message

还有三种状态：
- EndpointState：维护宿主数据的全局 version，并封装了 HeartBeat 和 ApplicationState ；
- HeartBeat：心跳信息；
- ApplicationState：系统负载信息（磁盘使用率）；

Cassandra 主要是使用 Gossip 完成三方面的功能：
- 失败检测
- 动态负载均衡
- 去中心化的弹性扩展

## 总结

Gossip 是一种去中心化、容错而又最终一致性的绝妙算法，其收敛性不但得到证明还具有指数级的收敛速度。使用 Gossip 的系统可以很容易的把 Server 扩展到更多的节点，满足弹性扩展轻而易举。

唯一的缺点是收敛是最终一致性，不使用那些强一致性的场景，比如 2pc 。