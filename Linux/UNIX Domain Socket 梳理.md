# UNIX Domain Socket 梳理

有同事问 “Unix Domain Socket 上能否进行 UDP 通信”，我答“应该能，因为其和普通 socket 的行为是一致的”；

然而，基于 Unix Domain Socket 进行 UDP 通信，我之前确实没有研究过，故整理相关资料如下；


----------


## UNIX Domain Socket IPC

socket API 原本是为网络通讯设计的，但后来在 socket 框架上发展出一种 IPC 机制，就是 **UNIX Domain Socket** 。虽然网络 socket 也可用于同一台主上机的进程间通讯（即通过 loopback 地址 127.0.0.1），但是 UNIX Domain Socket 用于 IPC 更有效率：**不需要经过网络协议栈，不需要打包拆包、计算校验和、维护序号和应答等，只是将应用层数据从一个进程拷贝到另一个进程**。这是因为 **IPC 机制本质上是可靠的通讯**，而网络协议是为不可靠的通讯设计的。UNIX Domain Socket 也提供面向流和面向数据包两种 API 接口，类似于 TCP 和 UDP ，但是**面向消息的 UNIX Domain Socket 也是可靠的，消息既不会丢失也不会顺序错乱**。

UNIX Domain Socket 是全双工的，API 接口语义丰富，相比其它 IPC 机制有明显的优越性，目前已成为使用最广泛的 IPC 机制，比如 X Window 服务器和 GUI 程序之间就是通过 UNIX Domain Socket 通讯的。

使用 UNIX Domain Socket 的过程和网络 socket 十分相似，也要先调用 `socket()` 创建一个 socket 文件描述符，address family 指定为 `AF_UNIX` ，type 可以选择 `SOCK_DGRAM` 或 `SOCK_STREAM` ，protocol 参数仍然指定为 0 。

UNIX Domain Socket 与网络 socket 编程最明显的不同在于地址格式不同（UNIX Domain Socket 用结构体 `sockaddr_un` 表示），网络编程的 socket 地址是 IP 地址加端口号，而 UNIX Domain Socket 的地址是一个 socket 类型的文件在文件系统中的路径，这个 socket 文件由 `bind()` 调用创建，如果调用 `bind()` 时该文件已存在，则 `bind()` 错误返回。

以下程序将 UNIX Domain socket 绑定到一个地址。

```c
#include <stdlib.h>
#include <stdio.h>
#include <stddef.h>
#include <sys/socket.h>
#include <sys/un.h>

int main(void)
{
	int fd, size;
	struct sockaddr_un un;

	memset(&un, 0, sizeof(un));
	un.sun_family = AF_UNIX;
	strcpy(un.sun_path, "foo.socket");
	if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0) {
		perror("socket error");
		exit(1);
	}
	size = offsetof(struct sockaddr_un, sun_path) + strlen(un.sun_path);
	if (bind(fd, (struct sockaddr *)&un, size) < 0) {
		perror("bind error");
		exit(1);
	}
	printf("UNIX domain socket bound\n");
	exit(0);
}
```

注意程序中的 `offsetof` 宏，它在 `stddef.h` 头文件中定义：

```c
#define offsetof(TYPE, MEMBER) ((int)&((TYPE *)0)->MEMBER)
```

`offsetof(struct sockaddr_un, sun_path)` 就是取 `sockaddr_un` 结构体的 `sun_path` 成员在结构体中的偏移，也就是从结构体的第几个字节开始是 `sun_path` 成员。想一想，这个宏是如何实现这一功能的？

该程序的运行结果如下

```
$ ./a.out
UNIX domain socket bound
$ ls -l foo.socket
srwxrwxr-x 1 user        0 Aug 22 12:43 foo.socket
$ ./a.out
bind error: Address already in use
$ rm foo.socket
$ ./a.out
UNIX domain socket bound
```

以下是服务器的 **listen 模块**，与网络 socket 编程类似，在 `bind` 之后要 `listen` ，表示通过 `bind` 的地址（也就是 socket 文件）提供服务。

```c
#include <stddef.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>

#define QLEN 10

/*
 * Create a server endpoint of a connection.
 * Returns fd if all OK, <0 on error.
 */
int serv_listen(const char *name)
{
	int                 fd, len, err, rval;
	struct sockaddr_un  un;

	/* create a UNIX domain stream socket */
	if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0)
		return(-1);
	unlink(name);   /* in case it already exists */

	/* fill in socket address structure */
	memset(&un, 0, sizeof(un));
	un.sun_family = AF_UNIX;
	strcpy(un.sun_path, name);
	len = offsetof(struct sockaddr_un, sun_path) + strlen(name);

	/* bind the name to the descriptor */
	if (bind(fd, (struct sockaddr *)&un, len) < 0) {
		rval = -2;
		goto errout;
	}
	if (listen(fd, QLEN) < 0) { /* tell kernel we're a server */
		rval = -3;
		goto errout;
	}
	return(fd);

errout:
	err = errno;
	close(fd);
	errno = err;
	return(rval);
}
```

