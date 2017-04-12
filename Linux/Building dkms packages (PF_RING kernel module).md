# Building dkms packages (PF_RING kernel module)

标签（空格分隔）： PF_RING

---

> 原文：[这里](https://github.com/ntop/PF_RING/blob/dev/kernel/README)；

编译 PF_RING 内核模块前，首先需安装 kernel 头文件；

- RedHat/CentOS

```
# yum install kernel-devel
```

- Debian/Ubuntu

```
# apt-get install linux-headers-$(uname -r)
```

## 构建 dkms packages

- 将你自己的用户名（替换 XXXXXX 为你自己的账号信息）添加到 sudo 列表中，以便后续构建 package 时不必再次输入密码；

```
# sudo visudo
XXXXXX ALL=NOPASSWD: /usr/bin/make
```

- 以普通用户（非 sudo）执行 

```
# configure
```

- Ubuntu

```
# sudo make -f Makefile.dkms deb
```

- CentoOS/Fedora/RedHat

```
# sudo make -f Makefile.dkms rpm
```


----------

> 以下为基于 `Makefile.dkms` 文件进行 PF_RING 内核模块编译安装的过程；

`Makefile.dkms` 文件内容如下：

```shell
GIT_REV:=unknown
ifneq (, $(shell which git))
 ifeq (, $(shell echo ${SUBDIRS}))
  GIT_BRANCH=$(shell git branch --no-color|cut -d ' ' -f 2)
  GIT_HASH=$(shell git rev-parse HEAD)
  ifneq ($(strip $(GIT_BRANCH)),)
   GIT_REV:=${GIT_BRANCH}:${GIT_HASH}
  endif
 endif
endif

all: install

add: remove
        \/bin/rm -rf /usr/src/pfring-6.5.0
        mkdir /usr/src/pfring-6.5.0
        cp -r Makefile dkms.conf pf_ring.c linux/ /usr/src/pfring-6.5.0
        cat Makefile | sed -e "s/GIT_REV:=$$/GIT_REV:=${GIT_REV}/" > /usr/src/pfring-6.5.0/Makefile
        dkms add -m pfring -v 6.5.0

build: add
        dkms build -m pfring -v 6.5.0

install: build
        dkms install --force -m pfring -v 6.5.0

deb: add add_deb install
        dkms mkdeb -m pfring -v 6.5.0 --source-only

rpm: add add_rpm install
        dkms mkrpm -m pfring -v 6.5.0 --source-only

add_rpm:
        cp -r pfring-dkms-*.spec /usr/src/pfring-6.5.0/

add_deb:
        cp -r pfring-dkms-mkdeb /usr/src/pfring-6.5.0/

remove:
        -dkms remove -m pfring -v 6.5.0 --all

veryclean: remove
        \/bin/rm -fr /usr/src/pfring-6.5.0
```

具体执行过程：

- 系统环境

```shell
[root@xg-esm-data-4 ~]# lsb_release -d
Description:	CentOS Linux release 7.1.1503 (Core)
[root@xg-esm-data-4 ~]#
[root@xg-esm-data-4 ~]# uname -a
Linux xg-esm-data-4 3.10.0-229.el7.x86_64 #1 SMP Fri Mar 6 11:36:42 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
[root@xg-esm-data-4 ~]#

[root@xg-esm-data-4 kernel]# ll
total 2988
-rwxr-xr-x 1 root root   80070 Apr 10 18:22 configure
-rw-r--r-- 1 root root     551 Apr 10 18:22 configure.in
-rw-r--r-- 1 root root     184 Apr 10 18:22 dkms.conf.in
drwxr-xr-x 2 root root    4096 Apr 10 18:22 linux
-rw-r--r-- 1 root root    1945 Apr 10 18:22 Makefile
-rw-r--r-- 1 root root    1029 Apr 10 18:22 Makefile.dkms.in
-rw-r--r-- 1 root root      47 Apr 10 18:24 modules.order
-rw-r--r-- 1 root root     184 Apr 10 18:24 Module.symvers
-rw-r--r-- 1 root root  256302 Apr 10 18:22 pf_ring.c
drwxr-xr-x 3 root root    4096 Apr 10 18:22 pfring-dkms-mkdeb
-rw-r--r-- 1 root root    3024 Apr 10 18:22 pfring-dkms-mkrpm.spec.in
-rw-r--r-- 1 root root 1327554 Apr 10 18:24 pf_ring.ko
-rw-r--r-- 1 root root    6482 Apr 10 18:24 pf_ring.mod.c
-rw-r--r-- 1 root root   58992 Apr 10 18:24 pf_ring.mod.o
-rw-r--r-- 1 root root 1274256 Apr 10 18:24 pf_ring.o
-rw-r--r-- 1 root root     580 Apr 10 18:22 README
[root@xg-esm-data-4 kernel]#
```

- 执行 configure

```shell
[root@xg-esm-data-4 kernel]# ./configure
configure: creating ./config.status
config.status: creating Makefile.dkms
config.status: creating dkms.conf
config.status: creating pfring-dkms-mkrpm.spec
config.status: creating pfring-dkms-mkdeb/debian/changelog
config.status: creating pfring-dkms-mkdeb/debian/control
[root@xg-esm-data-4 kernel]#


[root@xg-esm-data-4 kernel]# ll
total 3028
-rw-r--r-- 1 root root    3439 Apr 12 10:47 config.log
-rwxr-xr-x 1 root root   23481 Apr 12 10:47 config.status   -- new
-rwxr-xr-x 1 root root   80070 Apr 10 18:22 configure
-rw-r--r-- 1 root root     551 Apr 10 18:22 configure.in
-rw-r--r-- 1 root root     169 Apr 12 10:47 dkms.conf   -- new
-rw-r--r-- 1 root root     184 Apr 10 18:22 dkms.conf.in
drwxr-xr-x 2 root root    4096 Apr 10 18:22 linux
-rw-r--r-- 1 root root    1945 Apr 10 18:22 Makefile
-rw-r--r-- 1 root root    1016 Apr 12 10:47 Makefile.dkms   -- new
-rw-r--r-- 1 root root    1029 Apr 10 18:22 Makefile.dkms.in
-rw-r--r-- 1 root root      47 Apr 10 18:24 modules.order
-rw-r--r-- 1 root root     184 Apr 10 18:24 Module.symvers
-rw-r--r-- 1 root root  256302 Apr 10 18:22 pf_ring.c
drwxr-xr-x 3 root root    4096 Apr 10 18:22 pfring-dkms-mkdeb   -- updated
-rw-r--r-- 1 root root    3063 Apr 12 10:47 pfring-dkms-mkrpm.spec   -- new
-rw-r--r-- 1 root root    3024 Apr 10 18:22 pfring-dkms-mkrpm.spec.in
-rw-r--r-- 1 root root 1327554 Apr 10 18:24 pf_ring.ko
-rw-r--r-- 1 root root    6482 Apr 10 18:24 pf_ring.mod.c
-rw-r--r-- 1 root root   58992 Apr 10 18:24 pf_ring.mod.o
-rw-r--r-- 1 root root 1274256 Apr 10 18:24 pf_ring.o
-rw-r--r-- 1 root root     580 Apr 10 18:22 README
[root@xg-esm-data-4 kernel]#
```

- make（发生错误）

```
[root@xg-esm-data-4 kernel]# make -f Makefile.dkms rpm
dkms remove -m pfring -v 6.5.0 --all
make: dkms: Command not found          -- 错误信息
make: [remove] Error 127 (ignored)
\/bin/rm -rf /usr/src/pfring-6.5.0
mkdir /usr/src/pfring-6.5.0
cp -r Makefile dkms.conf pf_ring.c linux/ /usr/src/pfring-6.5.0
cat Makefile | sed -e "s/GIT_REV:=$/GIT_REV:=dev:ce37e8d2bef24dfc8f9cf28b7b1387276928cbd5/" > /usr/src/pfring-6.5.0/Makefile
dkms add -m pfring -v 6.5.0
make: dkms: Command not found
make: *** [add] Error 127
[root@xg-esm-data-4 kernel]#
```

- 安装 dkms

```shell
[root@xg-esm-data-4 kernel]# yum -y install dkms
Loaded plugins: fastestmirror, langpacks, priorities
...
Installed:
  dkms.noarch 0:2.2.0.3-30.git.7c3e7c5.el7

Complete!
[root@xg-esm-data-4 kernel]#
```

- make

```shell
[root@xg-esm-data-4 kernel]# make -f Makefile.dkms rpm

## 从 dkms 中移除 pfring 6.5.0
dkms remove -m pfring -v 6.5.0 --all
Error! There are no instances of module: pfring
6.5.0 located in the DKMS tree.
make: [remove] Error 3 (ignored)

## 目录组织形式需要满足 dkms 要求
\/bin/rm -rf /usr/src/pfring-6.5.0
mkdir /usr/src/pfring-6.5.0
cp -r Makefile dkms.conf pf_ring.c linux/ /usr/src/pfring-6.5.0
cat Makefile | sed -e "s/GIT_REV:=$/GIT_REV:=dev:ce37e8d2bef24dfc8f9cf28b7b1387276928cbd5/" > /usr/src/pfring-6.5.0/Makefile

## 向 dkms 中添加 pfring 6.5.0
dkms add -m pfring -v 6.5.0

Creating symlink /var/lib/dkms/pfring/6.5.0/source ->
                 /usr/src/pfring-6.5.0

DKMS: add completed.

## 用于构建 RPM 的模版文件
cp -r pfring-dkms-*.spec /usr/src/pfring-6.5.0/

## 构建 pfring 6.5.0
dkms build -m pfring -v 6.5.0

Kernel preparation unnecessary for this kernel.  Skipping...

Building module:
cleaning build area...
make KERNELRELEASE=3.10.0-229.el7.x86_64....
cleaning build area...

DKMS: build completed.

## 安装 pfring 6.5.0 到 kernel
dkms install --force -m pfring -v 6.5.0

pf_ring:
Running module version sanity check.
 - Original module
   - Found /lib/modules/3.10.0-229.el7.x86_64/kernel/net/pf_ring/pf_ring.ko
   - Storing in /var/lib/dkms/pfring/original_module/3.10.0-229.el7.x86_64/x86_64/
   - Archiving for uninstallation purposes
 - Installation
   - Installing to /lib/modules/3.10.0-229.el7.x86_64/extra/
Adding any weak-modules

## 生成 modules.dep 和 map 文件
depmod...

DKMS: install completed.

## 构建 RPM
dkms mkrpm -m pfring -v 6.5.0 --source-only

copying legacy postinstall template...
Copying source tree...

## 底层构建 RPM 包命令
rpmbuild...

Wrote: /var/lib/dkms/pfring/6.5.0/rpm/pfring-dkms-6.5.0-1187.src.rpm
Wrote: /var/lib/dkms/pfring/6.5.0/rpm/pfring-dkms-6.5.0-1187.noarch.rpm

DKMS: mkrpm completed.
[root@xg-esm-data-4 kernel]#
```


----------


## dkms 说明

dkms - Dynamic Kernel Module Support

```
dkms [action] [options] [module/module-version] [/path/to/source-tree] [/path/to/tarball.tar] [/path/to/driver.rpm]
```

- dkms 是一种框架；
- 该框架允许针对存在于目标系统中的每一种 kernel 进行 kernel modules 的动态构建的；
- 该框架保证构建方式简单且易于组织；

### ACTIONS

- **`add [module/module-version] [/path/to/source-tree] [/path/to/tarball.tar]`**

添加一种 `module/module-version` 组合到 dkms tree 中用于构建和安装；如果使用了 `module/module-version`, `-m module/module-version`, 或者 `-m module -v version` 作为了命令行选项，则该命令会使用 `/usr/src/<module>-<module-version>/` 目录下到源码文件，以及一个格式正确的 `dkms.conf` 文件；如果 `/path/to/source-tree` 作为选项被指定了，并且 source-tree 中包含了一个 `dkms.conf` 文件，那么该命令会将整个 `/path/to/source-tree` 拷贝为 `/usr/src/module-module-version` ；如果指定的是 `/path/to/tarball.tar` ，则该命令的行为同 `ldtarball` 命令；

- **`remove [module/module-version] [-k kernel/arch] [--all]`**

移除一种 `module/version` 或 `module/version/kernel/arch` 组合出 dkms tree ；如果目标模块当前已安装，则要先卸载安装，若模块已应用，则使用 `original_module` 中的内容进行替代；可以使用 `--all` 选项一次性移除针对每种内核的所有实例；

- **`build [module/module-version] [-k kernel/arch]`**

针对指定 `kernel/arch`构建指定的 `module/version` 组合；如果没有指定 `-k` 选项，则针对当前运行的 kernel 和 arch 进行构建；所有构建内容都位于目录 `/var/lib/dkms/<module>/<module-version>/build/` 中；如果 `module/module-version` 组合之前未被添加（到 dkms 中），则 dkms 将尝试先进行添加，此时 build 命令将接受部分 add 命令的参数；

- **`install [module/module-version] [-k kernel/arch] [/path/to/driver.rpm]`**

安装构建好的 `module/version` 组合到目标 kernel 中；如果没有通过 -k 指定 kernel 选项，则假定为当前运行 kernel 版本；如果 `module` 未被构建好，则 dkms 将尝试先进行 build ；如果 `module` 之前未被添加（到 dkms 中），则 dkms 将尝试先进行添加；在上述两种情况下，install 命令都将接受部分 build 或 add 命令的参数；如果你传入了一个 `.rpm` 文件，则 dkms 会尝试使用 `rpm -Uvh` 命令进行安装，并且会执行一个 `autoinstall` action 以便在 RPM 成功安装后，确认（mesure）你的 kernel 所需的任何东东都已构建好；
...

- **`mkrpm [module/module-version] [-k kernel/arch] [--source-only] [--binaries-only]`**

该 action 允许你针对指定的 `module/version` 创建一个 RPM package ；其使用了一个名为 `/etc/dkms/template-dkms-mkrpm.spec` 的 `.spec` 模板文件作为构建 RPM 的基础；除此之外，如果 DKMS 发现存在名为 `/usr/src/<module>-<module-version>/<module>-dkms-mkrpm.spec` 的文件，则会使用该 `.spec` 文件；通常来讲，如果一个 DKMS tarball 作为内容保存在 RPM 中，则可以基于 RPM 自身调用各种 DKMS 命令，在端用户系统上加载该 tarball ，build 和 install 相应的模块；如果你不想令 RPM 中包含任何预先构建好的二进制文件，请确保调用 mkrpm 命令时指定 `--source-only` 选项；

- **`mkdeb [module/module-version] [-k kernel/arch] [--binaries-only] [--source-only]`**

该 action 允许你针对指定的 `module/version` 创建一个 debian 二进制 package ；其使用一个名为 `/etc/dkms/template-dkms-mkdeb` 的保存了 debian 模板的目录作为 package 构建的基础；除此之外，如果 DKMS 发现存在名为 `/usr/src/<module>-<module-version>/<module>-dkms-mkdeb` 的文件夹，则会使用该文件夹替代；通常来讲，如果一个 DKMS tarball 作为内容放入 package 中，则可以通过 package 自身调用各种 DKMS 命令，在端用户系统上加载该 tarball ，build 和 install 相应的模块；如果你不想令 debian package 中包含任何预构建好的二进制文件，请确保调用 mkdeb 命令时指定 `--source-only` 选项；


### OPTIONS

- **`-m <module>/<module-version>`**

指定特定 action 的目标模块名字和模块版本；`-m` 选项时可选的；

- **`-v <module-version>`**

指定特定 action 的目标模块版本；该选项仅在基于 `-m` 选项指定模块名但未指定 `<module-version>` 时需要使用；

...

- **`--force`**

该选项和 `ldtarball` 一起使用，用于强制拷贝忽略已存在文件；

...

- **`--source-only`**

该选项用于结合 `mktarball` 或 `mkrpm` 或 `mkdeb` 命令创建 DKMS tarball ，确保不包含任何预先构建的内核模块二进制文件；该选项的用处在于：你可以简单的针对源码进行归档（tar），而不包含任何预构建的内容在其中；同样的，如果你在使用 `mkrpm` 命令时，不希望其中包含任何预构建模块，使用递该选项可以确保 RPM 内部的 DKMS tarball 不会包含任何预构建的模块；

- **`--all`**

该选项用于为 `module/module-version` 自动指定全部相关的 `kernels/arches` ；该选项对于诸如 `remove` ，`mktarball` 等命令非常有用；该命令的使用避免了针对每一种 kernel 构建相应模块时，均进行单独指定的麻烦，如 `-k kernel1 -a arch1 -k kernel2 -a arch2` ；


----------


## depmod 说明

depmod - Generate modules.dep and map files.

- Linux 内核模块能够为其它模块提供称作 "symbols" 的服务（services）供其使用（在代码中使用 EXPORT_SYMBOL 变体的一种）；
- 如果存在第二个模块同样使用到了该 symbol ，则第二个模块则明确依赖第一个模块；这些依赖很可能非常复杂；
- depmod 通过读取 `/lib/modules/version` 目录下的每一个模块信息，得到 symbols 导出和导入的具体情况，进而创建出模块依赖关系列表；默认情况下，该列表被写入 `modules.dep` 文件中，同时保存一个二进制哈希版本的 `modules.dep.bin` 文件在同一个目录中；如果在命令行中指定了文件名，则只有对应名字的模块会被处理（这种用法几乎没用，除非列出了所有模块）；
- depmod 还会创建一个由 modules 得到的 symbols 列表，保存到 `modules.symbols` 文件中，同时创建一个二进制哈希版本 `modules.symbols.bin` ；
- 最后，若模块中提供了指定的 device 名（devname），那么 depmod 将会输出名为 `modules.devname` 的文件；
- 如果提供了版本信息，则相应内核版本的模块目录被使用，而不是当前内核版本（即通过 `uname -r` 得到的版本信息）；


----------


## rpmbuild 说明

rpmbuild - Build RPM Package(s)

- rpmbuild 用于构建二进制和源码软件包；
- 生成的 RPM package 由 archive files 和用于 install 和 erase 这些 files 的 meta-data 构成；
- meta-data 包含帮助脚本，文件属性，关于当前包的描述性信息；
- Packages 有两种变体：
    - 二进制 packages ，用于封装待安装的软件；
    - 源码 packages ，包含源码以及用于生成二进制 packages 的 recipe ；
