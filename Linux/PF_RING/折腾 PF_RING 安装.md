# 折腾 PF_RING 安装

申请了一台线上物理机用来折腾 PF_RING ；

## 机器初始状态

```
[root@wg-esm-hc-1 ~]# uname -a
Linux wg-esm-hc-1 3.10.0-229.el7.x86_64 #1 SMP Fri Mar 6 11:36:42 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# cat /etc/*release
CentOS Linux release 7.1.1503 (Core)
NAME="CentOS Linux"
VERSION="7 (Core)"
ID="centos"
ID_LIKE="rhel fedora"
VERSION_ID="7"
PRETTY_NAME="CentOS Linux 7 (Core)"
ANSI_COLOR="0;31"
CPE_NAME="cpe:/o:centos:centos:7"
HOME_URL="https://www.centos.org/"
BUG_REPORT_URL="https://bugs.centos.org/"

CENTOS_MANTISBT_PROJECT="CentOS-7"
CENTOS_MANTISBT_PROJECT_VERSION="7"
REDHAT_SUPPORT_PRODUCT="centos"
REDHAT_SUPPORT_PRODUCT_VERSION="7"

CentOS Linux release 7.1.1503 (Core)
CentOS Linux release 7.1.1503 (Core)
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# locate libpcap.so
/usr/lib64/libpcap.so.1
/usr/lib64/libpcap.so.1.5.3
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# lsmod |grep pf_ring
[root@wg-esm-hc-1 ~]# lsmod |grep ixgbe
ixgbe                 290931  0
mdio                   13807  1 ixgbe
ptp                    18933  1 ixgbe
dca                    15130  2 ixgbe,ioatdma
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# lspci -v | less
...
04:00.0 Ethernet controller: Intel Corporation 82599ES 10-Gigabit SFI/SFP+ Network Connection (rev 01)
	Subsystem: Hewlett-Packard Company Ethernet 10Gb 2-port 560FLR-SFP+ Adapter
	Flags: bus master, fast devsel, latency 0, IRQ 17
	Memory at 92c00000 (32-bit, non-prefetchable) [size=1M]
	I/O ports at 2020 [size=32]
	Memory at 92d04000 (32-bit, non-prefetchable) [size=16K]
	[virtual] Expansion ROM at 92d80000 [disabled] [size=512K]
	Capabilities: [40] Power Management version 3
	Capabilities: [50] MSI: Enable- Count=1/1 Maskable+ 64bit+
	Capabilities: [70] MSI-X: Enable+ Count=64 Masked-
	Capabilities: [a0] Express Endpoint, MSI 00
	Capabilities: [e0] Vital Product Data
	Capabilities: [100] Advanced Error Reporting
	Capabilities: [140] Device Serial Number 14-02-ec-ff-ff-82-58-cc
	Capabilities: [150] Alternative Routing-ID Interpretation (ARI)
	Capabilities: [160] Single Root I/O Virtualization (SR-IOV)
	Kernel driver in use: ixgbe

04:00.1 Ethernet controller: Intel Corporation 82599ES 10-Gigabit SFI/SFP+ Network Connection (rev 01)
	Subsystem: Hewlett-Packard Company Ethernet 10Gb 2-port 560FLR-SFP+ Adapter
	Flags: bus master, fast devsel, latency 0, IRQ 16
	Memory at 92b00000 (32-bit, non-prefetchable) [size=1M]
	I/O ports at 2000 [size=32]
	Memory at 92d00000 (32-bit, non-prefetchable) [size=16K]
	Capabilities: [40] Power Management version 3
	Capabilities: [50] MSI: Enable- Count=1/1 Maskable+ 64bit+
	Capabilities: [70] MSI-X: Enable+ Count=64 Masked-
	Capabilities: [a0] Express Endpoint, MSI 00
	Capabilities: [e0] Vital Product Data
	Capabilities: [100] Advanced Error Reporting
	Capabilities: [140] Device Serial Number 14-02-ec-ff-ff-82-58-cc
	Capabilities: [150] Alternative Routing-ID Interpretation (ARI)
	Capabilities: [160] Single Root I/O Virtualization (SR-IOV)
	Kernel driver in use: ixgbe
...
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# modinfo ixgbe
filename:       /lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/net/ethernet/intel/ixgbe/ixgbe.ko
version:        4.0.1-k-rh7.1
license:        GPL
description:    Intel(R) 10 Gigabit PCI Express Network Driver
author:         Intel Corporation, <linux.nics@intel.com>
rhelversion:    7.1
srcversion:     1CFEC34DC017FBEDC0500B1
alias:          pci:v00008086d000015ABsv*sd*bc*sc*i*
alias:          pci:v00008086d000015AAsv*sd*bc*sc*i*
alias:          pci:v00008086d00001563sv*sd*bc*sc*i*
alias:          pci:v00008086d00001560sv*sd*bc*sc*i*
alias:          pci:v00008086d0000154Asv*sd*bc*sc*i*
alias:          pci:v00008086d00001557sv*sd*bc*sc*i*
alias:          pci:v00008086d00001558sv*sd*bc*sc*i*
alias:          pci:v00008086d0000154Fsv*sd*bc*sc*i*
alias:          pci:v00008086d0000154Dsv*sd*bc*sc*i*
alias:          pci:v00008086d00001528sv*sd*bc*sc*i*
alias:          pci:v00008086d000010F8sv*sd*bc*sc*i*
alias:          pci:v00008086d0000151Csv*sd*bc*sc*i*
alias:          pci:v00008086d00001529sv*sd*bc*sc*i*
alias:          pci:v00008086d0000152Asv*sd*bc*sc*i*
alias:          pci:v00008086d000010F9sv*sd*bc*sc*i*
alias:          pci:v00008086d00001514sv*sd*bc*sc*i*
alias:          pci:v00008086d00001507sv*sd*bc*sc*i*
alias:          pci:v00008086d000010FBsv*sd*bc*sc*i*
alias:          pci:v00008086d00001517sv*sd*bc*sc*i*
alias:          pci:v00008086d000010FCsv*sd*bc*sc*i*
alias:          pci:v00008086d000010F7sv*sd*bc*sc*i*
alias:          pci:v00008086d00001508sv*sd*bc*sc*i*
alias:          pci:v00008086d000010DBsv*sd*bc*sc*i*
alias:          pci:v00008086d000010F4sv*sd*bc*sc*i*
alias:          pci:v00008086d000010E1sv*sd*bc*sc*i*
alias:          pci:v00008086d000010F1sv*sd*bc*sc*i*
alias:          pci:v00008086d000010ECsv*sd*bc*sc*i*
alias:          pci:v00008086d000010DDsv*sd*bc*sc*i*
alias:          pci:v00008086d0000150Bsv*sd*bc*sc*i*
alias:          pci:v00008086d000010C8sv*sd*bc*sc*i*
alias:          pci:v00008086d000010C7sv*sd*bc*sc*i*
alias:          pci:v00008086d000010C6sv*sd*bc*sc*i*
alias:          pci:v00008086d000010B6sv*sd*bc*sc*i*
depends:        mdio,ptp,dca
intree:         Y
vermagic:       3.10.0-229.el7.x86_64 SMP mod_unload modversions
signer:         CentOS Linux kernel signing key
sig_key:        A6:2A:0E:1D:6A:6E:48:4E:9B:FD:73:68:AF:34:08:10:48:E5:35:E5
sig_hashalgo:   sha256
parm:           max_vfs:Maximum number of virtual functions to allocate per physical function - default is zero and maximum value is 63 (uint)
parm:           allow_unsupported_sfp:Allow unsupported and untested SFP+ modules on 82599-based adapters (uint)
parm:           debug:Debug level (0=none,...,16=all) (int)
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# ifconfig
bond0: flags=5187<UP,BROADCAST,RUNNING,MASTER,MULTICAST>  mtu 1500
        inet 10.200.16.27  netmask 255.255.254.0  broadcast 10.200.17.255
        ether 14:02:ec:82:58:cc  txqueuelen 0  (Ethernet)
        RX packets 2771397  bytes 952268989 (908.1 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 7442390  bytes 4021974277 (3.7 GiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

eno49: flags=6211<UP,BROADCAST,RUNNING,SLAVE,MULTICAST>  mtu 1500
        ether 14:02:ec:82:58:cc  txqueuelen 1000  (Ethernet)
        RX packets 669872  bytes 510127517 (486.4 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 1653218  bytes 324321796 (309.2 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

eno50: flags=6211<UP,BROADCAST,RUNNING,SLAVE,MULTICAST>  mtu 1500
        ether 14:02:ec:82:58:cc  txqueuelen 1000  (Ethernet)
        RX packets 2101525  bytes 442141472 (421.6 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 5789172  bytes 3697652481 (3.4 GiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

lo: flags=73<UP,LOOPBACK,RUNNING>  mtu 65536
        inet 127.0.0.1  netmask 255.0.0.0
        loop  txqueuelen 0  (Local Loopback)
        RX packets 25570  bytes 2306222 (2.1 MiB)
        RX errors 0  dropped 0  overruns 0  frame 0
        TX packets 25570  bytes 2306222 (2.1 MiB)
        TX errors 0  dropped 0 overruns 0  carrier 0  collisions 0

[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# cat /proc/interrupts |grep eno49
 147:         15      44958          0          0          0          0       2224       5214          0          0          0          0          0          0          0          0     151680          0     108157       2743     233116      35352      72048    5784399          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-0
 148:     783488          0         36     142133     253244      25082    1267745     550408          0          0          0          0          0          0          0          0     883455     452800     198122     320072     554383     394113     431253          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-1
 149:      24294    1381206     348493     308052     130943     703404          0     119652          0          0          0          0          0          0          0          0     379836      11715       2350     316074     331709    1005354    1016513          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-2
 150:      10882       2870    1568255      39723     122094     555053       7988      39260          0          0          0          0          0          0          0          0      50400     146813     615491     670551    1173329     575179     441216          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-3
 151:      54478      20796     412914    2154449     228440      10384     582931      15482          0          0          0          0          0          0          0          0     592544       7317     583830     173425     694384     536114          0          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-4
 152:       4556      20180     352073     120238    1586797      40535     188863     163513          0          0          0          0          0          0          0          0      94712       4666     345606     908949     653761    1446059     110751          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-5
 153:      62739      77227          0       1207      22491    2079755     753827    1019152          0          0          0          0          0          0          0          0     241800       9896     272734     687663      13434     297042     484261          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-6
 154:      38408        480     913870          0       9768          0    1266034     147647          2          0          0          0          0          0          0          0     490261     310663     225193     805442      67519     761192     668292     402628          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-7
 155:      55616       2070     180244      65321     660231     969255     230960    1420248          0          3          0          0          0          0          0          0     133707       1895     294213     173568     163688     579554    1087383          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-8
 156:          8      21161          0     247358     157749          0      83006     188647          0          0          2          0          0          0          0          0     159117          0     534146    1033015          0     366769     322863    2888564          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-9
 157:      58665      44017     276940     125593     221230     169855      16255     439791          0          0          0          5          0          0          0          0      27561    2267135     587762     441652     197840     722126     449941          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-10
 158:     387164      79641          0      64567      78937     149554     259298          0          0          0          0          0          3          0          0          0     519833     139756     519859     842755     481668     941600    1541365      40224          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-11
 159:     176877      13533      33959     107215     182169     224291      77162      71961          0          0          0          0          0         28          0          0     272786     506882     754759     816596     892943     760195    1074761     227913          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-12
 160:      90613     360496          0     161927     146887     509154     742019       2981          0          0          0          0          0          0         38          0      30818      70729     656539     336634    1643940     110344    1172487          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-13
 161:     166213          0          0      29421          0        863          0       3114          0          0          0          0          0          0          0         42      91011          0      86188          0      88512     564551          0    5015895          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-14
 162:          7        518      36695      30560      59594     561693     806277     353995          0          0          0          0          0          0          0          0     414588     293315     297165     637873     647445     975379     797445     102089          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-15
 163:     310648     387508     189819     179584     360643      69393      28215      29647          0          0          0          0          0          0          0          0     662622     297661     593714     533727     125691    1074051     823999     323287          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-16
 164:    2896538    1740225          0      22572     713088      29699      16969     135107          0          0          0          0          0          0          0          0      87056      10526     140167      41405      24665     218369       1052          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-17
 165:     266552    1155629    1181434       1477     197336     540369     302838     118696          0          0          0          0          0          0          0          0     379694       1229      57661     368143     449965     324065     675879          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-18
 166:      57252       6153    1573681      27208     708416     400738     136233    1243845          0          0          0          0          0          0          0          0     414203      18978     105231     230153     260825     200082     696980          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-19
 167:      55007      30040      88928     805954     565204     126053     291361    1879730          0          0          0          0          0          0          0          0     179025      47639     142115     332399     589988     595624     257657      11555          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-20
 168:          7      52347      34106       3616    1355080     122230          0    1970515          0          0          0          0          0          0          0          0     335249          0     420341     363426     276133     400721     128568     613925          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-21
 169:     232923     316433     107015      49911     271011     448159     112730     837878          2          0          0          0          0          0          0          0     446886     769422     205089     691811     291570     842319     476410          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-22
 170:      12349      22532      60903      55130     741131     699718     418684     117458          0          2          0          0          0          0          0          0     204810     393131     631865     398633     966435     329164     878059      90419          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-23
 171:     735826      59702     207015    1601840     259276     244694     252906     636408          0          0          2          0          0          0          0          0     354824      65012     252547     278273     217924     408027     411042          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-24
 172:      15432       8125          0      60111     564359     972220    1146455     498656          0          0          0          3          0          0          0          0     233271      28508     399946     457614     380757     408152     710551     158362          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-25
 173:       4363      87093          0     442105     255535     232911     499268    1538861          0          0          0          0          2          0          0          0     159269    1537263     445639     245274     221734     288555      75593          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-26
 174:     358250       2002     134186     225806     373993    1075050    1449450     419708          0          0          0          0          0          2          0          0     403803       8498     352962     398227     315487     251355     239010          0          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-27
 175:     978007      35680     599244     690278     277699     340592      88393     188618          0          0          0          0          0          0          4          0     159319     561939     410758     631431     491711     247848     117690     178143          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-28
 176:     611652      15310          0     134799      23209      62766    1798845     372987          0          0          0          0          0          0          0          2      73964    1554740     215314     114840     349786     440685     101887     119590          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-29
 177:     636542     725587    1306267      25424      98663      53228     249081    1054862          0          0          0          0          0          0          0          0     477543      53857     237379     290530     250092     169219     166221     192334          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-30
 178:     126260       2755          2     171488     355442     725237     961404    1087556          0          0          0          0          0          0          0          0     392502     819213     261163     432404     223399     125189     299263       2395          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49-TxRx-31
 179:          2          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          0          8          0          0          0          0          0          0          0          0  IR-PCI-MSI-edge      eno49
[root@wg-esm-hc-1 ~]#
```

