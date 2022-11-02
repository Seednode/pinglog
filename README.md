# A slightly improved `ping(8)`

![example output](https://git.seedno.de/seednode/pinglog/raw/branch/master/example.png)

## About
Over the years, I've written a number of bash wrappers on top of `ping(8)`, to add things like timestamps, logging, or trying to determine which packets in a sequence were dropped or otherwise lost. 

As part of a recent project, I've been converting my old shell scripts to Go, and this was next in the list.

Built using [go-ping](https://pkg.go.dev/github.com/go-ping/ping).

Builds available [here](https://cdn.seedno.de/builds/pinglog).

## Features
Added features compared to `ping(8)` include:
- Prepending timestamps
- Displaying dropped packets
- Specifying intervals with units (h,m,s,ms,...ns)
- Logging to a file
- Colorized output
- View current statistics with \<Return\>

## Color
For colorized output to work on Windows 10 with Powershell prior to v7.2.2, you need to enable VT support.

This can be done (persistently and globally) in the following ways:
- From Powershell, by running `Set-ItemProperty HKCU:\Console VirtualTerminalLevel -Type DWORD 1`
- From Command Prompt, by running `reg add HKCU\Console /v VirtualTerminalLevel /t REG_DWORD /d 1`

Either method requires a new terminal session to take effect.

Support on older versions is not guaranteed.

Colors can be stripped from log files via the `strip` subcommand, e.g. `pinglog strip file.log`.

## Linux
You may need to run `sudo sysctl -w net.ipv4.ping_group_range="0 2147483647"` on Linux hosts.
(See [here](https://github.com/go-ping/ping#supported-operating-systems) for details)

## Usage output
```
Usage:
  pinglog [flags] <host>
  pinglog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  loss        Calculate periods of packet loss from log file(s)
  strip       Strip ANSI color codes from log file
  version     Print version

Flags:
  -b, --beep                               enable audible bell for exceeded max-rtt
  -C, --color                              enable colorized output
  -c, --count uint                         number of pings to send
  -d, --dropped                            log dropped pings
  -f, --force                              overwrite log file without prompting
  -h, --help                               help for pinglog
  -i, --interval duration                  time between pings (default 1s)
  -4, --ipv4                               force dns resolution to ipv4
  -6, --ipv6                               force dns resolution to ipv6
  -m, --max-rtt duration                   colorize pings over this rtt (default 1h0m0s)
  -o, --output string[="<hostname>.log"]   write to the specified file as well as stdout
  -q, --quiet                              only display summary at end
  -s, --size uint16                        size of payload, in bytes (default 56)
  -w, --timeout duration                   connection timeout (default 2562047h47m16.854775807s)
  -t, --timestamp                          prepend timestamps to output
  -T, --ttl uint16                         maximum time-to-live (default 128)

Use "pinglog [command] --help" for more information about a command.
```
