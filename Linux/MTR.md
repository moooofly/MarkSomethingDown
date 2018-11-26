# MTR

> MTR is a powerful **network diagnostic tool** that enables administrators to **diagnose** and **isolate** networking errors and provide helpful reports of network status to upstream providers. MTR represents an evolution of the `traceroute` command by providing a greater data sample, as if augmenting `traceroute` with `ping` output. This document provides an in depth overview of MTR, the data it generates, and how to properly interpret and draw conclusions based on the data provided by it.

MTR 是网络诊断工具；MTR 可以看作 `traceroute` 和 `ping` 的组合；

## Network Diagnostics Background

> Networking diagnostic tools including `ping`, `traceroute`, and `mtr` use “**ICMP**” packets to test contention and traffic between two points on the Internet. When a user pings a host on the Internet, a series of ICMP packets are sent to the host, which responds by sending packets in return. The user’s client is then able to compute the **round trip time** between two points on the Internet.

`ping`, `traceroute` 和 `mtr` 默认都使用 ICMP 包进行**拥塞**和**网络状况**测试；

> By contrast, tools such as `traceroute` and MTR send ICMP packets with incrementally increasing TTLs in order to view the route or series of hops that the packet makes between the origin and its destination. The TTL, or time to live, controls how many “hops” a packet will make before “dying” and returning to the host. By sending a series of packets and causing them to die and return after one hop, then two, then three, the client machine is able to assemble the route that traffic takes between hosts on the Internet.

`traceroute` 和 MTR 会使用逐渐增加的 TTL 的策略发送 ICMP ，以便定位出源端和目的端之间的路由关系；

> Rather than provide a simple outline of the route that traffic takes across the Internet, MTR collects additional information regarding the state, connection, and responsiveness of the **intermediate hosts**. Because of this additional information, it is recommended that you use MTR whenever possible to provide the most complete overview of the connection between two hosts on the Internet. The following sections outline how to install the MTR software and how to interpret the results provided by this tool.

MTR 的牛逼之处在于**能够提供关于中间网络设备状况的统计信息**；

## Generating an MTR Report

> Because MTR provides an image of the route traffic takes from one host to another, you can think of it as a **directional** tool. Furthermore, the route taken between two points on the Internet can vary a great deal based on location and the routers that are located upstream of you. For this reason it is often recommended that you collect MTR reports in both directions for all hosts that are experiencing connectivity issues, or as many hosts as possible.

应该认为 MTR 是有方向性的；建议在所有存在连接性问题的主机上进行两个方向的 MTR 报告输出；

> Linode support will often request “mtr reports” both to and from your Linode if you are experiencing networking issues. This is because, from time to time, MTR reports will not point to errors from one direction when there is still packet loss from the opposite direction. Having both reports is helpful as it can aid in the identification of issues and will be needed if a problem must be reported.

一个方向上 MTR 报告正常，另外一个方向上报告不正常是可能的；因此需要两个方向的报告；

> For the sake of clarity, when referring to MTR reports this document refers to the host running mtr as the source host and the host targeted by the query as the destination host.

## Using MTR on Unix-based Systems

可以通过如下命令生成 MTR 报告：

```
mtr -rw [destination_host]
```

如果报告中显示没有 packet loss 发生，则可以以 **faster interval** 再运行一次：

```
mtr -rwc 50 -i 0.2 -rw 12.34.56.78
```

参数说明：

- The `r` option flag generates the report (short for `--report`).
- The `w` option flag uses the long-version of the hostname so our technicians and you can see the full hostname of each hop (short for `--report-wide`).
- The `c` option flag sets how many packets are sent and recorded in the report. When not used, the **default** will generally be **10**, but for faster intervals you may want to set it to **50** or **100**. The report can take longer to finish when doing this.
- The `i` option flag runs the report at a faster rate to reveal packet loss that can occur only during network congestion. This flag instructs MTR to **send one packet every n seconds**. The default is 1 second, so setting it to a few tenths of a second (0.1, 0.2, etc.) is generally helpful.

