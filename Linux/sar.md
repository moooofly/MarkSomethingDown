# sar

## 常用选项

```
# equivalent to specifying -bBdHqrRSuvwWy -I SUM -I XALL -m ALL -n ALL -u ALL -P ALL
sar -A

# Report CPU utilization.
sar -u

# Report memory statistics.
sar -R

# Report memory utilization statistics.
sar -r

# Report I/O and transfer rate statistics.
sar -b

# Report paging statistics.
sar -B

# Report hugepages utilization statistics.
sar -H

# Report network statistics.
sar -n DEV

# Report queue length and load averages.
sar -q

# Report swapping statistics.
sar -W

# Report task creation and system switching activity.
sar -w
```

## 和 sar 相关的 cron 任务

```
[root@wg-redis-corvus-71: ~]# cat /etc/cron.d/sysstat
# Run system activity accounting tool every 10 minutes
*/10 * * * * root /usr/lib64/sa/sa1 1 1
# 0 * * * * root /usr/lib64/sa/sa1 600 6 &
# Generate a daily summary of process accounting at 23:53
53 23 * * * root /usr/lib64/sa/sa2 -A

[root@wg-redis-corvus-71: ~]#
[root@wg-redis-corvus-71: ~]# ll /var/log/sa
total 65896
-rw-r--r-- 1 root root 1318496 Aug  1 23:50 sa01
-rw-r--r-- 1 root root  595488 Aug  2 10:40 sa02
-rw-r--r-- 1 root root 1318496 Jul  4 23:50 sa04
-rw-r--r-- 1 root root 1318496 Jul  5 23:50 sa05
-rw-r--r-- 1 root root 1318496 Jul  6 23:50 sa06
-rw-r--r-- 1 root root 1318496 Jul  7 23:50 sa07
-rw-r--r-- 1 root root 1318496 Jul  8 23:50 sa08
-rw-r--r-- 1 root root 1318496 Jul  9 23:50 sa09
-rw-r--r-- 1 root root 1318496 Jul 10 23:50 sa10
-rw-r--r-- 1 root root 1318496 Jul 11 23:50 sa11
-rw-r--r-- 1 root root 1318496 Jul 12 23:50 sa12
-rw-r--r-- 1 root root 1318496 Jul 13 23:50 sa13
-rw-r--r-- 1 root root 1318496 Jul 14 23:50 sa14
-rw-r--r-- 1 root root 1318496 Jul 15 23:50 sa15
-rw-r--r-- 1 root root 1318496 Jul 16 23:50 sa16
-rw-r--r-- 1 root root 1318496 Jul 17 23:50 sa17
-rw-r--r-- 1 root root 1318496 Jul 18 23:50 sa18
-rw-r--r-- 1 root root 1318496 Jul 19 23:50 sa19
-rw-r--r-- 1 root root 1318496 Jul 20 23:50 sa20
-rw-r--r-- 1 root root 1318496 Jul 21 23:50 sa21
-rw-r--r-- 1 root root 1318496 Jul 22 23:50 sa22
-rw-r--r-- 1 root root 1318496 Jul 23 23:50 sa23
-rw-r--r-- 1 root root 1318496 Jul 24 23:50 sa24
-rw-r--r-- 1 root root 1318496 Jul 25 23:50 sa25
-rw-r--r-- 1 root root 1318496 Jul 26 23:50 sa26
-rw-r--r-- 1 root root 1318496 Jul 27 23:50 sa27
-rw-r--r-- 1 root root 1318496 Jul 28 23:50 sa28
-rw-r--r-- 1 root root 1318496 Jul 29 23:50 sa29
-rw-r--r-- 1 root root 1318496 Jul 30 23:50 sa30
-rw-r--r-- 1 root root 1318496 Jul 31 23:50 sa31
-rw-r--r-- 1 root root  981679 Aug  1 23:53 sar01
-rw-r--r-- 1 root root  981679 Jul  4 23:53 sar04
-rw-r--r-- 1 root root  981679 Jul  5 23:53 sar05
-rw-r--r-- 1 root root  981679 Jul  6 23:53 sar06
-rw-r--r-- 1 root root  981679 Jul  7 23:53 sar07
-rw-r--r-- 1 root root  981679 Jul  8 23:53 sar08
-rw-r--r-- 1 root root  981679 Jul  9 23:53 sar09
-rw-r--r-- 1 root root  981679 Jul 10 23:53 sar10
-rw-r--r-- 1 root root  981679 Jul 11 23:53 sar11
-rw-r--r-- 1 root root  981679 Jul 12 23:53 sar12
-rw-r--r-- 1 root root  981679 Jul 13 23:53 sar13
-rw-r--r-- 1 root root  981679 Jul 14 23:53 sar14
-rw-r--r-- 1 root root  981679 Jul 15 23:53 sar15
-rw-r--r-- 1 root root  981679 Jul 16 23:53 sar16
-rw-r--r-- 1 root root  981679 Jul 17 23:53 sar17
-rw-r--r-- 1 root root  981679 Jul 18 23:53 sar18
-rw-r--r-- 1 root root  981679 Jul 19 23:53 sar19
-rw-r--r-- 1 root root  981679 Jul 20 23:53 sar20
-rw-r--r-- 1 root root  981679 Jul 21 23:53 sar21
-rw-r--r-- 1 root root  981679 Jul 22 23:53 sar22
-rw-r--r-- 1 root root  981679 Jul 23 23:53 sar23
-rw-r--r-- 1 root root  981679 Jul 24 23:53 sar24
-rw-r--r-- 1 root root  981679 Jul 25 23:53 sar25
-rw-r--r-- 1 root root  981679 Jul 26 23:53 sar26
-rw-r--r-- 1 root root  981679 Jul 27 23:53 sar27
-rw-r--r-- 1 root root  981679 Jul 28 23:53 sar28
-rw-r--r-- 1 root root  981679 Jul 29 23:53 sar29
-rw-r--r-- 1 root root  981679 Jul 30 23:53 sar30
-rw-r--r-- 1 root root  981679 Jul 31 23:53 sar31
[root@wg-redis-corvus-71: ~]#
```

