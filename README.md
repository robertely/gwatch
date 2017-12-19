# gwatch [![Build Status](https://travis-ci.org/robertely/gwatch.svg?branch=master)](https://travis-ci.org/robertely/gwatch)

Attempts to mimic the cli behavior of procps-ng `$ watch` but the output is graphed (technically plotted) instead of printed.
`gwatch` will graph the first number it is able to find and discard every thing else.

![Gif of gwatch demo](https://i.imgur.com/10x23py.gif)

## Installation

##### OSX
```
brew install robertely/gwatch/gwatch
```

##### debian-ish
```
wget https://github.com/robertely/gwatch/releases/download/0.0.3/gwatch_0.0.3_amd64.deb
sudo dpkg -i gwatch_0.0.3_amd64.deb
```

##### cent-ish
```
wget https://github.com/robertely/gwatch/releases/download/0.0.3/gwatch-0.0.3-1.x86_64.rpm
sudo rpm -ivh gwatch-0.0.3-1.x86_64.rpm
```

##### from source
```
git clone git@github.com:robertely/gwatch.git
cd gwatch
go get ./...
go build
```

## Usage
```
graphing watch: execute a program periodically, graphing the output fullscreen


Usage: gwatch [-behtvx] [-n value] [parameters ...]
 -b, --beep      beep if command has a non-zero exit
 -e, --errexit   exit if command has a non-zero exit
 -h, --help      display this help and exit
 -n, --interval=value
                 seconds to wait between updates
 -t, --no-title  turn off header
 -v, --version   output version information and exit
 -x, --exec      pass command to exec instead of "sh -c"

```
## Examples
###### OSX Load:
`$ gwatch 'sysctl -n vm.loadavg | cut -d" " -f2'`

###### OSX Pages free:
`$ gwatch "vm_stat | grep 'Pages free'"`

###### Used memory percent:
`$ gwatch -n .1 "free | grep Mem | awk '{print \$3/\$2 * 100.0}'"`

###### Total time for get request via curl:
`$ gwatch "curl -w "%{time_total}" -o /dev/null -s 'https://google.com'"`
