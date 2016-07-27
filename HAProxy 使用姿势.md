

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
make TARGET=linux2628 CPU=native USE_PCRE=1 USE_OPENSSL=1 USE_ZLIB=1 USE_LUA=1 LUA_LIB_NAME=lua LUA_INC=/usr/include/lua5.2
make install
```




