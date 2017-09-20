# Socket Buffer 相关 

## [Learning on socket buffer - 1](https://madalanarayana.wordpress.com/2014/12/08/learning-on-socket-buffer-1/)

In the previous post we have seen how the **stack size limit** can affect program execution. In this post I would like to explain the issue of **socket buffer sizes** with a real world problem I faced.

During one of my previous projects, we were using `syslogd` daemon for logging, all the log messages from different processes are collected on a remote `syslogd` server. Things were running smoothly until QA started load testing; Under heavy load, QA observed that some of the log messages are missing.

I decided to debug the issue, my initial thought was to doubt the user programs, as **there was no pattern in data loss**, I started to doubt the UDP data delivery. I used `netstat` to confirm that this is indeed due packet loss or errors at UDP. Let us see this with an example.

### How to be sure of packet loss ?

To better understand the issue of packet loss, I created a similar scenario faced by me using `iperf` tool. I ran the server program on my PI and the client on windows running on i7 core.

```
root@MPI1:/home/pi# iperf -u -s 192.168.0.9
iperf: ignoring extra argument -- 192.168.0.9
------------------------------------------------------------
Server listening on UDP port 5001
Receiving 1470 byte datagrams
UDP buffer size: 10.00 KByte (default)
------------------------------------------------------------
```

The client having powerful CPU, was sending packets at a higher rate than the ability of the server to receive and process. Observe the option “`-P 10`”, this starts 10 parallel threads, which will ensure that there is a real burst of packets on network.

```
c:\Perl\iperf-2.0.5-3-win32>iperf -u -c 192.168.0.9 -P 10
------------------------------------------------------------
Client connecting to 192.168.0.9, UDP port 5001
Sending 1470 byte datagrams
UDP buffer size: 63.0 KByte (default)
------------------------------------------------------------
[  3] local 192.168.0.4 port 54350 connected with 192.168.0.9 port 5001
......................................
```

Running `netstat` command on PI, tells us that is packet losses, “`-su`” option provides UDP statistics. These `packet receive errors`, could be because of **buffer overflows** or **wrong UDP check-sum**. In my case as the nodes are next to each other and no other check-sum errors are seen, we can deduce that the errors are due to buffer overflow.

```
root@MPI1:/home/pi# netstat -su
IcmpMsg:
    InType3: 18
    InType8: 79
    OutType0: 79
    OutType3: 59
    OutType11: 68
Udp:
    46686 packets received
    917 packets to unknown port received.
    1216 packet receive errors      -- here
    2112 packets sent
UdpLite:
```

Ok now we know there are errors in packet reception, Let us try to increase the receive buffer size and see how this will affect data reception.

### How to increase UDP socket buffer sizes ?

Socket buffer size can start with minimal value, and can grow till the maximum allowed size based on system load and options set using `setsockopt`, Linux is very good at dynamically changing the buffer sizes from minimum to maximum allowed.

```
root@MPI1:/home/pi# sysctl -a |grep mem
net.core.optmem_max = 10240
net.core.rmem_default = 163840
net.core.rmem_max = 163840
net.core.wmem_default = 163840
net.core.wmem_max = 163840
net.ipv4.udp_mem = 10491        13991   20982
net.ipv4.udp_rmem_min = 4096
net.ipv4.udp_wmem_min = 4096
root@MPI1:/home/pi#
```

The upper limits of socket buffer are controlled by the parameters **rmem_default**, **rmem_max**. These values can be modified using `sysctl` command.

```
root@MPI1:/home/pi# sysctl -w net.core.rmem_default=163840
root@MPI1:/home/pi# sysctl -w net.core.rmem_max=163840,
```

To preserve the changes across reboots, don’t forget to change the values in file `/etc/sysctl.conf`; Otherwise the values changed using `sysctl` command are lost upon reboot.

### How to confirm that the values are changed ?

In user program one can confirm that the values are changed using `setsockopt` function.
In my case I confirmed by running the program `iperf` again.

```
iperf: ignoring extra argument -- 192.168.0.9
------------------------------------------------------------
Server listening on UDP port 5001
Receiving 1470 byte datagrams
UDP buffer size:  160 KByte (default)
------------------------------------------------------------
```

The value of 160k confirms that the UDP buffer is modified on the system.

### Any other option to handle packet loss ?

The basic issue with UDP is the unreliable nature of packet delivery, By changing the communication protocol from UDP to TCP one can hope to eliminate the issue of packet loss. TCP provides flow control so one can be assured that, if the receive window is full the sender will be slowed down which should eliminate the buffer overflows.

```
# provides UDP syslog reception
#$ModLoad imudp
#$UDPServerRun 514

# provides TCP syslog reception
#$ModLoad imtcp
#$InputTCPServerRun 514
```

syslog can provide both UDP and TCP based data transfers, hence we should bee safe in using TCP.

In the next post I will explain the problems 😦 I faced even after I changed the communication protocol.


## [Learning on socket buffer - 2](https://madalanarayana.wordpress.com/2014/12/08/learning-on-socket-buffer-2/)

