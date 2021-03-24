# PXE Boot Server

Based upon https://www.golinuxcloud.com/configure-kickstart-pxe-boot-server-linux/

## Create a CentOS 8 VM

Follow the procedure to create a minimal CentOS 8 VM.

Use Paritions and not LVMs to keep things simple. No Swap space

Add `net.ifnames=0` to the `GRUB_CMDLINE_LINUX`.

```
[root@pxe ~]# cat /etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="$(sed 's, release .*$,,g' /etc/system-release)"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="crashkernel=auto net.ifnames=0"
GRUB_DISABLE_RECOVERY="true"
GRUB_ENABLE_BLSCFG=true
```

Install grub-tools

```
yum -y install grub2-tools
grub2-mkconfig -o /boot/grub2/grub.cfg
```

Reboot and then change network config
```
reboot
cd /etc/sysconfig/network-scripts/
mv ifcfg-ens33 ifcfg-eth0
```

Edit `ifcfg-eth0` to change the `bootproto` to `dhcp` and `device` and `name` to `eth0` and `onboot` to `yes`.
```
[root@pxe ~]# cat /etc/sysconfig/network-scripts/ifcfg-eth0
TYPE=Ethernet
PROXY_METHOD=none
BROWSER_ONLY=no
BOOTPROTO=dhcp
DEFROUTE=yes
IPV4_FAILURE_FATAL=no
IPV6INIT=yes
IPV6_AUTOCONF=yes
IPV6_DEFROUTE=yes
IPV6_FAILURE_FATAL=no
IPV6_ADDR_GEN_MODE=stable-privacy
NAME=eth0
UUID=c5a28042-a445-4530-95ce-0f11903a387a
DEVICE=eth0
ONBOOT=yes
```

Reboot

```
reboot
```

### Install rsyslog

```
yum -y install rsyslog
systemctl enable rsyslog
systemctl start rsyslog
```

### Configure firewall

```
firewall-cmd --permanent --zone=public --change-interface=eth0
firewall-cmd --permanent --zone=public --add-port=80/tcp
firewall-cmd --permanent --zone=public --add-port=69/udp
firewall-cmd --permanent --zone=public --add-service=tftp 
firewall-cmd --reload
```

## Copy RHEL DVD contents to /images

Downbload RHEL 8.2 DVD from RedHat site and assign to CDROM in VMware Workstation. Enable Connected to true in VM Ware Workstation so that the ISO is mounted.

```
mkdir -p /mnt/dvd /images
mount /dev/sr0 /mnt/dvd
cp -apr /mnt/dvd/* /images
cp -apr /mnt/dvd/.discinfo /mnt/dvd/.treeinfo /images
```

### selinux 

```
chcon -R -t public_content_rw_t /images
restorecon -r /images
restorecon -RvF /var/lib/tftpboot/
```

## Install tftp server 

To transfer PXE images files to the client for network based installation

Install tftp-server and xinetd

```
yum -y install tftp-server xinetd
```

Fork `tftp` process using `systemd`

```
[root@pxe ~]# cat /usr/lib/systemd/system/tftp.service
[Unit]
Description=Tftp Server
Requires=tftp.socket
Documentation=man:in.tftpd

[Service]
ExecStart=/usr/sbin/in.tftpd -s /var/lib/tftpboot
StandardInput=socket

[Install]
Also=tftp.socket
```

Notice that the defualt location for TFTP server is `/var/lib/tftpboot`.

Enable and start tftp server.
```
systemctl enable tftp
systemctl start tftp

[root@pxe ~]# systemctl status -l tftp
● tftp.service - Tftp Server
   Loaded: loaded (/usr/lib/systemd/system/tftp.service; indirect; vendor preset: disabled)
   Active: active (running) since Sun 2020-09-27 14:54:16 EDT; 6s ago
     Docs: man:in.tftpd
 Main PID: 1488 (in.tftpd)
    Tasks: 1 (limit: 4858)
   Memory: 220.0K
   CGroup: /system.slice/tftp.service
           └─1488 /usr/sbin/in.tftpd -s /var/lib/tftpboot

Sep 27 14:54:16 pxe systemd[1]: Started Tftp Server.
```

TFTP uses `tftp.socket` to serve TFTP requests.

