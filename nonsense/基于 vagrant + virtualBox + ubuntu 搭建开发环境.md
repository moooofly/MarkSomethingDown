# 基于 vagrant + virtualBox + ubuntu 搭建开发环境

## 软件安装

手动安装

- VirtualBox 下载地址：[https://www.virtualbox.org/wiki/Downloads](https://www.virtualbox.org/wiki/Downloads)
- vagrant 下载地址：[https://www.vagrantup.com/downloads.html](https://www.vagrantup.com/downloads.html)

自动安装

```
brew cask install virtualbox
brew cask install vagrant
```

> 或使用 reinstall 命令


.box 文件下载：

- [vagrantbox](http://www.vagrantbox.es/)
- [Ubuntu Cloud Images (RELEASED)](https://cloud-images.ubuntu.com/releases/)


## 环境准备

创建目录

```shell
mkdir -p /Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64
```

将之前下载好的 .box 文件放到上述开发目录中

```shell
➜ ubuntu-14.04-amd64 cp ~/Downloads/ubuntu-14.04-amd64.box ./
➜ ubuntu-14.04-amd64 ll
total 785608
-rw-r----- 1 sunfei staff 384M 7 3 12:46 ubuntu-14.04-amd64.box
➜ ubuntu-14.04-amd64
```

# 添加 box 到 vagrant

```shell
➜ ubuntu-14.04-amd64 vagrant box add "ubuntu-14.04-amd64" ubuntu-14.04-amd64.box
==> box: Box file was not detected as metadata. Adding it directly...
==> box: Adding box 'ubuntu-14.04-amd64' (v0) for provider:
box: Unpacking necessary files from: file:///Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64/ubuntu-14.04-amd64.box
==> box: Successfully added box 'ubuntu-14.04-amd64' (v0) for 'virtualbox'!
```

成功添加后，可以通过 `vagrant box list` 查看；


# 初始化

```shell
➜ ubuntu-14.04-amd64 vagrant init "ubuntu-14.04-amd64"
A `Vagrantfile` has been placed in this directory. You are now
ready to `vagrant up` your first virtual environment! Please read
the comments in the Vagrantfile as well as documentation on
`vagrantup.com` for more information on using Vagrant.
➜ ubuntu-14.04-amd64
➜ ubuntu-14.04-amd64 ll
total 785616
-rw-r--r-- 1 sunfei staff 3.0K 7 3 12:56 Vagrantfile
-rw-r----- 1 sunfei staff 384M 7 3 12:46 ubuntu-14.04-amd64.box
➜ ubuntu-14.04-amd64
```

初始化后，会生成默认配置文件，内容如下

```shell
# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure(2) do |config|
# The most common configuration options are documented and commented below.
# For a complete reference, please see the online documentation at
# https://docs.vagrantup.com.

# Every Vagrant development environment requires a box. You can search for
# boxes at https://atlas.hashicorp.com/search.
config.vm.box = "ubuntu-14.04-amd64"

# Disable automatic box update checking. If you disable this, then
# boxes will only be checked for updates when the user runs
# `vagrant box outdated`. This is not recommended.
# config.vm.box_check_update = false

# Create a forwarded port mapping which allows access to a specific port
# within the machine from a port on the host machine. In the example below,
# accessing "localhost:8080" will access port 80 on the guest machine.
# config.vm.network "forwarded_port", guest: 80, host: 8080

# Create a private network, which allows host-only access to the machine
# using a specific IP.
# config.vm.network "private_network", ip: "192.168.33.10"

# Create a public network, which generally matched to bridged network.
# Bridged networks make the machine appear as another physical device on
# your network.
# config.vm.network "public_network"

# Share an additional folder to the guest VM. The first argument is
# the path on the host to the actual folder. The second argument is
# the path on the guest to mount the folder. And the optional third
# argument is a set of non-required options.
# config.vm.synced_folder "../data", "/vagrant_data"

# Provider-specific configuration so you can fine-tune various
# backing providers for Vagrant. These expose provider-specific options.
# Example for VirtualBox:
#
# config.vm.provider "virtualbox" do |vb|
# # Display the VirtualBox GUI when booting the machine
# vb.gui = true
#
# # Customize the amount of memory on the VM:
# vb.memory = "1024"
# end
#
# View the documentation for the provider you are using for more
# information on available options.

# Define a Vagrant Push strategy for pushing to Atlas. Other push strategies
# such as FTP and Heroku are also available. See the documentation at
# https://docs.vagrantup.com/v2/push/atlas.html for more information.
# config.push.define "atlas" do |push|
# push.app = "YOUR_ATLAS_USERNAME/YOUR_APPLICATION_NAME"
# end

# Enable provisioning with a shell script. Additional provisioners such as
# Puppet, Chef, Ansible, Salt, and Docker are also available. Please see the
# documentation for more information about their specific syntax and use.
# config.vm.provision "shell", inline: <<-SHELL
# sudo apt-get update
# sudo apt-get install -y apache2
# SHELL
end
```

# 启动虚拟机

## 未调整 Vagrantfile 文件内容前

```shell
➜ ubuntu-14.04-amd64 vagrant up
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Importing base box 'ubuntu-14.04-amd64'...
==> default: Matching MAC address for NAT networking...
==> default: Setting the name of the VM: ubuntu-1404-amd64_default_1467530807542_36039
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
default: Adapter 1: nat
==> default: Forwarding ports...
default: 22 (guest) => 2222 (host) (adapter 1)
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
default: SSH address: 127.0.0.1:2222
default: SSH username: vagrant
default: SSH auth method: private key
default: Warning: Remote connection disconnect. Retrying...
default:
default: Vagrant insecure key detected. Vagrant will automatically replace
default: this with a newly generated keypair for better security.
default:
default: Inserting generated public key within guest...
default: Removing insecure key from the guest if it's present...
default: Key inserted! Disconnecting and reconnecting using new SSH key...
==> default: Machine booted and ready!
==> default: Checking for guest additions in VM...
default: No guest additions were detected on the base box for this VM! Guest
default: additions are required for forwarded ports, shared folders, host only
default: networking, and more. If SSH fails on this machine, please install
default: the guest additions and repackage the box to continue.
default:
default: This is not an error message; everything may continue to work properly,
default: in which case you may ignore this message.
==> default: Mounting shared folders...
default: /vagrant => /Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64
Failed to mount folders in Linux guest. This is usually because
the "vboxsf" file system is not available. Please verify that
the guest additions are properly installed in the guest and
can work properly. The command attempted was:

mount -t vboxsf -o uid=`id -u vagrant`,gid=`getent group vagrant | cut -d: -f3` vagrant /vagrant
mount -t vboxsf -o uid=`id -u vagrant`,gid=`id -g vagrant` vagrant /vagrant

The error output from the last command was:

stdin: is not a tty
mount: unknown filesystem type 'vboxsf'

➜ ubuntu-14.04-amd64
```

此时进入到虚拟机后，可以看到没有名为 `/vagrant` 的数据共享目录；

```shell
➜ ubuntu-14.04-amd64 vagrant ssh
Welcome to Ubuntu 14.04 LTS (GNU/Linux 3.13.0-24-generic x86_64)

* Documentation: https://help.ubuntu.com/
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ mount
/dev/sda1 on / type ext4 (rw,errors=remount-ro)
proc on /proc type proc (rw,noexec,nosuid,nodev)
sysfs on /sys type sysfs (rw,noexec,nosuid,nodev)
none on /sys/fs/cgroup type tmpfs (rw)
none on /sys/fs/fuse/connections type fusectl (rw)
none on /sys/kernel/debug type debugfs (rw)
none on /sys/kernel/security type securityfs (rw)
udev on /dev type devtmpfs (rw,mode=0755)
devpts on /dev/pts type devpts (rw,noexec,nosuid,gid=5,mode=0620)
tmpfs on /run type tmpfs (rw,noexec,nosuid,size=10%,mode=0755)
none on /run/lock type tmpfs (rw,noexec,nosuid,nodev,size=5242880)
none on /run/shm type tmpfs (rw,nosuid,nodev)
none on /run/user type tmpfs (rw,noexec,nosuid,nodev,size=104857600,mode=0755)
none on /sys/fs/pstore type pstore (rw)
rpc_pipefs on /run/rpc_pipefs type rpc_pipefs (rw)
systemd on /sys/fs/cgroup/systemd type cgroup (rw,noexec,nosuid,nodev,none,name=systemd)
vagrant@vagrant-ubuntu-trusty:~$
```

虚拟机的网络配置情况如下

```shell
vagrant@vagrant-ubuntu-trusty:~$ ifconfig
eth0 Link encap:Ethernet HWaddr 08:00:27:49:78:3b
inet addr:10.0.2.15 Bcast:10.0.2.255 Mask:255.255.255.0
inet6 addr: fe80::a00:27ff:fe49:783b/64 Scope:Link
UP BROADCAST RUNNING MULTICAST MTU:1500 Metric:1
RX packets:239 errors:0 dropped:0 overruns:0 frame:0
TX packets:170 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:1000
RX bytes:27080 (27.0 KB) TX bytes:23782 (23.7 KB)

lo Link encap:Local Loopback
inet addr:127.0.0.1 Mask:255.0.0.0
inet6 addr: ::1/128 Scope:Host
UP LOOPBACK RUNNING MTU:65536 Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:0
RX bytes:0 (0.0 B) TX bytes:0 (0.0 B)

vagrant@vagrant-ubuntu-trusty:~$
```

> 注意：上面的共享目录问题和 guest additions 有关，请参考《[Mac 和 VirtualBox 之间的文件共享](https://github.com/moooofly/MarkSomethingDown/blob/master/Mac%20%E5%92%8C%20VirtualBox%20%E4%B9%8B%E9%97%B4%E7%9A%84%E6%96%87%E4%BB%B6%E5%85%B1%E4%BA%AB.md)》中方法进行解决；


在解决了 guest additions 问题之后，我们可以按照需要调整 Vagrantfile 文件内容如下：

```shell
➜ ubuntu-14.04-amd64 diff Vagrantfile Vagrantfile_new
29c29
< # config.vm.network "private_network", ip: "192.168.33.10"
---
> config.vm.network "private_network", ip: "11.11.11.12"
40c40
< # config.vm.synced_folder "../data", "/vagrant_data"
---
> config.vm.synced_folder "../data", "/vagrant_data"
➜ ubuntu-14.04-amd64
```

重新启动虚拟机
```shell
➜ ubuntu-14.04-amd64 vagrant up
Bringing machine 'default' up with 'virtualbox' provider...
There are errors in the configuration of this machine. Please fix
the following errors and try again:

vm:
* The host path of the shared folder is missing: ../data

➜ ubuntu-14.04-amd64
```
**错误原因**：需要事先创建好 ../data 目录；

创建好 data 目录后，再次重新启动
```shell
➜ ubuntu-14.04-amd64 vagrant reload
==> default: Attempting graceful shutdown of VM...
==> default: Clearing any previously set forwarded ports...
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
default: Adapter 1: nat
default: Adapter 2: hostonly
==> default: Forwarding ports...
default: 22 (guest) => 2222 (host) (adapter 1)
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
default: SSH address: 127.0.0.1:2222
default: SSH username: vagrant
default: SSH auth method: private key
default: Warning: Remote connection disconnect. Retrying...
==> default: Machine booted and ready!
==> default: Checking for guest additions in VM...
==> default: Configuring and enabling network interfaces...
==> default: Mounting shared folders...
default: /vagrant => /Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64
default: /vagrant_data => /Users/sunfei/workspace/vagrant/data
==> default: Machine already provisioned. Run `vagrant provision` or use the `--provision`
==> default: flag to force provisioning. Provisioners marked to run always will still run.
➜ ubuntu-14.04-amd64
```
此时，共享目录问题已经正确了～

此时进入虚拟机，可以看到数据共享目录和网络配置均已经 ok ；

```shell
➜ ubuntu-14.04-amd64 vagrant ssh
Welcome to Ubuntu 14.04 LTS (GNU/Linux 3.13.0-24-generic x86_64)

* Documentation: https://help.ubuntu.com/
Last login: Sun Jul 3 08:01:09 2016 from 10.0.2.2
vagrant@vagrant-ubuntu-trusty:~$
vagrant@vagrant-ubuntu-trusty:~$ mount
/dev/sda1 on / type ext4 (rw,errors=remount-ro)
proc on /proc type proc (rw,noexec,nosuid,nodev)
sysfs on /sys type sysfs (rw,noexec,nosuid,nodev)
none on /sys/fs/cgroup type tmpfs (rw)
none on /sys/fs/fuse/connections type fusectl (rw)
none on /sys/kernel/debug type debugfs (rw)
none on /sys/kernel/security type securityfs (rw)
udev on /dev type devtmpfs (rw,mode=0755)
devpts on /dev/pts type devpts (rw,noexec,nosuid,gid=5,mode=0620)
tmpfs on /run type tmpfs (rw,noexec,nosuid,size=10%,mode=0755)
none on /run/lock type tmpfs (rw,noexec,nosuid,nodev,size=5242880)
none on /run/shm type tmpfs (rw,nosuid,nodev)
none on /run/user type tmpfs (rw,noexec,nosuid,nodev,size=104857600,mode=0755)
none on /sys/fs/pstore type pstore (rw)
systemd on /sys/fs/cgroup/systemd type cgroup (rw,noexec,nosuid,nodev,none,name=systemd)
rpc_pipefs on /run/rpc_pipefs type rpc_pipefs (rw)
vagrant on /vagrant type vboxsf (uid=1000,gid=1000,rw)
vagrant_data on /vagrant_data type vboxsf (uid=1000,gid=1000,rw)
vagrant@vagrant-ubuntu-trusty:~$
```
```shell
vagrant@vagrant-ubuntu-trusty:~$ ifconfig
eth0 Link encap:Ethernet HWaddr 08:00:27:49:78:3b
inet addr:10.0.2.15 Bcast:10.0.2.255 Mask:255.255.255.0
inet6 addr: fe80::a00:27ff:fe49:783b/64 Scope:Link
UP BROADCAST RUNNING MULTICAST MTU:1500 Metric:1
RX packets:491 errors:0 dropped:0 overruns:0 frame:0
TX packets:310 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:1000
RX bytes:52052 (52.0 KB) TX bytes:42254 (42.2 KB)

eth1 Link encap:Ethernet HWaddr 08:00:27:f7:de:77
inet addr:11.11.11.12 Bcast:11.11.11.255 Mask:255.255.255.0
inet6 addr: fe80::a00:27ff:fef7:de77/64 Scope:Link
UP BROADCAST RUNNING MULTICAST MTU:1500 Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:8 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:1000
RX bytes:0 (0.0 B) TX bytes:648 (648.0 B)

lo Link encap:Local Loopback
inet addr:127.0.0.1 Mask:255.0.0.0
inet6 addr: ::1/128 Scope:Host
UP LOOPBACK RUNNING MTU:65536 Metric:1
RX packets:0 errors:0 dropped:0 overruns:0 frame:0
TX packets:0 errors:0 dropped:0 overruns:0 carrier:0
collisions:0 txqueuelen:0
RX bytes:0 (0.0 B) TX bytes:0 (0.0 B)

vagrant@vagrant-ubuntu-trusty:~$
```

## 问题

### mac 电量不足导致 vagrant 直接退出后

```
➜  ubuntu-16.04-server-cloudimg-amd64-vagrant vup
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Resuming suspended VM...
==> default: Booting VM...
There was an error while executing `VBoxManage`, a CLI used by Vagrant
for controlling VirtualBox. The command and stderr is shown below.

Command: ["startvm", "a1dd6670-8ee3-4f9a-bb54-f10742d82ce3", "--type", "headless"]

Stderr: VBoxManage: error: Failed to load unit 'lsilogicscsi' (VERR_SSM_LOADED_TOO_LITTLE)
VBoxManage: error: Details: code NS_ERROR_FAILURE (0x80004005), component ConsoleWrap, interface IConsole

➜  ubuntu-16.04-server-cloudimg-amd64-vagrant
```

ref: https://github.com/hashicorp/vagrant/issues/1809

```
vagrant reload
```


### macOS 升级到 Mojave 后出现无法启动情况

```
➜  ubuntu-16.04-server-cloudimg-amd64-vagrant vup
Bringing machine 'default' up with 'virtualbox' provider...
==> default: Clearing any previously set network interfaces...
There was an error while executing `VBoxManage`, a CLI used by Vagrant
for controlling VirtualBox. The command and stderr is shown below.

Command: ["hostonlyif", "create"]

Stderr: 0%...
Progress state: NS_ERROR_FAILURE
VBoxManage: error: Failed to create the host-only adapter
VBoxManage: error: VBoxNetAdpCtl: Error while adding new interface: failed to open /dev/vboxnetctl: No such file or directory
VBoxManage: error: Details: code NS_ERROR_FAILURE (0x80004005), component HostNetworkInterfaceWrap, interface IHostNetworkInterface
VBoxManage: error: Context: "RTEXITCODE handleCreate(HandlerArg *)" at line 94 of file VBoxManageHostonly.cpp

➜  ubuntu-16.04-server-cloudimg-amd64-vagrant
```

ref: https://github.com/hashicorp/vagrant/issues/1671#issuecomment-424657289

```
brew cask reinstall virtualbox
brew cask reinstall vagrant
vagrant plugin update
```
