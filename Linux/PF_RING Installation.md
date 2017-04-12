# PF_RING 安装

> 原文地址：[这里](https://github.com/ntop/PF_RING/blob/dev/doc/README.install.md)

PF_RING 的安装既可以基于从 [GIT](https://github.com/ntop/PF_RING/) 上下载的源码，也可以使用 [Ubuntu/CentOS 仓库](http://packages.ntop.org)中的包进行安装，详见 **README.apt_rpm_packages** 中的解释说明；

当你下载了 PF_RING 后，实际获取到如下组件：

* PF_RING 用户空间 SDK（`libpfring.so` 和 `libpfring.a`）；
* 一个加强版的 `libpcap` 库，能够透明利用 PF_RING 功能（如果安装了的话），或回退到标准版本行为（如果未安装）；
* PF_RING 内核模块（`pf_ring.ko`）；
* PF_RING ZC drivers ；

## Linux 内核模块安装（PF_RING）

> 以下内容基于 `PF_RING/kernel/Makefile` 文件内容；

`PF_RING/kernel/Makefile` 文件内容如下：

```shell
#
# (C) 2009-15 - ntop.org
#

obj-m := pf_ring.o

# assigned by Makefile.dkms
GIT_REV:=

ifndef GIT_REV
 ifneq (, $(shell which git))
  ifeq (, $(shell echo ${SUBDIRS}))
   GIT_BRANCH=$(shell git branch --no-color|cut -d ' ' -f 2)
   GIT_HASH=$(shell git rev-parse HEAD)
   ifneq ($(strip $(GIT_BRANCH)),)
    GIT_REV:=${GIT_BRANCH}:${GIT_HASH}
   endif
  endif
 endif
endif

ifneq ($(strip $(GIT_REV)),)
 GITDEF:=-DGIT_REV="\"${GIT_REV}\""
endif

ifeq (,$(BUILD_KERNEL))
 BUILD_KERNEL=$(shell uname -r)
endif

PWD := $(shell pwd)
EXTRA_CFLAGS += -I${PWD} ${GITDEF}

HERE=${PWD}

# set the install path
INSTDIR := $(DESTDIR)/lib/modules/$(BUILD_KERNEL)/kernel/net/pf_ring
TARGETDIR := $(DESTDIR)/usr/src/$(BUILD_KERNEL)/include/linux/

all: Makefile pf_ring.c linux/pf_ring.h
#   @if test "$(USER)" = "root"; then \
#       echo "********** WARNING WARNING WARNING **********"; \
#       echo "*"; \
#       echo "* Compiling PF_RING as root might lead you to compile errors"; \
#       echo "* Please compile PF_RING as unpriviliged user"; \
#       echo "*"; \
#       echo "*********************************************"; \
#   fi
    make -C /lib/modules/$(BUILD_KERNEL)/build SUBDIRS=${HERE} EXTRA_CFLAGS='${EXTRA_CFLAGS}' modules

dkms-deb:
    sudo make -f Makefile.dkms deb

dkms-rpm:
    sudo make -f Makefile.dkms rpm

clean:
    make -C /lib/modules/$(BUILD_KERNEL)/build SUBDIRS=$(HERE) clean
    \rm -f *~ Module.symvers  Module.markers  modules.order *#

install:
    mkdir -p $(INSTDIR)
    cp *.ko $(INSTDIR)
    mkdir -p $(DESTDIR)/usr/include/linux
    cp linux/pf_ring.h $(DESTDIR)/usr/include/linux
    @if test -d ${TARGETDIR}; then \
        cp linux/pf_ring.h ${TARGETDIR}; \
    fi
ifeq (,$(DESTDIR))
    /sbin/depmod $(BUILD_KERNEL)
else
    @echo "*****NOTE:";
    @echo "pf_ring,ko kernel module installed in ${DESTDIR}";
    @echo "/sbin/depmod not run.  modprobe pf_ring won't work " ;
    @echo "You can load the kernel module directly using" ;
    @echo "insmod <path>/pf_ring.ko" ;
    @echo "*****";
endif
```

> 为了编译 PF_RING 内核模块，你需要安装 linux kernel 头文件（或 kernel 源码文件）；

编译安装 PF_RING 内核模块 `pf_ring.ko` ；

```shell
cd <PF_RING PATH>/kernel
make
make install
```

执行输出如下：

```shell
[root@xg-esm-data-4 kernel]# make install
mkdir -p /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring
cp *.ko /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring
mkdir -p /usr/include/linux
cp linux/pf_ring.h /usr/include/linux
/sbin/depmod 3.10.0-229.el7.x86_64
[root@xg-esm-data-4 kernel]#
```

需要注意的是，kernel 模块的安装（通过 `make install` 命令）需要 root 权限；

## 加载 PF_RING 内核模块

在使用任何基于 PF_RING 的应用前，内核模块 `pf_ring.ko` 需要先被加载（以超级用户身份）：

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

如果是想要达到 10 Gigabit 或之上的线速 packet 捕获速度，你应该使用 ZC drivers ；ZC drivers 属于 PF_RING 发布的一部分，可以在 `drivers/` 中找到；详情参考 **[README.ZC](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/PF_RING%20ZC.md)** 的说明；

## 安装 Libpfring 和 Libpcap

`libpfring`（PF_RING 对应的用户空间库）和 `libpcap` （加强版）均以源码格式发布，可以按照如下方式进行编译：

```shell 
## 编译安装 libpfring
cd <PF_RING PATH>/userland/lib
./configure
make
sudo make install
## 编译 libpcap
cd ../libpcap
./configure
make
```

> 默认安装到 `/usr/local/lib` 目录下；

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

## PF_RING 支持的其它模块

PF_RING 库采用模块化架构（modular architecture），允许使用额外的组件，而不仅仅是标准的 PF_RING 内核模块；在基于 `configure` 脚本检测到支持时，相应的组件会被编译到库中；

PF_RING 模块当前包含了对 Accolade, Endace DAG, Exablaze, Myricom, Napatech 等等的支持；


------

## insmod 说明

insmod - Simple program to insert a module into the Linux Kernel

insmod 是用于将指定模块插入 kernel 的简便程序；大多数用户更愿意使用 `modprobe(8)` ，因为后者更加智能，并能处理模块间的依赖；

该程序仅会报告一些最通用的 error 消息：因为当前尝试进行模块链接的工作已经由 kernel 负责，因此通常请款下 dmesg 能够提供关于 error 的更多信息；

## modprobe 说明

modprobe - Add and remove modules from the Linux Kernel

modprobe 能够智能的从 Linux 内核中添加和删除模块：需要注意的是，为了方便，在模块名字中出现 _ 和 - 时将不做区别处理（会进行自动 underscore 转换）；

modprobe 会在模块目录 `/lib/modules/`uname -r`` 中查找所有模块和其它相关文件；但仍会在 `/etc/modprobe.d` 目录中查找可选配置文件（详见 `modprobe.d(5)`）；  

modprobe 同样会使用在内核命令行中指定的模块选项，形式为 `<module>.<option>` 和黑名单 `modprobe.blacklist=<module>` ；

需要注意的是，与 2.4 系列的 Linux 内核不同（该工具不支持这个版本系列），该版本的 modprobe 不会对模块本身进行任何操作：因为 symbols 解析和 parameters 解析均在 kernel 内部完成；因此，有些时候模块失败信息会由内核消息提供，详见 `dmesg(8)` ；

modprobe 永远期望访问一个最新的（up-to-date）`modules.dep.bin` 文件（或者一个作为 fallback 的人类可读的 modules.dep 文件）；这些文件可由配套 modprobe 的 depmod 实用程序生成（详见 `depmod(8)`）；文件中会列出了每一个模块所需要的其他模块信息（如果存在这种依赖的话），而 modprobe 基于其内容自动添加和删除依赖关系；

如果在 modulename 后指定了参数，则会被传入内核（除了列在配置文件中的那些选项外）；





