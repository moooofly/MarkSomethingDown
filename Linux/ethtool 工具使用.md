# ethtool 工具使用总结

## 环境信息


```
[root@xg-esm-data-2 ~]# uname -a
Linux xg-esm-data-2 3.10.0-229.el7.x86_64 #1 SMP Fri Mar 6 11:36:42 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
[root@xg-esm-data-2 ~]# lsb_release -a
LSB Version:	:core-4.1-amd64:core-4.1-noarch
Distributor ID:	CentOS
Description:	CentOS Linux release 7.1.1503 (Core)
Release:	7.1.1503
Codename:	Core
[root@xg-esm-data-2 ~]#
[root@xg-esm-data-2: ~]# ifconfig
bond0: flags=5187<UP,BROADCAST,RUNNING,MASTER,MULTICAST>  mtu 1500
        inet 10.0.38.106  netmask 255.255.255.0  broadcast 10.0.38.255
        ether 24:6e:96:0f:d5:74  txqueuelen 0  (Ethernet)
        RX packets 407863904654  bytes 590106724429195 (536.6 TiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 96967721535  bytes 13799933981768 (12.5 TiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

em1: flags=6211<UP,BROADCAST,RUNNING,SLAVE,MULTICAST>  mtu 1500
        ether 24:6e:96:0f:d5:74  txqueuelen 1000  (Ethernet)
        RX packets 209559928561  bytes 295230966970283 (268.5 TiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 48158384801  bytes 6455713935938 (5.8 TiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
        device memory 0x91d00000-91dfffff

em2: flags=6211<UP,BROADCAST,RUNNING,SLAVE,MULTICAST>  mtu 1500
        ether 24:6e:96:0f:d5:74  txqueuelen 1000  (Ethernet)
        RX packets 198303976093  bytes 294875757458912 (268.1 TiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 48809336734  bytes 7344220045830 (6.6 TiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
        device memory 0x91c00000-91cfffff

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        loop  txqueuelen 0  (Local Loopback)
        RX packets 12819176413  bytes 45476883532640 (41.3 TiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 12819176413  bytes 45476883532640 (41.3 TiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

[root@xg-esm-data-2: ~]#
```

## 确定网卡多队列具体数量

- 从**中断**维度确定多队列情况