以下是服务器的 **accept 模块**，通过 accept 得到的客户端地址也应该是一个 socket 文件，如果不是 socket 文件就返回错误码；**如果是 socket 文件，在建立连接后这个文件就没有用了，调用 `unlink` 把它删掉**，通过传出参数 uidptr 返回客户端程序的 user id 。

```c
#include <stddef.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>

int serv_accept(int listenfd, uid_t *uidptr)
{
	int                 clifd, len, err, rval;
	time_t              staletime;
	struct sockaddr_un  un;
	struct stat         statbuf;

	len = sizeof(un);
	if ((clifd = accept(listenfd, (struct sockaddr *)&un, &len)) < 0)
		return(-1);     /* often errno=EINTR, if signal caught */

	/* obtain the client's uid from its calling address */
	len -= offsetof(struct sockaddr_un, sun_path); /* len of pathname */
	un.sun_path[len] = 0;           /* null terminate */

	if (stat(un.sun_path, &statbuf) < 0) {
		rval = -2;
		goto errout;
	}

	if (S_ISSOCK(statbuf.st_mode) == 0) {
		rval = -3;      /* not a socket */
		goto errout;
	}

	if (uidptr != NULL)
		*uidptr = statbuf.st_uid;   /* return uid of caller */
	unlink(un.sun_path);        /* we're done with pathname now */
	return(clifd);

errout:
	err = errno;
	close(clifd);
	errno = err;
	return(rval);
}
```

以下是客户端的 connect 模块，与网络 socket 编程不同的是，**UNIX Domain Socket 客户端一般要显式调用 bind 函数，而不依赖系统自动分配的地址**。客户端 bind 一个自己指定的 socket 文件名的好处是，该文件名可以包含客户端的 pid 以便服务器区分不同的客户端。

```c
#include <stdio.h>
#include <stddef.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>

#define CLI_PATH    "/var/tmp/"      /* +5 for pid = 14 chars */

/*
 * Create a client endpoint and connect to a server.
 * Returns fd if all OK, <0 on error.
 */
int cli_conn(const char *name)
{
	int                fd, len, err, rval;
	struct sockaddr_un un;

	/* create a UNIX domain stream socket */
	if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0)
		return(-1);

	/* fill socket address structure with our address */
	memset(&un, 0, sizeof(un));
	un.sun_family = AF_UNIX;
	sprintf(un.sun_path, "%s%05d", CLI_PATH, getpid());
	len = offsetof(struct sockaddr_un, sun_path) + strlen(un.sun_path);

	unlink(un.sun_path);        /* in case it already exists */
	if (bind(fd, (struct sockaddr *)&un, len) < 0) {
		rval = -2;
		goto errout;
	}

	/* fill socket address structure with server's address */
	memset(&un, 0, sizeof(un));
	un.sun_family = AF_UNIX;
	strcpy(un.sun_path, name);
	len = offsetof(struct sockaddr_un, sun_path) + strlen(name);
	if (connect(fd, (struct sockaddr *)&un, len) < 0) {
		rval = -4;
		goto errout;
	}
	return(fd);

errout:
	err = errno;
	close(fd);
	errno = err;
	return(rval);
}
```

下面是自己动手时间，请利用以上模块编写完整的客户端/服务器通讯的程序。

client 代码

