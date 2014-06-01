package main

import (
	"io"
	"net"
	"fmt"
	"sync"
	"errors"
	"strings"
)

var (
	Proxies ProxyList
)

func init() {
	Proxies.m = make(map[string]Proxy)
}

type ProxyList struct {
	sync.Mutex
	m map[string]Proxy
}

type Proxy struct {
	Spec string // tcp@1.2.3.4:8080->2.3.4.5:9090
	Local, Remote, Proto string
	Listener net.Listener
}

func (p *Proxy) Parse() error {
	parts1 := strings.Split(p.Spec, "@")
	if len(parts1) != 2 {
		parts1 = strings.Split("tcp@" + p.Spec, "@")
	}

	parts2 := strings.Split(parts1[1], "->")
	if len(parts2) != 2 {
		return errors.New("Wrong number of parts")
	}

	if _, _, err := net.SplitHostPort(parts2[0]); err != nil {
		return err
	}
	if _, _, err := net.SplitHostPort(parts2[1]); err != nil {
		return err
	}

	p.Proto = parts1[0]
	p.Local = parts2[0]
	p.Remote = parts2[1]

	return nil
}

func copyAndClose(w io.WriteCloser, r io.Reader) {
	io.Copy(w, r)
	if err := w.Close(); err != nil {
		fmt.Println("Error closing", err)
	}
}

func (p *Proxy) Run() {
	for {
		conn, err := p.Listener.Accept()
		if err != nil {
			if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
				fmt.Println("Temporary error in Accept", err)
				continue
			}
			fmt.Println("Permanent error in Accept", err)
			return
		}

		upstream, err := net.Dial(p.Proto, p.Remote)
		if err != nil {
			err = conn.Close()
			if err != nil {
				fmt.Println("Error closing socket", err)
			}
			fmt.Println("Error contacting remote", err)
			continue
		}

		go copyAndClose(conn, upstream)
		go copyAndClose(upstream, conn)
	}
}

func (p *Proxy) Connect() error {
	err := p.Parse()
	if err != nil {
		fmt.Println("Error parsing", err)
		return err
	}

	p.Listener, err = net.Listen(p.Proto, p.Local)
	if err != nil {
		fmt.Println("Error listening", err)
		return err
	}

	go p.Run()

	return nil
}

func (p *Proxy) Disconnect() error {
	return p.Listener.Close()
}

func CreateProxy(spec string) error {
	Proxies.Lock(); defer Proxies.Unlock()
	fmt.Println("creating", spec)
	p := &Proxy{Spec: spec}

	err := p.Connect()
	if err != nil {
		return err
	}
	Proxies.m[p.Spec] = *p
	return nil
}

func DeleteProxy(spec string) error {
	Proxies.Lock(); defer Proxies.Unlock()
	fmt.Println("deleting", spec)

	p, ok := Proxies.m[spec]
	if !ok {
		return errors.New("No proxy for this spec")
	}
	delete(Proxies.m, spec)	
	return p.Disconnect()
}

func ListProxies() (string, error) {
	Proxies.Lock(); defer Proxies.Unlock()
	var keys []string
	for k := range Proxies.m {
    	keys = append(keys, k)
	}
	return strings.Join(keys, "\n"), nil
}