```
[root@xg-esm-data-2: ~]# cat /proc/interrupts |grep "TxRx"
 129:  203654804         18  206954686          0  215703676          0  335972439          0  628033582          0  424121089          0  206385325          0 1247619717          0  691761717          0  211704785          0  602973131          0  498648391          0  393914682          0  387396157          0  298290226          0  301022906          0  IR-PCI-MSI-edge      em1-TxRx-0
 130:  197151430          0  548643182          0  474846981          0  429220613          0  563450823          0  246333538          0  721795690          0  219778629          0  162090939          0  249848486          0  419139013          0  494461714          0 1181659916          0  601369251          0  238803345          0  274469748          0  IR-PCI-MSI-edge      em1-TxRx-1
 131:  255435410          0  944492961         10  219919116          0  487410123          0  169770851          0  312747391          0  203237561          0  306189841          0  359706743          0  190954852          0  339493219          0 1410957417          0  546659822          0  232187453          0  261697302          0  709379609          0  IR-PCI-MSI-edge      em1-TxRx-2
 132:  216543277          0  428194809          0  302621636          0  333406502          0  340412909          0  186417009          0  436072572          0  452122282          0  422288495          0 1655179215          0  471480267          0  353840929          0  347674272          0  376689381          0  209448327          0  388945458          0  IR-PCI-MSI-edge      em1-TxRx-3
 133:  286036460          0  251770801          0  186366814          8  238694487          0  842617254          0  327821007          0  397138688          0  297162098          0 1489971604          0  620849273          0  231491097          0  269823628          0  273910247          0  444446945          0  275529986          0  470445591          0  IR-PCI-MSI-edge      em1-TxRx-4
 134:  393281469          0  160321906          0  786470716          0  183271989          0  222968304          0  254243179          0 1203583341          0  455897878          0  237356540          0  577079669          0  757415928          0  299506856          0  577446174          0  330092900          0  215152245          0  674291291          0  IR-PCI-MSI-edge      em1-TxRx-5
 135:  587760377          0  510031608          0  317578275          0 1041984283          8  584780979          0  292704619          0  180393272          0  355245120          0  610359422          0  290285846          0  463845511          0  206767097          0  221330219          0  371802005          0  400180148          0  458091034          0  IR-PCI-MSI-edge      em1-TxRx-6
 136:  927146256          0  278145564          0  239691655          0  432468291          0 1308845882          0  305586376          0  450682885          0  592920757          0  291274065          0  305124576          0  271763081          0  215111932          0  609097982          0  187131822          0  247813898          0  236610864          0  IR-PCI-MSI-edge      em1-TxRx-7
 139:  255361839          0  268919025          0  331268423          0  511030441          0  509239757         12  313857141          0  280124156          0  347529849          0  267811009          0  371678049          0   96020944          0  393941074          0  167872002          0  676087872          0  581322311          0  251915186          0  IR-PCI-MSI-edge      em2-TxRx-0
 140:  787597399          0  251808582          0  447270981          0  521379810          0  272971511          0  233064302          0  279951045          0  218605378          0  334399745          0  310380724          0  256087833          0  421201351          0  284148251          0  766998461          0  333088579          0  439742425          0  IR-PCI-MSI-edge      em2-TxRx-1
 141:  190244260          0  636607272          0  375882380          0  290459109          0  188079208          0  426086164          4  222813734          0  354461392          0  590663736          0  301255649          0  322645348          0  261392213          0  253247582          0  391208684          0  533083483          0  325357928          0  IR-PCI-MSI-edge      em2-TxRx-2
 142:  351890310          0  372176357          0  649916836          0  588650063          0  317963727          0  222376648          0  549636644          0  250197307          0  166430254          0  391837638          0  294609139          0  345560854          0  241110527          0  350806928          0  365608562          0  211251370          0  IR-PCI-MSI-edge      em2-TxRx-3
 143:  255417520          0  358010060          0  218218458          0  265349688          0  240727051          0 1241966025          0  657104242          5  377323908          0  148701231          0  216167093          0  190309912          0  225165334          0  301327577          0  237080149          0  477911396          0  353253152          0  IR-PCI-MSI-edge      em2-TxRx-4
 144:  364039001          0  172881510          0  258796565          0  155948748          0  195076074          0  551312022          0  402691870          0  285230393          0  239072732          0  426985290          0  556294845          0 1479456237          0  338164490          0  353041938          0  264202081          0  364634889          0  IR-PCI-MSI-edge      em2-TxRx-5
 145:  316597877          0  219492695          0  662755215          0  720726461          0  265329070          0  192182686          0  418046658          0  301316041         28  430398970          0  458684741          0  225842187          0  177044140          0  797823284          0   34369186          0  176496873          0  341573223          0  IR-PCI-MSI-edge      em2-TxRx-6
 146:  250823292          6  219724008          0  275970610          0  217514074          0  211403222          0  328630421          0  207289982          0  699268037          0  394690000          0  300832285          0  383496657          0  242844973          0  510332294          0  280400835          0  964921433          0  262955775          0  IR-PCI-MSI-edge      em2-TxRx-7
```

> 可以看到
> 
> - em1-TxRx-[0-7]
> - em2-TxRx-[0-7]


- 从**发送队列**（tx_queue）和**接收队列**（tx_queue）维度查看多队列情况 

