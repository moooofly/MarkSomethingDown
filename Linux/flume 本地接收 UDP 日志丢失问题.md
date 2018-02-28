# flume 本地接收 UDP 日志丢失问题

## 思路

- 通过 127.0.0.1 来发送
- 数据交互环节
    - 应用发日志到内核协议栈
    - 内核协议栈处理
    - 内核协议栈发日志到 flume
- 通过 dropwatch 观察发包时丢包路径分布
- 针对丢包最多的点查看内核源码进行分析
- 通过 systemtap 脚本跟踪内核源码中和丢包相关的关键参数
- 确定 sysctl 配置的内核参数和上述变量的关系
- 分析 systemtap 输出结果（得出：只要发送端跟接收端的速度不匹配，情况一直会出现）
- 通过 strace 查看 flume 的处理性能（估算真实能够处理的能力接近 8000/s 左右）

## dropwatch 相关

- dropwatch 的 man 手册

```
NAME
       dropwatch - kernel dropped packet monitoring utility

SYNOPSIS
       dropwatch [-l <method> | list]

DESCRIPTION
       dropwatch dropwatch is an interactive utility for monitoring and recording packets that are dropped by the kernel

OPTIONS
       -l <method> | list
              Select the translation method to use when a drop alert arrives.  By default the raw instruction pointer of a drop location is output, but by the use of the -l option, we can assign a translation method so that the instruction pointer can be translated
              into function names.  Currently supported lookup methods are:

       kas - use /proc/kallsyms to lookup instruction pointers to function mappings

INTERACTIVE COMMANDS
       start  Tells the kernel to start reporting dropped packets

       stop   Tells the kernel to discontinue reporting dropped packets

       exit   Exits the dropmonitor program

       help   Displays summary of all commands

       set alertlimit <value>
              Sets a triggerpoint to stop monitoring for dropped packets after <value> alerts have been received
```

- dropwatch 使用

```
[root@xg-ops-elk-lvs-1 ~]# dropwatch -l kas
Initalizing kallsyms db
dropwatch> start
Enabling monitoring...
Waiting for activation ack....
Waiting for activation ack....
Waiting for activation ack....
Waiting for activation ack....
Failed activation request, error: Resource temporarily unavailable
Shutting down ...
[root@xg-ops-elk-lvs-1 ~]#
```


> TO BE CONTINUED

------



# linux 系统 UDP 丢包问题分析思路

> http://cizixs.com/2018/01/13/linux-udp-packet-drop-debug
