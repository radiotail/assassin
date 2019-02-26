package server

import (
	"../socks"
	"io"
	"log"
	"net"
	"time"
)

func forward(client, connection net.Conn) (int64, int64, error) {
	type res struct {
		N   int64
		Err error
	}
	ch := make(chan res)

	go func() {
		n, err := io.Copy(client, connection)
		client.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
		connection.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
		ch <- res{n, err}
	} ()

	n, err := io.Copy(connection, client)
	connection.SetDeadline(time.Now()) // wake up the other goroutine blocking on right
	client.SetDeadline(time.Now())  // wake up the other goroutine blocking on left
	rs := <-ch

	if err == nil {
		err = rs.Err
	}
	return n, rs.N, err
}

func newClient(client net.Conn) {
	defer client.Close()

	for {
		buffer := make([]byte, 1024)
		_, err := client.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("client closed!")
				break
			}
			log.Printf("failed to read: %v\n", err)
		}

		if err := socks.Negotiation(client); err != nil {
			log.Printf("socks negotiation fail: %v\n", err)
			return
		}

		targetAddr, err := socks.Request(client)
		if err != nil {
			log.Printf("socks Request fail: %v\n", err)
			return
		}

		connection, err := net.Dial("tcp", targetAddr)
		if err != nil {
			log.Printf("failed to connect to target: %v", err)
			return
		}
		defer connection.Close()

		log.Printf("build proxy from %s to %s through %s\n", client.RemoteAddr(), targetAddr, client.LocalAddr())
		_, err = client.Write(buffer)
		if err != nil {
			log.Printf("failed to send: %v\n", err)
		}

		_, _, err = forward(client, connection)
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				return // ignore i/o timeout
			}
			log.Printf("forward error: %v\n", err)
		}
	}
}

func Start(addr string)  {
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
		return
	}

	log.Printf("sever listening on %s\n", addr)
	for {
		client, err := ls.Accept()
		if nil != err {
			log.Printf("failed to accept: %v\n", err)
			continue
		}

		go newClient(client)
	}
}

