

# 安装一些必要的包
```shell
apt-get -y install libpcre3 libpcre3-dev
apt-get -y install openssl libssl-dev
apt-get -y install zlib1g-dev
```

> 注意：若想在 haproxy-1.6.7 中使能 lua 功能，则必须安装 lua5.3 版本，否则编译不通过；

 在 Ubuntu 15.04 中默认只能如下 lua 版本，不符合要求；
```shell
apt-get -y install lua5.2 liblua5.2-dev
```

# 源码安装 lua5.3
```shell
wget http://www.lua.org/ftp/lua-5.3.3.tar.gz
tar zxvf lua-5.3.3.tar.gz
cd lua-5.3.3/
make linux
make install
```

> 注意：按照上面的方式，默认会将 lua5.3 安装到如下路径中
```shell
cd src && mkdir -p /usr/local/bin /usr/local/include /usr/local/lib /usr/local/man/man1 /usr/local/share/lua/5.3 /usr/local/lib/lua/5.3
cd src && install -p -m 0755 lua luac /usr/local/bin
cd src && install -p -m 0644 lua.h luaconf.h lualib.h lauxlib.h lua.hpp /usr/local/include
cd src && install -p -m 0644 liblua.a /usr/local/lib
cd doc && install -p -m 0644 lua.1 luac.1 /usr/local/man/man1
```

# 源码编译 haproxy-1.6.7
```shell
wget http://haproxy.1wt.eu/download/1.6/src/haproxy-1.6.7.tar.gz
tar zxvf haproxy-1.6.7.tar.gz
cd haproxy-1.6.7/
make TARGET=linux2628 CPU=native USE_PCRE=1 USE_OPENSSL=1 USE_ZLIB=1 USE_LUA=1
make install
```

> 注意：按照上面的方式，默认会将 haproxy 安装到如下路径中
```shell
install -d "/usr/local/sbin"
install haproxy  "/usr/local/sbin"
install -d "/usr/local/share/man"/man1
install -m 644 doc/haproxy.1 "/usr/local/share/man"/man1
install -d "/usr/local/doc/haproxy"
for x in configuration proxy-protocol management architecture cookie-options lua linux-syn-cookies network-namespaces close-options intro; do \
	install -m 644 doc/$x.txt "/usr/local/doc/haproxy" ; \
done
```

安装成功后通过 `haproxy -vv` 可以看到如下输出信息
```shell
HA-Proxy version 1.6.7 2016/07/13
Copyright 2000-2016 Willy Tarreau <willy@haproxy.org>

Build options :
  TARGET  = linux2628
  CPU     = native
  CC      = gcc
  CFLAGS  = -O2 -march=native -g -fno-strict-aliasing -Wdeclaration-after-statement
  OPTIONS = USE_ZLIB=1 USE_OPENSSL=1 USE_LUA=1 USE_PCRE=1

Default settings :
  maxconn = 2000, bufsize = 16384, maxrewrite = 1024, maxpollevents = 200

Encrypted password support via crypt(3): yes
Built with zlib version : 1.2.8
Compression algorithms supported : identity("identity"), deflate("deflate"), raw-deflate("deflate"), gzip("gzip")
Built with OpenSSL version : OpenSSL 1.0.1f 6 Jan 2014
Running on OpenSSL version : OpenSSL 1.0.1f 6 Jan 2014
OpenSSL library supports TLS extensions : yes
OpenSSL library supports SNI : yes
OpenSSL library supports prefer-server-ciphers : yes
Built with PCRE version : 8.35 2014-04-04
PCRE library supports JIT : no (USE_PCRE_JIT not set)
Built with Lua version : Lua 5.3.3
Built with transparent proxy support using: IP_TRANSPARENT IPV6_TRANSPARENT IP_FREEBIND

Available polling systems :
      epoll : pref=300,  test result OK
       poll : pref=200,  test result OK
     select : pref=150,  test result OK
Total: 3 (3 usable), will use epoll.

```


基于本地 RabbitMQ 节点构建集群的配置文件

```shell
# HAProxy Config for Local RabbitMQ Cluster

global
        log 127.0.0.1   local0 info
        maxconn 4096
        stats socket /tmp/haproxy.socket uid haproxy mode 770 level admin
        daemon

defaults
        log     global
        mode    tcp
        option  tcplog
        option  dontlognull
        retries 3
        option redispatch
        maxconn 2000
        timeout connect 5s
        timeout client 120s
        timeout server 120s

listen rabbitmq_local_cluster
    bind :5670
    mode tcp
    balance roundrobin
    server rabbit 127.0.0.1:5672 check inter 5000 rise 2 fall 3
    server rabbit_1 127.0.0.1:5673 check inter 5000 rise 2 fall 3
    server rabbit_2 127.0.0.1:5674 check inter 5000 rise 2 fall 3

listen private_monitoring
    bind :8100
    mode http
    option httplog
    stats enable
    stats uri   /stats
    stats refresh 5s

```

# 启动
```shell
haproxy -f /etc/haproxy/haproxy_rmq_cluster.cfg
```

# 重新加载
```shell
haproxy -f /etc/haproxy/haproxy_rmq_cluster.cfg -p $(pidof haproxy) -sf $(pidof haproxy)
```


