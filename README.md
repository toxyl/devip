# devip
Makes it easy to add and remove IPs for local development. It does so by wrapping the `ip` and `ping` commands on Linux. When run without arguments it lists all IPs assigned to the local machine.

A simple use-case could be to locally test systems with their "real" IP. 

# Usage
First of all, build and install it:
```
sudo CGO_ENABLED=0 go build -o  /usr/local/bin/devip  .
```

## Add IPs
```
devip add [IP 1] <IP 2> .. <IP n>
```
For example:
```
devip add 10.10.10.1 20.20.20.1 30.30.30.1
```
This will make `localhost` serve the IPs `10.10.10.1`, `20.20.20.1` and `30.30.30.1` by creating loopback aliases for them, so you can use them for development. Be aware that this means you won't be able to connect to the real IPs while the aliases are up.

## Remove IPs
```
devip remove [IP 1] <IP 2> .. <IP n>
```
For example:
```
devip remove 10.10.10.1 20.20.20.1 30.30.30.1
```
This will stop `localhost` serving the IPs `10.10.10.1`, `20.20.20.1` and `30.30.30.1` by removing the loopback aliases. 

## List IPs
```
devip 
```
This will print all IPs currently served by `localhost`. 

