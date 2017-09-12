
当 linux 系统上有多个单独网卡时，若想充分利用这些网卡，同时能够对外提供统一的网络地址，增大网络的吞吐量，提高网络的可用性，这时就需要 bond 来帮助我们解决这个问题。

bond 的配置实在很简单，但是配置不好，很容易造成严重的网络问题。bonding 功能由 linux 内核自带了，因此通常不需要安装，只需要把 bonding 模块加载到内核里即可。

# 七种 bond 模式说明

```
#define BOND_MODE_ROUNDROBIN 0   （balance-rr模式）网卡的负载均衡模式 
#define BOND_MODE_ACTIVEBACKUP 1 （active-backup模式）网卡的容错模式 
#define BOND_MODE_XOR 2          （balance-xor模式）需要交换机支持 
#define BOND_MODE_BROADCAST 3    （broadcast模式） 
#define BOND_MODE_8023AD 4       （IEEE 802.3ad动态链路聚合模式）需要交换机支持 
#define BOND_MODE_TLB 5           自适应传输负载均衡模式 
#define BOND_MODE_ALB 6           网卡虚拟化方式 
```

## 第一种模式：mod=0 ，即：(balance-rr) Round-robin policy（平衡轮询策略）

特点：数据包传输顺序是依次传输（即第 1 个包走 eth0 ，下一个包就走 eth1 ，一直循环下去，直到最后一个传输完毕），此模式提供**负载平衡**和**容错能力**；但是我们知道如果一个连接或者会话的数据包从不同的接口发出的话，若中途再经过不同的链路，在客户端很有可能会出现**数据包无序到达的问题**，而无序到达的数据包需要重新要求被发送，这样网络的吞吐量就会下降；

## 第二种模式：mod=1，即： (active-backup) Active-backup policy（主-备份策略）

特点：**只有一个设备处于活动状态**，当一个宕掉另一个马上由备份转换为主设备。mac 地址是外部可见的，**从外面看来，bond 的 MAC 地址是唯一的**，以避免 switch (交换机)发生混乱。此模式只提供了**容错能力**；由此可见，此算法的优点是可以提供高网络连接的可用性，但是它的资源利用率较低，只有一个接口处于工作状态，在有 N 个网络接口的情况下，资源利用率为 1/N ；

## 第三种模式：mod=2，即：(balance-xor) XOR policy（平衡策略）

特点：基于指定的**传输 hash 策略**传输数据包。

缺省策略是：

```
(源 MAC 地址 XOR 目标 MAC 地址) % slave 数量
```

其他的传输策略可以通过 xmit_hash_policy 选项指定，此模式提供**负载平衡**和**容错能力**；

## 第四种模式：mod=3，即：broadcast（广播策略）

特点：在每个 slave 接口上传输每个数据包，此模式提供了**容错能力**

## 第五种模式：mod=4，即：(802.3ad) IEEE 802.3ad Dynamic link aggregation（IEEE 802.3ad 动态链接聚合）

特点：创建一个**聚合组**，它们共享同样的速率和双工设定。根据 802.3ad 规范将多个 slave 工作在同一个激活的聚合体下。

外出流量的 slave 选举是基于传输 hash 策略，该策略可以通过 xmit_hash_policy 选项从缺省的 XOR 策略改变到其他策略。需要注意的是，并不是所有的传输策略都是 802.3ad 适应的，尤其考虑到在 802.3ad 标准 43.2.4 章节提及的包乱序问题。不同的实现可能会有不同的适应性。

必要条件：

- 条件1：ethtool 支持获取每个 slave 的速率和双工设定；
- 条件2：switch (交换机)支持 IEEE 802.3ad Dynamic link aggregation；
- 条件3：大多数 switch(交换机)需要经过特定配置才能支持 802.3ad 模式；

## 第六种模式：mod=5，即：(balance-tlb) Adaptive transmit load balancing（适配器传输负载均衡）

特点：不需要任何特别的 switch (交换机)支持的**通道 bonding** 。在每个 slave 上根据当前的负载（根据速度计算）分配外出流量。如果正在接受数据的 slave 出故障了，另一个 slave 接管失败的 slave 的 MAC 地址。

该模式的必要条件：ethtool 支持获取每个 slave 的速率；