bond 相关内容（不确定是否有影响，故留存）

```
[root@wg-esm-hc-1 ~]# locate bond
/etc/modprobe.d/bond.conf
/etc/sysconfig/network-scripts/ifcfg-bond0
/opt/ops/cmdb/check/physical/network/bond_conf.sh
/opt/ops/cmdb/check/physical/network/bond_traffic.sh
/srv/esm-agent/plugins/bond_mode.py
/usr/bin/bond2team
/usr/include/linux/if_bonding.h
/usr/lib/dracut/modules.d/40network/parse-bond.sh
/usr/lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/media/rc/winbond-cir.ko
/usr/lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/net/bonding
/usr/lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/net/bonding/bonding.ko
/usr/lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/net/ethernet/dec/tulip/winbond-840.ko
/usr/share/doc/initscripts-9.49.24/examples/networking/ifcfg-bond-802.3ad
/usr/share/doc/initscripts-9.49.24/examples/networking/ifcfg-bond-activebackup-arpmon
/usr/share/doc/initscripts-9.49.24/examples/networking/ifcfg-bond-activebackup-miimon
/usr/share/doc/initscripts-9.49.24/examples/networking/ifcfg-bond-slave
/usr/share/doc/iputils-20121221/README.bonding
/usr/share/man/man1/bond2team.1.gz
/usr/src/kernels/3.10.0-229.el7.x86_64/drivers/net/bonding
/usr/src/kernels/3.10.0-229.el7.x86_64/drivers/net/bonding/Makefile
/usr/src/kernels/3.10.0-229.el7.x86_64/drivers/staging/winbond
/usr/src/kernels/3.10.0-229.el7.x86_64/drivers/staging/winbond/Kconfig
/usr/src/kernels/3.10.0-229.el7.x86_64/drivers/staging/winbond/Makefile
/usr/src/kernels/3.10.0-229.el7.x86_64/include/config/bonding.h
/usr/src/kernels/3.10.0-229.el7.x86_64/include/config/winbond
/usr/src/kernels/3.10.0-229.el7.x86_64/include/config/ir/winbond
/usr/src/kernels/3.10.0-229.el7.x86_64/include/config/ir/winbond/cir.h
/usr/src/kernels/3.10.0-229.el7.x86_64/include/config/winbond/840.h
/usr/src/kernels/3.10.0-229.el7.x86_64/include/uapi/linux/if_bonding.h
[root@wg-esm-hc-1 ~]#
```

