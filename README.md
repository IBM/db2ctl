# DB2CTL

This tool generates Linbit (Storage), Pacemaker (HA) and Db2 (Warehouse) deployment scripts.

The minimum deployment unit for production is four servers with storage replication to provide resiliency from a node failue.

The failure domain is one machine per four servers. It means that you can lose one server out of four servers and Db2 warehouse will continue to run within a recovery time objective under 120 seconds.

# Licenses

The `db2ctl` is free to use. The LinBit software defined storage requires a support licennse on RHEL/CentOS and Db2 requires a commercial license from IBM.

It is possible to evaluate the solution by requesting the trial support license from LinBit and trial license for Db2 from IBM.

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [DB2CTL](#db2ctl)
  - [Terminology](#terminology)
  - [Commands](#commands)
  - [Building the application](#building-the-application)
  - [Generate yaml configuration file](#generate-yaml-configuration-file)
  - [Validate config file](#validate-config-file)
  - [Generate all bash configuration files](#generate-all-bash-configuration-files)
  - [Install config files](#install-config-files)
  - [Cleanup](#cleanup)

<!-- /code_chunk_output -->

## Terminology

There are 2 different configuration files that are mentioned in this repo:

1. The initial `yaml` config file. This is similar to `values.yaml` in helm, and it is used by the tool to create the second set of configuration files.
2. The second set of configuration files are the `bash configuration scripts` that are created by this tool. These files are used for setting up the software defined storage (Linbit), High Availability (Pacemaker), Db2 Warehouse).

In short, the `db2ctl` tool uses a `yaml` config file to create a set of `bash` config files for Linbit, Pacemaker and Db2 for any Db2 production database cluster from 4 nodes to all the way to the 32 nodes.

## Commands

```
db2ctl commands

- init          (creates yaml config file)
- parse         (optional, parses config file for error checking)
- generate
  - all         (optional, generates all shell script files)
  - mapping     (optional, generate mapping file)
  - binpacking  (optional, generate binpacking file)
  - linbit      Generate linbit install scripts
  - pacemaker   Generate pacemaker install scripts
  generate flags
    -b, --bin           (optional) use custom bin packing file when generating config scripts, run 'generate' first to generate bin-packing file
    -m, --map           (optional) use custom mapping yaml file when generating config scripts, run 'generate' first to generate mapping file

- install
  - all         (install all shell scripts created)
  - prereq      (install pre-req modules)
  - linbit      (install linbit modules)
  - pacemaker   (install pacemaker modules)
  - db2         (install db2 modules)

- cleanup
  - all         (cleans all resources) - TODO
  - linbit      (cleans linbit)
  - pacemaker   (cleans pacemaker)
  - db2         (cleans db2) - TODO

- state        (displays install state of all components, accepts optional args)
  - install     ((displays install state of all installed components)
  - cleanup     ((displays install state of all cleanup components)

- version       (displays version info for the application)
  --json                (gives output in JSON)

Global flags
  -c           (configuration file, defaults to 'db2pc-sample.yaml')
  -v           (prints verbosely, useful for debugging)
  -d, --dry-run       (optional) shows what scripts will run, but does not run the scripts
  -n, --no-generate   (optional) do not generate bash scripts as part of install/cleanup, instead use the ones in generated folder. Useful for running local change to the scripts
  -r, --re-run        (optional) re-run script from initial state, ignoring previously saved state

```

## Building the application

Use the following assuming that you have a MacBook. You may to tweak this for Linux and Windows.

1. Clone the repo
2. `brew install go`

Make sure your `PATH` env variable has `go/bin` included (`export PATH=$PATH:$(go env GOPATH)/bin)`)

3. `git checkout develop`
4. `go get github.com/rakyll/statik`
5. `brew install goreleaser`
6. `make install` to test the tool on your MacBook or Windows laptop.
7. `make build-linux` to build the GO binary for Linux.
8. `make send-linux` to scp the tool to your target RHEL/CentOS machine.
9. `upload.sh` - used by the developers of this tool to release the binray to GitHub.  

## Generate yaml configuration file

`db2ctl init`

This will create the `db2ctl-sample.yaml`. It is a configuration file that the tool uses to generate linbit, pacemaker and db2 configuration filess. 

In the spirit of the open-source, please put up a PR if you are making changes to the deployment scripts so that it is useful to the community.

Copy `db2ctl-sample.yaml` to `db2ctl.yaml` and make changes. 

## Validate config file

`db2ctl parse -c db2ctl.yaml`

This will check if the `db2ctl.yaml` is valid or not.
(_`Run with -v` to print the whole parsed object_)

_Note:_ `-c` defaults to `db2ctl.yaml`, so you can omit to specify `-c db2ctl.yaml`, if the file is in the same directory.

## Generate all bash configuration files

`db2ctl generate all -c db2ctl.yaml`

This will generate all the configuration files needed for the application.

_Note:_ `-c` defaults to `db2ctl.yaml`, so it can be ignored if the file is in the same directory.

## Install config files

`db2ctl install <command> -c db2ctl.yaml`

_Note:_ `-c` defaults to `db2ctl.yaml`, so it can be ignored if the file is in the same directory.

## Cleanup

`db2ctl cleanup <command> -c db2ctl.yaml`

_Note:_ `-c` defaults to `db2ctl.yaml`, so it can be ignored if the file is in the same directory.
