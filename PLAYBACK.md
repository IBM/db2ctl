# Using go-lang to create your own CLI + templating engine for generating any kinds of scripts

When I say scripts, I mean scripts that install software, upgrade software, basically they do anything essential on a server.

The scripts can be bash, python, it doesn't matter.

## Old way - create and run all scripts manually

_Script1_

```
ssh root@10.20.30.40
cd /root/test/script1

do some things...
```

_Script2_

```
ssh root@10.20.30.40
cd /root/test/script2

do some other things...
```

## New way

### 1. Have a CLI - very easy to do

For example:

```
$ my_sample_application

This is a sample to generate configuration scripts

Usage:
  my_sample_application [command]

Available Commands:
  cleanup     Cleanup commands
  generate    Generates editable files and configurations needed for application
  help        Help about any command
  init        Generate the initial config yaml file
  install     Install commands
  parse       Parse and validate configuration from yaml file

Flags:
  -c, --conf string   configuration yaml file needed for application
  -h, --help          help for my_sample_application
  -v, --verbose       print verbosely

Use "my_sample_application [command] --help" for more information about a command.

```

### 2. Have a templating engine - similar to HELM - moderate difficulty

For example: Instead of having a script like this:

```
ssh root@10.20.30.40
cd /root/test

do some other things...
```

You have a **template** like this:

```
ssh root@{{IP_ADDRESS}}
cd {{FOLDER_LOCATION}}

do some other things...
```

And have **a common configuration yaml file**:

_config.yaml_

```
IP_ADDRESS: 10.20.30.40
FOLDER_LOCATION: /root/test
```

You see that the script has the hardcoded IP address of the server. Now if you want to run this script against another server, you will manually need to change the IP address inside the script, and then run it against the new server.

This seems trivial, but what if you had many places where the IP address was used? You would have to do a complex search and replace, possibly using `sed`, which is hard to use.

## Technologies used

- **Cobra** - Golang's most-used CLI creation library. https://github.com/spf13/cobra.
  [Cobra generator](https://github.com/spf13/cobra/blob/master/cobra/README.md#cobra-generator) is extremely powerful.
- **text/template** - Golang's templating library - https://golang.org/pkg/text/template/
