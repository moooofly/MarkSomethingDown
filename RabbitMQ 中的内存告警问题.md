


# Memory Alarms

在启动并且执行过 `rabbitmqctlset_vm_memory_high_watermark fraction` 命令后，RabbitMQ server 会探测计算机上安装的 RAM 总量；默认情况下，当 RabbitMQ server 使用了超过 40% 的 RAM 内存时，会触发 memory alarm 并阻塞住所有 connection ；一旦 memory alarm 被清除（例如，由于 server 将消息 page out 到磁盘时，或者将消息 delivery 到客户端时）常规服务能力就能够恢复了；

默认的内存阈值被设置为已安装 RAM 的 40% ；注意，该值并不会真正阻止 RabbitMQ server 使用超过 40% 的内存，而只是一个会令 publisher 开始受限的点；在最坏的情况下，Erlang 的垃圾回收器能够导致内存使用量的 double（默认情况下为 RAM 的 80%）；因此，强烈建议开启 OS 自身的swap 或 page 文件功能；

在 32-bit 体系架构中，每个进程可用内存的限制为 2GB ；而常规实现的 64-bit 体系架构（即AMD64 和Intel EM64T）仅允许每个进程使用区区 256TB ；64-bit Windows 更进一步将其限制为8TB ；另外需要注意的是，就算在 64-bit 操作系统中，一个 32-bit 进程通常也只会使用最大 2GB 的地址空间；

## Configuring the Memory Threshold


可以通过配置文件调整能够触发流控机制的内存阈值；下面的示例中将阈值设置成默认的 0.4 ：
```shell
[{rabbit, [{vm_memory_high_watermark, 0.4}]}].
```

默认值 0.4 代表了 40% 的已安装 RAM 或者 40% 的可用虚拟地址空间（比前者更小）；例如在一个32-bit 平台上，如果你安装了 4GB 的 RAM ，那么 40% 的 4GB 为 1.6GB ，但是在 32-bit 的 Windows 上，通常会限制进程只能使用 2GB ，因此这里实际的阈值为 2GB 的 40% ，即 820MB ；

另一种方案为，直接设置节点可用的内存阈值为一个具体数值；下面的例子中设置阈值为 1073741824 字节（1024 MB） ：
```shell
[{rabbit, [{vm_memory_high_watermark, {absolute, 1073741824}}]}].
```

相同的例子，但使用了内存自己的单位：
```shell
[{rabbit, [{vm_memory_high_watermark, {absolute, "1024MiB"}}]}].
```

如果上面设置的绝对数值超过了实际安装的 RAM 大小，或者可用的虚拟地址空间大小，阈值会被自动调整为两者中较小的那个值；

在 RabbitMQ server 启动时，内存使用限制信息会输出到 RABBITMQ_NODENAME.log 文件中：
```shell
=INFO REPORT==== 29-Oct-2009::15:43:27 ===
Memory limit set to 2048MB.
```

The memory limit may also be queried using the rabbitmqctl status command.
The threshold can be changed while the broker is running using the rabbitmqctl set_vm_memory_high_watermark fraction command or rabbitmqctl set_vm_memory_high_watermark absolute memory_limit command. Memory units can also be used in this command. This command will take effect until the broker shuts down. The corresponding configuration setting should also be changed when the effects should survive a broker restart. The memory limit may change on systems with hot-swappable RAM when this command is executed without altering the threshold, due to the fact that the total amount of system RAM is queried.

Disabling all publishing
A value of 0 makes the memory alarm go off immediately and thus disables all publishing (this may be useful if you wish to disable publishing globally); use rabbitmqctl set_vm_memory_high_watermark 0.

Limited Address Space

When running RabbitMQ inside a 32 bit Erlang VM in a 64 bit OS (or a 32 bit OS with PAE), the addressable memory is limited. The server will detect this and log a message like:

=WARNING REPORT==== 19-Dec-2013::11:27:13 ===
Only 2048MB of 12037MB memory usable due to limited address space.
Crashes due to memory exhaustion are possible - see
http://www.rabbitmq.com/memory.html#address-space
The memory alarm system is not perfect. While stopping publishing will usually prevent any further memory from being used, it is quite possible for other things to continue to increase memory use. Normally when this happens and the physical memory is exhausted the OS will start to swap. But when running with a limited address space, running over the limit will cause the VM to crash.

It is therefore strongly recommended that when running on a 64 bit OS you use a 64 bit Erlang VM.

Configuring the Paging Threshold

Before the broker hits the high watermark and blocks publishers, it will attempt to free up memory by instructing queues to page their contents out to disc. Both persistent and transient messages will be paged out (the persistent messages will already be on disc but will be evicted from memory).

By default this starts to happen when the broker is 50% of the way to the high watermark (i.e. with a default high watermark of 0.4, this is when 20% of memory is used). To change this value, modify the vm_memory_high_watermark_paging_ratio configuration from its default value of 0.5. For example:

[{rabbit, [{vm_memory_high_watermark_paging_ratio, 0.75},
         {vm_memory_high_watermark, 0.4}]}].
The above configuration starts paging at 30% of memory used, and blocks publishers at 40%.

It is possible to set vm_memory_high_watermark_paging_ratio to a greater value than 1.0. In this case queues will not page their contents to disc. If this causes the memory alarm to go off, then producers will be blocked as explained above.

Unrecognised platforms

If the RabbitMQ server is unable to recognise your system, it will append a warning to the RABBITMQ_NODENAME.log file. It then assumes than 1GB of RAM is installed:

=WARNING REPORT==== 29-Oct-2009::17:23:44 ===
Unknown total memory size for your OS {unix,magic_homebrew_os}. Assuming memory size is 1024MB.
In this case, the vm_memory_high_watermark configuration value is used to scale the assumed 1GB RAM. With the default value of vm_memory_high_watermark set to 0.4, RabbitMQ's memory threshold is set to 410MB, thus it will throttle producers whenever RabbitMQ is using more than 410MB memory. Thus when RabbitMQ can't recognize your platform, if you actually have 8GB RAM installed and you want RabbitMQ to throttle producers when the server is using above 3GB, set vm_memory_high_watermark to 3.

For guidelines on recommended RAM watermark settings, see Production Checklist.