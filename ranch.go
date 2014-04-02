package gutils

import (
	"io"
	"net"
	"fmt"
	"time"
	"sync"
	"bufio"
	"errors"
	"net/http"
)

type RanchServer struct {
	Network string
	HostPort string
	Concurrancy int
	HTTP bool
	Handler func(net.Conn)
	HTTPHandler func(http.Request, net.Conn)
	listener net.Listener
	wg sync.WaitGroup
}

func (s *RanchServer) Start () error {
	fmt.Println("start", s)

	if s.listener != nil {
		return errors.New("Server already started.")
	}

	if s.Network == "" {
		s.Network = "tcp"
	}

	if s.Concurrancy == 0 {
		s.Concurrancy = 10
	}

	var err error

	fmt.Println("start", s)
	s.listener, err = net.Listen(s.Network, s.HostPort)

	if err != nil {
		return err
	}

	for i := 0; i < s.Concurrancy; i++ {
		go s.worker()
	}

	return nil
}

func (s *RanchServer) Stop () error {
	if s.listener == nil {
		return errors.New("Server not running.")
	}

	err := s.listener.Close()

	if err != nil {
		return err
	}

	s.wg.Wait()

	return nil
}

const noLimit int64 = (1 << 63) - 1

func (s *RanchServer) worker() {
	fmt.Println("worker started")
	s.wg.Add(1)
	var tempDelay time.Duration
	reader := new(bufio.Reader)
	var lr *io.LimitedReader

	for {
		fmt.Println("worker waiting")
		conn, e := s.listener.Accept()
		fmt.Println("conn", conn)
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				fmt.Printf("ranch: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			fmt.Printf("ranch: error during accept,", e)
			// log the error and shut down the listener?
			s.wg.Done()
		}
		tempDelay = 0

		if s.HTTP {
			reader.Reset(conn)
			lr = io.LimitReader(reader, noLimit).(*io.LimitedReader)
			br := bufio.NewReader(lr)
			bw := bufio.NewWriterSize(conn, 4<<10)
			buf := bufio.NewReadWriter(br, bw)
			req, err := http.ReadRequest(buf.Reader)
			if err != nil {
				fmt.Println(err)
				conn.Close()
				continue
			}
			s.HTTPHandler(*req, conn)
		} else {
			s.Handler(conn)
		}
	}
}