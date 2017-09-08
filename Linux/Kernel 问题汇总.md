# Kernel 问题汇总

- "unable to handle kernel NULL pointer dereference at 0000000000000010"
- "kernel: EPT: Misconfiguration"


----------


## "unable to handle kernel NULL pointer dereference at 0000000000000010"

### 故障信息

vmcore-dmesg.txt 信息

```
...
[47113590.457363] BUG: unable to handle kernel NULL pointer dereference at 0000000000000010
[47113590.458030] IP: [<ffffffff812db1e1>] rb_next+0x1/0x50
[47113590.458377] PGD 1b8f60067 PUD 460d95067 PMD 0
[47113590.458723] Oops: 0000 [#1] SMP
[47113590.459060] Modules linked in: iptable_filter tcp_diag inet_diag binfmt_misc bridge stp llc ipmi_si ipmi_devintf ipmi_msghandler bonding dm_mirror dm_region_hash dm_log dm_mod coretemp intel_rapl kvm_intel kvm crct10dif_pclmul crc32_pclmul crc32c_intel ghash_clmulni_intel aesni_intel lrw gf128mul glue_helper iTCO_wdt ablk_helper cryptd iTCO_vendor_support sg pcspkr lpc_ich mei_me shpchp mei sb_edac mfd_core edac_core i2c_i801 acpi_power_meter ip_tables xfs libcrc32c sd_mod crc_t10dif crct10dif_common igb ahci libahci ptp pps_core libata i2c_algo_bit i2c_core megaraid_sas dca
[47113590.461921] CPU: 3 PID: 16179 Comm: java Not tainted 3.10.0-229.el7.x86_64 #1
[47113590.462585] Hardware name: Huawei RH1288 V3/BC11HGSC0, BIOS 1.71 12/03/2015
[47113590.463137] task: ffff88085f0838e0 ti: ffff8803bd398000 task.ti: ffff8803bd398000
[47113590.463690] RIP: 0010:[<ffffffff812db1e1>]  [<ffffffff812db1e1>] rb_next+0x1/0x50
[47113590.464256] RSP: 0018:ffff8803bd39bc28  EFLAGS: 00010046
[47113590.464587] RAX: 0000000000000000 RBX: 0000000000000000 RCX: 0000000000000000
[47113590.465256] RDX: 0000000000000001 RSI: ffff880470073728 RDI: 0000000000000010
[47113590.465899] RBP: ffff8803bd39bc70 R08: 0000000000000000 R09: 0000000000000001
[47113590.466540] R10: 0000000000000001 R11: 0000000000000000 R12: ffff8804691dfe00
[47113590.467180] R13: 0000000000000000 R14: 0000000000000000 R15: 0000000000000000
[47113590.467842] FS:  00007f8f39f83700(0000) GS:ffff880470060000(0000) knlGS:0000000000000000
[47113590.468486] CS:  0010 DS: 0000 ES: 0000 CR0: 0000000080050033
[47113590.468821] CR2: 0000000000000010 CR3: 000000045539c000 CR4: 00000000001407e0
[47113590.469476] DR0: 0000000000000000 DR1: 0000000000000000 DR2: 0000000000000000
[47113590.470137] DR3: 0000000000000000 DR6: 00000000ffff0ff0 DR7: 0000000000000400
[47113590.470798] Stack:
[47113590.471117]  ffff8803bd39bc70 ffffffff810aff99 00000001bd39bc60 ffff880470073680
[47113590.471792]  ffff88085f083ec0 ffff880470073680 0000000000000003 ffff8803bd39bde0
[47113590.472467]  ffffc900135db000 ffff8803bd39bcd0 ffffffff81609782 ffff8803bd39bfd8
[47113590.473142] Call Trace:
[47113590.473473]  [<ffffffff810aff99>] ? pick_next_task_fair+0x129/0x1d0
[47113590.473820]  [<ffffffff81609782>] __schedule+0x122/0x7b0
[47113590.474159]  [<ffffffff81609e39>] schedule+0x29/0x70
[47113590.474497]  [<ffffffff810d242e>] futex_wait_queue_me+0xde/0x140
[47113590.474838]  [<ffffffff810d2fa9>] futex_wait+0x179/0x280
[47113590.475173]  [<ffffffff8109af80>] ? hrtimer_get_res+0x50/0x50
[47113590.475513]  [<ffffffff8109b914>] ? hrtimer_start_range_ns+0x14/0x20
[47113590.475848]  [<ffffffff810d508e>] do_futex+0xfe/0x5b0
[47113590.476178]  [<ffffffff810d55c0>] SyS_futex+0x80/0x180
[47113590.476517]  [<ffffffff81614a29>] system_call_fastpath+0x16/0x1b
[47113590.476842] Code: 89 06 48 8b 47 08 48 89 46 08 48 8b 47 10 48 89 46 10 c3 0f 1f 80 00 00 00 00 48 89 32 eb b2 0f 1f 00 48 89 70 10 eb a9 66 90 55 <48> 8b 17 48 89 e5 48 39 d7 74 3b 48 8b 47 08 48 85 c0 75 0e eb
[47113590.477909] RIP  [<ffffffff812db1e1>] rb_next+0x1/0x50
[47113590.478196]  RSP <ffff8803bd39bc28>
[47113590.478477] CR2: 0000000000000010
```

