package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

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

		log.Printf("read byte: %s\n", buffer)
		_, err = client.Write(buffer)
		if err != nil {
			log.Printf("failed to send: %v\n", err)
		}
	}
}

func server(addr string)  {
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

func main() {
	var flags struct {
		Addr		string
		Cipher		string
		Key       	string
		Password  	string
		Socks     	string
		TCPTun    	string
	}

	log.Println("server start...")

	flag.StringVar(&flags.Addr, "addr", "", "server listen address or url")
	flag.StringVar(&flags.Key, "key", "", "base64url-encoded key (derive from password if empty)")
	flag.StringVar(&flags.Password, "password", "", "password")
	flag.StringVar(&flags.Socks, "socks", "", "(client-only) SOCKS listen address")
	flag.StringVar(&flags.TCPTun, "tcptun", "", "(client-only) TCP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.Parse()

	if !strings.HasPrefix(flags.Addr, "ss://") {
		log.Fatal("error addr " + flags.Addr)
	}

	urlInfo, err := url.Parse(flags.Addr)
	if err != nil {
		log.Fatal(err)
	}

	go server(urlInfo.Host)

	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, syscall.SIGINT, syscall.SIGTERM)
	<-quitCH

	log.Println("server quit...")
}
