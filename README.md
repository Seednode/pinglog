# A slightly improved `ping(8)`

![example output](https://git.seedno.de/seednode/pinglog/raw/branch/master/example.png)

## About
Over the years, I've written a number of bash wrappers on top of `ping(8)`, to add things like timestamps, logging, or trying to determine which packets in a sequence were dropped or otherwise lost. 

As part of a recent project, I've been converting my old shell scripts to Go, and this was next in the list.

Built using [go-ping](https://pkg.go.dev/github.com/go-ping/ping).

## Features
Added features compared to `ping(8)` include:
- Prepending timestamps
- Displaying dropped packets
- Specifying intervals in (s|ms|us|ns)
- Logging to a file
- Colorized output

## Color
For colorized output to work on Windows 10 with Powershell prior to v7.2.2, you need to enable VT support.

This can be done (persistently and globally) in the following ways:
- From Powershell, by running `Set-ItemProperty HKCU:\Console VirtualTerminalLevel -Type DWORD 1`
- From Command Prompt, by running `reg add HKCU\Console /v VirtualTerminalLevel /t REG_DWORD /d 1`

Either method requires a new terminal session to take effect.

Support on older versions is not guaranteed.

## Usage output
```
Usage:
  pinglog [flags] <host>
  pinglog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print version

Flags:
  -C, --color                              enable colorized output
  -c, --count int                          number of packets to send (default -1)
  -d, --dropped                            log dropped packets
  -f, --force                              overwrite log file without prompting
  -h, --help                               help for pinglog
  -i, --interval duration                  time between packets (default 1s)
  -n, --no-rtt                             do not record RTTs (reduces memory use for long sessions)
  -o, --output string[="<hostname>.log"]   write to the specified file as well as stdout
  -p, --privileged                         run as privileged user (always enabled on Windows)
  -q, --quiet                              only display summary at end
  -s, --size int                           size of packets, in bytes (default 56)
  -w, --timeout duration                   connection timeout (default 2562047h47m16.854775807s)
  -t, --timestamp                          prepend timestamps to output

Use "pinglog [command] --help" for more information about a command.
```