# gutils: Utilities for go

[![Join the chat at https://gitter.im/amitu/gutils](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/amitu/gutils?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## udpcat

This utility listens on a UDP host:port, and prints whatever it receives 
on stdout.

### To install

```shell
go get github.com/amitu/gutils/cmds/udpcat
```

### Usage

```shell
$ udpcat -help 
Usage of udpcat:
  -listen="127.0.0.1:4443": Listen on this address.
  -quite=false: Quite mode.
  -statsonly=false: Only prints stats on stdout and discard data.

$ udpcat -listen 127.0.0.1:3334
UDP: Server started on 127.0.0.1:3334.
hello there

# from a different shell

$ echo hello there | nc -4 -u localhost 3334
```

## udpflood

`udpflood` writes to UDP destination as fast as it can. The speed to write 
can be specified on command lines.

`udpflood` writes summary about data transmission rate on stdout every second.

### To install

```shell
go get github.com/amitu/gutils/cmds/udpflood
```

### Usage

```shell
$ udpflood -help 
Usage of udpflood:
  -server="127.0.0.1:4443": Send to this address.
  -file="somefile.txt": Send the content of this file to server.
```

## udp2redis

`udp2redis` listens on udp, writes every packet to to redis in a list.

### To install

```shell
go get github.com/amitu/gutils/cmds/stdin2redis
```

### Usage

```shell
$ stdin2redis -help 
Usage of udpflood:
  -redis="127.0.0.1:6379": Redis server host:port.
  -server="127.0.0.1:4443": UDP host:port to listen on.
  -drops="/dev/stderr": File to write dropped packets in when redis is down.
First argument is name of redis list to write packets in.
```

At startup redis server must be connectable. If redis goes down while the
service is running, it will drop UDP packats, and start pushing it to
redis once it comes back live.




