# 吞吐量 Throughput

> When used in the context of communication networks, such as Ethernet or packet radio, throughput or network throughput is the rate of successful message delivery over a communication channel. The data these messages belong to may be delivered over a physical or logical link, or it can pass through a certain network node. Throughput is usually measured in **bits per second** (`bit/s` or `bps`), and sometimes in **data packets per second** (`p/s` or `pps`) or **data packets per time slot**.

吞吐量的度量单位有三种：

- bps
- pps
- ppt

> The throughput of a communication system may be affected by various factors, including the **limitations of underlying analog physical medium**, **available processing power of the system components**, and **end-user behavior**. When various protocol overheads are taken into account, useful rate of the transferred data can be significantly lower than the maximum achievable throughput; the useful part is usually referred to as `goodput`.

影响吞吐量的因素：

- **底层物理介质的限制**；
- **系统组件的可用处理能力**；
- **端用户本身的行为**；

考虑到上述因素可能导致的额外开销，实际吞吐量可能远低于理论上可达的最大吞吐量，前者通常被称为“goodput”（有效吞吐量）；


----------

> Sometimes, when we talk about **device performance** we are talking in terms of **packets per second** (`pps`) and **bits per second** (`bps`). But in latter case it's not quite correct to say "this device can do one hundred megabits ber second" **because router/switch/whatever performance is greatly depends on packet size** and if you want to mention device performance in a more accurate and professional way you would say "this device can do one hundred megabits per second at 64 bytes packet size"

设备性能和包大小有关系；因此正确描述设备性能的表达方式为：在 xxx 字节包大小情况下，可以达到 xxx bps ；

> Often vendors such as our favorite **Cisco** specify device performance as **packets per second**, so we don't need to bother about packet size mentioning because pps is rather a characteristic of device's (processor, bus, ASICs) computing power. Packets per second more or less still the same with different packets size. But it is not very convinient to deal with pps in a real life because we have to know "real" device performance in our network. So we have to do two things:
> 
> 1) Determine the average packet size which is specific for our network. For example traffic profile for our network could be 30% ftp-data (large packet at 1500 bytes) and 70% VoIP-data (a lot of small packets at 64 bytes) so our average packet size is about 800 bytes.
> 
> 2) Calculate with simple formula how much there will be Megabits per second (Mbps) if our average packet size is 800 bytes and device performance is, lets say, 100 kpps (one hundred thousand of packets per second)
>
> The second step is not a big deal for a real professional, but we live in 21st century, aren't we? Unfortunately I didn't found any bps to pps converter/calculator anywhere online so I decided to make it myself (though I'm not a programmer).  

知名设备厂商通常会基于 pps 数值展示设备性能，而这是从计算能力角度出发的；而有些场景中，我们更想要确定真正的设备性能，此时需要

- **确定目标网络上的平均包大小**；例如，经过网络分析后可知，存在 30% 的 ftp 数据（即对应 1500 字节的大数据包）和 70% 的 VoIP 数据（即对应大量 64 字节的小数据包），因此，平均包大小在 800 字节左右；
- **基于简单的公式换算**，就能得出当前吞吐量为多少 Mbps ，假如当前平均包大小为 800 字节，且设备性能假定为 100kpps 的话；则有 `100 * 1024 * 800 * 8 / 1024 = 640000 kbps = 62.5 Mbps`

需要注意的是：上述计算不能当作公式使用，因为设备性能（pps）是和包大小相关的，且这种相关性并非线性关系；


> P.S. There is one more thing I need to say. There are at least three well know packet size: the least one - **64 bytes** (toughest case for device, usually referred with **router/switch** performance), the biggest one **1500 bytes** (sometimes 1400 bytes) usually referred with **firewall/VPN** performance and the so-called "real" one - **IMIX** at **427 bytes**, which represents an average packet size somewhere in the Internet (but I saw values in between 300-900 bytes)

三种常用作基准的包大小：

- **64** 字节：常用于讨论 router/switch 的设备性能；
- **1500** 字节（有时也用 1400 字节）：常用于讨论 firewall/VPN 的设备性能；
- **427** 字节：常用于讨论 [IMIX](https://en.wikipedia.org/wiki/Internet_Mix) ，代表了 typical Internet traffic ；

> 在线计算工具：[这里](http://www.ccievault.net/index.php/tools)


----------


参考：

- [wiki/Throughput](https://en.wikipedia.org/wiki/Throughput)
- [Bits per second to packets per second converter](http://www.ccievault.net/index.php/articles/37-cvnarticles/58-bps2pps)
- [How many Packets per Second per port are needed to achieve Wire-Speed?](https://kb.juniper.net/InfoCenter/index?page=content&id=KB14737)