```
[root@pxe ~]# systemctl status -l tftp.socket
● tftp.socket - Tftp Server Activation Socket
   Loaded: loaded (/usr/lib/systemd/system/tftp.socket; enabled; vendor preset: disabled)
   Active: active (running) since Sun 2020-09-27 14:54:16 EDT; 1min 12s ago
   Listen: [::]:69 (Datagram)
   CGroup: /system.slice/tftp.socket

Sep 27 14:54:16 pxe systemd[1]: Listening on Tftp Server Activation Socket.
```

The TFTP process may become inactive when there are no incoming TFTP requests but the tftp.socket will start the service.

## Setup PXE Boot server

To perform PXE network based installation, we will configure PXE boot server. We will need Linux Boot Image to boot RHEL 8 OS. The boot will need `initrd` and `vmlinuz` to load necessary drivers from the memory to boot up the server.

Create directory `pxelinux` under `/var/lib/tftpboot` to store PXE images.

```
mkdir -p /var/lib/tftpboot/pxelinux
```

### Extract syslinux-tftboot

The `pxelinux` file is part of `syslinux-tftpboot` rpm so we will copy this from RHEL ISO to a temporary location.

```
cp /mnt/dvd/BaseOS/Packages/syslinux-tftpboot-6.04-4.el8.noarch.rpm /tmp
```

Extract the file.

```
cd /tmp
rpm2cpio syslinux-tftpboot-6.04-4.el8.noarch.rpm | cpio -idm
```

The above will extract `syslinux-tfpboot` under `/tmp`. We will copy this file to `/var/lib/tfpboot/pxelinux`.

```
cp /tmp/tftpboot/ldlinux.c32 /var/lib/tftpboot/pxelinux/
cp /tmp/tftpboot/pxelinux.0 /var/lib/tftpboot/pxelinux/
```

### Copy initrd and vmlinuz

We also need to copy other pxe boot images `initrd.img` and `vmlinuz` to `/var/lib/tftpboot/pxelinux/`

```
cp /images/isolinux/initrd.img /var/lib/tftpboot/pxelinux/
cp /images/isolinux/vmlinuz /var/lib/tftpboot/pxelinux/
```

List all files.
```
# ls -l /var/lib/tftpboot/pxelinux/
total 73820
-r--r--r--. 1 root root 66509268 Sep 27 15:08 initrd.img
-rw-r--r--. 1 root root   116096 Sep 27 15:06 ldlinux.c32
-rw-r--r--. 1 root root    42821 Sep 27 15:06 pxelinux.0
-r-xr-xr-x. 1 root root  8913760 Sep 27 15:09 vmlinuz
```

### Create Boot Menu

The PXE boot server can be used to create multiple install images for OS. We will create a boot menu.

Create a file `boot.msg` under `/var/lib/tftpboot/pxelinux`
```
cat << EOF > /var/lib/tftpboot/pxelinux/boot.msg
Available boot options:

  1 - Install Red Hat Enterprise Linux 8
  2 - Boot from media
EOF
```

### Create PXE Configuration file

* Once the client retrieves and executes `pxelinux.0`, it is hard-coded to look for a file from the pxelinux.cfg/ sub directory relative to where pxelinux.0 was found
* In large deployments we create individual PXE configuration file per node.
* The naming syntax of these PXE configuration file is very important
* First, it will look for a file named after the MAC address, in the form 01-xx-xx-xx-xx-xx-xx
* For example, in this example my client't NIC MAC Address is 08:00:27:83:1e:2a so my PXE configuration file will be 01-08-00-27-83-1e-2a
* Next, it will look for a file named by the IP address as provided by the DHCP server.
* The IP address is looked up in hexadecimal format.
* You can use printf to get the hexadecimal format of an IP Address, for example to get the hexadecimal value of 10.10.10.12


```
# printf "%02x%02x%02x%02xn" 10 10 10 12
0a0a0a0cn
```

* Since we are not assigning any static IP Address, we cannot use the hexadecimal format file. Although we can use MAC address based file.
* If both MAC Address and hexadecimal format file are not found under pxelinux.cfg then the installer will look for "default" file
* In this example we will use default file but I have also verified the PXE network installation using MAC based file 01-08-00-27-83-1e-2a
* We will create a new directory pxelinux.cfg

