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
	quite *bool
	statsonly *bool
)

func main() {
	listen = flag.String("listen", "127.0.0.1:4443", "Listen on this address.")
	quite = flag.Bool("quite", false, "Quite mode.")
	statsonly = flag.Bool(
		"statsonly", false, "Only prints stats on stdout and discard data.",
	)
	flag.Parse()

	conn, err := net.ListenPacket("udp4", *listen)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !*quite {
		if *statsonly {
			fmt.Printf(
				"UDP: Server started on %s in statsonly mode.\n", *listen,
			)
		} else {
			fmt.Printf("UDP: Server started on %s.\n", *listen)
		}
	}

	obytes := make([]byte, 64 * 1024)
	start := time.Now()
	var bcount, count time.Duration

	for {
		bytes := obytes
		n, _, err := conn.ReadFrom(bytes)

		if err != nil {
			fmt.Println(err)
			continue
		}

		if *statsonly {
			count += 1
			bcount += time.Duration(n)

			now := time.Now()
			diff := now.Sub(start)

			if diff > 1e9 {
				fmt.Printf(
					"bps = %sps, pps = %d.\n",
					gutils.FormatBytes(float64(bcount)), count,
				)
				start = now
				bcount, count = 0, 0
			}

		} else {
			bytes = bytes[:n]
			fmt.Printf(string(bytes))
		}
	}
}