```
[root@xg-esm-data-2: ~]# ethtool -S em1|grep tx_queue
     tx_queue_0_packets: 5976826882
     tx_queue_0_bytes: 770712461785
     tx_queue_0_restart: 0
     tx_queue_1_packets: 6038851951
     tx_queue_1_bytes: 852563009780
     tx_queue_1_restart: 0
     tx_queue_2_packets: 6089253123
     tx_queue_2_bytes: 829915801995
     tx_queue_2_restart: 0
     tx_queue_3_packets: 5991010002
     tx_queue_3_bytes: 820035881850
     tx_queue_3_restart: 0
     tx_queue_4_packets: 6097729769
     tx_queue_4_bytes: 851920657128
     tx_queue_4_restart: 0
     tx_queue_5_packets: 5959722442
     tx_queue_5_bytes: 763875366912
     tx_queue_5_restart: 0
     tx_queue_6_packets: 5977695265
     tx_queue_6_bytes: 788066761158
     tx_queue_6_restart: 0
     tx_queue_7_packets: 6028495226
     tx_queue_7_bytes: 779049377966
     tx_queue_7_restart: 0
[root@xg-esm-data-2: ~]# ethtool -S em1|grep rx_queue
     rx_queue_0_packets: 25845282908
     rx_queue_0_bytes: 36470719812033
     rx_queue_0_drops: 0
     rx_queue_0_csum_err: 9
     rx_queue_0_alloc_failed: 0
     rx_queue_1_packets: 26077784112
     rx_queue_1_bytes: 36661959202853
     rx_queue_1_drops: 0
     rx_queue_1_csum_err: 5
     rx_queue_1_alloc_failed: 0
     rx_queue_2_packets: 25930245876
     rx_queue_2_bytes: 36649676028754
     rx_queue_2_drops: 0
     rx_queue_2_csum_err: 6
     rx_queue_2_alloc_failed: 0
     rx_queue_3_packets: 25904770485
     rx_queue_3_bytes: 36515945604179
     rx_queue_3_drops: 0
     rx_queue_3_csum_err: 5
     rx_queue_3_alloc_failed: 0
     rx_queue_4_packets: 25910651699
     rx_queue_4_bytes: 36581817360386
     rx_queue_4_drops: 0
     rx_queue_4_csum_err: 7
     rx_queue_4_alloc_failed: 0
     rx_queue_5_packets: 28014134734
     rx_queue_5_bytes: 39018002580595
     rx_queue_5_drops: 0
     rx_queue_5_csum_err: 5
     rx_queue_5_alloc_failed: 0
     rx_queue_6_packets: 26043625399
     rx_queue_6_bytes: 36815697930826
     rx_queue_6_drops: 0
     rx_queue_6_csum_err: 4
     rx_queue_6_alloc_failed: 0
     rx_queue_7_packets: 25839085692
     rx_queue_7_bytes: 36525549115895
     rx_queue_7_drops: 0
     rx_queue_7_csum_err: 6
     rx_queue_7_alloc_failed: 0
[root@xg-esm-data-2: ~]#
```

- 从 **RSS (Receive Side Scaling)** 维度查看多队列情况 

> RSS 采用 IP-based 或 IP/Port-based (TCP) 哈希函数在指定数目的 RX queues 上进行负载分发，需要会结合一个 indirection 表：`queue = indirection_table[hash(packet)]` ；
> 可以通过如下命令查看 indirection table 的内容：

```
[root@xg-esm-data-2: ~]# ethtool -x em1
RX flow hash indirection table for em1 with 8 RX ring(s):
    0:      0     0     0     0     0     0     0     0
    8:      0     0     0     0     0     0     0     0
   16:      1     1     1     1     1     1     1     1
   24:      1     1     1     1     1     1     1     1
   32:      2     2     2     2     2     2     2     2
   40:      2     2     2     2     2     2     2     2
   48:      3     3     3     3     3     3     3     3
   56:      3     3     3     3     3     3     3     3
   64:      4     4     4     4     4     4     4     4
   72:      4     4     4     4     4     4     4     4
   80:      5     5     5     5     5     5     5     5
   88:      5     5     5     5     5     5     5     5
   96:      6     6     6     6     6     6     6     6
  104:      6     6     6     6     6     6     6     6
  112:      7     7     7     7     7     7     7     7
  120:      7     7     7     7     7     7     7     7
[root@xg-esm-data-2: ~]# ethtool -x em2
RX flow hash indirection table for em2 with 8 RX ring(s):
    0:      0     0     0     0     0     0     0     0
    8:      0     0     0     0     0     0     0     0
   16:      1     1     1     1     1     1     1     1
   24:      1     1     1     1     1     1     1     1
   32:      2     2     2     2     2     2     2     2
   40:      2     2     2     2     2     2     2     2
   48:      3     3     3     3     3     3     3     3
   56:      3     3     3     3     3     3     3     3
   64:      4     4     4     4     4     4     4     4
   72:      4     4     4     4     4     4     4     4
   80:      5     5     5     5     5     5     5     5
   88:      5     5     5     5     5     5     5     5
   96:      6     6     6     6     6     6     6     6
  104:      6     6     6     6     6     6     6     6
  112:      7     7     7     7     7     7     7     7
  120:      7     7     7     7     7     7     7     7
[root@xg-esm-data-2: ~]#
```