## 第七种模式：mod=6，即：(balance-alb) Adaptive load balancing（适配器适应性负载均衡）

特点：该模式包含了 balance-tlb 模式，同时加上**针对 IPV4 流量的接收负载均衡**(receive load balance, rlb)，而且不需要任何 switch (交换机)的支持。接收负载均衡是通过 **ARP 协商**实现的。bonding 驱动截获本机发送的 ARP 应答，并把源硬件地址改写为 bond 中某个 slave 的唯一硬件地址，从而使得不同的对端使用不同的硬件地址进行通信。

# 模式选择

常用的有三种：

- `mode=0`：**平衡负载模式**，有自动备援，但需要”Switch”支援及设定。
- `mode=1`：**自动备援模式**，其中一条线若断线，其他线路将会自动备援。
- `mode=6`：**平衡负载模式**，有自动备援，不必”Switch”支援及设定。

需要说明的是：如果想做成 mode 0 的负载均衡，仅仅设置这里 `options bond0 miimon=100 mode=0` 是不够的，与网卡相连的交换机必须做特殊配置（这两个端口应该采取聚合方式），因为做 bonding 的这两块网卡是使用同一个 MAC 地址；从原理分析一下（bond 运行在 mode 0 下）：

> mode 0 下 bond 所绑定的网卡的 IP 都被修改成**相同的 mac 地址**，如果这些网卡都被接在同一个交换机，那么交换机的 arp 表里这个 mac 地址对应的端口就有多个，那么交换机接收到发往这个 mac 地址的包应该往哪个端口转发呢？正常情况下 mac 地址是全球唯一的，一个 mac 地址对应多个端口肯定使交换机迷惑了。所以 mode 0 下的 bond 如果连接到交换机，交换机这几个端口应该采取聚合方式（cisco 称为 ethernetchannel，foundry 称为 portgroup），因为交换机做了聚合后，聚合下的几个端口也被捆绑成一个 mac 地址；我们的解决办法是，两个网卡接入不同的交换机即可。
> 
> mode 6 模式下无需配置交换机，因为做 bonding 的这两块网卡是使用不同的 MAC 地址。



------


Linux bonding 模式中：

- 模式 0 （balance-rr）存在收发乱序问题；
- 模式 2（balance-xor）仅仅基于二层作为定义域的hash算法对带宽的利用不充分；
- 模式 5（balance-tlb）属于单向均衡；