```c
#include <stdio.h>
#include <stddef.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>

#include <stdlib.h>

#define CLI_PATH    "/var/tmp/"      /* +5 for pid = 14 chars */
//#define CLI_PATH    "./"      /* +5 for pid = 14 chars */

/*
 * Create a client endpoint and connect to a server.
 * Returns fd if all OK, <0 on error.
 */
int cli_conn(const char *name)
{
    int                fd, len, err, rval;
    struct sockaddr_un un;

    /* create a UNIX domain stream socket */
    if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0)
        return(-1);

    /* fill socket address structure with our address */
    memset(&un, 0, sizeof(un));
    un.sun_family = AF_UNIX;
    sprintf(un.sun_path, "%s%05d", CLI_PATH, getpid());
    len = offsetof(struct sockaddr_un, sun_path) + strlen(un.sun_path);

    unlink(un.sun_path);        /* in case it already exists */
    if (bind(fd, (struct sockaddr *)&un, len) < 0) {
        rval = -2;
        goto errout;
    }

    /* fill socket address structure with server's address */
    memset(&un, 0, sizeof(un));
    un.sun_family = AF_UNIX;
    strcpy(un.sun_path, name);
    len = offsetof(struct sockaddr_un, sun_path) + strlen(name);
    if (connect(fd, (struct sockaddr *)&un, len) < 0) {
        rval = -4;
        goto errout;
    }
    return(fd);

errout:
    err = errno;
    close(fd);
    errno = err;
    return(rval);
}

int main(void)
{
    int res = -1;

    res = cli_conn("server.sock");
    printf("client: res = %d\n", res);
    printf("client: UNIX domain socket bound\n");
    exit(0);
}
```

server 代码

```c
#include <stddef.h>
#include <sys/stat.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <errno.h>

#include <stdio.h>
#include <stdlib.h>

#define QLEN 10

/*
 * Create a server endpoint of a connection.
 * Returns fd if all OK, <0 on error.
 */
int serv_listen(const char *name)
{
    int                 fd, len, err, rval;
    struct sockaddr_un  un;

    /* create a UNIX domain stream socket */
    if ((fd = socket(AF_UNIX, SOCK_STREAM, 0)) < 0)
        return(-1);
    unlink(name);   /* in case it already exists */

    /* fill in socket address structure */
    memset(&un, 0, sizeof(un));
    un.sun_family = AF_UNIX;
    strcpy(un.sun_path, name);
    len = offsetof(struct sockaddr_un, sun_path) + strlen(name);

    /* bind the name to the descriptor */
    if (bind(fd, (struct sockaddr *)&un, len) < 0) {
        rval = -2;
        goto errout;
    }
    if (listen(fd, QLEN) < 0) { /* tell kernel we're a server */
        rval = -3;
        goto errout;
    }
    return(fd);

errout:
    err = errno;
    close(fd);
    errno = err;
    return(rval);
}

int serv_accept(int listenfd, uid_t *uidptr)
{
    int                 clifd, len, err, rval;
    time_t              staletime;
    struct sockaddr_un  un;
    struct stat         statbuf;

    len = sizeof(un);
    if ((clifd = accept(listenfd, (struct sockaddr *)&un, &len)) < 0)
        return(-1);     /* often errno=EINTR, if signal caught */

    /* obtain the client's uid from its calling address */
    len -= offsetof(struct sockaddr_un, sun_path); /* len of pathname */
    un.sun_path[len] = 0;           /* null terminate */

    if (stat(un.sun_path, &statbuf) < 0) {
        rval = -2;
        goto errout;
    }

    if (S_ISSOCK(statbuf.st_mode) == 0) {
        rval = -3;      /* not a socket */
        goto errout;
    }

    if (uidptr != NULL)
        *uidptr = statbuf.st_uid;   /* return uid of caller */
    unlink(un.sun_path);        /* we're done with pathname now */
    return(clifd);

errout:
    err = errno;
    close(clifd);
    errno = err;
    return(rval);
}


int main(void)
{
    int fd = -1;
    int cfd = -1 ;
    uid_t uid;

    fd =  serv_listen("server.sock");
    cfd = serv_accept(fd, &uid);
    fprintf(stderr, "server: fd=%d, cfd=%d\n", fd, cfd);
    exit(0);
}
```


----------

## Unix Domain Socket 的一些小结

找了大半天的资料，收获也不多，其实还是自己思考更靠谱一些。

### Unix 域的数据报服务是否可靠

从 man 手册中可以看到，**Unix Domain Socket 的数据报既不会丢失也不会乱序**（据我所知，在 Linux 下的确是这样）。不过最新版本的内核，仍然又提供了一个保证次序的类型 "kernel 2.6.4 SOCK_SEQPACKET" 。

### STREAM 和 DGRAM 的主要区别

Unix Domain Socket 的 SOCK_DGRAM 既然不会丢失数据，那不是和 SOCK_STREAM 很类似么？我理解也确实是这样，而且我觉得 SOCK_DGRAM 相对还要更好一些，因为**发送的数据可以带边界**。二者另外的区别在于**收发时的数据量不一样**，基于 SOCK_STREAM 的套接字，`send` 时可以传入超过 SO_SNDBUF 长的数据，`recv` 时，同普通 TCP 类似，会存在数据粘连。