## 运行时变更网卡队列数量

```
[root@xg-esm-data-2: ~]# ethtool --set-channels em1 combined 1
[root@xg-esm-data-2: ~]# ethtool -S em1|grep rx_queue
     rx_queue_0_packets: 25845412988
     rx_queue_0_bytes: 36470913708708
     rx_queue_0_drops: 0
     rx_queue_0_csum_err: 9
     rx_queue_0_alloc_failed: 0
[root@xg-esm-data-2: ~]#
[root@xg-esm-data-2: ~]# ethtool -S em1|grep tx_queue
     tx_queue_0_packets: 5976910043
     tx_queue_0_bytes: 770729638530
     tx_queue_0_restart: 0
[root@xg-esm-data-2: ~]#
```

## 运行时 Ring Buffer 大小调整

```
[root@xg-esm-data-2: ~]# ethtool -g em1
Ring parameters for em1:
Pre-set maximums:
RX:		4096
RX Mini:	0
RX Jumbo:	0
TX:		4096
Current hardware settings:
RX:		4096
RX Mini:	0
RX Jumbo:	0
TX:		4096

[root@xg-esm-data-2: ~]#
[root@xg-esm-data-2: ~]# ethtool -G em1 rx 1024 tx 2048
[root@xg-esm-data-2: ~]# ethtool -g em1
Ring parameters for em1:
Pre-set maximums:
RX:		4096
RX Mini:	0
RX Jumbo:	0
TX:		4096
Current hardware settings:
RX:		1024
RX Mini:	0
RX Jumbo:	0
TX:		2048

[root@xg-esm-data-2: ~]#
[root@xg-esm-data-2: ~]# ethtool -G em1 rx 4096 tx 4096
[root@xg-esm-data-2: ~]#
[root@xg-esm-data-2: ~]# ethtool -g em1
Ring parameters for em1:
Pre-set maximums:
RX:		4096
RX Mini:	0
RX Jumbo:	0
TX:		4096
Current hardware settings:
RX:		4096
RX Mini:	0
RX Jumbo:	0
TX:		4096

[root@xg-esm-data-2: ~]#
```

## 接收方向的 fifo 丢包数（Ring Buffer 丢包）

```
[root@xg-esm-data-2: ~]# ethtool -S em1|grep rx_fifo
     rx_fifo_errors: 0
```

- 查看网卡当前的 offloading 配置情况

```
[root@xg-esm-data-2: ~]# ethtool -k em1|grep offload
tcp-segmentation-offload: on
udp-fragmentation-offload: off [fixed]
generic-segmentation-offload: on
generic-receive-offload: on
large-receive-offload: off [fixed]
rx-vlan-offload: on
tx-vlan-offload: on
[root@xg-esm-data-2: ~]#
```


----------


## 案例

一台机器经常收到丢包的报警，检查步骤：

- **Speed**+**Duplex**+**CRC** 检查

> 排除物理层面的干扰

```
[root@xg-esm-data-2 ~]# ethtool em2 | egrep 'Speed|Duplex'
	Speed: 1000Mb/s
	Duplex: Full
[root@xg-esm-data-2 ~]# ethtool -S em2 | grep crc
     rx_crc_errors: 0
```

- 观察 ifconfig 中的 **errors**+**dropped**+**overruns**

> 确定网卡整体错误信息

