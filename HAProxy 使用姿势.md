

```shell
wget http://haproxy.1wt.eu/download/1.6/src/haproxy-1.6.7.tar.gz
tar zxvf haproxy-1.6.7.tar.gz
cd haproxy-1.6.7/
apt-cache search pcre
apt-get -y install libpcre3 libpcre3-dev
apt-get -y install openssl libssl-dev
apt-get -y install zlib1g-dev
apt-get -y install lua5.2 liblua5.2-dev
make TARGET=linux2628 CPU=native USE_PCRE=1 USE_OPENSSL=1 USE_ZLIB=1 USE_LUA=1 LUA_LIB_NAME=lua LUA_INC=/usr/include/lua5.2
make install
```