```
mkdir /var/lib/tftpboot/pxelinux/pxelinux.cfg
```

The MAC address of the NIC of the machine to install RHEL 8 : 00:50:56:38:13:21

We will now create a file 01-00-50-56-38-13-21

Create the PXE configuration file `/var/lib/tftpboot/pxelinux/pxelinux.cfg/01-00-50-56-38-13-21`

## Mac Addresses
```
RHEl8.2 = 00:50:56:38:13:21
VM1     = 00:50:56:36:83:33
VM2     = 00:50:56:38:39:3A
VM3     = 00:50:56:23:69:A0
VM4     = 00:50:56:3F:F3:92
VM5     = 00:50:56:2F:A7:54
```
```
cat << EOF > /var/lib/tftpboot/pxelinux/pxelinux.cfg/01-00-50-56-38-13-21
timeout 300
default 1
prompt  1

display boot.msg

label 1
  menu label ^Install Red Hat Enterprise Linux 8
  menu default
  kernel vmlinuz
  append initrd=initrd.img showopts ks=http://192.168.142.128/kickstart/kickstart.conf ip=dhcp net.ifnames=0 biosdevname=0

label local
  menu label Boot from ^local media
  localboot 0

menu end
EOF
```

Likewise, we can create the above file for each NIC of the each host on which we want to start the RHEL 8 install.

## Create a kickstart file

You can create an encrypted password for root.

```
python3 -c 'import crypt,getpass;pw=getpass.getpass();print(crypt.crypt(pw) if (pw==getpass.getpass("Confirm: ")) else exit())'
```

Kickstart file

```
mkdir -p /var/www/html/kickstart
cat << EOF > /var/www/html/kickstart/kickstart.conf 
#version=RHEL8

# System authorization information
auth --enableshadow --passalgo=sha512

# Use Network installation media
url --url="http://192.168.142.128/rhel/"

# Run the Setup Agent on first boot
firstboot --enable

# Run the Setup Agent on first boot
ignoredisk --only-use=nvme0n1

# Partition clearing information
clearpart --all --initlabel

# Disk partitioning information
autopart --type=plain --fstype=xfs --nohome

# EULA
eula --agreed

# Firewall
firewall --disabled

# Use text install
text --non-interactive

# Keyboard layouts
keyboard --vckeymap=us --xlayouts='us'

# SELinux
selinux --disabled

# System language
lang en_US.UTF-8

# Network information
network  --bootproto=dhcp --device=eth0 --ipv6=ignore --activate
network  --bootproto=dhcp --device=eth1 --ipv6=ignore --activate
network  --hostname=node01.ibm.local

# Root password
rootpw --plaintext password 

# Do not configure the X Window System
skipx

# System services
services --enabled="chronyd"

# System timezone
timezone America/New_York --isUtc

# Reboot after installation completes
reboot

%packages
@^minimal-environment
kexec-tools
%end

%addon com_redhat_kdump --enable --reserve-mb='auto'
%end

%anaconda
pwpolicy root --minlen=6 --minquality=1 --notstrict --nochanges --notempty
pwpolicy user --minlen=6 --minquality=1 --notstrict --nochanges --emptyok
pwpolicy luks --minlen=6 --minquality=1 --notstrict --nochanges --notempty
%end

%post --interpreter=/bin/bash --log=/root/postinstall1.log
mkdir -p /root/.ssh
chmod 700 /root/.ssh
echo "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDceMiteuxrqpsotfX/df7o142HaNOGBIT5HT9A5yLwHaS1sywH/xnWsDTqf3G5laS1uWAtPqVXJ9rR5OhBRzHmfHC9CUr4RAIq5mgwJOGCEnX4wkhTdvy9ZMk5WY6I4SuLKaPCJmNEKyOTQFu6jjkYwCAJTge3hDVniawfSx+4OmCJCqG+QOarpOyb5lNLEHfwEmKfBKkVOk0Z/7gyPx7vNx2Joq0l6Fu6m6xDy3chEeh9FTlJjoJnyevjBsZiL1j4/DVYO849jxRiELhHDzHfQ7ELxyh7mkpmqRSH4NEzBLyERxzfwF1EkbdMmIIT/Eh6z0HvDz253D6HkUycFvpJ vikram@Vikram-MB-2018.local" > /root/.ssh/authorized_keys
chmod 600 /root/.ssh/authorized_keys
cat << EOT > /root/.ssh/id_rsa
-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
EOT
cat << EOT > /root/.ssh/id_rsa.pub
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC510gmtDuX6KPptxauGjyBriejEhGyPfVgpRTgNFE1+uNkymybSfRT8V5LRS4Kg4834xgyTSHQgsNVEZsAmH6n84Px8zoOwSZE6EOaksA6HgNMDy6zD0p09KbKLF3fZQqH8Naeh/8vVMlirmUx2TPKIXmVhbBwGySLf6tJZqL+TVlGTogWEhlfVehTO5VRzKw9bPednrqNGjQ9j0SsAS5EvtBakfOOTXFUynZ1n0y92MYKT7cvH1mWQyeFvdFZ5OXsGXHbwtQxbox/7zCladYDGueTWjXAVt/+GLIqTzeyU2ZOk27GdkxLWYDd2M1m4ppI9Wc97ez9dGfb+TPZwnBsvPgzwS2b/LzBNPYkJY9hyQ0teOdoKTMtMFML1V2SiiZ4oiLvuqhvSLNuqNhIyxOAuVPwsJ1wlRve/2ferUA6xLcRErQOLn2SHBN1F23mlMmFSofiOFlcEmmGNFnIenG7vzZ+wk1aw98kSFg2AzDfMshWMeyYyD0DV9Q9b4RgTTM= root@node01.ibm.local
EOT
echo 'host *' > /root/.ssh/config
echo '     StrictHostChecking no' >> /root/.ssh/config
chmod 600 /root/.ssh/id_rsa
%end
EOF
```

