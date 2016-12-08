# Ubuntu 15.04 上源码安装 mysql-5.7.16

标签（空格分隔）： mysql ubuntu

---

# 源码安装

官网参考：[这里](http://dev.mysql.com/doc/refman/5.7/en/installing-source-distribution.html)

官网给出的 “basic installation command sequence” 如下：

## 预配置阶段

```shell
# Preconfiguration setup
shell> groupadd mysql
shell> useradd -r -g mysql -s /bin/false mysql
```

> 问题：执行 `-s /bin/false` 命令的作用是？

## 源码安装阶段

```shell
# Beginning of source-build specific instructions
shell> tar zxvf mysql-VERSION.tar.gz
shell> cd mysql-VERSION
shell> cmake .
shell> make
shell> make install
# End of source-build specific instructions
```

> 问题：仅执行 `cmake .` 命令显然不够，如何定制化 cmake 参数？

## 安装后的处理

```shell
# Postinstallation setup
shell> cd /usr/local/mysql
shell> chown -R mysql .
shell> chgrp -R mysql .
shell> bin/mysql_install_db --user=mysql    # Before MySQL 5.7.6
shell> bin/mysqld --initialize --user=mysql # MySQL 5.7.6 and up
shell> bin/mysql_ssl_rsa_setup              # MySQL 5.7.6 and up
shell> chown -R root .
shell> chown -R mysql data
shell> bin/mysqld_safe --user=mysql &
# Next command is optional
shell> cp support-files/mysql.server /etc/init.d/mysql.server
```

> 问题：执行 `bin/mysql_ssl_rsa_setup` 命令的作用是？ 

## 小结

官网给出的安装命令过于粗糙：

- 相关目录是在何时，通过何种方式被创建出来的没有说明；
- 命令 chown 和 chgrp 明显可以进行简化组合；
- `bin/mysqld_safe` 在 5.7.xx 版本中已经被移除；
- 将 `support-files/mysql.server` 用作 init.d 下的脚本用于维护 mysql 服务的方式在已经支持 systemd 的系统中已经不合适了；


# 安装过程遇到的问题

## 可能用到的命令

```shell
# To list the configuration options
shell> cmake . -L   # overview
shell> cmake . -LH  # overview with help text
shell> cmake . -LAH # all params with help text
shell> ccmake .     # interactive display

# To prevent old object files or configuration information from being used, run these commands before re-running CMake
shell> make clean
shell> rm CMakeCache.txt


```

## boost 版本低问题

在源码安装系统要求中提到：

> The Boost C++ libraries are required to build MySQL (but not to use it). **Boost 1.59.0 must be installed**. To obtain Boost and its installation instructions, visit the [official site](http://www.boost.org/). After Boost is installed, tell the build system where the Boost files are located by defining the [WITH_BOOST](http://dev.mysql.com/doc/refman/5.7/en/source-configuration-options.html#option_cmake_with_boost) option when you invoke **CMake**. For example:
>
> `shell> cmake . -DWITH_BOOST=/usr/local/boost_1_59_0`
>
> Adjust the path as necessary to match your installation.

这种方式采用了先安装后配置的方式；事实上，可以基于 cmake 直接进行下载安装；

```shell
root@vagrant-ubuntu-trusty:~/workspace/WGET/mysql-5.7.16# cmake . -L
-- Running cmake version 3.0.2
-- Found Git: /usr/bin/git (found version "2.1.4")
-- Configuring with MAX_INDEXES = 64U
...
-- MySQL 5.7.16
-- Packaging as: mysql-5.7.16-Linux-x86_64
-- Found /usr/include/boost/version.hpp
-- BOOST_VERSION_NUMBER is #define BOOST_VERSION 105500
CMake Warning at cmake/boost.cmake:266 (MESSAGE):
  Boost minor version found is 55 we need 59
Call Stack (most recent call first):
  CMakeLists.txt:455 (INCLUDE)

-- BOOST_INCLUDE_DIR /usr/include
-- LOCAL_BOOST_DIR
-- LOCAL_BOOST_ZIP
-- Could not find (the correct version of) boost.
-- MySQL currently requires boost_1_59_0

CMake Error at cmake/boost.cmake:81 (MESSAGE):
  You can download it with -DDOWNLOAD_BOOST=1 -DWITH_BOOST=<directory>

  This CMake script will look for boost in <directory>.  If it is not there,
  it will download and unpack it (in that directory) for you.

  If you are inside a firewall, you may need to use an http proxy:

  export http_proxy=http://example.com:80

Call Stack (most recent call first):
  cmake/boost.cmake:269 (COULD_NOT_FIND_BOOST)
  CMakeLists.txt:455 (INCLUDE)

-- Configuring incomplete, errors occurred!
See also "/root/workspace/WGET/mysql-5.7.16/CMakeFiles/CMakeOutput.log".
See also "/root/workspace/WGET/mysql-5.7.16/CMakeFiles/CMakeError.log".
-- Cache values
...
root@vagrant-ubuntu-trusty:~/workspace/WGET/mysql-5.7.16#
```

可以看到，在 ubuntu 15.04 系统中已安装的 boost 版本不满足 “Boost minor version found is 55 we need 59” 要求，同时给出了修复办法 “You can download it with -DDOWNLOAD_BOOST=1 -DWITH_BOOST=<directory>”；

基于 cmake 直接更新 boost 版本的方法如下：
```shell
root@vagrant-ubuntu-trusty:~/workspace/WGET/mysql-5.7.16# cmake . -DDOWNLOAD_BOOST=1 -DWITH_BOOST=/tmp
```

之后就 ok 了；


## systemd 相关 cmake 配置问题

参考信息：[这里](http://dev.mysql.com/doc/refman/5.7/en/using-systemd.html)

> 关键信息：
>
> - As of MySQL 5.7.6, if you install MySQL using an RPM distribution on the following Linux platforms, server startup and shutdown is managed by systemd
> - To obtain systemd support if you install from a source distribution, configure the distribution using the [-DWITH_SYSTEMD=1](http://dev.mysql.com/doc/refman/5.7/en/source-configuration-options.html#option_cmake_with_systemd) CMake option.
> - manual server management using the systemctl command or service command.
> - On platforms for which systemd support is installed, scripts such as mysqld_safe and the System V initialization script are not installed because they are unnecessary. For example, mysqld_safe can handle server restarts, but systemd provides the same capability, and does so in a manner consistent with management of other services rather than using an application-specific program.
> - As of MySQL 5.7.13, on platforms for which systemd support is installed, systemd has the capability of managing multiple MySQL instances. 
> - Because mysqld_safe is not installed when systemd is used, options previously specified for that program (for example, in an [mysqld_safe] option group) must be specified another way.


**`-DSYSTEMD_PID_DIR=dir_name`**

> The name of the directory in which to create the PID file when MySQL is managed by systemd. The default is /var/run/mysqld; this might be changed implicitly according to the INSTALL_LAYOUT value.
> 
> This option is ignored unless WITH_SYSTEMD is enabled. It was added in MySQL 5.7.6.

**`-DSYSTEMD_SERVICE_NAME=name`**

> The name of the MySQL service to use when MySQL is managed by systemd. The default is mysqld; this might be changed implicitly according to the INSTALL_LAYOUT value.
> 
> This option is ignored unless WITH_SYSTEMD is enabled. It was added in MySQL 5.7.6.

**`-DWITH_SYSTEMD=bool`**

> Whether to enable installation of systemd support files. By default, this option is disabled. When enabled, systemd support files are installed, and scripts such as mysqld_safe and the System V initialization script are not installed. On platforms where systemd is not available, enabling WITH_SYSTEMD results in an error from CMake.
> 
> For more information about using systemd, see Section 2.5.10, “Managing MySQL Server with systemd”. That section also includes information about specifying options previously specified in [mysqld_safe] option groups. Because mysqld_safe is not installed when systemd is used, such options must be specified another way.
> 
> This option was added in MySQL 5.7.6.

## “c++: internal compiler error: Killed (program cc1plus)” 问题

在执行 make 命令时，会出现如下错误

```shell
root@vagrant-ubuntu-trusty:~/workspace/WGET/mysql-5.7.16# make
...
[ 38%] Building CXX object sql/CMakeFiles/sql.dir/item_func.cc.o
[ 39%] Building CXX object sql/CMakeFiles/sql.dir/item_geofunc.cc.o
c++: internal compiler error: Killed (program cc1plus)
Please submit a full bug report,
with preprocessed source if appropriate.
See <file:///usr/share/doc/gcc-4.9/README.Bugs> for instructions.
sql/CMakeFiles/sql.dir/build.make:907: recipe for target 'sql/CMakeFiles/sql.dir/item_geofunc.cc.o' failed
make[2]: *** [sql/CMakeFiles/sql.dir/item_geofunc.cc.o] Error 4
CMakeFiles/Makefile2:8443: recipe for target 'sql/CMakeFiles/sql.dir/all' failed
make[1]: *** [sql/CMakeFiles/sql.dir/all] Error 2
Makefile:137: recipe for target 'all' failed
make: *** [all] Error 2
root@vagrant-ubuntu-trusty:~/workspace/WGET/mysql-5.7.16#
```

可以看到 cc1plus 是由于 “Out of memory” 的原因被杀了；

```shell
root@vagrant-ubuntu-trusty:~# dmesg |grep -i kill
[29269.576976] automount invoked oom-killer: gfp_mask=0x201da, order=0, oom_score_adj=0
[29269.577051]  [<ffffffff8117d04b>] oom_kill_process+0x22b/0x390
[29269.577294] Out of memory: Kill process 4094 (cc1plus) score 866 or sacrifice child
[29269.577339] Killed process 4094 (cc1plus) total-vm:965712kB, anon-rss:453088kB, file-rss:0kB
root@vagrant-ubuntu-trusty:~#
```

问题的根源在于，我的 ubuntu 15.04 虚拟机环境只分配了 500M 的内存，结果表明是不够用的；

在 stackoverflow 上找了个比较好的[说明](http://stackoverflow.com/questions/30887143/make-j-8-g-internal-compiler-error-killed-program-cc1plus)，摘录如下：

> Most likely that is your problem. The problem above occurs when your system runs out of memory. In this case rather than the whole system falling over, the operating systems runs a process to score each process on the system. The one that scores the highest gets killed by the operating system to free up memory. If the process that is killed is **cc1plus**, gcc (perhaps incorrectly) interprets this as the process crashing and hence assumes that it must be a compiler bug. But it isn't really, the problem is the OS killed cc1plus, rather than it crashed.
> 
> If this is the case, you are running out of memory. 


解决办法可以参考[这里](http://xwsoul.com/posts/684)；摘录如下：

```
# 创建分区文件, 大小 2G
dd if=/dev/zero of=/swapfile bs=1k count=2048000
# 生成 swap 文件系统
mkswap /swapfile
# 激活 swap 文件
swapon /swapfile
# 在系统重启的时候自动挂载交换分区，需要在 /etc/fstab 文件中增加如下内容
/swapfile  swap  swap    defaults 0 0
```

通过上面的办法确实能够避免由于 oom 而被杀的问题，与此同时，因为使用了 swap 所以编译速度也慢的可以～～


## `--explicit_defaults_for_timestamp` 选项问题

在进行 mysql 初始化时，会有如下告警信息；

```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# bin/mysqld --initialize --user=mysql
2016-12-08T02:39:01.362885Z 0 [Warning] TIMESTAMP with implicit DEFAULT value is deprecated. Please use --explicit_defaults_for_timestamp server option (see documentation for more details).
2016-12-08T02:39:01.739983Z 0 [Warning] InnoDB: New log files created, LSN=45790
2016-12-08T02:39:01.804289Z 0 [Warning] InnoDB: Creating foreign key constraint system tables.
2016-12-08T02:39:01.875989Z 0 [Warning] No existing UUID has been found, so we assume that this is the first time that this server has been started. Generating a new UUID: 7a7039cf-bcef-11e6-a01e-0800274ac42f.
2016-12-08T02:39:01.878411Z 0 [Warning] Gtid table is not ready to be used. Table 'mysql.gtid_executed' cannot be opened.
2016-12-08T02:39:01.880601Z 1 [Note] A temporary password is generated for root@localhost: 7rc!uu#i!%Az
root@vagrant-ubuntu-trusty:/usr/local/mysql#
```

原则上该告警可以不解决；官网给出的[参数说明](http://dev.mysql.com/doc/refman/5.6/en/server-system-variables.html#sysvar_explicit_defaults_for_timestamp)如下：

**`explicit_defaults_for_timestamp`**

```shell
Introduced	5.6.6
Deprecated	5.6.6
Command-Line Format	--explicit_defaults_for_timestamp=#
System Variable	Name	explicit_defaults_for_timestamp
Variable Scope	Global, Session
Dynamic Variable	No
Permitted Values	Type	boolean
Default	FALSE
In MySQL, the TIMESTAMP data type differs in nonstandard ways from other data types:

TIMESTAMP columns not explicitly declared with the NULL attribute are assigned the NOT NULL attribute. (Columns of other data types, if not explicitly declared as NOT NULL, permit NULL values.) Setting such a column to NULL sets it to the current timestamp.

The first TIMESTAMP column in a table, if not declared with the NULL attribute or an explicit DEFAULT or ON UPDATE clause, is automatically assigned the DEFAULT CURRENT_TIMESTAMP and ON UPDATE CURRENT_TIMESTAMP attributes.

TIMESTAMP columns following the first one, if not declared with the NULL attribute or an explicit DEFAULT clause, are automatically assigned DEFAULT '0000-00-00 00:00:00' (the “zero” timestamp). For inserted rows that specify no explicit value for such a column, the column is assigned '0000-00-00 00:00:00' and no warning occurs.

Those nonstandard behaviors remain the default for TIMESTAMP but as of MySQL 5.6.6 are deprecated and this warning appears at startup:

[Warning] TIMESTAMP with implicit DEFAULT value is deprecated.
Please use --explicit_defaults_for_timestamp server option (see
documentation for more details).
As indicated by the warning, to turn off the nonstandard behaviors, enable the explicit_defaults_for_timestamp system variable at server startup. With this variable enabled, the server handles TIMESTAMP as follows instead:

TIMESTAMP columns not explicitly declared as NOT NULL permit NULL values. Setting such a column to NULL sets it to NULL, not the current timestamp.

No TIMESTAMP column is assigned the DEFAULT CURRENT_TIMESTAMP or ON UPDATE CURRENT_TIMESTAMP attributes automatically. Those attributes must be explicitly specified.

TIMESTAMP columns declared as NOT NULL and without an explicit DEFAULT clause are treated as having no default value. For inserted rows that specify no explicit value for such a column, the result depends on the SQL mode. If strict SQL mode is enabled, an error occurs. If strict SQL mode is not enabled, the column is assigned the implicit default of '0000-00-00 00:00:00' and a warning occurs. This is similar to how MySQL treats other temporal types such as DATETIME.

Note
explicit_defaults_for_timestamp is itself deprecated because its only purpose is to permit control over now-deprecated TIMESTAMP behaviors that will be removed in a future MySQL release. When that removal occurs, explicit_defaults_for_timestamp will have no purpose and will be removed as well.

This variable was added in MySQL 5.6.6.
```

简单的处理方式如下
```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# bin/mysqld --initialize --user=mysql --explicit_defaults_for_timestamp=true
```

## 关于 my.cnf 和 mysqld.service 文件

### my.cnf

```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# cp support-files/my-default.cnf /etc/my.cnf
```

### mysqld.service

```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# cp lib/systemd/system/mysqld.service /usr/lib/systemd/system/
```

## mysqld 启动失败问题

启动 mysql 服务，失败；

```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# systemctl start mysqld.service
Job for mysqld.service failed. See "systemctl status mysqld.service" and "journalctl -xe" for details.
root@vagrant-ubuntu-trusty:/usr/local/mysql#
root@vagrant-ubuntu-trusty:/usr/local/mysql# systemctl status mysqld.service
● mysqld.service - MySQL Server
   Loaded: loaded (/usr/lib/systemd/system/mysqld.service; disabled; vendor preset: enabled)
   Active: failed (Result: start-limit) since Thu 2016-12-08 14:34:27 CST; 17s ago
  Process: 24527 ExecStart=/usr/local/mysql/bin/mysqld --daemonize --pid-file=/var/run/mysqld/mysqld.pid $MYSQLD_OPTS (code=exited, status=1/FAILURE)
  Process: 24512 ExecStartPre=/usr/local/mysql/bin/mysqld_pre_systemd (code=exited, status=0/SUCCESS)

Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: mysqld.service: control process exited, code=exited status=1
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: Failed to start MySQL Server.
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: Unit mysqld.service entered failed state.
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: mysqld.service failed.
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: mysqld.service holdoff time over, scheduling restart.
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: start request repeated too quickly for mysqld.service
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: Failed to start MySQL Server.
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: Unit mysqld.service entered failed state.
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: mysqld.service failed.
root@vagrant-ubuntu-trusty:/usr/local/mysql#
```

查看 systemd 的日志，可以看到服务启动的详细信息；

```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# journalctl -xe
...
Dec 08 14:34:27 vagrant-ubuntu-trusty mysqld[24527]: 2016-12-08T06:34:27.628759Z 0 [ERROR] /usr/local/mysql/bin/mysqld: Can't create/write to file '/var/run/mysqld/mysqld.pid' (Errcode: 2 - No such file or directory)
Dec 08 14:34:27 vagrant-ubuntu-trusty mysqld[24527]: 2016-12-08T06:34:27.629114Z 0 [ERROR] Can't start server: can't create PID file: No such file or directory
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: mysqld.service: control process exited, code=exited status=1
Dec 08 14:34:27 vagrant-ubuntu-trusty systemd[1]: Failed to start MySQL Server.
-- Subject: Unit mysqld.service has failed
-- Defined-By: systemd
-- Support: http://lists.freedesktop.org/mailman/listinfo/systemd-devel
--
-- Unit mysqld.service has failed.
--
-- The result is failed.
...
```

从上面可以看出，失败的原因是由于 “Can't create/write to file '/var/run/mysqld/mysqld.pid'” ；

可以通过如下命令解决
```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# mkdir -p /var/run/mysqld
root@vagrant-ubuntu-trusty:/usr/local/mysql# chown mysql:mysql /var/run/mysqld
```

再次重新启动，成功；

```shell
root@vagrant-ubuntu-trusty:/usr/local/mysql# systemctl start mysqld.service
root@vagrant-ubuntu-trusty:/usr/local/mysql# systemctl status mysqld.service
● mysqld.service - MySQL Server
   Loaded: loaded (/usr/lib/systemd/system/mysqld.service; disabled; vendor preset: enabled)
   Active: active (running) since Thu 2016-12-08 14:44:33 CST; 17s ago
  Process: 24760 ExecStart=/usr/local/mysql/bin/mysqld --daemonize --pid-file=/var/run/mysqld/mysqld.pid $MYSQLD_OPTS (code=exited, status=0/SUCCESS)
  Process: 24745 ExecStartPre=/usr/local/mysql/bin/mysqld_pre_systemd (code=exited, status=0/SUCCESS)
 Main PID: 24764 (mysqld)
   CGroup: /system.slice/mysqld.service
           └─24764 /usr/local/mysql/bin/mysqld --daemonize --pid-file=/var/run/mysqld/mysqld.pid

Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.384861Z 0 [Note] InnoDB: Loading buffer pool(s) from /usr/local/mysql/data/ib_buffer_pool
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.386295Z 0 [Note] InnoDB: Buffer pool(s) load completed at 161208 14:44:33
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.386705Z 0 [Note] Server hostname (bind-address): '*'; port: 3306
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.386936Z 0 [Note] IPv6 is available.
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.387219Z 0 [Note]   - '::' resolves to '::';
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.387398Z 0 [Note] Server socket created on IP: '::'.
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.401092Z 0 [Note] Event Scheduler: Loaded 0 events
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: 2016-12-08T06:44:33.403270Z 0 [Note] /usr/local/mysql/bin/mysqld: ready for connections.
Dec 08 14:44:33 vagrant-ubuntu-trusty mysqld[24760]: Version: '5.7.16'  socket: '/tmp/mysql.sock'  port: 3306  Source distribution
Dec 08 14:44:33 vagrant-ubuntu-trusty systemd[1]: Started MySQL Server.
root@vagrant-ubuntu-trusty:/usr/local/mysql#
```

## 关于 mysql_secure_installation 命令

详细信息参考官网文章：[这里](https://dev.mysql.com/doc/refman/5.7/en/mysql-secure-installation.html)

This program enables you to improve the security of your MySQL installation in the following ways:

- You can set a password for **root** accounts.
- You can remove root accounts that are accessible from outside the local host.
- You can remove anonymous-user accounts.
- You can remove the test database (which by default can be accessed by all users, even anonymous users), and privileges that permit anyone to access databases with names that start with test_.

需要注意的是，若要允许远端机器进行访问，则需要在执行该命令时进行相关设置；之后还需要为 mysql 添加远程访问账号和密码；

```shell
mysql> GRANT ALL PRIVILEGES ON *.* TO 'root'@'%' IDENTIFIED BY '123456' WITH GRANT OPTION;
```

----------


最终用于生成 Makefile 文件的 cmake 指令

```shell
root@vagrant-ubuntu-trusty:~/workspace/WGET/mysql-5.7.16# cmake . -DCMAKE_INSTALL_PREFIX=/usr/local/mysql -DMYSQL_DATADIR=/usr/local/mysql/data -DDEFAULT_CHARSET=utf8 -DDEFAULT_COLLATION=utf8_general_ci -DEXTRA_CHARSETS=all -DENABLED_LOCAL_INFILE=1 -DDOWNLOAD_BOOST=1 -DWITH_BOOST=/tmp -DWITH_SYSTEMD=1 -DENABLE_DOWNLOADS=1
```

详细 cmake 配置说明见[这里](http://dev.mysql.com/doc/refman/5.7/en/source-configuration-options.html)；










