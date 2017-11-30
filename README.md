# gwatch

Attempts to mimic the cli behavior of procps-ng `$ watch` but output is graphed (technically plotted) instead of printed.
`gwatch` will graph the first number it is able to find and discard every thing else.

```
graphing watch: expects numerical values, graphs the first one it sees.

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
## Recipes:
###### OSX Load:
`$ gwatch 'sysctl -n vm.loadavg | cut -d" " -f2'`

###### OSX Pages free:
`$ gwatch "vm_stat | grep 'Pages free'"`

###### Used memory percent:
`$ gwatch -n .1 "free | grep Mem | awk '{print \$3/\$2 * 100.0}'"`

###### Total time for get request via curl:
`$ gwatch "curl -w "%{time_total}" -o /dev/null -s 'https://google.com'"`


### TODO
[x] Handle basic errors that like to crash

[ ] reasonable error handling

[ ] print a description with Usage

[ ] document...

[ ] literally any testing

[ ] write man Pages
