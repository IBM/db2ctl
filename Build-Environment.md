# Build test/production environments

## Minimum 4 machines

We need a minimum of 4 VMs or bare metal to test/deploy the LinBit, Pacemaker and Db2 to build a Db2 WH cluster using software defined storage (LinBit).

## Prerequisites

Enable password less SSH between all machines. If you have not done it before, you can consult https://github.com/vikramkhatri/sshsetup to use an automated procedure.


To use software defined storage, we need a set of NVMe drives on each machine. If you are testing the solution in a VM environment, create equal sized volumes on each machine.

For example:

Check if you have 4 volumes to simulate 4 NVMe drives. For example: You need to have `vdb`, `vdc`, `vdd` and `vde` - which are 25 GB in size each.

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

After update, reboot all the VMs and make sure that your OS is at the latest level.

```
# cat /etc/redhat-release
CentOS Linux release 8.2.2004 (Core)
```

Note: You can also use CentOS 8.2 as well. 

There is no Python on CentOS 8, use the following commands to install Python3 in all VMs.

```
dnf -y install python3
alternatives --set python /usr/bin/python3
```

## Deploy LinBit rpm

Contact LinBit to get a trial license for the LinBit software defined storage. 

https://www.linbit.com/contact-us/


```
wget --no-check-certificate https://my.linbit.com/linbit-manage-node.py
chmod u+x linbit-manage-node.py
./linbit-manage-node.py --> To log in, use the user id and password given by LinBit.
```

After providing user id and the key. Press `Y` to write.

```
Writing registration data:
--> Write to file (/var/lib/drbd-support/registration.json)? [y/N]y
```

Press 1, 2 and then 0 to enable the above.

```
  Here are the repositories you can enable:

    1) pacemaker-2(Disabled)
    2) drbd-9.0(Disabled)

  Enter the number of the repository `1` and `2` to enable. Type `0` when you are done.

  Enable/Disable:
```

Press Y to all prompts.

## Get db2ctl

```
curl -Ls https://git.io/db2ctl | /bin/bash
```

## Generate `db2ctl-sample.yaml` file

```
db2ctl init
```

Copy `db2ctl-sample.yaml` to `db2ctl.yaml` and make changes for your test or production cluster.

You can keep default values in `db2ctl.yaml` for all parameters except the following:

```
      numNodes: 4 #Example: 4, 8, 12 16 etc.
      nvmeList: [[{name: /dev/xvdc, size: 25GB}, {name: /dev/xvde, size: 25GB}, {name: /dev/xvdf, size: 25GB}, {name: /dev/xvdg, size: 25GB}], [{name: /dev/xvdc, size: 25GB}, {name: /dev/xvde, size: 25GB}, {name: /dev/xvdf, size: 25GB}, {name: /dev/xvdg, size: 25GB}], [{name: /dev/xvdc, size: 25GB}, {name: /dev/xvde, size: 25GB}, {name: /dev/xvdf, size: 25GB}, {name: /dev/xvdg, size: 25GB}], [{name: /dev/xvdc, size: 25GB}, {name: /dev/xvde, size: 25GB}, {name: /dev/xvdf, size: 25GB}, {name: /dev/xvdg, size: 25GB}]] # GB or TB only
      ipAddresses: [10.11.1.134, 10.11.1.135, 10.11.1.242, 10.11.2.20] #Example: [10.xx.xx.xx, 10.xx.xx.xx, ...]
      names: [p01.zinox.com, p02.zinox.com, p03.zinox.com, p04.zinox.com] #Example: [e1n1, e2n1, ...]
      partitions: 24
```

`numNodes` - it has to be in the multiples of 4. For example: 4,8,12,16,20,24,28,32 etc.
`nvmeList`: Each host must have a minimum of 4 NVMe drives (recommended). Specify the name of each NVMe and its size and repeat for all hosts. For a development environment, you may go with a minimum of 4x25GB volumes. For a production environment, you should consider a minimum of 4 NVMe drives of a minimum 4 TB per drive. 
`ipAddress` - Use the internal IP address of the high speed internal network that connects each host with a switch. It is recommended that you use a minimum of Mellanox ConnectX 4-En/5-en card of 100 Gbps having dual port. You could also use 25 Gbps cards but the price difference is minor between the two. It is suggested that you use the team network or bonding network if using more than one NIC per machine. Do not use public network IPs which might be using a slow 1 Gbps/10 Gbps network.
`names` - These are the hostnames associated with the internal network IP addresses that you defined in the previous step. Do not use public hostnames.
`partitions` - In a development environment, you can use 6 database partitions per host and have a minimum of 4 cores and 16 GB RAM. This will give you 24 total partitions in a 4 node cluster. For a production environment, you should consider a ratio of minimum of 64GB of RAM per database partition. For example, if you have 192GB RAM with 16 physical Cores, you could go with 3 database partitions per server. If your server has 32 cores and 384 GB RAM, you may consider using 6 database partitions. You should have a minimum of 12 GB/CPU physical core of memory.   

The other parameter that you must change is the virtual IP address used for automatic failover of the NFS server.

```
  nfs:
    server:
      required: 
        nfsVirtualIP: 10.11.1.2
        cidrNetmask: 20
        nic: eth0
```

In above the `nfsVirtualIP` is the one which can be used for the NFS server so that there is no need to unmount/mount the NFS volume when a server that has the NFS server fails. The pacemaker will fail the NFS server to a host that has a replicated copy and the virtual IP address will be moved to the new NFS server.

If you need to run Db2 warehouse where each or group of a machines run in a different subnet, we have a solution where we can use other method. This solution uses IP tables for using a loopback address and it can work across several subnets. Please open an issue if you want to test the solution which does not use virtual IP address.


## Generate Scripts

```
db2ctl generate all
```

With above, you can generate scripts for linbit, pacemaker and db2 deployments. The scripts are generated automatically when you use something like `db2ctl install all` or `db2ctl cleanup all`. If you make changes in the generated script, you can use `-n` switch to disable generation of the script at run time if you use `db2ctl install` or `db2ctl cleanup` commands.

## Install LinBit Software Defined Storage

This will install LinBit software defined storage with automatic replication bin packing designed specially for Db2 Warehouse to ensure localization of the data when a failover occurs.

The above will install 4-way replicated NFS volume which will be used to host shared SQLLIB directory on all hosts.

```
db2ctl install linbit
```

## Install Pacemaker for LinBit volume failures

```
db2ctl install pacemaker
```

The above will install Pacemaker for DRBD volumes and their automatic replication bin packing so that we can sustain one machine failure per four machines and still be up and running as far as Db2 warehouse is concerned.

## Install Db2 Warehouse

```
db2ctl install db2
```

In order to use the above, download the trial license of Db2 from https://www.ibm.com/analytics/db2/trials.

If you already have Db2 entitlements, you can copy the db2 tar file to `/tmp/db2ctl-download` directory. Also copy the license file to the same directory.

The name of the db2 binary and the license file should be defined in the `db2ctl.yaml` file.

```
      db2Binary: v11.5.4_linuxx64_server.tar.gz # Db2 binary that you get from IBM using try and buy
      db2License: db2adv_vpc.lic # The db2 license file.
```

If you are using a different version of db2, specify the proper file name and the license file name in the `db2ctl.yaml` file. 


