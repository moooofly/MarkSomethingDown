

# Disk Alarms

RabbitMQ will block producers when free disk space drops below a certain limit. This is a good idea since even transient messages can be paged to disk at any time, and running out of disk space can cause the server to crash. By default RabbitMQ will block producers, and prevent memory-based messages from being paged to disk, when free disk space drops below 50MB. This will reduce but not eliminate the likelihood of a crash due to disk space being exhausted. In particular, if messages are being paged out rapidly it is possible to run out of disk space and crash in the time between two runs of the disk space monitor. A more conservative approach would therefore be to set the limit to the same as the amount of memory installed on the system (see the configuration below).

Global flow control will be triggered if the amount of free disk space drops below a configured limit. The free space of the drive or partition that the broker database uses will be monitored at least every 10 seconds to determine whether the alarm should be raised or cleared. Monitoring will start as soon as the broker starts up, causing an entry in the broker logfile:

```shell
=INFO REPORT==== 23-Jun-2012::14:52:41 ===
Disk free limit set to 953MB
```

Monitoring will be disabled on unrecognised platforms, causing an entry such as the one below:

```shell
=WARNING REPORT==== 23-Jun-2012::15:45:29 ===
Disabling disk free space monitoring
```

When running RabbitMQ in a cluster, the disk alarm is cluster-wide; if one node goes under the limit then all nodes will block connections.

RabbitMQ will periodically check the amount of free disk space. The frequency with which disk space is checked is related to the amount of space at the last check (in order to ensure that the disk alarm goes off in a timely manner when space is exhausted). Normally disk space is checked every 10 seconds, but as the limit is approached the frequency increases. When very near the limit RabbitMQ will check as frequently as 10 times per second. This may have some effect on system load.

# Configuring the Disk Free Space Limit

The disk free space limit is configured with the disk_free_limit setting. By default 50MB is required to be free on the database partition. (See the description of file locations for the default location). This configuration file sets the disk free space limit to 1GB:

[{rabbit, [{disk_free_limit, 1000000000}]}].
Or you can use memory units (kB, kiB, MB, MiB, GB, GiB etc.) like this:
[{rabbit, [{disk_free_limit, "1GB"}]}].
It is also possible to set a free space limit relative to the RAM in the machine. This configuration file sets the disk free space limit to the same as the amount of RAM on the machine:
[{rabbit, [{disk_free_limit, {mem_relative, 1.0}}]}].
The limit can be changed while the broker is running using the rabbitmqctl set_disk_free_limit disk_limit command or rabbitmqctl set_disk_free_limit mem_relative fraction command. This command will take effect until the broker shuts down. The corresponding configuration setting should also be changed when the effects should survive a broker restart.