> 具体 Linux Bonding 模式及选择信息可以参见：[Documentation/networking/bonding.txt](https://www.kernel.org/doc/Documentation/networking/bonding.txt) ；

Linux Bonding 模式一共有 7 种，基于带宽利用率考量一般会选择模式 4 （802.3ad）；802.3ad 模式是业界标准，通过创建一个聚合组，确保组内所有链路的速率和工作模式一致；

802.3ad 有三种 xmit_hash_policy 可供选择，默认缺省是 **layer2** ，还有 **layer2+3** 和 **layer3+4** 可选，参考 IBM 知识库、基于负载分配的最优性选择 layer3+4 方式。具体参见《[针对结合方式 Round-robin 策略（方式 0）、balance-xor（方式 2）和 802.3ad（方式 4）的交换机端口链路聚集和均衡算法](https://www.ibm.com/support/knowledgecenter/zh/ST5Q4U_1.5.0/com.ibm.storwize.v7000.unified.150.doc/mng_t_pub_netw_bondingmodes2_4.html)》；



------


## mode 4 模式下的 bond 示例

- 配置 1

```
# vi /etc/modprobe.d/bond.conf
alias bond0 bonding
options bond0 miimon=100 mode=4 xmit_hash_policy=layer3+4
```

- 配置 2

```
# vi /etc/sysconfig/network-scripts/ifcfg-eth0
DEVICE="eth0"
BOOTPROTO="none"
IPV6INIT="no"
ONBOOT="yes"
USERCTL=yes
MASTER=bond0
SLAVE=yes

# vi /etc/sysconfig/network-scripts/ifcfg-eth1
DEVICE="eth1"
BOOTPROTO="none"
IPV6INIT="no"
ONBOOT="yes"
USERCTL=yes
MASTER=bond0
SLAVE=yes

# vi /etc/sysconfig/network-scripts/ifcfg-bond0
DEVICE=bond0
NAME=bond0
BOOTPROTO=none
ONBOOT=yes
USERCTL=yes
TYPE=Ethernet
```


```
# vi /etc/sysconfig/network-scripts/ifcfg-eth2
DEVICE=eth2
NAME=eth2
BOOTPROTO=none
ONBOOT=yes
USERCTL=yes
TYPE=Ethernet
IPADDR=10.0.mmm.nn
NETMASK=255.255.255.0
GATEWAY=10.0.mmm.jj
```


- 配置 3

```
# ifconfig
bond0     Link encap:Ethernet  HWaddr E8:4D:D0:B9:95:13
          UP BROADCAST RUNNING MASTER MULTICAST  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:0
          RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)

eth0      Link encap:Ethernet  HWaddr E8:4D:D0:B9:95:13
          UP BROADCAST SLAVE MULTICAST  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)
          Memory:92d00000-92dfffff

eth1      Link encap:Ethernet  HWaddr E8:4D:D0:B9:95:13
          UP BROADCAST SLAVE MULTICAST  MTU:1500  Metric:1
          RX packets:0 errors:0 dropped:0 overruns:0 frame:0
          TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:0 (0.0 b)  TX bytes:0 (0.0 b)
          Memory:92c00000-92cfffff

eth2      Link encap:Ethernet  HWaddr 74:9D:8F:88:D5:30
          inet addr:10.0.mmm.nn  Bcast:10.0.mmm.255  Mask:255.255.255.0
          inet6 addr: fe80::----:----:----:d530/64 Scope:Link
          UP BROADCAST RUNNING MULTICAST  MTU:1500  Metric:1
          RX packets:38023859390 errors:0 dropped:585819 overruns:0 frame:0
          TX packets:163594627051 errors:0 dropped:0 overruns:0 carrier:0
          collisions:0 txqueuelen:1000
          RX bytes:196224133077121 (178.4 TiB)  TX bytes:218936209284531 (199.1 TiB)

# ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <NO-CARRIER,BROADCAST,MULTICAST,SLAVE,UP> mtu 1500 qdisc mq master bond0 state DOWN qlen 1000
    link/ether e8:4d:d0:b9:95:13 brd ff:ff:ff:ff:ff:ff
3: eth1: <NO-CARRIER,BROADCAST,MULTICAST,SLAVE,UP> mtu 1500 qdisc mq master bond0 state DOWN qlen 1000
    link/ether e8:4d:d0:b9:95:13 brd ff:ff:ff:ff:ff:ff
4: eth2: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP qlen 1000
    link/ether 74:9d:8f:88:d5:30 brd ff:ff:ff:ff:ff:ff
    inet 10.0.mmm.nn/24 brd 10.0.mmm.255 scope global eth2
    inet6 fe80::769d:8fff:fe88:d530/64 scope link
       valid_lft forever preferred_lft forever
5: eth3: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN qlen 1000
    link/ether 74:9d:8f:88:d5:31 brd ff:ff:ff:ff:ff:ff
6: bond0: <BROADCAST,MULTICAST,MASTER,UP> mtu 1500 qdisc noqueue state UNKNOWN
    link/ether e8:4d:d0:b9:95:13 brd ff:ff:ff:ff:ff:ff
```



```
# cat /proc/net/bonding/bond0
Ethernet Channel Bonding Driver: v3.6.0 (September 26, 2009)

Bonding Mode: IEEE 802.3ad Dynamic link aggregation
Transmit Hash Policy: layer3+4 (1)
MII Status: down
MII Polling Interval (ms): 100
Up Delay (ms): 0
Down Delay (ms): 0

802.3ad info
LACP rate: slow
Aggregator selection policy (ad_select): stable
bond bond0 has no active aggregator

Slave Interface: eth0
MII Status: down
Speed: Unknown
Duplex: Unknown
Link Failure Count: 0
Permanent HW addr: e8:4d:d0:b9:95:13
Aggregator ID: 1
Slave queue ID: 0

Slave Interface: eth1
MII Status: down
Speed: Unknown
Duplex: Unknown
Link Failure Count: 0
Permanent HW addr: e8:4d:d0:b9:95:14
Aggregator ID: 2
Slave queue ID: 0
```



