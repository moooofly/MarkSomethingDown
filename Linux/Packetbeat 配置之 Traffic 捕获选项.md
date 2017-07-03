# Packetbeat 之 Traffic 捕获选项设置

> 原文地址：[这里](https://www.elastic.co/guide/en/beats/packetbeat/current/capturing-options.html)

存在两种部署 Packetbeat 的方式：

- 在使用专用服务器（dedicated servers）时，可以通过镜像端口或 tap 设备来获取 traffic ；
- 直接部署在你的应用服务器（application servers）上；

第一种方式的巨大优势在于：针对你的应用服务器不会产生任何额外的开销；但其要求专用网络设备（dedicated networking gear）的支持，而这也正是在云端部署时通常不具备的；

在上述两种情况下，sniffing 性能（从网络上进行被动包读取）非常重要；对于专用服务器的情况，更好的 sniffing 性能意味着对硬件本身的要求可以更低；对于 Packetbeat 被安装在已经部署好的应用服务器的情况，更好的 sniffing 性能则意味着更低的额外开销；

当前，Packetbeat 支持多种选项控制 traffic 捕获行为：

- `pcap` - 底层使用 **libpcap** 库，在大多数平台上均能工作，但不是三种选项中最快的；
- `af_packet` - 底层使用 memory mapped sniffing 机制（即 mmap() + socket API ，注册一种新的 socket 类型 `PACKET_MMAP`）；该选项比 **libpcap** 更快，并且对 kernel 模块没有额外要求，但该选项是特定于 Linux 的；
- `pf_ring` - 基于 ntop.org 的[项目](http://www.ntop.org/products/packet-capture/pf_ring/)实现；该选项提供了最快的 sniffing 速度；但是其要求应用程序的重新编译，且对 kernel 模块有要求；另外该方案是也特定于 Linux 的；

选项 `pf_ring` 实现了采用标准硬件上就可以达到 Gbps 级别的 sniffing 速度；但其要求基于 ntop 库重新编译 Packetbeat ，因此当前并未被 Elastic 以官方形式支持；

选项 `af_packet` 被称作 "memory-mapped sniffing" ，其利用 Linux 自身的[特性](http://lxr.free-electrons.com/source/Documentation/networking/packet_mmap.txt)实现；该选项对于专用服务器和将 Packetbeat 部署在当前应用服务器上的情况来说，可能是最佳 sniffing 模式；

该选项（`af_packet`）的实现方式为：程序针对 kernel 空间和用户空间的相同内存区域进行映射，并在该映射区域中实现一个简单的 circular buffer ，之后 kernel 会将 packets 写入该 circular buffer 之中，用户空间程序再从中进行读取；（用户空间程序会）先通过 `poll` 系统调用获取首个可用包的通知，之后的其它包就可以简单的通过内存访问方式进行读取；

该选项没有 `pf_ring` 方式快（在其开始丢包前，最快可达 200k pps），但该选项的使用不要求（程序的）重新编译和 kernel 模块的支持；因此仍可以将其视作针对 **libpcap** 的重要改进；

基于 `af_packet` 的 sniffer 还可以进一步调优：使用更多的内存以达到更高的性能；因为 circular buffer 的 size 越大，系统调用的需求量就会越少，也就意味着更少的 CPU 周期被消耗；该 buffer 的默认大小为 30 MB ，但是你可以按如下方式进行增大调整：

```
packetbeat.interfaces.device: eth0
packetbeat.interfaces.type: af_packet
packetbeat.interfaces.buffer_size_mb: 100
```

更多配置信息详见 [Network Device Configuration](https://www.elastic.co/guide/en/beats/packetbeat/current/configuration-interfaces.html) ；
