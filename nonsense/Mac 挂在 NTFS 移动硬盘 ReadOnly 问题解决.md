

参考：[这里](http://blog.csdn.net/sunbiao0526/article/details/8566317)


```shell
sunfeideMacBook-Pro:data sunfei$ diskutil info /Volumes/Seagate\ Backup\ Plus\ Drive/

   Device Identifier:        disk2s1
   Device Node:              /dev/disk2s1
   Whole:                    No
   Part of Whole:            disk2
   Device / Media Name:      Untitled 1
   Volume Name:              Seagate Backup Plus Drive
   Mounted:                  Yes
   Mount Point:              /Volumes/Seagate Backup Plus Drive
   File System Personality:  NTFS
   Type (Bundle):            ntfs
   Name (User Visible):      Windows NT File System (NTFS)
   Partition Type:           Windows_NTFS
   OS Can Be Installed:      No
   Media Type:               Generic
   Protocol:                 USB
   SMART Status:             Not Supported
   Volume UUID:              F7F4131A-1567-4A15-9D2F-68A69496DA0D
   Total Size:               1.0 TB (1000203836928 Bytes) (exactly 1953523119 512-Byte-Units)
   Volume Free Space:        785.6 GB (785610182656 Bytes) (exactly 1534394888 512-Byte-Units)
   Device Block Size:        512 Bytes
   Allocation Block Size:    4096 Bytes
   Read-Only Media:          No
   Read-Only Volume:         Yes
   Device Location:          External
   Removable Media:          No

sunfeideMacBook-Pro:data sunfei$ 


sunfeideMacBook-Pro:data sunfei$ hdiutil eject /Volumes/Seagate\ Backup\ Plus\ Drive/

"disk2" unmounted.
"disk2" ejected.

sunfeideMacBook-Pro:data sunfei$ 


sunfeideMacBook-Pro:data sunfei$ sudo mkdir /Volumes/MyHD
Password:
sunfeideMacBook-Pro:data sunfei$ ll /Volumes/
total 8
lrwxr-xr-x  1 root  admin   1  6  5 12:04 Macintosh HD -> /
drwxr-xr-x+ 2 root  admin  68  6  6 00:28 MyHD
sunfeideMacBook-Pro:data sunfei$ 


sunfeideMacBook-Pro:data sunfei$ sudo mount_ntfs -o rw,nobrowse /dev/disk2s1 /Volumes/MyHD/
sunfeideMacBook-Pro:data sunfei$


sunfeideMacBook-Pro:data sunfei$ mount

/dev/disk1 on / (hfs, local, journaled)
devfs on /dev (devfs, local, nobrowse)
map -hosts on /net (autofs, nosuid, automounted, nobrowse)
map auto_home on /home (autofs, automounted, nobrowse)
/dev/disk2s1 on /Volumes/MyHD (ntfs, local, noowners, nobrowse)

sunfeideMacBook-Pro:data sunfei$ 

```
