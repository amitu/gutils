package main

import (
	"flag"
	"fmt"
	"net"
	"time"
	"github.com/amitu/gutils"
)

var (
	listen *string
	redis *string
	drops *string
	key string
)

func UDPReader(conn net.PacketConn) {
	obytes := make([]byte, 64 * 1024)
	start := time.Now()
	var bcount, count time.Duration

	fmt.Printf("UDP:%s -> REDIS:%s@%s.\n", *listen, key, *redis)

	for {
		bytes := obytes
		n, _, err := conn.ReadFrom(bytes)

		if err != nil {
			fmt.Println(err)
			continue
		}

		count += 1
		bcount += time.Duration(n)

		now := time.Now()
		diff := now.Sub(start)

		if diff > 1e9 {
			fmt.Printf(
				"bps = %s, pps = %d.\n", gutils.FormatBytes(float64(bcount)),
				 count,
			)
			start = now
			bcount, count = 0, 0
		}
	}
}

func RedisWriter(conn net.Conn) {

}

func main() {
	listen = flag.String("listen", "127.0.0.1:4443", "Listen on this address.")
	redis = flag.String("redis", "127.0.0.1:6379", "Redis server host:port.")
	drops = flag.String(
		"drops", "/dev/stderr",
		"File to write dropped packets in when redis is down.",
	)
	flag.Parse()

	key = flag.Arg(0)

	if key == "" {
		fmt.Println("Please pass key as argument.")
		return
	}

	conn, err := net.ListenPacket("udp4", *listen)

	if err != nil {
		fmt.Println(err)
		return
	}

	rconn, err := net.Dial("tcp4", *redis)
	rconn = rconn
	if err != nil {
		fmt.Println(err)
		return
	}

	go RedisWriter(rconn)
	go UDPReader(conn)

	gutils.WaitForCtrlC()
}