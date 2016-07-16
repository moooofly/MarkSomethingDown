
# Netlink Protocol Library Suite (libnl)

## Summary

The libnl suite is a collection of libraries providing APIs to netlink protocol based Linux kernel interfaces.
libnl 套装由一组库的集合构成，其基于 Linux 内核接口提供了访问 netlink 协议的 API ；

Netlink is a IPC mechanism primarly between the kernel and user space processes. It was designed to be a more flexible successor to ioctl to provide mainly networking related kernel configuration and monitoring interfaces.
netlink 是一种 IPC 机制，主要用于内核空间进程与用户空间进程的交互；
netlink 是作为更为灵活的、ioctl 的后继者被设计出来的，主要用于提供与网络相关的内核配置能力和监控接口；


## Libraries

The interfaces are split into several small libraries to not force applications to link against a single, bloated library.
接口按照功能被拆分成许多个短小精悍的库，以便应用程序可以按需链接；

### libnl
Core library implementing the fundamentals required to use the netlink protocol such as socket handling, message construction and parsing, and sending and receiving of data. This library is kept small and minimalistic. Other libraries of the suite depend on this library.
libnl 作为核心库实现了 netlink 协议使用所需的基础功能，例如 socket 处理，消息构建和解析，数据发送和接收；该库的设计目标就是保持短小精悍和足够抽象；libnl 套装中的其它库都依赖此库；

### libnl-route
API to the configuration interfaces of the NETLINK_ROUTE family including network interfaces, routes, addresses, neighbours, and traffic control.
提供可以访问 NETLINK_ROUTE 族配置接口的 API ，包括 network interfaces, routes, addresses, neighbours 和 traffic control 几个部分；

### libnl-genl
API to the generic netlink protocol, an extended version of the netlink protocol.
提供了针对 generic netlink 协议的 API ；

### libnl-nf
API to netlink based netfilter configuration and monitoring interfaces (conntrack, log, queue)
提供针对 netfilter 配置的 API（基于 netlink），以及针对 conntrack, log, queue 的监控接口；

## Installation

The easiest method of installing the libnl library suite is to use the existing packages of your linux distribution. 

## Release
The latest stable release is: 3.2.25 (Released on Jul 16, 2014)


官网地址：[这里](http://www.infradead.org/~tgr/libnl/)