### 系统信息

```
[root@xg-napos-shop-service-11 ]# uname -a
Linux xg-napos-shop-service-11 3.10.0-229.el7.x86_64 #1 SMP Fri Mar 6 11:36:42 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
[root@xg-napos-shop-service-11 ]# lsb_release -a
LSB Version:	:core-4.1-amd64:core-4.1-noarch
Distributor ID:	CentOS
Description:	CentOS Linux release 7.1.1503 (Core)
Release:	7.1.1503
Codename:	Core
[root@xg-napos-shop-service-11 ~]# lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                24
On-line CPU(s) list:   0-23
Thread(s) per core:    2
Core(s) per socket:    6
Socket(s):             2
NUMA node(s):          2
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 63
Model name:            Intel(R) Xeon(R) CPU E5-2620 v3 @ 2.40GHz
Stepping:              2
CPU MHz:               2599.968
BogoMIPS:              4793.14
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              15360K
NUMA node0 CPU(s):     0-5,12-17
NUMA node1 CPU(s):     6-11,18-23
```

### 故障原因

（取自这里：[RHBA-2016:2966 - Bug Fix Advisory]()）：

> * Previously, a "NULL pointer dereference" problem in the pick_next_task_fair()
function occurred. This update fixes the bug by applying a set of patches on the
Completely Fair Scheduler (CFS) group scheduling. As a result, the "NULL pointer
dereference" problem no longer occurs. (BZ#1373820)

### 解决方案

升级内核到 kernel-3.10.0-327.44.2 以上可解决；



----------


## "kernel: EPT: Misconfiguration"

### 故障信息

在 `/var/log/messages` 或 dmesg 中可以看到

```
Mar 7 13:38:04 xg-app-zstack-118 kernel: EPT: Misconfiguration.
Mar 7 13:38:04 xg-app-zstack-118 kernel: EPT: GPA: 0x108ba4538
Mar 7 13:38:04 xg-app-zstack-118 kernel: ept_misconfig_inspect_spte: spte 0x84c236107 level 4
Mar 7 13:38:04 xg-app-zstack-118 kernel: ept_misconfig_inspect_spte: spte 0x10515c7107 level 3
Mar 7 13:38:04 xg-app-zstack-118 kernel: ept_misconfig_inspect_spte: spte 0xde4657107 level 2
Mar 7 13:38:04 xg-app-zstack-118 kernel: ept_misconfig_inspect_spte: spte 0x1028125d77 level 1
```

关键：

- **EPT**: Extended Page Table
- **GPA**: Guest Physical address

> 说明上述错误和 KVM 虚拟化有关；


可能的表现：

