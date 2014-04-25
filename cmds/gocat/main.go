/*
gocat - dynamically configurable tcp proxy

This guys listens on a http port and has the following API:

/ -> list the currently configured tcp proxies
/add?proxy=127.0.0.1:8000->1.2.3.4:80 add a new proxy
/del?proxy=127.0.0.1:8000->1.2.3.4:80 remove existing proxy

The http server url can be assined on command line. HTTP server can be disabled
too by setting -http=none.

Initial set of proxy servers can be passed on command line as args.
*/
package main

import (
	"fmt"
	"flag"
	"github.com/amitu/gutils"
)

func main() {
	hostport := flag.String("http", ":3354", "HTTP Server, can also be 'none'.")
	flag.Parse()

	for _, arg := range flag.Args() {
		if err := CreateProxy(arg); err != nil {
			fmt.Println(err)
		}
	}

	if (*hostport != "none") {
		fmt.Println("listening on", *hostport)
		go HttpServer(*hostport)
	}

	gutils.WaitForCtrlC()
}