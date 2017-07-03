# 关于 Recv-Q 和 Send-Q 的说明

## 引子

```
root@vagrant-ubuntu-trusty:~] $ uname -a
Linux vagrant-ubuntu-trusty 3.19.0-15-generic #15-Ubuntu SMP Thu Apr 16 23:32:37 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $ netstat -nat
Active Internet connections (servers and established)
Proto Recv-Q Send-Q Local Address           Foreign Address         State
tcp        0      0 127.0.0.1:27017         0.0.0.0:*               LISTEN
tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN
tcp        0      0 10.0.2.15:22            10.0.2.2:63158          ESTABLISHED
tcp        0      0 10.0.2.15:22            10.0.2.2:57560          ESTABLISHED
tcp6       0      0 :::111                  :::*                    LISTEN
tcp6       0      0 :::22                   :::*                    LISTEN
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $
root@vagrant-ubuntu-trusty:~] $ ss -nat
State       Recv-Q Send-Q       Local Address:Port                  Peer Address:Port
LISTEN      0      128              127.0.0.1:27017                            *:*
LISTEN      0      128                      *:111                              *:*
LISTEN      0      128                      *:22                               *:*
ESTAB       0      0                10.0.2.15:22                        10.0.2.2:63158
ESTAB       0      0                10.0.2.15:22                        10.0.2.2:57560
LISTEN      0      128                     :::111                             :::*
LISTEN      0      128                     :::22                              :::*
root@vagrant-ubuntu-trusty:~] $
```

从上面的输出可以看出：状态为 `LISTEN` 的 socket 其 `Send-Q` 的内容有所不同（其实 `Recv-Q` 的内容也可能不同）；

## 对比

在 `netstat` 输出中：

- `Recv-Q` means the count of bytes not copied by the user program connected to this socket. `Send-Q` means the count of bytes not acknowledged by the remote host. These should always be zero; if they're not you might have a problem. **Packets should not be piling up in either queue.** A brief queuing of outgoing packets is normal behavior. If the receiving queue is consistently jamming up, you might be experiencing a denial-of-service attack. If the sending queue does not clear quickly, you might have an application that is sending them out too fast, or the receiver cannot accept them quickly enough.


在 `ss` 输出中：

- **LISTEN 状态**：`Recv-Q` 表示当前 listen backlog 队列中的连接数目（等待用户调用 `accept()` 获取的、已完成 3 次握手的 socket 连接数量），而 `Send-Q` 表示了 listen socket 最大能容纳的 backlog ，即 `min(backlog, somaxconn)` 值。
- 非 LISTEN 状态：`Recv-Q` 表示了 receive queue 中存在的字节数目；`Send-Q` 表示 send queue 中存在的字节数；


## 分析

### netstat

通过如下命令可以确定 `Recv-Q`/`Send-Q` 的值是从 `/proc/net/tcp` 文件中读取出来的；

```
root@vagrant-ubuntu-trusty:~] $ strace -s 128 netstat -nat
...
write(1, "Active Internet connections (servers and established)\n", 54) = 54
write(1, "Proto Recv-Q Send-Q Local Address           Foreign Address         State      \n", 80) = 80
open("/proc/net/tcp", O_RDONLY)         = 3
read(3, "  sl  local_address rem_address   st tx_queue rx_queue tr tm->when retrnsmt   uid  timeout inode                                "..., 4096) = 900
write(1, "tcp        0      0 127.0.0.1:27017         0.0.0.0:*               LISTEN     \n", 80) = 80
write(1, "tcp        0      0 0.0.0.0:111             0.0.0.0:*               LISTEN     \n", 80) = 80
write(1, "tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN     \n", 80) = 80
write(1, "tcp        0      0 10.0.2.15:22            10.0.2.2:63158          ESTABLISHED\n", 80) = 80
write(1, "tcp        0      0 10.0.2.15:22            10.0.2.2:57560          ESTABLISHED\n", 80) = 80
read(3, "", 4096)                       = 0
close(3)                                = 0
open("/proc/net/tcp6", O_RDONLY)        = 3
read(3, "  sl  local_address                         remote_address                        st tx_queue rx_queue tr tm->when retrnsmt   ui"..., 4096) = 499
write(1, "tcp6       0      0 :::111                  :::*                    LISTEN     \n", 80) = 80
write(1, "tcp6       0      0 :::22                   :::*                    LISTEN     \n", 80) = 80
read(3, "", 4096)                       = 0
close(3)                                = 0
exit_group(0)                           = ?
+++ exited with 0 +++
root@vagrant-ubuntu-trusty:~] $
```

### ss

通过如下命令可以确定 `Recv-Q`/`Send-Q` 的值是基于 `PF_NETLINK` 类型 socket 读取到的；

