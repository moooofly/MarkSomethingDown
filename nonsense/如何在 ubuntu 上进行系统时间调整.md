

本文简单介绍如何在 Ubuntu 上进行系统时间调整；

----------


## 查看当前时间显示（变更前）
```shell
# date
Wed Jul 27 07:52:32 GMT 2016
```

## 进行时区调整

例如将时区设置成 `Asia/Shanghai` ；

### 确定目标时区

```shell
# tzselect
Please identify a location so that time zone rules can be set correctly.
Please select a continent, ocean, "coord", or "TZ".
 1) Africa
 2) Americas
 3) Antarctica
 4) Arctic Ocean
 5) Asia
 6) Atlantic Ocean
 7) Australia
 8) Europe
 9) Indian Ocean
10) Pacific Ocean
11) coord - I want to use geographical coordinates.
12) TZ - I want to specify the time zone using the Posix TZ format.
#? 5
Please select a country whose clocks agree with yours.
 1) Afghanistan		  18) Israel		    35) Palestine
 2) Armenia		  19) Japan		    36) Philippines
 3) Azerbaijan		  20) Jordan		    37) Qatar
 4) Bahrain		  21) Kazakhstan	    38) Russia
 5) Bangladesh		  22) Korea (North)	    39) Saudi Arabia
 6) Bhutan		  23) Korea (South)	    40) Singapore
 7) Brunei		  24) Kuwait		    41) Sri Lanka
 8) Cambodia		  25) Kyrgyzstan	    42) Syria
 9) China		  26) Laos		    43) Taiwan
10) Cyprus		  27) Lebanon		    44) Tajikistan
11) East Timor		  28) Macau		    45) Thailand
12) Georgia		  29) Malaysia		    46) Turkmenistan
13) Hong Kong		  30) Mongolia		    47) United Arab Emirates
14) India		  31) Myanmar (Burma)	    48) Uzbekistan
15) Indonesia		  32) Nepal		    49) Vietnam
16) Iran		  33) Oman		    50) Yemen
17) Iraq		  34) Pakistan
#? 9
Please select one of the following time zone regions.
1) Beijing Time
2) Xinjiang Time
#? 1

The following information has been given:

	China
	Beijing Time

Therefore TZ='Asia/Shanghai' will be used.
Local time is now:	Wed Jul 27 17:25:12 CST 2016.
Universal Time is now:	Wed Jul 27 09:25:12 UTC 2016.
Is the above information OK?
1) Yes
2) No
#? 1

You can make this change permanent for yourself by appending the line
	TZ='Asia/Shanghai'; export TZ
to the file '.profile' in your home directory; then log out and log in again.

Here is that TZ value again, this time on standard output so that you
can use the /usr/bin/tzselect command in shell scripts:
Asia/Shanghai
#
```


### 基于图形界面进行时区调整

只需要在界面上进行简单选择即可；

```shell
# dpkg-reconfigure tzdata

Current default time zone: 'Asia/Shanghai'
Local time is now:      Wed Jul 27 17:04:51 CST 2016.
Universal Time is now:  Wed Jul 27 09:04:51 UTC 2016.

#
```

### 基于 timedatectl 进行时区调整

如果你使用的 Linux 系统支持 `Systemd`，则可以使用 `timedatectl` 命令进行系统范围时区设置。

在 Systemd 下有一个名为 `systemd-timedated` 的系统服务负责调整系统时钟和时区，可以使用 timedatectl 命令对此系统服务进行配置。

设置时区
```shell
timedatectl set-timezone 'Asia/Shanghai'
```

查看当前时间设定
```shell
# timedatectl
Warning: Ignoring the TZ variable. Reading the system's time zone setting only.

      Local time: Wed 2016-07-27 17:45:54 CST
  Universal time: Wed 2016-07-27 09:45:54 UTC
        RTC time: Wed 2016-07-27 09:45:54
       Time zone: Asia/Shanghai (CST, +0800)
     NTP enabled: yes
NTP synchronized: yes
 RTC in local TZ: no
      DST active: n/a
#
```

## 防止系统重启后时区改变

### 用户级别设置时区

在 `.profile` 文件中添加
```shell
TZ='Asia/Shanghai'; export TZ
```
之后执行 `source .profile` 令设置生效；也可以重新进行 shell 登录；


### 系统级别设置时区

```shell
rm -f /etc/localtime
ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```
或者（等效命令）
```shell
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```
> 注意：`Asia/Shanghai` 是根据上面 `tzselect` 确定的；


## 通过公网 NTP 服务器进行时间校准

```shell
# ntpdate pool.ntp.org
27 Jul 15:56:36 ntpdate[8265]: adjust time server 115.28.122.198 offset 0.010742 sec
```

## 查看当前时间显示（变更后）

```shell
# date
Wed Jul 27 15:58:03 CST 2016
```

## 设置硬件时间和系统时间一致

```shell
# /sbin/hwclock --systohc
```
