# DB2 orchestrator

This will help generate pacemaker and corosync configurations, along with linbit.

<!-- @import "[TOC]" {cmd="toc" depthFrom=2 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Terminology](#terminology)
- [Commands](#commands)
- [Building the application](#building-the-application)
- [Generate yaml configuration file](#generate-yaml-configuration-file)
- [Validate config file](#validate-config-file)
- [Generate all bash configuration files](#generate-all-bash-configuration-files)
- [Install config files](#install-config-files)
- [Cleanup](#cleanup)
- [Helpful resources](#helpful-resources)
  - [Good links for templating](#good-links-for-templating)
  - [Examples for cobra cli](#examples-for-cobra-cli)
  - [Bash scripts through go](#bash-scripts-through-go)
  - [Versioning with go](#versioning-with-go)

<!-- /code_chunk_output -->

## Terminology

There are 2 different configuration files that are mentioned in this repo:

1. The initial `yaml` config file. This is similar to `values.yaml` in helm, and it is used by the orchestrator to create the second set of configuration files.
1. The second set of configuration files are the `bash configuration scripts` that are created by this orchestrator. These files are used for setting up the desired component (it can be linbit, pacemaker, db2, etc.)

In short, it uses a `yaml` config file to create a set of `bash` config files.

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

- status        (displays install status of all components, accepts optional args)
  - install     ((displays install status of all installed components)
  - cleanup     ((displays install status of all cleanup components)

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

1. Clone the repo
1. `git checkout mvp`
1. `make install`
1. Make sure your `PATH` env variable has `go/bin` included (`export PATH=$PATH:$(go env GOPATH)/bin)`)

## Generate yaml configuration file

`db2ctl init`

This will create the `db2pc-sample.yaml`. It is a sample configuration file that the application needs to generate linbit and pacemaker/corosync configurations. You can change it according to your specifications.

## Validate config file

`db2ctl parse -c db2pc-sample.yaml`

This will check if the `db2pc-sample.yaml` is valid or not. It will throw an error if the config file is invalid.
(_`Run with -v` to print the whole parsed object_)

_Note:_ `-c` defaults to `db2pc-sample.yaml`, so it can be ignored if the file is in the same directory.

## Generate all bash configuration files

`db2ctl generate all -c db2pc-sample.yaml`

This will generate all the configuration files needed for the application.

_Note:_ `-c` defaults to `db2pc-sample.yaml`, so it can be ignored if the file is in the same directory.

## Install config files

`db2ctl install <command> -c db2pc-sample.yaml`

_Note:_ `-c` defaults to `db2pc-sample.yaml`, so it can be ignored if the file is in the same directory.

## Cleanup

`db2ctl cleanup <command> -c db2pc-sample.yaml`

_Note:_ `-c` defaults to `db2pc-sample.yaml`, so it can be ignored if the file is in the same directory.