## 编译安装 pf_ring.ko 内核模块

```
[root@wg-esm-hc-1 fei.sun]# git clone https://github.com/ntop/PF_RING.git
[root@wg-esm-hc-1 fei.sun]# cd PF_RING/kernel/
[root@wg-esm-hc-1 kernel]# make
make -C /lib/modules/3.10.0-229.el7.x86_64/build SUBDIRS=/home/fei.sun/PF_RING/kernel EXTRA_CFLAGS='-I/home/fei.sun/PF_RING/kernel -DGIT_REV="\"dev:05ec8f1b17468d84a27366610e0a226e0a64e55d\""' modules
make[1]: Entering directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
  CC [M]  /home/fei.sun/PF_RING/kernel/pf_ring.o
  Building modules, stage 2.
  MODPOST 1 modules
  CC      /home/fei.sun/PF_RING/kernel/pf_ring.mod.o
  LD [M]  /home/fei.sun/PF_RING/kernel/pf_ring.ko
make[1]: Leaving directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
[root@wg-esm-hc-1 kernel]#
[root@wg-esm-hc-1 kernel]# make install
mkdir -p /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring
cp *.ko /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring
mkdir -p /usr/include/linux
cp linux/pf_ring.h /usr/include/linux
/sbin/depmod 3.10.0-229.el7.x86_64
[root@wg-esm-hc-1 kernel]#
[root@wg-esm-hc-1 kernel]# lsmod |grep pf_ring
[root@wg-esm-hc-1 kernel]#
[root@wg-esm-hc-1 kernel]# insmod pf_ring.ko
[root@wg-esm-hc-1 kernel]# lsmod |grep pf_ring
pf_ring              1234009  0
[root@wg-esm-hc-1 kernel]#
```

