# RSS (Receive Side Scaling)

几乎所有 Intel（以及其它厂商）的 NICs 都支持 RSS ，这意味着**可以基于硬件针对 packets 进行哈希运算以便在多条 RX queues 中进行负载分布（distribute）**；

为了配置 queues 数量，你可以在调用 `insmod` 时使用 `RSS` 参数（如果你是基于 packages 安装的 PF_RING ZC drivers ，则可以使用配置文件进行配置，详情参见 `README.apt_rpm_packages` 说明），并传递一个逗号分隔的数字（每个接口上的 queue 数目）列表（每个接口对应一个列表），例如：

针对每个接口使能与处理器数量相同的 queue ：

``` 
insmod ixgbe.ko RSS=0,0,0,0
``` 

每个接口使能 4 个 RX queues ：

``` 
insmod ixgbe.ko RSS=4,4,4,4
``` 

去使能 multiqueue 功能（每个接口 1 个 RX queue）：

``` 
insmod ixgbe.ko RSS=1,1,1,1
``` 

另外，还可以运行时基于 `ethtool` 配置 RX queues 数量：

``` 
ethtool --set-channels <if> combined 1
``` 

RSS 基于 IP-based 或 IP/Port-based (TCP) 哈希函数在指定数目的 RX queues 上进行负载分发，同时会结合一个 indirection 表：`queue = indirection_table[hash(packet)]` ；你可以通过如下命令查看 indirection table 的内容：

``` 
ethtool -x <if>
``` 

可以使用 `ethtool` 进行 indirection table 配置：即针对每一个 RX queue 简单设置一个 weights 值；例如，如果我们希望全部 traffic 都走向 queue 0 中，则我们可以通过如下命令配置具有 4 个 RX queues 的卡：

``` 
ethtool -X <if> weight 1 0 0 0
``` 

## Naming convention

为了打开指定接口上的 queue ，你必须通过 "@<ID>" 后缀指定 queue ID ，例如：

``` 
pfcount -i zc:eth1@0
``` 

需要注意的是，如果你对一个具有多个 RSS queues 的接口进行配置，并且使用 `zc:eth1` 形式按照使用 ZC 模式打开，则等同于 `zc:eth1@0` 形式打开；这种方式不适用于标准内核模式，因为其对接口进行了抽象并从 eth1 上进行捕获，而这意味着从所有 queues 进行捕获；发生这种情况对原因在于 ZC 是一种 kernel-bypass 技术，因此不存在抽象层，故应用将直接打开接口上的 queue ，即对应了当 RSS=1 时的 full interface ；