## Reading MTR Reports

MTR 报告中包含了大量有价值的信息，以如下输出为例进行说明（本地连接 google.com）：

```
$ mtr --report google.com
HOST: example                  Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. inner-cake                    0.0%    10    2.8   2.1   1.9   2.8   0.3
  2. outer-cake                    0.0%    10    3.2   2.6   2.4   3.2   0.3
  3. 68.85.118.13                  0.0%    10    9.8  12.2   8.7  18.2   3.0
  4. po-20-ar01.absecon.nj.panjde  0.0%    10   10.2  10.4   8.9  14.2   1.6
  5. be-30-crs01.audubon.nj.panjd  0.0%    10   10.8  12.2  10.1  16.6   1.7
  6. pos-0-12-0-0-ar01.plainfield  0.0%    10   13.4  14.6  12.6  21.6   2.6
  7. pos-0-6-0-0-cr01.newyork.ny.  0.0%    10   15.2  15.3  13.9  18.2   1.3
  8. pos-0-4-0-0-pe01.111eighthav  0.0%    10   16.5  16.2  14.5  19.3   1.3
  9. as15169-3.111eighthave.ny.ib  0.0%    10   16.0  17.1  14.2  27.7   3.9
 10. 72.14.238.232                 0.0%    10   19.1  22.0  13.9  43.3  11.1
 11. 209.85.241.148                0.0%    10   15.1  16.2  14.8  20.2   1.6
 12. lga15s02-in-f104.1e100.net    0.0%    10   15.6  16.9  15.2  20.6   1.7
```

该测试中会发送 10 packets 到 google.com 并产生输出信息；在未指定 `--report` 选项时，`mtr` 将在交互环境下持续运行；在交互模式中会展现当前出到达每一个主机的 RTT ；

报告中包含了 12 hops ；“Hops” 对应了网络上的 **nodes** 或 **routers** ，即 packets 到达目的地前所经过的设备；在上例中，packets 经过了 “inner-cake” 和 “outer-cake” 两个本地网络设备，之后到达 “68.85.118.13” 和一系列命名主机；这些主机的名字由反向 DNS 查询所确定；可以通过指定 `--no-dns` 选项取消 rDNS 查询；

```
% mtr --no-dns --report google.com
HOST: deleuze                     Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 192.168.1.1                   0.0%    10    2.2   2.2   2.0   2.7   0.2
  2. 68.85.118.13                  0.0%    10    8.6  11.0   8.4  17.8   3.0
  3. 68.86.210.126                 0.0%    10    9.1  12.1   8.5  24.3   5.2
  4. 68.86.208.22                  0.0%    10   12.2  15.1  11.7  23.4   4.4
  5. 68.85.192.86                  0.0%    10   17.2  14.8  13.2  17.2   1.3
  6. 68.86.90.25                   0.0%    10   14.2  16.4  14.2  20.3   1.9
  7. 68.86.86.194                  0.0%    10   17.6  16.8  15.5  18.1   0.9
  8. 75.149.230.194                0.0%    10   15.0  20.1  15.0  33.8   5.6
  9. 72.14.238.232                 0.0%    10   15.6  18.7  14.1  32.8   5.9
 10. 209.85.241.148                0.0%    10   16.3  16.9  14.7  21.2   2.2
 11. 66.249.91.104                 0.0%    10   22.2  18.6  14.2  36.0   6.5
```