## 编译安装用户模块 libpfring 和 libpcap

```
[root@wg-esm-hc-1 kernel]# cd ../userland/lib
[root@wg-esm-hc-1 lib]# ./configure
...
configure: creating ./config.status
config.status: creating lib/Makefile
config.status: creating lib/pfring_config
config.status: creating examples/Makefile
config.status: creating examples_zc/Makefile
config.status: creating c++/Makefile
config.status: creating nbpf/Makefile
config.status: creating wireshark/extcap/Makefile
config.status: creating lib/config.h
[root@wg-esm-hc-1 lib]#
[root@wg-esm-hc-1 lib]# make
cd ../nbpf; make
make[1]: Entering directory `/home/fei.sun/PF_RING/userland/nbpf'
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o nbpf_mod_rdif.o nbpf_mod_rdif.c
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o rules.o rules.c
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o tree_match.o tree_match.c
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o parser.o parser.c
bison -d grammar.y
lex scanner.l
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o lex.yy.o lex.yy.c
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o grammar.tab.o grammar.tab.c
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o nbpf_mod_fiberblaze.o nbpf_mod_fiberblaze.c
gcc -Wall -fPIC -I../lib -I../../kernel   -O2    -c -o nbpf_mod_napatech.o nbpf_mod_napatech.c
ar rs libnbpf.a nbpf_mod_rdif.o rules.o tree_match.o parser.o lex.yy.o grammar.tab.o nbpf_mod_fiberblaze.o nbpf_mod_napatech.o
ar: creating libnbpf.a
ranlib libnbpf.a
gcc -Wall -fPIC -I../lib -I../../kernel   -O2  -g nbpftest.c -o nbpftest libnbpf.a `../lib/pfring_config --libs` -lpthread
make[1]: Leaving directory `/home/fei.sun/PF_RING/userland/nbpf'
ar x ../nbpf/libnbpf.a
cp ../nbpf/nbpf.h .
ar x libs/libpfring_zc_x86_64_core-avx2.a
ar x libs/libpfring_nt_x86_64_core-avx2.a
ar x libs/libpfring_myricom_x86_64_core-avx2.a
ar x libs/libpfring_dag_x86_64_core-avx2.a
ar x libs/libpfring_fiberblaze_x86_64_core-avx2.a
ar x libs/libpfring_accolade_x86_64_core-avx2.a
ar x libs/libnpcap_x86_64_core-avx2.a
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring.c -o pfring.o
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring_mod.c -o pfring_mod.o
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring_utils.c -o pfring_utils.o
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring_mod_stack.c -o pfring_mod_stack.o
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring_hw_filtering.c -o pfring_hw_filtering.o
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring_hw_timestamp.c -o pfring_hw_timestamp.o
gcc -march=native -mtune=native  -Wall -fPIC -I../../kernel -I../libpcap  -D HAVE_PF_RING_ZC -D HAVE_DAG -D HAVE_FIBERBLAZE -D HAVE_ACCOLADE -DHAVE_MYRICOM   -D ENABLE_BPF -D ENABLE_HW_TIMESTAMP -D HAVE_NT  -D HAVE_NPCAP -O2  -c pfring_mod_sysdig.c -o pfring_mod_sysdig.o
=*= making library libpfring.a =*=
ar rs libpfring.a pfring.o pfring_mod.o pfring_utils.o pfring_mod_stack.o pfring_hw_filtering.o pfring_hw_timestamp.o pfring_mod_sysdig.o pfring_zc_dev_e1000.o pfring_zc_dev_e1000e.o pfring_zc_dev_ixgbe.o pfring_zc_dev_igb.o pfring_zc_dev_i40e.o pfring_zc_dev_fm10k.o pfring_zc_dev_rss.o pfring_zc_dev_sal.o pfring_mod_zc.o pfring_mod_zc_dev.o pfring_mod_zc_spsc.o pfring_zc_cluster.o pfring_zc_mm.o uio_lib.o hugetlb_lib.o numa_lib.o pfring_zc_kvm.o pfring_zc_kvm_utils.o silicom_ts_card.o  pfring_mod_dag.o  pfring_mod_fiberblaze.o  pfring_mod_nt.o   pfring_mod_accolade.o  pfring_mod_myricom.o    pfring_mod_timeline.o npcapextract_lib.o index_match.o npcap_utils.o npcap_compression.o lzf_c.o lzf_d.o  `ar t ../nbpf/libnbpf.a | grep -F .o | tr '\n' ' '`
ar: creating libpfring.a
ranlib libpfring.a
=*= making library libpfring.so =*=
gcc -g -shared pfring.o pfring_mod.o pfring_utils.o pfring_mod_stack.o pfring_hw_filtering.o pfring_hw_timestamp.o pfring_mod_sysdig.o pfring_zc_dev_e1000.o pfring_zc_dev_e1000e.o pfring_zc_dev_ixgbe.o pfring_zc_dev_igb.o pfring_zc_dev_i40e.o pfring_zc_dev_fm10k.o pfring_zc_dev_rss.o pfring_zc_dev_sal.o pfring_mod_zc.o pfring_mod_zc_dev.o pfring_mod_zc_spsc.o pfring_zc_cluster.o pfring_zc_mm.o uio_lib.o hugetlb_lib.o numa_lib.o pfring_zc_kvm.o pfring_zc_kvm_utils.o silicom_ts_card.o  pfring_mod_dag.o  pfring_mod_fiberblaze.o  pfring_mod_nt.o   pfring_mod_accolade.o  pfring_mod_myricom.o    pfring_mod_timeline.o npcapextract_lib.o index_match.o npcap_utils.o npcap_compression.o lzf_c.o lzf_d.o  `ar t ../nbpf/libnbpf.a | grep -F .o | tr '\n' ' '` -lpthread  -lrt -ldl -lm -ldl -lm -ldl   -o libpfring.so
[root@wg-esm-hc-1 lib]#
[root@wg-esm-hc-1 lib]# make install
ar x ../nbpf/libnbpf.a
cp ../nbpf/nbpf.h .
ar x libs/libpfring_zc_x86_64_core-avx2.a
ar x libs/libpfring_nt_x86_64_core-avx2.a
ar x libs/libpfring_myricom_x86_64_core-avx2.a
ar x libs/libpfring_dag_x86_64_core-avx2.a
ar x libs/libpfring_fiberblaze_x86_64_core-avx2.a
ar x libs/libpfring_accolade_x86_64_core-avx2.a
ar x libs/libnpcap_x86_64_core-avx2.a
=*= making library libpfring.a =*=
ar rs libpfring.a pfring.o pfring_mod.o pfring_utils.o pfring_mod_stack.o pfring_hw_filtering.o pfring_hw_timestamp.o pfring_mod_sysdig.o pfring_zc_dev_e1000.o pfring_zc_dev_e1000e.o pfring_zc_dev_ixgbe.o pfring_zc_dev_igb.o pfring_zc_dev_i40e.o pfring_zc_dev_fm10k.o pfring_zc_dev_rss.o pfring_zc_dev_sal.o pfring_mod_zc.o pfring_mod_zc_dev.o pfring_mod_zc_spsc.o pfring_zc_cluster.o pfring_zc_mm.o uio_lib.o hugetlb_lib.o numa_lib.o pfring_zc_kvm.o pfring_zc_kvm_utils.o silicom_ts_card.o  pfring_mod_dag.o  pfring_mod_fiberblaze.o  pfring_mod_nt.o   pfring_mod_accolade.o  pfring_mod_myricom.o    pfring_mod_timeline.o npcapextract_lib.o index_match.o npcap_utils.o npcap_compression.o lzf_c.o lzf_d.o  `ar t ../nbpf/libnbpf.a | grep -F .o | tr '\n' ' '`
ranlib libpfring.a
mkdir -p //usr/local/include
cp pfring.h pfring_mod_sysdig.h pfring_zc.h pfring_zc.h ../nbpf/nbpf.h //usr/local/include/
cp: warning: source file ‘pfring_zc.h’ specified more than once
mkdir -p //usr/local/lib
cp libpfring.a //usr/local/lib/
=*= making library libpfring.so =*=
gcc -g -shared pfring.o pfring_mod.o pfring_utils.o pfring_mod_stack.o pfring_hw_filtering.o pfring_hw_timestamp.o pfring_mod_sysdig.o pfring_zc_dev_e1000.o pfring_zc_dev_e1000e.o pfring_zc_dev_ixgbe.o pfring_zc_dev_igb.o pfring_zc_dev_i40e.o pfring_zc_dev_fm10k.o pfring_zc_dev_rss.o pfring_zc_dev_sal.o pfring_mod_zc.o pfring_mod_zc_dev.o pfring_mod_zc_spsc.o pfring_zc_cluster.o pfring_zc_mm.o uio_lib.o hugetlb_lib.o numa_lib.o pfring_zc_kvm.o pfring_zc_kvm_utils.o silicom_ts_card.o  pfring_mod_dag.o  pfring_mod_fiberblaze.o  pfring_mod_nt.o   pfring_mod_accolade.o  pfring_mod_myricom.o    pfring_mod_timeline.o npcapextract_lib.o index_match.o npcap_utils.o npcap_compression.o lzf_c.o lzf_d.o  `ar t ../nbpf/libnbpf.a | grep -F .o | tr '\n' ' '` -lpthread  -lrt -ldl -lm -ldl -lm -ldl   -o libpfring.so
mkdir -p //usr/local/lib
cp libpfring.so //usr/local/lib/
[root@wg-esm-hc-1 lib]#


[root@wg-esm-hc-1 lib]# cd ../libpcap
[root@wg-esm-hc-1 libpcap]#
[root@wg-esm-hc-1 libpcap]# ./configure
...
configure: creating ./config.status
config.status: creating Makefile
config.status: creating pcap-filter.manmisc
config.status: creating pcap-linktype.manmisc
config.status: creating pcap-tstamp.manmisc
config.status: creating pcap-savefile.manfile
config.status: creating pcap.3pcap
config.status: creating pcap_compile.3pcap
config.status: creating pcap_datalink.3pcap
config.status: creating pcap_dump_open.3pcap
config.status: creating pcap_get_tstamp_precision.3pcap
config.status: creating pcap_list_datalinks.3pcap
config.status: creating pcap_list_tstamp_types.3pcap
config.status: creating pcap_open_dead.3pcap
config.status: creating pcap_open_offline.3pcap
config.status: creating pcap_set_tstamp_precision.3pcap
config.status: creating pcap_set_tstamp_type.3pcap
config.status: creating config.h
config.status: executing default-1 commands
[root@wg-esm-hc-1 libpcap]#
[root@wg-esm-hc-1 libpcap]# make
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./pcap-linux.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./pcap-usb-linux.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./pcap-can-linux.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./pcap-netfilter-linux.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./fad-getad.c
if grep GIT ./VERSION >/dev/null; then \
	read ver <./VERSION; \
	echo $ver | tr -d '\012'; \
	date +_%Y_%m_%d; \
else \
	cat ./VERSION; \
fi | sed -e 's/.*/static const char pcap_version_string[] = "libpcap version &";/' > version.h
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./pcap.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./inet.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./gencode.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./optimize.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./nametoaddr.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./etherent.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./savefile.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./sf-pcap.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./sf-pcap-ng.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./pcap-common.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./bpf_image.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c ./bpf_dump.c
./runlex.sh flex -Ppcap_ -oscanner.c scanner.l
mv scanner.c scanner.c.bottom
cat ./scanner.c.top scanner.c.bottom > scanner.c
bison -y -p pcap_ -d grammar.y
conflicts: 38 shift/reduce
mv y.tab.c grammar.c
mv y.tab.h tokdefs.h
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c scanner.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -Dyylval=pcap_lval -c grammar.c
rm -f bpf_filter.c
ln -s ./bpf/net/bpf_filter.c bpf_filter.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c bpf_filter.c
if grep GIT ./VERSION >/dev/null; then \
	read ver <./VERSION; \
	echo $ver | tr -d '\012'; \
	date +_%Y_%m_%d; \
else \
	cat ./VERSION; \
fi | sed -e 's/.*/char pcap_version[] = "&";/' > version.c
gcc -fpic -I. -I ../../kernel -I ../lib  -DHAVE_CONFIG_H  -D_U_="__attribute__((unused))" -DHAVE_PF_RING   -g -O2    -c version.c
ar rc libpcap.a pcap-linux.o pcap-usb-linux.o pcap-can-linux.o pcap-netfilter-linux.o fad-getad.o pcap.o inet.o gencode.o optimize.o nametoaddr.o etherent.o savefile.o sf-pcap.o sf-pcap-ng.o pcap-common.o bpf_image.o bpf_dump.o  scanner.o grammar.o bpf_filter.o version.o
ranlib libpcap.a
VER=`cat ./VERSION`; \
MAJOR_VER=`sed 's/\([0-9][0-9]*\)\..*/\1/' ./VERSION`; \
gcc -shared -Wl,-soname,libpcap.so.$MAJOR_VER  \
    -o libpcap.so.$VER pcap-linux.o pcap-usb-linux.o pcap-can-linux.o pcap-netfilter-linux.o fad-getad.o pcap.o inet.o gencode.o optimize.o nametoaddr.o etherent.o savefile.o sf-pcap.o sf-pcap-ng.o pcap-common.o bpf_image.o bpf_dump.o  scanner.o grammar.o bpf_filter.o version.o   ../lib/libpfring.a -lpthread -lrt   -lrt -ldl
./config.status --file=pcap-config.tmp:./pcap-config.in
config.status: creating pcap-config.tmp
mv pcap-config.tmp pcap-config
chmod a+x pcap-config
[root@wg-esm-hc-1 libpcap]#
```

