package main

import (
	"fmt"
	"net/http"
)

func AddHandler(res http.ResponseWriter, req *http.Request) {
	if err := CreateProxy(req.FormValue("proxy")); err != nil {
		fmt.Fprintf(res, "%s", err)
		return
	}

	fmt.Fprintf(res, "ok")
}

func DelHandler(res http.ResponseWriter, req *http.Request) {
	if err := DeleteProxy(req.FormValue("proxy")); err != nil {
		fmt.Fprintf(res, "%s", err)
		return
	}

	fmt.Fprintf(res, "ok")
}

func ListHandler(res http.ResponseWriter, req *http.Request) {
	proxies, err := ListProxies()

	if err != nil {
		fmt.Fprintf(res, "%s", err)
		return
	}

	fmt.Fprintf(res, proxies)
}

func HealthHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "ok")
}

func HttpServer(hostport string) {
	http.HandleFunc("/", ListHandler)	
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/del", DelHandler)
	http.HandleFunc("/health", HealthHandler)

	http.ListenAndServe(hostport, nil)
}