- `Loss%` 列展示了每一 hop 上的 packet loss 百分比；
- `Snt` 列统计了 packets sent 的数量；
- `--report` 默认发送 10 packets 除非通过 `--report-cycles=[number-of-packets]` 另外指定；
- `Last`, `Avg`, `Best`, and `Wrst` 全部时以毫秒 (ms) 为单位的延迟度量；
- `Last` 是 last packet sent 的 latency ；
- `Avg` 是 all packets 的平均 latency ；
- `Best` 和 `Wrst` 展示了一个数据包的 best (shortest) 和 worst (longest) RTT 数值；
- 在大多数情况下，`Avg` 列应该是你的主要关注点；
- `StDev` 提供了针对每个主机的 standard deviation of latency；standard deviation 的数值越大，表示 measurements 之间的 latency 差值越大；如果均值恰好位于数据集中值，或者由于某种现象或 measurement 错误而导致数据不准，则可以使用 Standard deviation ；例如，如果 standard deviation 很高，则表示 latency 测量值非常不一致（起伏很大）；尽管其中某些值可能很低（例如 25ms），其它值可能非常高（例如 350ms）：在对 10 packets 的延迟求平均后，均值看起来可能很正常，但事实上却无法很好的代表数据的实际情况；如果 standard deviation 很高，则可以查看下 best 和 worst latency 度量值，以确保均值能够很好的表示真实 latency 而不是大量波动产生的效果；

在大多数情况下，你可以认为 MTR 输出主要分为三个段：取决于具体配置，前 2 或 3 hops 通常代表源主机的 ISP ，而最后 2 或 3 hops 则代表目的主机的 ISP ；在这两者之间的 hops 则对应了 packet 传递过程中所经过的 routers ；

例如，如果 MTR 运行在你的 home PC 上，目的地为你的 Linode 主机，则前 2 或 3 hops 属于你的 ISP ；最后的 3 hops 则属于你的 Linode 所位于的数据中心；位于中间位置的任何 hops 均为 intermediate hops ；当你在本地运行 MTR 时，如果你在 source 附近的前几 hops 上发现异常，则联系你本地服务提供商，或者调查一下你的本地网络配置情况；相反的，如果你看到 destination 端有异常，你可能会要联系目标服务器的管理员，或者目标机器的网络支持人员（例如 Linode）；不幸的是，当问题出现在中间 hops 时，两端的服务提供者在处理问题时的能力有限；

## Analyzing MTR Reports

### Verifying Packet Loss

当分析 MTR 输出时，主要寻找两方面内容：**loss** 和 **latency** ；首先，我们先讨论 loss 问题；如果你看到在任意的 hop 上存在一定百分比的 loss ，则可能就意味着对应的特定 router 存在问题；然而，还有一种常见的情况，就是一些服务提供商会针对 ICMP traffic 进行速率限制，而 ICMP 正是 MTR 所使用的；这可能会给出 packet loss 的假象，而实际上并没有丢包发生；为了确定是否真的存在丢包，还是由速率限制导致的丢包假象，可以看看后续 hop 的情况；如果后续的 hop 显示 loss 为 0.0% ，则你可以确信针对 ICMP 的 rate limiting 确实在发生，而非真正的丢包；详见下面的例子：

```
root@localhost:~# mtr --report www.google.com
HOST: example               Loss%   Snt   Last   Avg  Best  Wrst StDev
1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
2. 63.247.64.157                50.0%    10    0.4   1.0   0.4   6.1   1.8
3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
5. 72.14.233.56                  0.0%    10    7.2   8.3   7.1  16.4   2.9
6. 209.85.254.247                0.0%    10   39.1  39.4  39.1  39.7   0.2
7. 64.233.174.46                 0.0%    10   39.6  40.4  39.4  46.9   2.3
8. gw-in-f147.1e100.net          0.0%    10   39.6  40.5  39.5  46.7   2.2
```

