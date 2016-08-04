


# Memory Alarms

在启动并且执行过 `rabbitmqctlset_vm_memory_high_watermark fraction` 命令后，RabbitMQ server 会探测计算机上安装的 RAM 总量；默认情况下，当 RabbitMQ server 使用了超过 40% 的 RAM 内存时，会触发 memory alarm 并阻塞住所有 connection ；一旦 memory alarm 被清除（例如，由于 server 将消息 page out 到磁盘时，或者将消息 delivery 到客户端时）常规服务能力就能够恢复了；

默认的内存阈值被设置为已安装 RAM 的 40% ；注意，该值并不会真正阻止 RabbitMQ server 使用超过 40% 的内存，而只是一个会令 publisher 开始受限的点；在最坏的情况下，Erlang 的垃圾回收器能够导致内存使用量的 double（默认情况下为 RAM 的 80%）；因此，强烈建议开启 OS 自身的swap 或 page 文件功能；

在 32-bit 体系架构中，每个进程可用内存的限制为 2GB ；而常规实现的 64-bit 体系架构（即AMD64 和Intel EM64T）仅允许每个进程使用区区 256TB ；64-bit Windows 更进一步将其限制为8TB ；另外需要注意的是，就算在 64-bit 操作系统中，一个 32-bit 进程通常也只会使用最大 2GB 的地址空间；

## Configuring the Memory Threshold

The memory threshold at which the flow control is triggered can be adjusted by editing the configuration file. The example below sets the threshold to the default value of 0.4:

[{rabbit, [{vm_memory_high_watermark, 0.4}]}].
The default value of 0.4 stands for 40% of installed RAM or 40% of available virtual address space, whichever is smaller. E.g. on a 32-bit platform, if you have 4GB of RAM installed, 40% of 4GB is 1.6GB, but 32-bit Windows normally limits processes to 2GB, so the threshold is actually to 40% of 2GB (which is 820MB).

Alternatively, the memory threshold can be adjusted by setting an absolute limit of RAM used by the node. The example below sets the threshold to 1073741824 bytes (1024 MB):

[{rabbit, [{vm_memory_high_watermark, {absolute, 1073741824}}]}].
Same example, but using memory units:
[{rabbit, [{vm_memory_high_watermark, {absolute, "1024MiB"}}]}].
If the absolute limit is larger than the installed RAM or available virtual address space, the threshold is set to whichever limit is smaller.

The memory limit is appended to the RABBITMQ_NODENAME.log file when the RabbitMQ server starts:

=INFO REPORT==== 29-Oct-2009::15:43:27 ===
Memory limit set to 2048MB.
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