从输出中可以看到：

- cron 任务中使用了 `/usr/lib64/sa/sa1` 和 `/usr/lib64/sa/sa2` ；
- sar 日志生成在 /var/log/sa 目录下，以天为单位按月滚动；

## sa1

`sa1` - Collect and store binary data in the system activity daily data file.

```
/usr/lib64/sa/sa1 [ --boot | interval count ]
```

The  `sa1`  command is a shell procedure variant of the `sadc` command and handles all of the flags and parameters of that command. The `sa1` command collects and stores binary data in the `/var/log/sa/sadd` file, where the dd  parameter  indicates  the current  day.  The interval and count parameters specify that the record should be written count times at interval seconds.

If no arguments are given to `sa1` then a single record is written.

The `sa1` command is designed to be started automatically by the `cron` command.
       
`--boot` This option tells `sa1` that the `sadc` command should be called without specifying the interval and count parameters in order to insert a dummy record, marking the time when the counters restarts from 0.

`/var/log/sa/sadd` Indicate the daily data file, where the dd parameter is a number representing the day of the month.

```
[root@wg-public-rediscluster-119: ~]# cat /usr/lib64/sa/sa1
#!/bin/sh
# /usr/lib64/sa/sa1
# (C) 1999-2012 Sebastien Godard (sysstat <at> orange.fr)
#
#@(#) sysstat-10.1.5
#@(#) sa1: Collect and store binary data in system activity data file.
#

# Set default value for some variables.
# Used only if ${SYSCONFIG_DIR}/sysstat doesn't exist!
HISTORY=0
SADC_OPTIONS=""
DDIR=/var/log/sa
DATE=`date +%d`
CURRENTFILE=sa${DATE}
CURRENTDIR=`date +%Y%m`
SYSCONFIG_DIR=/etc/sysconfig
umask 0022
[ -r ${SYSCONFIG_DIR}/sysstat ] && . ${SYSCONFIG_DIR}/sysstat
if [ ${HISTORY} -gt 28 ]
then
	cd ${DDIR} || exit 1
	[ -d ${CURRENTDIR} ] || mkdir -p ${CURRENTDIR}
	# If ${CURRENTFILE} exists and is a regular file, then make sure
       	# the file was modified this day (and not e.g. month ago)
	# and move it to ${CURRENTDIR}
	[ ! -L ${CURRENTFILE} ] &&
		[ -f ${CURRENTFILE} ] &&
		[ "`date +%Y%m%d -r ${CURRENTFILE}`" = "${CURRENTDIR}${DATE}" ] &&
		mv -f ${CURRENTFILE} ${CURRENTDIR}/${CURRENTFILE}
	touch ${CURRENTDIR}/${CURRENTFILE}
	# Remove the "compatibility" link and recreate it to point to
	# the (new) current file
	rm -f ${CURRENTFILE}
	ln -s ${CURRENTDIR}/${CURRENTFILE} ${CURRENTFILE}
else
	# If ${CURRENTFILE} exists, is a regular file and is from a previous
	# month then delete it so that it is recreated by sadc afresh
	[ -f ${CURRENTFILE} ] && [ "`date +%Y%m -r ${CURRENTFILE}`" -lt "${CURRENTDIR}" ] && rm -f ${CURRENTFILE}
fi
ENDIR=/usr/lib64/sa
cd ${ENDIR}
[ "$1" = "--boot" ] && shift && BOOT=y || BOOT=n
if [ $# = 0 ] && [ "${BOOT}" = "n" ]
then
# Note: Stats are written at the end of previous file *and* at the
# beginning of the new one (when there is a file rotation) only if
# outfile has been specified as '-' on the command line...
	exec ${ENDIR}/sadc -F -L ${SADC_OPTIONS} 1 1 -
else
	exec ${ENDIR}/sadc -F -L ${SADC_OPTIONS} $* -
fi

[root@wg-public-rediscluster-119: ~]#
```