或者直接在 PF_RING/userland 目录下执行

```
[root@wg-esm-hc-1 userland]# make
[root@wg-esm-hc-1 userland]# make install
```

此时

```
[root@wg-esm-hc-1 ~]# ll /usr/local/lib/
total 4252
-rw-r--r-- 1 root root 1683210 Apr 17 14:55 libpcap.a
lrwxrwxrwx 1 root root      12 Apr 17 14:55 libpcap.so -> libpcap.so.1
lrwxrwxrwx 1 root root      16 Apr 17 14:55 libpcap.so.1 -> libpcap.so.1.7.4
-rwxr-xr-x 1 root root 1447327 Apr 17 14:55 libpcap.so.1.7.4
-rw-r--r-- 1 root root  696898 Apr 17 14:55 libpfring.a
-rwxr-xr-x 1 root root  517730 Apr 17 14:55 libpfring.so
[root@wg-esm-hc-1 ~]#
```

此时

```
[root@wg-esm-hc-1 ~]# locate libpcap.so
/home/fei.sun/PF_RING/userland/libpcap-1.7.4/libpcap.so.1.7.4
/usr/lib64/libpcap.so.1
/usr/lib64/libpcap.so.1.5.3
/usr/local/lib/libpcap.so
/usr/local/lib/libpcap.so.1
/usr/local/lib/libpcap.so.1.7.4
[root@wg-esm-hc-1 ~]#
```

