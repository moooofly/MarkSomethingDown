
当linux系统上有多个单独网卡，若想充分利用这些网卡，同时能够对外提供统一的网络地址，增大网络的吞吐量，同时提高网络的可用性，这时就需要bond来帮助我们解决这个问题。

bond的配置实在很简单，但是配置不好，很容易造成严重的网络问题。bonding功能是linux内核就自带了，因此，通常不需要安装它，只需要把bonding模块加载到内核里即可。

# 七种bond模式说明

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

特点：数据包传输顺序是依次传输（即第1个包走eth0，下一个包就走eth1，一直循环下去，直到最后一个传输完毕），此模式提供**负载平衡**和**容错能力**；但是我们知道如果一个连接或者会话的数据包从不同的接口发出的话，中途再经过不同的链路，在客户端很有可能会出现**数据包无序到达的问题**，而无序到达的数据包需要重新要求被发送，这样网络的吞吐量就会下降；

## 第二种模式：mod=1，即： (active-backup) Active-backup policy（主-备份策略）

特点：**只有一个设备处于活动状态**，当一个宕掉另一个马上由备份转换为主设备。mac地址是外部可见得，**从外面看来，bond的MAC地址是唯一的**，以避免switch(交换机)发生混乱。此模式只提供了**容错能力**；由此可见此算法的优点是可以提供高网络连接的可用性，但是它的资源利用率较低，只有一个接口处于工作状态，在有 N 个网络接口的情况下，资源利用率为1/N ；

## 第三种模式：mod=2，即：(balance-xor) XOR policy（平衡策略）

特点：基于指定的传输HASH策略传输数据包。缺省的策略是：(源MAC地址 XOR 目标MAC地址) % slave数量。其他的传输策略可以通过xmit_hash_policy选项指定，此模式提供**负载平衡**和**容错能力**；

## 第四种模式：mod=3，即：broadcast（广播策略）

特点：在每个slave接口上传输每个数据包，此模式提供了**容错能力**

## 第五种模式：mod=4，即：(802.3ad) IEEE 802.3adDynamic link aggregation（IEEE 802.3ad 动态链接聚合）

特点：创建一个聚合组，它们共享同样的速率和双工设定。根据802.3ad规范将多个slave工作在同一个激活的聚合体下。

外出流量的slave选举是基于传输hash策略，该策略可以通过xmit_hash_policy选项从缺省的XOR策略改变到其他策略。需要注意的是，并不是所有的传输策略都是802.3ad适应的，尤其考虑到在802.3ad标准43.2.4章节提及的包乱序问题。不同的实现可能会有不同的适应性。

必要条件：

- 条件1：ethtool支持获取每个slave的速率和双工设定；
- 条件2：switch(交换机)支持IEEE 802.3ad Dynamic link aggregation；
- 条件3：大多数switch(交换机)需要经过特定配置才能支持802.3ad模式；

## 第六种模式：mod=5，即：(balance-tlb) Adaptive transmit load balancing（适配器传输负载均衡）

特点：不需要任何特别的switch(交换机)支持的**通道bonding**。在每个slave上根据当前的负载（根据速度计算）分配外出流量。如果正在接受数据的slave出故障了，另一个slave接管失败的slave的MAC地址。

该模式的必要条件：ethtool支持获取每个slave的速率；

## 第七种模式：mod=6，即：(balance-alb) Adaptive load balancing（适配器适应性负载均衡）

特点：该模式包含了balance-tlb模式，同时加上**针对IPV4流量的接收负载均衡**(receive load balance, rlb)，而且不需要任何switch(交换机)的支持。接收负载均衡是通过**ARP协商**实现的。bonding驱动截获本机发送的ARP应答，并把源硬件地址改写为bond中某个slave的唯一硬件地址，从而使得不同的对端使用不同的硬件地址进行通信。

# 模式选择

常用的有三种：

- `mode=0`：**平衡负载模式**，有自动备援，但需要”Switch”支援及设定。
- `mode=1`：**自动备援模式**，其中一条线若断线，其他线路将会自动备援。
- `mode=6`：**平衡负载模式**，有自动备援，不必”Switch”支援及设定。

需要说明的是：如果想做成mode 0的负载均衡，仅仅设置这里`options bond0 miimon=100 mode=0`是不够的，与网卡相连的交换机必须做特殊配置（这两个端口应该采取聚合方式），因为做bonding的这两块网卡是使用同一个MAC地址.从原理分析一下（bond运行在mode 0下）：

> mode 0下bond所绑定的网卡的IP都被修改成**相同的mac地址**，如果这些网卡都被接在同一个交换机，那么交换机的arp表里这个mac地址对应的端口就有多个，那么交换机接收到发往这个mac地址的包应该往哪个端口转发呢？正常情况下mac地址是全球唯一的，一个mac地址对应多个端口肯定使交换机迷惑了。所以mode0下的bond如果连接到交换机，交换机这几个端口应该采取聚合方式（cisco称为 ethernetchannel，foundry称为portgroup），因为交换机做了聚合后，聚合下的几个端口也被捆绑成一个mac地址.我们的解决办法是，两个网卡接入不同的交换机即可。
> 
> mode 6模式下无需配置交换机，因为做bonding的这两块网卡是使用不同的MAC地址。