Install kickstart validator and check the validity of the file.

```
yum -y install pykickstart
ksvalidator /var/www/html/kickstart/kickstart.conf
echo $?
```

## Configure HTTP server

```
yum -y install httpd
chown -R apache:apache /var/www/html/

systemctl enable httpd
systemctl start httpd
systemctl status -l httpd
```
 Check the file

 ```
 curl http://192.168.142.128/kickstart/kickstart.conf
 ```

Create entry for images as rhel alias

```
cat << EOF > /etc/httpd/conf.d/images.conf
Alias /rhel /images

<Directory "/images">
    Options Indexes MultiViews FollowSymlinks
    AllowOverride all
    Order deny,allow
    Allow from all
    Require all granted
</Directory>
```

Create SE Linux entries

```
chcon -R -t public_content_rw_t /images
restorecon -r /images
chown -R apache.apache /images
```

Restart httpd server

```
systemctl restart httpd
```

Check both URL
```
curl -L http://192.168.142.128/rhel/
curl -L http://192.168.142.128/kickstart/kickstart.conf
```

## Install and Configure DHCP Server

Install dhcp

```
yum -y install dhcp-server
```

DHCP server configuration

```
cat << EOF > /etc/dhcp/dhcpd.conf
allow bootp;
allow booting;
max-lease-time 1200;
default-lease-time 900;
log-facility local7;

option ip-forwarding    false;
option mask-supplier    false;

   subnet 192.168.142.0 netmask 255.255.255.0 {
       option  routers   192.168.142.2;
       option  domain-name-servers  192.168.142.2;
       range 192.168.142.150 192.168.142.225;
       next-server 192.168.142.128;
       filename "pxelinux/pxelinux.0";
   }

   host node01.ibm.local {
       hardware ethernet 00:50:56:38:13:21;
       fixed-address 192.168.142.129;
   }
EOF
```

* The PXE file name is defined with filename. Since the tftp is configured to use `/var/lib/tftpboot` as default location we have provided `pxelinux/pxelinux.0`
* `next-server` defines the IP address of the TFTP server
* range is used to assign IP address for DHCP requests

Enable and start DHCP

```
systemctl enable dhcpd
systemctl start dhcpd
systemctl status dhcpd
```

## Put PXE boot install option at the top 

Put PXE or LAN boot at the top of the boot order

In VM environment

```
VM > Power > Power on to Firmware
```

And, keep Network Boot from the NIC at the top

Using impitool in Lenovo hardware

```
```

