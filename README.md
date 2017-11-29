# gwatch
graphing watch



## recipes
* `./gwatch 'sysctl -n vm.loadavg | cut -d" " -f2'`
* `./gwatch "vm_stat | grep 'Pages free'"`
* `./gwatch "curl -w "%{time_total}" -o /dev/null -s 'http://google.com'"`
