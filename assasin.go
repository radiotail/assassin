package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	var flags struct {
		Client    	string
		Server    	string
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
	flag.StringVar(&flags.Server, "s", "", "server listen address or url")
	flag.StringVar(&flags.Client, "c", "", "client connect address or url")
	flag.Parse()

	if !strings.HasPrefix(flags.Addr, "ss://") {
		log.Fatal("error addr " + flags.Addr)
	}

	urlInfo, err := url.Parse(flags.Addr)
	if err != nil {
		log.Fatal(err)
	}

	if flags.Client != "" {
		go local.Start(urlInfo.Host)
	}

	if flags.Server != "" {
		go server.Start(urlInfo.Host)
	}

	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, syscall.SIGINT, syscall.SIGTERM)
	<-quitCH

	log.Println("server quit...")
}