在这个例子中，在 hops 1 和 2 之间报告出的丢包情况，非常可能是由于第二 hop 上进行了 rate limiting 导致；尽管到达剩余 hops 的 traffic 均经过第二 hop ，但都没有出现 packet loss 的情况；如果 loss 连续发生在不止一个 hop 上，则可能存在 packet loss 或路由问题；需要知道的是，rate limiting 和 loss 可能会同时发生；在这种情况下，可以将一组 loss 值中的最低百分比值当作实际的 loss 值；例如，考虑如下输出：

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
1. 63.247.74.43                   0.0%    10    0.3   0.6   0.3   1.2   0.3
2. 63.247.64.157                  0.0%    10    0.4   1.0   0.4   6.1   1.8
3. 209.51.130.213                60.0%    10    0.8   2.7   0.8  19.0   5.7
4. aix.pr1.atl.google.com        60.0%    10    6.7   6.8   6.7   6.9   0.1
5. 72.14.233.56                  50.0%   10    7.2   8.3   7.1  16.4   2.9
6. 209.85.254.247                40.0%   10   39.1  39.4  39.1  39.7   0.2
7. 64.233.174.46                 40.0%   10   39.6  40.4  39.4  46.9   2.3
8. gw-in-f147.1e100.net          40.0%   10   39.6  40.5  39.5  46.7   2.2
```

在这个场景下，你能看到 60% 的 loss 发生在 hops 2 和 3 之间，以及 hops 3 和 4 之间；由此你可以认为 3 和 4 hop 之间可能发生了 traffic 的丢失，因为后续 host 的报告中没有再出现 zero loss 的情况；然而，其中一些 loss 应该是由于 rate limiting 的原因，因为最后几个 hops 仅有 40% 的 loss 值；当存在不同的 loss 值被输出时，一个原则就是越位于后面的 hop 输出越应该被相信；

有些 loss 可以解释为在回程路由 (return route) 中发生的问题；Packets 在到达其目的地过程中没有发生错误，但是在回程时却遇到了问题；这在报告中会很明显，但缺很难从 MTR 的输出中推断出来；因此，通常都建议从两个方向上获取 MTR 报告；

> Additionally, resist the temptation to investigate or report all incidences of packet loss in your connections. The Internet protocols are designed to be resilient to some network degradation, and the routes that data takes across the Internet can fluctuate in response to load, brief maintenance events, and other routing issues. If your MTR report shows small amounts of loss in the neighborhood of 10%, there is no cause for real concern as the application layer will compensate for the loss which is likely transient.

- 偶发性丢包是正常的；
- 协议本身已被设计为对网络降级情况能够弹性应对；
- 数据途径的路由针对各种实际情况会出现波动是正常的；
- 若 MTR 报告中给出少量的 loss ，例如 10% 左右，则无需担心；


### Understanding Network Latency

> In addition to helping you assess packet loss, MTR will also help you assess the latency of a connection between your host and the target host. By virtue of physical constraints, latency always increases with the number of hops in a route. However, the increases should be consistent and linear. Unfortunately, latency is often relative and very dependent on the quality of both host’s connections and their physical distance. When evaluating MTR reports for potentially problematic connections, consider earlier fully functional reports as context in addition to known connection speeds between other hosts in a given area.

- latency 总是随着路由中包含的 hop 数量增长；
- latency 的增长应该是一致的和线性的；
- latency 通常是相对的，并且依赖于主机连接的质量和物理距离；
- 当基于 MTR 报告评估潜在的连接问题时，可以参考之前获取的全功能报告，以及主机之间已知连接速度作为参考和对比；

> The connection quality may also affect the amount of latency you experience for a particular route. Predictably, dial-up connections will have much higher latency than cable modem connections to the same destination. Consider the following MTR report which shows a high latency:

连接的质量同样会影响延迟的具体数值；下面给出的是 high latency 示例：

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
  2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
  3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
  4. aix.pr1.atl.google.com        0.0%    10  388.0 360.4 342.1 396.7   0.2
  5. 72.14.233.56                  0.0%    10  390.6 360.4 342.1 396.7   0.2
  6. 209.85.254.247                0.0%    10  391.6 360.4 342.1 396.7   0.4
  7. 64.233.174.46                 0.0%    10  391.8 360.4 342.1 396.7   2.1
  8. gw-in-f147.1e100.net          0.0%    10  392.0 360.4 342.1 396.7   1.2
```

