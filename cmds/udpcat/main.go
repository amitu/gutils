package main 

import (
	"flag"
	"fmt"
	"net"
)

var (
	listen *string
	quite *bool
)

func main() {
	listen = flag.String("listen", "127.0.0.1:4443", "Listen on this address.")
	quite = flag.Bool("quite", false, "Quite mode.")
	flag.Parse()

	conn, err := net.ListenPacket("udp4", *listen)

	if err != nil {
		fmt.Println(err)
		return
	}

	if !*quite {
		fmt.Printf("UDP: Server started on %s.\n", *listen)
	}


	for {
		bytes := make([]byte, 512)
		n, _, err := conn.ReadFrom(bytes)

		if err != nil {
			fmt.Println(err)
			continue
		}

		bytes = bytes[:n]
		fmt.Printf(string(bytes))
	}
}