## sa2


`sa2` - Write a daily report in the `/var/log/sa` directory.

```
/usr/lib64/sa/sa2
```

The  `sa2` command is a shell procedure variant of the `sar` command which writes a daily report in the `/var/log/sa/sardd` file, where the dd parameter indicates the current day. The `sa2` command handles all of the flags and parameters of the `sar` command.

The `sa2` command is designed to be started automatically by the cron command.

```
[root@wg-public-rediscluster-119: ~]# cat /usr/lib64/sa/sa2
#!/bin/sh
# /usr/lib64/sa/sa2
# (C) 1999-2012 Sebastien Godard (sysstat <at> orange.fr)
#
#@(#) sysstat-10.1.5
#@(#) sa2: Write a daily report
#
S_TIME_FORMAT=ISO ; export S_TIME_FORMAT
umask 0022
prefix=/usr
exec_prefix=/usr
# Add a trailing slash so that 'find' can go through this directory if it's a symlink
DDIR=/var/log/sa/
SYSCONFIG_DIR=/etc/sysconfig
YESTERDAY=
DATE=`date ${YESTERDAY} +%d`
CURRENTFILE=sa${DATE}
CURRENTRPT=sar${DATE}
HISTORY=28
COMPRESSAFTER=31
ZIP="bzip2"
[ -r ${SYSCONFIG_DIR}/sysstat ] && . ${SYSCONFIG_DIR}/sysstat
if [ ${HISTORY} -gt 28 ]
then
	CURRENTDIR=`date ${YESTERDAY} +%Y%m`
	cd ${DDIR} || exit 1
	[ -d ${CURRENTDIR} ] || mkdir -p ${CURRENTDIR}
	# Check if ${CURRENTFILE} is the correct file created at ${DATE}
	# Note: using `-ge' instead of `=' since the file could have
	# the next day time stamp because of the file rotating feature of sadc
	[ -f ${CURRENTFILE} ] &&
		[ "`date +%Y%m%d -r ${CURRENTFILE}`" -ge "${CURRENTDIR}${DATE}" ] || exit 0
	# If the file is a regular file, then move it to ${CURRENTDIR}
	[ ! -L ${CURRENTFILE} ] &&
		mv -f ${CURRENTFILE} ${CURRENTDIR}/${CURRENTFILE} &&
			ln -s ${CURRENTDIR}/${CURRENTFILE} ${CURRENTFILE}
	touch ${CURRENTDIR}/${CURRENTRPT}
	# Remove the "compatibility" link and recreate it to point to
	# the (new) current file
	rm -f ${CURRENTRPT}
	ln -s ${CURRENTDIR}/${CURRENTRPT} ${CURRENTRPT}
	CURRENTDIR=${DDIR}/${CURRENTDIR}
else
	CURRENTDIR=${DDIR}
fi
RPT=${CURRENTDIR}/${CURRENTRPT}
ENDIR=/usr/bin
DFILE=${CURRENTDIR}/${CURRENTFILE}
[ -f "$DFILE" ] || exit 0
cd ${ENDIR}
[ -L ${RPT} ] && rm -f ${RPT}
${ENDIR}/sar $* -f ${DFILE} > ${RPT}
find ${DDIR} \( -name 'sar??' -o -name 'sa??' -o -name 'sar??.xz' -o -name 'sa??.xz' -o -name 'sar??.gz' -o -name 'sa??.gz' -o -name 'sar??.bz2' -o -name 'sa??.bz2' \) \
	-mtime +"${HISTORY}" -exec rm -f {} \;
find ${DDIR} \( -name 'sar??' -o -name 'sa??' \) -type f -mtime +"${COMPRESSAFTER}" \
	-exec ${ZIP} {} \; > /dev/null 2>&1
# Remove broken links
for f in `find ${DDIR} \( -name 'sar??' -o -name 'sa??' \) -type l`; do
	[ -e $f ] || rm -f $f
done
cd ${DDIR}
rmdir [0-9]????? > /dev/null 2>&1
exit 0

[root@wg-public-rediscluster-119: ~]#
```

## sadc

sadc - System activity data collector.

```
/usr/lib64/sa/sadc  [  -C comment ] [ -F ] [ -L ] [ -V ] [ -S { INT | DISK | SNMP | IPV6 | POWER | XDISK | ALL | XALL } ] [interval [ count ] ] [ outfile ]
```

The sadc command is intended to be used as a backend to the sar command.


## sadf

sadf - Display data collected by sar in multiple formats.
