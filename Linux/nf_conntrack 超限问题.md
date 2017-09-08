# nf_conntrack 超限问题

> 背景：基于 C/S 模型进行大并发连接测试；

## 问题表现

当 server 侧连接数达到 262138 后，无法再建立新的连接，此时通过 dmesg 或者 `/var/log/kern.log` 可以看到

```
[35470.975678] nf_conntrack: table full, dropping packet
[35704.617801] nf_conntrack: table full, dropping packet
[35706.621871] nf_conntrack: table full, dropping packet
[35710.633918] nf_conntrack: table full, dropping packet
[35718.649976] nf_conntrack: table full, dropping packet
```

在 client 侧也是一样的现象；

发生问题后，通过其他机器 ping 该机器会时通时不通，使用 ssh 登录的情况也一样；


## 故障原因

**用来跟踪连接信息的哈希表满了**；

在 Ubuntu 16.04.3 LTS (Linux 4.4.0-87-generic x86_64) 上，未进行配置调整前，`nf_conntrack` 的默认值为（对应之前建立了 26w+ 连接后报错的情况）：

```
$ cat /proc/sys/net/netfilter/nf_conntrack_max
262144
```

可以通过如下命令查看 `nf_conntrack` 的当前使用状况

```
$ cat /proc/sys/net/netfilter/nf_conntrack_count
```

背景知识：

- `nf_conntrack`/`ip_conntrack` 跟 NAT 有关，用来**跟踪连接条目**，它会使用一个**哈希表**来记录 established 的记录。`nf_conntrack` 在 2.6.15 被引入，而 `ip_conntrack` 在 2.6.22 被移除；如果该哈希表满了，就会出现
"nf_conntrack: table full, dropping packet" 信息；
- `nf_conntrack` 工作在 3 层，支持 IPv4 和 IPv6，而 `ip_conntrack` 只支持 IPv4。目前，大多的 ip_conntrack_* 已被 nf_conntrack_* 取代，很多 ip_conntrack_* 仅仅是个 alias ；


## 系统配置

```
$ vi /etc/sysctl.conf
...
net.nf_conntrack_max = 1048576
net.netfilter.nf_conntrack_max = 1048576
```


----------


## 其他

> 以下内容针对 CentOS6/CentOS7

- 内核参数 `net.nf_conntrack_max` 系统默认值为 "65536" ，当 `nf_conntrack` 模块被装置且服务器上连接超过这个设定的值时，系统会主动丢掉新连接包，直到连接小于此设置值才会恢复。同时内核参数 `net.netfilter.nf_conntrack_tcp_timeout_established` 系统默认值为 "432000" ，代表 `nf_conntrack` 中保存的 TCP 连接记录的默认时间是 5 天，（保存时间过长）致使 `nf_conntrack` 的值减不下来，丢包持续时间长；
- `nf_conntrack` 模块在**首次装载**或**重新装载**时，内核参数 `net.nf_conntrack_max` 会重新设置为默认值 "65536" ，并且不会调用 `sysctl` 设置为我们的预设值；
- 触发 `nf_conntrack` 模块**首次装载**比较隐蔽，任何调用 iptables NAT 功能的操作都会触发。当系统没有挂载 `nf_conntrack` 模块时（即 `lsmod |grep conntrack` 时无输出时），调用 `iptables` NAT 相关命令（`iptables -L -t nat`）就会触发 nf_conntrack 模块装置，致使 `net.nf_conntrack_max` 重设为 "65536" 。
- 触发 `nf_conntrack` 模块重新装载的操作很多，CentOS6 中 `service iptables restart` ，CentOS7 中 `systemctl restart firewalld` 都会触发设置重置，致使 `net.nf_conntrack_max` 重设为 "65536" 。

### 解决办法

- 通过系统初始化脚本创建配置文件 `/etc/modprobe.d/nf_conntrack.conf` ，写入内容为 `options nf_conntrack hashsize=262144` ；（这种方法）通过设置 `nf_conntrack` 模块的挂接参数 `hashsize` ，进而设置 net.nf_conntrack_max 的值为 "2097152"（因为 nf_conntrack_max=hashsize*8），保证后续新初始化服务器配置正确；
- 通过xxx方式将配置文件 `/etc/modprobe.d/nf_conntrack.conf` 推送到所有目标机器，内容为 `options nf_conntrack hashsize=262144` ，保证 `nf_conntrack` 模块在首次装载或重新装载时，`net.nf_conntrack_max` 内核参数被设置为预期的 "2097152" ；
- 更新系统初始化脚本，设置 `net.netfilter.nf_conntrack_tcp_timeout_established=1800` ，减少 `nf_conntrack` 连接表中 TCP 连接的记录维持时间；


### 测试命令

- `nf_conntrack` 模块在首次装载时初始化默认值为 "65536"

```
[root@localhost ~]# systemctl stop firewalld
[root@localhost ~]# lsmod |grep nf_conntrack
[root@localhost ~]# sysctl -a |grep nf_conntrack
[root@localhost ~]# iptables -L -t nat
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination
 
Chain INPUT (policy ACCEPT)
target     prot opt source               destination
 
Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination
 
Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination
[root@localhost ~]# lsmod |grep nf_conntrack
nf_conntrack_ipv4      19108  1
nf_defrag_ipv4         12729  1 nf_conntrack_ipv4
nf_conntrack          111302  3 nf_nat,nf_nat_ipv4,nf_conntrack_ipv4
[root@localhost ~]# sysctl -a |grep net.nf_conntrack_max
net.nf_conntrack_max = 65536
```

- 设置 `net.nf_conntrack_max = 2097152` ，重启 firewalld 服务，`nf_conntrack_max` 重新被初始化为 "65536" ；

```
[root@localhost ~]# systemctl start firewalld
[root@localhost ~]# sysctl net.nf_conntrack_max=2097152
net.nf_conntrack_max = 2097152
[root@localhost ~]# sysctl -a |grep net.nf_conntrack_max
net.nf_conntrack_max = 2097152
 
[root@localhost ~]# systemctl restart firewalld
[root@localhost ~]# sysctl -a |grep net.nf_conntrack_max
net.nf_conntrack_max = 65536
```

- 设置配置文件 `/etc/modprobe.d/nf_conntrack.conf` ，内容为 `options nf_conntrack hashsize=262144` ，保证 `net.nf_conntrack_max =2097152` 

```
[root@localhost ~]# cat /etc/modprobe.d/nf_conntrack.conf
options nf_conntrack hashsize=262144
 
[root@localhost ~]# systemctl stop firewalld
 
[root@localhost ~]# iptables -L -t nat
Chain PREROUTING (policy ACCEPT)
target     prot opt source               destination
 
Chain INPUT (policy ACCEPT)
target     prot opt source               destination
 
Chain OUTPUT (policy ACCEPT)
target     prot opt source               destination
 
Chain POSTROUTING (policy ACCEPT)
target     prot opt source               destination
 
[root@localhost ~]# sysctl -a |grep net.nf_conntrack_max
net.nf_conntrack_max = 2097152
 
[root@localhost ~]# systemctl restart firewalld
[root@localhost ~]# sysctl -a |grep net.nf_conntrack_max
net.nf_conntrack_max = 2097152
```
