# PF_RING User's Guide

> 原文地址：[这里](https://github.com/ntop/PF_RING/wiki)

PF_RING 的定位是高速的 packet 捕获库，其可将常见的商用（commodity）PC 转变成高效、廉价的网络测量设备（network measurement box），适用于 packet 和活动 traffic 的分析和操控；更进一步，PF_RING 打开了一个全新的市场，因为其实现了基于数行代码令诸如 traffic balancers 或 packet filters 这类高效率应用程序的创建成为了可能；

更多信息详见 [Vanilla PF_RING](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/Vanilla%20PF_RING.md) 和 [PF_RING ZC](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/PF_RING%20ZC.md) ；

针对 PF_RING 所支持的 NICs 的比较，详见[这里](https://github.com/ntop/PF_RING/wiki/Drivers-and-Modules)；


PF_RING 用户手册目录：

* [Vanilla PF_RING](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/Vanilla%20PF_RING.md)
* [PF_RING ZC](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/PF_RING%20ZC.md)
* [Source Code Compilation/Installation](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/PF_RING%20Installation.md)
* [Packages Installation](https://github.com/ntop/PF_RING/blob/dev/doc/README.apt_rpm_packages.md)
* [Drivers/Modules](https://github.com/ntop/PF_RING/wiki/Drivers-and-Modules)
    * [Accolade](https://github.com/ntop/PF_RING/blob/dev/doc/README.accolade.md)
    * [Endace DAG](https://github.com/ntop/PF_RING/blob/dev/doc/README.dag.md)
    * [Exablaze](https://github.com/ntop/PF_RING/blob/dev/doc/README.exablaze.md)
    * [Fiberblaze](https://github.com/ntop/PF_RING/blob/dev/doc/README.fiberblaze.md)
    * [Myricom](https://github.com/ntop/PF_RING/blob/dev/doc/README.myricom.md)
    * [Napatech](https://github.com/ntop/PF_RING/blob/dev/doc/README.napatech.md)
    * [Stack](https://github.com/ntop/PF_RING/blob/dev/doc/README.stack.md)
    * [Timeline](https://github.com/ntop/PF_RING/blob/dev/doc/README.timeline.md)
* [nBPF](https://github.com/ntop/PF_RING/blob/dev/userland/nbpf/README.md)
    * [FM10K](https://github.com/ntop/PF_RING/blob/dev/userland/nbpf/README.fm10k.md)
* VM Support
    * [PCI Passthrough](https://github.com/ntop/PF_RING/blob/dev/doc/README.virsh_hostdev.md)
    * [ZC on QEMU/KVM](https://github.com/ntop/PF_RING/blob/dev/doc/README.kvm_zc.md)
* 3rd Party Integration
    * [Bro](https://github.com/ntop/PF_RING/blob/dev/doc/README.bro.md)
    * [Suricata](https://github.com/ntop/PF_RING/blob/dev/doc/README.suricata.md)
    * [Wireshark](https://github.com/ntop/PF_RING/blob/dev/userland/wireshark/extcap/README.md)
* Misc
    * [RSS](https://github.com/moooofly/MarkSomethingDown/blob/master/Linux/RSS_Receive%20Side%20Scaling.md)
    * [Hugepages](https://github.com/ntop/PF_RING/blob/dev/doc/README.hugepages.md)
    * [Docker](https://github.com/ntop/PF_RING/blob/dev/doc/README.docker.md)
    * [DNA to ZC Migration](https://github.com/ntop/PF_RING/blob/dev/doc/README.DNA_to_ZC.md)
* [API Documentation](http://www.ntop.org/pfring_api/files.html)