> The amount of latency jumps significantly between hops 3 and 4 and remains high. This may point to a network latency issue as round trip times remain high after the fourth hop. From this report, it is impossible to determine the cause although a saturated peering session, a poorly configured router, or a congested link are frequent causes.

- 在 hops 3 和 4 之间 latency 发生了跳变，并在后续 hop 上保持高值，代表存在网络延迟问题；
- 从上述报告中，无法确定根因；
- 导致延迟的常见的情况有：
    - **a saturated peering session**
    - **a poorly configured router**
    - **a congested link**

> Unfortunately, high latency does not always mean a problem with the current route. A report like the one above means that despite some sort of issue with the 4th hop, traffic is still reaching the destination host and returning to the source host. Latency could be caused by a problem with the return route as well. The return route will not be seen in your MTR report, and packets can take completely different routes to and from a particular destination.

- 高 latency 并不总是表示当前路由存在问题；
- 上述报告表明：尽管 hop 4 存在某种问题，但是 traffic 仍能到达目的主机，并返回源主机；
- Latency 可能由返程路由 (return route) 的问题引起；而返程路由的问题在上述 MTR 报告中是看不出来的，因为网络包可能走的是完全不同的路由；

> In the above example, while there is a large jump in latency between hosts 3 and 4 the latency does not increase unusually in any subsequent hops. From this it is logical to assume that there is some issue with the 4th router.

在上面的例子中，尽管在 hop 3 和 4 之间发生了 latency 的剧烈跳变，但是后续 hop 中并没有再次出现异常的增长；从这点来看，认为 hop 4 上的 router 存在问题是符合逻辑的；

> ICMP rate limiting can also create the appearance of latency, similar to the way that it can create the appearance of packet loss. Consider the following example:

针对 ICMP 的 rate limiting 同样会导致 latency 的出现，和导致 packet loss 的出现如出一辙；详见如下示例：

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
  2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
  3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
  4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
  5. 72.14.233.56                  0.0%    10  254.2 250.3 230.1 263.4   2.9
  6. 209.85.254.247                0.0%    10   39.1  39.4  39.1  39.7   0.2
  7. 64.233.174.46                 0.0%    10   39.6  40.4  39.4  46.9   2.3
  8. gw-in-f147.1e100.net          0.0%    10   39.6  40.5  39.5  46.7   2.2
```

> At first glance, the latency between hops 4 and 5 draws attention. However after the fifth hop, the latency drops drastically. The actual latency measured here is about 40ms. In cases like this, MTR draws attention to an issue which does not affect the service. Consider the latency to the final hop when evaluating an MTR report.

- 只看第一眼，hops 4 和 5 之间的 latency 就能够引起注意；但是，在 hop 5 之后，latency 的值又极大的降低了，实际测量到的 latency 大约在 40ms 左右；
- 在这个例子中，MTR 将我们的注意引到了一个“问题”上，而这个“问题”却不会影响服务；
- 建议在评估 MTR 报告时，重点考虑最后 hop 给出的 latency 值；

## Common MTR Reports

> Some networking issues are novel and require escalation to the operators of the upstream networks. However, there are a selection of common MTR reports that describe common networking issues. If you’re experiencing some sort of networking issue and want to diagnose your problem, consider the following examples.

有一些网络问题是很新奇的，需要位于上游的网络管理人员的关注和支持才能定位和解决；然而，存在一些通用的 MTR 报告描述了一些常见网络问题；

### Destination Host Networking Improperly Configured

> In the next example, it appears that there is 100% loss to a the destination host because of an incorrectly configured router. At first glance it appears that the packets are not reaching the host but this is not the case.

从如下输出中可以看到，在到达目的主机的最后 hop 上存在 100% 丢包情况，正是由于错误的路由配置导致；第一眼看来，似乎数据包没有到达主机，然而这并非真实情况；

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
  2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
  3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
  4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
  5. 72.14.233.56                  0.0%    10    7.2   8.3   7.1  16.4   2.9
  6. 209.85.254.247                0.0%    10   39.1  39.4  39.1  39.7   0.2
  7. 64.233.174.46                 0.0%    10   39.6  40.4  39.4  46.9   2.3
  8. gw-in-f147.1e100.net         100.0    10    0.0   0.0   0.0   0.0   0.0
```

