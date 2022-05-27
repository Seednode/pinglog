## A slightly improved ping(8)

[![example output](https://git.seedno.de/seednode/pinglog/raw/branch/master/example.png)]

### About
Over the years, I've written a number of bash wrappers on top of ping(8), to add things like timestamps, logging, or trying to determine which packets in a sequence were dropped or otherwise lost. 

As part of a recent project, I've been converting my old shell scripts to Go, and this was next in the list.

### Features
Added features compared to ping(8) include:
- Prepending timestamps
- Displaying dropped packets
- Specifying intervals in (s|ms|us|ns)
- Logging to a file
- Colorized output

### Usage
```
Usage:
  pinglog [flags] <host>
  pinglog [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Print version

Flags:
  -c, --count int                          number of packets to send (default -1)
  -d, --dropped                            log dropped packets
  -f, --force                              overwrite log file without prompting
  -h, --help                               help for pinglog
  -i, --interval duration                  time between packets (default 1s)
  -x, --no-color                           disable colorized output
  -n, --no-rtt                             do not record RTTs (reduces memory use for long sessions)
  -o, --output string[="<hostname>.log"]   write to the specified file as well as stdout
  -p, --privileged                         run as privileged user (needed on Windows)
  -q, --quiet                              only display summary at end
  -s, --size int                           size of packets, in bytes (default 56)
  -w, --timeout duration                   connection timeout (default 15m0s)
  -t, --timestamp                          prepend timestamps to output

Use "pinglog [command] --help" for more information about a command.
```