## 编译安装 intel 网卡驱动

```
[root@wg-esm-hc-1 PF_RING]# cd drivers/
[root@wg-esm-hc-1 drivers]# ll
total 12
drwxr-xr-x 9 root root 4096 Apr 17 14:40 intel
-rw-r--r-- 1 root root  198 Apr 17 14:40 Makefile
-rw-r--r-- 1 root root  300 Apr 17 14:40 README
[root@wg-esm-hc-1 drivers]# cat README
This directory contains drivers for popular 1/10/40/100 Gbit adapters with PF_RING ZC support.
They can be used as standard drivers, or in zero-copy mode
opening the device with the 'zc:' prefix (e.g. zc:eth1).

Please note that opening the device in zero-copy mode disconnects it from
the system.

[root@wg-esm-hc-1 drivers]#
[root@wg-esm-hc-1 drivers]#
[root@wg-esm-hc-1 drivers]# cat Makefile
all:
#	cd broadcom; make
	cd intel; make
#	cd chelsio/cxgb4; make
#	cd myricom; make

clean:
#	cd broadcom; make clean
	cd intel; make clean
#	cd chelsio/cxgb4; make clean
#	cd myricom; make clean

[root@wg-esm-hc-1 drivers]#
[root@wg-esm-hc-1 drivers]# make
...
cd ixgbe/ixgbe-4.1.5-zc/src; make
make[2]: Entering directory `/home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src'
make -C /lib/modules/3.10.0-229.el7.x86_64/build SUBDIRS=/home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src modules
make[3]: Entering directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_main.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_common.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_api.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_param.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_lib.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_ethtool.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/kcompat.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_82598.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_82599.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_x540.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_x550.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_sriov.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_mbx.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb_82598.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb_82599.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_sysfs.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_procfs.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_phy.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb_nl.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_fcoe.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_debugfs.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_ptp.o
  LD [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe.o
  Building modules, stage 2.
  MODPOST 1 modules
  CC      /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe.mod.o
  LD [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe.ko
make[3]: Leaving directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
make[2]: Leaving directory `/home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src'
...
[root@wg-esm-hc-1 drivers]#
```

或者直接

```
[root@wg-esm-hc-1 drivers]# cd intel/ixgbe/ixgbe-4.1.5-zc/src/
[root@wg-esm-hc-1 src]# make
make -C /lib/modules/3.10.0-229.el7.x86_64/build SUBDIRS=/home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src modules
make[1]: Entering directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_main.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_common.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_api.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_param.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_lib.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_ethtool.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/kcompat.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_82598.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_82599.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_x540.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_x550.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_sriov.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_mbx.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb_82598.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb_82599.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_sysfs.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_procfs.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_phy.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_dcb_nl.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_fcoe.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_debugfs.o
  CC [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe_ptp.o
  LD [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe.o
  Building modules, stage 2.
  MODPOST 1 modules
  CC      /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe.mod.o
  LD [M]  /home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src/ixgbe.ko
make[1]: Leaving directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
[root@wg-esm-hc-1 src]#
[root@wg-esm-hc-1 src]#
[root@wg-esm-hc-1 src]# make install
make -C /lib/modules/3.10.0-229.el7.x86_64/build SUBDIRS=/home/fei.sun/PF_RING/drivers/intel/ixgbe/ixgbe-4.1.5-zc/src modules
make[1]: Entering directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
  Building modules, stage 2.
  MODPOST 1 modules
make[1]: Leaving directory `/usr/src/kernels/3.10.0-229.el7.x86_64'
gzip -c ../ixgbe.7 > ixgbe.7.gz
# remove all old versions of the driver
find /lib/modules/3.10.0-229.el7.x86_64 -name ixgbe.ko -exec rm -f {} \; || true
find /lib/modules/3.10.0-229.el7.x86_64 -name ixgbe.ko.gz -exec rm -f {} \; || true
install -D -m 644 ixgbe.ko /lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/net/ethernet/intel/ixgbe/ixgbe.ko
/sbin/depmod -a 3.10.0-229.el7.x86_64 || true
install -D -m 644 ixgbe.7.gz /usr/share/man/man7/ixgbe.7.gz
[root@wg-esm-hc-1 src]#
```

到此 `ixgbe.ko` 网卡驱动内核模块已经保存到正确的系统目录中了（**注意：会覆盖之前的 `ixgbe.ko` 版本**）

最后，需要通过 `load_driver.sh` 脚本进行内核模块和系统配置（**注意：需要在脚本的最后增加 `systemctl restart network` 命令，否则执行脚本后将无法连接到目标机器**）；会有短暂断网发生；

执行脚本时的调试打印

```
[root@wg-esm-hc-1 src]# cat ld.txt
+ FAMILY=ixgbe
+ rmmod ixgbe
+ rmmod pf_ring
+ insmod ../../../../../kernel/pf_ring.ko
+ modprobe ptp
+ modprobe vxlan
+ modprobe dca
+ insmod ./ixgbe.ko RSS=1,1,1,1
+ sleep 1
+ killall irqbalance
irqbalance: no process found
++ cat /proc/net/dev
++ grep :
++ grep -v lo
++ grep -v sit
++ awk -F: '{print $1}'
++ tr -d ' '
+ INTERFACES='bond0
eno49
eno50'
+ for IF in '$INTERFACES'
++ ethtool -i bond0
++ grep ixgbe
++ wc -l
+ TOCONFIG=0
+ '[' 0 -eq 1 ']'
+ for IF in '$INTERFACES'
++ ethtool -i eno49
++ grep ixgbe
++ wc -l
+ TOCONFIG=1
+ '[' 1 -eq 1 ']'
+ printf 'Configuring %s\n' eno49
Configuring eno49
+ ifconfig eno49 up
+ sleep 1
+ bash ../scripts/set_irq_affinity eno49
IFACE CORE MASK -> FILE
=======================
eno49 0 1 -> /proc/irq/147/smp_affinity
+ ethtool -G eno49 rx 32768
+ ethtool -G eno49 tx 32768
+ ethtool -K eno49 rxvlan off
+ ethtool -A eno49 rx off
+ ethtool -A eno49 tx off
+ for IF in '$INTERFACES'
++ ethtool -i eno50
++ grep ixgbe
++ wc -l
+ TOCONFIG=1
+ '[' 1 -eq 1 ']'
+ printf 'Configuring %s\n' eno50
Configuring eno50
+ ifconfig eno50 up
+ sleep 1
+ bash ../scripts/set_irq_affinity eno50
IFACE CORE MASK -> FILE
=======================
eno50 0 1 -> /proc/irq/149/smp_affinity
+ ethtool -G eno50 rx 32768
+ ethtool -G eno50 tx 32768
+ ethtool -K eno50 rxvlan off
+ ethtool -A eno50 rx off
+ ethtool -A eno50 tx off
+ HUGEPAGES_NUM=1024
+ HUGEPAGES_PATH=/dev/hugepages
+ sync
+ echo 3
+ echo 1024
++ cat /proc/mounts
++ grep hugetlbfs
++ grep /dev/hugepages
++ wc -l
+ '[' 1 -eq 0 ']'
++ grep HugePages_Total /sys/devices/system/node/node0/meminfo
++ cut -d : -f 2
++ sed 's/ //g'
+ HUGEPAGES_AVAIL=512
+ '[' 512 -ne 1024 ']'
+ printf 'Warning: %s hugepages available, %s requested\n' 512 1024
Warning: 512 hugepages available, 1024 requested
+ systemctl restart network
[root@wg-esm-hc-1 src]#
```

执行后系统信息如下

```
[root@wg-esm-hc-1 ~]# lsmod |grep pf_ring
pf_ring              1234009  2
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# lsmod |grep ixgbe
ixgbe                 302621  0
vxlan                  37409  1 ixgbe
ptp                    18933  1 ixgbe
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# modinfo pf_ring
filename:       /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring/pf_ring.ko
alias:          net-pf-27
version:        6.5.0
description:    Packet capture acceleration and analysis
author:         ntop.org
license:        GPL
rhelversion:    7.1
srcversion:     414F094C8FD5E8D55A89517
depends:
vermagic:       3.10.0-229.el7.x86_64 SMP mod_unload modversions
parm:           min_num_slots:Min number of ring slots (uint)
parm:           perfect_rules_hash_size:Perfect rules hash size (uint)
parm:           enable_tx_capture:Set to 1 to capture outgoing packets (uint)
parm:           enable_frag_coherence:Set to 1 to handle fragments (flow coherence) in clusters (uint)
parm:           enable_ip_defrag:Set to 1 to enable IP defragmentation(only rx traffic is defragmentead) (uint)
parm:           quick_mode:Set to 1 to run at full speed but with upto one socket per interface (uint)
parm:           force_ring_lock:Set to 1 to force ring locking (automatically enable with rss) (uint)
parm:           enable_debug:Set to 1 to enable PF_RING debug tracing into the syslog, 2 for more verbosity (uint)
parm:           transparent_mode:(deprecated) (uint)
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# modinfo ixgbe
filename:       /lib/modules/3.10.0-229.el7.x86_64/kernel/drivers/net/ethernet/intel/ixgbe/ixgbe.ko
version:        4.1.5
license:        GPL
description:    Intel(R) 10 Gigabit PCI Express Network Driver
author:         Intel Corporation, <linux.nics@intel.com>
rhelversion:    7.1
srcversion:     368F6527E52D2C4A14E4CAD
alias:          pci:v00008086d000015ADsv*sd*bc*sc*i*
alias:          pci:v00008086d00001560sv*sd*bc*sc*i*
alias:          pci:v00008086d00001558sv*sd*bc*sc*i*
alias:          pci:v00008086d0000154Asv*sd*bc*sc*i*
alias:          pci:v00008086d00001557sv*sd*bc*sc*i*
alias:          pci:v00008086d0000154Fsv*sd*bc*sc*i*
alias:          pci:v00008086d0000154Dsv*sd*bc*sc*i*
alias:          pci:v00008086d00001528sv*sd*bc*sc*i*
alias:          pci:v00008086d000010F8sv*sd*bc*sc*i*
alias:          pci:v00008086d0000151Csv*sd*bc*sc*i*
alias:          pci:v00008086d00001529sv*sd*bc*sc*i*
alias:          pci:v00008086d0000152Asv*sd*bc*sc*i*
alias:          pci:v00008086d000010F9sv*sd*bc*sc*i*
alias:          pci:v00008086d00001514sv*sd*bc*sc*i*
alias:          pci:v00008086d00001507sv*sd*bc*sc*i*
alias:          pci:v00008086d000010FBsv*sd*bc*sc*i*
alias:          pci:v00008086d00001517sv*sd*bc*sc*i*
alias:          pci:v00008086d000010FCsv*sd*bc*sc*i*
alias:          pci:v00008086d000010F7sv*sd*bc*sc*i*
alias:          pci:v00008086d00001508sv*sd*bc*sc*i*
alias:          pci:v00008086d000010DBsv*sd*bc*sc*i*
alias:          pci:v00008086d000010F4sv*sd*bc*sc*i*
alias:          pci:v00008086d000010E1sv*sd*bc*sc*i*
alias:          pci:v00008086d000010F1sv*sd*bc*sc*i*
alias:          pci:v00008086d000010ECsv*sd*bc*sc*i*
alias:          pci:v00008086d000010DDsv*sd*bc*sc*i*
alias:          pci:v00008086d0000150Bsv*sd*bc*sc*i*
alias:          pci:v00008086d000010C8sv*sd*bc*sc*i*
alias:          pci:v00008086d000010C7sv*sd*bc*sc*i*
alias:          pci:v00008086d000010C6sv*sd*bc*sc*i*
alias:          pci:v00008086d000010B6sv*sd*bc*sc*i*
depends:        ptp,vxlan
vermagic:       3.10.0-229.el7.x86_64 SMP mod_unload modversions
parm:           allow_tap_1g:Allow 1Gbit/s TAP disabling atonegotiation on 82599 based adapters (uint)
parm:           InterruptType:Change Interrupt Mode (0=Legacy, 1=MSI, 2=MSI-X), default IntMode (deprecated) (array of int)
parm:           IntMode:Change Interrupt Mode (0=Legacy, 1=MSI, 2=MSI-X), default 2 (array of int)
parm:           MQ:Disable or enable Multiple Queues, default 1 (array of int)
parm:           DCA:Disable or enable Direct Cache Access, 0=disabled, 1=descriptor only, 2=descriptor and data (array of int)
parm:           RSS:Number of Receive-Side Scaling Descriptor Queues, default 0=number of cpus (array of int)
parm:           VMDQ:Number of Virtual Machine Device Queues: 0/1 = disable, 2-16 enable (default=8) (array of int)
parm:           max_vfs:Number of Virtual Functions: 0 = disable (default), 1-63 = enable this many VFs (array of int)
parm:           VEPA:VEPA Bridge Mode: 0 = VEB (default), 1 = VEPA (array of int)
parm:           InterruptThrottleRate:Maximum interrupts per second, per vector, (0,1,956-488281), default 1 (array of int)
parm:           LLIPort:Low Latency Interrupt TCP Port (0-65535) (array of int)
parm:           LLIPush:Low Latency Interrupt on TCP Push flag (0,1) (array of int)
parm:           LLISize:Low Latency Interrupt on Packet Size (0-1500) (array of int)
parm:           LLIEType:Low Latency Interrupt Ethernet Protocol Type (array of int)
parm:           LLIVLANP:Low Latency Interrupt on VLAN priority threshold (array of int)
parm:           FdirPballoc:Flow Director packet buffer allocation level:
			1 = 8k hash filters or 2k perfect filters
			2 = 16k hash filters or 4k perfect filters
			3 = 32k hash filters or 8k perfect filters (array of int)
parm:           AtrSampleRate:Software ATR Tx packet sample rate (array of int)
parm:           FCoE:Disable or enable FCoE Offload, default 1 (array of int)
parm:           LRO:Large Receive Offload (0,1), default 1 = on (array of int)
parm:           allow_unsupported_sfp:Allow unsupported and untested SFP+ modules on 82599 based adapters, default 0 = Disable (array of int)
parm:           dmac_watchdog:DMA coalescing watchdog in microseconds (0,41-10000), default 0 = off (array of int)
parm:           vxlan_rx:VXLAN receive checksum offload (0,1), default 1 = Enable (array of int)
parm:           numa_cpu_affinity:Comma separated list of core ids where per-adapter memory will be allocated (array of int)
parm:           low_latency_tx:Set to 1 to reduce transmission latency, minimize PCIe overhead otherwise (uint)
parm:           enable_debug:Set to 1 to enable debug tracing into the syslog (uint)
[root@wg-esm-hc-1 ~]#
[root@wg-esm-hc-1 ~]# ethtool -i eno49
driver: ixgbe
version: 4.1.5
firmware-version: 0x80000887, 1.1200.0
bus-info: 0000:04:00.0
supports-statistics: yes
supports-test: yes
supports-eeprom-access: yes
supports-register-dump: yes
supports-priv-flags: no
[root@wg-esm-hc-1 ~]#
```

之后就可以进行 PF_RING 功能测试了；