> The traffic does reach the destination host however, the MTR report shows loss because the destination host is not sending a reply. This may be the result of improperly configured networking or firewall (iptables) rules that cause the host to drop ICMP packets.

事实上，数据包确实到达了目的主机，然而 MTR 报告却显示有丢包，根本原因在于目的主机没有发送应答信息；这可能是由于错误配置网络或防火墙 (iptables) 导致；

> The way you can tell that the loss is due to a misconfigured host is to look at the hop which shows 100% loss. From previous reports, you see that this is the final hop and that MTR does not try additional hops. While it is difficult to isolate this issue without a baseline measurement, these kinds of errors are quite common.

能够确认丢包是由于主机错误配置导致的办法是，查看给出 100% 丢包的 hop ；从报告中可以知道，这已经是最后的 hop 了，由于 MTR 报告中没有额外 hop 信息可以提供，因此难以将该问题进行隔离；虽然无法得到一个基准线进行比对，但这种类型的错误是经常会发生的；

### Residential or Business Router

> Oftentimes residential gateways will cause MTR reports to look a little misleading.

常常会听说，家庭住宅用的网关经常会导致 MTR 报告看起来很奇怪；

```
% mtr --no-dns --report google.com
HOST: deleuze                     Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 192.168.1.1                   0.0%    10    2.2   2.2   2.0   2.7   0.2
  2. ???                          100.0    10    8.6  11.0   8.4  17.8   3.0
  3. 68.86.210.126                 0.0%    10    9.1  12.1   8.5  24.3   5.2
  4. 68.86.208.22                  0.0%    10   12.2  15.1  11.7  23.4   4.4
  5. 68.85.192.86                  0.0%    10   17.2  14.8  13.2  17.2   1.3
  6. 68.86.90.25                   0.0%    10   14.2  16.4  14.2  20.3   1.9
  7. 68.86.86.194                  0.0%    10   17.6  16.8  15.5  18.1   0.9
  8. 75.149.230.194                0.0%    10   15.0  20.1  15.0  33.8   5.6
  9. 72.14.238.232                 0.0%    10   15.6  18.7  14.1  32.8   5.9
 10. 209.85.241.148                0.0%    10   16.3  16.9  14.7  21.2   2.2
 11. 66.249.91.104                 0.0%    10   22.2  18.6  14.2  36.0   6.5
```

> Do not be alarmed by the 100% loss reported. This does not indicate that there is a problem. You can see that there is no loss on subsequent hops.

不要被其中的 100% 丢包所吓到，此处的数值并不表明存在问题；你可以看到后续 hop 上并没有丢包；

### An ISP Router Is Not Configured Properly

> Sometimes a router on the route your packet takes is incorrectly configured and your packets may never reach their destination. Consider the following example:

有些时候会遇到数据包路径上的 router 没有配置正确的情况，由此导致数据包到达不了目的地；

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
  2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
  3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
  4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
  5. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
  6. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
  7. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
  8. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
  9. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
 10. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
