# PF_RING Installation

PF_RING 的安装既可以基于从 [GIT](https://github.com/ntop/PF_RING/) 上下载的源码，也可以使用 [Ubuntu/CentOS 仓库](http://packages.ntop.org)中的包进行安装，详见 **README.apt_rpm_packages** 中的解释说明；

当你下载了 PF_RING 后，实际获取到如下组件：

* 用户空间 PF_RING SDK ；
* 一个加强版的 `libpcap` 库，能够透明利用 PF_RING 功能（如果安装了的话），或回退到标准版本行为（如果未安装）；
* 内核模块 PF_RING ；
* PF_RING ZC drivers ；

## Linux Kernel Module Installation

为了编译 PF_RING 内核模块，你需要安装 linux kernel 头文件（或 kernel 源码文件）；

```shell
cd <PF_RING PATH>/kernel
make
make install
``` 

需要注意的是，kernel 模块的安装（通过 `make install` 命令）需要 root 权限；

## Running PF_RING

在使用任何基于 PF_RING 的应用前，内核模块 pf_ring 应该先被加载（以超级用户身份）：

```shell
insmod <PF_RING PATH>/kernel/pf_ring.ko [min_num_slots=x][enable_tx_capture=1|0] [ enable_ip_defrag=1|0] [quick_mode=1|0]
``` 

其中：

* `min_num_slots` ：用于指定 ring slots 的最小值（默认为 4096）；
* `enable_tx_capture` ：若设置为 1 则捕获 outgoing packets ，若设置为 0 则不捕获 outgoing packets（默认为 RX+TX）；
* `enable_ip_defrag` ：若设置为 1 ，则使能 IP 重组（defragmentation）功能，只有 rx traffic 才能进行重组；
* `quick_mode` ：若设置为 1 则以全速运行，但至多每个接口使用一个 socket ；

使用示例：

```shell
cd <PF_RING PATH>/kernel
insmod pf_ring.ko min_num_slot=8192 enable_tx_capture=0 quick_mode=1
``` 

如果是想要达到 10 Gigabit 或之上的线速 packet 捕获速度，你应该使用 ZC drivers ；ZC drivers 属于 PF_RING 发布的一部分，可以在 `drivers/` 中找到；详情参考 **README.ZC** 的说明；

## Libpfring and Libpcap Installation

`libpfring`（用户空间 PF_RING 库）和 `libpcap` 均以源码格式发布；可以按照如下方式进行编译：

```shell 
cd <PF_RING PATH>/userland/lib
./configure
make
sudo make install
cd ../libpcap
./configure
make
``` 

需要注意的是：

* lib 是可重入的（reentrant），因此非常有必要令你的 PF_RING-enabled 应用程序还链接到 `-lpthread` 库上；
* 已停止更新的、基于静态链接的、pcap-based 应用程序，需要基于新生成的 PF_RING-enabled `libpcap.a` 库重新编译，以便利用 `PF_RING` 特性；不要期望使用 PF_RING 的时候，可以不重新编译你已存在的应用程序；

## Application Examples

如果是你 PF_RING 新手，你可以从一些例子开始；在 `userland/examples` 文件夹下有很多可使用的 PF_RING 应用程序：
	
```shell
cd <PF_RING PATH>/userland/examples 
ls *.c
alldevs.c      pfcount_82599.c	         pflatency.c  pfwrite.c
pcap2nspcap.c  pfcount.c	         pfsend.c     preflect.c
pcount.c       pfcount_multichannel.c    pfsystest.c
pfbridge.c     pfdump.c		         pfutils.c
make
``` 

例如，`pfcount` 允许在接收 packets 时打印一些统计信息： 

```shell
./pfcount -i zc:eth1
...
=========================
Absolute Stats: [64415543 pkts rcvd][0 pkts dropped]
Total Pkts=64415543/Dropped=0.0 %
64'415'543 pkts - 5'410'905'612 bytes [4'293'748.94 pkt/sec - 2'885.39 Mbit/sec]
=========================
Actual Stats: 14214472 pkts [1'000.03 ms][14'214'017.15 pps/9.55 Gbps]
=========================
``` 

另外一个例子是 `pfsend` ，允许你以指定的速度发送 packets（人工合成 packets 或使用 .pcap 文件) ：

```shell
./pfsend -f 64byte_packets.pcap -n 0 -i zc:eth1 -r 5
...
TX rate: [current 7'508'239.00 pps/5.05 Gbps][average 7'508'239.00 pps/5.05 Gbps][total 7'508'239.00 pkts]
``` 

## PF_RING Additional Modules

PF_RING 库采用模块化架构（modular architecture），允许使用额外的组件，而不仅仅是标准的 PF_RING 内核模块；在基于 `configure` 脚本检测到支持时，相应的组件会被编译到库中；

PF_RING 模块当前包含了对 Accolade, Endace DAG, Exablaze, Myricom, Napatech 等等的支持；