```
root@vagrant-ubuntu-trusty:~] $ strace -s 128 ss -nat
...
write(1, "State      Recv-Q Send-Q                                                                                   Local Address:Port   "..., 229) = 229
socket(PF_NETLINK, SOCK_RAW, 4)         = 3
sendmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"H\0\0\0\24\0\1\3@                                                                                                                              \342\1\0\0\0\0\0\2\6\0\0\377\17\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 72}], msg_controllen=0, msg_flags=0}, 0) = 72
recvmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"`\0\0\0\24\0\2\0@                                                                                                                              \342\1\0p3\0\0\2\n\0\0i\211\0\0\177\0\0\1\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\300\225\224\32\0\210\377\377\0\0\0\0\0\0\0\0\200\0\0\0n\0\0\0\2643\0\0\5\0\10\0\0\0\0\0`\0\0\0\24\0\2\0@                    \342\1\0p3\0\0\2\n\0\0\0o\0\0\0\0\0\0\0\0\0\0"..., 8192}], msg_controllen=0, msg_flags=0}, 0) = 480
write(1, "LISTEN     0      128                                                                                          127.0.0.1:27017  "..., 230) = 230
write(1, "LISTEN     0      128                                                                                                  *:111    "..., 230) = 230
write(1, "LISTEN     0      128                                                                                                  *:22     "..., 230) = 230
write(1, "ESTAB      0      0                                                                                            10.0.2.15:22     "..., 230) = 230
write(1, "ESTAB      0      36                                                                                           10.0.2.15:22     "..., 230) = 230
recvmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"\24\0\0\0\3\0\2\0@\342\1\0p3\0\0\0\0\0\0", 8192}], msg_controllen=0, msg_flags=0}, 0) = 20
sendmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"H\0\0\0\24\0\1\3@                                                                                                                              \342\1\0\0\0\0\0\n\6\0\0\377\17\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0", 72}], msg_controllen=0, msg_flags=0}, 0) = 72
recvmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"`\0\0\0\24\0\2\0@\342\1\0p3\0\0\n\n\0\0\0o\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0\0@                      \225\32\0\210\377\377\0\0\0\0\0\0\0\0\200\0\0\0\0\0\0\0006,\0\0\5\0\10\0\0\0\0\0`\0\0\0\24\0\2\0@\342\1\0p3\0\0\n\n\0\0\0\26\0\0\0\0\0\0\0\0\0\0"..., 8192}], msg_controllen=0, msg_flags=0}, 0) = 192
write(1, "LISTEN     0      128                                                                                                 :::111    "..., 230) = 230
write(1, "LISTEN     0      128                                                                                                 :::22     "..., 230) = 230
recvmsg(3, {msg_name(12)={sa_family=AF_NETLINK, pid=0, groups=00000000}, msg_iov(1)=[{"\24\0\0\0\3\0\2\0@\342\1\0p3\0\0\0\0\0\0", 8192}], msg_controllen=0, msg_flags=0}, 0) = 20
close(3)                                = 0
exit_group(0)                           = ?
+++ exited with 0 +++
```


源码地址：[linux/net/ipv4/tcp_diag.c](http://elixir.free-electrons.com/linux/v3.19/source/net/ipv4/tcp_diag.c)

```
static void tcp_diag_get_info(struct sock *sk, struct inet_diag_msg *r, void *_info)
{
	const struct tcp_sock *tp = tcp_sk(sk);
	struct tcp_info *info = _info;

	if (sk->sk_state == TCP_LISTEN) {
		r->idiag_rqueue = sk->sk_ack_backlog;
		r->idiag_wqueue = sk->sk_max_ack_backlog;
	} else {
		// 等待接收的下一个 tcp 段的序号 - 尚未从内核空间 copy 到用户空间的段最前面的一个序号
		r->idiag_rqueue = max_t(int, tp->rcv_nxt - tp->copied_seq, 0);
		// 已加入发送队列中的 tcp 段的最后一个序号 - 已发送但尚未确认的最早一个序号
		r->idiag_wqueue = tp->write_seq - tp->snd_una;
	}
	if (info != NULL)
		tcp_get_info(sk, info);
}
```


----------


参考：

- [Keep an Eye on Your Linux Systems with Netstat
](http://www.enterprisenetworkingplanet.com/netos/article.php/3430561/Keep-an-Eye-on-Your-Linux-Systems-with-Netstat.htm)
- [nc, netstat, ss tools in Linux](https://madalanarayana.wordpress.com/2013/08/07/nm-netstat-ss-tools-in-linux/)
- [TCP queue 的一些问题](http://jaseywang.me/2014/07/20/tcp-queue-%E7%9A%84%E4%B8%80%E4%BA%9B%E9%97%AE%E9%A2%98/?spm=5176.100239.blogcont79972.33.ipKT4W)
- [Linux下查看socket状态 netstat升级版ss](http://blog.chinaunix.net/uid-20662820-id-3509532.html)