```
[root@xg-esm-data-2 ~]# while true; do ifconfig em2 | grep RX | grep overruns; sleep 1; done
        RX errors 0  dropped 0  overruns 0  frame 0
        RX errors 0  dropped 0  overruns 0  frame 0
        ...
```

> 确定每一个 rx_queue 上的 drop 情况

```
[root@xg-esm-data-2 ~]# ethtool -S em2 | grep drop
     dropped_smbus: 0
     tx_dropped: 0
     rx_queue_0_drops: 0
     rx_queue_1_drops: 0
     rx_queue_2_drops: 0
     rx_queue_3_drops: 0
     rx_queue_4_drops: 0
     rx_queue_5_drops: 0
     rx_queue_6_drops: 0
     rx_queue_7_drops: 0
```

两种方式观察网卡信息

```
[root@xg-esm-data-2 ~]# ifconfig em2; netstat -i | column -t
em2: flags=6211<UP,BROADCAST,RUNNING,SLAVE,MULTICAST>  mtu 1500
        ether 24:6e:96:0f:d5:74  txqueuelen 1000  (Ethernet)
        RX packets 198464716743  bytes 295111769404191 (268.4 TiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 48856869916  bytes 7365554557653 (6.6 TiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0
        device memory 0x91c00000-91cfffff

Kernel  Interface  table
Iface   MTU        RX-OK         RX-ERR  RX-DRP  RX-OVR  TX-OK        TX-ERR  TX-DRP  TX-OVR  Flg
bond0   1500       408179592404  0       0       0       97062107603  0       0       0       BMmRU
em1     1500       209714875661  0       0       0       48205237687  0       0       0       BMsRU
em2     1500       198464716743  0       0       0       48856869916  0       0       0       BMsRU
lo      65536      12819980798   0       0       0       12819980798  0       0       0       LRU
[root@xg-esm-data-2 ~]#
```

针对 RX 来说

- RX `errors`: 表示总的收包相关错误数量，包括
    - too-long-frames errors
    - ring-buffer overflow errors
    - crc errors
    - frame alignment errors
    - fifo overruns
    - missed packets
- RX `dropped`: 表示数据包**已经进入了 Ring Buffer**，但是由于内存不够等系统原因，导致**在拷贝到系统内存的过程中被丢弃**。
- RX `overruns`: 表示**发生了 fifo 的 overrun**，这是由于 Ring Buffer (aka Driver Queue) 传输的 IO 大于 kernel 能够处理的 IO 导致的，而 Ring Buffer 则是指在发起 IRQ 请求之前的那块 buffer (The ring-buffer refers to a buffer that the NIC transfers frames to before raising an IRQ with the kernel)。很明显，**overruns 的增大意味着数据包没到 Ring Buffer 就被网卡物理层给丢弃了**，而 CPU 无法即使的处理中断是造成 Ring Buffer 满的原因之一（有问题的机器可能室因为 interruprs 分布的不均匀（都压在 core0），没有做 affinity 而造成的丢包；
- RX `frame`: 表示 misaligned 的 frames ；

针对 TX 来说，出现计数值增大的原因主要包括

- errors due to the transmission being aborted
- errors due to the carrier
- fifo errors
- heartbeat errors
- window errors

而 collisions 则表示由于 CSMA/CD 造成的传输中断；


不管是使用何种工具，最终的数据无外乎是从下面这两个地方获取到的:

- `/sys/class/net/<if>/statistics/*`
- `/proc/net/dev`

小图一张：

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/%E8%A7%82%E5%AF%9F%E7%BD%91%E5%8D%A1%20RingBuffer%20%E4%B8%A2%E5%8C%85.png "")


----------

参考：

- [The Missing Man Page for ifconfig](http://blog.hyfather.com/blog/2013/03/04/ifconfig/)
- [ifconfig 下面的一些字段(errors, dropped, overruns)](http://jaseywang.me/2014/08/16/ifconfig-%E4%B8%8B%E9%9D%A2%E7%9A%84%E4%B8%80%E4%BA%9B%E5%AD%97%E6%AE%B5errors-dropped-overruns/)


