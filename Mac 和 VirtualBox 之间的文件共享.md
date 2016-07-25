


本文针对 Mac 和 VirtualBox 之间的文件共享进行说明；

众所周知，跨系统文件共享问题有几种解决办法：

- 通过 samba 协议解决
- 通过 Web 服务解决
- 通过 ftp 类协议解决

本文针对另外一种方法进行说明，即 `Guest Additions` ；

以下内容部分参考自：《[How to install VirtualBox Guest Additions for Linux](http://www.tuicool.com/articles/U3U73u)》

------

> 下面内容假设你已经基于 vagrant ＋ VirtualBox ＋ Ubuntu 成功搭建起了虚拟机系统；

# 下载对应 VirtualBox 版本的 guest additions

```shell
vagrant@vagrant-ubuntu-trusty:~/workspace/WGET$ wget http://download.virtualbox.org/virtualbox/5.0.20/VBoxGuestAdditions_5.0.20.iso
```

# 安装必要的软件包

```shell
vagrant@vagrant-ubuntu-trusty:~/workspace/WGET$ sudo apt-get install dkms gcc
```

# 挂载 guest additions 对应的 ISO 文件

```shell
vagrant@vagrant-ubuntu-trusty:~/workspace/WGET$ sudo mount -o loop VBoxGuestAdditions_5.0.20.iso /mnt
mount: block device /home/vagrant/workspace/WGET/VBoxGuestAdditions_5.0.20.iso is write-protected, mounting read-only
```

查看挂载情况

```shell
vagrant@vagrant-ubuntu-trusty:~/workspace/WGET$ mount
…
/home/vagrant/workspace/WGET/VBoxGuestAdditions_5.0.20.iso on /mnt type iso9660 (ro)
vagrant@vagrant-ubuntu-trusty:~/workspace/WGET$ 
```

# 运行脚本进行安装

```shell
vagrant@vagrant-ubuntu-trusty:~/workspace/WGET$ cd /mnt
vagrant@vagrant-ubuntu-trusty:/mnt$ sudo ./VBoxLinuxAdditions.run 
Verifying archive integrity... All good.
Uncompressing VirtualBox 5.0.20 Guest Additions for Linux............
VirtualBox Guest Additions installer
Copying additional installer modules ...
Installing additional modules ...
Removing existing VirtualBox DKMS kernel modules ...done.
Removing existing VirtualBox non-DKMS kernel modules ...done.
Building the VirtualBox Guest Additions kernel modules ...done.
Doing non-kernel setup of the Guest Additions ...done.
Starting the VirtualBox Guest AdditionsInstalling the Window System drivers
Could not find the X.Org or XFree86 Window System, skipping.

 ...done.

vagrant@vagrant-ubuntu-trusty:/mnt$ 
```

# 查看安装成功后，内核模块中增加的和 vbox 相关的内容

```shell
vagrant@vagrant-ubuntu-trusty:/mnt$ lsmod | grep vbox
vboxvideo              45696  1 
ttm                    85115  1 vboxvideo
drm_kms_helper         52758  1 vboxvideo
drm                   302817  3 ttm,drm_kms_helper,vboxvideo
syscopyarea            12529  1 vboxvideo
sysfillrect            12701  1 vboxvideo
sysimgblt              12640  1 vboxvideo
vboxsf                 43802  0 
vboxguest             276728  3 vboxsf,vboxvideo

vagrant@vagrant-ubuntu-trusty:/mnt$ 
```

------

上面给出的一切顺利的情况下，你能够看到的安装过程～～

然而，实际操作中你可能会遇以下情况：

- 针对 lucid64.box ，即使 guest additions 版本不一致也可以提供文件共享功能

```shell
sunfeideMacBook-Pro:lucid64 sunfei$ vagrant reload

==> default: Attempting graceful shutdown of VM...
==> default: Clearing any previously set forwarded ports...
==> default: Fixed port collision for 22 => 2222. Now on port 2200.
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
    default: Adapter 1: nat
    default: Adapter 2: hostonly
==> default: Forwarding ports...
    default: 80 (guest) => 8080 (host) (adapter 1)
    default: 22 (guest) => 2200 (host) (adapter 1)
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
    default: SSH address: 127.0.0.1:2200
    default: SSH username: vagrant
    default: SSH auth method: private key
    default: Warning: Remote connection disconnect. Retrying...
==> default: Machine booted and ready!
==> default: Checking for guest additions in VM...
    default: The guest additions on this VM do not match the installed version of    －－ 这里可以看到存在版本不一致问题
    default: VirtualBox! In most cases this is fine, but in rare cases it can
    default: prevent things such as shared folders from working properly. If you see
    default: shared folder errors, please make sure the guest additions within the
    default: virtual machine match the version of VirtualBox you have installed on
    default: your host and reload your VM.
    default: 
    default: Guest Additions Version: 4.2.0
    default: VirtualBox Version: 5.0
==> default: Configuring and enabling network interfaces...
==> default: Mounting shared folders...
    default: /vagrant => /Users/sunfei/workspace/vagrant/lucid64                －－ 但依然能够成功挂载共享目录
    default: /vagrant_data => /Users/sunfei/workspace/vagrant/lucid64/data      －－ 但依然能够成功挂载共享目录
==> default: Machine already provisioned. Run `vagrant provision` or use the `--provision`
==> default: flag to force provisioning. Provisioners marked to run always will still run.

sunfeideMacBook-Pro:lucid64 sunfei$ 
```

- 针对 ubuntu-14.04-amd64.box 的情况，则必须安装正确版本的 guest additions 才能共享文件

```shell
sunfeideMacBook-Pro:ubuntu-14.04-amd64 sunfei$ vagrant up

Bringing machine 'default' up with 'virtualbox' provider...
==> default: Importing base box 'ubuntu-14.04-amd64'...
==> default: Matching MAC address for NAT networking...
==> default: Setting the name of the VM: ubuntu-1404-amd64_default_1465201090522_12444
==> default: Fixed port collision for 22 => 2222. Now on port 2200.
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
    default: Adapter 1: nat
    default: Adapter 2: hostonly
==> default: Forwarding ports...
    default: 22 (guest) => 2200 (host) (adapter 1)
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
    default: SSH address: 127.0.0.1:2200
    default: SSH username: vagrant
    default: SSH auth method: private key
    default: Warning: Remote connection disconnect. Retrying...
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
    default: No guest additions were detected on the base box for this VM! Guest  － 这里看到当前 box 未提供 guest additions
    default: additions are required for forwarded ports, shared folders, host only
    default: networking, and more. If SSH fails on this machine, please install
    default: the guest additions and repackage the box to continue.
    default: 
    default: This is not an error message; everything may continue to work properly,
    default: in which case you may ignore this message.
==> default: Configuring and enabling network interfaces...
==> default: Mounting shared folders...
    default: /vagrant => /Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64
Failed to mount folders in Linux guest. This is usually because       － 这里会看到共享文件夹挂载失败
the "vboxsf" file system is not available. Please verify that
the guest additions are properly installed in the guest and
can work properly. The command attempted was:

mount -t vboxsf -o uid=`id -u vagrant`,gid=`getent group vagrant | cut -d: -f3` vagrant /vagrant
mount -t vboxsf -o uid=`id -u vagrant`,gid=`id -g vagrant` vagrant /vagrant

The error output from the last command was:
stdin: is not a tty
mount: unknown filesystem type 'vboxsf'

sunfeideMacBook-Pro:ubuntu-14.04-amd64 sunfei$ 
```

- 针对 ubuntu-14.04-amd64.box 安装完 guest additions 的情况

```shell
sunfeideMacBook-Pro:ubuntu-14.04-amd64 sunfei$ vagrant reload

==> default: Attempting graceful shutdown of VM...
==> default: Clearing any previously set forwarded ports...
==> default: Fixed port collision for 22 => 2222. Now on port 2200.
==> default: Clearing any previously set network interfaces...
==> default: Preparing network interfaces based on configuration...
    default: Adapter 1: nat
    default: Adapter 2: hostonly
==> default: Forwarding ports...
    default: 22 (guest) => 2200 (host) (adapter 1)
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
    default: SSH address: 127.0.0.1:2200
    default: SSH username: vagrant
    default: SSH auth method: private key
    default: Warning: Remote connection disconnect. Retrying...
==> default: Machine booted and ready!
==> default: Checking for guest additions in VM...                  — 安装后，这里不再出现其他信息
==> default: Configuring and enabling network interfaces...
==> default: Mounting shared folders...

    default: /vagrant => /Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64   － 根据配置文件成功挂载共享目录
    default: /vagrant_data => /Users/sunfei/workspace/vagrant/ubuntu-14.04-amd64/data

sunfeideMacBook-Pro:ubuntu-14.04-amd64 sunfei$ 
```

此时，登录到虚拟机系统中可以看到

```shell
sunfeideMacBook-Pro:ubuntu-14.04-amd64 sunfei$ vagrant ssh
Welcome to Ubuntu 14.04 LTS (GNU/Linux 3.13.0-24-generic x86_64)

 * Documentation:  https://help.ubuntu.com/

Last login: Mon Jun  6 09:05:32 2016 from 10.0.2.2
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
vagrant on /vagrant type vboxsf (uid=1000,gid=1000,rw)      — 共享目录成功挂载
vagrant_data on /vagrant_data type vboxsf (uid=1000,gid=1000,rw)     — 共享目录成功挂载

vagrant@vagrant-ubuntu-trusty:~$ 
```

> 补充：Guest Additions 除了能够解决文件共享问题外，还能解决一些其它宿主机和虚拟机之间的问题，详细内容可以查阅相关资料～