**采用阻塞方式使用 API** ；在 Unix Domain Socket 下调用 `sendto` 时，如果缓冲队列已满，会阻塞。而普通 UDP 因为不是可靠的，无法感知对端的情况，即使对端没有及时收取数据，基本上 sendto 都能立即返回成功（如果发端疯狂 `sendto` 就另当别论，因为过快地调用 `sendto` 在慢速网络的环境下，可能撑爆套接字的缓冲区，导致 `sendto` 阻塞）。

### SO_SNDBUF 和 SO_REVBUF

**对于 Unix Domain Socket 来说，设置 SO_SNDBUF 会影响 `sendto` 最大的报文长度，但是任何针对 SO_RCVBUF 的设置都是无效的**。实际上 Unix Domain Socket 的数据报（SOCK_DGRAM）处理过程还是得将数据放入内核所申请的内存块里面，再由另一个进程通过 `recvfrom` 从内核读取，因此具体**可以发送的数据报长度会受限于内核的 `slab` 策略**。在 linux 平台下，早先版本（如 2.6.2）可发送最大数据报长度约为 128 k ，新版本的内核支持更大的长度。

### 使用 SOCK_DGRAM 时，缓冲队列的长度

有几个因素会影响缓冲队列的长度，一个是上面提到的 `slab` 策略，另一个则是系统的内核参数 `/proc/sys/net/unix/max_dgram_qlen` 。缓冲队列长度是这二者共同决定的。

如 `max_dgram_qlen` 默认为 10，在数据报较小时（如 1k），先挂起接收数据的进程后，仍可以成功执行 `sendto` 十次并顺利返回；但是如果数据报较大（如 120k）时，就要看 slab 的 “size-131072” 的 limit 了。

### 使用 Unix Domain Socket 进行进程间通信 v.s. 其他方式

- 需要先**确定操作系统类型**，以及其所对应的**最大数据报长度**；如果有需要传送超过该长度的数据报，建议拆分成几个发送，接收后组装即可（不会乱序，个人觉得这样做比用 STREAM 再切包方便得多）；
- 同**管道**相比，Unix 域的数据报不但可以维持数据的边界，还不会碰到在写入管道时的原子性问题；
- 同**共享内存**相比，不能实现独立于进程的、缓存大量数据的功能，但是却避免了同步互斥的考量；
- 同普通 socket 相比，开销相对较小（不用计算报头），Unix Domain Socket 数据报的报文长度可以大于 64k，不过不能像普通 socket 那样将进程切换到不同机器 。

### 其他

其实在本机 IPC 时，同普通 socket 的 UDP 相比，Unix Domain Socket 的数据报只不过是在收发时分别少计算了一下校验和而已，**本机的 UDP 会走 loopback 接口**，**不会进行 IP 分片**，也**不会占用网卡硬件**（不会真正跑到网卡链路层上去）。也就是说，在本机上使用普通的 socket UDP，只是多耗了一些 CPU（之所以说一些，是因为校验和的计算很简单），此外本机的 UDP 也可以保证数据不丢失、不乱序 。

从我个人的经验来看，即便是高并发的网络服务器，单纯因为收发包造成的 CPU 占用其实并不算多（其中收包占用的 CPU 从 `%si` 可见一斑，因为收包需通过软中断实现的），倒是网卡带宽、磁盘IO、后台逻辑、内存使用等问题往往成为主要矛盾。

所以，在没有长时缓存通信数据的需求时，可以考虑通过 UDP 来实现本地进程间 IPC，这样也便于切换机器。对于较长的报文，可以切分成小的，再重新组装，不过这样做仅适用于本机通信，如果还要考虑以后迁移机器，那还是老老实实地 TCP 吧。


----------


## Datagrams in the UNIX domain

Unlike the previous examples, which dealt with streams sockets, the following two sample programs send and receive data on datagram sockets. These examples are for the **UNIX domain**; for the Internet domain equivalents, see "[Datagrams in the Internet Domain](http://osr507doc.sco.com/en/netguide/disockT.datagram_codesamples.html)".

First, create a server that can receive UNIX domain datagrams:

### Reading UNIX domain datagrams

