# Programming with pcap

标签（空格分隔）： c network pcap

---

原文地址：[这里](http://www.tcpdump.org/pcap.htm)

本文档写给谁（基本要求）：

- some basic knowledge of C
- some basic understanding of networking

> All of the code examples presented here have been tested on FreeBSD 4.3 with a default kernel.

## Getting Started: The format of a pcap application

### The **layout** of a pcap sniffer

1. We begin by **determining which interface we want to sniff on**. In Linux this may be something like eth0, in BSD it may be xl1, etc. We can either define this device in a string, or we can ask pcap to provide us with the name of an interface that will do the job.
2. **Initialize pcap**. This is where we actually tell pcap what device we are sniffing on. We can, if we want to, sniff on multiple devices. How do we differentiate between them? Using file handles. Just like opening a file for reading or writing, we must name our sniffing "session" so we can tell it apart from other such sessions.
3. In the event that we only want to sniff specific traffic (e.g.: only TCP/IP packets, only packets going to port 23, etc) we must **create a rule set, "compile" it, and apply it**. This is a three phase process, all of which is closely related. The rule set is kept in a string, and is converted into a format that pcap can read (hence compiling it.) The compilation is actually just done by calling a function within our program; it does not involve the use of an external application. Then we tell pcap to apply it to whichever session we wish for it to filter.
4. Finally, we **tell pcap to enter it's primary execution loop**. In this state, pcap waits until it has received however many packets we want it to. Every time it gets a new packet in, it calls another function that we have already defined. The function that it calls can do anything we want; it can dissect the packet and print it to the user, it can save it in a file, or it can do nothing at all.
5. After our sniffing needs are satisfied, we **close our session** and are complete.

This is actually a very simple process. Five steps total, one of which is optional (step 3, in case you were wondering.) Let's take a look at each of the steps and how to implement them.

## Setting the device

可以通过两种方式指定进行 sniff 的设备；

第一种，由用户进行指定；
```c
	#include <stdio.h>
	#include <pcap.h>

	int main(int argc, char *argv[])
	{
		 char *dev = argv[1];

		 printf("Device: %s\n", dev);
		 return(0);
	}
```
> 用户必须指定正确的 interface 名字；

第二种：
```c
	#include <stdio.h>
	#include <pcap.h>

	int main(int argc, char *argv[])
	{
		char *dev, errbuf[PCAP_ERRBUF_SIZE];

		dev = pcap_lookupdev(errbuf);
		if (dev == NULL) {
			fprintf(stderr, "Couldn't find default device: %s\n", errbuf);
			return(2);
		}
		printf("Device: %s\n", dev);
		return(0);
	}
```
这种方式是由 pcap 自行设置进行 sniff 的设备；

> `errbuf` - In the event that the command fails, it will populate the string with a description of the error. 

## Opening the device for sniffing

创建 sniffing session 使用如下函数： 
```c
	pcap_t *pcap_open_live(char *device, int snaplen, int promisc, int to_ms,
	    char *ebuf)
```

参数说明：

- **device** - 上面刚刚得到的设备名字； 
- **snaplen** - 定义由 pcap 捕获的最大字节数目；
- **promisc** - 若设置为 true 则将 interface 切换到混杂（promiscuous）模式（在有些情况下，即使设置为 false，仍旧可能将 interface 切换到混杂模式）；
- **to_ms** - 以毫秒为单位的读超时时间（0 代表永不超时；在一些平台上，这意味着那你可能要等待足够多的数据包到达后才能（在用户层）看到数据包，因此建议设置非 0 超时值）；
- **ebuf** - 保存错误信息的 buffer ；

该函数将返回创建的 session handler ；

To demonstrate, consider this code snippet:
```c
	 #include <pcap.h>
	 ...
	 pcap_t *handle;

	 handle = pcap_open_live(dev, BUFSIZ, 1, 1000, errbuf);
	 if (handle == NULL) {
		 fprintf(stderr, "Couldn't open device %s: %s\n", dev, errbuf);
		 return(2);
	 }
```
> BUFSIZ 定义于 pcap.h 文件；

**A note about promiscuous vs. non-promiscuous sniffing**:
> The two techniques are very different in style. In standard, **non-promiscuous** sniffing, a host is sniffing only traffic that is directly related to it. Only traffic to, from, or routed through the host will be picked up by the sniffer. **Promiscuous** mode, on the other hand, sniffs all traffic on the wire. In a **non-switched** environment, this could be all network traffic. The obvious advantage to this is that it provides more packets for sniffing, which may or may not be helpful depending on the reason you are sniffing the network. 
> 
> However, there are regressions. 
> 
> - Promiscuous mode sniffing is detectable; a host can test with strong reliability to determine if another host is doing promiscuous sniffing. 
> - Second, it only works in a **non-switched** environment (such as a hub, or a switch that is being ARP flooded). 
> - Third, on high traffic networks, the host can become quite taxed for system resources.

**Not all devices provide the same type of link-layer headers in the packets you read.** Ethernet devices, and some non-Ethernet devices, might provide Ethernet headers, but other device types, such as **loopback** devices in BSD and OS X, **PPP** interfaces, and **Wi-Fi** interfaces when capturing in monitor mode, don't.

你需要确认目标设备提供的 link-layer headers 的类型，并基于该类型进行包内容处理；函数 `pcap_datalink()` 会返回表示 link-layer headers 类型的值；详见 [list of link-layer header type values](http://www.tcpdump.org/linktypes.html) 的说明； 其返回值带有 `DLT_` 前缀；

If your program doesn't support the link-layer header type provided by the device, it has to give up; this would be done with code such as
```c
	if (pcap_datalink(handle) != DLT_EN10MB) {
		fprintf(stderr, "Device %s doesn't provide Ethernet headers - not supported\n", dev);
		return(2);
	}
```
which fails if the device doesn't supply Ethernet headers. This would be appropriate for the code below, as it assumes Ethernet headers.

## Filtering traffic

当需要针对特定 traffic 进行 sniff 时需要使用 `pcap_compile()` 和 `pcap_setfilter()` 函数；

在调用 `pcap_open_live()` 之后，我们得到了一个 sniffing session ，然后就可以对该 session 应用我们指定的 filter 了；Why not just use our own if/else if statements? Two reasons. First, pcap's filter is far more efficient, because it does it directly with the **BPF filter**; we eliminate numerous steps by having the BPF driver do it directly. Second, this is a lot easier :)

Before applying our filter, we must "**compile**" it. filter 表达式的语法在 tcpdump 的 man 手册中有详细说明；

To compile the program we call `pcap_compile()`. 原型如下：
```c
	int pcap_compile(pcap_t *p, struct bpf_program *fp, char *str, int optimize, 
	    bpf_u_int32 netmask)
```

参数说明：

- p - 前一个函数创建的 session handle ；
- fp - a reference to the place we will store the compiled version of our filter.
- str - expression 本身；
- optimize - 标识是否对 expression 进行 "optimized" ；类型为整型，0 表示 false，1 表示 true；
- netmask - 用于限制 filter 应用范围的 network mask ；

The function returns -1 on failure; all other values imply success.

After the expression has been compiled, it is time to apply it. `pcap_setfilter()` 原型如下：
```c
	int pcap_setfilter(pcap_t *p, struct bpf_program *fp)
```
参数说明：

- p - session handler
- fp - a reference to the compiled version of the expression 

Perhaps another code sample would help to better understand:
```c
	 #include <pcap.h>
	 ...
	 pcap_t *handle;		/* Session handle */
	 char dev[] = "rl0";		/* Device to sniff on */
	 char errbuf[PCAP_ERRBUF_SIZE];	/* Error string */
	 struct bpf_program fp;		/* The compiled filter expression */
	 char filter_exp[] = "port 23";	/* The filter expression */
	 bpf_u_int32 mask;		/* The netmask of our sniffing device */
	 bpf_u_int32 net;		/* The IP of our sniffing device */

	 if (pcap_lookupnet(dev, &net, &mask, errbuf) == -1) {
		 fprintf(stderr, "Can't get netmask for device %s\n", dev);
		 net = 0;
		 mask = 0;
	 }
	 handle = pcap_open_live(dev, BUFSIZ, 1, 1000, errbuf);
	 if (handle == NULL) {
		 fprintf(stderr, "Couldn't open device %s: %s\n", dev, errbuf);
		 return(2);
	 }
	 if (pcap_compile(handle, &fp, filter_exp, 0, net) == -1) {
		 fprintf(stderr, "Couldn't parse filter %s: %s\n", filter_exp, pcap_geterr(handle));
		 return(2);
	 }
	 if (pcap_setfilter(handle, &fp) == -1) {
		 fprintf(stderr, "Couldn't install filter %s: %s\n", filter_exp, pcap_geterr(handle));
		 return(2);
	 }
```
代码功能描述：
> This program preps the sniffer to sniff all traffic coming from or going to port 23, in promiscuous mode, on the device rl0.

`pcap_lookupnet()` 功能说明：指定设备名字，返回 one of its IPv4 network numbers 和 corresponding network mask 
> **network number** is the IPv4 address ANDed with the network mask, so it contains only the network part of the address
> This was essential because we needed to know the network mask in order to apply the filter. 

> ⚠ It has been my experience that this filter does not work across all operating systems. In my test environment, I found that OpenBSD 2.9 with a default kernel does support this type of filter, but FreeBSD 4.3 with a default kernel does not. Your mileage may vary.

## The actual sniffing

There are two main techniques for capturing packets. 
- We can either capture a single packet at a time, or
- we can enter a loop that waits for n number of packets to be sniffed before being done. 

针对一次抓一个包的情况，使用函数 `pcap_next()` ，原型如下：
```c
	u_char *pcap_next(pcap_t *p, struct pcap_pkthdr *h)
```
参数说明：

- p - session handler
- h - a pointer to a structure that holds general information about the packet, specifically the **time** in which it was sniffed, the **length of this packet**, and the **length of his specific portion** (incase it is **fragmented**, for example.)

`pcap_next()` 返回一个 u_char 指针指向具体的数据包结构；

Here is a simple demonstration of `using pcap_next()` to sniff a packet.
```c
	 #include <pcap.h>
	 #include <stdio.h>

	 int main(int argc, char *argv[])
	 {
		pcap_t *handle;			/* Session handle */
		char *dev;			/* The device to sniff on */
		char errbuf[PCAP_ERRBUF_SIZE];	/* Error string */
		struct bpf_program fp;		/* The compiled filter */
		char filter_exp[] = "port 23";	/* The filter expression */
		bpf_u_int32 mask;		/* Our netmask */
		bpf_u_int32 net;		/* Our IP */
		struct pcap_pkthdr header;	/* The header that pcap gives us */
		const u_char *packet;		/* The actual packet */

		/* Define the device */
		dev = pcap_lookupdev(errbuf);
		if (dev == NULL) {
			fprintf(stderr, "Couldn't find default device: %s\n", errbuf);
			return(2);
		}
		/* Find the properties for the device */
		if (pcap_lookupnet(dev, &net, &mask, errbuf) == -1) {
			fprintf(stderr, "Couldn't get netmask for device %s: %s\n", dev, errbuf);
			net = 0;
			mask = 0;
		}
		/* Open the session in promiscuous mode */
		handle = pcap_open_live(dev, BUFSIZ, 1, 1000, errbuf);
		if (handle == NULL) {
			fprintf(stderr, "Couldn't open device %s: %s\n", dev, errbuf);
			return(2);
		}
		/* Compile and apply the filter */
		if (pcap_compile(handle, &fp, filter_exp, 0, net) == -1) {
			fprintf(stderr, "Couldn't parse filter %s: %s\n", filter_exp, pcap_geterr(handle));
			return(2);
		}
		if (pcap_setfilter(handle, &fp) == -1) {
			fprintf(stderr, "Couldn't install filter %s: %s\n", filter_exp, pcap_geterr(handle));
			return(2);
		}
		/* Grab a packet */
		packet = pcap_next(handle, &header);
		/* Print its length */
		printf("Jacked a packet with length of [%d]\n", header.len);
		/* And close the session */
		pcap_close(handle);
		return(0);
	 }
```
代码功能描述：
> This application sniffs on whatever device is returned by `pcap_lookupdev()` by putting it into **promiscuous** mode. It finds the first packet to come across port 23 (telnet) and tells the user the size of the packet (in bytes). Again, this program includes a new call, `pcap_close()`, which we will discuss later (although it really is quite self explanatory).

针对循环抓取指定数目个包的情况（更有实际价值），将使用 `pcap_loop()` 或 `pcap_dispatch()`

> 几乎没有 sniffers (if any) 会真正使用 pcap_next() 进行抓包处理；

**Callback functions** are not anything new, and are very common in many API's.  Callbacks are used in pcap, they are called when pcap sniffs a packet. The two functions that one can use to define their callback is `pcap_loop()` and `pcap_dispatch()`. Both of them call a callback function every time a packet is sniffed that meets our filter requirements (if any filter exists, of course. If not, then all packets that are sniffed are sent to the callback.)

pcap_loop() 函数原型如下：
```c
	int pcap_loop(pcap_t *p, int cnt, pcap_handler callback, u_char *user)
```
参数说明：

- p - session handle
- cnt - how many packets `pcap_loop()` should sniff for before returning（若指定负值则表示直到错误发生才停止 sniff）； 
- callback - callback function 名字；
- user - 用于传入用户指定数据；

`pcap_dispatch()` 和 `pcap_loop()` 的用法差别在于：`pcap_dispatch()` 将仅处理首批从系统接收到的数据包；而 `pcap_loop()` 将 continue processing packets or batches of packets 直到达到要求的包数量；

callback 函数 prototype 如下：
```c
	void got_packet(u_char *args, const struct pcap_pkthdr *header,
	    const u_char *packet);
```
几点说明：

- 返回值为 void ；
- 参数 args 对应 pcap_loop() 的最后一个参数；
- 参数 header 对应 **pcap header** ，其中包含了包被 sniff 的时间，包大小等信息；

pcap_pkthdr 结构体定义于 `pcap.h` 文件：
```c
	struct pcap_pkthdr {
		struct timeval ts; /* time stamp */
		bpf_u_int32 caplen; /* length of portion present */
		bpf_u_int32 len; /* length this packet (off wire) */
	};
```

- 最后一个参数 packet 对于很多 pcap 编程人员来说都是比较困惑的：It is another pointer to a u_char, and it points to the first byte of a chunk of data containing the entire packet, as sniffed by pcap_loop().

But how do you make use of this variable (named "packet" in our prototype)? A packet contains many attributes, so as you can imagine, it is not really a string, but actually a collection of structures (for instance, a TCP/IP packet would have an Ethernet header, an IP header, a TCP header, and lastly, the packet's payload). This u_char pointer points to the serialized version of these structures. To make any use of it, we must do some interesting typecasting.

首先，我们需要定义一个实际结构以便进行 **typecast** 转换；下面的结构定义用于描述 ** TCP/IP packet over Ethernet**.
```c
/* Ethernet addresses are 6 bytes */
#define ETHER_ADDR_LEN	6

	/* Ethernet header */
	struct sniff_ethernet {
		u_char ether_dhost[ETHER_ADDR_LEN]; /* Destination host address */
		u_char ether_shost[ETHER_ADDR_LEN]; /* Source host address */
		u_short ether_type; /* IP? ARP? RARP? etc */
	};

	/* IP header */
	struct sniff_ip {
		u_char ip_vhl;		/* version << 4 | header length >> 2 */
		u_char ip_tos;		/* type of service */
		u_short ip_len;		/* total length */
		u_short ip_id;		/* identification */
		u_short ip_off;		/* fragment offset field */
	#define IP_RF 0x8000		/* reserved fragment flag */
	#define IP_DF 0x4000		/* dont fragment flag */
	#define IP_MF 0x2000		/* more fragments flag */
	#define IP_OFFMASK 0x1fff	/* mask for fragmenting bits */
		u_char ip_ttl;		/* time to live */
		u_char ip_p;		/* protocol */
		u_short ip_sum;		/* checksum */
		struct in_addr ip_src,ip_dst; /* source and dest address */
	};
	#define IP_HL(ip)		(((ip)->ip_vhl) & 0x0f)
	#define IP_V(ip)		(((ip)->ip_vhl) >> 4)

	/* TCP header */
	typedef u_int tcp_seq;

	struct sniff_tcp {
		u_short th_sport;	/* source port */
		u_short th_dport;	/* destination port */
		tcp_seq th_seq;		/* sequence number */
		tcp_seq th_ack;		/* acknowledgement number */
		u_char th_offx2;	/* data offset, rsvd */
	#define TH_OFF(th)	(((th)->th_offx2 & 0xf0) >> 4)
		u_char th_flags;
	#define TH_FIN 0x01
	#define TH_SYN 0x02
	#define TH_RST 0x04
	#define TH_PUSH 0x08
	#define TH_ACK 0x10
	#define TH_URG 0x20
	#define TH_ECE 0x40
	#define TH_CWR 0x80
	#define TH_FLAGS (TH_FIN|TH_SYN|TH_RST|TH_ACK|TH_URG|TH_ECE|TH_CWR)
		u_short th_win;		/* window */
		u_short th_sum;		/* checksum */
		u_short th_urp;		/* urgent pointer */
};
```
So how does all of this relate to pcap and our mysterious `u_char` pointer? Well, those structures define the headers that appear in the data for the packet. So how can we break it apart? 

Again, we're going to assume that we are dealing with a TCP/IP packet over Ethernet. This same technique applies to any packet; the only difference is the structure types that you actually use. So let's begin by defining the variables and compile-time definitions we will need to deconstruct the packet data.
```c
/* ethernet headers are always exactly 14 bytes */
#define SIZE_ETHERNET 14

	const struct sniff_ethernet *ethernet; /* The ethernet header */
	const struct sniff_ip *ip; /* The IP header */
	const struct sniff_tcp *tcp; /* The TCP header */
	const char *payload; /* Packet payload */

	u_int size_ip;
	u_int size_tcp;
And now we do our magical typecasting:

	ethernet = (struct sniff_ethernet*)(packet);
	ip = (struct sniff_ip*)(packet + SIZE_ETHERNET);
	size_ip = IP_HL(ip)*4;
	if (size_ip < 20) {
		printf("   * Invalid IP header length: %u bytes\n", size_ip);
		return;
	}
	tcp = (struct sniff_tcp*)(packet + SIZE_ETHERNET + size_ip);
	size_tcp = TH_OFF(tcp)*4;
	if (size_tcp < 20) {
		printf("   * Invalid TCP header length: %u bytes\n", size_tcp);
		return;
	}
	payload = (u_char *)(packet + SIZE_ETHERNET + size_ip + size_tcp);
```

> 指针偏移计算（略）

The IP header, unlike the Ethernet header, does not have a fixed length; The minimum length of that header is 20 bytes.

The TCP header also has a variable length; and its minimum length is also 20 bytes.

So let's make a chart:

| Variable | Location (in bytes) |
| ---- | ---- |
| sniff_ethernet | X |
| sniff_ip | X + SIZE_ETHERNET |
| sniff_tcp | X + SIZE_ETHERNET + {IP header length} |
| payload | X + SIZE_ETHERNET + {IP header length} + {TCP header length} |

The `sniff_ethernet` structure, being the first in line, is simply at location X. `sniff_ip`, who follows directly after `sniff_ethernet`, is at the location X, plus however much space the Ethernet header consumes (14 bytes, or SIZE_ETHERNET). `sniff_tcp` is after both `sniff_ip` and `sniff_ethernet`, so it is location at X plus the sizes of the Ethernet and IP headers (14 bytes, and 4 times the IP header length, respectively). Lastly, the payload (which doesn't have a single structure corresponding to it, as its contents depends on the protocol being used atop TCP) is located after all of them.

So at this point, we know how to set our callback function, call it, and find out the attributes about the packet that has been sniffed. It's now the time you have been waiting for: writing a useful packet sniffer. Because of the length of the source code, I'm not going to include it in the body of this document. Simply download [sniffex.c](http://www.tcpdump.org/sniffex.c) and try it out.

## Wrapping Up

At this point you should be able to write a sniffer using pcap. You have learned the basic concepts behind opening a pcap session, learning general attributes about it, sniffing packets, applying filters, and using callbacks. Now it's time to get out there and sniff those wires!