```

> The question marks appear when there is no additional route information. The following report displays the same issue:

??? 出现的原因正是由于没有获取到额外的路由信息导致的；

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
   1. 63.247.74.43                 0.0%    10    0.3   0.6   0.3   1.2   0.3
   2. 63.247.64.157                0.0%    10    0.4   1.0   0.4   6.1   1.8
   3. 209.51.130.213               0.0%    10    0.8   2.7   0.8  19.0   5.7
   4. aix.pr1.atl.google.com       0.0%    10    6.7   6.8   6.7   6.9   0.1
   5. 172.16.29.45                 0.0%    10    0.0   0.0   0.0   0.0   0.0
   6. ???                          0.0%    10    0.0   0.0   0.0   0.0   0.0
   7. ???                          0.0%    10    0.0   0.0   0.0   0.0   0.0
   8. ???                          0.0%    10    0.0   0.0   0.0   0.0   0.0
   9. ???                          0.0%    10    0.0   0.0   0.0   0.0   0.0
  10. ???                          0.0%    10    0.0   0.0   0.0   0.0   0.0
```

> Sometimes, a poorly configured router will send packets in a loop. You can see that in the following example:

有时错误配置的 router 可能会导致数据包在环形路径中被发送；

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
  2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
  3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
  4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
  5. 12.34.56.79                   0.0%    10    0.0   0.0   0.0   0.0   0.0
  6. 12.34.56.78                   0.0%    10    0.0   0.0   0.0   0.0   0.0
  7. 12.34.56.79                   0.0%    10    0.0   0.0   0.0   0.0   0.0
  8. 12.34.56.78                   0.0%    10    0.0   0.0   0.0   0.0   0.0
  9. 12.34.56.79                   0.0%    10    0.0   0.0   0.0   0.0   0.0
 10. 12.34.56.78                   0.0%    10    0.0   0.0   0.0   0.0   0.0
 11. 12.34.56.79                   0.0%    10    0.0   0.0   0.0   0.0   0.0
 12. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
 13. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
 14. ???                           0.0%    10    0.0   0.0   0.0   0.0   0.0
```

> All of these reports show that the router at hop 4 is not properly configured. When these situations happen, the only way to resolve the issue is to contact the network administrator’s team of operators at the source host.

报告信息表明，位于 hop 4 的 router 配置存在问题；此时唯一的办法就是联系相应的负责人员；


### ICMP Rate Limiting

> **ICMP rate limiting** can cause apparent packet loss as described below. When there is packet loss to one hop that doesn’t persist to subsequent hops, the loss is caused by ICMP limiting. See the following example:

ICMP rate limiting 能够导致明显的 packet loss 出现；当只有单一 hop 出现 packet loss ，而后续 hop 上并没有 packet loss 时，该 packet loss 可能就是由于 ICMP limiting 导致；

```
root@localhost:~# mtr --report www.google.com
 HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
   1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
   2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
   3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
   4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
   5. 72.14.233.56                 60.0%    10   27.2  25.3  23.1  26.4   2.9
   6. 209.85.254.247                0.0%    10   39.1  39.4  39.1  39.7   0.2
   7. 64.233.174.46                 0.0%    10   39.6  40.4  39.4  46.9   2.3
   8. gw-in-f147.1e100.net          0.0%    10   39.6  40.5  39.5  46.7   2.2
```

> In situations like this there is no cause for concern. Rate limiting is a common practice and it reduces congestion to prioritizes more important traffic.

在上述情况下，没有什么具体问题可担心；因为 Rate limiting 作为一种常规实践，能够减少拥塞以保障更重要的网络流量的处理优先级；

### Timeouts

> Timeouts can happen for various reasons. Some routers will discard ICMP and no replies will be shown on the output as **timeouts (???)**. Alternatively there may be a problem with the return route:

- 超时无处不在；
- 一些路由器会丢弃 ICMp 报，因此在输出中不会有回复信息被显示（此时显示为 ???）；
- 还有一种可能，就是返程路由确实存在问题；

```
root@localhost:~# mtr --report www.google.com
HOST: localhost                   Loss%   Snt   Last   Avg  Best  Wrst StDev
  1. 63.247.74.43                  0.0%    10    0.3   0.6   0.3   1.2   0.3
  2. 63.247.64.157                 0.0%    10    0.4   1.0   0.4   6.1   1.8
  3. 209.51.130.213                0.0%    10    0.8   2.7   0.8  19.0   5.7
  4. aix.pr1.atl.google.com        0.0%    10    6.7   6.8   6.7   6.9   0.1
  5. ???                           0.0%    10    7.2   8.3   7.1  16.4   2.9
  6. ???                           0.0%    10   39.1  39.4  39.1  39.7   0.2
  7. 64.233.174.46                 0.0%    10   39.6  40.4  39.4  46.9   2.3
  8. gw-in-f147.1e100.net          0.0%    10   39.6  40.5  39.5  46.7   2.2