```c
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>

/*
 * In the included file <sys/un.h> a sockaddr_un is defined as follows
 * struct sockaddr_un {
 *  short   sun_family;
 *  char    sun_path[108];
 * };
 */


#include <stdio.h>

#define NAME "socket"

/*
 * This program creates a UNIX domain datagram socket, binds a name to it,
 * then reads from the socket.
 */
main()
{
    int sock, length;
    struct sockaddr_un name;
    char buf[1024];


    /* Create socket from which to read. */
    sock = socket(AF_UNIX, SOCK_DGRAM, 0);
    if (sock < 0) {
        perror("opening datagram socket");
        exit(1);
    }


    /* Create name. */
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, NAME);


    /* Bind the UNIX domain address to the created socket */
    if (bind(sock, (struct sockaddr *) &name, sizeof(struct sockaddr_un))) {
        perror("binding name to datagram socket");
        exit(1);
    }
    printf("socket -->%s\n", NAME);


    /* Read from the socket */
    if (read(sock, buf, 1024) < 0)
        perror("receiving datagram packet");
    printf("-->%s\n", buf);
    close(sock);
    unlink(NAME);
}
```

The following sample code creates a client and sends datagrams to a server like the one created in the previous example. For the Internet domain equivalent example, see "[Sending an Internet domain datagram](http://osr507doc.sco.com/en/netguide/disockT.datagram_codesamples.html#disockD.datagram_client)".

### Sending UNIX domain datagrams

```c
#include <sys/types.h>
#include <sys/socket.h>
#include <sys/un.h>
#include <stdio.h>

#define DATA "The sea is calm tonight, the tide is full . . ."

/*
 * Send a datagram to a receiver whose name is specified in the command
 * line arguments.  The form of the command line is <programname> <pathname>
 */

main(argc, argv)
    int argc;
    char *argv[];
{
    int sock;
    struct sockaddr_un name;


    /* Create socket on which to send. */
    sock = socket(AF_UNIX, SOCK_DGRAM, 0);
    if (sock < 0) {
        perror("opening datagram socket");
        exit(1);
    }
    /* Construct name of socket to send to. */
    name.sun_family = AF_UNIX;
    strcpy(name.sun_path, argv[1]);
    /* Send message. */
    if (sendto(sock, DATA, sizeof(DATA), 0,
        &name, sizeof(struct sockaddr_un)) < 0) {
        perror("sending datagram message");
    }
    close(sock);
}
```


----------

## Difference between UNIX domain STREAM and DATAGRAM sockets?

### Q1: When I create a UNIX domain socket which is a local socket, how would it matter if the socket is STREAM socket or DATAGRAM socket. This type of socket would write the data to the socket file, would the protocol matter in this case since I am not transmitting data over a network? Is there any chance of data loss in this case if I use UNIX-based DATAGRAM sockets? 

> A1

Just as the [manual page](http://man7.org/linux/man-pages/man7/unix.7.html) says Unix sockets are always reliable. The difference between `SOCK_STREAM` and `SOCK_DGRAM` is in the semantics of consuming data out of the socket.

**Stream socket** allows for reading arbitrary number of bytes, but still preserving byte sequence. In other words, a sender might write 4K of data to the socket, and the receiver can consume that data byte by byte. The other way around is true too - sender can write several small messages to the socket that the receiver can consume in one read. **Stream socket does not preserve message boundaries**.

**Datagram socket**, on the other hand, does **preserve these boundaries** - one write by the sender always corresponds to one read by the receiver (even if receiver's buffer given to `read(2)` or `recv(2)` is smaller then that message).

> A2

The main difference is that one is **connection based (STREAM)** and the other is **connection-less (DGRAM)** - the difference between stream and packet oriented communication is usually much less important.

With `SOCK_STREAM` you still get all the connection handling, i.e. `listen`/`accept` and you can tell if a connection is closed by the other side.

### Q2: **Does UNIX DATAGRAM sockets provide better performance than UNIX STREAM sockets?** 

Performance should be the same since both types just go through local in-kernel memory, just the buffer management is different.

### Q3: **How to decide for a STREAM/DATAGRAM UNIX based socket in my application?**

So if your application protocol has **small messages with known upper bound on message size** you are better off with `SOCK_DGRAM` since that's easier to manage.

If your protocol calls for **arbitrary long message payloads**, or is **just an unstructured stream** (like raw audio or something), then pick `SOCK_STREAM` and do the required buffering.


----------


## 参考

- [UNIX Domain Socket IPC](https://akaedu.github.io/book/ch37s04.html)
- [Unix domain socket 的一些小结](http://blog.csdn.net/wlh_flame/article/details/6358795)
- [Datagrams in the UNIX domain](http://osr507doc.sco.com/en/netguide/dusockT.datagram_code_samples.html)
- [Difference between UNIX domain STREAM and DATAGRAM sockets?](https://stackoverflow.com/questions/13953912/difference-between-unix-domain-stream-and-datagram-sockets)

