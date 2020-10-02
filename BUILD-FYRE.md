# Build test environment on Fyre

## Minimum 5 machines

We need a minimum of 5 VMs to test the LinBit, Pacemaker and Db2 to simulate delivery on Cloud Pak for Data System.

We will use the first 4 machines as a LinStor, Pacemaker and Db2 cluster.

The 5th machine will be used as a CPDS control node to deploy everything on the 4 machines.

Assume that password less ssh will be available in CPDS and we will deploy the full Db2 WH stack using the control node.

## Prerequisites

Check if you have 4 volumes to simulate 4 NVMe drives in CPDS. For example: You need to have `vdb`, `vdc`, `vdd` and `vde` - which are 25 GB in size each.

```
# lsblk
NAME        MAJ:MIN RM  SIZE RO TYPE MOUNTPOINT
vda         252:0    0  250G  0 disk
├─vda1      252:1    0    1G  0 part /boot
└─vda2      252:2    0  249G  0 part
  ├─cl-root 253:0    0  233G  0 lvm  /
  └─cl-swap 253:1    0   16G  0 lvm  [SWAP]
vdb         252:16   0   25G  0 disk
vdc         252:32   0   25G  0 disk
vdd         252:48   0   25G  0 disk
vde         252:64   0   25G  0 disk
```

Update OS on all VMs

```
dnf -y update
```

After update, reboot all the VMs and make sure that your OS is at

```
# cat /etc/redhat-release
CentOS Linux release 8.2.2004 (Core)
```

We will now fudge CentOS to make it look like RedHat. Run these on all VMs.

```

```

There is no Python on CentOS 8, use the following commands to install Python3 in all VMs.

```
dnf -y install python3
alternatives --set python /usr/bin/python3
```

## Deploy LinBit rpm

As of now, we need to get the LinBit RPMs from LinBit repo but that will change going forward in which the software will be installed using our own repo.

We need to follow this procedure to register the LinBit repo on all 4 machines. Do not do this on the 5th machine as that is the control node which does not require LinBit rpm.

Get the LinBit repo starter kit

Repeat the following on first 4 machines. (We do not need to run on the 5th or control node)

```
wget --no-check-certificate https://my.linbit.com/linbit-manage-node.py
chmod u+x linbit-manage-node.py
./linbit-manage-node.py --> To log in, use the following key  `ibmperformanceserver/ZD87fN6F`
```

Type user id and the key. Press Y to write.

```
Writing registration data:
--> Write to file (/var/lib/drbd-support/registration.json)? [y/N]y
```

Press 1, 2 and then 0 to enable the above.

```
  Here are the repositories you can enable:

    1) pacemaker-2(Disabled)
    2) drbd-9.0(Disabled)

  Enter the number of the repository you wish to enable/disable. Hit 0 when you are done.

  Enable/Disable:
```

Press Y to all prompts.

## Get db2ctl

```
curl -Ls https://git.io/db2ctl | /bin/bash
```

## Generate `db2pc-sample.yaml` file

```
db2ctl init
```

Edit `db2pc-sample.yaml` to specify values for the cluster.

For example:

Add the IP address for each host:

Assume that you will use eth0 network interface for the communication.

```
 ip a s eth0 | grep ine[t] | awk '{print $2}'
```

Change the host names for each host.

Take care: Use a floating IP address in the same CIDR range which is not assigned to any host.

## Generate Scripts (Optional step)

```
db2ctl generate all
```

## Install components

```
db2ctl install linbit
```