```

> Timeouts are not necessarily an indication of packet loss. Packets still reach their destination without significant packet loss or latency. Timeouts may be attributable to routers dropping packets for QoS (quality of service) purposes or there may be some issue with return routes causing the timeouts. This is another false positive.

- 超时本身不能作为得出 packet loss 结论的必要条件；
- 包在到达目的地时可能并没有严重丢包或延迟，超时原因可能就是 router 基于 QoS 的缘故主动丢弃包导致；也可能是由于返程路由存在问题导致；

## Advanced MTR techniques

> Newer versions of MTR are now capable of running in **TCP mode** on a specified TCP port, compared to the default use of the ICMP (ping) protocol. In some instances network degradation will only affect certain ports or misconfigured firewall rules on a router may block a certain protocol. Running MTR over a certain port can show packet loss where the default ICMP report may not.

- 更新版本的 MTR 已经支持了 TCP 模式；
- 在有些场景中，网络降级问题只会影响一些特定的端口，而 router 上错误配置的防火墙规则可能只会阻塞特定的协议；
- 因此在特定端口上运行 MTR 可能会显示 packet loss ，但基于默认的 ICMP 可能看不出来；

> Running MTR in TCP mode will require super-user privileges on most machines:

```
sudo mtr -P 80 -i 0.5 -rw50 example.com
sudo mtr -P 22 -i 0.5 -rw50 example.com
```

## Resolving Routing and Networking Issues Identified in your MTR report

> A majority of routing issues displayed by MTR reports are temporary. Most issues will clear up by themselves within 24 hours. In most cases, by the time you are able to notice a problem with a route, the Internet service provider’s monitoring has already reported the problem and administrators are working to fix the issue. In cases where you are experiencing degraded service for an extended period of time, you may choose to alert a provider of the issues you’re experiencing. When contacting a service provider, send MTR reports and any other relevant data you may have. Without usable data, providers have no way to verify or fix problems.

这里纯扯蛋；

> While routing errors and issues account for a percentage of network-related slowness, they are by no means the only cause of degraded performance. Network congestion, particularly over long distances during peak times, can become quite congested. Transatlantic and transpacific traffic can be quite variable and are subject to general network congestion. In these cases, it is recommended that you position hosts and resources as geographically close to their targeted audience as possible.

- 尽管 routing errors 和 issues 占据网络相关慢速问题的一定百分比，但其绝非导致性能降级的唯一原因；网络拥塞，尤其是面对高峰期长距离传输时，可能会更加严重；
- 跨大西洋和跨太平洋的网络通信状况复杂多变，常常会遭遇网络拥塞问题；在这种情况下，建议将主机和相关资源布置在地理位置尽可能贴近目标用户的位置；

> If you are experiencing connectivity issues and are unable to interpret your MTR report, you may open a support ticket including the output of your report in the “Support” tab of the Linode Manager and our technicians can help analyze your issue.

若对 MTR 报告存在不理解，可以求助；


----------


## 参考

- [Diagnosing Network Issues with MTR](https://www.linode.com/docs/networking/diagnostics/diagnosing-network-issues-with-mtr/)