In the previous post we had seen the issue of UDP **buffer overflow**, and tried to solve the problem using TCP. In this post I will explain the problems faced while changing the communication protocol from UDP to TCP.

**Flow control** is the root cause of this problem, We can approach to solve in multiple ways; Till now we tried to fix the issue by changing the size of reception buffer and changing the protocol used for communication.

### Do you see any issues with TCP based transfers ?

Use of TCP should solve the flow control problem, **If the sender code is not written properly this change of protocol can some time move the problem from receiver to sender**. Let us see this with an example.

Most of the real world programs creates a non-blocking socket, By creating a non-blocking socket the program doesn’t have to wait while calling `send`/`recv` on socket. Using a non-blocking socket with `select`/`epoll` system call is the best form of using sockets in any real world application, Some times this can result in unforeseen situations if the code is not written properly.

When a socket is created and moved to non-blocking mode, Use of “send” on this socket can result in errors under high load situations. The man page of `send` states, for a TCP socket

> "When the message does not fit into the send buffer of the socket, `send()` normally blocks, unless the socket has been placed in non-blocking I/O mode.  In non-blocking mode it would fail with the error EAGAIN or EWOULDBLOCK in this case. The `select(2)` call may be used to determine when it is possible to send more data."

Here in lies the problem, If the application treats it as an communication error, this can lead to other kind of problems.

Use of TCP can result in error situations if the data overflows send buffer.

### How to handle send buffer overflows in TCP?

We can mitigate the send buffer overflows by increasing the buffer size of sender. Under Linux TCP maintains buffers of three sizes.

```
pi@MPI1:~sysctl -a |grep mem |grep tcp
net.ipv4.tcp_mem = 10398        13866   20796
net.ipv4.tcp_rmem = 4096        87380   3581856
net.ipv4.tcp_wmem = 4096        16384   3581856
```

The values shown above are on my PI, **tcp_rmem** denotes the buffer sizes of receiver and **tcp_wmem** denotes the buffer sizes of sender. One might have noticed that the problem can be fixed by properly tuning the buffer sizes at both sender and receiver, By modifying the values at both nodes we are allowing for packets to stay in buffer little longer.

When I started `iperf` tool on my PI, the size of TCP receive buffer is fixed at median value.

```
pi@MPI1:~$ iperf -s 192.168.0.9
iperf: ignoring extra argument -- 192.168.0.9
------------------------------------------------------------
Server listening on TCP port 5001
TCP window size: 85.3 KByte (default)
```

### How to change the buffer sizes ?

The sizes of send and receive buffers can be modified using sysctl command.

```
root@MPI1:/home/pi# sysctl -w net.ipv4.tcp_rmem='4096    873800   3581856'
net.ipv4.tcp_rmem = 4096        873800   3581856
```

I tried to modify the buffer sizes using command `sysctl`.  Though I incremented the TCP buffer to 870K this didn’t take affect, This is because overall system limits are not modified. The overall system limit is fixed at 160k.

```
root@MPI1:/home/pi# iperf -u -s 192.168.0.9
iperf: ignoring extra argument -- 192.168.0.9
------------------------------------------------------------
Server listening on UDP port 5001
Receiving 1470 byte datagrams
UDP buffer size:  160 KByte (default)
------------------------------------------------------------
```

**Operating system level buffer sizes** are controlled by the following variables

- `net.core.rmem_max` – This sets the **maximum** receive buffer size for all types of connections which includes both TCP and UDP.
- `net.core.wmem_max` – This sets the **maximum** send buffer size for all types of connections which includes both TCP and UDP .
- `net.core.rmem_default` – This sets the **default** receive buffer size for all types of connections.
- `net.core.wmem_default` – This sets the **default** send buffer size for all types of connections.

Don’t forget to change the values in file `/etc/sysctl.conf`, Changes here will preserves the values across system reboots.

```
root@MPI1:/home/pi# sysctl -a |grep mem
net.core.rmem_default = 163840
net.core.rmem_max = 163840
net.core.wmem_default = 163840
```

- `net.ipv4.tcp_rmem` – This sets **minimum** size, **Default** size and the **maximum** size of receive buffers allocated for TCP.
- `net.ipv4.tcp_wmem` – This sets **minimum** size, **Default** size and the **maximum** size of receive buffers allocated for TCP.

### Any other solution to TCP send buffer overflow?

Increasing the buffer sizes can help us tide over the problem for some time, In case of TCP if you are using blocking sockets Linux will take care of flow control for you. In the case of non-blocking sockets it is always better to check available buffer size before sending the packet.

Available buffer size for a socket can be found using `ioctl` function, “man 7 tcp” provides a detail list of problem and solutions. Please read this man page for further solutions.

Before writing data to a non-blocking TCP socket,

- Check for available buffer size using `ioctl(.., SIOCOUTQ, &val)`
- Only send data of this size, and wait till the buffer is free by using select


----------


小结：这两篇文章中，关于 UDP 和 TCP 的 buffer size 的调节可以借鉴，关于 `iperf` 和 `iptraf` 工具的使用可以借鉴；但关于 non-blocking TCP socket 的 buffer 处理问题，持保留意见；
