# Packet Analysis in Golang

标签（空格分隔）： golang packet

---

# 基于 golang.org/x/net 这个 package 能行么？

原文：[这里](https://groups.google.com/forum/#!topic/Golang-nuts/fxm9k10fi5A)

It's not so easy, and `golang.org/x/net` is not very helpful in this case.

You need to:

1. Set your network interface in promiscuous mode. This is OS dependent and will need privileges.
2. Read from the raw socket.
3. Do what you want with the packets.

Otherwise, on Linux you can use C or port examples like [these](https://stackoverflow.com/questions/114804/reading-from-a-promiscuous-network-device) to Go.


----------


# 基于 github.com/google/gopacket

一段[历史](https://www.reddit.com/r/golang/comments/174s3e/go_library_for_decoding_packet_data/)：

gopacket 作者高兴的对吃瓜群众说：我 TMD 的基于 gopcap 开源了更牛逼的 gopacket ！！

> I've just open-sourced a packet decoding library I've been working on for a bit, and it's now up at [http://code.google.com/p/gopacket](http://code.google.com/p/gopacket). (see [http://godoc.org/code.google.com/p/gopacket](http://godoc.org/code.google.com/p/gopacket) for docs). This library provides a framework for setting up packet decoders and using them to decode packets in real-time (or from pcap files).
> 
> There's a few subprojects that are really important:
> 
> - **http://godoc.org/code.google.com/p/gopacket/layers** - Decoding logic for all of the currently supported protocols (ARP, 802.1Q, EAP, EAPOL, Ethernet, EtherIP, FDDI, GRE, ICMP4/6, IP4/6, IPSec, LLC/SNAP, MPLS, PPP, PPPoE, RUDP, SCTP, TCP, UDP, UDPLite, and a few more)
> - **http://godoc.org/code.google.com/p/gopacket/pcap** - libpcap cgo bindings
> - **http://godoc.org/code.google.com/p/gopacket/pfring** - PF_RING cgo bindings
> 
> This project was forked from [https://github.com/akrennmair/gopcap](https://github.com/akrennmair/gopcap) as a more extensible method of decoding that was much simpler to add new protocol decoders to. Among other features, it provides:
> 
> - **Lazy packet decoding** - only decode layers when they're requested
> - **Modifiable decoding rules** - you can easily override any decoder or set of decoders, if you've got a better TCP decoder, for example, while keeping the rest of the decoding chain intact
> - **Extensible** - easy to add new decoders. You can even plug in new decoders without modifying the gopacket code, through a simple decoder registration system.

之后，gopcap 作者跳出来问到：NND，我的 gopcap 还在继续改进中，你凭啥认为 fork 我代码后，换个名字搞搞是一件有意义的事？

> As the original author of gopcap, I need to ask: why do you think was that fork really necessary? I'm always open to improvements.

gopacket 作者婉转的作了回答：（此处省略 10000 字）

> There's a bunch of reasons, and believe me, none of them reflect poorly on your code.
> 
> - The first is maintaining backwards compatibility... since this exposes an entirely new API, merging this into gopcap would break all current gopcap users. Since I am also currently a gopcap user (a bunch of my older code uses it), I'd really not like to see this happen. That said, I think gopacket's API is one of it's strongest points, especially its use of interfaces and its pluggability, so extending the current gopcap interface while maintaining backwards compatibility was a non-starter.
> - The second is speed. gopacket is very fast, but gopcap is still faster, due to some of the assumptions it makes (it assumes ethernet packets, for one thing, and it isn't pluggable, meaning it drops some layers of indirection). I've worked really hard on making gopacket fast, but for a vanilla IPv4/TCP packet, gopacket is still a good deal faster. At its fastest, gopacket was taking around double the time to decode one of those packets. There have been some recent regressions due to stack growing/shrinking in tight loops (runtime.morestack/lessstack) that have made gopacket even slower. It's no slouch: 1.77 us per packet on a single 3.2GHZ CPU. But gopcap is around .6-.8, if I remember from when I benchmarked it.
> - The third is simply naming: It seemed weird to throw pfring support into something named 'pcap' :) Note, though, that the pfring support is general enough that it can be used with gopcap's decoding (it just returns a byte slice, which can be passed into a gopcap.Packet).
> 
> To sum up, I'm very happy with how gopacket's turning out, but I'm also very pleased with some of gopcap's characteristics, and I don't want to break the current users of it.

gopcap 作者回复说：好吧，你说服我了（不然累～）

> OK, approved. ;-) It looks like you want to go into a different direction that gopcap, so that would be a case where a fork is totally warranted.

gopacket 作者最后又撒了一次盐：

> If you drop me the URL to that fork, I can work on merging it.

gopcap 作者流着泪默默的走开了（纯属本人歪歪～）

---

# [Packet Capture, Injection, and Analysis with Gopacket](http://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket)

> 此文被当作基于 gopacket 进行开发的基础教程；



----------


# [Reading from a promiscuous network device](https://stackoverflow.com/questions/114804/reading-from-a-promiscuous-network-device)

> 此文主要贡献了如何用 C 基于 raw socket 进行抓包；

On Linux you use a **PF_PACKET** socket to read data from a raw device, such as an ethernet interface running in promiscuous mode:
```c
s = socket(PF_PACKET, SOCK_RAW, htons(ETH_P_ALL))
```
This will send copies of every packet received up to your socket. It is quite likely that you don't really want every packet, though. The kernel can perform a first level of filtering using **BPF**, the [Berkeley Packet Filter](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter). BPF is essentially a stack-based virtual machine: it handles a small set of instructions such as:

```
ldh = load halfword (from packet)  
jeq = jump if equal  
ret = return with exit code 
```

**BPF's exit code tells the kernel whether to copy the packet to the socket or not.** It is possible to write relatively small BPF programs directly, using `setsockopt(s, SOL_SOCKET, SO_ATTACH_FILTER, )`. (WARNING: The kernel takes a `struct sock_fprog`, not a `struct bpf_program`, do not mix those up or your program will not work on some platforms).

For anything reasonably complex, you really want to use `libpcap`. BPF is limited in what it can do, in particular in the number of instructions it can execute per packet. [libpcap](http://www.tcpdump.org/pcap.htm) will take care of splitting a complex filter up into two pieces, with the kernel performing a first level of filtering and the more-capable user-space code dropping the packets it didn't actually want to see.

libpcap also abstracts the kernel interface out of your application code. Linux and BSD use similar APIs, but Solaris requires DLPI and Windows uses something else.

不同意见：
> Actually, no, libpcap won't split the filter into two pieces. But, yes, libpcap is what you want to use - it knows how to put an interface into promiscuous mode on different platforms (just using **ETH_P_ALL** isn't sufficient on Linux, for example; that's "**SAP(service access point) promiscuous**", in that you get all packets regardless of protocol type, but it's not "**physically promiscuous**", in that the adapter won't deliver to the host unicast packets not being sent to the host). 


I once had to listen on raw ethernet frames and ended up creating a wrapper for this. By calling the function with the device name, ex **eth0** I got a socket in return that was in promiscuous mode. What you need to do is to create a raw socket and then put it into promiscuous mode. Here is how I did it.

```c
int raw_init (const char *device)
{
    struct ifreq ifr;
    int raw_socket;

    memset (&ifr, 0, sizeof (struct ifreq));

    /* Open A Raw Socket */
    if ((raw_socket = socket (PF_PACKET, SOCK_RAW, htons (ETH_P_ALL))) < 1)
    {
        printf ("ERROR: Could not open socket, Got #?\n");
        exit (1);
    }

    /* Set the device to use */
    strcpy (ifr.ifr_name, device);

    /* Get the current flags that the device might have */
    if (ioctl (raw_socket, SIOCGIFFLAGS, &ifr) == -1)
    {
        perror ("Error: Could not retrive the flags from the device.\n");
        exit (1);
    }

    /* Set the old flags plus the IFF_PROMISC flag */
    ifr.ifr_flags |= IFF_PROMISC;
    if (ioctl (raw_socket, SIOCSIFFLAGS, &ifr) == -1)
    {
        perror ("Error: Could not set flag IFF_PROMISC");
        exit (1);
    }
    printf ("Entering promiscuous mode\n");

    /* Configure the device */

    if (ioctl (raw_socket, SIOCGIFINDEX, &ifr) < 0)
    {
        perror ("Error: Error getting the device index.\n");
        exit (1);
    }

    return raw_socket;
}
```

Then when you have your socket you can just use select to handle packets as they arrive.


----------


# 杂七杂八

- Here is a good example: [github.com/google/stenographer](https://github.com/google/stenographer)
and a good tutorial: [http://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket](http://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket )

- [EtherApe](http://etherape.sourceforge.net/): EtherApe is a graphical network monitor for Unix modeled after etherman. Featuring link layer, IP and TCP modes, it displays network activity graphically. Hosts and links change in size with traffic. Color coded protocols display. It supports Ethernet, FDDI, Token Ring, ISDN, PPP, SLIP and WLAN devices, plus several encapsulation formats. It can filter traffic to be shown, and can read packets from a file as well as live from the network. Node statistics can be exported.

- yes using **raw sockets** to capture packets works fine in Linux... but not on BSDs. that is... unless you are throwing around the term "raw socket" to mean "whatever capture method is available such as the BSD **BPF** etc.". 

- **gopacket** can use `libpcap`... but i don't see any reason to use libpcap since gopacket also supports other packet capture methods that are much faster and safer (without linking to an old C library)... such as AF_PACKET which is Linux only. And if you want to capture packets on a BSD system then gopacket supports `BPF`. I've tested it on OpenBSD, NetBSD and FreeBSD. 

- if you do care about having super duper performance then you must write your go library that uses faster capture methods such as `netmap` or `DPDK`; oh but then you have memory safety issues with DPDK. Do you want it fast, safe or correct? 

- You might be able to have all three properties; Safe, Correct, and Fast Low-Level Networking by Robert Clipsham [http://octarineparrot.com/assets/msci_paper.pdf](http://octarineparrot.com/assets/msci_paper.pdf) 