- Found VM in paused state
- Unable to get it running again
- 服务没有响应，目标 KVM 虚拟机无法登陆



### 系统信息

kvm 虚拟机

```
[ops@xg-account-core-srv-10 ~]$ uname -a
Linux xg-account-core-srv-10 3.10.0-229.el7.x86_64 #1 SMP Fri Mar 6 11:36:42 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
[ops@xg-account-core-srv-10 ~]$ lsb_release -a
LSB Version:	:core-4.1-amd64:core-4.1-noarch
Distributor ID:	CentOS
Description:	CentOS Linux release 7.1.1503 (Core)
Release:	7.1.1503
Codename:	Core
[root@xg-account-core-srv-10 ~]# lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                8
On-line CPU(s) list:   0-7
Thread(s) per core:    1
Core(s) per socket:    8
Socket(s):             1
NUMA node(s):          1
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 13
Model name:            QEMU Virtual CPU version 2.3.0
Stepping:              3
CPU MHz:               2599.996
BogoMIPS:              5199.99
Hypervisor vendor:     KVM
Virtualization type:   full
L1d cache:             32K
L1i cache:             32K
L2 cache:              4096K
NUMA node0 CPU(s):     0-7
```

宿主机

```
[root@xg-app-zstack-118 ~]# uname -a
Linux xg-app-zstack-118 3.10.0-229.el7.x86_64 #1 SMP Fri Mar 6 11:36:42 UTC 2015 x86_64 x86_64 x86_64 GNU/Linux
[root@xg-app-zstack-118 ~]# lsb_release -a
LSB Version:	:core-4.1-amd64:core-4.1-noarch
Distributor ID:	CentOS
Description:	CentOS Linux release 7.1.1503 (Core)
Release:	7.1.1503
Codename:	Core
[root@xg-app-zstack-118 ~]# lscpu
Architecture:          x86_64
CPU op-mode(s):        32-bit, 64-bit
Byte Order:            Little Endian
CPU(s):                32
On-line CPU(s) list:   0-31
Thread(s) per core:    2
Core(s) per socket:    8
Socket(s):             2
NUMA node(s):          2
Vendor ID:             GenuineIntel
CPU family:            6
Model:                 63
Model name:            Intel(R) Xeon(R) CPU E5-2640 v3 @ 2.60GHz
Stepping:              2
CPU MHz:               2799.976
BogoMIPS:              5204.37
Virtualization:        VT-x
L1d cache:             32K
L1i cache:             32K
L2 cache:              256K
L3 cache:              20480K
NUMA node0 CPU(s):     0,2,4,6,8,10,12,14,16,18,20,22,24,26,28,30
NUMA node1 CPU(s):     1,3,5,7,9,11,13,15,17,19,21,23,25,27,29,31
```


### 故障原因

（取自这里：[RHEV: VM paused, cannot resume](https://access.redhat.com/solutions/1758133)）：

> A MMIO Page Fault related bug in KVM causes this. This is a Hypervisor side bug and any VM OS might hit it.
The probability of hitting it appears to be quite low (many VMs running for days).


### 解决方案

- **RHEL 7.3**: Upgrade to kernel-3.10.0-514.el7 from Errata [RHSA-2016:2574](https://rhn.redhat.com/errata/RHSA-2016-2574.html) or later
- **RHEL 7.2.z**: Upgrade to kernel-3.10.0-327.3.1.el7 from Errata [RHSA-2015:2552](https://access.redhat.com/errata/RHSA-2015:2552) or later


### 其他

- Red Hat Enterprise Virtualization (**RHEV**) is Red Hat Inc.'s server virtualization platfor.
- ZStack 最新推出的混合云产品，领先业内、第一家实现了控制面和数据面的完全打通，给用户提供无缝混合云的体验。ZStack 通过标准的混合云产品，提供“互连、灾备、服务、一键迁云”四大场景。



----------


## more kernel bugs

### 故障信息

### 系统信息

### 故障原因

### 解决方案


----------
